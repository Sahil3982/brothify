package main

import (
	"log"
	"net/http"
	"github.com/brothify/internal/router"	
)

func main() {
	mux := router.NewRouter()
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
	log.Fatal("Server failed:", err)
	}

}
