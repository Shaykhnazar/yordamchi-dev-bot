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
    VALUES (?, ?, ?, ?)
    ON CONFLICT(telegram_id) DO UPDATE SET
        username = excluded.username,
        first_name = excluded.first_name,
        last_name = excluded.last_name,
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
    userIDQuery := "SELECT id FROM users WHERE telegram_id = ?"
    var userID int
    err := db.conn.QueryRow(userIDQuery, telegramID).Scan(&userID)
    if err != nil {
        return fmt.Errorf("foydalanuvchi ID topilmadi: %w", err)
    }

    activityQuery := "INSERT INTO user_activity (user_id, command) VALUES (?, ?)"
    _, err = db.conn.Exec(activityQuery, userID, command)
    if err != nil {
        return fmt.Errorf("faollik yozishda xatolik: %w", err)
    }

    return nil
}

func (db *DB) GetUserStats() (int, error) {
    query := "SELECT COUNT(*) FROM users"
    var count int
    err := db.conn.QueryRow(query).Scan(&count)
    return count, err
}

func (db *DB) Close() error {
    return db.conn.Close()
}