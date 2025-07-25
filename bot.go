package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
)

type Bot struct {
    Token string
    URL   string
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

func NewBot(token string) *Bot {
    return &Bot{
        Token: token,
        URL:   fmt.Sprintf("https://api.telegram.org/bot%s", token),
    }
}

func (b *Bot) Start() error {
    http.HandleFunc("/webhook", b.handleWebhook)
    log.Println("Server 8080 portda ishlamoqda...")
    return http.ListenAndServe(":8080", nil)
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

    chatID := msg.Chat.ID
    text := strings.ToLower(msg.Text)

    switch {
    case text == "/start":
        b.sendMessage(chatID, "🎉 Assalomu alaykum! Men Yordamchi Dev Bot. Go dasturlash tilini o'rganishingizda yordam beraman!\n\n/help - barcha buyruqlar ro'yxati")
    case text == "/help":
        helpText := `🤖 Yordamchi Dev Bot Buyruqlari:

/start - Botni ishga tushirish
/help - Bu yordam xabari
/ping - Bot ishlaganligini tekshirish
/hazil - Tasodifiy hazil
/iqtibos - Motivatsion iqtibos

Keyingi haftalarda ko'proq funksiyalar qo'shiladi! 🚀`
        b.sendMessage(chatID, helpText)
    case text == "/ping":
        b.sendMessage(chatID, "🏓 Pong! Bot ishlayapti ✅")
    case text == "/hazil":
        b.sendMessage(chatID, "😄 Dasturchi nima uchun ko'zoynak kiyadi? Chunki Java ko'ra olmaydi! ☕")
    case text == "/iqtibos":
        b.sendMessage(chatID, "💭 \"Birinchi kod ishlamasa, console.log qo'sh\" - Har bir dasturchi")
    default:
        if strings.HasPrefix(text, "/") {
            b.sendMessage(chatID, "❓ Noma'lum buyruq. /help yozing")
        }
    }
}

func (b *Bot) sendMessage(chatID int, text string) {
    url := fmt.Sprintf("%s/sendMessage", b.URL)
    
    payload := fmt.Sprintf(`{
        "chat_id": %d,
        "text": "%s",
        "parse_mode": "HTML"
    }`, chatID, text)

    resp, err := http.Post(url, "application/json", strings.NewReader(payload))
    if err != nil {
        log.Println("Xabar yuborishda xatolik:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        log.Printf("Telegram API xatolik: %d", resp.StatusCode)
    }
}