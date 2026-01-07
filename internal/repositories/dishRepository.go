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

func (r *DishRepository) GetDishByID(id int) (*models.Dish, error) {
	ctx := context.Background()
	query := `SELECT dish_id, dish_name, cat_id, price, description, dish_url, availability, rating, highlight, created_at, updated_at 
			  FROM dishes WHERE dish_id = $1`
	var d models.Dish
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.NAME, &d.CATID, &d.PRICE, &d.DESCRIPTION,
		&d.DISHURL, &d.AVAILABILITY, &d.RATING, &d.HIGHLIGHT,
		&d.CREATEDAT, &d.UPDATEDAT,
	)
	if err != nil {
		log.Println("❌ GetDishByID error:", err)
		return nil, err
	}	
	return &d, nil
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

func (r *DishRepository) CreateDish(d *models.Dish) (*models.Dish, error) {
	ctx := context.Background()
	query := `
		INSERT INTO dishes (dish_name, description, price, cat_id, dish_url, availability, rating, highlight)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING dish_id, dish_name, cat_id, price, description, dish_url, availability, rating, highlight, created_at, updated_at
	`
	row := r.DB.QueryRow(ctx, query, d.NAME, d.DESCRIPTION, d.PRICE, d.CATID, d.DISHURL, d.AVAILABILITY, d.RATING, d.HIGHLIGHT	)

	var newDish models.Dish
	err := row.Scan(
		&newDish.ID, &newDish.NAME, &newDish.CATID, &newDish.PRICE, &newDish.DESCRIPTION,
		&newDish.DISHURL, &newDish.AVAILABILITY, &newDish.RATING, &newDish.HIGHLIGHT,
		&newDish.CREATEDAT, &newDish.UPDATEDAT,
	)
	if err != nil {
		log.Println("❌ Insert Scan error:", err)
		return nil, err
	}
	return &newDish, nil
}

func (r *DishRepository) UpdateDish(id string, d *models.Dish) error {
	ctx := context.Background()
	query := `
		UPDATE dishes 
		SET dish_name = $1, description = $2, price = $3, cat_id = $4, dish_url = $5, availability = $6, rating = $7, highlight = $8
		WHERE dish_id = $9
	`
	_, err := r.DB.Exec(ctx, query, d.NAME, d.DESCRIPTION, d.PRICE, d.CATID, d.DISHURL, d.AVAILABILITY, d.RATING, d.HIGHLIGHT, id)
	return err
}

func (r *DishRepository) DeleteDish(id string) error {
	ctx := context.Background()
	_, err := r.DB.Exec(ctx, "DELETE FROM dishes WHERE dish_id = $1", id)
	return err
}
