package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"yordamchi-dev-bot/database"
	"yordamchi-dev-bot/internal/domain"
)

// ProjectCommand handles project management operations
type ProjectCommand struct {
	db     *database.DB
	logger domain.Logger
}

// NewProjectCommand creates a new project command handler
func NewProjectCommand(db *database.DB, logger domain.Logger) *ProjectCommand {
	return &ProjectCommand{
		db:     db,
		logger: logger,
	}
}

// CanHandle checks if this handler can process the command
func (c *ProjectCommand) CanHandle(command string) bool {
	return command == "/create_project"
}

// Description returns the command description
func (c *ProjectCommand) Description() string {
	return "📝 Create a new development project for task management"
}

// Usage returns the command usage instructions
func (c *ProjectCommand) Usage() string {
	return "/create_project <name> - Create new development project"
}

// Handle processes the create_project command
func (c *ProjectCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	c.logger.Info("Processing create_project command", "user_id", cmd.User.TelegramID, "chat_id", cmd.Chat.ID)

	// Extract project name from command text (skip the command itself)
	cmdText := strings.TrimPrefix(cmd.Text, "/create_project")
	cmdText = strings.TrimSpace(cmdText)
	
	if cmdText == "" {
		return &domain.Response{
			Text: "❌ Please provide a project name.\n\n" +
				"**Example:** `/create_project E-commerce Platform`\n\n" +
				"**Tips:**\n" +
				"• Use descriptive names (E-commerce Platform, Mobile App, etc.)\n" +
				"• Keep it concise but clear\n" +
				"• Avoid special characters",
			ParseMode: "Markdown",
		}, nil
	}

	projectName := cmdText
	
	// Generate project ID
	projectID := generateProjectID()
	
	// Create project
	project := &domain.Project{
		ID:          projectID,
		Name:        projectName,
		Description: fmt.Sprintf("Project created via Telegram bot by @%s", cmd.User.Username),
		TeamID:      fmt.Sprintf("chat_%d", cmd.Chat.ID),
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// In the future, this would save to database
	// For now, we'll just simulate the creation
	c.logger.Info("Project created", 
		"project_id", project.ID,
		"name", project.Name,
		"created_by", cmd.User.TelegramID)
	
	response := fmt.Sprintf("✅ **Project Created Successfully!**\n\n"+
		"📝 **Name:** %s\n"+
		"🆔 **Project ID:** `%s`\n"+
		"👤 **Created by:** @%s\n"+
		"📅 **Created:** %s\n"+
		"📊 **Status:** %s\n\n"+
		"**Next Steps:**\n"+
		"• Use `/analyze <requirement>` to break down features\n"+
		"• Use `/add_member @user skills` to build your team\n"+
		"• Use `/list_projects` to see all your projects",
		projectName, 
		project.ID, 
		cmd.User.Username,
		project.CreatedAt.Format("Jan 2, 2006 15:04"),
		strings.Title(project.Status))
	
	return &domain.Response{
		Text:      response,
		ParseMode: "Markdown",
	}, nil
}

// Helper function to generate project IDs
func generateProjectID() string {
	return fmt.Sprintf("proj_%d", time.Now().UnixNano()%1000000)
}