package router

import (
	"github.com/brothify/internal/services"
	"net/http"
)

type CategoryHandler struct {
	service *services.CategoryService
}

// ServeHTTP implements http.Handler.
func (c *CategoryHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("unimplemented")
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}
