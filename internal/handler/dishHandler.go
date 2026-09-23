package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/brothify/internal/dto"
	"github.com/brothify/internal/helpers"
	"github.com/brothify/internal/models"
	"github.com/brothify/internal/services"

	// "github.com/brothify/internal/validators"
	"github.com/google/uuid"
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
		basePath := strings.TrimSuffix(r.URL.Path, "/")
		if basePath != "/v1/api/dishes" {
			h.GetDishById(w, r)
			return
		}
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

func (h *DishHandler) GetDishById(w http.ResponseWriter, r *http.Request) {
	id := helpers.ExtractIDFromPath(r)

	dishID, err := uuid.Parse(id)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid dish ID")
		return
	}
	dish, err := h.service.GetDishByID(dishID)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to retrieve dish")
		return
	}
	helpers.JSON(w, http.StatusOK, "Dish fetched successfully", dish)
}

func (h *DishHandler) getAllDishes(w http.ResponseWriter, r *http.Request) {
	dishes, err := h.service.GetAllDishes()
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to retrieve dishes")
		return
	}
	helpers.JSON(w, http.StatusOK, "dishes fetch successfully", dishes)
}

func (h *DishHandler) createDish(w http.ResponseWriter, r *http.Request) {
	var req dto.DishRequest

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	req.NAME = r.FormValue("dish_name")
	if req.NAME == "" {
		req.NAME = r.FormValue("name")
	}
	if req.NAME == "" {
		helpers.Error(w, http.StatusBadRequest, "Dish name is required")
		return
	}
	req.DESCRIPTION = r.FormValue("description")
	categoryID := r.FormValue("category_id")
	if categoryID == "" {
		helpers.Error(w, http.StatusBadRequest, "Category ID is required")
		return
	}
	req.CATEGORYID, err = uuid.Parse(categoryID)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	priceStr := r.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid price")
		return
	}
	req.PRICE = price

	file, fileHeader, err := r.FormFile("dish_url")
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Failed to retrieve image file")
		return
	}
	defer file.Close()

	createdDish, err := h.service.CreateDish(r.Context(), &req, file, fileHeader)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, err.Error())
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

	// Ensure the dish ID is set from the path
	parsedID, err := uuid.Parse(id)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid dish ID format")
		return
	}

	var d dto.DishRequest
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			helpers.Error(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		d.NAME = r.FormValue("dish_name")
		d.DESCRIPTION = r.FormValue("description")
		d.DISHURL = r.FormValue("dish_url")
		d.CATEGORYID, err = uuid.Parse(r.FormValue("category_id"))
		if err != nil {
			helpers.Error(w, http.StatusBadRequest, "Invalid category ID")
			return
		}
		d.PRICE, err = strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			helpers.Error(w, http.StatusBadRequest, "Invalid price")
			return
		}
	} else if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if d.NAME == "" || d.PRICE <= 0 {
		helpers.Error(w, http.StatusBadRequest, "Name and Price are required")
		return
	}

	m := models.Dish{
		ID:           parsedID,
		NAME:         d.NAME,
		CATEGORYID:   &d.CATEGORYID,
		PRICE:        d.PRICE,
		DESCRIPTION:  d.DESCRIPTION,
		DISHURL:      d.DISHURL,
		AVAILABILITY: d.AVAILABILITY,
		RATING:       d.RATING,
		HIGHLIGHT:    d.HIGHLIGHT,
	}

	updated, err := h.service.UpdateDish(parsedID, &m)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to update dish")
		return
	}
	if !updated {
		helpers.Error(w, http.StatusNotFound, "Dish not found with given ID")
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

	dishID, err := uuid.Parse(id)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid dish ID format")
		return
	}

	deleted, err := h.service.DeleteDish(dishID)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to delete dish")
		return
	}
	if !deleted {
		helpers.Error(w, http.StatusNotFound, "Dish not found with given ID")
		return
	}

	helpers.JSON(w, http.StatusOK, "Dish deleted successfully")

}
