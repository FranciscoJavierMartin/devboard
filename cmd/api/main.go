// Package main
package main

import (
	"os"

	"github.com/FranciscoJavierMartin/devboard"
	"github.com/FranciscoJavierMartin/devboard/internal/handler"
	"github.com/FranciscoJavierMartin/devboard/internal/logger"
	"github.com/FranciscoJavierMartin/devboard/internal/middlewares"
	"github.com/FranciscoJavierMartin/devboard/internal/server"
)

// @title Devboard API
// @version 1.0
// @description API REST for project management
// @host  localhost:8080
// @BasePath /
func main() {
	logger := logger.New(logger.DefaultConfig())

	srv := server.New(":8080", logger)

	srv.Use(middlewares.Recovery(logger))
	srv.Use(middlewares.Logger(logger))

	healthHandler := handler.NewHealthHandler()

	docsHandler, err := handler.NewDocsHandler(devboard.SwaggerJSON)
	if err != nil {
		logger.Error("Fatal error", "error", err)
		os.Exit(1)
	}

	srv.RegisterRoutes("GET /docs", docsHandler)
	srv.RegisterRoutes("GET /health", healthHandler)

	if err := srv.Start(); err != nil {
		logger.Error("Fatal error", "error", err)
		os.Exit(1)
	}

}
