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
)

// @title Devboard API
// @version 1.0
// @description API REST for project management
// @host  localhost:8080
// @BasePath /api/v1
func main() {

	logger := setupLogger()

	userRepo := memory.NewUserRepository()
	taskRepo := memory.NewTaskRepository()

	userUC := usecases.NewUserUseCase(userRepo)
	taskUC := usecases.NewTaskUseCase(taskRepo, taskRepo)

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
	taskHandler := handler.NewTaskHandler(taskUC, validate, logger)

	docsHandler, err := handler.NewDocsHandler(devboard.SwaggerJSON)
	if err != nil {
		logger.Error("Fatal error", "error", err)
		os.Exit(1)
	}

	srv.RegisterRoutes("GET /docs", docsHandler)
	srv.RegisterRoutes("GET /api/v1/health", healthHandler)
	// User routes
	srv.RegisterRoutes("POST /api/v1/users", http.HandlerFunc(userHandler.Create))
	srv.RegisterRoutes("GET /api/v1/users/{id}", http.HandlerFunc(userHandler.Get))
	// Task routes
	srv.RegisterRoutes("POST /api/v1/tasks", http.HandlerFunc(taskHandler.Create))
	srv.RegisterRoutes("GET /api/v1/tasks/{id}", http.HandlerFunc(taskHandler.Get))
	srv.RegisterRoutes("GET /api/v1/projects/{id}/tasks", http.HandlerFunc(taskHandler.ListByProject))
	srv.RegisterRoutes("PUT /api/v1/tasks/{id}/assign", http.HandlerFunc(taskHandler.Assign))
	srv.RegisterRoutes("PUT /api/v1/tasks/{id}/status", http.HandlerFunc(taskHandler.UpdateStatus))

	srv.UseChain(
		middlewares.Recovery(logger),
		middlewares.Logger(logger),
	)

	if err := srv.Start(); err != nil {
		logger.Error("Fatal error", "error", err)
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
