package models

type Category struct {
	ID          int    `json:"cat_id"`
	NAME        string `json:"name"`
	DESCRIPTION string `json:"description"`
	SLUG		string `json:"slug"`
}
