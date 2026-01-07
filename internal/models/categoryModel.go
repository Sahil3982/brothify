package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `json:"cat_id"`
	NAME        string    `json:"name"`
	DESCRIPTION string    `json:"description"`
	SLUG        string    `json:"slug"`
	UPDATEDAT   time.Time `json:"updated_at"`
	CREATEDAT   time.Time `json:"created_at"`
}
