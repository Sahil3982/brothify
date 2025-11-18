package repositories

import (
	"context"

	"github.com/brothify/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) LoginUser(d *models.User) (*models.User, error) {
    ctx := context.Background()

    query := `
        SELECT id, name, email, password
        FROM users
        WHERE email=$1 AND password=$2
    `

    row := r.DB.QueryRow(ctx, query, d.EMAIL, d.PASSWORD)

    var newUser models.User

    err := row.Scan(
        &newUser.ID,
        &newUser.NAME,
        &newUser.EMAIL,
        &newUser.PASSWORD,
    )
    if err != nil {
        return nil, err
    }

    return &newUser, nil
}
