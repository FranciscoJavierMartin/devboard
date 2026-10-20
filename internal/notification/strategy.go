package notification

import (
	"context"
	"log/slog"
)

type Notifier interface {
	Notify(ctx context.Context, userID string, message string) error
}

type LogNotifier struct {
	logger *slog.Logger
}

func NewLogNotifier(logger *slog.Logger) *LogNotifier {
	return &LogNotifier{
		logger: logger,
	}
}

func (notifier *LogNotifier) Notify(ctx context.Context, userID string, message string) error {
	notifier.logger.InfoContext(
		ctx,
		"sent notification",
		slog.String("user_id", userID),
		slog.String("message", message),
	)

	return nil
}

type NoOpNotifier struct{}

func (notifier *NoOpNotifier) Notify(_ context.Context, _ string, _ string) error {
	return nil
}
