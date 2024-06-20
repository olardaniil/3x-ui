package repository

import (
	"gorm.io/gorm"
	"x-ui/database/model"
	"x-ui/xray"
)

type Inbound interface {
	Get() ([]model.Inbound, error)
	GetById(inboundId int) (model.Inbound, error)
}

type InboundClient interface {
	Get(inboundId int) ([]model.Client, error)
	Add(inboundId int, newClient *model.Client) (bool, error)
}

type User interface {
	CheckUser(username string, password string) (model.User, error)
}

type Repository struct {
	Inbound
	InboundClient
	User
	xrayApi xray.XrayAPI
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Inbound:       NewInboundRepo(db),
		InboundClient: NewInboundClientsRepo(db),
		User:          NewUserRepo(db),
	}
}
