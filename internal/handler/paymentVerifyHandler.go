package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/brothify/internal/config"
	"github.com/brothify/internal/helpers"
	"github.com/brothify/internal/repositories"
)

type VerifyPaymentRequest struct {
	RazorpayPaymentID string `json:"razorpay_payment_id"`
	RazorpayOrderID   string `json:"razorpay_order_id"`
	RazorpaySignature string `json:"razorpay_signature"`
	ReservationID     int    `json:"reservation_id"`
	Email             string `json:"email"`
}

type PaymentHandler struct {
	ResRepo *repositories.ReservationRepository
}

func NewPaymentHandler(repo *repositories.ReservationRepository) *PaymentHandler {
	return &PaymentHandler{ResRepo: repo}
}


func (h *PaymentHandler) VerifyRazorpayPayment(w http.ResponseWriter, r *http.Request) {
	var req VerifyPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	secret := os.Getenv("RAZORPAY_KEY_SECRET")
	if secret == "" {
		helpers.Error(w, http.StatusInternalServerError, "Payment configuration error")
		return
	}

	// Verify signature
	msg := req.RazorpayOrderID + "|" + req.RazorpayPaymentID
	hmacObj := hmac.New(sha256.New, []byte(secret))
	hmacObj.Write([]byte(msg))
	expectedSignature := hex.EncodeToString(hmacObj.Sum(nil))
	
	if !hmac.Equal([]byte(expectedSignature), []byte(req.RazorpaySignature)) {
		helpers.Error(w, http.StatusBadRequest, "Invalid signature")
		return
	}

	// Convert reservation ID to string for the repository method
	res, err := h.ResRepo.GetReservationByID(strconv.Itoa(req.ReservationID))
	if err != nil {
		helpers.Error(w, http.StatusNotFound, "Reservation not found")
		return
	}

	html, err := helpers.BuildInvoiceHTML(res)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to generate invoice")
		return
	}

	url, err := config.UploadInvoiceToS3(html, req.ReservationID)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to upload invoice")
		return
	}

	// Save invoice URL in DB
	err = h.ResRepo.SaveInvoiceURL(req.ReservationID, req.RazorpayPaymentID, req.RazorpaySignature, url)
	if err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to save invoice URL")
		return
	}

	// Send Email
	if err := config.SendEmailWithInvoice(req.Email, html); err != nil {
		helpers.Error(w, http.StatusInternalServerError, "Failed to send email")
		return
	}

	helpers.JSON(w, http.StatusOK,"Payment verified successfully")
}