package middleware

import (
	"context"
	"testing"

	"yordamchi-dev-bot/internal/domain"
)

// MockLogger for testing
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, args ...interface{}) {}
func (m *MockLogger) Info(msg string, args ...interface{})  {}
func (m *MockLogger) Warn(msg string, args ...interface{})  {}
func (m *MockLogger) Error(msg string, args ...interface{}) {}
func (m *MockLogger) With(args ...interface{}) domain.Logger { return m }

// MockHandler for testing middleware
func mockHandler(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	return &domain.Response{
		Text:      "Mock response",
		ParseMode: "HTML",
	}, nil
}

func TestValidationMiddleware_ValidCommands(t *testing.T) {
	logger := &MockLogger{}
	validation := NewValidationMiddleware(logger)

	tests := []struct {
		name    string
		command string
		valid   bool
	}{
		{"Valid weather command", "/weather London", true},
		{"Valid ob-havo command", "/ob-havo Toshkent", true},
		{"Valid repo command", "/repo microsoft/vscode", true},
		{"Valid user command", "/user octocat", true},
		{"Basic command without validation", "/start", true},
		{"Basic command without validation", "/help", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := &domain.Command{
				Text: test.command,
				User: &domain.User{
					TelegramID: 12345,
					FirstName:  "Test",
				},
			}

			handler := validation.Process(context.Background(), mockHandler)
			response, err := handler(context.Background(), cmd)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if test.valid && response.Text != "Mock response" {
				t.Errorf("Expected valid command to pass through, got: %s", response.Text)
			}
		})
	}
}

func TestValidationMiddleware_InvalidCommands(t *testing.T) {
	logger := &MockLogger{}
	validation := NewValidationMiddleware(logger)

	tests := []struct {
		name    string
		command string
	}{
		{"Weather without city", "/weather"},
		{"Weather with invalid chars", "/weather 123$#@"},
		{"Repo wrong format", "/repo invalidformat"},
		{"Repo with spaces", "/repo user name/repo name"},
		{"User with invalid chars", "/user user@#$"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := &domain.Command{
				Text: test.command,
				User: &domain.User{
					TelegramID: 12345,
					FirstName:  "Test",
				},
			}

			handler := validation.Process(context.Background(), mockHandler)
			response, err := handler(context.Background(), cmd)

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Invalid commands should return validation error messages
			if response.Text == "Mock response" {
				t.Error("Expected validation to fail, but command passed through")
			}

			// Check if error message contains expected elements
			if !contains(response.Text, "❌") {
				t.Error("Expected error message to contain error indicator")
			}
		})
	}
}

func TestValidationMiddleware_CommandTooLong(t *testing.T) {
	logger := &MockLogger{}
	validation := NewValidationMiddleware(logger)

	// Create a very long command
	longCommand := "/start " + string(make([]byte, 600)) // Over 500 char limit

	cmd := &domain.Command{
		Text: longCommand,
		User: &domain.User{
			TelegramID: 12345,
			FirstName:  "Test",
		},
	}

	handler := validation.Process(context.Background(), mockHandler)
	response, err := handler(context.Background(), cmd)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if response.Text == "Mock response" {
		t.Error("Expected long command to be rejected")
	}

	if !contains(response.Text, "juda uzun") {
		t.Error("Expected error message about command being too long")
	}
}

func TestValidationMiddleware_InputSanitization(t *testing.T) {
	logger := &MockLogger{}
	validation := NewValidationMiddleware(logger)

	originalCommand := "/start    with   extra   spaces   "
	expectedSanitized := "/start with extra spaces"

	cmd := &domain.Command{
		Text: originalCommand,
		User: &domain.User{
			TelegramID: 12345,
			FirstName:  "Test",
		},
	}

	handler := validation.Process(context.Background(), func(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
		// Check if command was sanitized
		if cmd.Text != expectedSanitized {
			t.Errorf("Expected sanitized command '%s', got '%s'", expectedSanitized, cmd.Text)
		}
		return mockHandler(ctx, cmd)
	})

	_, err := handler(context.Background(), cmd)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}