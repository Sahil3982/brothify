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
		helpers.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	res := r.Body.Close()
	if res != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to close request body")
		return
	}
	//validate and create reservation logic here
	if d.NUMBEROFGUESTS <= 0 {
		helpers.Error(w, http.StatusBadRequest, "Number of guests must be greater than zero")
		return
	}
	if d.RESERVATIONTIME == "" {
		helpers.Error(w, http.StatusBadRequest, "Reservation time is required")
		return
	}
	if d.RESERVATIONPERSONNAME == "" {
		helpers.Error(w, http.StatusBadRequest, "Reservation person name is required")
		return
	}
	if !utils.ValidateEmail(d.RESERVATIONPERSONEMAIL) {
		helpers.Error(w, http.StatusBadRequest, "Reservation person email is required")
		return
	}
	if !utils.ValidateMobileNumber(d.RESERVATIONPERSONMOBILENUMBER) {
		helpers.Error(w, http.StatusBadRequest, "Reservation person mobile number is required")
		return
	}
	
	// Call service to create reservation

	reservationData, err := h.service.CreateReservation(&d)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to create reservation")
		return
	}

	helpers.JSON(w, http.StatusCreated, "Reservation created successfully", reservationData)
}
