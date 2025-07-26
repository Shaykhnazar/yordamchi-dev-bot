# 🤖 Yordamchi Dev Bot

A Telegram bot built with Go (Golang) to assist developers with daily tasks. This project is part of the GoBot Challenge - a 30-day learning journey to master Go programming.

![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Status](https://img.shields.io/badge/status-Active-brightgreen.svg)

## 🎯 Features

- ✅ **Basic Commands**: Start, help, ping functionality
- ✅ **Entertainment**: Random jokes and motivational quotes
- ✅ **Utility**: Current time, bot information
- ✅ **JSON Configuration**: Easily customizable messages and content
- ✅ **Polling & Webhook**: Support for both development and production modes
- ✅ **Logging**: Comprehensive activity tracking
- ✅ **Error Handling**: Robust error management

## 📋 Available Commands

| Command      | Description                                |
| ------------ | ------------------------------------------ |
| `/start`   | Initialize the bot and get welcome message |
| `/help`    | Display all available commands             |
| `/ping`    | Check if bot is online                     |
| `/hazil`   | Get a random programming joke              |
| `/iqtibos` | Get a motivational programming quote       |
| `/haqida`  | Get information about the bot              |
| `/vaqt`    | Get current timestamp                      |

## 🛠 Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.21+**: [Download Go](https://golang.org/dl/)
- **Git**: [Download Git](https://git-scm.com/)
- **ngrok**: [Download ngrok](https://ngrok.com/download) (for webhook setup)
- **Postman**: [Download Postman](https://www.postman.com/downloads/) (for API testing)
- **Telegram Account**: For creating and testing the bot

## 🚀 Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/yordamchi-dev-bot.git
cd yordamchi-dev-bot
```

### 2. Initialize Go Module

```bash
go mod init yordamchi-dev-bot
go mod tidy
```

### 3. Create Your Telegram Bot

1. Open Telegram and find [@BotFather](https://t.me/botfather)
2. Send `/newbot` command
3. Choose a name: `Yordamchi Dev Bot`
4. Choose a username: `your_bot_username_bot`
5. Copy the bot token (format: `123456789:ABCdef-YourTokenHere`)

### 4. Environment Configuration

#### Create .env file from example:

```bash
# Copy the example file
cp .env.example .env
```

#### Edit .env file:

```bash
# .env
BOT_TOKEN=your_bot_token_from_botfather
BOT_MODE=polling
PORT=8080
DB_TYPE=sqlite
DEBUG=true
```

**Important**: Replace `your_bot_token_from_botfather` with your actual bot token!

### 5. Verify Configuration Files

Ensure these files exist in your project root:

- `config.json` - Bot configuration and messages
- `.env` - Environment variables
- `main.go` - Application entry point
- `bot.go` - Bot logic and handlers
- `config.go` - Configuration management

## 🔧 Running the Bot

### Using Webhook

For webhook setup, you'll need to expose your local server to the internet.

#### Step 1: Install and Setup ngrok

```bash
# Download ngrok from https://ngrok.com/download
# Or install via package manager:

# macOS
brew install ngrok

# Windows (Chocolatey)
choco install ngrok

# Linux
sudo snap install ngrok
```

#### Step 2: Expose Local Server

```bash
# Terminal 1: Start your bot in webhook mode
export BOT_TOKEN="your_actual_bot_token"
export BOT_MODE="webhook"
go run .

# Terminal 2: Expose port 8080
ngrok http 8080
```

**ngrok output example:**

```
Forwarding    https://abc123.ngrok.io -> http://localhost:8080
```

#### Step 3: Set Webhook via Postman

##### Method 1: Using Postman

1. Open Postman
2. Create a new POST request
3. **URL**: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook`
4. **Headers**:
   - `Content-Type`: `application/json`
5. **Body** (raw JSON):

```json
{
    "url": "https://abc123.ngrok.io/webhook"
}
```

6. Click **Send**

**Expected response:**

```json
{
    "ok": true,
    "result": true,
    "description": "Webhook was set"
}
```

##### Method 2: Using cURL

```bash
curl -X POST "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook" \
     -H "Content-Type: application/json" \
     -d '{"url": "https://abc123.ngrok.io/webhook"}'
```

#### Step 4: Verify Webhook

**Using Postman:**

- **GET** request to: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getWebhookInfo`

**Using cURL:**

```bash
curl "https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getWebhookInfo"
```

**Expected response:**

```json
{
    "ok": true,
    "result": {
        "url": "https://abc123.ngrok.io/webhook",
        "has_custom_certificate": false,
        "pending_update_count": 0
    }
}
```

## 🧪 Testing Your Bot

### 1. Find Your Bot in Telegram

- Search for your bot username in Telegram
- Start a conversation

### 2. Test Basic Commands

Send these messages to test functionality:

```
/start
Expected: Welcome message with instructions

/help  
Expected: List of all available commands

/ping
Expected: "🏓 Pong! Bot ishlayapti ✅"

/hazil
Expected: Random programming joke

/iqtibos  
Expected: Motivational programming quote

/haqida
Expected: Information about the bot

/vaqt
Expected: Current timestamp
```

### 3. Monitor Logs

In your terminal running the bot, you should see:

```
👤 YourName (@yourusername): /start
👤 YourName (@yourusername): /help
👤 YourName (@yourusername): /ping
```

## 📁 Project Structure

```
yordamchi-dev-bot/
├── main.go              # Application entry point
├── bot.go               # Bot logic and message handling
├── config.go            # Configuration management
├── config.json          # Bot settings and messages
├── .env.example         # Environment variables template
├── .env                 # Your environment variables (create this)
├── go.mod               # Go module definition
├── go.sum               # Go module dependencies
├── README.md            # This file
└── .gitignore           # Git ignore rules
```

## 🔍 Troubleshooting

### Common Issues and Solutions

#### ❌ "BOT_TOKEN environment variable topilmadi"

**Solution:**

```bash
# Check if .env file exists and contains BOT_TOKEN
cat .env

# Or set manually:
export BOT_TOKEN="your_actual_bot_token"
```

#### ❌ "config.json faylini ochishda xatolik"

**Solution:**

```bash
# Verify config.json exists and is valid
ls -la config.json
cat config.json | jq '.'  # Validates JSON format
```

#### ❌ Bot doesn't respond to messages

**Possible causes:**

1. **Wrong bot token**: Verify token from @BotFather
2. **Bot not running**: Check terminal for "Bot polling started" message
3. **Webhook conflicts**: Delete webhook if using polling mode

**Solutions:**

```bash
# Test bot token validity
curl "https://api.telegram.org/bot<YOUR_TOKEN>/getMe"

# Delete webhook (if using polling)
curl -X POST "https://api.telegram.org/bot<YOUR_TOKEN>/deleteWebhook"
```

#### ❌ "Telegram API xatolik: 409"

**Cause:** Another bot instance is running or webhook is set while using polling.

**Solution:**

```bash
# Delete webhook
curl -X POST "https://api.telegram.org/bot<YOUR_TOKEN>/deleteWebhook"

# Kill other bot processes
pkill -f "go run"
```

#### ❌ ngrok tunnel not working

**Solutions:**

```bash
# Check if ngrok is running
ps aux | grep ngrok

# Restart ngrok
ngrok http 8080

# Verify webhook URL matches ngrok URL
curl "https://api.telegram.org/bot<YOUR_TOKEN>/getWebhookInfo"
```

### Debug Mode

Enable debug logging:

```bash
export DEBUG=true
go run .
```

## 🔧 Bot Management Commands

### Using cURL

```bash
# Get bot information
curl "https://api.telegram.org/bot<TOKEN>/getMe"

# Get webhook info
curl "https://api.telegram.org/bot<TOKEN>/getWebhookInfo"

# Set webhook
curl -X POST "https://api.telegram.org/bot<TOKEN>/setWebhook" \
     -H "Content-Type: application/json" \
     -d '{"url": "https://your-ngrok-url.ngrok.io/webhook"}'

# Delete webhook
curl -X POST "https://api.telegram.org/bot<TOKEN>/deleteWebhook"
```

### Using Postman Collection

Import this collection to Postman for easy bot management:

**Collection variables:**

- `bot_token`: Your bot token
- `webhook_url`: Your ngrok URL + /webhook

**Requests:**

1. **GET Bot Info**: `{{base_url}}/bot{{bot_token}}/getMe`
2. **GET Webhook Info**: `{{base_url}}/bot{{bot_token}}/getWebhookInfo`
3. **SET Webhook**: `{{base_url}}/bot{{bot_token}}/setWebhook`
4. **DELETE Webhook**: `{{base_url}}/bot{{bot_token}}/deleteWebhook`

Where `base_url` = `https://api.telegram.org`

## 📈 Next Steps (Week 2)

This bot is designed for progressive enhancement:

- 🔄 **GitHub API Integration**: Repository statistics and issue search
- 🔍 **Stack Overflow Search**: Programming Q&A lookup
- 🎨 **Code Formatter**: Go code formatting service
- 📚 **Documentation Lookup**: Quick access to Go docs
- 🔄 **Concurrent Processing**: Goroutines and channels
- 💾 **Database Integration**: PostgreSQL user management

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🎯 Learning Goals

By completing this project, you will master:

- ✅ **Go Fundamentals**: Packages, structs, functions, error handling
- ✅ **HTTP & JSON**: REST API integration and data parsing
- ✅ **Concurrency**: Goroutines for polling and webhook handling
- ✅ **Configuration Management**: JSON and environment variables
- ✅ **Testing**: Unit tests and integration testing
- ✅ **Deployment**: Docker containerization and production setup

## 🙋‍♂️ Support

If you encounter any issues:

1. Check the [Troubleshooting](#-troubleshooting) section
2. Review your `.env` and `config.json` files
3. Verify your bot token with @BotFather
4. Check the [Issues](https://github.com/yourusername/yordamchi-dev-bot/issues) page
5. Create a new issue with detailed error logs

---

**Happy Coding! 🚀**

Built with ❤️ using Go • Part of the 30-Day GoBot Challenge
