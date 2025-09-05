# Week 1: Clean Architecture Patterns (Toza Arxitektura Namunalari)

## 🎯 O'rganish Maqsadlari

Bu darsda siz professional darajadagi clean architecture pattern'larini va scalable bot yaratish usullarini o'rganasiz.

## 📚 Asosiy Tushunchalar

### 1. Clean Architecture nima?

Clean Architecture - bu Robert Martin (Uncle Bob) tomonidan taklif qilingan arxitektura pattern'i bo'lib, kod ixchamligi, testlash va kengaytirilishni osonlashtiradi.

**Asosiy qoidalar:**
- Dependency Inversion (Bog'liqlikni teskari aylantirish)
- Separation of Concerns (Mas'uliyatlarni ajratish)
- Interface-based design (Interface asosida dizayn)

### 2. Loyihaning Yangi Strukturasi

```
yordamchi-dev-bot/
├── cmd/bot/main.go                 # Application entry point
├── internal/
│   ├── domain/                     # Business entities va interfaces
│   │   ├── command.go
│   │   ├── user.go
│   │   └── context.go
│   ├── app/                        # Application layer
│   │   ├── bot.go
│   │   ├── router.go
│   │   ├── dependencies.go
│   │   └── user_service.go
│   ├── handlers/commands/          # Command handlers
│   │   ├── start.go
│   │   ├── help.go
│   │   ├── ping.go
│   │   └── github.go
│   ├── middleware/                 # Middleware layer
│   │   ├── logging.go
│   │   ├── auth.go
│   │   └── ratelimit.go
│   └── services/                   # External services
│       ├── http_client.go
│       ├── github_service.go
│       └── weather_service.go
```

### 3. Domain Layer (Domen Qatlami)

Domain layer - bu biznes mantiqining markaziy qismi:

```go
// Command - buyruq entity'si
type Command struct {
    ID        string
    Text      string
    User      *User
    Chat      *Chat
    Timestamp time.Time
}

// CommandHandler - buyruq ishlovchi interface
type CommandHandler interface {
    Handle(ctx context.Context, cmd *Command) (*Response, error)
    CanHandle(command string) bool
    Description() string
    Usage() string
}
```

**Muhim tushunchalar:**
- Entity'lar biznes obyektlarini ifodalaydi
- Interface'lar kontraktlarni belgilaydi
- Hech qanday tashqi bog'liqlik yo'q

### 4. Interface Pattern

Go'da interface'lar implicit tarzda implement qilinadi:

```go
// Interface ta'rifi
type CommandHandler interface {
    Handle(ctx context.Context, cmd *Command) (*Response, error)
    CanHandle(command string) bool
    Description() string
    Usage() string
}

// Implementation
type StartCommand struct {
    welcomeMessage string
    logger         Logger
}

// Interface metodlarini implement qilish
func (h *StartCommand) Handle(ctx context.Context, cmd *Command) (*Response, error) {
    // Implementation
}

func (h *StartCommand) CanHandle(command string) bool {
    return command == "/start"
}
```

### 5. Middleware Pattern

Middleware - bu so'rovni qayta ishlash zanjirini yaratish uchun:

```go
type Middleware interface {
    Process(ctx context.Context, next HandlerFunc) HandlerFunc
}

type LoggingMiddleware struct {
    logger Logger
}

func (m *LoggingMiddleware) Process(ctx context.Context, next HandlerFunc) HandlerFunc {
    return func(ctx context.Context, cmd *Command) (*Response, error) {
        start := time.Now()
        
        // So'rov boshlanishini log qilish
        m.logger.Info("Command processing started", "command", cmd.Text)
        
        // Keyingi handler'ga o'tish
        response, err := next(ctx, cmd)
        
        duration := time.Since(start)
        m.logger.Info("Command processing completed", "duration", duration)
        
        return response, err
    }
}
```

**Middleware'ning foydasi:**
- Cross-cutting concerns (logging, auth, rate limiting)
- Code reuse (qayta foydalanish)
- Separation of concerns

### 6. Dependency Injection

Dependency Injection - bu bog'liqliklarni tashqaridan berish:

```go
// Dependencies container
type Dependencies struct {
    Logger         Logger
    Config         *Config
    DB             *database.DB
    Router         Router
    GitHubService  *services.GitHubService
    WeatherService *services.WeatherService
}

// Constructor bilan bog'liqliklarni yaratish
func NewDependencies(config *Config, db *database.DB) (*Dependencies, error) {
    logger := NewStructuredLogger()
    
    githubService := services.NewGitHubService(logger)
    weatherService := services.NewWeatherService(logger)
    
    router := NewCommandRouter(logger)
    
    // Command handler'larni ro'yxatga olish
    router.RegisterHandler(commands.NewStartCommand(config.Messages.Welcome, logger))
    router.RegisterHandler(commands.NewHelpCommand(router, config.Messages.Help, logger))
    
    return &Dependencies{
        Logger:         logger,
        Config:         config,
        DB:             db,
        Router:         router,
        GitHubService:  githubService,
        WeatherService: weatherService,
    }, nil
}
```

### 7. Command Router Pattern

Router - bu buyruqlarni mos handler'ga yo'naltirish:

```go
type CommandRouter struct {
    handlers    []CommandHandler
    middlewares []Middleware
    logger      Logger
}

func (r *CommandRouter) Route(ctx context.Context, cmd *Command) (*Response, error) {
    // Mos handler topish
    var handler CommandHandler
    for _, h := range r.handlers {
        if h.CanHandle(cmd.Text) {
            handler = h
            break
        }
    }
    
    if handler == nil {
        return &Response{Text: "Noma'lum buyruq"}, nil
    }
    
    // Middleware zanjirini qurish
    handlerFunc := r.buildMiddlewareChain(handler.Handle)
    
    // Bajarish
    return handlerFunc(ctx, cmd)
}
```

### 8. Context Pattern

Go'da context - bu so'rov ma'lumotlarini uzatish uchun:

```go
// Context'ga ma'lumot qo'shish
func WithUser(ctx context.Context, user *User) context.Context {
    return context.WithValue(ctx, UserContextKey, user)
}

// Context'dan ma'lumot olish
func GetUserFromContext(ctx context.Context) (*User, bool) {
    user, ok := ctx.Value(UserContextKey).(*User)
    return user, ok
}

// Middleware'da context ishlatish
func (m *AuthMiddleware) Process(ctx context.Context, next HandlerFunc) HandlerFunc {
    return func(ctx context.Context, cmd *Command) (*Response, error) {
        // Foydalanuvchini ro'yxatga olish
        user := m.authenticateUser(cmd.User)
        
        // Context'ga qo'shish
        ctx = WithUser(ctx, user)
        
        return next(ctx, cmd)
    }
}
```

## 🔧 Handler Implementation Pattern

Har bir command uchun alohida handler:

```go
type StartCommand struct {
    welcomeMessage string
    logger         Logger
}

func NewStartCommand(welcomeMessage string, logger Logger) *StartCommand {
    return &StartCommand{
        welcomeMessage: welcomeMessage,
        logger:         logger,
    }
}

func (h *StartCommand) Handle(ctx context.Context, cmd *Command) (*Response, error) {
    user, _ := GetUserFromContext(ctx)
    
    message := h.welcomeMessage
    if user != nil && user.FirstName != "" {
        message += "\n\n👋 Salom, " + user.FirstName + "!"
    }
    
    h.logger.Info("Start command processed", "user_id", cmd.User.TelegramID)
    
    return &Response{
        Text:      message,
        ParseMode: "HTML",
    }, nil
}
```

## 💡 Clean Architecture'ning Afzalliklari

### 1. Testability (Testlash)
```go
// Mock interface yaratish oson
type MockLogger struct{}

func (m *MockLogger) Info(msg string, args ...interface{}) {}

// Test yozish
func TestStartCommand(t *testing.T) {
    logger := &MockLogger{}
    cmd := NewStartCommand("Welcome", logger)
    
    response, err := cmd.Handle(context.Background(), &Command{...})
    // assertions
}
```

### 2. Flexibility (Moslashuvchanlik)
- Interface'lar orqali implementation'ni osongina almashtirish mumkin
- Yangi feature'lar qo'shish oson
- Middleware orqali xususiyatlar qo'shish

### 3. Maintainability (Saqlanuvchanllik)
- Har bir komponent o'z mas'uliyatiga ega
- Bog'liqliklar aniq belgilangan
- Kod o'qishga oson

## 📈 Scalability Patterns

### 1. Async Processing
```go
// Asynchronous message processing
go b.processUpdate(&update)
```

### 2. Rate Limiting
```go
rateLimitMiddleware := NewRateLimitMiddleware(10, time.Minute, logger)
router.RegisterMiddleware(rateLimitMiddleware)
```

### 3. Background Tasks
```go
// Cleanup task
go func() {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        rateLimitMiddleware.Cleanup()
    }
}()
```

## 🎯 Keyingi Qadamlar

Clean Architecture pattern'i bilan:
- Yangi command'lar qo'shish oson
- Testing framework yozish mumkin
- Monitoring va observability qo'shish
- Horizontal scaling uchun tayyorlash

Bu arxitektura high-load va scalable application'lar yaratish uchun professional standart hisoblanadi!

## 📝 Vazifalar

1. Yangi command handler yarating
2. Custom middleware yozing
3. Interface pattern'ini tushunish
4. Dependency injection'ni qo'llash
5. Test yozish uchun mock'lar yarating

Clean Architecture orqali siz professional darajadagi Go application'lar yarata olasiz!