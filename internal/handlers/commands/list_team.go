package commands

import (
	"context"
	"fmt"
	"strings"

	"yordamchi-dev-bot/database"
	"yordamchi-dev-bot/internal/domain"
)

// ListTeamCommand handles listing team members
type ListTeamCommand struct {
	db     *database.DB
	logger domain.Logger
}

// NewListTeamCommand creates a new list team command handler
func NewListTeamCommand(db *database.DB, logger domain.Logger) *ListTeamCommand {
	return &ListTeamCommand{
		db:     db,
		logger: logger,
	}
}

// CanHandle checks if this handler can process the command
func (c *ListTeamCommand) CanHandle(command string) bool {
	return command == "/list_team"
}

// Description returns the command description
func (c *ListTeamCommand) Description() string {
	return "👥 List all team members and their current workload status"
}

// Usage returns the command usage instructions
func (c *ListTeamCommand) Usage() string {
	return "/list_team - Show all team members and workload"
}

// Handle processes the list_team command
func (c *ListTeamCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	c.logger.Info("Processing list_team command", "user_id", cmd.User.TelegramID, "chat_id", cmd.Chat.ID)

	teamID := fmt.Sprintf("chat_%d", cmd.Chat.ID)

	// For MVP, show mock team data
	// In production, this would query the database for team members
	mockMembers := c.getMockTeamMembers(teamID)

	if len(mockMembers) == 0 {
		return &domain.Response{
			Text: "👥 **No Team Members Found**\n\n" +
				"This chat doesn't have any team members yet.\n\n" +
				"**Get Started:**\n" +
				"• Use `/add_member @username skills` to add team members\n" +
				"• Skills examples: `go,react,docker` or `python,vue,aws`\n" +
				"• Use `/workload` to analyze team capacity after adding members\n\n" +
				"**Example:** `/add_member @alice go,postgresql,docker`",
			ParseMode: "MarkdownV2",
		}, nil
	}

	response := c.formatTeamList(mockMembers)

	c.logger.Info("Team listed",
		"chat_id", cmd.Chat.ID,
		"members_count", len(mockMembers))

	return &domain.Response{
		Text:      response,
		ParseMode: "MarkdownV2",
	}, nil
}

// formatTeamList formats team members for display
func (c *ListTeamCommand) formatTeamList(members []domain.TeamMember) string {
	response := "👥 **Team Members & Workload**\n\n"

	totalCapacity := 0.0
	totalCurrent := 0.0

	for _, member := range members {
		utilization := 0.0
		if member.Capacity > 0 {
			utilization = member.Current / member.Capacity
		}

		statusEmoji := getTeamStatusEmoji(utilization)
		roleEmoji := getRoleEmoji(member.Role)
		utilizationBar := getTeamUtilizationBar(utilization)

		response += fmt.Sprintf("%s **@%s** (%s)\n", roleEmoji, member.Username, strings.Title(member.Role))
		response += fmt.Sprintf("├── 🛠️ Skills: %s\n", strings.Join(member.Skills, ", "))
		response += fmt.Sprintf("├── %s Capacity: %.0fh/week\n", utilizationBar, member.Capacity)
		response += fmt.Sprintf("├── Current: %.0fh (%.0f%% utilization) %s\n", member.Current, utilization*100, statusEmoji)
		response += fmt.Sprintf("└── Status: %s\n\n", getStatusText(utilization))

		totalCapacity += member.Capacity
		totalCurrent += member.Current
	}

	// Team summary
	teamUtilization := 0.0
	if totalCapacity > 0 {
		teamUtilization = totalCurrent / totalCapacity
	}

	teamStatusEmoji := getTeamStatusEmoji(teamUtilization)
	response += fmt.Sprintf("📊 **Team Summary:**\n")
	response += fmt.Sprintf("├── Total Capacity: %.0fh/week\n", totalCapacity)
	response += fmt.Sprintf("├── Current Workload: %.0fh/week\n", totalCurrent)
	response += fmt.Sprintf("└── %s Team Utilization: %.0f%%", teamStatusEmoji, teamUtilization*100)

	// Add utilization guidance
	if teamUtilization > 0.85 {
		response += " (Near capacity ⚠️)"
	} else if teamUtilization < 0.6 {
		response += " (Available capacity ✅)"
	} else {
		response += " (Optimal range ✅)"
	}

	response += "\n\n"

	// Recommendations
	response += "**Team Management:**\n"
	response += "• `/workload` - Detailed workload analysis\n"
	response += "• `/add_member @user skills` - Add more team members\n"
	response += "• `/analyze requirement` - Get smart task assignments\n"

	// Capacity recommendations
	if teamUtilization > 0.85 {
		response += "\n💡 **Recommendation:** Team is near capacity. Consider:\n"
		response += "• Adding team members for upcoming work\n"
		response += "• Extending project timelines\n"
		response += "• Reducing project scope"
	} else if teamUtilization < 0.6 {
		response += "\n💡 **Opportunity:** Team has available capacity for:\n"
		response += "• Taking on additional features\n"
		response += "• Accelerating current projects\n"
		response += "• Training and skill development"
	}

	return response
}

// Mock data generator (would be replaced with database queries)
func (c *ListTeamCommand) getMockTeamMembers(teamID string) []domain.TeamMember {
	return []domain.TeamMember{
		{
			ID:       "member_1",
			TeamID:   teamID,
			Username: "alice",
			Skills:   []string{"go", "postgresql", "docker", "kubernetes"},
			Capacity: 40.0,
			Role:     "lead",
			Current:  34.0, // 85% utilization
		},
		{
			ID:       "member_2",
			TeamID:   teamID,
			Username: "bob",
			Skills:   []string{"react", "typescript", "css", "node.js"},
			Capacity: 40.0,
			Role:     "senior",
			Current:  37.0, // 92% utilization
		},
		{
			ID:       "member_3",
			TeamID:   teamID,
			Username: "carol",
			Skills:   []string{"kubernetes", "docker", "aws", "terraform"},
			Capacity: 40.0,
			Role:     "mid",
			Current:  24.0, // 60% utilization
		},
		{
			ID:       "member_4",
			TeamID:   teamID,
			Username: "david",
			Skills:   []string{"testing", "automation", "selenium", "jest"},
			Capacity: 35.0, // Part-time
			Role:     "junior",
			Current:  24.5, // 70% utilization
		},
	}
}

// Helper functions for formatting
func getTeamStatusEmoji(utilization float64) string {
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

func getRoleEmoji(role string) string {
	switch strings.ToLower(role) {
	case "lead":
		return "👑"
	case "senior":
		return "🎯"
	case "mid":
		return "⭐"
	case "junior":
		return "🌟"
	default:
		return "👤"
	}
}

func getTeamUtilizationBar(utilization float64) string {
	bars := int(utilization * 10)
	if bars > 10 {
		bars = 10
	}

	filled := strings.Repeat("█", bars)
	empty := strings.Repeat("░", 10-bars)
	return filled + empty
}

func getStatusText(utilization float64) string {
	if utilization > 0.9 {
		return "Overloaded - needs rebalancing"
	} else if utilization > 0.75 {
		return "High utilization - monitor closely"
	} else if utilization > 0.6 {
		return "Optimal workload - on track"
	} else if utilization > 0.3 {
		return "Available for more work"
	} else {
		return "Low utilization - needs tasks"
	}
}
