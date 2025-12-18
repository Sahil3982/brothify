package services

import (
	"github.com/brothify/internal/models"
	"github.com/brothify/internal/repositories"
	"github.com/google/uuid"
)

type ReservationService struct {
	repo *repositories.ReservationRepository
}

func NewReservationService(repo *repositories.ReservationRepository) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) GetReservationByID(id uuid.UUID) (*models.Reservation, error) {
	return s.repo.GetReservationByID(id)
}

func (s *ReservationService) GetAllReservations(search, status, date string, limit, offset int) ([]models.Reservation, error) {
	return s.repo.GetAllReservations(search, status, date, limit, offset)
}

func (s *ReservationService) CreateReservation(d *models.Reservation) (*models.Reservation, error) {
	return s.repo.CreateReservation(d)
}

func (s *ReservationService) UpdateReservation(d *models.Reservation, params string) (*models.Reservation, error) {
	return s.repo.UpdateReservation(d, params)
}

func (s *ReservationService) DeleteReservation(d *models.Reservation, params string) error {
	return s.repo.DeleteReservation(d, params)
}

func (s *ReservationService) GetDishPrice(dishID uuid.UUID) (float64, error) {
	return s.repo.GetDishPriceByID(dishID)
}
