package router

import (
	"net/http"

	"github.com/brothify/internal/services"
	"github.com/brothify/internal/helpers"
)

type DishHandler struct {
	service *services.DishService
}

func NewDishHandler(service *services.DishService) *DishHandler {
	return &DishHandler{service: service}
}

func (h *DishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllDishes(w, r)
	}
}

func (h *DishHandler) getAllDishes(w http.ResponseWriter, r *http.Request) {
	dishes, err := h.service.GetAllDishes()
	if err != nil {
		http.Error(w, "Failed to retrieve dishes", http.StatusInternalServerError)
		return
	}
	helpers.JSON(w, http.StatusOK, dishes)
}