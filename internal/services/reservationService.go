package services

import (
	"github.com/brothify/internal/repositories"
)

type ReservationService struct {
	repo *repositories.ReservationRepository
}

func NewReservationService(	repo *repositories.ReservationRepository) *ReservationService {
	return &ReservationService{repo: repo}
}