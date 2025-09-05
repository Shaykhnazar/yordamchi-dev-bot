# 📅 Week 1 Complete Implementation Guide

## 📅 Day 4: Code Structure & Handler Organization

### 🎯 Goals
- Implement clean code architecture
- Separate concerns with proper package structure
- Create reusable command handlers
- Implement middleware pattern

### 📁 Enhanced Project Structure
```
yordamchi-dev-bot/
├── cmd/
│   └── bot/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── bot.go
│   │   └── dependencies.go
│   ├── handlers/
│   │   ├── commands/
│   │   │   ├── start.go
│   │   │   ├── help.go
│   │   │   ├── ping.go
│   │   │   └── echo.go
│   │   └── middleware/
│   │       ├── auth.go
│   │       ├── logging.go
│   │       └── ratelimit.go
│   ├── services/
│   │   └── user/
│   │       ├── service.go
│   │       └── repository.go
│   └── domain/
│       ├── user.go
│       └── command.go
├── configs/
│   ├── bot.json
│   └── messages/
│       ├── en.json
│       └── uz.json
└── go.mod
```

### 🔧 Implementation

#### `internal/domain/command.go` - Domain Models
```go
package domain

import (
    "context"
    "time"
)

// Command represents a user command
type Command struct {
    ID        string
    Text      string
    User      *User
    Chat      *Chat
    Timestamp time.Time
}

// Response represents a bot response
type Response struct {
    Text         string
    ParseMode    string
    ReplyMarkup  interface{}
    DisablePreview bool
}

// CommandHandler defines the interface for command handling
type CommandHandler interface {
    Handle(ctx context.Context, cmd *Command) (*Response, error)
    CanHandle(command string) bool
    Description() string
    Usage() string
}

// Middleware defines the interface for processing pipeline
type Middleware interface {
    Process(ctx context.Context, next HandlerFunc) HandlerFunc
}

type HandlerFunc func(ctx context.Context, cmd *Command) (*Response, error)
```

#### `internal/domain/user.go` - User Entity
```go
package domain

import (
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
    IsBlocked   bool      `json:"is_blocked"`
    Preferences UserPrefs `json:"preferences"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// UserPrefs represents user preferences
type UserPrefs struct {
    Notifications bool   `json:"notifications"`
    Theme         string `json:"theme"`
    Timezone      string `json:"timezone"`
}

// Chat represents a Telegram chat
type Chat struct {
    ID   int64  `json:"id"`
    Type string `json:"type"`
}
```

#### `internal/handlers/commands/start.go` - Start Command Handler
```go
package commands

import (
    "context"
    "fmt"
    "strings"
    
    "yordamchi-dev-bot/internal/domain"
    "yordamchi-dev-bot/internal/services/user"
)

type StartHandler struct {
    userService user.Service
    messages    map[string]interface{}
}

func NewStartHandler(userService user.Service, messages map[string]interface{}) *StartHandler {
    return &StartHandler{
        userService: userService,
        messages:    messages,
    }
}

func (h *StartHandler) CanHandle(command string) bool {
    return strings.ToLower(command) == "/start"
}

func (h *StartHandler) Description() string {
    return "Start the bot and show welcome message"
}

func (h *StartHandler) Usage() string {
    return "/start - Initialize bot conversation"
}

func (h *StartHandler) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
    // Get or create user
    user, err := h.userService.GetOrCreateUser(ctx, &domain.User{
        TelegramID: cmd.User.TelegramID,
        Username:   cmd.User.Username,
        FirstName:  cmd.User.FirstName,
        LastName:   cmd.User.LastName,
        Language:   "en", // Default language
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get or create user: %w", err)
    }

    // Get localized welcome message
    welcomeMsg := h.getLocalizedMessage(user.Language, "commands.start.welcome", map[string]string{
        "name": cmd.User.FirstName,
    })

    return &domain.Response{
        Text:      welcomeMsg,
        ParseMode: "HTML",
    }, nil
}

func (h *StartHandler) getLocalizedMessage(lang, key string, params map[string]string) string {
    // Simplified message retrieval - in real implementation, use proper i18n
    template := `🎉 Welcome, {name}!

I'm DevMate Bot, your smart developer assistant built with Go!

🚀 Type /help to see what I can do for you.`

    message := template
    for key, value := range params {
        placeholder := fmt.Sprintf("{%s}", key)
        message = strings.ReplaceAll(message, placeholder, value)
    }

    return message
}
```

#### `internal/handlers/commands/help.go` - Help Command Handler
```go
package commands

