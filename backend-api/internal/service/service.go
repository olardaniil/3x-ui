package service

import (
	"x-ui/backend-api/internal/repository"
	"x-ui/database/model"
)

type Inbound interface {
	GetInbounds() ([]model.Inbound, error)
	GetInbound(inboundId int) (model.Inbound, error)
}

type InboundClient interface {
	GetInboundClients(inboundId int) ([]model.Client, error)
	AddInboundClient(inboundId int, newClient *model.Client) (string, bool, error)
}

type User interface {
	CheckUser(username string, password string) bool
}

type Service struct {
	Inbound
	InboundClient
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Inbound:       NewInboundService(repos.Inbound),
		InboundClient: NewInboundClientService(repos.InboundClient),
		User:          NewUserService(repos.User),
	}
}
