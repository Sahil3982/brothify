package models

import "github.com/google/uuid"

type Category struct {
	ID          uuid.UUID `json:"cat_id"`
	NAME        string    `json:"name"`
	DESCRIPTION string    `json:"description"`
	SLUG        string    `json:"slug"`
}
