package router

import (
	"fmt"
	"net/http"

	"github.com/brothify/internal/handler"
	"github.com/brothify/internal/helpers"
	"github.com/brothify/internal/middleware"
)

func NewRouter(dishHandler *handler.DishHandler, userHandler *handler.UserHandler, reservationHandler *handler.ReservationHandler, paymentHandler *handler.PaymentHandler, categoryHandler *handler.CategoryHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK Health is good!")
	})
	// These endpoints are part of the public checkout flow.
	mux.Handle("/v1/api/payment/order", helpers.PostMethod(handler.CreateRazorpayOrder))
	mux.Handle("/v1/api/payment/verify", helpers.PostMethod(paymentHandler.VerifyRazorpayPayment))
	mux.Handle("/v1/api/categories/", protectNonGet(categoryHandler))
	mux.Handle("/v1/api/dishes/", protectNonGet(dishHandler))
	mux.Handle("/v1/api/login/", userHandler)
	mux.Handle("/v1/api/reservations/", publicPost(reservationHandler))

	return mux
}

func publicPost(next http.Handler) http.Handler {
	protected := middleware.AuthMiddleware(next)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}
		protected.ServeHTTP(w, r)
	})
}

func protectNonGet(next http.Handler) http.Handler {
	protected := middleware.AuthMiddleware(next)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		protected.ServeHTTP(w, r)
	})
}
