package repositories

import (
	"context"

	"github.com/brothify/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepository struct {
	DB *pgxpool.Pool
}
func NewReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{DB: db}
}

func (r *ReservationRepository) CreateReservation(d *models.Reservation) (*models.Reservation, error) {

	query := `INSERT INTO reservations (reservation_person_name, reservation_person_email, number_of_guests, reservation_time) VALUES ($1, $2, $3, $4) RETURNING reservation_id, reservation_person_name, reservation_person_email, number_of_guests, reservation_time`
	row := r.DB.QueryRow(context.Background(),query)
	err := row.Scan(&d.ID, &d.RESERVATIONPERSONNAME, &d.RESERVATIONPERSONEMAIL, &d.NUMBEROFGUESTS, &d.RESERVATIONTIME)
	if err != nil {
		return nil, err
	}
	return d, nil

}
