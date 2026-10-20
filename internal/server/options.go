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

// Option is a function that update config
type Option func(*config)

// WithReadTimeout configure read timeout
func WithReadTimeout(duration time.Duration) Option {
	return func(c *config) {
		c.readTimeout = duration
	}
}

// WithWriteTimeout configure write timeout
func WithWriteTimeout(duration time.Duration) Option {
	return func(c *config) {
		c.writeTimeout = duration
	}
}

// WithIdleTimeout configure idle timeout
func WithIdleTimeout(duration time.Duration) Option {
	return func(c *config) {
		c.idleTimeout = duration
	}
}

// WithReadHeaderTimeout configure readHeader timeout
func WithReadHeaderTimeout(duration time.Duration) Option {
	return func(c *config) {
		c.readHeaderTimeout = duration
	}
}

// WithLogger configure logger
func WithLogger(logger *slog.Logger) Option {
	return func(c *config) {
		c.logger = logger
	}
}
