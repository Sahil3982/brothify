package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/brothify/internal/models"
	"github.com/brothify/internal/repositories"
	"github.com/google/uuid"
)

type EmailService interface {
	SendReservationStatusEmail(reservation *models.Reservation) error
	SendFeedbackReminderEmail(reservation *models.Reservation) error
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

func (s *ReservationService) Name() string {
	return "feedback-reminders"
}

func (s *ReservationService) Spec() string {
	return "@every 1m"
}

func (s *ReservationService) Run(ctx context.Context) error {
	return s.SendDueFeedbackReminders(ctx)
}

func (s *ReservationService) SendDueFeedbackReminders(ctx context.Context) error {
	location := time.Local
	if timezone := strings.TrimSpace(os.Getenv("RESERVATION_TIMEZONE")); timezone != "" {
		var err error
		location, err = time.LoadLocation(timezone)
		if err != nil {
			return fmt.Errorf("load reservation timezone %q: %w", timezone, err)
		}
	}

	reservations, err := s.repo.GetPendingFeedbackReminders(ctx)
	if err != nil {
		return fmt.Errorf("get pending feedback reminders: %w", err)
	}

	now := time.Now().In(location)
	var reminderErrors []error
	for i := range reservations {
		reservation := &reservations[i]
		reservationAt, err := parseReservationDateTime(
			reservation.RESERVATIONDATE,
			reservation.RESERVATIONTIME,
			location,
		)
		if err != nil {
			reminderErrors = append(reminderErrors, fmt.Errorf("reservation %s: %w", reservation.ID, err))
			continue
		}
		if !feedbackReminderDue(reservationAt, now) {
			continue
		}
		if err := s.email.SendFeedbackReminderEmail(reservation); err != nil {
			reminderErrors = append(reminderErrors, fmt.Errorf("send feedback reminder for %s: %w", reservation.ID, err))
			continue
		}
		if err := s.repo.MarkFeedbackReminderSent(ctx, reservation.ID); err != nil {
			reminderErrors = append(reminderErrors, fmt.Errorf("mark feedback reminder sent for %s: %w", reservation.ID, err))
		}
	}
	return errors.Join(reminderErrors...)
}

func feedbackReminderDue(reservationAt, now time.Time) bool {
	return !now.Before(reservationAt.Add(6 * time.Hour))
}

func parseReservationDateTime(date, clock string, location *time.Location) (time.Time, error) {
	parsedDate, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(date), location)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid reservation date %q: %w", date, err)
	}

	clock = strings.TrimSpace(clock)
	var parsedTime time.Time
	for _, layout := range []string{"15:04", "15:04:05", "3:04 PM", "3:04:05 PM"} {
		parsedTime, err = time.ParseInLocation(layout, clock, location)
		if err == nil {
			return time.Date(
				parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
				parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, location,
			), nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid reservation time %q", clock)
}
