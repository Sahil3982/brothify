package repositories

import (
	"context"
	"log"

	"github.com/brothify/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepository struct {
	DB *pgxpool.Pool
}

func NewReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{DB: db}
}

func (r *ReservationRepository) GetReservationByID(id uuid.UUID) (*models.Reservation, error) {

	query := `
        SELECT 
            reservation_id,
            user_id,
            table_number,
            reservation_person_name,
            reservation_person_email,
            reservation_person_mobile_number,
            number_of_guests,
            reservation_time,
            reservation_date,
            special_requests,
            status,
            created_at,
            updated_at
        FROM reservations 
        WHERE reservation_id = $1
    `

	var res models.Reservation

	err := r.DB.QueryRow(context.Background(), query, id).Scan(
		&res.ID,
		&res.USERID,
		&res.TABLENUMBER,
		&res.RESERVATIONPERSONNAME,
		&res.RESERVATIONPERSONEMAIL,
		&res.RESERVATIONPERSONMOBILENUMBER,
		&res.NUMBEROFGUESTS,
		&res.RESERVATIONTIME,
		&res.RESERVATIONDATE,
		&res.SPECIALREQUESTS,
		&res.STATUS,
		&res.CREATEDAT,
		&res.UPDATEDAT,
	)

	if err != nil {
		log.Println("❌ GetReservationByID error:", err)
		return nil, err
	}

	dishQuery := `SELECT d.dish_id, d.dish_name, d.cat_id, d.price, d.description, d.dish_url, d.availability, d.rating, d.highlight 
              FROM reservation_dishes rd 
              JOIN dishes d ON rd.dish_id = d.dish_id 
              WHERE rd.reservation_id = $1`
	dishRows, _ := r.DB.Query(context.Background(), dishQuery, res.ID)

	var dishes []models.Dish
	for dishRows.Next() {
		var d models.Dish
		dishRows.Scan(
			&d.ID,
			&d.NAME,
			&d.CATID,
			&d.PRICE,
			&d.DESCRIPTION,
			&d.DISHURL,
			&d.AVAILABILITY,
			&d.RATING,
			&d.HIGHLIGHT,
		)
		dishes = append(dishes, d)
	}
	dishRows.Close()

	res.DISHDETAILS = dishes

	return &res, nil
}

func (r *ReservationRepository) GetAllReservations(search string, status string, date string, limit int, offset int) ([]models.Reservation, error) {
	query := `SELECT 
	reservation_id,
	user_id,
	table_number,
	reservation_person_name,
	reservation_person_email,
	reservation_person_mobile_number,
	number_of_guests,
	reservation_time,
	reservation_date,
	special_requests,
	status,
	created_at,
	updated_at FROM reservations
	  WHERE 
            ($1 = '' OR 
                reservation_person_name ILIKE '%' || $1 || '%' OR
                reservation_person_email ILIKE '%' || $1 || '%' OR
                reservation_person_mobile_number ILIKE '%' || $1 || '%')
        AND ($2 = '' OR LOWER(status) = LOWER($2))
        AND ($3 = '' OR reservation_date = $3)
        ORDER BY created_at DESC
        LIMIT $4 OFFSET $5
	`
	rows, err := r.DB.Query(context.Background(), query, search, status, date, limit, offset)
	if err != nil {
		log.Println("❌ Query error:", err)
		return nil, err
	}
	defer rows.Close()
	var reservations []models.Reservation

	for rows.Next() {
		var res models.Reservation
		err := rows.Scan(
			&res.ID,
			&res.USERID,
			&res.TABLENUMBER,
			&res.RESERVATIONPERSONNAME,
			&res.RESERVATIONPERSONEMAIL,
			&res.RESERVATIONPERSONMOBILENUMBER, // 5 mobile number
			&res.NUMBEROFGUESTS,                // 6 guests
			&res.RESERVATIONTIME,               // 7 reservation_time
			&res.RESERVATIONDATE,               // 8 reservation_date
			&res.SPECIALREQUESTS,               // 9 special_requests
			&res.STATUS,                        // 10 status
			&res.CREATEDAT,                     // 11 created_at
			&res.UPDATEDAT,                     // 12 updated_at
		)

		if err != nil {
			return nil, err
		}
		dishQuery := `SELECT d.dish_id, d.dish_name, d.cat_id, d.price, d.description, d.dish_url, d.availability, d.rating, d.highlight FROM reservation_dishes rd JOIN dishes d ON rd.dish_id = d.dish_id WHERE rd.reservation_id = $1`
		dishRows, err := r.DB.Query(context.Background(), dishQuery, res.ID)
		if err != nil {
			return nil, err
		}
		var dishes []models.Dish
		for dishRows.Next() {
			var d models.Dish
			err := dishRows.Scan(
				&d.ID,
				&d.NAME,
				&d.CATID,
				&d.PRICE,
				&d.DESCRIPTION,
				&d.DISHURL,
				&d.AVAILABILITY,
				&d.RATING,
				&d.HIGHLIGHT,
			)
			if err != nil {
				return nil, err
			}

			dishes = append(dishes, d)

		}
		dishRows.Close()
		res.DISHDETAILS = dishes
		reservations = append(reservations, res)

	}
	return reservations, nil

}

