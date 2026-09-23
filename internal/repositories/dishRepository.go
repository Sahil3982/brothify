package repositories

import (
	"context"
	"log"

	"github.com/brothify/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DishRepository struct {
	DB *pgxpool.Pool
}

func NewDishRepository(db *pgxpool.Pool) *DishRepository {
	return &DishRepository{DB: db}
}

func (r *DishRepository) GetDishByID(id uuid.UUID) (*models.Dish, error) {
	ctx := context.Background()
	query := `SELECT dish_id, dish_name, category_id, price, description, dish_url, availability, rating, highlight, created_at, updated_at 
			  FROM dishes WHERE dish_id = $1`
	var d models.Dish
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&d.ID, &d.NAME, &d.CATEGORYID, &d.PRICE, &d.DESCRIPTION,
		&d.DISHURL, &d.AVAILABILITY, &d.RATING, &d.HIGHLIGHT,
		&d.CREATEDAT, &d.UPDATEDAT,
	)
	if err != nil {
		log.Println("GetDishByID error:", err)
		return nil, err
	}
	return &d, nil
}

func (r *DishRepository) GetAllDishes() ([]models.Dish, error) {
	ctx := context.Background()
	rows, err := r.DB.Query(ctx, `
		SELECT dish_id, dish_name, category_id, price, description, dish_url, availability, rating, highlight, created_at, updated_at 
		FROM dishes
	`)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}
	defer rows.Close()

	var dishes []models.Dish

	for rows.Next() {
		var d models.Dish
		err := rows.Scan(
			&d.ID, &d.NAME, &d.CATEGORYID, &d.PRICE, &d.DESCRIPTION,
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
		log.Println("Rows iteration error:", err)
		return nil, err
	}

	return dishes, nil
}

func (r *DishRepository) CreateDish(ctx context.Context, d *models.Dish) (*models.Dish, error) {
	log.Printf("Creating dish: %+v\n", d)
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		log.Println("Transaction begin error:", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	dishQuery := `
		INSERT INTO dishes (
		dish_name, category_id, description, price, dish_url,
		 availability, rating, highlight
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING dish_id, dish_name, category_id, price, description, 
		dish_url, availability, rating, highlight, created_at, updated_at
	`
	var newDish models.Dish
	err = tx.QueryRow(
		ctx,
		dishQuery,
		d.NAME,
		d.CATEGORYID,
		d.DESCRIPTION,
		d.PRICE,
		d.DISHURL,
		d.AVAILABILITY,
		d.RATING,
		d.HIGHLIGHT,
	).Scan(
		&newDish.ID,
		&newDish.NAME,
		&newDish.CATEGORYID,
		&newDish.PRICE,
		&newDish.DESCRIPTION,
		&newDish.DISHURL,
		&newDish.AVAILABILITY,
		&newDish.RATING,
		&newDish.HIGHLIGHT,
		&newDish.CREATEDAT,
		&newDish.UPDATEDAT,
	)
	if err != nil {
		log.Println("Insert Scan error:", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Transaction commit error:", err)
		return nil, err
	}

	return &newDish, nil
}

func (r *DishRepository) UpdateDish(id uuid.UUID, d *models.Dish) (bool, error) {
	ctx := context.Background()
	query := `
		UPDATE dishes 
		SET dish_name = $1, description = $2, price = $3, category_id = $4, dish_url = $5, availability = $6, rating = $7, highlight = $8
		WHERE dish_id = $9
	`
	result, err := r.DB.Exec(ctx, query, d.NAME, d.DESCRIPTION, d.PRICE, d.CATEGORYID, d.DISHURL, d.AVAILABILITY, d.RATING, d.HIGHLIGHT, id)
	if err != nil {
		return false, err
	}
	return result.RowsAffected() > 0, nil
}

func (r *DishRepository) DeleteDish(id uuid.UUID) (bool, error) {
	ctx := context.Background()
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "DELETE FROM reservation_dishes WHERE dish_id = $1", id); err != nil {
		return false, err
	}

	result, err := tx.Exec(ctx, "DELETE FROM dishes WHERE dish_id = $1", id)
	if err != nil {
		return false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}
