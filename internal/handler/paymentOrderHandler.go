package handler

import (
	"encoding/json"
	"net/http"

	"github.com/brothify/internal/config"
	"github.com/brothify/internal/helpers"
)

type CreateOrderRequest struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
	Receipt  string `json:"receipt"`
}

func CreateRazorpayOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	data := map[string] interface{}{
		"amount" : req.Amount,
		"currency": req.Currency,
		"receipt": req.Receipt,
	}

	order,err := config.RazorpayClient.Order.Create(data, nil)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to create order")
		return
	}
	
	helpers.JSON(w, http.StatusOK, "Order created successfully", order)
}
