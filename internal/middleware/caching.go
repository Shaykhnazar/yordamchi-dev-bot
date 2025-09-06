package middleware

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"yordamchi-dev-bot/internal/cache"
	"yordamchi-dev-bot/internal/domain"
)

// CachingMiddleware provides response caching for expensive operations
type CachingMiddleware struct {
	cache        *cache.MemoryCache
	logger       domain.Logger
	cacheTTL     time.Duration
	cacheableCommands map[string]bool
}

// NewCachingMiddleware creates a new caching middleware
func NewCachingMiddleware(logger domain.Logger) *CachingMiddleware {
	// Commands that should be cached (expensive operations)
	cacheableCommands := map[string]bool{
		"/weather":   true,
		"/repo":      true,
		"/user":      true,
	}

	return &CachingMiddleware{
		cache:             cache.NewMemoryCache(10 * time.Minute), // 10 minute default TTL
		logger:            logger,
		cacheTTL:          10 * time.Minute,
		cacheableCommands: cacheableCommands,
	}
}

// Process implements the Middleware interface
func (m *CachingMiddleware) Process(ctx context.Context, next domain.HandlerFunc) domain.HandlerFunc {
	return func(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
		// Check if command should be cached
		commandParts := strings.Fields(strings.ToLower(cmd.Text))
		if len(commandParts) == 0 {
			return next(ctx, cmd)
		}

		baseCommand := commandParts[0]
		if !m.cacheableCommands[baseCommand] {
			return next(ctx, cmd)
		}

		// Generate cache key from user ID and full command text
		cacheKey := m.generateCacheKey(cmd.User.TelegramID, cmd.Text)

		// Try to get from cache first
		if cachedResponse, found := m.cache.Get(cacheKey); found {
			if response, ok := cachedResponse.(*domain.Response); ok {
				m.logger.Debug("Cache hit", 
					"command", cmd.Text, 
					"user_id", cmd.User.TelegramID,
					"cache_key", cacheKey)

				// Add cache indicator to response
				response.Text = "🔄 " + response.Text
				return response, nil
			}
		}

		// Execute command
		response, err := next(ctx, cmd)
		if err != nil {
			return response, err
		}

		// Cache successful responses
		if response != nil && response.Text != "" {
			// Set cache TTL based on command type
			ttl := m.getCacheTTL(baseCommand)
			m.cache.SetWithTTL(cacheKey, response, ttl)

			m.logger.Debug("Response cached", 
				"command", cmd.Text,
				"user_id", cmd.User.TelegramID,
				"cache_key", cacheKey,
				"ttl", ttl)
		}

		return response, err
	}
}

// generateCacheKey creates a unique cache key for the request
func (m *CachingMiddleware) generateCacheKey(userID int64, command string) string {
	// Create hash of user + command for cache key
	data := fmt.Sprintf("user:%d:cmd:%s", userID, strings.ToLower(command))
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("cache:%x", hash)
}

// getCacheTTL returns appropriate TTL for different command types
func (m *CachingMiddleware) getCacheTTL(command string) time.Duration {
	switch command {
	case "/weather":
		return 15 * time.Minute // Weather changes more frequently
	case "/repo", "/user":
		return 30 * time.Minute // GitHub data changes less frequently
	default:
		return m.cacheTTL
	}
}

// GetCacheStats returns cache statistics for monitoring
func (m *CachingMiddleware) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"cache_size": m.cache.Size(),
		"ttl_minutes": int(m.cacheTTL.Minutes()),
		"cacheable_commands": len(m.cacheableCommands),
	}
}

// ClearCache clears all cached responses
func (m *CachingMiddleware) ClearCache() {
	m.cache.Clear()
	m.logger.Info("Cache cleared")
}