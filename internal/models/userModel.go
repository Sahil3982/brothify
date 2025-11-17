package models

import "time"

type User struct {
	ID        int       `json:"user_id"`
	NAME      string    `json:"name"`
	EMAIL     string    `json:"email"`
	PASSWORD  string    `json:"password"`
	CREATEDAT time.Time `json:"created_at"`
	UPDATEDAT time.Time `json:"updated_at"`
}
