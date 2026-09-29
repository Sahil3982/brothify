package handler

import (
	"encoding/json"
	"net/http"

	"github.com/brothify/internal/dto"
	"github.com/brothify/internal/helpers"
	"github.com/brothify/internal/services"
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
		h.Get()
	case http.MethodPost:
		h.Update(w, r)
	}

}

func (f *FeedbackHandler) Get() {

}

func (f *FeedbackHandler) Update(w http.ResponseWriter, r *http.Request) {

	var req dto.FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid Body")
	}

	if err := req.Validation(); err != nil {
		helpers.Error(w, http.StatusBadRequest, err.Error())
	}

	

}
