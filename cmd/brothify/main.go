package main

import (
	"log"
	"net/http"
	"os"

	"github.com/brothify/internal/repositories"
	"github.com/brothify/internal/router"
	"github.com/brothify/internal/services"
	"github.com/brothify/pkg/database"
	"github.com/brothify/internal/middleware"
)

func main() {
	db := database.ConnectingDb()
	defer db.Close()

	dishRepo := repositories.NewDishRepository(db)
	userRepo := repositories.NewUserRepository(db)

	userService := services.NewUserService(userRepo)	
	dishService := services.NewDishService(dishRepo)

	userHandler := router.NewUserHandler(userService)
	dishHandler := router.NewDishHandler(dishService)

	mux := router.NewRouter(dishHandler,userHandler)
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
