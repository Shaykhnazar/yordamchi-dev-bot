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
	// File attachments
	Document *TelegramDocument `json:"document,omitempty"`
	Photo    []TelegramPhoto   `json:"photo,omitempty"`
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

// Telegram file types for file upload support

// TelegramDocument represents a document file sent via Telegram
type TelegramDocument struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileName     string `json:"file_name,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int    `json:"file_size,omitempty"`
	Thumbnail    *TelegramPhotoSize `json:"thumb,omitempty"`
}

// TelegramPhoto represents a photo sent via Telegram
type TelegramPhoto struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size,omitempty"`
}

// TelegramPhotoSize represents different sizes of photos/thumbnails
type TelegramPhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size,omitempty"`
}

// TelegramFile represents file info from getFile API
type TelegramFile struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}