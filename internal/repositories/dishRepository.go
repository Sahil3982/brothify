package repositories

import (
	"github.com/brothify/internal/models"
	"github.com/jackc/pgx/v5"
	"context"
)

type DishRepository struct {
	DB *pgx.Conn
}

func NewDishRepository(db *pgx.Conn) *DishRepository {
	return &DishRepository{DB: db}
}

func (r *DishRepository) GetAllDishes() ([]models.Dish, error) {
	rows, err := r.DB.Query(context.Background(), "SELECT * FROM dishs")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dishs []models.Dish

	for rows.Next() {
		var d models.Dish
		if err := rows.Scan(&d.ID, &d.NAME, &d.DESCRIPTION, &d.PRICE, &d.CATID, &d.RATING, &d.HIGHLIGHT, &d.DISHURL, &d.AVAILABILITY, &d.CREATEDAT, &d.UPDATEDAT); err != nil {
			return nil, err
		}
		dishs = append(dishs, d)
	}

	return dishs, nil
}
