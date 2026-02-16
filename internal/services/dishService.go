package services

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/brothify/internal/config"
	"github.com/brothify/internal/dto"
	"github.com/brothify/internal/models"
	"github.com/brothify/internal/repositories"
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

func (s *DishService) CreateDish(ctx context.Context, req *dto.DishRequest, file multipart.File, fileHeader *multipart.FileHeader) (*dto.DishResponse, error) {

	imageURL, err := config.UploadImageToS3(file, fileHeader, os.Getenv("AWS_S3_BUCKET"))
	if err != nil {
		return nil, err
	}

	dishModel := &models.Dish{
		NAME:         req.NAME,
		DESCRIPTION:  req.DESCRIPTION,
		PRICE:        req.PRICE,
		DISHURL:      imageURL,
		AVAILABILITY: true,
		RATING:       0,
		HIGHLIGHT:    false,
	}

	newDish, err := s.repo.CreateDish(ctx, dishModel, req.CATEGORYID)
	if err != nil {
		return nil, err
	}

	response := &dto.DishResponse{
		ID:           newDish.ID,
		NAME:         newDish.NAME,
		PRICE:        float64(newDish.PRICE), 
		DESCRIPTION:  newDish.DESCRIPTION,
		DISH_URL:     newDish.DISHURL,
		AVAILABILITY: newDish.AVAILABILITY,
		RATING:       newDish.RATING,
		HIGHLIGHT:    newDish.HIGHLIGHT,
		CREATEDAT:    newDish.CREATEDAT,
		UPDATEDAT:    newDish.UPDATEDAT,
	}

	for _, cat := range newDish.CATEGORIES {
		response.CATEGORIES = append(response.CATEGORIES, dto.CategoryResponse{
			ID:          cat.ID,
			NAME:        cat.NAME,
			DESCRIPTION: cat.DESCRIPTION,
		})
	}

	return response, nil
}

func (s *DishService) UpdateDish(id string, d *models.Dish) error {
	return s.repo.UpdateDish(id, d)
}

func (s *DishService) DeleteDish(id string) error {
	return s.repo.DeleteDish(id)
}
