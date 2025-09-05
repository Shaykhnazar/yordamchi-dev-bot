package domain

import (
	"context"
	"time"
)

// Command represents a user command
type Command struct {
	ID        string
	Text      string
	User      *User
	Chat      *Chat
	Timestamp time.Time
}

// Response represents a bot response
type Response struct {
	Text           string
	ParseMode      string
	ReplyMarkup    interface{}
	DisablePreview bool
}

// CommandHandler defines the interface for command handling
type CommandHandler interface {
	Handle(ctx context.Context, cmd *Command) (*Response, error)
	CanHandle(command string) bool
	Description() string
	Usage() string
}

// Middleware defines the interface for processing pipeline
type Middleware interface {
	Process(ctx context.Context, next HandlerFunc) HandlerFunc
}

type HandlerFunc func(ctx context.Context, cmd *Command) (*Response, error)

// Router manages command routing and middleware
type Router interface {
	RegisterHandler(handler CommandHandler)
	RegisterMiddleware(middleware Middleware)
	Route(ctx context.Context, cmd *Command) (*Response, error)
	GetHandlers() []CommandHandler
}

// Logger interface for structured logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	With(args ...interface{}) Logger
}