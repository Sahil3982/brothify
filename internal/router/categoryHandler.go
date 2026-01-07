package router

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/brothify/internal/helpers"
	"github.com/brothify/internal/models"
	"github.com/brothify/internal/services"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// ServeHTTP implements http.Handler.
func (h *CategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAllCategories(w, r)
	case http.MethodPost:
		h.createCategory(w, r)
	case http.MethodPut:
		h.updateCategory(w, r)
	case http.MethodDelete:
		h.deleteCategory(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) getAllCategories(w http.ResponseWriter, r *http.Request) {
	// Implementation for getting all categories
	category, err := h.service.GetAllCategories()
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to retrieve categories")
		return
	}

	helpers.JSON(w, http.StatusOK, "Categories fetched successfully", category)
}

func (h *CategoryHandler) createCategory(w http.ResponseWriter, r *http.Request) {
	var d models.Category
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	category, err := h.service.CreateCategory(&d)
	if err != nil {
		log.Println("Failed to create category", err)
		helpers.Error(w, http.StatusInternalServerError, "Failed to create category")
		return
	}
	helpers.JSON(w, http.StatusCreated, "Category created successfully", category)

}

func (h *CategoryHandler) updateCategory(w http.ResponseWriter, r *http.Request) {
	var d models.Category
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	id := helpers.ExtractIDFromPath(r)
	if id == "" {
		helpers.Error(w, http.StatusBadRequest, "Category id not provided")
		return
	}

	category, err := h.service.UpdateCategory(&d, id)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to update category")
		return
	}
	helpers.JSON(w, http.StatusOK, "Category updated successfully", category)
}

func (h *CategoryHandler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	id := helpers.ExtractIDFromPath(r)
	if id == "" {
		helpers.Error(w, http.StatusBadRequest, "Category id not provided")
		return
	}
	err := h.service.DeleteCategory(id)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to delete category")
		return
	}
	helpers.JSON(w, http.StatusOK, "Category deleted successfully", nil)
}
