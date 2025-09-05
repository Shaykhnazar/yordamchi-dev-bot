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
		"🤖 <b>%s</b>\n\n"+
		"📊 <b>Versiya:</b> %s\n"+
		"📝 <b>Tavsif:</b> %s\n"+
		"👨‍💻 <b>Yaratuvchi:</b> %s\n\n"+
		"🏗️ <b>Arxitektura:</b> Clean Architecture\n"+
		"🚀 <b>Til:</b> Go (Golang)\n"+
		"📦 <b>Ma'lumotlar bazasi:</b> SQLite/PostgreSQL\n\n"+
		"💡 <i>/help buyrug'i bilan barcha imkoniyatlarni ko'ring!</i>",
		h.config.Bot.Name,
		h.config.Bot.Version,
		h.config.Bot.Description,
		h.config.Bot.Author,
	)

	h.logger.Info("Haqida command processed", "user_id", cmd.User.TelegramID)

	return &domain.Response{
		Text:      botInfo,
		ParseMode: "HTML",
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