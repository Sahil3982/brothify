package services

import (
	"github.com/brothify/internal/dto"
	"github.com/brothify/internal/repositories"
	"github.com/google/uuid"
)

type FeedbackService struct {
	repo *repositories.FeedbackRepository
}

func NewFeedbackService(repo *repositories.FeedbackRepository) *FeedbackService {
	return &FeedbackService{
		repo: repo,
	}
}

func (s *FeedbackService) CreateFeedback(req dto.FeedbackRequest) error {
	return s.repo.CreateFeedback(req)
}

func (s *FeedbackService) UpdateFeedback(req dto.FeedbackRequest, id uuid.UUID) error {
	return s.repo.UpdateFeedback(req, id)
}

func (s *FeedbackService) GetFeedback() (string, error) {
	return s.repo.GetFeedback()
}

func (s *FeedbackService) GetFeedbackByID(id uuid.UUID) (string, error) {
	return s.repo.GetFeedbackByID(id)
}

func (s *FeedbackService) DeleteFeedbackByID(id uuid.UUID) error {
	return s.repo.DeleteFeedbackByID(id)
}