import (
    "context"
    "strings"
    
    "yordamchi-dev-bot/internal/domain"
)

type HelpHandler struct {
    handlers []domain.CommandHandler
    messages map[string]interface{}
}

func NewHelpHandler(handlers []domain.CommandHandler, messages map[string]interface{}) *HelpHandler {
    return &HelpHandler{
        handlers: handlers,
        messages: messages,
    }
}

func (h *HelpHandler) CanHandle(command string) bool {
    return strings.ToLower(command) == "/help"
}

func (h *HelpHandler) Description() string {
    return "Show available commands and their usage"
}

func (h *HelpHandler) Usage() string {
    return "/help - Display this help message"
}

func (h *HelpHandler) Handle(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
    var helpText strings.Builder
    
    helpText.WriteString("🤖 <b>DevMate Bot Commands</b>\n\n")
    helpText.WriteString("<b>Available Commands:</b>\n")
    
    for _, handler := range h.handlers {
        helpText.WriteString(fmt.Sprintf("• %s - %s\n", 
            handler.Usage(), 
            handler.Description()))
    }
    
    helpText.WriteString("\n<i>More features coming soon! 🚀</i>")

    return &domain.Response{
        Text:      helpText.String(),
        ParseMode: "HTML",
    }, nil
}
```

#### `internal/handlers/middleware/logging.go` - Logging Middleware
```go
package middleware

import (
    "context"
    "log"
    "time"
    
    "yordamchi-dev-bot/internal/domain"
)

type LoggingMiddleware struct {
    logger *log.Logger
}

func NewLoggingMiddleware(logger *log.Logger) *LoggingMiddleware {
    return &LoggingMiddleware{logger: logger}
}

func (m *LoggingMiddleware) Process(ctx context.Context, next domain.HandlerFunc) domain.HandlerFunc {
    return func(ctx context.Context, cmd *domain.Command) (*domain.Response, error) {
        start := time.Now()
        
        // Log incoming command
        m.logger.Printf("📥 Command received: %s from user %s (@%s)", 
            cmd.Text, 
            cmd.User.FirstName, 
            cmd.User.Username)
        
        // Execute next handler
        response, err := next(ctx, cmd)
        
        duration := time.Since(start)
        
        // Log result
        if err != nil {
            m.logger.Printf("❌ Command failed: %s (took %v) - Error: %v", 
                cmd.Text, duration, err)
        } else {
            m.logger.Printf("✅ Command completed: %s (took %v)", 
                cmd.Text, duration)
        }
        
        return response, err
    }
}
```

#### `internal/app/bot.go` - Main Bot Application
```go
package app

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
    
    "yordamchi-dev-bot/internal/domain"
    "yordamchi-dev-bot/internal/handlers/commands"
    "yordamchi-dev-bot/internal/handlers/middleware"
)

type Bot struct {
    token       string
    baseURL     string
    client      *http.Client
    handlers    []domain.CommandHandler
    middleware  []domain.Middleware
    logger      *log.Logger
}

func NewBot(deps *Dependencies) *Bot {
    bot := &Bot{
        token:   deps.Config.Bot.Token,
        baseURL: "https://api.telegram.org/bot" + deps.Config.Bot.Token,
        client:  &http.Client{Timeout: 30 * time.Second},
        logger:  deps.Logger,
    }
    
    // Initialize handlers
    bot.initializeHandlers(deps)
    
    // Initialize middleware
    bot.initializeMiddleware(deps)
    
    return bot
}

func (b *Bot) initializeHandlers(deps *Dependencies) {
    b.handlers = []domain.CommandHandler{
        commands.NewStartHandler(deps.UserService, deps.Messages),
        commands.NewHelpHandler(nil, deps.Messages), // Will be updated with all handlers
        commands.NewPingHandler(),
        commands.NewEchoHandler(),
    }
    
    // Update help handler with all handlers
    if len(b.handlers) > 1 {
        if helpHandler, ok := b.handlers[1].(*commands.HelpHandler); ok {
            helpHandler.SetHandlers(b.handlers)
        }
    }
}

func (b *Bot) initializeMiddleware(deps *Dependencies) {
    b.middleware = []domain.Middleware{
        middleware.NewLoggingMiddleware(b.logger),
        middleware.NewAuthMiddleware(deps.UserService),
        middleware.NewRateLimitMiddleware(deps.RateLimiter),
    }
}

func (b *Bot) Start() error {
    http.HandleFunc("/webhook", b.handleWebhook)
    http.HandleFunc("/health", b.handleHealth)
    
    b.logger.Println("🚀 DevMate Bot started on port 8080")
    return http.ListenAndServe(":8080", nil)
}

