package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/brothify/pkg/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(tokenString, " ")

		token, err := auth.VerifyToken(parts[1])
		if err != nil || !token.Valid {
			fmt.Println("err", err)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}
		fmt.Println("token", token)
		next.ServeHTTP(w, r)

	})
}

// func isAdmin(next http.Handler) http.Handler {

// }

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Allow any domain
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allow common headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Allow methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
