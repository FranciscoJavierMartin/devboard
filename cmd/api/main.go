// Package main
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/FranciscoJavierMartin/devboard"
	"github.com/FranciscoJavierMartin/devboard/internal/handler"
	"github.com/FranciscoJavierMartin/devboard/internal/logger"
	"github.com/FranciscoJavierMartin/devboard/internal/middlewares"
	pgrepo "github.com/FranciscoJavierMartin/devboard/internal/repository/postgres"
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

	dbURL := os.Getenv("DATABASE_URL")
	pool, err := pgrepo.NewPool(context.Background(), pgrepo.Config{
		URL:             dbURL,
		MaxConns:        12,
		MinConns:        5,
		MaxConnLifeTime: 1 * time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	})

	if err != nil {
		logger.Error("Could not connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	defer pool.Close()

	logger.Info("Connection to database")

	// Memory
	// userRepo := memory.NewUserRepository()
	// taskRepo := memory.NewTaskRepository()

	userRepo := pgrepo.NewUserRepository(pool)
	taskRepo := pgrepo.NewTaskRepository(pool)

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
