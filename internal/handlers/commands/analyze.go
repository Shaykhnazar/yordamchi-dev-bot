package commands

import (
	"context"
	"fmt"
	"strings"

	"yordamchi-dev-bot/internal/domain"
	"yordamchi-dev-bot/internal/services"
)

// AnalyzeCommand handles AI-powered task analysis
type AnalyzeCommand struct {
	taskAnalyzer *services.TaskAnalyzer
	logger       domain.Logger
}

// NewAnalyzeCommand creates a new analyze command handler
func NewAnalyzeCommand(taskAnalyzer *services.TaskAnalyzer, logger domain.Logger) *AnalyzeCommand {
	return &AnalyzeCommand{
		taskAnalyzer: taskAnalyzer,
		logger:       logger,
	}
}

// Handle processes the analyze command
func (c *AnalyzeCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	c.logger.Info("Processing analyze command", "user_id", cmd.User.TelegramID)

	// Parse command arguments (everything after /analyze)
	parts := strings.Fields(cmd.Text)
	if len(parts) < 2 {
		return &domain.Response{
			Text: "❌ Please provide a requirement to analyze.\n\n" +
				"**Example:** `/analyze Build user authentication system with OAuth`\n\n" +
				"**Tips:**\n" +
				"• Be specific about technologies (React, Go, PostgreSQL)\n" +
				"• Include scope (backend, frontend, full-stack)\n" +
				"• Mention integrations (GitHub, Stripe, etc.)",
			ParseMode: "Markdown",
		}, nil
	}

	requirement := strings.Join(parts[1:], " ")
	
	// Create analysis request with default team skills
	req := domain.TaskBreakdownRequest{
		Requirement: requirement,
		TeamSkills:  []string{"go", "react", "python", "docker", "postgresql", "javascript", "typescript", "kubernetes"}, 
		ProjectType: "web", // Default to web project
	}
	
	// Analyze with TaskAnalyzer
	result, err := c.taskAnalyzer.AnalyzeRequirement(req)
	if err != nil {
		c.logger.Error("Task analysis failed", "error", err, "requirement", requirement)
		return &domain.Response{
			Text: "❌ **Analysis failed.** Please try again with a clearer requirement.\n\n" +
				"Make sure to:\n" +
				"• Use clear, specific language\n" +
				"• Include technology stack details\n" +
				"• Specify project scope",
			ParseMode: "Markdown",
		}, nil
	}
	
	// Format and send results
	responseText := c.formatTaskBreakdown(result)
	
	c.logger.Info("Task analysis completed", 
		"user_id", cmd.User.TelegramID,
		"tasks_count", len(result.Tasks),
		"total_estimate", result.TotalEstimate,
		"confidence", result.Confidence)
	
	return &domain.Response{
		Text:      responseText,
		ParseMode: "Markdown",
	}, nil
}

// formatTaskBreakdown formats the analysis results for display
func (c *AnalyzeCommand) formatTaskBreakdown(result *domain.TaskBreakdownResponse) string {
	var response strings.Builder
	
	response.WriteString("📋 **Task Breakdown Analysis**\n\n")
	
	// Group tasks by category
	categories := make(map[string][]domain.Task)
	for _, task := range result.Tasks {
		categories[task.Category] = append(categories[task.Category], task)
	}
	
	// Format each category with appropriate icons
	categoryIcons := map[string]string{
		"backend":  "🔐",
		"frontend": "🎨", 
		"qa":       "🧪",
		"devops":   "⚙️",
	}
	
	categoryNames := map[string]string{
		"backend":  "Backend Development",
		"frontend": "Frontend Development", 
		"qa":       "Quality Assurance",
		"devops":   "DevOps & Infrastructure",
	}
	
	for category, tasks := range categories {
		icon := categoryIcons[category]
		if icon == "" {
			icon = "📝"
		}
		
		categoryName := categoryNames[category]
		if categoryName == "" {
			categoryName = strings.Title(category) + " Tasks"
		}
		
		categoryTotal := 0.0
		response.WriteString(fmt.Sprintf("%s **%s** (Est: %.1fh)\n", icon, categoryName, getCategoryTotal(tasks)))
		
		for _, task := range tasks {
			priorityIcon := getPriorityIcon(task.Priority)
			response.WriteString(fmt.Sprintf("├── %s %s - %.1fh\n", priorityIcon, task.Title, task.EstimateHours))
			categoryTotal += task.EstimateHours
		}
		
		response.WriteString(fmt.Sprintf("└── **Subtotal: %.1f hours**\n\n", categoryTotal))
	}
	
	// Total estimate with developer days calculation
	devDays := result.TotalEstimate / 8
	response.WriteString(fmt.Sprintf("⏱️ **Total Estimate: %.1f hours (%.1f developer days)**\n\n", 
		result.TotalEstimate, devDays))
	
	// Recommended team
	if len(result.RecommendedTeam) > 0 {
		response.WriteString("👥 **Recommended Team:**\n")
		for _, member := range result.RecommendedTeam {
			response.WriteString(fmt.Sprintf("• %s\n", member))
		}
		response.WriteString("\n")
	}
	
	// Critical path tasks
	if len(result.CriticalPath) > 0 {
		response.WriteString("🎯 **Critical Path Tasks:** ")
		response.WriteString(fmt.Sprintf("%d high-priority items\n\n", len(result.CriticalPath)))
	}
	
	// Risk factors
	if len(result.RiskFactors) > 0 {
		response.WriteString("⚠️ **Risk Factors & Considerations:**\n")
		for _, risk := range result.RiskFactors {
			response.WriteString(fmt.Sprintf("• %s\n", risk))
		}
		response.WriteString("\n")
	}
	
	// Analysis confidence and next steps
	confidenceEmoji := getConfidenceEmoji(result.Confidence)
	response.WriteString(fmt.Sprintf("📊 **Analysis Confidence:** %s %.0f%%\n\n", confidenceEmoji, result.Confidence*100))
	
	response.WriteString("**Next Steps:**\n")
	response.WriteString("• Use `/create_project <name>` to create a project\n")
	response.WriteString("• Use `/add_member @user skills` to build your team\n")
	response.WriteString("• Use `/workload` to check team capacity")
	
	return response.String()
}

// Helper functions for formatting
func getCategoryTotal(tasks []domain.Task) float64 {
	total := 0.0
	for _, task := range tasks {
		total += task.EstimateHours
	}
	return total
}

func getPriorityIcon(priority int) string {
	switch priority {
	case 1:
		return "🔴" // High priority
	case 2:
		return "🟡" // Medium priority
	case 3:
		return "🟢" // Low priority
	default:
		return "⚪" // Unknown priority
	}
}

func getConfidenceEmoji(confidence float64) string {
	if confidence >= 0.9 {
		return "🎯" // Very high confidence
	} else if confidence >= 0.75 {
		return "✅" // High confidence
	} else if confidence >= 0.6 {
		return "⚠️"  // Medium confidence
	} else {
		return "❓" // Low confidence
	}
}

// CanHandle checks if this handler can process the command
func (c *AnalyzeCommand) CanHandle(command string) bool {
	return command == "/analyze"
}

// Description returns the command description
func (c *AnalyzeCommand) Description() string {
	return "🔍 Break down development requirements into actionable tasks with AI analysis"
}

// Usage returns the command usage instructions
func (c *AnalyzeCommand) Usage() string {
	return "/analyze <requirement> - Analyze development requirements and break them down into tasks"
}