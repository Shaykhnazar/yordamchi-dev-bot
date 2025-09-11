package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
_   "github.com/lib/pq"
)

func NewPostgresDB() (*DB, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        // Local development uchun
        dbURL = "postgres://default:secret@localhost/yordamchi_bot?sslmode=disable"
        log.Println("⚠️ DATABASE_URL topilmadi, local DB ishlatilmoqda")
    }

    conn, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, fmt.Errorf("PostgreSQL'ga ulanishda xatolik: %w", err)
    }

    // Connection'ni test qilish
    if err := conn.Ping(); err != nil {
        return nil, fmt.Errorf("PostgreSQL ping xatoligi: %w", err)
    }

    db := &DB{conn: conn}
    
    if err := db.createPostgresTables(); err != nil {
        return nil, fmt.Errorf("PostgreSQL jadvallar yaratishda xatolik: %w", err)
    }

    log.Println("✅ PostgreSQL muvaffaqiyatli sozlandi")
    return db, nil
}

func (db *DB) createPostgresTables() error {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        telegram_id BIGINT UNIQUE NOT NULL,
        username VARCHAR(255),
        first_name VARCHAR(255),
        last_name VARCHAR(255),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS user_activity (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
        command TEXT,
        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS teams (
        id TEXT PRIMARY KEY,
        chat_id BIGINT UNIQUE NOT NULL,
        name TEXT NOT NULL,
        description TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS projects (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT,
        team_id TEXT NOT NULL,
        status TEXT DEFAULT 'active',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (team_id) REFERENCES teams (id)
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
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (team_id) REFERENCES teams (id)
    );

    CREATE TABLE IF NOT EXISTS tasks (
        id TEXT PRIMARY KEY,
        project_id TEXT NOT NULL,
        title TEXT NOT NULL,
        description TEXT,
        category TEXT,
        estimate_hours REAL DEFAULT 0.0,
        actual_hours REAL DEFAULT 0.0,
        status TEXT DEFAULT 'todo',
        priority TEXT DEFAULT 'medium',
        assigned_to TEXT,
        dependencies TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        completed_at TIMESTAMP,
        FOREIGN KEY (project_id) REFERENCES projects (id),
        FOREIGN KEY (assigned_to) REFERENCES team_members (id)
    );
    `

    _, err := db.conn.Exec(query)
    return err
}

