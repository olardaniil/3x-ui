package service

import (
	"x-ui/backend-api/internal/repository"
)

type UserService struct {
	UserRepo repository.User
}

func NewUserService(userRepo repository.User) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) CheckUser(username string, password string) bool {
	user, err := s.UserRepo.CheckUser(username, password)
	if err != nil {
		return false
	}
	return user.Username != "" && user.Password != ""
}
