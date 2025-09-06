package commands

import (
	"context"
	"testing"

	"yordamchi-dev-bot/internal/domain"
)

// MockLogger implements domain.Logger for testing
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, args ...interface{}) {}
func (m *MockLogger) Info(msg string, args ...interface{})  {}
func (m *MockLogger) Warn(msg string, args ...interface{})  {}
func (m *MockLogger) Error(msg string, args ...interface{}) {}
func (m *MockLogger) With(args ...interface{}) domain.Logger { return m }

func TestStartCommand_Handle(t *testing.T) {
	logger := &MockLogger{}
	welcomeMsg := "Welcome to test bot!"
	
	startCmd := NewStartCommand(welcomeMsg, logger)

	// Create test command
	cmd := &domain.Command{
		ID:   "test-1",
		Text: "/start",
		User: &domain.User{
			TelegramID: 12345,
			FirstName:  "Test",
			Username:   "testuser",
		},
		Chat: &domain.Chat{
			ID:   67890,
			Type: "private",
		},
	}

	ctx := context.Background()
	response, err := startCmd.Handle(ctx, cmd)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if response == nil {
		t.Fatal("Expected response, got nil")
	}

	if response.Text == "" {
		t.Error("Expected non-empty response text")
	}

	if !contains(response.Text, welcomeMsg) {
		t.Errorf("Expected response to contain welcome message '%s'", welcomeMsg)
	}

	if response.ParseMode != "HTML" {
		t.Errorf("Expected parse mode 'HTML', got '%s'", response.ParseMode)
	}
}

func TestStartCommand_CanHandle(t *testing.T) {
	logger := &MockLogger{}
	startCmd := NewStartCommand("Welcome", logger)

	tests := []struct {
		command  string
		expected bool
	}{
		{"/start", true},
		{"/START", false}, // Case sensitive
		{"start", false},
		{"/help", false},
		{"", false},
		{"/start extra", false}, // Exact match required
	}

	for _, test := range tests {
		result := startCmd.CanHandle(test.command)
		if result != test.expected {
			t.Errorf("CanHandle('%s'): expected %v, got %v", test.command, test.expected, result)
		}
	}
}

func TestStartCommand_Description(t *testing.T) {
	logger := &MockLogger{}
	startCmd := NewStartCommand("Welcome", logger)

	description := startCmd.Description()
	if description == "" {
		t.Error("Expected non-empty description")
	}
}

func TestStartCommand_Usage(t *testing.T) {
	logger := &MockLogger{}
	startCmd := NewStartCommand("Welcome", logger)

	usage := startCmd.Usage()
	if usage == "" {
		t.Error("Expected non-empty usage")
	}

	if !contains(usage, "/start") {
		t.Error("Expected usage to contain '/start'")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		 findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}