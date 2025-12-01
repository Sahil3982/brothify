package services

import (

	"github.com/brothify/internal/models"
	"github.com/brothify/internal/repositories"
	
)

type ReservationService struct {
	repo *repositories.ReservationRepository
}

func NewReservationService(repo *repositories.ReservationRepository) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) GetAllReservations(search, status, date string, limit, offset int) ([] models.Reservation, error){
	return s.repo.GetAllReservations(search, status, date, limit, offset)
}

func (s *ReservationService) CreateReservation(d *models.Reservation) (*models.Reservation, error) {
	
	return s.repo.CreateReservation(d)
}

