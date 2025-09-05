# 📅 Week 3: Advanced Bot Features & Database Integration

## 🎯 Learning Objectives

By the end of Week 3, you will master:
- Database integration with both SQLite and PostgreSQL
- User management and activity tracking
- Advanced command routing and architecture
- Environment variable management
- Production deployment patterns
- Code organization and modular design

## 📊 Week Overview

| Day | Focus | Key Features | Database Integration |
|-----|-------|-------------|---------------------|
| 15 | Database Foundation | SQLite setup, basic queries | Local SQLite database |
| 16 | User Management | User registration, profiles | User data persistence |
| 17 | Activity Tracking | Command logging, statistics | Activity analytics |
| 18 | Multi-DB Support | PostgreSQL integration | Database abstraction |
| 19 | Advanced Commands | Statistics, user commands | Query optimization |
| 20 | Production Setup | Environment configuration | Production database |
| 21 | Testing & Refactoring | Code cleanup, testing | Database testing |

---

## 📅 Day 15: Database Foundation

### 🎯 Goals
- Set up SQLite database for local development
- Create user and activity tables
- Implement basic database operations
- Add user registration functionality

### 🔧 Database Setup Implementation

#### `database/db.go` - SQLite Implementation

```go
package database

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    
    _ "github.com/mattn/go-sqlite3"
)

type DB struct {
    conn *sql.DB
}

type User struct {
    ID        int64     `json:"id"`
    TelegramID int64    `json:"telegram_id"`
    Username  string    `json:"username"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type UserActivity struct {
    ID        int64     `json:"id"`
    UserID    int64     `json:"user_id"`
    Command   string    `json:"command"`
    CreatedAt time.Time `json:"created_at"`
}
```

**Code Explanation:**

1. **Brief Summary**: This code sets up the foundation for database operations using SQLite, defining data structures for users and their activities.

2. **Step-by-Step Breakdown**:
   - **Package Declaration**: Defines the database package for modular organization
   - **Import Statements**: Includes necessary packages for SQL operations and SQLite driver
   - **DB Struct**: Wrapper around sql.DB for our database operations
   - **User Struct**: Represents user data with JSON tags for API responses
   - **UserActivity Struct**: Tracks user command usage and timestamps

3. **Key Programming Concepts**:
   - Struct tags for JSON serialization
   - Database driver imports with blank identifier
   - Type definitions for data modeling
   - Time handling for timestamps

4. **Complexity Level**: **Beginner** - Basic struct definitions and imports

5. **Suggestions for Improvement**:
   - Add validation tags for struct fields
   - Consider using UUID for IDs instead of int64
   - Add database connection pooling configuration

6. **Related Examples**: 
   - PostgreSQL implementation
   - NoSQL database alternatives
   - ORM frameworks like GORM

#### Database Connection and Initialization

```go
func NewDB() (*DB, error) {
    conn, err := sql.Open("sqlite3", "./yordamchi.db")
    if err != nil {
        return nil, fmt.Errorf("ma'lumotlar bazasini ochishda xatolik: %w", err)
    }

    db := &DB{conn: conn}
    
    if err := db.createTables(); err != nil {
        return nil, fmt.Errorf("jadvallar yaratishda xatolik: %w", err)
    }

    log.Println("✅ SQLite ma'lumotlar bazasi muvaffaqiyatli ulandi")
    return db, nil
}

func (db *DB) createTables() error {
    userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        telegram_id INTEGER UNIQUE NOT NULL,
        username TEXT,
        first_name TEXT NOT NULL,
        last_name TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

    activityTable := `
    CREATE TABLE IF NOT EXISTS user_activities (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        command TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (telegram_id)
    );`

    if _, err := db.conn.Exec(userTable); err != nil {
        return fmt.Errorf("users jadvali yaratishda xatolik: %w", err)
    }

    if _, err := db.conn.Exec(activityTable); err != nil {
        return fmt.Errorf("user_activities jadvali yaratishda xatolik: %w", err)
    }

    return nil
}
```

**Code Explanation:**

1. **Brief Summary**: Establishes database connection and creates necessary tables for user management and activity tracking.

2. **Step-by-Step Breakdown**:
   - **Database Opening**: Opens SQLite connection with file-based storage
   - **Error Handling**: Wraps errors with context for better debugging
   - **Table Creation**: Defines schema for users and activities with proper constraints
   - **Foreign Key Relationships**: Links activities to users via telegram_id

3. **Key Programming Concepts**:
   - Error wrapping with fmt.Errorf
   - SQL DDL (Data Definition Language) statements
   - Database constraints and relationships
   - Constructor pattern for database initialization

4. **Complexity Level**: **Intermediate** - Database schema design and SQL execution

5. **Suggestions for Improvement**:
   - Add database migrations system
   - Implement connection retry logic
   - Add database health checks
   - Use prepared statements for better security

#### User Management Operations

```go
func (db *DB) CreateOrUpdateUser(telegramID int64, username, firstName, lastName string) error {
    query := `
    INSERT INTO users (telegram_id, username, first_name, last_name) 
    VALUES (?, ?, ?, ?)
    ON CONFLICT(telegram_id) DO UPDATE SET
        username = excluded.username,
        first_name = excluded.first_name,
        last_name = excluded.last_name,
        updated_at = CURRENT_TIMESTAMP;`

    _, err := db.conn.Exec(query, telegramID, username, firstName, lastName)
    if err != nil {
        return fmt.Errorf("foydalanuvchini saqlashda xatolik: %w", err)
    }

    return nil
}

