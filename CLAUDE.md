# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Telegram bot built with Go called "Yordamchi Dev Bot" - a developer assistant bot that provides jokes, quotes, and utility commands. The project is part of a 30-day Go learning challenge.

## Development Commands

### Running the Bot
```bash
# Run in development mode
go run .

# Build the binary
go build -o yordamchi-dev-bot .

# Install dependencies
go mod tidy
```

### Environment Setup
The bot requires a `.env` file with:
```
BOT_TOKEN=your_bot_token_from_botfather
BOT_MODE=webhook
PORT=8080
```

The bot can run in two modes:
- **Polling mode**: Set `BOT_MODE=polling` (for development)
- **Webhook mode**: Set `BOT_MODE=webhook` (for production with ngrok)

## Architecture

### Core Structure
- `main.go` - Application entry point, loads config and starts bot
- `bot.go` - Bot implementation with webhook handling and message processing
- `handlers/` - Command handlers and configuration management
  - `config.go` - Configuration loading and random content functions
  - `commands.go` - Command handler implementation (currently duplicates bot.go logic)
- `config.json` - Bot configuration including messages, jokes, and quotes

### Key Components
- **Bot struct**: Main bot implementation with token, URL, config, and handler
- **Update/Message structs**: Telegram API message types
- **Config struct**: Configuration data from config.json
- **CommandHandler**: Modular command processing (handlers/commands.go)

### Dependencies
- `github.com/joho/godotenv` - Environment variable loading from .env files

## Bot Commands
- `/start` - Welcome message
- `/help` - List all commands  
- `/ping` - Health check
- `/hazil` - Random programming joke
- `/iqtibos` - Random motivational quote
- `/haqida` - Bot information
- `/vaqt` - Current timestamp
- `/salom` - Personalized greeting

## Configuration
Bot behavior is configured via `config.json` which contains:
- Bot metadata (name, version, description, author)
- Message templates (welcome, help, unknown command)
- Arrays of jokes and quotes for random selection

## Development Notes
- The bot currently has duplicate command handling logic in both `bot.go` and `handlers/commands.go`
- No test files exist yet - this is identified as a future enhancement
- The project uses webhook mode for production deployment with ngrok tunnel setup
- All user interactions are logged with username and message content