package repository

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"x-ui/database"
	"x-ui/database/model"
	"x-ui/logger"
	"x-ui/util/common"
	"x-ui/web/service"
	"x-ui/xray"
)

type InboundClientsRepo struct {
	db          *gorm.DB
	xrayApi     xray.XrayAPI
	InboundRepo InboundRepo
}

func NewInboundClientsRepo(db *gorm.DB) *InboundClientsRepo {
	return &InboundClientsRepo{
		db: db,
	}
}

func (r *InboundClientsRepo) Get(inboundId int) ([]model.Client, error) {
	var inbound model.Inbound
	err := r.db.Model(model.Inbound{}).Preload("ClientStats").First(&inbound, inboundId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	settings := map[string][]model.Client{}
	err = json.Unmarshal([]byte(inbound.Settings), &settings)
	if settings == nil {
		return nil, fmt.Errorf("setting is null")
	}

	clients := settings["clients"]
	if clients == nil {
		return nil, nil
	}
	return clients, nil
}

func (r *InboundClientsRepo) Add(inboundId int, newClient *model.Client) (bool, error) {
	clients, err := r.Get(inboundId)
	if err != nil {
		return false, err
	}
	log.Println("step 2")
	interfaceClients := []interface{}{newClient}
	existEmail, err := r.checkEmailsExistForClients([]model.Client{*newClient})
	if err != nil {
		return false, err
	}
	if existEmail != "" {
		log.Println(existEmail)
		return false, common.NewError("Duplicate email:", existEmail)
	}
	log.Println("step 3")
	var oldInbound model.Inbound
	err = r.db.Model(model.Inbound{}).Preload("ClientStats").First(&oldInbound, inboundId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	// Secure client ID
	for _, client := range clients {
		if oldInbound.Protocol == "trojan" {
			if client.Password == "" {
				return false, common.NewError("empty client ID")
			}
		} else if oldInbound.Protocol == "shadowsocks" {
			if client.Email == "" {
				return false, common.NewError("empty client ID")
			}
		} else {
			if client.ID == "" {
				return false, common.NewError("empty client ID")
			}
		}
	}

	var oldSettings map[string]interface{}
	err = json.Unmarshal([]byte(oldInbound.Settings), &oldSettings)
	if err != nil {
		return false, err
	}

	oldClients := oldSettings["clients"].([]interface{})
	oldClients = append(oldClients, interfaceClients...)

	oldSettings["clients"] = oldClients

	newSettings, err := json.MarshalIndent(oldSettings, "", "  ")
	if err != nil {
		return false, err
	}

	oldInbound.Settings = string(newSettings)

	db := database.GetDB()
	tx := db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	needRestart := false

	r.xrayApi.Init(service.P.GetAPIPort())
	for _, client := range clients {
		if len(client.Email) > 0 {
			r.AddClientStat(tx, inboundId, &client)
			if client.Enable {
				cipher := ""
				if oldInbound.Protocol == "shadowsocks" {
					cipher = oldSettings["method"].(string)
				}
				err1 := r.xrayApi.AddUser(string(oldInbound.Protocol), oldInbound.Tag, map[string]interface{}{
					"email":    client.Email,
					"id":       client.ID,
					"flow":     client.Flow,
					"password": client.Password,
					"cipher":   cipher,
				})
				if err1 == nil {
					logger.Debug("Client added by api:", client.Email)
				} else {
					logger.Debug("Error in adding client by api:", err1)
					needRestart = true
				}
			}
		} else {
			needRestart = true
		}
	}
	r.xrayApi.Close()

	return needRestart, tx.Save(oldInbound).Error
}

func (r *InboundClientsRepo) AddClientStat(tx *gorm.DB, inboundId int, client *model.Client) error {
	clientTraffic := xray.ClientTraffic{}
	clientTraffic.InboundId = inboundId
	clientTraffic.Email = client.Email
	clientTraffic.Total = client.TotalGB
	clientTraffic.ExpiryTime = client.ExpiryTime
	clientTraffic.Enable = true
	clientTraffic.Up = 0
	clientTraffic.Down = 0
	clientTraffic.Reset = client.Reset
	result := tx.Create(&clientTraffic)
	err := result.Error
	return err
}

func (r *InboundClientsRepo) checkEmailsExistForClients(clients []model.Client) (string, error) {
	allEmails, err := r.getAllEmails()
	if err != nil {
		return "", err
	}
	var emails []string
	for _, client := range clients {
		if client.Email != "" {
			if contains(emails, client.Email) {
				return client.Email, nil
			}
			if contains(allEmails, client.Email) {
				return client.Email, nil
			}
			emails = append(emails, client.Email)
		}
	}
	return "", nil
}

func (r *InboundClientsRepo) getAllEmails() ([]string, error) {
	db := database.GetDB()
	var emails []string
	err := db.Raw(`
		SELECT JSON_EXTRACT(client.value, '$.email')
		FROM inbounds,
			JSON_EACH(JSON_EXTRACT(inbounds.settings, '$.clients')) AS client
		`).Scan(&emails).Error
	if err != nil {
		return nil, err
	}
	return emails, nil
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
