package dto

import (
	"github.com/google/uuid"
)

type DishRequest struct {
	NAME         string    `json:"dish_name"`
	PRICE        float64   `json:"price"`
	DESCRIPTION  string    `json:"description"`
	DISHURL      string   `json:"dish_url"`
	AVAILABILITY bool      `json:"availability"`
	RATING       float64   `json:"rating"`
	HIGHLIGHT    bool      `json:"highlight"`
	CATEGORYID   uuid.UUID `json:"category_id"`
}
