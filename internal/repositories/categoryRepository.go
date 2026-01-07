package repositories

import (
	"context"

	"github.com/brothify/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	DB *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

func (r *CategoryRepository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	ctx := context.Background()

	query := `SELECT cat_id, name, description, slug FROM categories`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.NAME, &category.DESCRIPTION, &category.SLUG); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *CategoryRepository) CreateCategory(c *models.Category) (*models.Category, error) {
	ctx := context.Background()
	query := `INSERT INTO categories (name, description, slug) VALUES ($1, $2, $3) RETURNING cat_id, name, description, slug, created_at, updated_at`
	err := r.DB.QueryRow(ctx, query, c.NAME, c.DESCRIPTION, c.SLUG).Scan(&c.ID, &c.NAME, &c.DESCRIPTION, &c.SLUG, &c.CREATEDAT, &c.UPDATEDAT)	
	if err != nil {
		return nil, err
	}
	return c, nil
}	

func (r *CategoryRepository) UpdateCategory(c *models.Category, params string) (*models.Category, error) {
	ctx := context.Background()
	query := `UPDATE categories SET name = $1, description = $2, slug = $3 WHERE cat_id = $4 RETURNING cat_id, name, description, slug`
	err := r.DB.QueryRow(ctx, query, c.NAME, c.DESCRIPTION, c.SLUG, params).Scan(&c.ID, &c.NAME, &c.DESCRIPTION, &c.SLUG)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CategoryRepository) DeleteCategory(id string) error {
	ctx := context.Background()
	query := `DELETE FROM categories WHERE cat_id = $1`
	_, err := r.DB.Exec(ctx, query, id)
	return err
}