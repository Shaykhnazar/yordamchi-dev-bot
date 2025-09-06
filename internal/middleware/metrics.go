package middleware

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

// MetricsMiddleware collects performance and usage metrics
type MetricsMiddleware struct {
	logger           domain.Logger
	totalRequests    int64
	successfulRequests int64
	failedRequests   int64
	commandMetrics   map[string]*CommandMetrics
	mutex           sync.RWMutex
	startTime       time.Time
}

// CommandMetrics tracks metrics for individual commands
type CommandMetrics struct {
	Count            int64
	TotalDuration    time.Duration
	AverageDuration  time.Duration
	LastUsed         time.Time
	ErrorCount       int64
	mutex           sync.RWMutex
}

// NewMetricsMiddleware creates a new metrics collection middleware
func NewMetricsMiddleware(logger domain.Logger) *MetricsMiddleware {
	return &MetricsMiddleware{
		logger:         logger,
		commandMetrics: make(map[string]*CommandMetrics),
		startTime:      time.Now(),
	}
}

// Process implements the Middleware interface
func (m *MetricsMiddleware) Process(ctx context.Context, next domain.HandlerFunc) domain.HandlerFunc {
	return func(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
		startTime := time.Now()
		
		// Increment total requests
		atomic.AddInt64(&m.totalRequests, 1)

		// Execute command
		response, err := next(ctx, cmd)
		
		// Calculate execution time
		duration := time.Since(startTime)
		
		// Update metrics
		m.updateCommandMetrics(cmd.Text, duration, err)
		
		// Update success/failure counters
		if err != nil {
			atomic.AddInt64(&m.failedRequests, 1)
		} else {
			atomic.AddInt64(&m.successfulRequests, 1)
		}

		// Log performance metrics for slow commands
		if duration > 2*time.Second {
			m.logger.Warn("Slow command execution",
				"command", cmd.Text,
				"user_id", cmd.User.TelegramID,
				"duration", duration,
				"threshold", "2s")
		}

		return response, err
	}
}

// updateCommandMetrics updates metrics for a specific command
func (m *MetricsMiddleware) updateCommandMetrics(command string, duration time.Duration, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.commandMetrics[command] == nil {
		m.commandMetrics[command] = &CommandMetrics{}
	}

	metrics := m.commandMetrics[command]
	metrics.mutex.Lock()
	defer metrics.mutex.Unlock()

	// Update counters
	metrics.Count++
	metrics.LastUsed = time.Now()
	metrics.TotalDuration += duration

	// Calculate average duration
	metrics.AverageDuration = time.Duration(int64(metrics.TotalDuration) / metrics.Count)

	// Update error count
	if err != nil {
		metrics.ErrorCount++
	}
}

// GetMetrics returns comprehensive metrics data
func (m *MetricsMiddleware) GetMetrics() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	totalRequests := atomic.LoadInt64(&m.totalRequests)
	successfulRequests := atomic.LoadInt64(&m.successfulRequests)
	failedRequests := atomic.LoadInt64(&m.failedRequests)

	uptime := time.Since(m.startTime)
	
	var successRate float64
	if totalRequests > 0 {
		successRate = float64(successfulRequests) / float64(totalRequests) * 100
	}

	metrics := map[string]interface{}{
		"uptime_seconds":       int64(uptime.Seconds()),
		"total_requests":       totalRequests,
		"successful_requests":  successfulRequests,
		"failed_requests":      failedRequests,
		"success_rate_percent": successRate,
		"requests_per_minute":  m.getRequestsPerMinute(totalRequests, uptime),
		"command_metrics":      m.getCommandMetricsData(),
		"top_commands":         m.getTopCommands(5),
	}

	return metrics
}

// getRequestsPerMinute calculates requests per minute
func (m *MetricsMiddleware) getRequestsPerMinute(totalRequests int64, uptime time.Duration) float64 {
	minutes := uptime.Minutes()
	if minutes == 0 {
		return 0
	}
	return float64(totalRequests) / minutes
}

// getCommandMetricsData returns formatted command metrics
func (m *MetricsMiddleware) getCommandMetricsData() map[string]interface{} {
	cmdMetrics := make(map[string]interface{})
	
	for cmd, metrics := range m.commandMetrics {
		metrics.mutex.RLock()
		cmdMetrics[cmd] = map[string]interface{}{
			"count":             metrics.Count,
			"average_duration":  metrics.AverageDuration.Milliseconds(),
			"total_duration":    metrics.TotalDuration.Milliseconds(),
			"error_count":       metrics.ErrorCount,
			"last_used":         metrics.LastUsed.Format(time.RFC3339),
			"error_rate":        m.calculateErrorRate(metrics),
		}
		metrics.mutex.RUnlock()
	}
	
	return cmdMetrics
}

// getTopCommands returns most frequently used commands
func (m *MetricsMiddleware) getTopCommands(limit int) []map[string]interface{} {
	type CommandCount struct {
		Command string
		Count   int64
	}

	var commands []CommandCount
	for cmd, metrics := range m.commandMetrics {
		metrics.mutex.RLock()
		commands = append(commands, CommandCount{
			Command: cmd,
			Count:   metrics.Count,
		})
		metrics.mutex.RUnlock()
	}

	// Simple bubble sort for top commands
	for i := 0; i < len(commands)-1; i++ {
		for j := 0; j < len(commands)-i-1; j++ {
			if commands[j].Count < commands[j+1].Count {
				commands[j], commands[j+1] = commands[j+1], commands[j]
			}
		}
	}

	var result []map[string]interface{}
	maxItems := limit
	if len(commands) < limit {
		maxItems = len(commands)
	}

	for i := 0; i < maxItems; i++ {
		result = append(result, map[string]interface{}{
			"command": commands[i].Command,
			"count":   commands[i].Count,
		})
	}

	return result
}

// calculateErrorRate calculates error rate for a command
func (m *MetricsMiddleware) calculateErrorRate(metrics *CommandMetrics) float64 {
	if metrics.Count == 0 {
		return 0
	}
	return float64(metrics.ErrorCount) / float64(metrics.Count) * 100
}

// ResetMetrics resets all metrics (useful for testing)
func (m *MetricsMiddleware) ResetMetrics() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	atomic.StoreInt64(&m.totalRequests, 0)
	atomic.StoreInt64(&m.successfulRequests, 0)
	atomic.StoreInt64(&m.failedRequests, 0)
	m.commandMetrics = make(map[string]*CommandMetrics)
	m.startTime = time.Now()

	m.logger.Info("Metrics reset")
}