func (r *ReservationRepository) CreateReservation(d *models.Reservation) (*models.Reservation, error) {
	query := `INSERT INTO reservations
 (user_id,
  table_number,
  reservation_person_name,
  reservation_person_email,
  reservation_person_mobile_number,
  number_of_guests,
  reservation_time,
  reservation_date,
  special_requests,
  status)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
RETURNING 
    reservation_id,
    user_id,
    table_number,
    reservation_person_name,
    reservation_person_email,
    reservation_person_mobile_number,
    number_of_guests,
    reservation_time,
    reservation_date,
    special_requests,
    status,
    created_at,
    updated_at`

	row := r.DB.QueryRow(context.Background(), query,
		d.USERID,
		d.TABLENUMBER,
		d.RESERVATIONPERSONNAME,
		d.RESERVATIONPERSONEMAIL,
		d.RESERVATIONPERSONMOBILENUMBER,
		d.NUMBEROFGUESTS,
		d.RESERVATIONTIME,
		d.RESERVATIONDATE,
		d.SPECIALREQUESTS,
		d.STATUS,
	)

	err := row.Scan(
		&d.ID,
		&d.USERID,
		&d.TABLENUMBER,
		&d.RESERVATIONPERSONNAME,
		&d.RESERVATIONPERSONEMAIL,
		&d.RESERVATIONPERSONMOBILENUMBER,
		&d.NUMBEROFGUESTS,
		&d.RESERVATIONTIME,
		&d.RESERVATIONDATE,
		&d.SPECIALREQUESTS,
		&d.STATUS,
		&d.CREATEDAT,
		&d.UPDATEDAT,
	)

	if err != nil {
		return nil, err
	}

	// insert dish items
	for _, dishId := range d.DISHITEMS {
		_, err = r.DB.Exec(context.Background(),
			`INSERT INTO reservation_dishes (reservation_id, dish_id) VALUES ($1, $2)`,
			d.ID, dishId,
		)
		if err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (r *ReservationRepository) UpdateReservation(d *models.Reservation, id string) (*models.Reservation, error) {
	query := `
        UPDATE reservations SET
            user_id                        = $1,
            table_number                   = $2,
            reservation_person_name        = $3,
            reservation_person_email       = $4,
            reservation_person_mobile_number = $5,
            number_of_guests               = $6,
            reservation_time               = $7,
            reservation_date               = $8,
            special_requests               = $9,
            status                         = $10,
            updated_at                     = NOW()
        WHERE reservation_id = $11
        RETURNING
            reservation_id,
            user_id,
            table_number,
            reservation_person_name,
            reservation_person_email,
            reservation_person_mobile_number,
            number_of_guests,
            reservation_time,
            reservation_date,
            special_requests,
            status,
            created_at,
            updated_at;
    `

	var res models.Reservation

	err := r.DB.QueryRow(
		context.Background(),
		query,
		d.USERID,
		d.TABLENUMBER,
		d.RESERVATIONPERSONNAME,
		d.RESERVATIONPERSONEMAIL,
		d.RESERVATIONPERSONMOBILENUMBER,
		d.NUMBEROFGUESTS,
		d.RESERVATIONTIME,
		d.RESERVATIONDATE,
		d.SPECIALREQUESTS,
		d.STATUS,
		id,
	).Scan(
		&res.ID,
		&res.USERID,
		&res.TABLENUMBER,
		&res.RESERVATIONPERSONNAME,
		&res.RESERVATIONPERSONEMAIL,
		&res.RESERVATIONPERSONMOBILENUMBER,
		&res.NUMBEROFGUESTS,
		&res.RESERVATIONTIME,
		&res.RESERVATIONDATE,
		&res.SPECIALREQUESTS,
		&res.STATUS,
		&res.CREATEDAT,
		&res.UPDATEDAT,
	)

	if err != nil {
		log.Println("❌ UpdateReservation error:", err)
		return nil, err
	}

	return &res, nil
}

func (r *ReservationRepository) DeleteReservation(d *models.Reservation, id string) error {
	query := `DELETE FROM reservations WHERE reservation_id = $1`
	_, err := r.DB.Exec(context.Background(), query, id)
	return err
}

func (r *ReservationRepository) GetDishPriceByID(id uuid.UUID) (float64, error) {
	var price float64
	err := r.DB.QueryRow(context.Background(),
		"SELECT price FROM dishes WHERE dish_id = $1", id,
	).Scan(&price)

	if err != nil {
		return 0, err
	}
	return price, nil
}

func (r *ReservationRepository) SaveInvoiceURL(reservationID uuid.UUID, paymentID, signature, invoiceURL string) error {
	_, err := r.DB.Exec(context.Background(),
		`UPDATE reservations 
         SET payment_id = $1,
             signature = $2,
             payment_status = 'PAID',
             invoice_url = $3,
             updated_at = NOW() 
         WHERE reservation_id = $4`,
		paymentID, signature, invoiceURL, reservationID,
	)
	return err
}
