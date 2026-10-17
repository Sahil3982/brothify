package handler

import (
	"encoding/json"
	"net/http"

	"github.com/brothify/internal/dto"
	"github.com/brothify/internal/helpers"
	"github.com/brothify/internal/services"
	"github.com/google/uuid"
)

type FeedbackHandler struct {
	service *services.FeedbackService
}

func NewFeedbackHandler(service *services.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{
		service: service,
	}
}

func (h *FeedbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	case http.MethodPost:
		h.Create(w, r)
	case http.MethodPatch:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		helpers.Error(w, http.StatusMethodNotAllowed, "Invalid request method")
	}
}

func (f *FeedbackHandler) Get(w http.ResponseWriter, _ *http.Request) {
	feedback, err := f.service.GetFeedback()
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to retrieve feedback: "+err.Error())
		return
	}
	helpers.JSON(w, http.StatusOK, "Feedback fetched successfully", feedback)
}

func (f *FeedbackHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid Body")
		return
	}

	if err := req.Validation(); err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	err := f.service.CreateFeedback(req)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to create feedback: "+err.Error())
		return
	}

	helpers.JSON(w, http.StatusCreated, "Feedback created successfully")
}

func (f *FeedbackHandler) Update(w http.ResponseWriter, r *http.Request) {

	var req dto.FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid Body")
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		helpers.Error(w, http.StatusBadRequest, "ID is required")
		return
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := req.Validation(); err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	err = f.service.UpdateFeedback(req, idUUID)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to update feedback: "+err.Error())
		return
	}

	helpers.JSON(w, http.StatusOK, "Feedback updated successfully")

}

func (f *FeedbackHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		helpers.Error(w, http.StatusBadRequest, "ID is required")
		return
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	err = f.service.DeleteFeedbackByID(idUUID)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to delete feedback: "+err.Error())
		return
	}

	helpers.JSON(w, http.StatusOK, "Feedback deleted successfully")

}
