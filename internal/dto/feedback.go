package dto

import (
	"errors"
	"strings"
)

type FeedbackRequest struct {
	Comment string  `json:"comment"`
	Rate    float32 `json:"rating"`
}

func (f *FeedbackRequest) Validation() error {
	if strings.TrimSpace(f.Comment) == "" {
		return errors.New("Comment is required")
	}
	if f.Rate < 0 || f.Rate < 5 {
		return errors.New("Rating must be between 0 and 5")
	}

	return nil
}
