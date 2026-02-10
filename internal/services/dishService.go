package services

import (
	"github.com/brothify/internal/models"
	"github.com/brothify/internal/repositories"
	"github.com/google/uuid"
)

type DishService struct {
	repo *repositories.DishRepository
}

func NewDishService(repo *repositories.DishRepository) *DishService {
	return &DishService{repo: repo}
}

func (s *DishService) GetDishByID(id int) (*models.Dish, error) {
	return s.repo.GetDishByID(id)
}

func (s *DishService) GetAllDishes() ([]models.Dish, error) {
	return s.repo.GetAllDishes()
}

func (s *DishService) CreateDish(d *models.Dish, catID uuid.UUID) (*models.Dish, error) {
	return s.repo.CreateDish(d, catID)
}


func (s *DishService) UpdateDish(id string, d *models.Dish) error {
	return s.repo.UpdateDish(id, d)
}

func (s *DishService) DeleteDish(id string) error {
	return s.repo.DeleteDish(id)
}