func (b *Bot) handleWebhook(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    body, err := io.ReadAll(r.Body)
    if err != nil {
        b.logger.Printf("❌ Error reading body: %v", err)
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()
    
    var update TelegramUpdate
    if err := json.Unmarshal(body, &update); err != nil {
        b.logger.Printf("❌ JSON parse error: %v", err)
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Process message asynchronously
    go b.processUpdate(context.Background(), &update)
    
    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, "OK")
}

func (b *Bot) processUpdate(ctx context.Context, update *TelegramUpdate) {
    if update.Message == nil || update.Message.Text == "" {
        return
    }
    
    // Convert to domain command
    cmd := &domain.Command{
        Text: update.Message.Text,
        User: &domain.User{
            TelegramID: int64(update.Message.From.ID),
            Username:   update.Message.From.Username,
            FirstName:  update.Message.From.FirstName,
            LastName:   update.Message.From.LastName,
        },
        Chat: &domain.Chat{
            ID:   update.Message.Chat.ID,
            Type: update.Message.Chat.Type,
        },
        Timestamp: time.Now(),
    }
    
    // Find appropriate handler
    var handler domain.CommandHandler
    for _, h := range b.handlers {
        if h.CanHandle(cmd.Text) {
            handler = h
            break
        }
    }
    
    if handler == nil {
        b.sendErrorResponse(cmd.Chat.ID, "❓ Unknown command. Type /help for available commands.")
        return
    }
    
    // Create handler function with middleware chain
    handlerFunc := handler.Handle
    
    // Apply middleware in reverse order
    for i := len(b.middleware) - 1; i >= 0; i-- {
        handlerFunc = b.middleware[i].Process(ctx, handlerFunc)
    }
    
    // Execute handler with middleware
    response, err := handlerFunc(ctx, cmd)
    if err != nil {
        b.logger.Printf("❌ Handler error: %v", err)
        b.sendErrorResponse(cmd.Chat.ID, "⚠️ An error occurred while processing your request.")
        return
    }
    
    // Send response
    if err := b.sendMessage(cmd.Chat.ID, response); err != nil {
        b.logger.Printf("❌ Failed to send message: %v", err)
    }
}

func (b *Bot) sendMessage(chatID int64, response *domain.Response) error {
    payload := map[string]interface{}{
        "chat_id":    chatID,
        "text":       response.Text,
        "parse_mode": response.ParseMode,
    }
    
    if response.DisablePreview {
        payload["disable_web_page_preview"] = true
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("marshal payload: %w", err)
    }
    
    url := fmt.Sprintf("%s/sendMessage", b.baseURL)
    resp, err := b.client.Post(url, "application/json", strings.NewReader(string(jsonData)))
    if err != nil {
        return fmt.Errorf("send HTTP request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("telegram API error: %d - %s", resp.StatusCode, body)
    }
    
    return nil
}

func (b *Bot) sendErrorResponse(chatID int64, message string) {
    response := &domain.Response{
        Text:      message,
        ParseMode: "HTML",
    }
    
    if err := b.sendMessage(chatID, response); err != nil {
        b.logger.Printf("❌ Failed to send error message: %v", err)
    }
}

// Telegram API types
type TelegramUpdate struct {
    UpdateID int              `json:"update_id"`
    Message  *TelegramMessage `json:"message"`
}

type TelegramMessage struct {
    MessageID int           `json:"message_id"`
    From      *TelegramUser `json:"from"`
    Chat      *TelegramChat `json:"chat"`
    Text      string        `json:"text"`
    Date      int64         `json:"date"`
}

type TelegramUser struct {
    ID        int    `json:"id"`
    IsBot     bool   `json:"is_bot"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Username  string `json:"username"`
}

type TelegramChat struct {
    ID   int64  `json:"id"`
    Type string `json:"type"`
}
```

---

## 📅 Day 5: SQLite Database Implementation

### 🎯 Goals
- Implement SQLite database layer
- Create user repository pattern
- Add database migrations
- Implement CRUD operations

### 🔧 Database Implementation

#### `internal/infrastructure/database/sqlite.go`
```go
package database

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    
    _ "github.com/mattn/go-sqlite3"
    
    "yordamchi-dev-bot/internal/domain"
)

type SQLiteDB struct {
    db     *sql.DB
    logger Logger
}

type Logger interface {
    Printf(format string, args ...interface{})
    Println(args ...interface{})
}

func NewSQLiteDB(dbPath string, logger Logger) (*SQLiteDB, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, fmt.Errorf("open database: %w", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("ping database: %w", err)
    }
    
    sqliteDB := &SQLiteDB{
        db:     db,
        logger: logger,
    }
    
    // Run migrations
    if err := sqliteDB.migrate(); err != nil {
        return nil, fmt.Errorf("migrate database: %w", err)
    }
    
    logger.Println("✅ SQLite database connected successfully")
    return sqliteDB, nil
}

func (db *SQLiteDB) migrate() error {
    migrations := []string{
        `CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            telegram_id INTEGER UNIQUE NOT NULL,
            username TEXT,
            first_name TEXT,
            last_name TEXT,
            language TEXT DEFAULT 'en',
            is_blocked BOOLEAN DEFAULT FALSE,
            preferences TEXT, -- JSON string
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )`,
        
        `CREATE TABLE IF NOT EXISTS user_sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            session_data TEXT, -- JSON string
            expires_at DATETIME,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users (id)
        )`,
        
        `CREATE TABLE IF NOT EXISTS command_history (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            command TEXT NOT NULL,
            parameters TEXT, -- JSON string
            response_type TEXT,
            executed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users (id)
        )`,
        
        `CREATE TABLE IF NOT EXISTS api_usage (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            endpoint TEXT NOT NULL,
            requests_count INTEGER DEFAULT 1,
            date DATE DEFAULT (date('now')),
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (user_id) REFERENCES users (id)
        )`,
        
        `CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id)`,
        `CREATE INDEX IF NOT EXISTS idx_command_history_user_id ON command_history(user_id)`,
        `CREATE INDEX IF NOT EXISTS idx_api_usage_user_date ON api_usage(user_id, date)`,
    }
    
    for _, migration := range migrations {
        if _, err := db.db.Exec(migration); err != nil {
            return fmt.Errorf("execute migration: %w", err)
        }
    }
    
    db.logger.Println("✅ Database migrations completed")
    return nil
}

// UserRepository implements user data operations
type UserRepository struct {
    db     *sql.DB
    logger Logger
}

func NewUserRepository(db *SQLiteDB) *UserRepository {
    return &UserRepository{
        db:     db.db,
        logger: db.logger,
    }
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (telegram_id, username, first_name, last_name, language, preferences)
        VALUES (?, ?, ?, ?, ?, ?)
    `
    
    prefsJSON := "{}" // Default empty preferences
    if user.Preferences != (domain.UserPrefs{}) {
        // In real implementation, properly marshal preferences to JSON
        prefsJSON = fmt.Sprintf(`{"notifications":%t,"theme":"%s","timezone":"%s"}`, 
            user.Preferences.Notifications, 
            user.Preferences.Theme, 
            user.Preferences.Timezone)
    }
    
    result, err := r.db.ExecContext(ctx, query,
        user.TelegramID,
        user.Username,
        user.FirstName,
        user.LastName,
        user.Language,
        prefsJSON,
    )
    if err != nil {
        return fmt.Errorf("insert user: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return fmt.Errorf("get last insert ID: %w", err)
    }
    
    user.ID = id
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    
    r.logger.Printf("📝 User created: ID=%d, TelegramID=%d, Name=%s", 
        user.ID, user.TelegramID, user.FirstName)
    
    return nil
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
    query := `
        SELECT id, telegram_id, username, first_name, last_name, language, 
               is_blocked, preferences, created_at, updated_at
        FROM users 
        WHERE telegram_id = ?
    `
    
    var user domain.User
    var prefsJSON string
    
    err := r.db.QueryRowContext(ctx, query, telegramID).Scan(
        &user.ID,
        &user.TelegramID,
        &user.Username,
        &user.FirstName,
        &user.LastName,
        &user.Language,
        &user.IsBlocked,
        &prefsJSON,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found: telegram_id=%d", telegramID)
        }
        return nil, fmt.Errorf("query user: %w", err)
    }
    
    // Parse preferences JSON (simplified)
    // In real implementation, use proper JSON unmarshaling
    user.Preferences = domain.UserPrefs{
        Notifications: true,
        Theme:         "light",
        Timezone:      "UTC",
    }
    
    return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
    query := `
        UPDATE users 
        SET username = ?, first_name = ?, last_name = ?, language = ?, 
            is_blocked = ?, preferences = ?, updated_at = CURRENT_TIMESTAMP
        WHERE id = ?
    `
    
    prefsJSON := fmt.Sprintf(`{"notifications":%t,"theme":"%s","timezone":"%s"}`, 
        user.Preferences.Notifications, 
        user.Preferences.Theme, 
        user.Preferences.Timezone)
    
    result, err := r.db.ExecContext(ctx, query,
        user.Username,
        user.FirstName,
        user.LastName,
        user.Language,
        user.IsBlocked,
        prefsJSON,
        user.ID,
    )
    if err != nil {
        return fmt.Errorf("update user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("user not found: id=%d", user.ID)
    }
    
    user.UpdatedAt = time.Now()
    
    r.logger.Printf("📝 User updated: ID=%d, TelegramID=%d", user.ID, user.TelegramID)
    return nil
}

