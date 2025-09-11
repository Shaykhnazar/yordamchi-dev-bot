package database

import (
    "database/sql"
    "fmt"
    "log"
    "strings"
    "time"
_ 	"github.com/mattn/go-sqlite3"
)

type User struct {
    ID        int       `json:"id"`
    TelegramID int64    `json:"telegram_id"`
    Username  string    `json:"username"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Project represents a project in the database
type Project struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    TeamID      string    `json:"team_id"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Task represents a task in the database
type Task struct {
    ID            string     `json:"id"`
    ProjectID     string     `json:"project_id"`
    Title         string     `json:"title"`
    Description   string     `json:"description"`
    Category      string     `json:"category"`
    EstimateHours float64    `json:"estimate_hours"`
    ActualHours   float64    `json:"actual_hours"`
    Status        string     `json:"status"`
    Priority      int        `json:"priority"`
    AssignedTo    string     `json:"assigned_to"`
    Dependencies  []string   `json:"dependencies"`
    CreatedAt     time.Time  `json:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at"`
    CompletedAt   *time.Time `json:"completed_at"`
}

// TeamMember represents a team member in the database
type TeamMember struct {
    ID       string   `json:"id"`
    TeamID   string   `json:"team_id"`
    UserID   int64    `json:"user_id"`
    Username string   `json:"username"`
    Role     string   `json:"role"`
    Skills   []string `json:"skills"`
    Capacity float64  `json:"capacity"`
    Current  float64  `json:"current"`
}

// ProjectStats represents project statistics
type ProjectStats struct {
    ProjectID        string  `json:"project_id"`
    TotalTasks       int     `json:"total_tasks"`
    CompletedTasks   int     `json:"completed_tasks"`
    Progress         float64 `json:"progress"`
    EstimatedHours   float64 `json:"estimated_hours"`
    ActualHours      float64 `json:"actual_hours"`
    EfficiencyRatio  float64 `json:"efficiency_ratio"`
}

type DB struct {
    conn *sql.DB
}

func NewDB() (*DB, error) {
    conn, err := sql.Open("sqlite3", "./yordamchi_bot.db")
    if err != nil {
        return nil, fmt.Errorf("ma'lumotlar bazasiga ulanishda xatolik: %w", err)
    }

    db := &DB{conn: conn}
    
    if err := db.createTables(); err != nil {
        return nil, fmt.Errorf("jadvallar yaratishda xatolik: %w", err)
    }

    log.Println("✅ Ma'lumotlar bazasi muvaffaqiyatli sozlandi")
    return db, nil
}

func (db *DB) createTables() error {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        telegram_id INTEGER UNIQUE NOT NULL,
        username TEXT,
        first_name TEXT,
        last_name TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS user_activity (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        command TEXT,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

    CREATE TABLE IF NOT EXISTS teams (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        chat_id INTEGER UNIQUE NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS team_members (
        id TEXT PRIMARY KEY,
        team_id TEXT NOT NULL,
        user_id INTEGER NOT NULL,
        username TEXT NOT NULL,
        role TEXT DEFAULT 'developer',
        skills TEXT,
        capacity REAL DEFAULT 40.0,
        current_workload REAL DEFAULT 0.0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (team_id) REFERENCES teams (id),
        FOREIGN KEY (user_id) REFERENCES users (telegram_id)
    );

    CREATE TABLE IF NOT EXISTS projects (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT,
        team_id TEXT NOT NULL,
        status TEXT DEFAULT 'active',
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (team_id) REFERENCES teams (id)
    );

    CREATE TABLE IF NOT EXISTS tasks (
        id TEXT PRIMARY KEY,
        project_id TEXT NOT NULL,
        title TEXT NOT NULL,
        description TEXT,
        category TEXT NOT NULL,
        estimate_hours REAL DEFAULT 0.0,
        actual_hours REAL DEFAULT 0.0,
        status TEXT DEFAULT 'todo',
        priority INTEGER DEFAULT 3,
        assigned_to TEXT,
        dependencies TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        completed_at DATETIME,
        FOREIGN KEY (project_id) REFERENCES projects (id),
        FOREIGN KEY (assigned_to) REFERENCES team_members (id)
    );
    `

    _, err := db.conn.Exec(query)
    return err
}

func (db *DB) CreateOrUpdateUser(telegramID int64, username, firstName, lastName string) error {
    query := `
    INSERT INTO users (telegram_id, username, first_name, last_name)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT(telegram_id) DO UPDATE SET
        username = EXCLUDED.username,
        first_name = EXCLUDED.first_name,
        last_name = EXCLUDED.last_name,
        updated_at = CURRENT_TIMESTAMP
    `

    _, err := db.conn.Exec(query, telegramID, username, firstName, lastName)
    if err != nil {
        return fmt.Errorf("foydalanuvchini saqlashda xatolik: %w", err)
    }

    log.Printf("👤 Foydalanuvchi saqlandi: %s (@%s)", firstName, username)
    return nil
}

func (db *DB) LogUserActivity(telegramID int64, command string) error {
    userIDQuery := "SELECT id FROM users WHERE telegram_id = $1"
    var userID int
    err := db.conn.QueryRow(userIDQuery, telegramID).Scan(&userID)
    if err != nil {
        return fmt.Errorf("foydalanuvchi ID topilmadi: %w", err)
    }

    activityQuery := "INSERT INTO user_activity (user_id, command) VALUES ($1, $2)"
    _, err = db.conn.Exec(activityQuery, userID, command)
    if err != nil {
        return fmt.Errorf("faollik yozishda xatolik: %w", err)
    }

    return nil
}

// UserActivity represents user activity data
type UserActivity struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    Command   string    `json:"command"`
    CreatedAt time.Time `json:"created_at"`
}

// GetUserStats returns total user count
func (db *DB) GetUserStats() (int, error) {
    query := "SELECT COUNT(*) FROM users"
    var count int
    err := db.conn.QueryRow(query).Scan(&count)
    return count, err
}

// GetUserActivities returns recent activities for a user
func (db *DB) GetUserActivities(telegramID int64, limit int) ([]UserActivity, error) {
    query := `
    SELECT ua.id, ua.user_id, ua.command, ua.timestamp 
    FROM user_activity ua
    JOIN users u ON ua.user_id = u.id
    WHERE u.telegram_id = ?
    ORDER BY ua.timestamp DESC 
    LIMIT ?`

    rows, err := db.conn.Query(query, telegramID, limit)
    if err != nil {
        return nil, fmt.Errorf("faollikni olishda xatolik: %w", err)
    }
    defer rows.Close()

    var activities []UserActivity
    for rows.Next() {
        var activity UserActivity
        err := rows.Scan(
            &activity.ID,
            &activity.UserID,
            &activity.Command,
            &activity.CreatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("faollik ma'lumotlarini o'qishda xatolik: %w", err)
        }
        activities = append(activities, activity)
    }

    return activities, nil
}

// GetPopularCommands returns most used commands
func (db *DB) GetPopularCommands(limit int) (map[string]int, error) {
    query := `
    SELECT command, COUNT(*) as count 
    FROM user_activity 
    GROUP BY command 
    ORDER BY count DESC 
    LIMIT ?`

    rows, err := db.conn.Query(query, limit)
    if err != nil {
        return nil, fmt.Errorf("populyar buyruqlarni olishda xatolik: %w", err)
    }
    defer rows.Close()

    commands := make(map[string]int)
    for rows.Next() {
        var command string
        var count int
        err := rows.Scan(&command, &count)
        if err != nil {
            return nil, fmt.Errorf("buyruq ma'lumotlarini o'qishda xatolik: %w", err)
        }
        commands[command] = count
    }

    return commands, nil
}

// GetDailyStats returns activity stats for today
func (db *DB) GetDailyStats() (map[string]int, error) {
    stats := make(map[string]int)
    
    // Total users today
    query := "SELECT COUNT(*) FROM users WHERE DATE(created_at) = DATE('now')"
    var newUsersToday int
    err := db.conn.QueryRow(query).Scan(&newUsersToday)
    if err != nil {
        return nil, fmt.Errorf("bugungi foydalanuvchilar sonini olishda xatolik: %w", err)
    }
    stats["new_users_today"] = newUsersToday
    
    // Activities today
    query = "SELECT COUNT(*) FROM user_activity WHERE DATE(timestamp) = DATE('now')"
    var activitiesToday int
    err = db.conn.QueryRow(query).Scan(&activitiesToday)
    if err != nil {
        return nil, fmt.Errorf("bugungi faollik sonini olishda xatolik: %w", err)
    }
    stats["activities_today"] = activitiesToday
    
    // Active users today
    query = "SELECT COUNT(DISTINCT user_id) FROM user_activity WHERE DATE(timestamp) = DATE('now')"
    var activeUsersToday int
    err = db.conn.QueryRow(query).Scan(&activeUsersToday)
    if err != nil {
        return nil, fmt.Errorf("bugungi faol foydalanuvchilar sonini olishda xatolik: %w", err)
    }
    stats["active_users_today"] = activeUsersToday
    
    return stats, nil
}

// Project methods
func (db *DB) CreateProject(project *Project) error {
    query := `
    INSERT INTO projects (id, name, description, team_id, status)
    VALUES (?, ?, ?, ?, ?)`
    
    _, err := db.conn.Exec(query, project.ID, project.Name, project.Description, project.TeamID, project.Status)
    if err != nil {
        return fmt.Errorf("loyiha yaratishda xatolik: %w", err)
    }
    
    log.Printf("📝 Loyiha yaratildi: %s (ID: %s)", project.Name, project.ID)
    return nil
}

func (db *DB) GetProjectsByChatID(chatID int64) ([]Project, error) {
    // First get the team for this chat
    teamQuery := "SELECT id FROM teams WHERE chat_id = ?"
    var teamID string
    err := db.conn.QueryRow(teamQuery, chatID).Scan(&teamID)
    if err != nil {
        // If no team exists, create one
        teamID = fmt.Sprintf("team_%d", chatID)
        createTeamQuery := "INSERT INTO teams (id, name, chat_id) VALUES (?, ?, ?)"
        _, err = db.conn.Exec(createTeamQuery, teamID, fmt.Sprintf("Chat %d Team", chatID), chatID)
        if err != nil {
            return nil, fmt.Errorf("jamoa yaratishda xatolik: %w", err)
        }
    }
    
    query := `
    SELECT id, name, description, team_id, status, created_at, updated_at 
    FROM projects 
    WHERE team_id = ?
    ORDER BY created_at DESC`
    
    rows, err := db.conn.Query(query, teamID)
    if err != nil {
        return nil, fmt.Errorf("loyihalarni olishda xatolik: %w", err)
    }
    defer rows.Close()
    
    var projects []Project
    for rows.Next() {
        var project Project
        err := rows.Scan(
            &project.ID,
            &project.Name,
            &project.Description,
            &project.TeamID,
            &project.Status,
            &project.CreatedAt,
            &project.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("loyiha ma'lumotlarini o'qishda xatolik: %w", err)
        }
        projects = append(projects, project)
    }
    
    return projects, nil
}

// Task methods
func (db *DB) CreateTask(task *Task) error {
    dependencies := ""
    if len(task.Dependencies) > 0 {
        dependencies = fmt.Sprintf("[%s]", strings.Join(task.Dependencies, ","))
    }
    
    query := `
    INSERT INTO tasks (id, project_id, title, description, category, estimate_hours, status, priority, assigned_to, dependencies)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
    
    _, err := db.conn.Exec(query, 
        task.ID, task.ProjectID, task.Title, task.Description, 
        task.Category, task.EstimateHours, task.Status, task.Priority, 
        task.AssignedTo, dependencies)
    
    if err != nil {
        return fmt.Errorf("vazifa yaratishda xatolik: %w", err)
    }
    
    return nil
}

func (db *DB) GetTasksByProjectID(projectID string) ([]Task, error) {
    query := `
    SELECT id, project_id, title, description, category, estimate_hours, actual_hours, 
           status, priority, assigned_to, dependencies, created_at, updated_at, completed_at
    FROM tasks 
    WHERE project_id = ?
    ORDER BY priority ASC, created_at ASC`
    
    rows, err := db.conn.Query(query, projectID)
    if err != nil {
        return nil, fmt.Errorf("vazifalarni olishda xatolik: %w", err)
    }
    defer rows.Close()
    
    var tasks []Task
    for rows.Next() {
        var task Task
        var dependencies sql.NullString
        var completedAt sql.NullTime
        
        err := rows.Scan(
            &task.ID,
            &task.ProjectID,
            &task.Title,
            &task.Description,
            &task.Category,
            &task.EstimateHours,
            &task.ActualHours,
            &task.Status,
            &task.Priority,
            &task.AssignedTo,
            &dependencies,
            &task.CreatedAt,
            &task.UpdatedAt,
            &completedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("vazifa ma'lumotlarini o'qishda xatolik: %w", err)
        }
        
        // Parse dependencies
        if dependencies.Valid && dependencies.String != "" {
            // Remove brackets and split by comma
            depStr := strings.Trim(dependencies.String, "[]")
            if depStr != "" {
                task.Dependencies = strings.Split(depStr, ",")
            }
        }
        
        if completedAt.Valid {
            task.CompletedAt = &completedAt.Time
        }
        
        tasks = append(tasks, task)
    }
    
    return tasks, nil
}

// Team Member methods
func (db *DB) CreateTeamMember(member *TeamMember) error {
    skillsJSON := strings.Join(member.Skills, ",")
    
    query := `
    INSERT INTO team_members (id, team_id, user_id, username, role, skills, capacity, current_workload)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
    
    _, err := db.conn.Exec(query, 
        member.ID, member.TeamID, member.UserID, member.Username, 
        member.Role, skillsJSON, member.Capacity, member.Current)
    
    if err != nil {
        return fmt.Errorf("jamoa a'zosini yaratishda xatolik: %w", err)
    }
    
    log.Printf("👥 Jamoa a'zosi qo'shildi: %s (@%s)", member.Username, member.Username)
    return nil
}

func (db *DB) GetTeamMembersByChatID(chatID int64) ([]TeamMember, error) {
    // First get the team for this chat
    teamQuery := "SELECT id FROM teams WHERE chat_id = ?"
    var teamID string
    err := db.conn.QueryRow(teamQuery, chatID).Scan(&teamID)
    if err != nil {
        // If no team exists, create one
        teamID = fmt.Sprintf("team_%d", chatID)
        createTeamQuery := "INSERT INTO teams (id, name, chat_id) VALUES (?, ?, ?)"
        _, err = db.conn.Exec(createTeamQuery, teamID, fmt.Sprintf("Chat %d Team", chatID), chatID)
        if err != nil {
            return nil, fmt.Errorf("jamoa yaratishda xatolik: %w", err)
        }
    }
    
    query := `
    SELECT id, team_id, user_id, username, role, skills, capacity, current_workload
    FROM team_members 
    WHERE team_id = ?
    ORDER BY role DESC, username ASC`
    
    rows, err := db.conn.Query(query, teamID)
    if err != nil {
        return nil, fmt.Errorf("jamoa a'zolarini olishda xatolik: %w", err)
    }
    defer rows.Close()
    
    var members []TeamMember
    for rows.Next() {
        var member TeamMember
        var skillsStr string
        
        err := rows.Scan(
            &member.ID,
            &member.TeamID,
            &member.UserID,
            &member.Username,
            &member.Role,
            &skillsStr,
            &member.Capacity,
            &member.Current,
        )
        if err != nil {
            return nil, fmt.Errorf("jamoa a'zosi ma'lumotlarini o'qishda xatolik: %w", err)
        }
        
        // Parse skills
        if skillsStr != "" {
            member.Skills = strings.Split(skillsStr, ",")
        }
        
        members = append(members, member)
    }
    
    return members, nil
}

func (db *DB) GetProjectStats(projectID string) (*ProjectStats, error) {
    query := `
    SELECT 
        COUNT(*) as total_tasks,
        COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_tasks,
        SUM(estimate_hours) as estimated_hours,
        SUM(actual_hours) as actual_hours
    FROM tasks 
    WHERE project_id = ?`
    
    var stats ProjectStats
    err := db.conn.QueryRow(query, projectID).Scan(
        &stats.TotalTasks,
        &stats.CompletedTasks,
        &stats.EstimatedHours,
        &stats.ActualHours,
    )
    
    if err != nil {
        return nil, fmt.Errorf("loyiha statistikasini olishda xatolik: %w", err)
    }
    
    stats.ProjectID = projectID
    if stats.TotalTasks > 0 {
        stats.Progress = float64(stats.CompletedTasks) / float64(stats.TotalTasks)
    }
    
    if stats.EstimatedHours > 0 {
        stats.EfficiencyRatio = stats.ActualHours / stats.EstimatedHours
    }
    
    return &stats, nil
}

func (db *DB) Close() error {
    return db.conn.Close()
}