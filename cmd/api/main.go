// Package main
package main

import (
	"log"
	"log/slog"

	"github.com/FranciscoJavierMartin/devboard/internal/handler"
	"github.com/FranciscoJavierMartin/devboard/internal/logger"
	"github.com/FranciscoJavierMartin/devboard/internal/middlewares"
	"github.com/FranciscoJavierMartin/devboard/internal/server"
)

func main() {
	logger := logger.New(logger.DefaultConfig())

	srv := server.New(":8080")

	srv.Use(middlewares.Recovery(logger))
	srv.Use(middlewares.Logger(logger))

	healthHandler := handler.NewHealthHandler()

	srv.RegisterRoutes("GET /health", healthHandler)

	logger.Info("Server running", slog.String("addr", ":8080"))

	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
