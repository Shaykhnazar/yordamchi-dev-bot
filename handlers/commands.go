package handlers

import (
    "fmt"
    "strings"
    "time"
)

type CommandHandler struct {
    Config *Config
}

func NewCommandHandler(config *Config) *CommandHandler {
    return &CommandHandler{Config: config}
}

func (h *CommandHandler) HandleCommand(command, username string) string {
    command = strings.ToLower(command)
    
    switch command {
    case "/start":
        return h.Config.Messages.Welcome + "\n\n/help - barcha buyruqlar ro'yxati"
    case "/help":
        return h.Config.Messages.Help
    case "/ping":
        return "🏓 Pong! Bot ishlayapti ✅"
    case "/hazil":
        return h.getRandomJoke()
    case "/iqtibos":
        return h.getRandomQuote()
    case "/haqida":
        return h.getAboutInfo()
    case "/vaqt":
        return h.getCurrentTime()
    case "/salom":
        return h.getGreeting(username)
    default:
        if strings.HasPrefix(command, "/") {
            return h.Config.Messages.UnknownCommand
        }
        return ""
    }
}

func (h *CommandHandler) getRandomJoke() string {
    if len(h.Config.Jokes) == 0 {
        return "😅 Hazillar yuklanmagan!"
    }
    // Now using the config parameter with the global function
    return GetRandomJoke(h.Config)
}

func (h *CommandHandler) getRandomQuote() string {
    if len(h.Config.Quotes) == 0 {
        return "💭 Iqtiboslar yuklanmagan!"
    }
    // Now using the config parameter with the global function
    return GetRandomQuote(h.Config)
}

func (h *CommandHandler) getAboutInfo() string {
    return fmt.Sprintf(`ℹ️ %s

🔸 Versiya: %s
🔸 Tavsif: %s
🔸 Yaratuvchi: %s
🔸 Til: Go (Golang)

Bu bot Go tilini o'rganish jarayonida yaratilmoqda! 🎯`,
        h.Config.Bot.Name,
        h.Config.Bot.Version,
        h.Config.Bot.Description,
        h.Config.Bot.Author)
}

func (h *CommandHandler) getCurrentTime() string {
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    return fmt.Sprintf("🕐 Hozirgi vaqt: %s", currentTime)
}

func (h *CommandHandler) getGreeting(username string) string {
    if username == "" {
        return "👋 Salom! Ismingizni bilmayman, lekin baribir xush kelibsiz! 😊"
    }
    return fmt.Sprintf("👋 Salom, %s! Go dasturlashni o'rganishga tayyormisiz? 🚀", username)
}