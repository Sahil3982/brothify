package router

import (
	"encoding/json"
	"net/http"

	"github.com/brothify/internal/helpers"
	"github.com/brothify/internal/models"
	"github.com/brothify/internal/services"
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
	case http.MethodPost:
		h.createDish(w, r)
	case http.MethodPatch:
		h.updateDish(w, r)
	case http.MethodDelete:
		h.deleteDish(w, r)
	default:
		helpers.Error(w, http.StatusBadRequest, "Invalid request method")
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

func (h *DishHandler) createDish(w http.ResponseWriter, r *http.Request) {
	var d models.Dish
	err := json.NewDecoder(r.Body).Decode(&d); if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if  d.NAME == "" || d.PRICE <= 0 {
		http.Error(w, "Missing or invalid dish fields", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateDish(&d); err != nil {
		http.Error(w, "Failed to create dish", http.StatusInternalServerError)
		return
	}

	helpers.JSON(w, http.StatusCreated, map[string]string{"message": "Dish created successfully"})
}

func (h *DishHandler) updateDish(w http.ResponseWriter, r *http.Request) {

}

func (h *DishHandler) deleteDish(w http.ResponseWriter, r *http.Request) {

}
