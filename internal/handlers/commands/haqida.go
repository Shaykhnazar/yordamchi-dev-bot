package commands

import (
	"context"
	"fmt"
	"strings"

	"yordamchi-dev-bot/handlers"
	"yordamchi-dev-bot/internal/domain"
)

// HaqidaCommand handles /haqida command for bot information
type HaqidaCommand struct {
	config *handlers.Config
	logger domain.Logger
}

// NewHaqidaCommand creates a new haqida command handler
func NewHaqidaCommand(config *handlers.Config, logger domain.Logger) *HaqidaCommand {
	return &HaqidaCommand{
		config: config,
		logger: logger,
	}
}

// Handle processes the /haqida command
func (h *HaqidaCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	botInfo := fmt.Sprintf(
		"🤖 **%s**\n\n"+
			"📊 **Versiya:** %s\n"+
			"📝 **Tavsif:** %s\n"+
			"👨‍💻 **Yaratuvchi:** %s\n\n"+
			"🏗️ **Arxitektura:** Clean Architecture\n"+
			"🚀 **Til:** Go (Golang)\n"+
			"📦 **Ma'lumotlar bazasi:** SQLite/PostgreSQL\n\n"+
			"💡 */help buyrug'i bilan barcha imkoniyatlarni ko'ring!*",
		h.config.Bot.Name,
		h.config.Bot.Version,
		h.config.Bot.Description,
		h.config.Bot.Author,
	)

	h.logger.Info("Haqida command processed", "user_id", cmd.User.TelegramID)

	return &domain.Response{
		Text:      botInfo,
		ParseMode: "MarkdownV2",
	}, nil
}

// CanHandle checks if this handler can process the command
func (h *HaqidaCommand) CanHandle(command string) bool {
	return strings.ToLower(strings.TrimSpace(command)) == "/haqida"
}

// Description returns the command description
func (h *HaqidaCommand) Description() string {
	return "Bot haqida ma'lumot"
}

// Usage returns the command usage
func (h *HaqidaCommand) Usage() string {
	return "/haqida - Bot haqida to'liq ma'lumot"
}
