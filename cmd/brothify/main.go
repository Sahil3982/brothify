package main

import (
	"log"
	"net/http"
	"os"

	"github.com/brothify/internal/repositories"
	"github.com/brothify/internal/router"
	"github.com/brothify/internal/services"
	"github.com/brothify/pkg/database"
)

func main() {
	db := database.ConnectingDb()
	defer db.Close()

	dishRepo := repositories.NewDishRepository(db)
	dishService := services.NewDishService(dishRepo)
	dishHandler := router.NewDishHandler(dishService)
	mux := router.NewRouter(dishHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default for local development
	}

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed:", err)
	}
}
