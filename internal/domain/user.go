package domain

import (
	"context"
	"time"
)

// User represents a bot user
type User struct {
	ID          int64     `json:"id"`
	TelegramID  int64     `json:"telegram_id"`
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Language    string    `json:"language"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Preferences map[string]interface{} `json:"preferences"`
}

// Chat represents a Telegram chat
type Chat struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Username string `json:"username"`
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, telegramID int64) error
	GetActiveUsers(ctx context.Context, limit int) ([]*User, error)
	GetUserStats(ctx context.Context) (*UserStats, error)
}

// UserStats represents user statistics
type UserStats struct {
	TotalUsers   int `json:"total_users"`
	ActiveUsers  int `json:"active_users"`
	NewToday     int `json:"new_today"`
	ActiveToday  int `json:"active_today"`
}

// UserService defines the interface for user business logic
type UserService interface {
	RegisterUser(ctx context.Context, telegramID int64, username, firstName, lastName string) (*User, error)
	GetUser(ctx context.Context, telegramID int64) (*User, error)
	UpdateUserActivity(ctx context.Context, telegramID int64) error
	GetStats(ctx context.Context) (*UserStats, error)
}

// DevTaskMaster Domain Models
// ===========================

// Project represents a development project
type Project struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	TeamID      string    `json:"team_id" db:"team_id"`
	Status      string    `json:"status" db:"status"` // active, completed, paused
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Task represents a development task with AI-generated breakdown
type Task struct {
	ID            string     `json:"id" db:"id"`
	ProjectID     string     `json:"project_id" db:"project_id"`
	Title         string     `json:"title" db:"title"`
	Description   string     `json:"description" db:"description"`
	Category      string     `json:"category" db:"category"`        // backend, frontend, qa, devops
	EstimateHours float64    `json:"estimate_hours" db:"estimate_hours"`
	ActualHours   float64    `json:"actual_hours" db:"actual_hours"`
	Status        string     `json:"status" db:"status"`            // todo, in_progress, completed, blocked
	Priority      int        `json:"priority" db:"priority"`        // 1-5
	AssignedTo    string     `json:"assigned_to" db:"assigned_to"`  // team member ID
	Dependencies  []string   `json:"dependencies" db:"dependencies"` // task IDs this depends on
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	CompletedAt   *time.Time `json:"completed_at" db:"completed_at"`
}

// Team represents a development team
type Team struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	ChatID    int64     `json:"chat_id" db:"chat_id"` // Telegram chat ID
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TeamMember represents a developer in a team
type TeamMember struct {
	ID       string   `json:"id" db:"id"`
	TeamID   string   `json:"team_id" db:"team_id"`
	UserID   int64    `json:"user_id" db:"user_id"`   // Telegram user ID
	Username string   `json:"username" db:"username"`
	Role     string   `json:"role" db:"role"`         // lead, senior, mid, junior
	Skills   []string `json:"skills" db:"skills"`     // go, react, python, etc.
	Capacity float64  `json:"capacity" db:"capacity"` // hours per week
	Current  float64  `json:"current" db:"current"`   // current workload hours
}

// TaskBreakdownRequest represents AI analysis request
type TaskBreakdownRequest struct {
	Requirement string   `json:"requirement"`
	TeamSkills  []string `json:"team_skills"`
	ProjectType string   `json:"project_type"` // web, mobile, api, etc.
}

// TaskBreakdownResponse represents AI analysis result
type TaskBreakdownResponse struct {
	Tasks           []Task   `json:"tasks"`
	TotalEstimate   float64  `json:"total_estimate"`
	RecommendedTeam []string `json:"recommended_team"`
	CriticalPath    []string `json:"critical_path"`
	RiskFactors     []string `json:"risk_factors"`
	Confidence      float64  `json:"confidence"` // 0-1
}

// ProjectStats represents project analytics
type ProjectStats struct {
	ProjectID        string  `json:"project_id"`
	TotalTasks       int     `json:"total_tasks"`
	CompletedTasks   int     `json:"completed_tasks"`
	Progress         float64 `json:"progress"`
	EstimatedHours   float64 `json:"estimated_hours"`
	ActualHours      float64 `json:"actual_hours"`
	EfficiencyRatio  float64 `json:"efficiency_ratio"`
	TeamUtilization  float64 `json:"team_utilization"`
	OnTimeCompletion float64 `json:"on_time_completion"`
}

// TeamWorkload represents current team capacity analysis
type TeamWorkload struct {
	TeamID      string           `json:"team_id"`
	Members     []MemberWorkload `json:"members"`
	Available   float64          `json:"available_hours"`
	Allocated   float64          `json:"allocated_hours"`
	Utilization float64          `json:"utilization"`
}

// MemberWorkload represents individual member workload
type MemberWorkload struct {
	MemberID    string  `json:"member_id"`
	Username    string  `json:"username"`
	Capacity    float64 `json:"capacity"`
	Current     float64 `json:"current"`
	Utilization float64 `json:"utilization"`
	Status      string  `json:"status"` // available, busy, overloaded
}