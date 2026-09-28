package services

import (
	"log"

	"github.com/brothify/internal/models"
	"github.com/brothify/internal/repositories"
	"github.com/google/uuid"
)

type EmailService interface {
	SendReservationStatusEmail(reservation *models.Reservation) error
}

type ReservationService struct {
	repo  *repositories.ReservationRepository
	email EmailService
}

func NewReservationService(repo *repositories.ReservationRepository, email EmailService) *ReservationService {
	return &ReservationService{
		repo:  repo,
		email: email,
	}
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

func (s *ReservationService) UpdateReservation(d *models.Reservation, params uuid.UUID) (*models.Reservation, error) {
	oldReservation, err := s.repo.GetReservationByID(params)
	if err != nil {
		return nil, err
	}
	oldStatus := oldReservation.STATUS
	if d.USERID == nil {
		d.USERID = oldReservation.USERID
	}
	if d.TABLENUMBER == 0 {
		d.TABLENUMBER = oldReservation.TABLENUMBER
	}
	if d.RESERVATIONPERSONNAME == "" {
		d.RESERVATIONPERSONNAME = oldReservation.RESERVATIONPERSONNAME
	}
	if d.RESERVATIONPERSONEMAIL == "" {
		d.RESERVATIONPERSONEMAIL = oldReservation.RESERVATIONPERSONEMAIL
	}
	if d.RESERVATIONPERSONMOBILENUMBER == "" {
		d.RESERVATIONPERSONMOBILENUMBER = oldReservation.RESERVATIONPERSONMOBILENUMBER
	}
	if d.NUMBEROFGUESTS == 0 {
		d.NUMBEROFGUESTS = oldReservation.NUMBEROFGUESTS
	}
	if d.RESERVATIONTIME == "" {
		d.RESERVATIONTIME = oldReservation.RESERVATIONTIME
	}
	if d.RESERVATIONDATE == "" {
		d.RESERVATIONDATE = oldReservation.RESERVATIONDATE
	}
	if d.SPECIALREQUESTS == nil {
		d.SPECIALREQUESTS = oldReservation.SPECIALREQUESTS
	}
	if d.STATUS == "" {
		d.STATUS = oldReservation.STATUS
	}
	updatedReservation, err := s.repo.UpdateReservation(d, params)
	if err != nil {
		return nil, err
	}
	if oldStatus != updatedReservation.STATUS {
		log.Printf("Reservation %s status changed from %q to %q; sending email to %s", params, oldStatus, updatedReservation.STATUS, updatedReservation.RESERVATIONPERSONEMAIL)
		if err := s.email.SendReservationStatusEmail(updatedReservation); err != nil {
			return nil, err
		}
		log.Printf("Reservation status email sent for %s", params)
	}

	return updatedReservation, nil
}

func (s *ReservationService) DeleteReservation(d *models.Reservation, params string) error {
	return s.repo.DeleteReservation(d, params)
}

func (s *ReservationService) GetDishPrice(dishID uuid.UUID) (float64, error) {
	return s.repo.GetDishPriceByID(dishID)
}
