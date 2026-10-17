package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/brothify/internal/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FeedbackRepository struct {
	DB *pgxpool.Pool
}

func NewFeedbackRepository(db *pgxpool.Pool) *FeedbackRepository {
	return &FeedbackRepository{
		DB: db,
	}
}

func (r *FeedbackRepository) CreateFeedback(feedback dto.FeedbackRequest) error {
	_, err := r.DB.Exec(context.Background(),
		`INSERT INTO feedback (comment, rating) VALUES ($1, $2)`,
		feedback.Comment, feedback.Rate,
	)
	return err
}

func (r *FeedbackRepository) UpdateFeedback(feedback dto.FeedbackRequest, id uuid.UUID) error {
	_, err := r.DB.Exec(context.Background(),
		`UPDATE feedback SET comment = $1, rating = $2 WHERE feedback_id = $3`,
		feedback.Comment, feedback.Rate, id,
	)
	return err
}

func (r *FeedbackRepository) GetFeedback() (string, error) {
	rows, err := r.DB.Query(context.Background(), `SELECT comment, rating FROM feedback`)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var feedback strings.Builder
	for rows.Next() {
		var comment string
		var rating int
		if err := rows.Scan(&comment, &rating); err != nil {
			return "", err
		}
		fmt.Fprintf(&feedback, "Comment: %s, Rating: %d\n", comment, rating)
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return feedback.String(), nil
}

func (r *FeedbackRepository) GetFeedbackByID(id uuid.UUID) (string, error) {
	row := r.DB.QueryRow(context.Background(),
		`SELECT comment, rating FROM feedback WHERE feedback_id = $1`, id,
	)

	var comment string
	var rate int

	err := row.Scan(&comment, &rate)
	if err != nil {
		return "", err
	}

	return "Comment: " + comment + ", Rating: " + fmt.Sprintf("%d", rate), nil
}

func (r *FeedbackRepository) DeleteFeedbackByID(id uuid.UUID) error {
	_, err := r.DB.Exec(context.Background(),
		`DELETE FROM feedback WHERE feedback_id = $1`, id,
	)
	return err
}