func (r *UserRepository) GetUserStats(ctx context.Context) (*UserStats, error) {
    query := `
        SELECT 
            COUNT(*) as total_users,
            COUNT(CASE WHEN created_at >= date('now', '-7 days') THEN 1 END) as new_users_week,
            COUNT(CASE WHEN is_blocked = TRUE THEN 1 END) as blocked_users
        FROM users
    `
    
    var stats UserStats
    err := r.db.QueryRowContext(ctx, query).Scan(
        &stats.TotalUsers,
        &stats.NewUsersThisWeek,
        &stats.BlockedUsers,
    )
    
    if err != nil {
        return nil, fmt.Errorf("query user stats: %w", err)
    }
    
    return &stats, nil
}

func (r *UserRepository) LogCommand(ctx context.Context, userID int64, command, responseType string) error {
    query := `
        INSERT INTO command_history (user_id, command, response_type)
        VALUES (?, ?, ?)
    `
    
    _, err := r.db.ExecContext(ctx, query, userID, command, responseType)
    if err != nil {
        return fmt.Errorf("log command: %w", err)
    }
    
    return nil
}

func (r *UserRepository) Close() error {
    return r.db.Close()
}

// UserStats represents user statistics
type UserStats struct {
    TotalUsers       int `json:"total_users"`
    NewUsersThisWeek int `json:"new_users_week"`
    BlockedUsers     int `json:"blocked_users"`
}
```

#### `internal/services/user/service.go` - User Service
```go
package user

