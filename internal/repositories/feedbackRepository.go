package repositories

import "github.com/jackc/pgx/v5/pgxpool"

type FeedbackRepository struct {
	DB *pgxpool.Pool
}

func NewFeedbackRepository(db *pgxpool.Pool) *FeedbackRepository {
	return &FeedbackRepository{
		DB: db,
	}
}

