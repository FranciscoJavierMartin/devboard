// Package notification strategy
package notification

import (
	"context"
	"log/slog"
)

// Notifier is the interface that the contract of notification
type Notifier interface {
	Notify(ctx context.Context, userID string, message string) error
}

// LogNotifier register notification on the logs
type LogNotifier struct {
	logger *slog.Logger
}

// NewLogNotifier is a contructor of LogNotifier
func NewLogNotifier(logger *slog.Logger) *LogNotifier {
	return &LogNotifier{
		logger: logger,
	}
}

// Notify is the method to send notifications in the logs
func (notifier *LogNotifier) Notify(ctx context.Context, userID string, message string) error {
	notifier.logger.InfoContext(
		ctx,
		"sent notification",
		slog.String("user_id", userID),
		slog.String("message", message),
	)

	return nil
}

// NoOpNotifier struct that doen't do nothing (for testing)
type NoOpNotifier struct{}

// Notify method for NoOpNotifier
func (notifier *NoOpNotifier) Notify(_ context.Context, _ string, _ string) error {
	return nil
}