func (db *DB) GetUser(telegramID int64) (*User, error) {
    query := `
    SELECT id, telegram_id, username, first_name, last_name, created_at, updated_at 
    FROM users WHERE telegram_id = ?`

    var user User
    err := db.conn.QueryRow(query, telegramID).Scan(
        &user.ID,
        &user.TelegramID,
        &user.Username,
        &user.FirstName,
        &user.LastName,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("foydalanuvchi topilmadi: %d", telegramID)
        }
        return nil, fmt.Errorf("foydalanuvchini olishda xatolik: %w", err)
    }

    return &user, nil
}
```

**Code Explanation:**

1. **Brief Summary**: Implements CRUD operations for user management with upsert functionality and proper error handling.

2. **Step-by-Step Breakdown**:
   - **Upsert Operation**: Uses ON CONFLICT to insert or update existing users
   - **Parameterized Queries**: Prevents SQL injection with placeholder parameters
   - **Row Scanning**: Maps database columns to struct fields
   - **Error Classification**: Distinguishes between not found and actual errors

3. **Key Programming Concepts**:
   - SQL UPSERT (INSERT ... ON CONFLICT)
   - Parameterized queries for security
   - Pointer receivers for methods
   - Error handling and classification

4. **Complexity Level**: **Intermediate** - Advanced SQL operations and error handling

5. **Suggestions for Improvement**:
   - Add bulk insert operations
   - Implement user search functionality
   - Add data validation before database operations
   - Consider using transactions for complex operations

---

## 📅 Day 16: User Activity Tracking

### 🎯 Goals
- Implement activity logging system
- Create statistics gathering functions
- Add user engagement tracking
- Build analytics foundations

#### Activity Tracking Implementation

```go
func (db *DB) LogUserActivity(userID int64, command string) error {
    query := `INSERT INTO user_activities (user_id, command) VALUES (?, ?)`
    
    _, err := db.conn.Exec(query, userID, command)
    if err != nil {
        return fmt.Errorf("faollik yozishda xatolik: %w", err)
    }

    return nil
}

func (db *DB) GetUserStats() (int, error) {
    query := `SELECT COUNT(*) FROM users`
    
    var count int
    err := db.conn.QueryRow(query).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("statistika olishda xatolik: %w", err)
    }

    return count, nil
}

func (db *DB) GetUserActivities(userID int64, limit int) ([]UserActivity, error) {
    query := `
    SELECT id, user_id, command, created_at 
    FROM user_activities 
    WHERE user_id = ? 
    ORDER BY created_at DESC 
    LIMIT ?`

    rows, err := db.conn.Query(query, userID, limit)
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
```

**Code Explanation:**

1. **Brief Summary**: Implements comprehensive activity tracking system for monitoring user interactions and generating statistics.

2. **Step-by-Step Breakdown**:
   - **Activity Logging**: Records each command execution with timestamp
   - **Statistics Gathering**: Provides user count and other metrics
   - **Activity History**: Retrieves user's recent command history
   - **Result Processing**: Handles multiple rows and proper resource cleanup

3. **Key Programming Concepts**:
   - SQL aggregation functions (COUNT)
   - Result set iteration with rows.Next()
   - Resource management with defer
   - Slice building from database results

4. **Complexity Level**: **Intermediate** - Database queries with result processing

5. **Suggestions for Improvement**:
   - Add time-based filtering for activities
   - Implement activity analytics (most used commands)
   - Add user engagement scoring
   - Create activity archiving for old records

---

## 📅 Day 17-18: Multi-Database Support

### 🎯 Goals
- Abstract database operations for multiple providers
- Implement PostgreSQL support
- Create environment-based database selection
- Ensure compatibility between database types

#### `database/postgres.go` - PostgreSQL Implementation

```go
package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    
    _ "github.com/lib/pq"
)

func NewPostgresDB() (*DB, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return nil, fmt.Errorf("DATABASE_URL environment variable topilmadi")
    }

    conn, err := sql.Open("postgres", dbURL)
    if err != nil {
        return nil, fmt.Errorf("PostgreSQL ma'lumotlar bazasini ochishda xatolik: %w", err)
    }

    // Test connection
    if err := conn.Ping(); err != nil {
        return nil, fmt.Errorf("PostgreSQL serveriga ulanishda xatolik: %w", err)
    }

    db := &DB{conn: conn}
    
    if err := db.createPostgresTables(); err != nil {
        return nil, fmt.Errorf("PostgreSQL jadvallar yaratishda xatolik: %w", err)
    }

    log.Println("✅ PostgreSQL ma'lumotlar bazasi muvaffaqiyatli ulandi")
    return db, nil
}

