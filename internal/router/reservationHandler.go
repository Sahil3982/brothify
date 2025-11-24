package router

import (
	"net/http"
	service "github.com/brothify/internal/services"
)

type ReservationHandler struct {
	service *service.ReservationService
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		GetAllReservations(w, r)
	case http.MethodPost:
		CreateReservation(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}
}

func GetAllReservations(w http.ResponseWriter, r *http.Request) {
	// Implementation for fetching all reservations
}

func CreateReservation(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating a new reservation
}
