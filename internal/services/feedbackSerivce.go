package services

import "github.com/brothify/internal/repositories"

type FeedbackService struct {
	repo *repositories.FeedbackRepository
}

func NewFeedbackService(repo *repositories.FeedbackRepository ) *FeedbackService {
	return &FeedbackService{
		repo: repo,
	}
}