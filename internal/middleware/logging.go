package middleware

import (
	"context"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

// LoggingMiddleware provides request logging
type LoggingMiddleware struct {
	logger domain.Logger
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger domain.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// Process implements the Middleware interface
func (m *LoggingMiddleware) Process(ctx context.Context, next domain.HandlerFunc) domain.HandlerFunc {
	return func(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
		start := time.Now()
		
		// Log request start
		m.logger.Info("Command processing started",
			"command", cmd.Text,
			"user_id", cmd.User.TelegramID,
			"username", cmd.User.Username,
			"timestamp", start)

		// Execute next handler
		response, err := next(ctx, cmd)
		
		duration := time.Since(start)
		
		if err != nil {
			m.logger.Error("Command processing failed",
				"command", cmd.Text,
				"user_id", cmd.User.TelegramID,
				"duration", duration,
				"error", err)
		} else {
			m.logger.Info("Command processing completed",
				"command", cmd.Text,
				"user_id", cmd.User.TelegramID,
				"duration", duration,
				"response_length", len(response.Text))
		}

		return response, err
	}
}