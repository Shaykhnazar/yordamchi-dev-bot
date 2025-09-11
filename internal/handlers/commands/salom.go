package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

// SalomCommand handles /salom command for personalized greetings
type SalomCommand struct {
	logger domain.Logger
}

// NewSalomCommand creates a new salom command handler
func NewSalomCommand(logger domain.Logger) *SalomCommand {
	return &SalomCommand{
		logger: logger,
	}
}

// Handle processes the /salom command
func (h *SalomCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	user := cmd.User
	now := time.Now()
	hour := now.Hour()

	// Determine greeting based on time of day
	var greeting string
	switch {
	case hour < 6:
		greeting = "🌙 Kech tunda ishlayapsizmi"
	case hour < 12:
		greeting = "🌅 Xayrli tong"
	case hour < 17:
		greeting = "☀️ Xayrli kun"
	case hour < 21:
		greeting = "🌇 Xayrli kech"
	default:
		greeting = "🌃 Xayrli tun"
	}

	// Create personalized message
	name := user.FirstName
	if name == "" {
		name = user.Username
	}
	if name == "" {
		name = "Do'st"
	}

	message := fmt.Sprintf(
		"%s, %s! 👋\n\n"+
			"🤖 Men sizning yordamchi botingizman.\n"+
			"📚 Dasturlashni o'rganishda yordam beraman!\n\n"+
			"💡 */help buyrug'i bilan nima qila olishimni bilib oling.*",
		greeting, name,
	)

	h.logger.Info("Salom command processed",
		"user_id", user.TelegramID,
		"username", user.Username,
		"first_name", user.FirstName)

	return &domain.Response{
		Text:      message,
		ParseMode: "Markdown",
	}, nil
}

// CanHandle checks if this handler can process the command
func (h *SalomCommand) CanHandle(command string) bool {
	return strings.ToLower(strings.TrimSpace(command)) == "/salom"
}

// Description returns the command description
func (h *SalomCommand) Description() string {
	return "Shaxsiylashtirilgan salom"
}

// Usage returns the command usage
func (h *SalomCommand) Usage() string {
	return "/salom - Shaxsiylashtirilgan salom olish"
}
