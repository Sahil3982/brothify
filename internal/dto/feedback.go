package dto

import (
	"errors"
	"strings"
)

type FeedbackRequest struct {
	Comment string `json:"comment"`
	Rate    int    `json:"rating"`
}

func (f *FeedbackRequest) Validation() error {
	if strings.TrimSpace(f.Comment) == "" {
		return errors.New("Comment is required")
	}
	if f.Rate < 1 || f.Rate > 5 {
		return errors.New("Rating must be between 1 and 5")
	}

	return nil
}
