package domain

import (
	"context"
	"time"
)

// User represents a bot user
type User struct {
	ID          int64     `json:"id"`
	TelegramID  int64     `json:"telegram_id"`
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Language    string    `json:"language"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Preferences map[string]interface{} `json:"preferences"`
}

// Chat represents a Telegram chat
type Chat struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Username string `json:"username"`
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, telegramID int64) error
	GetActiveUsers(ctx context.Context, limit int) ([]*User, error)
	GetUserStats(ctx context.Context) (*UserStats, error)
}

// UserStats represents user statistics
type UserStats struct {
	TotalUsers   int `json:"total_users"`
	ActiveUsers  int `json:"active_users"`
	NewToday     int `json:"new_today"`
	ActiveToday  int `json:"active_today"`
}

// UserService defines the interface for user business logic
type UserService interface {
	RegisterUser(ctx context.Context, telegramID int64, username, firstName, lastName string) (*User, error)
	GetUser(ctx context.Context, telegramID int64) (*User, error)
	UpdateUserActivity(ctx context.Context, telegramID int64) error
	GetStats(ctx context.Context) (*UserStats, error)
}