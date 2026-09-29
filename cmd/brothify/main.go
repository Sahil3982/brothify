package main

import (
	"log"
	"net/http"
	"os"

	"github.com/brothify/internal/config"
	"github.com/brothify/internal/handler"
	"github.com/brothify/internal/middleware"
	"github.com/brothify/internal/repositories"
	"github.com/brothify/internal/router"
	"github.com/brothify/internal/services"
	"github.com/brothify/pkg/database"
)

func main() {
	db := database.ConnectingDb()
	defer db.Close()
	config.InitAWS()
	config.InitS3()
	config.InitSES()
	config.InitRazorpay()
	emailService := &config.EmailSender{}


	if err := database.RunMigration(db); err != nil {
		log.Fatal("Migration Error:", err)
	}

	dishRepo := repositories.NewDishRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	userRepo := repositories.NewUserRepository(db)
	reservationRepo := repositories.NewReservationRepository(db)
	feedbackRepo := repositories.NewFeedbackRepository(db)
	paymentHandler := handler.NewPaymentHandler(reservationRepo)

	userService := services.NewUserService(userRepo)
	dishService := services.NewDishService(dishRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	reservationService := services.NewReservationService(reservationRepo,emailService)
	feedbackService := services.NewFeedbackService(feedbackRepo)

	userHandler := handler.NewUserHandler(userService)
	dishHandler := handler.NewDishHandler(dishService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	reservationHandler := handler.NewReservationHandler(reservationService)
	feedbackHandler := handler.NewFeedbackHandler(feedbackService)

	mux := router.NewRouter(feedbackHandler, dishHandler, userHandler, reservationHandler, paymentHandler, categoryHandler)
	handler := middleware.CorsMiddleware(mux)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Server failed:", err)
	}
}
