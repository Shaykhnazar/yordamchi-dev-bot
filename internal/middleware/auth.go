package middleware

import (
	"context"
	"fmt"

	"yordamchi-dev-bot/internal/domain"
)

// AuthMiddleware provides user authentication and registration
type AuthMiddleware struct {
	userService domain.UserService
	logger      domain.Logger
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(userService domain.UserService, logger domain.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
		logger:      logger,
	}
}

// Process implements the Middleware interface
func (m *AuthMiddleware) Process(ctx context.Context, next domain.HandlerFunc) domain.HandlerFunc {
	return func(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
		// Check if user exists in context
		if cmd.User == nil {
			return nil, fmt.Errorf("user information missing from command")
		}

		// Try to get existing user
		user, err := m.userService.GetUser(ctx, cmd.User.TelegramID)
		if err != nil {
			// User doesn't exist, register them
			m.logger.Info("Registering new user",
				"telegram_id", cmd.User.TelegramID,
				"username", cmd.User.Username)

			user, err = m.userService.RegisterUser(
				ctx, 
				cmd.User.TelegramID,
				cmd.User.Username,
				cmd.User.FirstName,
				cmd.User.LastName,
			)
			if err != nil {
				m.logger.Error("Failed to register user",
					"telegram_id", cmd.User.TelegramID,
					"error", err)
				return &domain.Response{
					Text:      "❌ Foydalanuvchini ro'yxatga olishda xatolik",
					ParseMode: "HTML",
				}, err
			}
		}

		// Update user activity
		err = m.userService.UpdateUserActivity(ctx, user.TelegramID)
		if err != nil {
			m.logger.Warn("Failed to update user activity",
				"telegram_id", user.TelegramID,
				"error", err)
		}

		// Add authenticated user to context
		ctx = domain.WithUser(ctx, user)

		// Update command with full user info
		cmd.User = user

		return next(ctx, cmd)
	}
}