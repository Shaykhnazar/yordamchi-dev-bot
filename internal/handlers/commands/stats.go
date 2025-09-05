package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

// StatsCommand handles /stats command for user statistics
type StatsCommand struct {
	userService domain.UserService
	startTime   time.Time
	logger      domain.Logger
}

// NewStatsCommand creates a new stats command handler
func NewStatsCommand(userService domain.UserService, startTime time.Time, logger domain.Logger) *StatsCommand {
	return &StatsCommand{
		userService: userService,
		startTime:   startTime,
		logger:      logger,
	}
}

// Handle processes the /stats command
func (h *StatsCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	// Get user statistics
	stats, err := h.userService.GetStats(ctx)
	if err != nil {
		h.logger.Error("Failed to get user stats", "error", err)
		return &domain.Response{
			Text:      "❌ Statistikani olishda xatolik yuz berdi",
			ParseMode: "HTML",
		}, nil
	}

	uptime := time.Since(h.startTime)
	
	message := fmt.Sprintf(
		"📊 <b>Bot Statistikasi</b>\n\n"+
		"👥 <b>Foydalanuvchilar:</b>\n"+
		"   • Jami: %d\n"+
		"   • Faol: %d\n"+
		"   • Bugun yangi: %d\n"+
		"   • Bugun faol: %d\n\n"+
		"⏱️ <b>Uptime:</b> %s\n"+
		"🔄 <b>Arxitektura:</b> Clean Architecture\n"+
		"🚀 <b>Versiya:</b> 1.0.0",
		stats.TotalUsers,
		stats.ActiveUsers, 
		stats.NewToday,
		stats.ActiveToday,
		uptime.Truncate(time.Second).String(),
	)

	h.logger.Info("Stats command processed", 
		"user_id", cmd.User.TelegramID,
		"total_users", stats.TotalUsers)

	return &domain.Response{
		Text:      message,
		ParseMode: "HTML",
	}, nil
}

// CanHandle checks if this handler can process the command
func (h *StatsCommand) CanHandle(command string) bool {
	return strings.ToLower(strings.TrimSpace(command)) == "/stats"
}

// Description returns the command description
func (h *StatsCommand) Description() string {
	return "Foydalanuvchilar statistikasi"
}

// Usage returns the command usage
func (h *StatsCommand) Usage() string {
	return "/stats - Bot statistikasini ko'rish"
}