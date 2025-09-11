package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"yordamchi-dev-bot/database"
	"yordamchi-dev-bot/internal/domain"
)

// StatsCommand handles /stats command for user statistics
type StatsCommand struct {
	userService domain.UserService
	db          *database.DB
	startTime   time.Time
	logger      domain.Logger
}

// NewStatsCommand creates a new stats command handler
func NewStatsCommand(userService domain.UserService, db *database.DB, startTime time.Time, logger domain.Logger) *StatsCommand {
	return &StatsCommand{
		userService: userService,
		db:          db,
		startTime:   startTime,
		logger:      logger,
	}
}

// Handle processes the /stats command
func (h *StatsCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	// Get basic user statistics
	totalUsers, err := h.db.GetUserStats()
	if err != nil {
		h.logger.Error("Failed to get user stats", "error", err)
		return &domain.Response{
			Text:      "❌ Statistikani olishda xatolik yuz berdi",
			ParseMode: "MarkdownV2",
		}, nil
	}

	// Get daily statistics
	dailyStats, err := h.db.GetDailyStats()
	if err != nil {
		h.logger.Error("Failed to get daily stats", "error", err)
		dailyStats = make(map[string]int) // Continue with empty stats
	}

	// Get popular commands
	popularCommands, err := h.db.GetPopularCommands(5)
	if err != nil {
		h.logger.Error("Failed to get popular commands", "error", err)
		popularCommands = make(map[string]int)
	}

	uptime := time.Since(h.startTime)

	message := fmt.Sprintf(
		"📊 **Bot Statistikasi**\n\n"+
			"👥 **Foydalanuvchilar:**\n"+
			"   • Jami: %d\n"+
			"   • Bugun yangi: %d\n"+
			"   • Bugun faol: %d\n\n"+
			"📈 **Faollik:**\n"+
			"   • Bugun buyruqlar: %d\n\n"+
			"⏱️ **Uptime:** %s\n"+
			"🔄 **Arxitektura:** Clean Architecture\n"+
			"🚀 **Versiya:** 1.0.0",
		totalUsers,
		dailyStats["new_users_today"],
		dailyStats["active_users_today"],
		dailyStats["activities_today"],
		uptime.Truncate(time.Second).String(),
	)

	// Add popular commands if available
	if len(popularCommands) > 0 {
		message += "\n\n🔥 **Populyar buyruqlar:**\n"
		for cmd, count := range popularCommands {
			message += fmt.Sprintf("   • %s: %d\n", cmd, count)
		}
	}

	h.logger.Info("Stats command processed",
		"user_id", cmd.User.TelegramID,
		"total_users", totalUsers)

	return &domain.Response{
		Text:      message,
		ParseMode: "MarkdownV2",
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
