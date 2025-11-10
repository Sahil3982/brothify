package main

import (
	"log"
	"net/http"

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

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Server failed:", err)
	}

}
