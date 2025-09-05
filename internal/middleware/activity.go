package middleware

import (
	"context"

	"yordamchi-dev-bot/database"
	"yordamchi-dev-bot/internal/domain"
)

// ActivityMiddleware logs user activity for analytics
type ActivityMiddleware struct {
	db     *database.DB
	logger domain.Logger
}

// NewActivityMiddleware creates a new activity tracking middleware
func NewActivityMiddleware(db *database.DB, logger domain.Logger) *ActivityMiddleware {
	return &ActivityMiddleware{
		db:     db,
		logger: logger,
	}
}

// Process implements the Middleware interface
func (m *ActivityMiddleware) Process(ctx context.Context, next domain.HandlerFunc) domain.HandlerFunc {
	return func(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
		// Execute the command first
		response, err := next(ctx, cmd)

		// Log activity after successful command execution
		if err == nil && cmd.User != nil {
			// Log user activity in background to avoid blocking response
			go func() {
				logErr := m.db.LogUserActivity(cmd.User.TelegramID, cmd.Text)
				if logErr != nil {
					m.logger.Warn("Failed to log user activity",
						"telegram_id", cmd.User.TelegramID,
						"command", cmd.Text,
						"error", logErr)
				} else {
					m.logger.Debug("User activity logged",
						"telegram_id", cmd.User.TelegramID,
						"command", cmd.Text)
				}
			}()
		}

		return response, err
	}
}