func (db *DB) createPostgresTables() error {
    userTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        telegram_id BIGINT UNIQUE NOT NULL,
        username VARCHAR(255),
        first_name VARCHAR(255) NOT NULL,
        last_name VARCHAR(255),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

    activityTable := `
    CREATE TABLE IF NOT EXISTS user_activities (
        id SERIAL PRIMARY KEY,
        user_id BIGINT NOT NULL,
        command TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (telegram_id)
    );`

    if _, err := db.conn.Exec(userTable); err != nil {
        return fmt.Errorf("users jadvali yaratishda xatolik: %w", err)
    }

    if _, err := db.conn.Exec(activityTable); err != nil {
        return fmt.Errorf("user_activities jadvali yaratishda xatolik: %w", err)
    }

    return nil
}
```

**Code Explanation:**

1. **Brief Summary**: Provides PostgreSQL-specific implementation while maintaining the same interface as SQLite for database operations.

2. **Step-by-Step Breakdown**:
   - **Environment Configuration**: Reads connection string from environment variables
   - **Connection Testing**: Validates database connectivity with Ping()
   - **Schema Differences**: Handles PostgreSQL-specific data types (SERIAL, BIGINT)
   - **Type Compatibility**: Maintains same struct interface for both databases

3. **Key Programming Concepts**:
   - Database abstraction and polymorphism
   - Environment variable configuration
   - Database driver abstraction
   - SQL dialect differences handling

4. **Complexity Level**: **Advanced** - Multi-database architecture and abstraction

5. **Suggestions for Improvement**:
   - Add connection pooling configuration
   - Implement database health monitoring
   - Add migration version tracking
   - Create database performance monitoring

---

## 📅 Day 19-21: Production Integration & Testing

### 🎯 Goals
- Integrate database with bot functionality
- Add comprehensive error handling
- Implement production logging
- Create testing framework

#### Updated Bot Integration

The updated `main.go:32-46` shows the production integration:

```go
// Ma'lumotlar bazasi turini aniqlash
var db *database.DB
dbType := os.Getenv("DB_TYPE")

switch dbType {
case "postgres":
    db, err = database.NewPostgresDB()
default:
    db, err = database.NewDB() // SQLite
}

if err != nil {
    log.Fatal("Ma'lumotlar bazasi xatoligi:", err)
}
defer db.Close()
```

**Code Explanation:**

1. **Brief Summary**: Production-ready database selection based on environment configuration with proper error handling and resource cleanup.

2. **Step-by-Step Breakdown**:
   - **Environment Detection**: Uses DB_TYPE to determine database provider
   - **Factory Pattern**: Creates appropriate database instance based on type
   - **Error Propagation**: Properly handles initialization errors
   - **Resource Management**: Ensures database connections are closed

3. **Key Programming Concepts**:
   - Factory pattern for object creation
   - Environment-driven configuration
   - Proper resource cleanup with defer
   - Production error handling

4. **Complexity Level**: **Intermediate** - Production patterns and configuration

Bot Integration in `bot.go:100-117`:

```go
// Foydalanuvchini ma'lumotlar bazasida saqlash
err := b.DB.CreateOrUpdateUser(
    int64(msg.From.ID),
    msg.From.Username,
    msg.From.FirstName,
    msg.From.LastName,
)
if err != nil {
    log.Printf("Foydalanuvchini saqlashda xatolik: %v", err)
}

// Faollikni yozish
if strings.HasPrefix(msg.Text, "/") {
    err = b.DB.LogUserActivity(int64(msg.From.ID), msg.Text)
    if err != nil {
        log.Printf("Faollik yozishda xatolik: %v", err)
    }
}
```

**Code Explanation:**

1. **Brief Summary**: Seamless integration of database operations within message processing flow with proper error handling and logging.

2. **Step-by-Step Breakdown**:
   - **User Registration**: Automatically registers or updates users on message
   - **Activity Tracking**: Logs command usage for analytics
   - **Error Resilience**: Continues operation despite database errors
   - **Conditional Logging**: Only tracks command activities (starting with /)

3. **Key Programming Concepts**:
   - Graceful error handling
   - Automated user lifecycle management
   - Activity pattern detection
   - Non-blocking error handling

4. **Complexity Level**: **Intermediate** - Production integration patterns

5. **Suggestions for Improvement**:
   - Add database operation retries
   - Implement circuit breaker pattern
   - Add metrics collection
   - Create database operation queuing for high load

## 🏆 Week 3 Achievements

By the end of Week 3, you have successfully implemented:

✅ **Database Foundation**: SQLite and PostgreSQL support with proper schema design
✅ **User Management**: Automatic user registration and profile management  
✅ **Activity Tracking**: Command logging and statistics gathering
✅ **Multi-DB Architecture**: Environment-based database selection
✅ **Production Integration**: Proper error handling and resource management
✅ **Statistics System**: User count and activity analytics

## 🚀 Next Steps for Week 4

Week 4 will focus on:
- Advanced bot features and command expansion
- API integrations and external services
- Performance optimization and caching
- Advanced error handling and monitoring
- Testing framework and quality assurance
- Deployment automation and CI/CD