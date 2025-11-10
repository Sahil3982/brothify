package repositories

import (
	"context"
	"log"

	"github.com/brothify/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DishRepository struct {
	DB *pgxpool.Pool
}

func NewDishRepository(db *pgxpool.Pool) *DishRepository {
	return &DishRepository{DB: db}
}

func (r *DishRepository) GetAllDishes() ([]models.Dish, error) {
	ctx := context.Background()
	rows, err := r.DB.Query(ctx, `
		SELECT dish_id, dish_name, cat_id, price, description, dish_url, availability, rating, highlight, created_at, updated_at 
		FROM dishes
	`)
	if err != nil {
		log.Println("❌ Query error:", err)
		return nil, err
	}
	defer rows.Close()

	var dishes []models.Dish

	for rows.Next() {
		var d models.Dish
		// Use *int for CATID to handle NULL
		err := rows.Scan(
			&d.ID, &d.NAME, &d.CATID, &d.PRICE, &d.DESCRIPTION,
			&d.DISHURL, &d.AVAILABILITY, &d.RATING, &d.HIGHLIGHT,
			&d.CREATEDAT, &d.UPDATEDAT,
		)
		if err != nil {
			log.Println("⚠️ Scan error:", err)
			continue
		}
		dishes = append(dishes, d)
	}

	if err = rows.Err(); err != nil {
		log.Println("❌ Rows iteration error:", err)
		return nil, err
	}

	log.Println("✅ Dishes retrieved:", len(dishes))
	return dishes, nil
}

func (r *DishRepository) CreateDish(d *models.Dish) error {
	ctx := context.Background()
	_, err := r.DB.Exec(ctx,
		`INSERT INTO dishes (dish_name, description, price) VALUES ($1, $2, $3)`,
		d.NAME, d.DESCRIPTION, d.PRICE)
	return err
}
