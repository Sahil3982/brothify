package router

import (
	"fmt"
	"net/http"
)

func NewRouter(dishHandler *DishHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK Health is good!")
	})
	mux.HandleFunc("/v1/api/category", dishCategory)
	mux.Handle("/v1/api/dishes", dishHandler)

	return mux
}
