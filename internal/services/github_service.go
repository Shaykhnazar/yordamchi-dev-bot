package services

import (
	"context"
	"fmt"
	"time"
)

// GitHubService provides GitHub API integration
type GitHubService struct {
	httpClient *HTTPClient
	logger     Logger
}

// GitHubRepository represents a GitHub repository
type GitHubRepository struct {
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	Description     string `json:"description"`
	Stars           int    `json:"stargazers_count"`
	Forks           int    `json:"forks_count"`
	Language        string `json:"language"`
	URL             string `json:"html_url"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	Owner           GitHubUser `json:"owner"`
	DefaultBranch   string `json:"default_branch"`
	OpenIssues      int    `json:"open_issues_count"`
	Topics          []string `json:"topics"`
}

// GitHubUser represents a GitHub user
type GitHubUser struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	Blog        string `json:"blog"`
	Location    string `json:"location"`
	Email       string `json:"email"`
	Bio         string `json:"bio"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
	Following   int    `json:"following"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	AvatarURL   string `json:"avatar_url"`
	URL         string `json:"html_url"`
}

// NewGitHubService creates a new GitHub service
func NewGitHubService(logger Logger) *GitHubService {
	httpClient := NewHTTPClient(30*time.Second, logger)
	
	return &GitHubService{
		httpClient: httpClient,
		logger:     logger,
	}
}

// GetRepository fetches repository information from GitHub
func (g *GitHubService) GetRepository(ctx context.Context, owner, repo string) (*GitHubRepository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	
	var repository GitHubRepository
	err := g.httpClient.GetJSON(ctx, url, nil, &repository)
	if err != nil {
		return nil, fmt.Errorf("GitHub repository ma'lumotlarini olishda xatolik: %w", err)
	}
	
	g.logger.Printf("📦 GitHub repository retrieved: %s/%s", owner, repo)
	return &repository, nil
}

// GetUser fetches user information from GitHub
func (g *GitHubService) GetUser(ctx context.Context, username string) (*GitHubUser, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	
	var user GitHubUser
	err := g.httpClient.GetJSON(ctx, url, nil, &user)
	if err != nil {
		return nil, fmt.Errorf("GitHub foydalanuvchi ma'lumotlarini olishda xatolik: %w", err)
	}
	
	g.logger.Printf("👤 GitHub user retrieved: %s", username)
	return &user, nil
}

// FormatRepository formats repository info for Telegram message
func (g *GitHubService) FormatRepository(repo *GitHubRepository) string {
	description := repo.Description
	if description == "" {
		description = "Tavsif mavjud emas"
	}
	
	language := repo.Language
	if language == "" {
		language = "Aniqlanmagan"
	}
	
	return fmt.Sprintf(`📦 <b>%s</b>

📝 <b>Tavsif:</b> %s
⭐ <b>Yulduzlar:</b> %d
🍴 <b>Forklar:</b> %d
💻 <b>Til:</b> %s
🔧 <b>Asosiy branch:</b> %s
🐛 <b>Ochiq muammolar:</b> %d

👤 <b>Egasi:</b> %s
🔗 <b>Havola:</b> <a href="%s">%s</a>

📅 <b>Yaratilgan:</b> %s
🔄 <b>Yangilangan:</b> %s`,
		repo.FullName,
		description,
		repo.Stars,
		repo.Forks,
		language,
		repo.DefaultBranch,
		repo.OpenIssues,
		repo.Owner.Login,
		repo.URL,
		repo.URL,
		g.formatDate(repo.CreatedAt),
		g.formatDate(repo.UpdatedAt))
}

// FormatUser formats user info for Telegram message
func (g *GitHubService) FormatUser(user *GitHubUser) string {
	name := user.Name
	if name == "" {
		name = user.Login
	}
	
	bio := user.Bio
	if bio == "" {
		bio = "Bio mavjud emas"
	}
	
	company := user.Company
	if company == "" {
		company = "Ko'rsatilmagan"
	}
	
	location := user.Location
	if location == "" {
		location = "Ko'rsatilmagan"
	}
	
	return fmt.Sprintf(`👤 <b>%s</b> (@%s)

📝 <b>Bio:</b> %s
🏢 <b>Kompaniya:</b> %s
📍 <b>Joylashuv:</b> %s
📦 <b>Ochiq repozitoriyalar:</b> %d
👥 <b>Obunachilar:</b> %d
➡️ <b>Obunalar:</b> %d

🔗 <b>Profil:</b> <a href="%s">%s</a>
📅 <b>Ro'yxatdan o'tgan:</b> %s`,
		name,
		user.Login,
		bio,
		company,
		location,
		user.PublicRepos,
		user.Followers,
		user.Following,
		user.URL,
		user.URL,
		g.formatDate(user.CreatedAt))
}

// formatDate formats GitHub date string to readable format
func (g *GitHubService) formatDate(dateStr string) string {
	if dateStr == "" {
		return "Noma'lum"
	}
	
	// Parse GitHub's ISO 8601 format
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	
	return t.Format("2006-01-02 15:04")
}