import (
    "context"
    "fmt"
    "time"
    
    "yordamchi-dev-bot/internal/domain"
)

type Repository interface {
    Create(ctx context.Context, user *domain.User) error
    GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    GetUserStats(ctx context.Context) (*UserStats, error)
    LogCommand(ctx context.Context, userID int64, command, responseType string) error
}

type Service struct {
    repo   Repository
    logger Logger
}

type Logger interface {
    Printf(format string, args ...interface{})
    Println(args ...interface{})
}

type UserStats struct {
    TotalUsers       int `json:"total_users"`
    NewUsersThisWeek int `json:"new_users_week"`
    BlockedUsers     int `json:"blocked_users"`
}

func NewService(repo Repository, logger Logger) *Service {
    return &Service{
        repo:   repo,
        logger: logger,
    }
}

func (s *Service) GetOrCreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
    // Try to get existing user
    existingUser, err := s.repo.GetByTelegramID(ctx, user.TelegramID)
    if err == nil {
        // User exists, update if needed
        updated := false
        
        if existingUser.Username != user.Username {
            existingUser.Username = user.Username
            updated = true
        }
        
        if existingUser.FirstName != user.FirstName {
            existingUser.FirstName = user.FirstName
            updated = true
        }
        
        if existingUser.LastName != user.LastName {
            existingUser.LastName = user.LastName
            updated = true
        }
        
        if updated {
            if err := s.repo.Update(ctx, existingUser); err != nil {
                s.logger.Printf("❌ Failed to update user: %v", err)
                // Continue with existing user data
            }
        }
        
        return existingUser, nil
    }
    
    // User doesn't exist, create new one
    newUser := &domain.User{
        TelegramID: user.TelegramID,
        Username:   user.Username,
        FirstName:  user.FirstName,
        LastName:   user.LastName,
        Language:   user.Language,
        IsBlocked:  false,
        Preferences: domain.UserPrefs{
            Notifications: true,
            Theme:         "light",
            Timezone:      "UTC",
        },
    }
    
    if err := s.repo.Create(ctx, newUser); err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }
    
    s.logger.Printf("👤 New user registered: %s (@%s)", newUser.FirstName, newUser.Username)
    return newUser, nil
}

