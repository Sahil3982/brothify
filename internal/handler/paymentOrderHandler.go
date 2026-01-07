package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/brothify/internal/config"
	"github.com/brothify/internal/helpers"
)

type CreateOrderRequest struct {
	Amount        int    `json:"amount"`
	Currency      string `json:"currency"`
	ReservationID string `json:"reservation_id"`
}

func CreateRazorpayOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	log.Printf("CreateOrderRequest: %+v", req)

	data := map[string]interface{}{
		"amount":         req.Amount,
		"currency":       req.Currency,
		"receipt": req.ReservationID,
	}

	order, err := config.RazorpayClient.Order.Create(data, nil)
	if err != nil {
		log.Println("Failed to create Razorpay order:", err)
		helpers.Error(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	helpers.JSON(w, http.StatusOK, "Order created successfully", order)
}
