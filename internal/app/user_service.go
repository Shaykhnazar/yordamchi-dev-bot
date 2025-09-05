package app

import (
	"context"
	"time"

	"yordamchi-dev-bot/database"
	"yordamchi-dev-bot/internal/domain"
)

// UserService implements the domain.UserService interface
type UserService struct {
	db     *database.DB
	logger domain.Logger
}

// NewUserService creates a new user service
func NewUserService(db *database.DB, logger domain.Logger) *UserService {
	return &UserService{
		db:     db,
		logger: logger,
	}
}

// RegisterUser registers a new user or returns existing user
func (s *UserService) RegisterUser(ctx context.Context, telegramID int64, username, firstName, lastName string) (*domain.User, error) {
	// Try to create or update user in database
	err := s.db.CreateOrUpdateUser(telegramID, username, firstName, lastName)
	if err != nil {
		s.logger.Error("Failed to register user", "telegram_id", telegramID, "error", err)
		return nil, err
	}

	// Return the user object
	user := &domain.User{
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
		Language:   "uz", // Default language
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Preferences: make(map[string]interface{}),
	}

	s.logger.Info("User registered successfully", "telegram_id", telegramID, "username", username)
	return user, nil
}

// GetUser retrieves a user by Telegram ID
func (s *UserService) GetUser(ctx context.Context, telegramID int64) (*domain.User, error) {
	// For now, we'll create a user object from the basic info we have
	// In a full implementation, this would fetch from database
	user := &domain.User{
		TelegramID: telegramID,
		Language:   "uz",
		IsActive:   true,
		UpdatedAt:  time.Now(),
		Preferences: make(map[string]interface{}),
	}

	return user, nil
}

// UpdateUserActivity updates user's last activity timestamp
func (s *UserService) UpdateUserActivity(ctx context.Context, telegramID int64) error {
	// This would update last_seen timestamp in database
	// For now, we'll just log it
	s.logger.Debug("User activity updated", "telegram_id", telegramID)
	return nil
}

// GetStats returns user statistics
func (s *UserService) GetStats(ctx context.Context) (*domain.UserStats, error) {
	count, err := s.db.GetUserStats()
	if err != nil {
		return nil, err
	}

	stats := &domain.UserStats{
		TotalUsers:  count,
		ActiveUsers: count, // Simplified - in reality would be different
		NewToday:    0,     // Would need to implement
		ActiveToday: 0,     // Would need to implement
	}

	return stats, nil
}