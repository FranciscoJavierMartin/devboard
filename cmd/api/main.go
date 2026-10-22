// Package main
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/FranciscoJavierMartin/devboard"
	"github.com/FranciscoJavierMartin/devboard/internal/handler"
	"github.com/FranciscoJavierMartin/devboard/internal/logger"
	"github.com/FranciscoJavierMartin/devboard/internal/middlewares"
	"github.com/FranciscoJavierMartin/devboard/internal/repository/memory"
	"github.com/FranciscoJavierMartin/devboard/internal/server"
	"github.com/FranciscoJavierMartin/devboard/internal/usecases"
	"github.com/FranciscoJavierMartin/devboard/internal/validator"
	"github.com/gofiber/fiber/v2/log"
)

// @title Devboard API
// @version 1.0
// @description API REST for project management
// @host  localhost:8080
// @BasePath /
func main() {

	logger := setupLogger()

	userRepo := memory.NewUserRepository()

	userUC := usecases.NewUserUseCase(userRepo)

	validate := validator.New()
	// notifier := notification.NewLogNotifier(log)
	// _ = notifier

	srv := server.New(
		":8080",
		server.WithLogger(logger),
		server.WithReadTimeout(15*time.Second),
		server.WithWriteTimeout(30*time.Second),
	)

	healthHandler := handler.NewHealthHandler()
	userHandler := handler.NewUserHandler(userUC, validate, logger)

	docsHandler, err := handler.NewDocsHandler(devboard.SwaggerJSON)
	if err != nil {
		log.Error("Fatal error", "error", err)
		os.Exit(1)
	}

	srv.RegisterRoutes("GET /docs", docsHandler)
	srv.RegisterRoutes("GET /health", healthHandler)
	srv.RegisterRoutes("POST /api/v1/users", http.HandlerFunc(userHandler.Create))
	srv.RegisterRoutes("GET /api/v1/users/{id}", http.HandlerFunc(userHandler.Get))

	srv.UseChain(
		middlewares.Recovery(logger),
		middlewares.Logger(logger),
	)

	if err := srv.Start(); err != nil {
		log.Error("Fatal error", "error", err)
		os.Exit(1)
	}

}

func setupLogger() *slog.Logger {
	env := os.Getenv("GO_ENV")
	slog.Info("Loaded environment", "GO_ENV", env)

	var log *slog.Logger

	if env == "production" {
		log = logger.New(logger.ProductionConfig())
	} else {
		log = logger.New(logger.DefaultConfig())
	}

	return log
}
