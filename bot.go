package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"yordamchi-dev-bot/database"
	"yordamchi-dev-bot/handlers"
)

type Bot struct {
	Token   string
	URL     string
	Config  *handlers.Config
	Handler *handlers.CommandHandler
	DB      *database.DB
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

func NewBot(token string, config *handlers.Config) *Bot {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal("Ma'lumotlar bazasi xatoligi:", err)
	}

	return &Bot{
		Token:   token,
		URL:     fmt.Sprintf("https://api.telegram.org/bot%s", token),
		Config:  config,
		Handler: handlers.NewCommandHandler(config),
		DB:      db,
	}
}

func NewBotWithDB(token string, config *handlers.Config, db *database.DB) *Bot {
	return &Bot{
		Token:   token,
		URL:     fmt.Sprintf("https://api.telegram.org/bot%s", token),
		Config:  config,
		Handler: handlers.NewCommandHandler(config),
		DB:      db,
	}
}

func (b *Bot) Start(appPort string) error {
	http.HandleFunc("/webhook", b.handleWebhook)
	log.Printf("Server %v portda ishlamoqda...", appPort)
	return http.ListenAndServe(":"+appPort, nil)
}

func (b *Bot) handleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Body o'qishda xatolik:", err)
		return
	}

	var update Update
	if err := json.Unmarshal(body, &update); err != nil {
		log.Println("JSON parse qilishda xatolik:", err)
		return
	}

	// Xabarlarni qayta ishlash
	b.processMessage(update.Message)
}

func (b *Bot) processMessage(msg Message) {
	if msg.Text == "" {
		return
	}

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

	// Foydalanuvchi faolligini logging qilish
	log.Printf("👤 %s (@%s): %s", msg.From.FirstName, msg.From.Username, msg.Text)

	chatID := msg.Chat.ID
	text := strings.ToLower(msg.Text)

	switch {
	case text == "/start":
		welcomeMsg := b.Config.Messages.Welcome + "\n\n/help - barcha buyruqlar ro'yxati"
		b.sendMessage(chatID, welcomeMsg)
	case text == "/help":
		b.sendMessage(chatID, b.Config.Messages.Help)
	case text == "/ping":
		b.sendMessage(chatID, "🏓 Pong! Bot ishlayapti ✅")
	case text == "/hazil":
		joke := handlers.GetRandomJoke(b.Config)
		b.sendMessage(chatID, joke)
	case text == "/iqtibos":
		quote := handlers.GetRandomQuote(b.Config)
		b.sendMessage(chatID, quote)
	case text == "/haqida":
		aboutText := fmt.Sprintf(`ℹ️ %s

                                🔸 Versiya: %s
                                🔸 Tavsif: %s
                                🔸 Yaratuvchi: %s
                                🔸 Til: Go (Golang)

                                Bu bot Go tilini o'rganish jarayonida yaratilmoqda! 🎯`,
			b.Config.Bot.Name,
			b.Config.Bot.Version,
			b.Config.Bot.Description,
			b.Config.Bot.Author)
		b.sendMessage(chatID, aboutText)
	case text == "/vaqt":
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		b.sendMessage(chatID, fmt.Sprintf("🕐 Hozirgi vaqt: %s", currentTime))
	case text == "/salom":
		greeting := fmt.Sprintf("👋 Salom, %s! Go dasturlashni o'rganishga tayyormisiz? 🚀", msg.From.FirstName)
		b.sendMessage(chatID, greeting)
	case text == "/stats":
		count, err := b.DB.GetUserStats()
		if err != nil {
			b.sendMessage(chatID, "❌ Statistika olishda xatolik")
		} else {
			b.sendMessage(chatID, fmt.Sprintf("📊 Jami foydalanuvchilar: %d", count))
		}

	default:
		if strings.HasPrefix(text, "/") {
			b.sendMessage(chatID, b.Config.Messages.UnknownCommand)
		}
	}
}

func (b *Bot) sendMessage(chatID int, text string) error {
	url := fmt.Sprintf("%s/sendMessage", b.URL)

	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("JSON yaratishda xatolik: %w", err)
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return fmt.Errorf("HTTP so'rov yuborishda xatolik: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Telegram API xatolik: %d, javob: %s", resp.StatusCode, string(body))
	}

	return nil
}
