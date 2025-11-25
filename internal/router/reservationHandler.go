package router

import (
	"encoding/json"
	"net/http"

	"github.com/brothify/internal/helpers"
	"github.com/brothify/internal/models"
	"github.com/brothify/internal/services"
	"github.com/brothify/pkg/utils"
)

type ReservationHandler struct {
	service *services.ReservationService
}

func NewReservationHandler(service *services.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

func (h *ReservationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		h.GetAllReservations(w, r)
	case http.MethodPost:
		h.CreateReservation(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}
}

func (h *ReservationHandler) GetAllReservations(w http.ResponseWriter, r *http.Request) {
	// Implementation for fetching all reservations
}

func (h *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var d models.Reservation
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	res := r.Body.Close()
	if res != nil {
		http.Error(w, "Failed to close request body", http.StatusInternalServerError)
		return
	}
	//validate and create reservation logic here
	if d.NUMBEROFGUESTS <= 0 {
		http.Error(w, "Number of guests must be greater than zero", http.StatusBadRequest)
		return
	}
	if d.RESERVATIONTIME == "" {
		http.Error(w, "Reservation time is required", http.StatusBadRequest)
		return
	}
	if d.RESERVATIONPERSONNAME == "" {
		http.Error(w, "Reservation person name is required", http.StatusBadRequest)
		return
	}
	if !utils.ValidateEmail(d.RESERVATIONPERSONEMAIL) {
		http.Error(w, "Reservation person email is required", http.StatusBadRequest)
		return
	}
	if !utils.ValidateMobileNumber(d.RESERVATIONPERSONMOBILENUMBER) {
		http.Error(w, "Reservation person mobile number is required", http.StatusBadRequest)
		return
	}

	// Call service to create reservation

	reservationData, err := h.service.CreateReservation(&d)
	if err != nil {
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		return
	}

	helpers.JSON(w, http.StatusCreated, "Reservation created successfully", reservationData)
}
