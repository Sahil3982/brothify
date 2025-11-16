package router

import (
	"fmt"
	"github.com/brothify/internal/middleware"
	"net/http"
)

func NewRouter(dishHandler *DishHandler, userHandler *UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK Health is good!")
	})
	mux.HandleFunc("/v1/api/category/", dishCategory)
	mux.Handle("/v1/api/dishes/", dishHandler)
	mux.Handle("/v1/api/users/", middleware.AuthMiddleware(userHandler))

	return mux
}
