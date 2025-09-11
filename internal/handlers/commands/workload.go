package commands

import (
	"context"
	"fmt"
	"strings"

	"yordamchi-dev-bot/database"
	"yordamchi-dev-bot/internal/domain"
	"yordamchi-dev-bot/internal/services"
)

// WorkloadCommand handles team workload analysis
type WorkloadCommand struct {
	db          *database.DB
	teamManager *services.TeamManager
	logger      domain.Logger
}

// NewWorkloadCommand creates a new workload command handler
func NewWorkloadCommand(db *database.DB, teamManager *services.TeamManager, logger domain.Logger) *WorkloadCommand {
	return &WorkloadCommand{
		db:          db,
		teamManager: teamManager,
		logger:      logger,
	}
}

// CanHandle checks if this handler can process the command
func (c *WorkloadCommand) CanHandle(command string) bool {
	return command == "/workload"
}

// Description returns the command description
func (c *WorkloadCommand) Description() string {
	return "📊 Analyze current team workload and get optimization recommendations"
}

// Usage returns the command usage instructions
func (c *WorkloadCommand) Usage() string {
	return "/workload - Analyze team workload and capacity"
}

// Handle processes the workload command
func (c *WorkloadCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	c.logger.Info("Processing workload command", "user_id", cmd.User.TelegramID, "chat_id", cmd.Chat.ID)

	teamID := fmt.Sprintf("chat_%d", cmd.Chat.ID)
	
	// For MVP, we'll show mock data since we don't have database integration yet
	// In production, this would fetch real data from database
	mockMembers := c.getMockTeamMembers(teamID)
	mockTasks := c.getMockTasks()
	
	if len(mockMembers) == 0 {
		return &domain.Response{
			Text: "❌ No team members found for this chat.\n\n" +
				"**Get Started:**\n" +
				"• Use `/add_member @username skills` to add team members\n" +
				"• Use `/create_project project_name` to create a project\n" +
				"• Use `/analyze requirement` to generate tasks\n\n" +
				"**Example:** `/add_member @alice go,react,docker`",
			ParseMode: "Markdown",
		}, nil
	}
	
	// Analyze workload using TeamManager
	workload := c.teamManager.AnalyzeWorkload(teamID, mockMembers, mockTasks)
	
	// Format and return results
	response := c.formatWorkloadAnalysis(workload)
	
	c.logger.Info("Workload analysis completed", 
		"team_id", teamID,
		"members_count", len(workload.Members),
		"total_utilization", workload.Utilization)
	
	return &domain.Response{
		Text:      response,
		ParseMode: "Markdown",
	}, nil
}

// formatWorkloadAnalysis formats workload data for display
func (c *WorkloadCommand) formatWorkloadAnalysis(workload *domain.TeamWorkload) string {
	var response strings.Builder
	
	response.WriteString("📊 **Team Workload Analysis**\n\n")
	
	// Team overview
	utilizationEmoji := getUtilizationEmoji(workload.Utilization)
	response.WriteString(fmt.Sprintf("**Team Overview:**\n"))
	response.WriteString(fmt.Sprintf("├── Total Capacity: %.1fh/week\n", workload.Available))
	response.WriteString(fmt.Sprintf("├── Currently Allocated: %.1fh/week\n", workload.Allocated))
	response.WriteString(fmt.Sprintf("└── %s Team Utilization: %.0f%%\n\n", utilizationEmoji, workload.Utilization*100))
	
	// Individual member workloads
	response.WriteString("👥 **Individual Workloads:**\n")
	
	for _, member := range workload.Members {
		statusEmoji := getStatusEmoji(member.Status)
		utilizationBar := getUtilizationBar(member.Utilization)
		
		response.WriteString(fmt.Sprintf("👤 **@%s**\n", member.Username))
		response.WriteString(fmt.Sprintf("├── %s Capacity: %.1fh/week\n", utilizationBar, member.Capacity))
		response.WriteString(fmt.Sprintf("├── Current: %.1fh (%.0f%% utilization)\n", member.Current, member.Utilization*100))
		response.WriteString(fmt.Sprintf("└── %s Status: %s\n\n", statusEmoji, strings.Title(member.Status)))
	}
	
	// Alerts and recommendations
	alerts := c.generateAlerts(workload)
	if len(alerts) > 0 {
		response.WriteString("⚠️ **Alerts:**\n")
		for _, alert := range alerts {
			response.WriteString(fmt.Sprintf("• %s\n", alert))
		}
		response.WriteString("\n")
	}
	
	recommendations := c.generateRecommendations(workload)
	if len(recommendations) > 0 {
		response.WriteString("💡 **Recommendations:**\n")
		for _, rec := range recommendations {
			response.WriteString(fmt.Sprintf("• %s\n", rec))
		}
		response.WriteString("\n")
	}
	
	// Impact analysis
	response.WriteString("📈 **Optimization Impact:**\n")
	if workload.Utilization > 0.85 {
		response.WriteString("• High utilization detected - consider timeline adjustment\n")
		response.WriteString("• Risk of burnout if sustained long-term\n")
		response.WriteString("• Consider adding team members or reducing scope\n")
	} else if workload.Utilization < 0.6 {
		response.WriteString("• Team has available capacity for additional work\n")
		response.WriteString("• Consider taking on new features or projects\n")
		response.WriteString("• Opportunity to accelerate timeline\n")
	} else {
		response.WriteString("• Team utilization is in optimal range (60-85%)\n")
		response.WriteString("• Good balance of productivity and sustainability\n")
		response.WriteString("• Continue current pace\n")
	}
	
	return response.String()
}

