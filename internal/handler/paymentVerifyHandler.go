package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/brothify/internal/helpers"
)

type VerifyPaymentRequest struct {
	RazorpayPaymentID string `json:"razorpay_payment_id"`
	RazorpayOrderID   string `json:"razorpay_order_id"`
	RazorpaySignature string `json:"razorpay_signature"`
}

func VerifyRazorpayPayment(w http.ResponseWriter, r *http.Request) {
	var req VerifyPaymentRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	secret := os.Getenv("RAZORPAY_KEY_SECRET")

	msg := req.RazorpayOrderID + "|" + req.RazorpayPaymentID
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(msg))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if expectedSignature != req.RazorpaySignature {
		helpers.Error(w, http.StatusBadRequest, "Invalid signature")
		return
	}

	helpers.JSON(w, http.StatusOK, "Payment verified successfully")

}
