package router

import (
	"encoding/json"
	"net/http"
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
