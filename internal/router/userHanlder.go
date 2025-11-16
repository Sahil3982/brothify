package router

import (
	"net/http"

	"github.com/brothify/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Implement user handling logic here
}
