package main

import (
    "log"
    "os"
    "yordamchi-dev-bot/handlers"
    "github.com/joho/godotenv"
    "yordamchi-dev-bot/database"
)

func main() {
    // .env faylini yuklash
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️ .env fayli topilmadi, environment variable'lardan o'qiladi")
    } else {
        log.Println("✅ .env fayli muvaffaqiyatli yuklandi")
    }

    // Konfiguratsiyani yuklash
    config, err := handlers.LoadConfig()
    if err != nil {
        log.Fatal("Konfiguratsiya yuklashda xatolik:", err)
    }
    log.Println("✅ Konfiguratsiya muvaffaqiyatli yuklandi")

    // Bot token'ni environment variable'dan olish
    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        log.Fatal("BOT_TOKEN environment variable topilmadi. .env faylini tekshiring!")
    }

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

    bot := NewBotWithDB(token, config, db)
    log.Printf("🤖 %s (v%s) %s bilan ishga tushdi!", 
        config.Bot.Name, config.Bot.Version, dbType)
    
    appPort := os.Getenv("APP_PORT")
    if appPort == "" {
        log.Fatal("APP_PORT environment variable topilmadi. .env faylini tekshiring!")
    }

    if err := bot.Start(appPort); err != nil {
        log.Fatal("Bot'ni ishga tushirishda xatolik:", err)
    }
}