func (s *Service) UpdateUserLanguage(ctx context.Context, telegramID int64, language string) error {
    user, err := s.repo.GetByTelegramID(ctx, telegramID)
    if err != nil {
        return fmt.Errorf("get user: %w", err)
    }
    
    user.Language = language
    
    if err := s.repo.Update(ctx, user); err != nil {
        return fmt.Errorf("update user language: %w", err)
    }
    
    s.logger.Printf("🌐 User language updated: ID=%d, Language=%s", user.ID, language)
    return nil
}

func (s *Service) BlockUser(ctx context.Context, telegramID int64) error {
    user, err := s.repo.GetByTelegramID(ctx, telegramID)
    if err != nil {
        return fmt.Errorf("get user: %w", err)
    }
    
    user.IsBlocked = true
    
    if err := s.repo.Update(ctx, user); err != nil {
        return fmt.Errorf("block user: %w", err)
    }
    
    s.logger.Printf("🚫 User blocked: ID=%d, TelegramID=%d", user.ID, user.TelegramID)
    return nil
}

func (s *Service) GetStats(ctx context.Context) (*UserStats, error) {
    return s.repo.GetUserStats(ctx)
}

func (s *Service) LogCommand(ctx context.Context, telegramID int64, command string) error {
    user, err := s.repo.GetByTelegramID(ctx, telegramID)
    if err != nil {
        return fmt.Errorf("get user for logging: %w", err)
    }
    
    return s.repo.LogCommand(ctx, user.ID, command, "success")
}
```

---

## 📅 Day 6: PostgreSQL Production Database

### 🎯 Goals
- Implement PostgreSQL support for production
- Add environment-based database selection
- Implement connection pooling
- Add database health checks

### 🔧 PostgreSQL Implementation

#### `internal/infrastructure/database/postgres.go`
```go
package database

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    
    _ "github.com/lib/pq"
    
    "yordamchi-dev-bot/internal/domain"
)

type PostgresDB struct {
    db     *sql.DB
    logger Logger
}

func NewPostgresDB(connectionString string, logger Logger) (*PostgresDB, error) {
    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        return nil, fmt.Errorf("open postgres database: %w", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("ping postgres database: %w", err)
    }
    
    postgresDB := &PostgresDB{
        db:     db,
        logger: logger,
    }
    
    // Run migrations
    if err := postgresDB.migrate(); err != nil {
        return nil, fmt.Errorf("migrate postgres database: %w", err)
    }
    
    logger.Println("✅ PostgreSQL database connected successfully")
    return postgresDB, nil
}

func (db *PostgresDB) migrate() error {
    migrations := []string{
        `CREATE TABLE IF NOT EXISTS users (
            id BIGSERIAL PRIMARY KEY,
            telegram_id BIGINT UNIQUE NOT NULL,
            username VARCHAR(255),
            first_name VARCHAR(255),
            last_name VARCHAR(255),
            language VARCHAR(10) DEFAULT 'en',
            is_blocked BOOLEAN DEFAULT FALSE,
            preferences JSONB DEFAULT '{}',
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )`,
        
        `CREATE TABLE IF NOT EXISTS user_sessions (
            id BIGSERIAL PRIMARY KEY,
            user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            session_data JSONB DEFAULT '{}',
            expires_at TIMESTAMP WITH TIME ZONE,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )`,
        
        `CREATE TABLE IF NOT EXISTS command_history (
            id BIGSERIAL PRIMARY KEY,
            user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            command TEXT NOT NULL,
            parameters JSONB DEFAULT '{}',
            response_type VARCHAR(50),
            execution_time_ms INTEGER,
            executed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )`,
        
        `CREATE TABLE IF NOT EXISTS api_usage (
            id BIGSERIAL PRIMARY KEY,
            user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            endpoint VARCHAR(255) NOT NULL,
            requests_count INTEGER DEFAULT 1,
            date DATE DEFAULT CURRENT_DATE,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            UNIQUE(user_id, endpoint, date)
        )`,
        
        // Indexes for performance
        `CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id)`,
        `CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at)`,
        `CREATE INDEX IF NOT EXISTS idx_command_history_user_id_executed_at ON command_history(user_id, executed_at)`,
        `CREATE INDEX IF NOT EXISTS idx_api_usage_user_date ON api_usage(user_id, date)`,
        `CREATE INDEX IF NOT EXISTS idx_api_usage_endpoint_date ON api_usage(endpoint, date)`,
        
        // Trigger for updated_at
        `CREATE OR REPLACE FUNCTION update_updated_at_column()
         RETURNS TRIGGER AS $$
         BEGIN
             NEW.updated_at = CURRENT_TIMESTAMP;
             RETURN NEW;
         END;
         $$ language 'plpgsql'`,
        
        `DROP TRIGGER IF EXISTS update_users_updated_at ON users`,
        `CREATE TRIGGER update_users_updated_at
         BEFORE UPDATE ON users
         FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()`,
    }
    
    for i, migration := range migrations {
        if _, err := db.db.Exec(migration); err != nil {
            return fmt.Errorf("execute migration %d: %w", i+1, err)
        }
    }
    
    db.logger.Println("✅ PostgreSQL migrations completed")
    return nil
}

