package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"yordamchi-dev-bot/internal/domain"
)

// TelegramBot represents the main bot application
type TelegramBot struct {
	token        string
	url          string
	dependencies *Dependencies
}

// TelegramUpdate represents Telegram webhook update
type TelegramUpdate struct {
	UpdateID int             `json:"update_id"`
	Message  *TelegramMessage `json:"message"`
}

// TelegramMessage represents Telegram message
type TelegramMessage struct {
	MessageID int           `json:"message_id"`
	From      *TelegramUser `json:"from"`
	Chat      *TelegramChat `json:"chat"`
	Text      string        `json:"text"`
	Date      int64         `json:"date"`
	// File attachments
	Document *domain.TelegramDocument `json:"document,omitempty"`
	Photo    []domain.TelegramPhoto   `json:"photo,omitempty"`
}

// TelegramUser represents Telegram user
type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	IsBot     bool   `json:"is_bot"`
}

// TelegramChat represents Telegram chat
type TelegramChat struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Username string `json:"username"`
}

// NewTelegramBot creates a new Telegram bot instance
func NewTelegramBot(token string, dependencies *Dependencies) *TelegramBot {
	return &TelegramBot{
		token:        token,
		url:          fmt.Sprintf("https://api.telegram.org/bot%s", token),
		dependencies: dependencies,
	}
}

// Start starts the bot HTTP server
func (b *TelegramBot) Start(port string) error {
	http.HandleFunc("/webhook", b.handleWebhook)
	http.HandleFunc("/health", b.handleHealth)
	
	b.dependencies.Logger.Info("Bot server starting", "port", port)
	return http.ListenAndServe(":"+port, nil)
}

// handleWebhook processes incoming Telegram webhooks
func (b *TelegramBot) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		b.dependencies.Logger.Error("Failed to read request body", "error", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var update TelegramUpdate
	if err := json.Unmarshal(body, &update); err != nil {
		b.dependencies.Logger.Error("Failed to unmarshal update", "error", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Process the update asynchronously
	go b.processUpdate(&update)

	// Respond to Telegram immediately
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// handleHealth provides health check endpoint
func (b *TelegramBot) handleHealth(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(b.dependencies.StartTime)
	
	health := map[string]interface{}{
		"status":  "healthy",
		"uptime":  uptime.String(),
		"version": "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// processUpdate processes a single Telegram update
func (b *TelegramBot) processUpdate(update *TelegramUpdate) {
	if update.Message == nil {
		return
	}
	
	// Allow messages with files even if they don't have text
	if update.Message.Text == "" && update.Message.Document == nil && len(update.Message.Photo) == 0 {
		return
	}

	// Convert Telegram structures to domain structures
	domainCmd := b.convertToDomainCommand(update.Message)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Route command through the application
	response, err := b.dependencies.Router.Route(ctx, domainCmd)
	if err != nil {
		b.dependencies.Logger.Error("Command routing failed", 
			"command", domainCmd.Text, 
			"user_id", domainCmd.User.TelegramID,
			"error", err)
		
		// Send error response
		b.sendTelegramMessage(update.Message.Chat.ID, "❌ Xatolik yuz berdi. Keyinroq urinib ko'ring.")
		return
	}

	// Send response back to Telegram
	if response != nil && response.Text != "" {
		err = b.sendTelegramMessageWithParseMode(update.Message.Chat.ID, response.Text, response.ParseMode)
		if err != nil {
			b.dependencies.Logger.Error("Failed to send Telegram message", 
				"chat_id", update.Message.Chat.ID,
				"error", err)
		}
	}
}

// convertToDomainCommand converts Telegram message to domain command
func (b *TelegramBot) convertToDomainCommand(msg *TelegramMessage) *domain.Command {
	cmd := &domain.Command{
		ID:   fmt.Sprintf("%d_%d", msg.Chat.ID, msg.MessageID),
		Text: strings.TrimSpace(msg.Text),
		User: &domain.User{
			TelegramID: msg.From.ID,
			Username:   msg.From.Username,
			FirstName:  msg.From.FirstName,
			LastName:   msg.From.LastName,
			Language:   "uz", // Default language
			IsActive:   true,
		},
		Chat: &domain.Chat{
			ID:       msg.Chat.ID,
			Type:     msg.Chat.Type,
			Title:    msg.Chat.Title,
			Username: msg.Chat.Username,
		},
		Timestamp: time.Unix(msg.Date, 0),
		// Include file attachments
		Document:  msg.Document,
		Photo:     msg.Photo,
	}
	
	// If there's no text but there's a file, set the text to /analyze for automatic processing
	if cmd.Text == "" && (msg.Document != nil || len(msg.Photo) > 0) {
		cmd.Text = "/analyze"
	}
	
	return cmd
}

// sendTelegramMessage sends a message to Telegram with HTML parse mode (for backwards compatibility)
func (b *TelegramBot) sendTelegramMessage(chatID int64, text string) error {
	return b.sendTelegramMessageWithParseMode(chatID, text, "HTML")
}

// sendTelegramMessageWithParseMode sends a message to Telegram with specified parse mode
func (b *TelegramBot) sendTelegramMessageWithParseMode(chatID int64, text string, parseMode string) error {
	// Default to HTML if parseMode is empty
	if parseMode == "" {
		parseMode = "HTML"
	}

	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": parseMode,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/sendMessage", b.url)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		
		// If it's a Markdown parsing error and we're using Markdown, fallback to plain text
		if parseMode == "Markdown" && strings.Contains(string(body), "can't parse entities") {
			b.dependencies.Logger.Warn("Markdown parsing failed, falling back to plain text", 
				"chat_id", chatID, 
				"error", string(body))
			
			// Strip Markdown formatting and retry with no parse mode
			plainText := stripMarkdown(text)
			return b.sendTelegramMessageWithParseMode(chatID, plainText, "")
		}
		
		return fmt.Errorf("telegram API error: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

// stripMarkdown removes Markdown formatting from text to create plain text fallback
func stripMarkdown(text string) string {
	// Remove bold formatting **text**
	text = regexp.MustCompile(`\*\*(.*?)\*\*`).ReplaceAllString(text, "$1")
	
	// Remove italic formatting *text*
	text = regexp.MustCompile(`\*(.*?)\*`).ReplaceAllString(text, "$1")
	
	// Remove inline code `text`
	text = regexp.MustCompile("`([^`]*)`").ReplaceAllString(text, "$1")
	
	// Remove links [text](url) -> text
	text = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)`).ReplaceAllString(text, "$1")
	
	// Remove escaped characters
	text = strings.ReplaceAll(text, `\-`, "-")
	text = strings.ReplaceAll(text, `\_`, "_")
	text = strings.ReplaceAll(text, `\!`, "!")
	text = strings.ReplaceAll(text, `\.`, ".")
	
	return text
}