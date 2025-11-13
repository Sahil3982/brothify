package services

import (
	"github.com/brothify/internal/models"
	"github.com/brothify/internal/repositories"
)

type DishService struct {
	repo *repositories.DishRepository
}

func NewDishService(repo *repositories.DishRepository) *DishService {
	return &DishService{repo: repo}
}

func (s *DishService) GetAllDishes() ([]models.Dish, error) {
	return s.repo.GetAllDishes()
}

func (s *DishService) CreateDish(d *models.Dish) (*models.Dish, error) {
	return s.repo.CreateDish(d)
}


func (s *DishService) UpdateDish(id string, d *models.Dish) error {
	return s.repo.UpdateDish(id, d)
}

func (s *DishService) DeleteDish(id string) error {
	return s.repo.DeleteDish(id)
}