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

	query := `INSERT INTO reservations
	 (user_id,table_number,reservation_person_name, reservation_person_email,reservation_person_mobile_number, number_of_guests, reservation_time,special_requests,status,note) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING reservation_id, user_id, table_number, reservation_person_name,
              reservation_person_email, number_of_guests, reservation_time, reservation_person_mobile_number,
              special_requests, status, note, created_at`
	row := r.DB.QueryRow(context.Background(), query,
		d.USERID,
		d.TABLENUMBER,
		d.RESERVATIONPERSONNAME,
		d.RESERVATIONPERSONEMAIL,
		d.RESERVATIONPERSONMOBILENUMBER,
		d.NUMBEROFGUESTS,
		d.RESERVATIONTIME,
		d.SPECIALREQUESTS,
		d.STATUS,
	)

	err := row.Scan(&d.ID,
		&d.USERID,
		&d.TABLENUMBER,
		&d.RESERVATIONPERSONNAME,
		&d.RESERVATIONPERSONEMAIL,
		&d.RESERVATIONPERSONMOBILENUMBER,
		&d.NUMBEROFGUESTS,
		&d.RESERVATIONTIME,
		&d.SPECIALREQUESTS,
		&d.STATUS,
		&d.CREATEDAT)
	if err != nil {
		return nil, err
	}

	d.DISHITEMS = []int{}

	for _, dishId := range d.DISHITEMS {
		_, err = r.DB.Exec(context.Background(),
			`INSERT INTO reservation_dishes (reservation_id, dish_id) VALUES ($1, $2)`,
			d.ID, dishId,
		)
		if err != nil {
			return nil, err
		}

		d.DISHITEMS = append(d.DISHITEMS, dishId)
	}

	return d, nil

}
