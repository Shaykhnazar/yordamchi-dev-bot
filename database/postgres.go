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
    `

    _, err := db.conn.Exec(query)
    return err
}

