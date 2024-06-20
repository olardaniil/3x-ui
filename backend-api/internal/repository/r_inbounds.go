package repository

import (
	"gorm.io/gorm"
	"x-ui/database/model"
)

type InboundRepo struct {
	db *gorm.DB
}

func NewInboundRepo(db *gorm.DB) *InboundRepo {
	return &InboundRepo{
		db: db,
	}
}

func (r *InboundRepo) Get() ([]model.Inbound, error) {
	var inbounds []model.Inbound
	err := r.db.Model(model.Inbound{}).Preload("ClientStats").Find(&inbounds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return inbounds, nil
}

func (r *InboundRepo) GetById(inboundId int) (model.Inbound, error) {
	var inbound model.Inbound
	err := r.db.Model(model.Inbound{}).Preload("ClientStats").First(&inbound, inboundId).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return inbound, err
	}
	return inbound, nil
}
