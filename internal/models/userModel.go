package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"user_id"`
	NAME      string    `json:"name"`
	EMAIL     string    `json:"email"`
	PASSWORD  string    `json:"password"`
	CREATEDAT time.Time `json:"created_at"`
	UPDATEDAT time.Time `json:"updated_at"`
}
