package commands

import (
	"context"

	"yordamchi-dev-bot/internal/domain"
)

// HelpCommand handles the /help command
type HelpCommand struct {
	router     domain.Router
	logger     domain.Logger
	staticHelp string
}

// NewHelpCommand creates a new help command handler
func NewHelpCommand(router domain.Router, staticHelp string, logger domain.Logger) *HelpCommand {
	return &HelpCommand{
		router:     router,
		staticHelp: staticHelp,
		logger:     logger,
	}
}

// Handle processes the help command
func (h *HelpCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	// Use static help message from config for now
	// In the future, this could dynamically generate from registered handlers
	helpMessage := h.staticHelp

	if helpMessage == "" {
		// Fallback to dynamic help
		helpMessage = h.generateDynamicHelp()
	}

	h.logger.Info("Help command processed", "user_id", cmd.User.TelegramID)

	return &domain.Response{
		Text:      helpMessage,
		ParseMode: "MarkdownV2",
	}, nil
}

// CanHandle checks if this handler can process the command
func (h *HelpCommand) CanHandle(command string) bool {
	return command == "/help"
}

// Description returns the command description
func (h *HelpCommand) Description() string {
	return "Help command - shows available commands"
}

// Usage returns the command usage instructions
func (h *HelpCommand) Usage() string {
	return "/help - Bu yordam xabari"
}

// generateDynamicHelp creates help message from registered handlers
func (h *HelpCommand) generateDynamicHelp() string {
	handlers := h.router.GetHandlers()

	helpText := "🤖 Mavjud buyruqlar:\n\n"
	for _, handler := range handlers {
		if handler.Usage() != "" {
			helpText += handler.Usage() + "\n"
		}
	}

	if len(handlers) == 0 {
		helpText = "Hech qanday buyruq mavjud emas"
	}

	return helpText
}
