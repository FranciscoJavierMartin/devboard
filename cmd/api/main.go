// Package main
package main

import (
	"log"
	"net/http"

	"github.com/FranciscoJavierMartin/devboard/internal/server"
)

func main() {
	srv := server.New(":8080")

	srv.RegisterRoutes("GET /health", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(`{"status": "ok"}`))
	}))

	log.Println("Server running on :8080")

	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
