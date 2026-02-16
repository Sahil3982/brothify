package dto

import (
	"time"

	"github.com/google/uuid"
)

type CategoryResponse struct {
	ID          uuid.UUID `json:"category_id"`
	NAME        string    `json:"name"`
	SLUG        string    `json:"slug"`
	DESCRIPTION string    `json:"description"`
}

type DishResponse struct {
	ID           uuid.UUID          `json:"dish_id"`
	NAME         string             `json:"dish_name"`
	PRICE        float64            `json:"price"`
	DESCRIPTION  string             `json:"description"`
	DISH_URL     string             `json:"dish_url"`
	AVAILABILITY bool               `json:"availability"`
	RATING       float64            `json:"rating"`
	HIGHLIGHT    bool               `json:"highlight"`
	UPDATEDAT    time.Time          `json:"updated_at"`
	CREATEDAT    time.Time          `json:"created_at"`
	CATEGORIES   []CategoryResponse `json:"categories"`
}
