package router

import (
	"encoding/json"
	"log"
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
	dishes, err := h.service.GetAllDishes()
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to retrieve dishes")
		return
	}
	helpers.JSON(w, http.StatusOK, dishes)
}

func (h *DishHandler) createDish(w http.ResponseWriter, r *http.Request) {
	log.Println("✅ CreateDish called")
	var d models.Dish
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if d.NAME == "" || d.PRICE <= 0 {
		http.Error(w, "Missing or invalid dish fields", http.StatusBadRequest)
		return
	}

	createdDish, err := h.service.CreateDish(&d)
	if err != nil {
		http.Error(w, "Failed to create dish", http.StatusInternalServerError)
		return
	}

	helpers.JSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Dish created successfully",
		"data":    createdDish,
	})
}

func (h *DishHandler) updateDish(w http.ResponseWriter, r *http.Request) {

	id := helpers.ExtractIDFromPath(r)
	if id == "" {
		helpers.Error(w, http.StatusBadRequest, "Dish ID not provided")
		return
	}
	log.Println("✅ Dishes_id:", id)

	var d models.Dish
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// validation logic here
	if d.NAME == "" || d.PRICE <= 0 {
		http.Error(w, "Missing or invalid dish fields", http.StatusBadRequest)
		return
	}

	// call serive to update dish

	if err := h.service.UpdateDish(id, &d); err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to update dish")
		return
	}

	helpers.JSON(w, http.StatusOK, map[string]string{"message": "Dish updated successfully"})

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

	helpers.JSON(w, http.StatusOK, map[string]string{"message": "Dish deleted successfully"})

}
