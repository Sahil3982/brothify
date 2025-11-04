package router

import (
	"net/http"
	"github.com/brothify/internal/helpers"
)

type Category struct {
	CATID   int    `json:"cat_id"`
	CATNAME string `json:"cat_name"`
}

func dishCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		categories := []Category{
			{CATID: 1, CATNAME: "Rice"},
			{CATID: 2, CATNAME: "Soup"},
			{CATID: 3, CATNAME: "Salad"},
			{CATID: 4, CATNAME: "Bread"},
			{CATID: 5, CATNAME: "Juice"},
		}

		helpers.JSON(w, http.StatusOK, categories)
		return
	}

	helpers.Error(w, http.StatusBadRequest, "Invalid request method")
}
