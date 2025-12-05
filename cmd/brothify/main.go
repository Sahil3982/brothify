package main

import (
	"log"
	"net/http"
	"os"

	"github.com/brothify/internal/config"
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

	if err := database.RunMigration(db); err != nil {
		log.Fatal("Migration Error:", err)
	}

	dishRepo := repositories.NewDishRepository(db)
	userRepo := repositories.NewUserRepository(db)
	reservationRepo := repositories.NewReservationRepository(db)

	userService := services.NewUserService(userRepo)
	dishService := services.NewDishService(dishRepo)
	reservationService := services.NewReservationService(reservationRepo)

	userHandler := router.NewUserHandler(userService)
	dishHandler := router.NewDishHandler(dishService)
	reservationHandler := router.NewReservationHandler(reservationService)

	mux := router.NewRouter(dishHandler, userHandler, reservationHandler)
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
