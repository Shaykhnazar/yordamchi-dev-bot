package main

import (
    "log"
    "os"
)

// main - Bot'ni ishga tushirish funksiyasi
//
// Bu funksiya BOT_TOKEN environment variable'idan bot token'ini olib,
// yangi bot ob'ektini yaratadi va bot'ni ishga tushiradi.
// Bot ishga tushirilmaganda xatolik kodini qaytaradi.
func main() {
    // Bot token'ni environment variable'dan olish
    token := os.Getenv("BOT_TOKEN")
    if token == "" {
        log.Fatal("BOT_TOKEN environment variable topilmadi")
    }

    // Bot'ni boshlash
    bot := NewBot(token)
    log.Println("🤖 Yordamchi Dev Bot ishga tushdi!")
    
    if err := bot.Start(); err != nil {
        log.Fatal("Bot'ni ishga tushirishda xatolik:", err)
    }
}