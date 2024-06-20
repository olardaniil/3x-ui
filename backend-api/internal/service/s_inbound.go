package service

import (
	"x-ui/backend-api/internal/repository"
	"x-ui/database/model"
)

type InboundService struct {
	InboundRepo repository.Inbound
}

func NewInboundService(inboundRepo repository.Inbound) *InboundService {
	return &InboundService{
		InboundRepo: inboundRepo,
	}
}

func (s *InboundService) GetInbounds() ([]model.Inbound, error) {
	inbounds, err := s.InboundRepo.Get()
	if err != nil {
		return nil, err
	}
	return inbounds, nil
}

func (s *InboundService) GetInbound(inboundId int) (model.Inbound, error) {
	inbound, err := s.InboundRepo.GetById(inboundId)
	if err != nil {
		return inbound, err
	}
	return inbound, nil
}