// PostgresUserRepository implements user data operations for PostgreSQL
type PostgresUserRepository struct {
    db     *sql.DB
    logger Logger
}

func NewPostgresUserRepository(db *PostgresDB) *PostgresUserRepository {
    return &PostgresUserRepository{
        db:     db.db,
        logger: db.logger,
    }
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (telegram_id, username, first_name, last_name, language, preferences)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `
    
    // Convert preferences to JSON
    prefsJSON := map[string]interface{}{
        "notifications": user.Preferences.Notifications,
        "theme":         user.Preferences.Theme,
        "timezone":      user.Preferences.Timezone,
    }
    
    err := r.db.QueryRowContext(ctx, query,
        user.TelegramID,
        user.Username,
        user.FirstName,
        user.LastName,
        user.Language,
        prefsJSON,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("insert user: %w", err)
    }
    
    r.logger.Printf("📝 User created: ID=%d, TelegramID=%d, Name=%s", 
        user.ID, user.TelegramID, user.FirstName)
    
    return nil
}

func (r *PostgresUserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
    query := `
        SELECT id, telegram_id, username, first_name, last_name, language, 
               is_blocked, preferences, created_at, updated_at
        FROM users 
        WHERE telegram_id = $1
    `
    
    var user domain.User
    var prefsJSON map[string]interface{}
    
    err := r.db.QueryRowContext(ctx, query, telegramID).Scan(
        &user.ID,
        &user.TelegramID,
        &user.Username,
        &user.FirstName,
        &user.LastName,
        &user.Language,
        &user.IsBlocked,
        &prefsJSON,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found: telegram_id=%d", telegramID)
        }
        return nil, fmt.Errorf("query user: %w", err)
    }
    
    // Parse preferences from JSON
    if prefsJSON != nil {
        if notifications, ok := prefsJSON["notifications"].(bool); ok {
            user.Preferences.Notifications = notifications
        }
        if theme, ok := prefsJSON["theme"].(string); ok {
            user.Preferences.Theme = theme
        }
        if timezone, ok := prefsJSON["timezone"].(string); ok {
            user.Preferences.Timezone = timezone
        }
    }
    
    return &user, nil
}

func (r *PostgresUserRepository) GetActiveUsers(ctx context.Context, limit int) ([]*domain.User, error) {
    query := `
        SELECT u.id, u.telegram_id, u.username, u.first_name, u.last_name, 
               u.language, u.is_blocked, u.preferences, u.created_at, u.updated_at
        FROM users u
        WHERE u.is_blocked = FALSE
        AND EXISTS (
            SELECT 1 FROM command_history ch 
            WHERE ch.user_id = u.id 
            AND ch.executed_at >= NOW() - INTERVAL '30 days'
        )
        ORDER BY u.updated_at DESC
        LIMIT $1
    `
    
    rows, err := r.db.QueryContext(ctx, query, limit)
    if err != nil {
        return nil, fmt.Errorf("query active users: %w", err)
    }
    defer rows.Close()
    
    var users []*domain.User
    for rows.Next() {
        var user domain.User
        var prefsJSON map[string]interface{}
        
        err := rows.Scan(
            &user.ID,
            &user.TelegramID,
            &user.Username,
            &user.FirstName,
            &user.LastName,
            &user.Language,
            &user.IsBlocked,
            &prefsJSON,
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("scan user row: %w", err)
        }
        
        // Parse preferences
        if prefsJSON != nil {
            // Simplified parsing - in production, use proper JSON library
            user.Preferences = domain.UserPrefs{
                Notifications: true,
                Theme:         "light",
                Timezone:      "UTC",
            }
        }
        
        users = append(users, &user)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("iterate user rows: %w", err)
    }
    
    return users, nil
}

func (r *PostgresUserRepository) GetUserStats(ctx context.Context) (*UserStats, error) {
    query := `
        SELECT 
            COUNT(*) as total_users,
            COUNT(CASE WHEN created_at >= NOW() - INTERVAL '7 days' THEN 1 END) as new_users_week,
            COUNT(CASE WHEN is_blocked = TRUE THEN 1 END) as blocked_users,
            COUNT(CASE WHEN updated_at >= NOW() - INTERVAL '1 day' THEN 1 END) as active_today
        FROM users
    `
    
    var stats UserStats
    err := r.db.QueryRowContext(ctx, query).Scan(
        &stats.TotalUsers,
        &stats.NewUsersThisWeek,
        &stats.BlockedUsers,
        &stats.ActiveToday,
    )
    
    if err != nil {
        return nil, fmt.Errorf("query user stats: %w", err)
    }
    
    return &stats, nil
}

func (r *PostgresUserRepository) HealthCheck(ctx context.Context) error {
    query := "SELECT 1"
    var result int
    
    err := r.db.QueryRowContext(ctx, query).Scan(&result)
    if err != nil {
        return fmt.Errorf("database health check failed: %w", err)
    }
    
    return nil
}

// Enhanced UserStats for PostgreSQL
type UserStats struct {
    TotalUsers       int `json:"total_users"`
    NewUsersThisWeek int `json:"new_users_week"`
    BlockedUsers     int `json:"blocked_users"`
    ActiveToday      int `json:"active_today"`
}
```

#### `internal/app/dependencies.go` - Dependency Injection
```go
package app

import (
    "log"
    "os"
    
    "yordamchi-dev-bot/internal/infrastructure/database"
    "yordamchi-dev-bot/internal/services/user"
    "yordamchi-dev-bot/pkg/config"
)

type Dependencies struct {
    Config      *config.Config
    Logger      *log.Logger
    Database    Database
    UserService user.Service
    Messages    map[string]interface{}
    RateLimiter RateLimiter
}

type Database interface {
    HealthCheck(ctx context.Context) error
    Close() error
}

type RateLimiter interface {
    Allow(userID int64) bool
    Reset(userID int64)
}

func NewDependencies() (*Dependencies, error) {
    // Initialize logger
    logger := log.New(os.Stdout, "[BOT] ", log.LstdFlags|log.Lshortfile)
    
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        return nil, fmt.Errorf("load config: %w", err)
    }
    
    // Initialize database based on configuration
    var db Database
    var userRepo user.Repository
    
    dbType := os.Getenv("DB_TYPE")
    if dbType == "" {
        dbType = cfg.Database.Type
    }
    
    switch dbType {
    case "postgres":
        connectionString := os.Getenv("DATABASE_URL")
        if connectionString == "" {
            connectionString = cfg.Database.ConnectionString
        }
        
        postgresDB, err := database.NewPostgresDB(connectionString, logger)
        if err != nil {
            return nil, fmt.Errorf("initialize postgres: %w", err)
        }
        
        db = postgresDB
        userRepo = database.NewPostgresUserRepository(postgresDB)
        
    default: // sqlite
        dbPath := os.Getenv("DATABASE_PATH")
        if dbPath == "" {
            dbPath = cfg.Database.ConnectionString
        }
        
        sqliteDB, err := database.NewSQLiteDB(dbPath, logger)
        if err != nil {
            return nil, fmt.Errorf("initialize sqlite: %w", err)
        }
        
        db = sqliteDB
        userRepo = database.NewUserRepository(sqliteDB)
    }
    
    // Initialize services
    userService := user.NewService(userRepo, logger)
    
    // Initialize rate limiter (simplified in-memory version)
    rateLimiter := NewInMemoryRateLimiter(cfg.Limits.RateLimitPerMinute)
    
    return &Dependencies{
        Config:      cfg,
        Logger:      logger,
        Database:    db,
        UserService: userService,
        Messages:    cfg.Messages,
        RateLimiter: rateLimiter,
    }, nil
}

func (d *Dependencies) Close() error {
    if d.Database != nil {
        return d.Database.Close()
    }
    return nil
}

// Simple in-memory rate limiter
type InMemoryRateLimiter struct {
    limits map[int64]int
    max    int
}

func NewInMemoryRateLimiter(maxRequests int) *InMemoryRateLimiter {
    return &InMemoryRateLimiter{
        limits: make(map[int64]int),
        max:    maxRequests,
    }
}

func (rl *InMemoryRateLimiter) Allow(userID int64) bool {
    current := rl.limits[userID]
    if current >= rl.max {
        return false
    }
    
    rl.limits[userID] = current + 1
    return true
}

func (rl *InMemoryRateLimiter) Reset(userID int64) {
    delete(rl.limits, userID)
}
```

This completes Days 4-6 of Week 1. Would you like me to continue with Day 7 (Testing & Deployment) and then move on to the other weeks' documentation?

