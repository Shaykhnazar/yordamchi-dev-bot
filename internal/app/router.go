package app

import (
	"context"
	"strings"

	"yordamchi-dev-bot/internal/domain"
)

// CommandRouter implements the Router interface
type CommandRouter struct {
	handlers    []domain.CommandHandler
	middlewares []domain.Middleware
	logger      domain.Logger
}

// NewCommandRouter creates a new command router
func NewCommandRouter(logger domain.Logger) *CommandRouter {
	return &CommandRouter{
		handlers:    make([]domain.CommandHandler, 0),
		middlewares: make([]domain.Middleware, 0),
		logger:      logger,
	}
}

// RegisterHandler registers a new command handler
func (r *CommandRouter) RegisterHandler(handler domain.CommandHandler) {
	r.handlers = append(r.handlers, handler)
	r.logger.Info("Command handler registered", "handler", handler.Description())
}

// RegisterMiddleware registers a new middleware
func (r *CommandRouter) RegisterMiddleware(middleware domain.Middleware) {
	r.middlewares = append(r.middlewares, middleware)
	r.logger.Info("Middleware registered")
}

// Route finds and executes the appropriate handler for a command
func (r *CommandRouter) Route(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	// Add command to context
	ctx = domain.WithCommand(ctx, cmd)

	// Find appropriate handler
	var handler domain.CommandHandler
	// Extract just the command part (first word) for matching
	parts := strings.Fields(cmd.Text)
	if len(parts) == 0 {
		return &domain.Response{
			Text:      "❓ Noma'lum buyruq. /help yozing",
			ParseMode: "Markdown",
		}, nil
	}
	command := parts[0]
	for _, h := range r.handlers {
		if h.CanHandle(command) {
			handler = h
			break
		}
	}

	if handler == nil {
		return &domain.Response{
			Text:      "❓ Noma'lum buyruq. /help yozing",
			ParseMode: "Markdown",
		}, nil
	}

	// Build middleware chain
	handlerFunc := r.buildMiddlewareChain(handler.Handle)

	// Execute with middleware chain
	response, err := handlerFunc(ctx, cmd)
	if err != nil {
		r.logger.Error("Command execution failed", "command", cmd.Text, "error", err)
		return &domain.Response{
			Text:      "❌ Buyruqni bajarishda xatolik yuz berdi",
			ParseMode: "Markdown",
		}, err
	}

	r.logger.Info("Command executed successfully",
		"command", cmd.Text,
		"user", cmd.User.TelegramID,
		"handler", handler.Description())

	return response, nil
}

// GetHandlers returns all registered handlers
func (r *CommandRouter) GetHandlers() []domain.CommandHandler {
	return r.handlers
}

// buildMiddlewareChain builds the middleware execution chain
func (r *CommandRouter) buildMiddlewareChain(handler domain.HandlerFunc) domain.HandlerFunc {
	// Start with the actual handler
	next := handler

	// Apply middlewares in reverse order (last registered, first executed)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		next = r.middlewares[i].Process(context.Background(), next)
	}

	return next
}

// GetAvailableCommands returns a formatted list of available commands
func (r *CommandRouter) GetAvailableCommands() string {
	var commands []string

	for _, handler := range r.handlers {
		usage := handler.Usage()
		if usage != "" {
			commands = append(commands, usage)
		}
	}

	if len(commands) == 0 {
		return "Hech qanday buyruq mavjud emas"
	}

	return "Mavjud buyruqlar:\n" + strings.Join(commands, "\n")
}
