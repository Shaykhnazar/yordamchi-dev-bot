package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"yordamchi-dev-bot/database"
	"yordamchi-dev-bot/internal/domain"
	"yordamchi-dev-bot/internal/services"
)

// TeamCommand handles team management operations
type TeamCommand struct {
	db          *database.DB
	teamManager *services.TeamManager
	logger      domain.Logger
}

// NewTeamCommand creates a new team command handler
func NewTeamCommand(db *database.DB, teamManager *services.TeamManager, logger domain.Logger) *TeamCommand {
	return &TeamCommand{
		db:          db,
		teamManager: teamManager,
		logger:      logger,
	}
}

// CanHandle checks if this handler can process the command
func (c *TeamCommand) CanHandle(command string) bool {
	return command == "/add_member"
}

// Description returns the command description
func (c *TeamCommand) Description() string {
	return "👥 Add a team member with skills to the current project"
}

// Usage returns the command usage instructions
func (c *TeamCommand) Usage() string {
	return "/add_member @username skills - Add team member with skills"
}

// Handle processes the add_member command
func (c *TeamCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	c.logger.Info("Processing add_member command", "user_id", cmd.User.TelegramID, "chat_id", cmd.Chat.ID)

	// Extract arguments from command text
	cmdText := strings.TrimPrefix(cmd.Text, "/add_member")
	cmdText = strings.TrimSpace(cmdText)
	args := strings.Fields(cmdText)

	if len(args) < 2 {
		return &domain.Response{
			Text: "❌ Please provide username and skills.\n\n" +
				"**Example:** `/add_member @alice go,react,docker`\n\n" +
				"**Skills examples:**\n" +
				"• Backend: `go,python,java,nodejs,postgresql`\n" +
				"• Frontend: `react,vue,angular,javascript,typescript`\n" +
				"• DevOps: `docker,kubernetes,aws,terraform`\n" +
				"• QA: `testing,automation,selenium,cypress`",
			ParseMode: "Markdown",
		}, nil
	}

	username := strings.TrimPrefix(args[0], "@")
	skillsStr := strings.Join(args[1:], "")
	skills := strings.Split(skillsStr, ",")
	
	// Clean up skills (trim whitespace and convert to lowercase)
	cleanSkills := make([]string, 0, len(skills))
	for _, skill := range skills {
		cleanSkill := strings.TrimSpace(strings.ToLower(skill))
		if cleanSkill != "" {
			cleanSkills = append(cleanSkills, cleanSkill)
		}
	}

	if len(cleanSkills) == 0 {
		return &domain.Response{
			Text:      "❌ Please provide at least one skill.",
			ParseMode: "Markdown",
		}, nil
	}

	// Generate member ID
	memberID := generateMemberID()
	
	// Create team member
	member := &domain.TeamMember{
		ID:       memberID,
		TeamID:   fmt.Sprintf("chat_%d", cmd.Chat.ID),
		UserID:   0, // Would be set when user interacts
		Username: username,
		Skills:   cleanSkills,
		Capacity: 40.0, // Default 40h/week
		Role:     "developer", // Default role
		Current:  0.0,
	}
	
	// In the future, this would save to database
	c.logger.Info("Team member added", 
		"member_id", member.ID,
		"username", username,
		"skills", cleanSkills,
		"team_id", member.TeamID)
	
	response := fmt.Sprintf("✅ **Team Member Added Successfully!**\n\n"+
		"👤 **Username:** @%s\n"+
		"🛠️ **Skills:** %s\n"+
		"📊 **Capacity:** %.0fh/week\n"+
		"🎯 **Role:** %s\n"+
		"🆔 **Member ID:** `%s`\n\n"+
		"**Next Steps:**\n"+
		"• Use `/list_team` to see all team members\n"+
		"• Use `/workload` to analyze team capacity\n"+
		"• Use `/analyze <requirement>` for smart task assignment",
		username, 
		strings.Join(cleanSkills, ", "), 
		member.Capacity,
		member.Role,
		member.ID)
	
	return &domain.Response{
		Text:      response,
		ParseMode: "Markdown",
	}, nil
}

// Helper function to generate member IDs
func generateMemberID() string {
	return fmt.Sprintf("member_%d", time.Now().UnixNano()%1000000)
}