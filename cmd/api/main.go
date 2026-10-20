// Package main
package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/FranciscoJavierMartin/devboard"
	"github.com/FranciscoJavierMartin/devboard/internal/handler"
	"github.com/FranciscoJavierMartin/devboard/internal/logger"
	"github.com/FranciscoJavierMartin/devboard/internal/middlewares"
	"github.com/FranciscoJavierMartin/devboard/internal/notification"
	"github.com/FranciscoJavierMartin/devboard/internal/server"
	"github.com/FranciscoJavierMartin/devboard/internal/validator"
)

// @title Devboard API
// @version 1.0
// @description API REST for project management
// @host  localhost:8080
// @BasePath /
func main() {
	env := os.Getenv("GO_ENV")

	var log *slog.Logger

	if env == "production" {
		log = logger.New(logger.ProductionConfig())
	} else {
		log = logger.New(logger.DefaultConfig())
	}

	validate := validator.New()
	notifier := notification.NewLogNotifier(log)
	_ = notifier

	srv := server.New(
		":8080",
		server.WithLogger(log),
		server.WithReadTimeout(15*time.Second),
		server.WithWriteTimeout(30*time.Second),
	)

	healthHandler := handler.NewHealthHandler()

	docsHandler, err := handler.NewDocsHandler(devboard.SwaggerJSON)
	if err != nil {
		log.Error("Fatal error", "error", err)
		os.Exit(1)
	}

	srv.RegisterRoutes("GET /docs", docsHandler)
	srv.RegisterRoutes("GET /health", healthHandler)

	srv.UseChain(
		middlewares.Recovery(log),
		middlewares.Logger(log),
	)

	if err := srv.Start(); err != nil {
		log.Error("Fatal error", "error", err)
		os.Exit(1)
	}

}
