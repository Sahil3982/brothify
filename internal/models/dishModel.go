package models

import (
	"time"

	"github.com/google/uuid"
)

type Dish struct {
	ID           uuid.UUID `json:"dish_id"`
	NAME         string    `json:"dish_name"`
	CATID        *int      `json:"cat_id"`
	PRICE        float64   `json:"price"`
	DESCRIPTION  string    `json:"description"`
	DISHURL      *string   `json:"dish_url"`
	AVAILABILITY bool      `json:"availability"`
	RATING       float64   `json:"rating"`
	HIGHLIGHT    bool      `json:"highlight"`
	UPDATEDAT    time.Time `json:"updated_at"`
	CREATEDAT    time.Time `json:"created_at"`
}
