package service

import (
	"x-ui/backend-api/internal/repository"
	"x-ui/database/model"
)

type InboundClientService struct {
	InboundClientRepo repository.InboundClient
}

func NewInboundClientService(inboundClientRepo repository.InboundClient) *InboundClientService {
	return &InboundClientService{
		InboundClientRepo: inboundClientRepo,
	}
}

func (s *InboundClientService) GetInboundClients(inboundId int) ([]model.Client, error) {
	return s.InboundClientRepo.Get(inboundId)
}

func (s *InboundClientService) AddInboundClient(inboundId int, newClient *model.Client) (bool, error) {
	return s.InboundClientRepo.Add(inboundId, newClient)
}
