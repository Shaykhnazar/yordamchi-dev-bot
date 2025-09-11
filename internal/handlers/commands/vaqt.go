package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

// VaqtCommand handles /vaqt command for current timestamp
type VaqtCommand struct {
	logger domain.Logger
}

// NewVaqtCommand creates a new vaqt command handler
func NewVaqtCommand(logger domain.Logger) *VaqtCommand {
	return &VaqtCommand{
		logger: logger,
	}
}

// Handle processes the /vaqt command
func (h *VaqtCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	now := time.Now()
	
	// Format time in a user-friendly way
	timeInfo := fmt.Sprintf(
		"🕐 **Hozirgi vaqt:**\n\n"+
		"📅 **Sana:** %s\n"+
		"⏰ **Vaqt:** %s\n"+
		"🌍 **UTC:** %s\n"+
		"📊 **Unix timestamp:** %d",
		now.Format("2006-01-02"),
		now.Format("15:04:05"),
		now.UTC().Format("2006-01-02 15:04:05"),
		now.Unix(),
	)

	h.logger.Info("Vaqt command processed", 
		"user_id", cmd.User.TelegramID,
		"timestamp", now.Unix())

	return &domain.Response{
		Text:      timeInfo,
		ParseMode: "Markdown",
	}, nil
}

// CanHandle checks if this handler can process the command
func (h *VaqtCommand) CanHandle(command string) bool {
	return strings.ToLower(strings.TrimSpace(command)) == "/vaqt"
}

// Description returns the command description
func (h *VaqtCommand) Description() string {
	return "Hozirgi vaqt va sana"
}

// Usage returns the command usage
func (h *VaqtCommand) Usage() string {
	return "/vaqt - Hozirgi vaqtni ko'rish"
}