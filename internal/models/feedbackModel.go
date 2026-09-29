package models

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	ID uuid.UUID `json:"id"`
	Comment string `json:"comment"`
	Rate float32 `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}