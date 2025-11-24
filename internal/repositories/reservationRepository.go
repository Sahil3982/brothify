package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type ReservationRepository struct {
	DB *pgxpool.Pool
}
func NewReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{DB: db}
}

