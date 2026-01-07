package router

import (
	"fmt"
	"net/http"

	"github.com/brothify/internal/handler"
	"github.com/brothify/internal/helpers"
)

func NewRouter(dishHandler *DishHandler, userHandler *UserHandler, reservationHandler *ReservationHandler , paymentHandler *handler.PaymentHandler, categoryHandler *CategoryHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK Health is good!")
	})
	mux.HandleFunc("/v1/api/payment/order", helpers.PostMethod(handler.CreateRazorpayOrder))
	mux.HandleFunc("/v1/api/payment/verify", helpers.PostMethod(paymentHandler.VerifyRazorpayPayment))
	mux.Handle("/v1/api/category/", categoryHandler)
	mux.Handle("/v1/api/dishes/", dishHandler)
	mux.Handle("/v1/api/login/", userHandler)
	mux.Handle("/v1/api/reservations/", reservationHandler)

	return mux
}
