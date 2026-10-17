package utils

import (
	"crypto/rand"
	"encoding/base64"
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

func RandomGenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
