package repository

import (
	"gorm.io/gorm"
	"x-ui/database/model"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) CheckUser(username string, password string) (model.User, error) {
	var user model.User
	err := r.db.Model(model.User{}).Where("username = ? AND password = ?", username, password).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, nil
}