// Mock data generators (would be replaced with database queries in production)
func (c *WorkloadCommand) getMockTeamMembers(teamID string) []domain.TeamMember {
	return []domain.TeamMember{
		{
			ID:       "member_1",
			TeamID:   teamID,
			Username: "alice",
			Skills:   []string{"go", "postgresql", "docker"},
			Capacity: 40.0,
			Role:     "lead",
			Current:  34.0,
		},
		{
			ID:       "member_2", 
			TeamID:   teamID,
			Username: "bob",
			Skills:   []string{"react", "typescript", "css"},
			Capacity: 40.0,
			Role:     "senior",
			Current:  37.0,
		},
		{
			ID:       "member_3",
			TeamID:   teamID,
			Username: "carol",
			Skills:   []string{"kubernetes", "docker", "aws"},
			Capacity: 40.0,
			Role:     "mid",
			Current:  24.0,
		},
	}
}

func (c *WorkloadCommand) getMockTasks() []domain.Task {
	return []domain.Task{
		{
			ID:            "task_1",
			AssignedTo:    "member_1",
			EstimateHours: 20.0,
			Status:        "in_progress",
		},
		{
			ID:            "task_2",
			AssignedTo:    "member_2", 
			EstimateHours: 25.0,
			Status:        "todo",
		},
		{
			ID:            "task_3",
			AssignedTo:    "member_3",
			EstimateHours: 15.0,
			Status:        "todo",
		},
	}
}

// Helper functions for formatting
func getUtilizationEmoji(utilization float64) string {
	if utilization > 0.9 {
		return "🔴" // Overloaded
	} else if utilization > 0.75 {
		return "🟡" // High utilization
	} else if utilization > 0.6 {
		return "🟢" // Optimal
	} else {
		return "🔵" // Under-utilized
	}
}

func getStatusEmoji(status string) string {
	switch status {
	case "overloaded":
		return "🚨"
	case "busy":
		return "⚠️"
	case "available":
		return "✅"
	default:
		return "ℹ️"
	}
}

func getUtilizationBar(utilization float64) string {
	bars := int(utilization * 10)
	filled := strings.Repeat("█", bars)
	empty := strings.Repeat("░", 10-bars)
	return filled + empty
}

func (c *WorkloadCommand) generateAlerts(workload *domain.TeamWorkload) []string {
	alerts := []string{}
	
	for _, member := range workload.Members {
		if member.Status == "overloaded" {
			alerts = append(alerts, fmt.Sprintf("@%s is overloaded (%.0f%% utilization)", member.Username, member.Utilization*100))
		}
	}
	
	if workload.Utilization > 0.85 {
		alerts = append(alerts, "Team approaching maximum capacity")
	}
	
	return alerts
}

func (c *WorkloadCommand) generateRecommendations(workload *domain.TeamWorkload) []string {
	recommendations := []string{}
	
	// Find overloaded and underloaded members
	var overloaded, underloaded []domain.MemberWorkload
	for _, member := range workload.Members {
		if member.Utilization > 0.9 {
			overloaded = append(overloaded, member)
		} else if member.Utilization < 0.6 {
			underloaded = append(underloaded, member)
		}
	}
	
	// Generate redistribution recommendations
	if len(overloaded) > 0 && len(underloaded) > 0 {
		recommendations = append(recommendations, 
			fmt.Sprintf("Reassign tasks from @%s to @%s", overloaded[0].Username, underloaded[0].Username))
	}
	
	if workload.Utilization < 0.6 {
		recommendations = append(recommendations, "Team has capacity for additional work")
	}
	
	if workload.Utilization > 0.85 {
		recommendations = append(recommendations, "Consider extending timeline by 0.5-1 days")
	}
	
	return recommendations
}