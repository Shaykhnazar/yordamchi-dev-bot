package app

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"yordamchi-dev-bot/database"
	"yordamchi-dev-bot/handlers"
	"yordamchi-dev-bot/internal/domain"
	"yordamchi-dev-bot/internal/handlers/commands"
	"yordamchi-dev-bot/internal/middleware"
	"yordamchi-dev-bot/internal/services"
)

// Dependencies holds all application dependencies
type Dependencies struct {
	// Core
	Logger domain.Logger
	Config *handlers.Config
	DB     *database.DB
	Router domain.Router

	// Services
	GitHubService  *services.GitHubService
	WeatherService *services.WeatherService
	UserService    domain.UserService

	// Bot
	StartTime time.Time
}

// NewDependencies creates and configures all application dependencies
func NewDependencies(config *handlers.Config, db *database.DB) (*Dependencies, error) {
	startTime := time.Now()
	
	// Create logger
	logger := NewStructuredLogger()

	// Create logger adapter for services
	serviceLogger := &loggerAdapter{logger: logger}
	
	// Create services
	githubService := services.NewGitHubService(serviceLogger)
	weatherService := services.NewWeatherService(serviceLogger)
	userService := NewUserService(db, logger)

	// Create router
	router := NewCommandRouter(logger)

	// Create and register middlewares
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	authMiddleware := middleware.NewAuthMiddleware(userService, logger)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(10, time.Minute, logger) // 10 requests per minute

	router.RegisterMiddleware(loggingMiddleware)
	router.RegisterMiddleware(authMiddleware)
	router.RegisterMiddleware(rateLimitMiddleware)

	// Create and register command handlers
	startCommand := commands.NewStartCommand(config.Messages.Welcome, logger)
	helpCommand := commands.NewHelpCommand(router, config.Messages.Help, logger)
	pingCommand := commands.NewPingCommand(logger, startTime)
	githubCommand := commands.NewGitHubCommand(githubService, logger)
	hazilCommand := commands.NewHazilCommand(config.Jokes, logger)
	iqtibosCommand := commands.NewIqtibosCommand(config.Quotes, logger)
	haqidaCommand := commands.NewHaqidaCommand(config, logger)
	vaqtCommand := commands.NewVaqtCommand(logger)
	salomCommand := commands.NewSalomCommand(logger)
	statsCommand := commands.NewStatsCommand(userService, startTime, logger)
	weatherCommand := commands.NewWeatherCommand(weatherService, logger)

	router.RegisterHandler(startCommand)
	router.RegisterHandler(helpCommand)
	router.RegisterHandler(pingCommand)
	router.RegisterHandler(githubCommand)
	router.RegisterHandler(hazilCommand)
	router.RegisterHandler(iqtibosCommand)
	router.RegisterHandler(haqidaCommand)
	router.RegisterHandler(vaqtCommand)
	router.RegisterHandler(salomCommand)
	router.RegisterHandler(statsCommand)
	router.RegisterHandler(weatherCommand)

	// Start background tasks
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			rateLimitMiddleware.Cleanup()
		}
	}()

	return &Dependencies{
		Logger:         logger,
		Config:         config,
		DB:             db,
		Router:         router,
		GitHubService:  githubService,
		WeatherService: weatherService,
		UserService:    userService,
		StartTime:      startTime,
	}, nil
}

// StructuredLogger implements domain.Logger interface
type StructuredLogger struct {
	logger *log.Logger
}

// NewStructuredLogger creates a new structured logger
func NewStructuredLogger() *StructuredLogger {
	return &StructuredLogger{
		logger: log.New(os.Stdout, "[BOT] ", log.LstdFlags|log.Lshortfile),
	}
}

// Debug logs debug messages
func (l *StructuredLogger) Debug(msg string, args ...interface{}) {
	l.logWithFields("DEBUG", msg, args...)
}

// Info logs info messages
func (l *StructuredLogger) Info(msg string, args ...interface{}) {
	l.logWithFields("INFO", msg, args...)
}

// Warn logs warning messages
func (l *StructuredLogger) Warn(msg string, args ...interface{}) {
	l.logWithFields("WARN", msg, args...)
}

// Error logs error messages
func (l *StructuredLogger) Error(msg string, args ...interface{}) {
	l.logWithFields("ERROR", msg, args...)
}

// logWithFields formats structured logging with key-value pairs
func (l *StructuredLogger) logWithFields(level, msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Printf("%s: %s", level, msg)
		return
	}
	
	// Format key-value pairs
	var fields []string
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields = append(fields, fmt.Sprintf("%v=%v", args[i], args[i+1]))
		} else {
			fields = append(fields, fmt.Sprintf("extra=%v", args[i]))
		}
	}
	
	l.logger.Printf("%s: %s %s", level, msg, strings.Join(fields, " "))
}

// With creates a new logger with additional context (simplified implementation)
func (l *StructuredLogger) With(args ...interface{}) domain.Logger {
	return l // For now, return same logger. In production, would add context
}

// loggerAdapter adapts domain.Logger to services.Logger interface
type loggerAdapter struct {
	logger domain.Logger
}

// Printf implements services.Logger interface
func (a *loggerAdapter) Printf(format string, args ...interface{}) {
	a.logger.Info(format, args...)
}

// Println implements services.Logger interface  
func (a *loggerAdapter) Println(args ...interface{}) {
	a.logger.Info("%v", args...)
}