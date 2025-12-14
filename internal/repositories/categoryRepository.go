package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type CategoryRepository struct {
	DB *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{DB: db}
}
