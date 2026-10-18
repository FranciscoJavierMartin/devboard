// Package main
package main

import (
	"log"

	"github.com/FranciscoJavierMartin/devboard/internal/server"
)

func main() {
	srv := server.New(":8080")
	log.Println("Server running on :8080")

	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
