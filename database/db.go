package database

import (
    "database/sql"
    "fmt"
    "log"
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

func (db *DB) Close() error {
    return db.conn.Close()
}