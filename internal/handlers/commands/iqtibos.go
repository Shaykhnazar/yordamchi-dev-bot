package commands

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

// IqtibosCommand handles /iqtibos command for motivational quotes
type IqtibosCommand struct {
	quotes []string
	logger domain.Logger
}

// NewIqtibosCommand creates a new iqtibos command handler
func NewIqtibosCommand(quotes []string, logger domain.Logger) *IqtibosCommand {
	return &IqtibosCommand{
		quotes: quotes,
		logger: logger,
	}
}

// Handle processes the /iqtibos command
func (h *IqtibosCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	if len(h.quotes) == 0 {
		return &domain.Response{
			Text:      "💭 Hech qanday iqtibos topilmadi!",
			ParseMode: "Markdown",
		}, nil
	}

	// Get random quote
	rand.Seed(time.Now().UnixNano())
	randomQuote := h.quotes[rand.Intn(len(h.quotes))]

	h.logger.Info("Iqtibos command processed", 
		"user_id", cmd.User.TelegramID,
		"quote_length", len(randomQuote))

	return &domain.Response{
		Text:      randomQuote,
		ParseMode: "Markdown",
	}, nil
}

// CanHandle checks if this handler can process the command
func (h *IqtibosCommand) CanHandle(command string) bool {
	return strings.ToLower(strings.TrimSpace(command)) == "/iqtibos"
}

// Description returns the command description
func (h *IqtibosCommand) Description() string {
	return "Motivatsion iqtibos"
}

// Usage returns the command usage
func (h *IqtibosCommand) Usage() string {
	return "/iqtibos - Motivatsion iqtibos olish"
}