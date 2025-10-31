package main

import (
	"log"
	"net/http"

	"github.com/brothify/internal/router"
	"github.com/brothify/pkg/database"
)

func main() {
	database.ConnectingDb()
	mux := router.NewRouter()
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
	log.Fatal("Server failed:", err)
	}

}
