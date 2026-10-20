package server

import (
	"log/slog"
	"time"
)

type config struct {
	readTimeout       time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration
	readHeaderTimeout time.Duration
	logger            *slog.Logger
}

func defaultConfig() config {
	return config{
		readTimeout:       10 * time.Second,
		writeTimeout:      30 * time.Second,
		idleTimeout:       60 * time.Second,
		readHeaderTimeout: 5 * time.Second,
		logger:            slog.Default(),
	}
}

type Option func(*config)

func WithReadTimeout(duration time.Duration) Option {
	return func(c *config) {
		c.readTimeout = duration
	}
}

func WithWriteTimeout(duration time.Duration) Option {
	return func(c *config) {
		c.writeTimeout = duration
	}
}

func WithIdleTimeout(duration time.Duration) Option {
	return func(c *config) {
		c.idleTimeout = duration
	}
}

func WithReadHeaderTimeout(duration time.Duration) Option {
	return func(c *config) {
		c.readHeaderTimeout = duration
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(c *config) {
		c.logger = logger
	}
}
