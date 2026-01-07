package services

import (
	"github.com/brothify/internal/models"
	"github.com/brothify/internal/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.repo.GetAllCategories()
}
func (s *CategoryService) CreateCategory(d *models.Category) (*models.Category, error) {
	return s.repo.CreateCategory(d)
}
func (s *CategoryService) UpdateCategory(d *models.Category, params string) (*models.Category, error) {
	return s.repo.UpdateCategory(d, params)
}
func (s *CategoryService) DeleteCategory(id string) error {
	return s.repo.DeleteCategory(id)
}