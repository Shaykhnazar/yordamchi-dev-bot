package commands

import (
	"context"
	"fmt"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

// PingCommand handles the /ping command for health checks
type PingCommand struct {
	logger    domain.Logger
	startTime time.Time
}

// NewPingCommand creates a new ping command handler
func NewPingCommand(logger domain.Logger, startTime time.Time) *PingCommand {
	return &PingCommand{
		logger:    logger,
		startTime: startTime,
	}
}

// Handle processes the ping command
func (h *PingCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	uptime := time.Since(h.startTime)

	response := fmt.Sprintf(`🏓 **Pong!**
	
✅ Bot ishlayapti
⏱ Uptime: %s
🕐 Server vaqti: %s
👤 Foydalanuvchi: %s (@%s)`,
		formatUptime(uptime),
		time.Now().Format("2006-01-02 15:04:05"),
		cmd.User.FirstName,
		cmd.User.Username)

	h.logger.Info("Ping command processed",
		"user_id", cmd.User.TelegramID,
		"uptime", uptime.String())

	return &domain.Response{
		Text:      response,
		ParseMode: "MarkdownV2",
	}, nil
}

// CanHandle checks if this handler can process the command
func (h *PingCommand) CanHandle(command string) bool {
	return command == "/ping"
}

// Description returns the command description
func (h *PingCommand) Description() string {
	return "Ping command - health check and uptime info"
}

// Usage returns the command usage instructions
func (h *PingCommand) Usage() string {
	return "/ping - Bot ishlaganligini tekshirish"
}

// formatUptime formats duration into human readable format
func formatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%d kun, %d soat, %d daqiqa", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%d soat, %d daqiqa, %d soniya", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%d daqiqa, %d soniya", minutes, seconds)
	} else {
		return fmt.Sprintf("%d soniya", seconds)
	}
}
