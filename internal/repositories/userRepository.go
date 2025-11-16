package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type UserRepository struct {
	DB *pgxpool.Pool
}
