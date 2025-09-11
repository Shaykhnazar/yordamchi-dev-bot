package commands

import (
	"context"
	"fmt"

	"yordamchi-dev-bot/database"
	"yordamchi-dev-bot/internal/domain"
)

// ListProjectsCommand handles listing user projects
type ListProjectsCommand struct {
	db     *database.DB
	logger domain.Logger
}

// NewListProjectsCommand creates a new list projects command handler
func NewListProjectsCommand(db *database.DB, logger domain.Logger) *ListProjectsCommand {
	return &ListProjectsCommand{
		db:     db,
		logger: logger,
	}
}

// CanHandle checks if this handler can process the command
func (c *ListProjectsCommand) CanHandle(command string) bool {
	return command == "/list_projects"
}

// Description returns the command description
func (c *ListProjectsCommand) Description() string {
	return "📊 List all your development projects and their status"
}

// Usage returns the command usage instructions
func (c *ListProjectsCommand) Usage() string {
	return "/list_projects - Show all development projects"
}

// Handle processes the list_projects command
func (c *ListProjectsCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	c.logger.Info("Processing list_projects command", "user_id", cmd.User.TelegramID, "chat_id", cmd.Chat.ID)

	// Get real projects from database
	projects, err := c.db.GetProjectsByChatID(cmd.Chat.ID)
	if err != nil {
		c.logger.Error("Failed to get projects", "error", err, "chat_id", cmd.Chat.ID)
		return &domain.Response{
			Text:      "❌ Failed to retrieve projects. Please try again.",
			ParseMode: "Markdown",
		}, nil
	}

	if len(projects) == 0 {
		return &domain.Response{
			Text: "📋 **No Projects Found**\n\n" +
				"You haven't created any projects yet.\n\n" +
				"**Get Started:**\n" +
				"• Use `/create_project project_name` to create your first project\n" +
				"• Use `/analyze requirement` to break down features\n" +
				"• Use `/add_member @user skills` to build your team\n\n" +
				"**Example:** `/create_project E-commerce Platform`",
			ParseMode: "Markdown",
		}, nil
	}

	response := c.formatProjectsList(projects)

	c.logger.Info("Projects listed",
		"user_id", cmd.User.TelegramID,
		"projects_count", len(projects))

	return &domain.Response{
		Text:      response,
		ParseMode: "Markdown",
	}, nil
}

// formatProjectsList formats projects for display
func (c *ListProjectsCommand) formatProjectsList(projects []database.Project) string {
	response := "📊 **Your Development Projects**\n\n"

	activeProjects := []database.Project{}
	completedProjects := []database.Project{}
	pausedProjects := []database.Project{}

	// Group projects by status
	for _, project := range projects {
		switch project.Status {
		case "active":
			activeProjects = append(activeProjects, project)
		case "completed":
			completedProjects = append(completedProjects, project)
		case "paused":
			pausedProjects = append(pausedProjects, project)
		}
	}

	// Display active projects
	if len(activeProjects) > 0 {
		response += "🟢 **Active Projects:**\n"
		for _, project := range activeProjects {
			progress := c.getProjectProgress(project.ID)
			response += fmt.Sprintf("├── **%s** (`%s`)\n", project.Name, project.ID)
			response += fmt.Sprintf("│   └── Progress: %s %.0f%% complete\n", getProgressBar(progress), progress*100)
		}
		response += "\n"
	}

	// Display paused projects
	if len(pausedProjects) > 0 {
		response += "🟡 **Paused Projects:**\n"
		for _, project := range pausedProjects {
			response += fmt.Sprintf("├── **%s** (`%s`)\n", project.Name, project.ID)
			response += "│   └── Status: On hold\n"
		}
		response += "\n"
	}

	// Display completed projects
	if len(completedProjects) > 0 {
		response += "✅ **Completed Projects:**\n"
		for _, project := range completedProjects {
			response += fmt.Sprintf("├── **%s** (`%s`)\n", project.Name, project.ID)
			response += fmt.Sprintf("│   └── Completed: %s\n", project.UpdatedAt.Format("Jan 2, 2006"))
		}
		response += "\n"
	}

	// Summary and next steps
	totalProjects := len(projects)
	response += fmt.Sprintf("📈 **Summary:** %d total project", totalProjects)
	if totalProjects != 1 {
		response += "s"
	}
	response += "\n\n"

	response += "**Available Actions:**\n"
	response += "• `/analyze requirement` - Break down new features\n"
	response += "• `/workload` - Check team capacity\n"
	response += "• `/project_stats <id>` - Detailed project analytics\n"
	response += "• `/create_project project_name` - Start a new project"

	return response
}

// Mock data generator (would be replaced with database queries)
func (c *ListProjectsCommand) getMockProjects(userID int64) []domain.Project {
	return []domain.Project{
		{
			ID:          "proj_123",
			Name:        "E-commerce Platform",
			Description: "Full-stack e-commerce solution with React and Go",
			Status:      "active",
		},
		{
			ID:          "proj_124",
			Name:        "Mobile App Backend",
			Description: "REST API backend for mobile application",
			Status:      "active",
		},
		{
			ID:          "proj_125",
			Name:        "Analytics Dashboard",
			Description: "Real-time analytics and reporting dashboard",
			Status:      "paused",
		},
		{
			ID:          "proj_126",
			Name:        "Authentication Service",
			Description: "OAuth2 authentication microservice",
			Status:      "completed",
		},
	}
}

// Get real project progress from database
func (c *ListProjectsCommand) getProjectProgress(projectID string) float64 {
	stats, err := c.db.GetProjectStats(projectID)
	if err != nil {
		c.logger.Error("Failed to get project stats", "error", err, "project_id", projectID)
		return 0.0
	}

	return stats.Progress
}

// Helper function to create progress bar
func getProgressBar(progress float64) string {
	bars := int(progress * 10)
	filled := ""
	empty := ""

	for i := 0; i < bars; i++ {
		filled += "█"
	}
	for i := bars; i < 10; i++ {
		empty += "░"
	}

	return filled + empty
}
