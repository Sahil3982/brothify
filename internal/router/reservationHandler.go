package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	case http.MethodPatch:
		h.UpdateReservation(w, r)
	case http.MethodDelete:
		h.DeleteReservation(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusBadRequest)
	}
}

func (h *ReservationHandler) GetAllReservations(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	search := q.Get("search")
	status := q.Get("status")
	date := q.Get("date")

	page, _ := strconv.Atoi(q.Get("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	reservations, err := h.service.GetAllReservations(search, status, date, limit, offset)
	if err != nil {
		log.Println("Failed to retrieve reservations", err)
		helpers.Error(w, http.StatusInternalServerError, "Failed to retrieve reservations")
		return
	}

	helpers.JSON(w, http.StatusOK, "Reservations fetched successfully", map[string]interface{}{
		"page":         page,
		"limit":        limit,
		"total_items":  len(reservations),
		"reservations": reservations,
	})
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
	if d.USERID <= 0 {
		helpers.Error(w, http.StatusBadRequest, "Please provide user Id")
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
		log.Panicln("Failed to create reservation", err)
		helpers.Error(w, http.StatusInternalServerError, "Failed to create reservation")
		return
	}

	helpers.JSON(w, http.StatusCreated, "Reservation created successfully", reservationData)
}

func (h *ReservationHandler) UpdateReservation(w http.ResponseWriter, r *http.Request) {
	var d models.Reservation
	params := helpers.ExtractIDFromPath(r)
	data, err := h.service.UpdateReservation(&d, params)
	if err != nil {
		log.Println("Reservation not updated", err)
		helpers.Error(w, http.StatusBadRequest, "Reservation not updated")
		return 
	}

	helpers.JSON(w, http.StatusOK, "Reservation Updated successfully", data)

}

func (h *ReservationHandler) DeleteReservation(w http.ResponseWriter, r *http.Request) {
		var d models.Reservation
		params := helpers.ExtractIDFromPath(r)
		err := h.service.DeleteReservation(&d,params)
		if err != nil {
			log.Println("Reservation not deleted")
			helpers.Error(w, http.StatusBadGateway,"Reservation not deleted")
			return 
		}
		helpers.JSON(w,http.StatusOK, "Reservation successfully deleted")
}
