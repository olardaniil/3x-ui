package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"x-ui/backend-api/internal/repository"
	"x-ui/database/model"
)

type InboundClientService struct {
	InboundClientRepo repository.InboundClient
	InboundRepo       repository.Inbound
}

func NewInboundClientService(inboundClientRepo repository.InboundClient, inboundRepo repository.Inbound) *InboundClientService {
	return &InboundClientService{
		InboundClientRepo: inboundClientRepo,
		InboundRepo:       inboundRepo,
	}
}

func (s *InboundClientService) GetInboundClients(inboundId int) ([]model.Client, error) {
	return s.InboundClientRepo.Get(inboundId)
}

func (s *InboundClientService) AddInboundClient(inboundId int, newClient *model.Client) (string, bool, error) {
	res, err := s.InboundClientRepo.Add(inboundId, newClient)
	if err != nil {
		return "", false, err
	}
	// Генерируем ключ
	inbound, err := s.InboundRepo.GetById(inboundId)
	if err != nil {
		return "", false, err
	}
	streamSettingsStr := inbound.StreamSettings
	// Преобразоване строки в структуру
	streamSettings := StreamSettingsModel{}
	err = json.Unmarshal([]byte(streamSettingsStr), &streamSettings)
	if err != nil {
		return "", false, err
	}

	ip := GetOutboundIP()

	//  Ключ
	key := fmt.Sprint(string(inbound.Protocol) + "://" + newClient.ID + "@" + ip.String() + ":" + strconv.Itoa(inbound.Port) + "?type=" + streamSettings.Network + "&security=" + streamSettings.Security + "&pbk=" + streamSettings.RealitySettings.Settings.PublicKey + "&fp=chrome&sni=" + streamSettings.RealitySettings.ServerNames[0] + "&sid=" + streamSettings.RealitySettings.ShortIds[0] + "&spx=%2F&flow=" + newClient.Flow + "#" + inbound.Remark + "-" + newClient.Email)

	return key, res, nil
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

type StreamSettingsModel struct {
	Network         string `json:"network"`
	Security        string `json:"security"`
	ExternalProxy   []any  `json:"externalProxy"`
	RealitySettings struct {
		Show        bool     `json:"show"`
		Xver        int      `json:"xver"`
		Dest        string   `json:"dest"`
		ServerNames []string `json:"serverNames"`
		PrivateKey  string   `json:"privateKey"`
		MinClient   string   `json:"minClient"`
		MaxClient   string   `json:"maxClient"`
		MaxTimediff int      `json:"maxTimediff"`
		ShortIds    []string `json:"shortIds"`
		Settings    struct {
			PublicKey   string `json:"publicKey"`
			Fingerprint string `json:"fingerprint"`
			ServerName  string `json:"serverName"`
			SpiderX     string `json:"spiderX"`
		} `json:"settings"`
	} `json:"realitySettings"`
	TCPSettings struct {
		AcceptProxyProtocol bool `json:"acceptProxyProtocol"`
		Header              struct {
			Type string `json:"type"`
		} `json:"header"`
	} `json:"tcpSettings"`
}
