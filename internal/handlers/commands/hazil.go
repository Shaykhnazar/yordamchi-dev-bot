package commands

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

// HazilCommand handles /hazil command for random programming jokes
type HazilCommand struct {
	jokes  []string
	logger domain.Logger
}

// NewHazilCommand creates a new hazil command handler
func NewHazilCommand(jokes []string, logger domain.Logger) *HazilCommand {
	return &HazilCommand{
		jokes:  jokes,
		logger: logger,
	}
}

// Handle processes the /hazil command
func (h *HazilCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	if len(h.jokes) == 0 {
		return &domain.Response{
			Text:      "😅 Hech qanday hazil topilmadi!",
			ParseMode: "Markdown",
		}, nil
	}

	// Get random joke
	rand.Seed(time.Now().UnixNano())
	randomJoke := h.jokes[rand.Intn(len(h.jokes))]

	h.logger.Info("Hazil command processed",
		"user_id", cmd.User.TelegramID,
		"joke_length", len(randomJoke))

	return &domain.Response{
		Text:      randomJoke,
		ParseMode: "Markdown",
	}, nil
}

// CanHandle checks if this handler can process the command
func (h *HazilCommand) CanHandle(command string) bool {
	return strings.ToLower(strings.TrimSpace(command)) == "/hazil"
}

// Description returns the command description
func (h *HazilCommand) Description() string {
	return "Tasodifiy dasturlash hazili"
}

// Usage returns the command usage
func (h *HazilCommand) Usage() string {
	return "/hazil - Tasodifiy hazil olish"
}
