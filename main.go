package main

import (
    "log"
    "os"
    "yordamchi-dev-bot/handlers"
    "github.com/joho/godotenv"
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

    // Bot'ni yaratish va konfiguratsiyani uzatish
    bot := NewBot(token, config)
    log.Printf("🤖 %s (v%s) ishga tushdi!", config.Bot.Name, config.Bot.Version)
    
    if err := bot.Start(); err != nil {
        log.Fatal("Bot'ni ishga tushirishda xatolik:", err)
    }
}