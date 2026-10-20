// Package server creates and start server
package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FranciscoJavierMartin/devboard/internal/middlewares"
)

// Server struct
type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
	logger     *slog.Logger
}

// New creates a new server
func New(addr string, opts ...Option) *Server {
	cfg := defaultConfig()

	for _, opt := range opts {
		opt(&cfg)
	}

	mux := http.NewServeMux() // HTTP Router

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,

		ReadTimeout:       cfg.readTimeout,
		WriteTimeout:      cfg.writeTimeout,
		IdleTimeout:       cfg.idleTimeout,
		ReadHeaderTimeout: cfg.readHeaderTimeout,
	}

	return &Server{
		httpServer: httpServer,
		mux:        mux,
		logger:     cfg.logger,
	}
}

// Start server
func (server *Server) Start() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	defer signal.Stop(quit)

	serverErr := make(chan error, 1)

	go func() {
		server.logger.Info("Server running", slog.String("addr", server.httpServer.Addr))
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case sig := <-quit:
		server.logger.Info("Signal received, shutdown initializing", slog.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		server.logger.Error("Error on shutdown", slog.Any("error", err))
		return err
	}

	server.logger.Info("Server shuted down successful")

	return nil
}

// RegisterRoutes method
func (server *Server) RegisterRoutes(pattern string, handler http.Handler) {
	server.mux.Handle(pattern, handler)
}

// ServeHTTP server
func (server *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server.mux.ServeHTTP(writer, request)
}

// Use function for middlewares
func (server *Server) Use(middleware func(http.Handler) http.Handler) {
	server.httpServer.Handler = middleware(server.httpServer.Handler)
}

// UseChain chain middlewares
func (server *Server) UseChain(middlewaresParams ...middlewares.Middleware) {
	server.httpServer.Handler = middlewares.Chain(server.httpServer.Handler, middlewaresParams...)
}
