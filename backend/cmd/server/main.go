package main

import (
	"log"
	"net/http"
	"os"

	"backend/internal/handler"
	"backend/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	calcService := service.NewCalculatorService()
	router := handler.NewRouter(calcService)

	addr := ":" + port
	log.Printf("Server running on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
