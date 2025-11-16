package auth

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var secretKey string

func init() {
	godotenv.Load()
	secretKey = os.Getenv("JWT_SECRET_KEY")
}

func GenerateToken(userID int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil

}

func VerifyToken(tokenStirng string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStirng, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return token, nil
}
