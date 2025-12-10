package utils

import (
	"strings"
)

func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func ValidateMobileNumber(mobile string) bool {
	if len(mobile) < 7 || len(mobile) > 15 {
		return false
	}
	for _, ch := range mobile {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
