package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	log.Println("✅ DishHandler ServeHTTP called with method:", r.Method)
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
	log.Println("✅ GetAllDishes called")
	id := helpers.ExtractIDFromPath(r)

	dishes, err := h.service.GetAllDishes()

	dishID, _ := strconv.Atoi(id)
	if id != "" {
		for _, dish := range dishes {
			if dish.ID == dishID {
				helpers.JSON(w, http.StatusOK, "dish fetch successfully", dish)
				return
			}
		}
	}
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to retrieve dishes")
		return
	}
	helpers.JSON(w, http.StatusOK, "dishes fetch successfully", dishes)
}

func (h *DishHandler) createDish(w http.ResponseWriter, r *http.Request) {
	var d models.Dish
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if d.NAME == "" || d.PRICE <= 0 {
		http.Error(w, "Name and Price are missing please provide", http.StatusBadRequest)
		return
	}

	createdDish, err := h.service.CreateDish(&d)
	if err != nil {
		http.Error(w, "Failed to create dish", http.StatusInternalServerError)
		return
	}

	helpers.JSON(w, http.StatusCreated, "Dish created successfully", createdDish)
}

func (h *DishHandler) updateDish(w http.ResponseWriter, r *http.Request) {
	id := helpers.ExtractIDFromPath(r)
	if id == "" {
		helpers.Error(w, http.StatusBadRequest, "Dish ID not provided")
		return
	}

	dishID, err := strconv.Atoi(id)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid dish ID")
		return
	}

	var d models.Dish
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if d.NAME == "" || d.PRICE <= 0 {
		helpers.Error(w, http.StatusBadRequest, "Name and Price are required")
		return
	}

	// ✅ First check if dish exists
	allDishes, err := h.service.GetAllDishes()
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to fetch dishes for validation")
		return
	}

	var exists bool
	for _, dish := range allDishes {
		if dish.ID == dishID {
			exists = true
			break
		}
	}

	if !exists {
		helpers.Error(w, http.StatusNotFound, "Dish not found with given ID")
		return
	}

	// Ensure the dish ID is set from the path
	d.ID = dishID

	// ✅ Proceed to update
	if err := h.service.UpdateDish(id, &d); err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to update dish")
		return
	}

	helpers.JSON(w, http.StatusOK, "Dish updated successfully", d)
}

func (h *DishHandler) deleteDish(w http.ResponseWriter, r *http.Request) {

	id := helpers.ExtractIDFromPath(r)
	if id == "" {
		helpers.Error(w, http.StatusBadRequest, "Dish ID not provided")
		return
	}

	if err := h.service.DeleteDish(id); err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to delete dish")
		return
	}

	helpers.JSON(w, http.StatusOK, "Dish deleted successfully")

}
