package domain

import "context"

// Context keys for request context
type contextKey string

const (
	UserContextKey    contextKey = "user"
	CommandContextKey contextKey = "command"
	LoggerContextKey  contextKey = "logger"
)

// GetUserFromContext extracts user from context
func GetUserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(UserContextKey).(*User)
	return user, ok
}

// WithUser adds user to context
func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

// GetCommandFromContext extracts command from context
func GetCommandFromContext(ctx context.Context) (*Command, bool) {
	cmd, ok := ctx.Value(CommandContextKey).(*Command)
	return cmd, ok
}

// WithCommand adds command to context
func WithCommand(ctx context.Context, cmd *Command) context.Context {
	return context.WithValue(ctx, CommandContextKey, cmd)
}

// GetLoggerFromContext extracts logger from context
func GetLoggerFromContext(ctx context.Context) (Logger, bool) {
	logger, ok := ctx.Value(LoggerContextKey).(Logger)
	return logger, ok
}

// WithLogger adds logger to context
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, LoggerContextKey, logger)
}