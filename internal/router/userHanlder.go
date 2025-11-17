package router

import (
	"encoding/json"
	"net/http"

	"github.com/brothify/internal/helpers"
	"github.com/brothify/internal/models"
	"github.com/brothify/internal/services"
	"github.com/brothify/pkg/auth"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.loginUser(w, r)
	}
}

func (h *UserHandler) loginUser(w http.ResponseWriter, r *http.Request) {
	var d models.User
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if d.EMAIL == "" {
		helpers.Error(w, http.StatusBadRequest, "Email is required")
		return
	}

	if d.PASSWORD == "" {
		helpers.Error(w, http.StatusBadRequest, "Password is required")
		return
	}

	user, err := h.service.LoginUser(&d)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Invailed Email and Password")
		return
	}

	token, err := auth.GenerateToken(d.ID)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to generate token: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"token": token,
		"user":  user,
	}
	helpers.JSON(w, http.StatusOK, "Login successful", response)

}
