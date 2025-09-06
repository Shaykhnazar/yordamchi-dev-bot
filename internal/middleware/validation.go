package middleware

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"yordamchi-dev-bot/internal/domain"
)

// ValidationMiddleware provides input validation and sanitization
type ValidationMiddleware struct {
	logger     domain.Logger
	maxLength  int
	validators map[string]*CommandValidator
}

// CommandValidator defines validation rules for specific commands
type CommandValidator struct {
	Pattern     *regexp.Regexp
	MinArgs     int
	MaxArgs     int
	Description string
	Usage       string
}

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware(logger domain.Logger) *ValidationMiddleware {
	validators := make(map[string]*CommandValidator)

	// Weather command validation
	validators["/weather"] = &CommandValidator{
		Pattern:     regexp.MustCompile(`^/weather\s+[a-zA-Z\s\-']{2,50}$`),
		MinArgs:     2,
		MaxArgs:     5,
		Description: "Weather command requires a city name",
		Usage:       "/weather <city_name>",
	}

	validators["/ob-havo"] = &CommandValidator{
		Pattern:     regexp.MustCompile(`^/ob-havo\s+[a-zA-Z\s\-']{2,50}$`),
		MinArgs:     2,
		MaxArgs:     5,
		Description: "Ob-havo command requires a city name",
		Usage:       "/ob-havo <shahar_nomi>",
	}

	// GitHub command validation
	validators["/repo"] = &CommandValidator{
		Pattern:     regexp.MustCompile(`^/repo\s+[a-zA-Z0-9\-_.]+/[a-zA-Z0-9\-_.]+$`),
		MinArgs:     2,
		MaxArgs:     2,
		Description: "Repository command requires owner/repo format",
		Usage:       "/repo <owner/repository>",
	}

	validators["/user"] = &CommandValidator{
		Pattern:     regexp.MustCompile(`^/user\s+[a-zA-Z0-9\-_.]+$`),
		MinArgs:     2,
		MaxArgs:     2,
		Description: "User command requires a GitHub username",
		Usage:       "/user <username>",
	}

	return &ValidationMiddleware{
		logger:     logger,
		maxLength:  500, // Maximum command length
		validators: validators,
	}
}

// Process implements the Middleware interface
func (m *ValidationMiddleware) Process(ctx context.Context, next domain.HandlerFunc) domain.HandlerFunc {
	return func(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
		// Basic length validation
		if len(cmd.Text) > m.maxLength {
			m.logger.Warn("Command too long", 
				"user_id", cmd.User.TelegramID,
				"command_length", len(cmd.Text),
				"max_length", m.maxLength)

			return &domain.Response{
				Text:      fmt.Sprintf("❌ Buyruq juda uzun. Maksimal uzunlik: %d belgi", m.maxLength),
				ParseMode: "HTML",
			}, nil
		}

		// Sanitize input
		cmd.Text = m.sanitizeInput(cmd.Text)

		// Get command parts
		parts := strings.Fields(strings.ToLower(cmd.Text))
		if len(parts) == 0 {
			return next(ctx, cmd)
		}

		baseCommand := parts[0]

		// Check if command needs validation
		validator, exists := m.validators[baseCommand]
		if !exists {
			// No specific validation rules, proceed
			return next(ctx, cmd)
		}

		// Validate argument count
		if len(parts) < validator.MinArgs {
			m.logger.Warn("Command has too few arguments",
				"user_id", cmd.User.TelegramID,
				"command", baseCommand,
				"args_provided", len(parts)-1,
				"min_required", validator.MinArgs-1)

			return &domain.Response{
				Text: fmt.Sprintf(
					"❌ %s\n\n💡 <b>To'g'ri format:</b>\n<code>%s</code>",
					validator.Description,
					validator.Usage,
				),
				ParseMode: "HTML",
			}, nil
		}

		if len(parts) > validator.MaxArgs {
			m.logger.Warn("Command has too many arguments",
				"user_id", cmd.User.TelegramID,
				"command", baseCommand,
				"args_provided", len(parts)-1,
				"max_allowed", validator.MaxArgs-1)

			return &domain.Response{
				Text: fmt.Sprintf(
					"❌ %s\n\n💡 <b>To'g'ri format:</b>\n<code>%s</code>",
					validator.Description,
					validator.Usage,
				),
				ParseMode: "HTML",
			}, nil
		}

		// Pattern validation
		if !validator.Pattern.MatchString(cmd.Text) {
			m.logger.Warn("Command pattern validation failed",
				"user_id", cmd.User.TelegramID,
				"command", cmd.Text,
				"pattern", validator.Pattern.String())

			return &domain.Response{
				Text: fmt.Sprintf(
					"❌ %s\n\n💡 <b>To'g'ri format:</b>\n<code>%s</code>\n\n📝 <b>Misol:</b>\n<code>%s</code>",
					validator.Description,
					validator.Usage,
					m.getExampleUsage(baseCommand),
				),
				ParseMode: "HTML",
			}, nil
		}

		m.logger.Debug("Command validation passed",
			"user_id", cmd.User.TelegramID,
			"command", baseCommand)

		return next(ctx, cmd)
	}
}

// sanitizeInput cleans and normalizes input text
func (m *ValidationMiddleware) sanitizeInput(input string) string {
	// Remove excessive whitespace
	input = regexp.MustCompile(`\s+`).ReplaceAllString(input, " ")
	
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Remove potentially dangerous characters (basic XSS prevention)
	input = strings.ReplaceAll(input, "<script>", "")
	input = strings.ReplaceAll(input, "</script>", "")
	input = strings.ReplaceAll(input, "javascript:", "")
	
	return input
}

// getExampleUsage returns example usage for commands
func (m *ValidationMiddleware) getExampleUsage(command string) string {
	examples := map[string]string{
		"/weather": "/weather London",
		"/ob-havo": "/ob-havo Toshkent",
		"/repo":    "/repo microsoft/vscode",
		"/user":    "/user octocat",
	}

	if example, exists := examples[command]; exists {
		return example
	}

	return command + " <parametr>"
}