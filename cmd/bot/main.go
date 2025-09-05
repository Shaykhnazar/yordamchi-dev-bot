package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	
	"yordamchi-dev-bot/database"
	"yordamchi-dev-bot/handlers"
	"yordamchi-dev-bot/internal/app"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, reading from environment variables")
	} else {
		log.Println("✅ .env file loaded successfully")
	}

	// Load configuration
	config, err := handlers.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("✅ Configuration loaded successfully")

	// Get bot token from environment
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN environment variable not found. Check your .env file!")
	}

	// Initialize database
	var db *database.DB
	dbType := os.Getenv("DB_TYPE")
	
	switch dbType {
	case "postgres":
		db, err = database.NewPostgresDB()
	default:
		db, err = database.NewDB() // SQLite
	}
	
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	log.Printf("✅ Database connected: %s", dbType)

	// Initialize application dependencies
	dependencies, err := app.NewDependencies(config, db)
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	// Create bot instance
	bot := app.NewTelegramBot(token, dependencies)

	// Get port from environment
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8090" // Default port
	}

	// Start bot server
	log.Printf("🤖 %s (v%s) starting with clean architecture on port %s", 
		config.Bot.Name, config.Bot.Version, appPort)
	
	if err := bot.Start(appPort); err != nil {
		log.Fatalf("Failed to start bot server: %v", err)
	}
}