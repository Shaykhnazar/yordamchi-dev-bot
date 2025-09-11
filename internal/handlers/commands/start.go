package commands

import (
	"context"

	"yordamchi-dev-bot/internal/domain"
)

// StartCommand handles the /start command
type StartCommand struct {
	welcomeMessage string
	logger         domain.Logger
}

// NewStartCommand creates a new start command handler
func NewStartCommand(welcomeMessage string, logger domain.Logger) *StartCommand {
	return &StartCommand{
		welcomeMessage: welcomeMessage,
		logger:         logger,
	}
}

// Handle processes the start command
func (h *StartCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	user, _ := domain.GetUserFromContext(ctx)

	message := h.welcomeMessage
	if user != nil && user.FirstName != "" {
		message += "\n\n👋 Salom, " + user.FirstName + "!"
	}
	message += "\n\n/help - barcha buyruqlar ro'yxati"

	h.logger.Info("Start command processed", "user_id", cmd.User.TelegramID)

	return &domain.Response{
		Text:      message,
		ParseMode: "MarkdownV2",
	}, nil
}

// CanHandle checks if this handler can process the command
func (h *StartCommand) CanHandle(command string) bool {
	return command == "/start"
}

// Description returns the command description
func (h *StartCommand) Description() string {
	return "Start command - welcomes users and provides introduction"
}

// Usage returns the command usage instructions
func (h *StartCommand) Usage() string {
	return "/start - Botni ishga tushirish"
}
