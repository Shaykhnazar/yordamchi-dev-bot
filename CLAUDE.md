# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Telegram bot built with Go called "Yordamchi Dev Bot" - a developer assistant bot that provides jokes, quotes, and utility commands. The project is part of a 30-day Go learning challenge with comprehensive documentation covering 4 weeks of progressive development.

## Documentation

Comprehensive learning documentation is available in `.claude/docs/`:

- **[README.md](.claude/docs/README.md)** - Master guide and learning path overview
- **[Week 1](.claude/docs/week1_complete.md)** - ✅ Foundation (implemented in current codebase)
- **[Week 2](.claude/docs/week2_documentation.md)** - 📖 External APIs (ready for implementation)
- **[Week 3](.claude/docs/week3.md)** - ✅ Database & Architecture (implemented in current codebase)
- **[Week 4](.claude/docs/week4.md)** - 📖 Advanced Features (ready for implementation)
- **[Architecture](.claude/docs/architecture_docs.md)** - System design patterns

The documentation includes detailed code explanations, step-by-step breakdowns, programming concepts, and implementation guidance.

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
APP_PORT=8090
DB_TYPE=sqlite
```

**Bot Modes:**
- **Polling mode**: Set `BOT_MODE=polling` (for development)
- **Webhook mode**: Set `BOT_MODE=webhook` (for production)

**Database Options:**
- **SQLite**: Set `DB_TYPE=sqlite` (default, for development)
- **PostgreSQL**: Set `DB_TYPE=postgres` (requires `DATABASE_URL` environment variable)

## Architecture

### Core Structure

- `main.go` - Application entry point, loads config, initializes database, and starts bot
- `bot.go` - Bot implementation with webhook handling, message processing, and database integration
- `handlers/` - Command handlers and configuration management
  - `config.go` - Configuration loading and random content functions
  - `commands.go` - Command handler implementation (currently duplicates bot.go logic)
- `database/` - Database layer with multi-provider support
  - `db.go` - SQLite implementation with user management and activity tracking
  - `postgres.go` - PostgreSQL implementation for production
- `config.json` - Bot configuration including messages, jokes, and quotes
- `.claude/docs/` - Comprehensive learning documentation (4-week curriculum)

### Key Components

- **Bot struct**: Main bot implementation with token, URL, config, database, and handler
- **Update/Message structs**: Telegram API message types
- **Config struct**: Configuration data from config.json
- **DB struct**: Database abstraction layer supporting SQLite and PostgreSQL
- **User/UserActivity structs**: Data models for user management and activity tracking
- **CommandHandler**: Modular command processing (handlers/commands.go)

### Dependencies

- `github.com/joho/godotenv` - Environment variable loading from .env files
- `github.com/mattn/go-sqlite3` - SQLite database driver (for SQLite mode)
- `github.com/lib/pq` - PostgreSQL database driver (for PostgreSQL mode)

## Bot Commands

### Core Commands (Implemented)
- `/start` - Welcome message and bot introduction
- `/help` - List all available commands
- `/ping` - Health check and connectivity test
- `/hazil` - Random programming joke from config
- `/iqtibos` - Random motivational quote from config
- `/haqida` - Bot information and version details
- `/vaqt` - Current timestamp display
- `/salom` - Personalized greeting with user's name
- `/stats` - User statistics and bot usage metrics

### Database Features (Implemented)
- **Automatic User Registration**: Users are automatically registered on first interaction
- **Activity Tracking**: All command usage is logged for analytics
- **Statistics**: Track total users and command usage patterns

## Configuration

Bot behavior is configured via `config.json` which contains:

- Bot metadata (name, version, description, author)
- Message templates (welcome, help, unknown command)
- Arrays of jokes and quotes for random selection

## Development Status

### ✅ Implemented Features
- **Core Bot Functionality**: Webhook handling, command routing, configuration management
- **Database Integration**: SQLite/PostgreSQL support with user management and activity tracking
- **Production Ready**: Environment-based configuration, proper error handling, logging
- **Multi-Database Support**: Automatic database selection based on environment variables

### 📖 Ready for Implementation (Documented)
- **External API Integration**: GitHub, Stack Overflow, weather services (Week 2 & 4 docs)
- **Advanced Architecture**: Interface-based command system with middleware (Week 4 docs)
- **Performance Features**: Caching, rate limiting, optimization (Week 4 docs)
- **Testing Framework**: Unit tests, integration tests, mocking (Week 4 docs)
- **Security Features**: Input validation, sanitization, access control (Week 4 docs)

### 📝 Technical Notes
- The bot currently has duplicate command handling logic in both `bot.go` and `handlers/commands.go`
- Database operations are non-blocking and include proper error handling
- All user interactions are logged with username and message content
- The project supports both development (SQLite) and production (PostgreSQL) databases
- Comprehensive documentation provides learning path and implementation guidance

## 🔄 Development Rules

### Progressive Implementation
When implementing features from the documentation:

1. **Step-by-Step Implementation**: Follow documented weeks progressively (Week 1 → Week 2 → Week 3 → Week 4)
2. **Commit After Each Part**: After completing each significant feature or day's implementation, commit changes before continuing
3. **Sync with Documentation**: Ensure implemented code matches the documented patterns and explanations

### Learning Documentation Generation
For each implementation step:

1. **Create Learning Files**: Generate Uzbek language code explanations and definitions in `learn/` directory
2. **File Structure**: Create separate files for each major concept or implementation step
3. **Content Focus**: Include Go concepts, programming patterns, and code explanations in Uzbek
4. **Progressive Learning**: Build upon previous concepts with each new file

### Example Learning File Structure:
```
learn/
├── week1_go_asoslari.md          # Week 1: Go basics and fundamentals
├── week1_struct_va_interface.md   # Week 1: Structs and interfaces  
├── week2_http_client.md          # Week 2: HTTP client concepts
├── week2_json_parsing.md         # Week 2: JSON parsing and handling
├── week3_database_patterns.md    # Week 3: Database design patterns
├── week4_middleware_patterns.md  # Week 4: Middleware and architecture
└── ...
```

### Implementation Order
1. ✅ **Week 1 & 3**: Already implemented (foundation + database)
2. 📖 **Week 2**: External API integrations (GitHub, Stack Overflow, weather)
3. 📖 **Week 4**: Advanced patterns (interfaces, middleware, testing, security)
