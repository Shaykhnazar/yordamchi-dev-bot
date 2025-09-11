package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

// MetricsProvider interface for getting metrics from middleware
type MetricsProvider interface {
	GetMetrics() map[string]interface{}
	GetCacheStats() map[string]interface{}
}

// MetricsCommand handles /metrics command for performance monitoring
type MetricsCommand struct {
	metricsProvider MetricsProvider
	logger          domain.Logger
}

// NewMetricsCommand creates a new metrics command handler
func NewMetricsCommand(metricsProvider MetricsProvider, logger domain.Logger) *MetricsCommand {
	return &MetricsCommand{
		metricsProvider: metricsProvider,
		logger:          logger,
	}
}

// Handle processes the /metrics command
func (h *MetricsCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	// Get performance metrics
	metrics := h.metricsProvider.GetMetrics()

	// Get cache metrics
	cacheStats := h.metricsProvider.GetCacheStats()

	message := h.formatMetricsMessage(metrics, cacheStats)

	h.logger.Info("Metrics command processed",
		"user_id", cmd.User.TelegramID)

	return &domain.Response{
		Text:      message,
		ParseMode: "Markdown",
	}, nil
}

// formatMetricsMessage formats metrics data into readable message
func (h *MetricsCommand) formatMetricsMessage(metrics, cacheStats map[string]interface{}) string {
	var message strings.Builder

	message.WriteString("📈 **Bot Performance Metrics**\n\n")

	// System metrics
	message.WriteString("🖥️ **System:**\n")
	if uptime, ok := metrics["uptime_seconds"].(int64); ok {
		message.WriteString(fmt.Sprintf("   • Uptime: %s\n", formatDuration(time.Duration(uptime)*time.Second)))
	}

	// Request metrics
	message.WriteString("\n📊 **Requests:**\n")
	if total, ok := metrics["total_requests"].(int64); ok {
		message.WriteString(fmt.Sprintf("   • Total: %d\n", total))
	}
	if successful, ok := metrics["successful_requests"].(int64); ok {
		message.WriteString(fmt.Sprintf("   • Successful: %d\n", successful))
	}
	if failed, ok := metrics["failed_requests"].(int64); ok {
		message.WriteString(fmt.Sprintf("   • Failed: %d\n", failed))
	}
	if rate, ok := metrics["success_rate_percent"].(float64); ok {
		message.WriteString(fmt.Sprintf("   • Success Rate: %.1f%%\n", rate))
	}
	if rpm, ok := metrics["requests_per_minute"].(float64); ok {
		message.WriteString(fmt.Sprintf("   • Req/min: %.1f\n", rpm))
	}

	// Top commands
	if topCommands, ok := metrics["top_commands"].([]map[string]interface{}); ok && len(topCommands) > 0 {
		message.WriteString("\n🔥 **Popular Commands:**\n")
		for i, cmdData := range topCommands {
			if i >= 5 { // Limit to top 5
				break
			}
			if cmd, ok := cmdData["command"].(string); ok {
				if count, ok := cmdData["count"].(int64); ok {
					message.WriteString(fmt.Sprintf("   %d. %s: %d\n", i+1, cmd, count))
				}
			}
		}
	}

	// Cache metrics
	if cacheStats != nil {
		message.WriteString("\n💾 **Cache:**\n")
		if size, ok := cacheStats["cache_size"].(int); ok {
			message.WriteString(fmt.Sprintf("   • Size: %d items\n", size))
		}
		if ttl, ok := cacheStats["ttl_minutes"].(int); ok {
			message.WriteString(fmt.Sprintf("   • TTL: %d minutes\n", ttl))
		}
		if commands, ok := cacheStats["cacheable_commands"].(int); ok {
			message.WriteString(fmt.Sprintf("   • Cached Commands: %d\n", commands))
		}
	}

	// Performance indicators
	message.WriteString("\n⚡ **Performance:**\n")
	if cmdMetrics, ok := metrics["command_metrics"].(map[string]interface{}); ok {
		avgDuration := h.calculateAverageResponseTime(cmdMetrics)
		message.WriteString(fmt.Sprintf("   • Avg Response: %dms\n", avgDuration))

		slowestCommand := h.findSlowestCommand(cmdMetrics)
		if slowestCommand.Command != "" {
			message.WriteString(fmt.Sprintf("   • Slowest: %s (%dms)\n",
				slowestCommand.Command, slowestCommand.Duration))
		}
	}

	message.WriteString("\n🤖 *Real-time performance monitoring*")

	return message.String()
}

type CommandPerformance struct {
	Command  string
	Duration int64
}

// calculateAverageResponseTime calculates overall average response time
func (h *MetricsCommand) calculateAverageResponseTime(cmdMetrics map[string]interface{}) int64 {
	var totalDuration int64
	var totalCommands int64

	for _, metricData := range cmdMetrics {
		if metrics, ok := metricData.(map[string]interface{}); ok {
			if avgDuration, ok := metrics["average_duration"].(int64); ok {
				if count, ok := metrics["count"].(int64); ok {
					totalDuration += avgDuration * count
					totalCommands += count
				}
			}
		}
	}

	if totalCommands == 0 {
		return 0
	}

	return totalDuration / totalCommands
}

// findSlowestCommand finds the command with highest average response time
func (h *MetricsCommand) findSlowestCommand(cmdMetrics map[string]interface{}) CommandPerformance {
	var slowest CommandPerformance

	for command, metricData := range cmdMetrics {
		if metrics, ok := metricData.(map[string]interface{}); ok {
			if avgDuration, ok := metrics["average_duration"].(int64); ok {
				if avgDuration > slowest.Duration {
					slowest.Command = command
					slowest.Duration = avgDuration
				}
			}
		}
	}

	return slowest
}

// formatDuration formats time duration into readable string
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	} else {
		return fmt.Sprintf("%.1fd", d.Hours()/24)
	}
}

// CanHandle checks if this handler can process the command
func (h *MetricsCommand) CanHandle(command string) bool {
	return strings.ToLower(strings.TrimSpace(command)) == "/metrics"
}

// Description returns the command description
func (h *MetricsCommand) Description() string {
	return "Bot performance metrics"
}

// Usage returns the command usage
func (h *MetricsCommand) Usage() string {
	return "/metrics - Bot performance va statistikasini ko'rish"
}
