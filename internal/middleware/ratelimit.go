package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

// RateLimitMiddleware provides rate limiting per user
type RateLimitMiddleware struct {
	limits   map[int64]*UserLimit
	mutex    sync.RWMutex
	logger   domain.Logger
	maxReqs  int
	window   time.Duration
}

// UserLimit tracks rate limiting for a specific user
type UserLimit struct {
	requests []time.Time
	mutex    sync.Mutex
}

// NewRateLimitMiddleware creates a new rate limiting middleware
func NewRateLimitMiddleware(maxRequests int, window time.Duration, logger domain.Logger) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limits:  make(map[int64]*UserLimit),
		logger:  logger,
		maxReqs: maxRequests,
		window:  window,
	}
}

// Process implements the Middleware interface
func (m *RateLimitMiddleware) Process(ctx context.Context, next domain.HandlerFunc) domain.HandlerFunc {
	return func(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
		userID := cmd.User.TelegramID

		// Check if user is rate limited
		if m.isRateLimited(userID) {
			m.logger.Warn("User rate limited",
				"user_id", userID,
				"username", cmd.User.Username,
				"command", cmd.Text)

			return &domain.Response{
				Text:      fmt.Sprintf("⚠️ Juda ko'p so'rov! %d soniyadan keyin qayta urinib ko'ring.", int(m.window.Seconds())),
				ParseMode: "HTML",
			}, nil
		}

		// Record this request
		m.recordRequest(userID)

		// Continue to next handler
		return next(ctx, cmd)
	}
}

// isRateLimited checks if user has exceeded rate limit
func (m *RateLimitMiddleware) isRateLimited(userID int64) bool {
	m.mutex.RLock()
	limit, exists := m.limits[userID]
	m.mutex.RUnlock()

	if !exists {
		return false
	}

	limit.mutex.Lock()
	defer limit.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-m.window)

	// Remove old requests
	validRequests := make([]time.Time, 0)
	for _, reqTime := range limit.requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	limit.requests = validRequests

	// Check if limit exceeded
	return len(limit.requests) >= m.maxReqs
}

// recordRequest records a new request for the user
func (m *RateLimitMiddleware) recordRequest(userID int64) {
	m.mutex.Lock()
	limit, exists := m.limits[userID]
	if !exists {
		limit = &UserLimit{
			requests: make([]time.Time, 0),
		}
		m.limits[userID] = limit
	}
	m.mutex.Unlock()

	limit.mutex.Lock()
	limit.requests = append(limit.requests, time.Now())
	limit.mutex.Unlock()
}

// Cleanup removes old rate limit data (should be called periodically)
func (m *RateLimitMiddleware) Cleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-m.window * 2) // Keep data for 2x window duration

	for userID, limit := range m.limits {
		limit.mutex.Lock()
		hasRecentRequests := false
		for _, reqTime := range limit.requests {
			if reqTime.After(cutoff) {
				hasRecentRequests = true
				break
			}
		}
		
		if !hasRecentRequests {
			delete(m.limits, userID)
		}
		limit.mutex.Unlock()
	}

	m.logger.Info("Rate limit cleanup completed", "remaining_users", len(m.limits))
}