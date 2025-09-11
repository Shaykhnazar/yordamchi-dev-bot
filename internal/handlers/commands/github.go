package commands

import (
	"context"
	"strings"
	"time"

	"yordamchi-dev-bot/internal/domain"
	"yordamchi-dev-bot/internal/services"
)

// GitHubCommand handles GitHub-related commands
type GitHubCommand struct {
	githubService *services.GitHubService
	logger        domain.Logger
}

// NewGitHubCommand creates a new GitHub command handler
func NewGitHubCommand(githubService *services.GitHubService, logger domain.Logger) *GitHubCommand {
	return &GitHubCommand{
		githubService: githubService,
		logger:        logger,
	}
}

// Handle processes GitHub commands
func (h *GitHubCommand) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
	parts := strings.Fields(cmd.Text)
	if len(parts) < 2 {
		return &domain.Response{
			Text:      h.getUsageMessage(),
			ParseMode: "Markdown",
		}, nil
	}

	command := strings.ToLower(parts[0])
	
	switch command {
	case "/repo":
		return h.handleRepoCommand(ctx, parts[1:])
	case "/user":
		return h.handleUserCommand(ctx, parts[1:])
	default:
		return &domain.Response{
			Text:      "❌ Noma'lum GitHub buyruq",
			ParseMode: "Markdown",
		}, nil
	}
}

// CanHandle checks if this handler can process the command
func (h *GitHubCommand) CanHandle(command string) bool {
	cmd := strings.ToLower(strings.Fields(command)[0])
	return cmd == "/repo" || cmd == "/user"
}

// Description returns the command description
func (h *GitHubCommand) Description() string {
	return "GitHub integration - repository and user lookup"
}

// Usage returns the command usage instructions
func (h *GitHubCommand) Usage() string {
	return "/repo owner/name - Repository ma'lumoti\n/user username - Foydalanuvchi profili"
}

// handleRepoCommand handles repository lookup
func (h *GitHubCommand) handleRepoCommand(ctx context.Context, args []string) (*domain.Response, error) {
	if len(args) != 1 {
		return &domain.Response{
			Text:      "❌ Format: /repo owner/repository\nMisol: /repo torvalds/linux",
			ParseMode: "Markdown",
		}, nil
	}

	repoParts := strings.Split(args[0], "/")
	if len(repoParts) != 2 {
		return &domain.Response{
			Text:      "❌ Format: /repo owner/repository\nMisol: /repo torvalds/linux",
			ParseMode: "Markdown",
		}, nil
	}

	owner := repoParts[0]
	repo := repoParts[1]

	ctxTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	repository, err := h.githubService.GetRepository(ctxTimeout, owner, repo)
	if err != nil {
		h.logger.Error("GitHub repository error", "error", err, "owner", owner, "repo", repo)
		return &domain.Response{
			Text:      "❌ Repository topilmadi yoki xatolik yuz berdi",
			ParseMode: "Markdown",
		}, nil
	}

	message := h.githubService.FormatRepository(repository)
	return &domain.Response{
		Text:      message,
		ParseMode: "Markdown",
	}, nil
}

// handleUserCommand handles user lookup
func (h *GitHubCommand) handleUserCommand(ctx context.Context, args []string) (*domain.Response, error) {
	if len(args) != 1 {
		return &domain.Response{
			Text:      "❌ Format: /user username\nMisol: /user torvalds",
			ParseMode: "Markdown",
		}, nil
	}

	username := args[0]

	ctxTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	user, err := h.githubService.GetUser(ctxTimeout, username)
	if err != nil {
		h.logger.Error("GitHub user error", "error", err, "username", username)
		return &domain.Response{
			Text:      "❌ Foydalanuvchi topilmadi yoki xatolik yuz berdi",
			ParseMode: "Markdown",
		}, nil
	}

	message := h.githubService.FormatUser(user)
	return &domain.Response{
		Text:      message,
		ParseMode: "Markdown",
	}, nil
}

// getUsageMessage returns usage instructions
func (h *GitHubCommand) getUsageMessage() string {
	return "🐙 **GitHub Commands**\n\n" +
		"**Repository ma'lumoti:**\n" +
		"`/repo owner/repository`\n" +
		"Misol: `/repo microsoft/vscode`\n\n" +
		"**Foydalanuvchi profili:**\n" +
		"`/user username`\n" +
		"Misol: `/user torvalds`"
}