package services

import (
	"github.com/brothify/internal/models"
	"github.com/brothify/internal/repositories"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) LoginUser(d *models.User) (*models.User, error) {
	return s.repo.LoginUser(d)
}
