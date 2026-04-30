# TeleRiHa - Telegram Bot Framework

Build Telegram bots at the speed of thought. TeleRiHa is a production-grade Telegram Bot API framework for Go that provides a clean, modular, and developer-friendly way to build Telegram bots.

## Features

- **Pure net/http**: Direct communication with Telegram Bot API using Go's built-in `net/http` package
- **Zero external Telegram dependencies**: No wrappers, no abstractions - just pure Telegram API
- **Gin-like API**: Familiar router and middleware pattern for building bot handlers
- **Modular design**: Plugins, conversation management, and state storage are all pluggable
- **Production-ready**: Built with best practices for error handling, logging, and testing
- **Extensible Plugin System**: Rate limiting, analytics, i18n, and scheduling out of the box

## Quick Start

### Install the CLI

```bash
# Clone the repository
git clone https://github.com/Yashwanth-Kumar-26/teleriha
cd TeleRiHa

# Build the CLI
cd cmd/riha
go build -o riha .

# Or install globally
go install github.com/Yashwanth-Kumar-26/teleriha/cmd/riha@latest
```

### Create a Bot Token

1. Talk to [@BotFather](https://t.me/BotFather) on Telegram
2. Create a new bot with `/newbot`
3. Copy the bot token

### Run the Starter Bot

```bash
# Create a .env file
cd examples/starter
echo "BOT_TOKEN=your_bot_token_here" > .env

# Run the bot
go run main.go
```

## Core Concepts

### Bot

The core of TeleRiHa. Creates and manages the connection to Telegram's Bot API.

```go
import "github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"

b := bot.New("YOUR_BOT_TOKEN")
```

### Context

Provides access to the current update, chat, sender, and convenience methods for replying.

```go
func handler(ctx *bot.Context) error {
    // Get message text
    text := ctx.Text()
    
    // Reply to the message
    _, err := ctx.Reply("You said: " + text)
    return err
}
```

### Router

Routes incoming updates to the appropriate handlers based on commands, regex, or message types.

```go
b.Router.On("/start", handleStart)
b.Router.On("/help", handleHelp)
b.Router.OnRegex(`^hello$`, handleHello)
b.Router.Default(handleDefault)
```

### Plugins

TeleRiHa comes with several powerful plugins:

- **Rate Limit**: Control request frequency per user
- **Analytics**: Track bot usage and message statistics
- **I18n**: Multi-language support with automatic detection
- **Scheduler**: Run background tasks on a schedule

```go
import (
    "github.com/Yashwanth-Kumar-26/teleriha/plugins/ratelimit"
    "github.com/Yashwanth-Kumar-26/teleriha/plugins/analytics"
)

// Add plugins
b.PluginRegistry.RegisterPlugin(ratelimit.New(ratelimit.Config{
    MaxRequests: 5,
    Interval:    time.Minute,
}))
b.PluginRegistry.RegisterPlugin(analytics.New())

b.PluginRegistry.LoadAll(b)
```

### State Storage

Choose between in-memory storage or persistent Redis storage.

```go
import "github.com/Yashwanth-Kumar-26/teleriha/pkg/store"
import "github.com/redis/go-redis/v9"

// Use Redis store
s := store.NewRedisStore(&redis.Options{
    Addr: "localhost:6379",
})
```

## Project Structure

```
TeleRiHa/
├── cmd/
│   └── riha/                # CLI tool
├── pkg/
│   └── bot/                 # Core framework
│   └── store/               # State storage (Memory, Redis)
├── plugins/
│   ├── ratelimit/           # Rate limiting
│   ├── analytics/           # Message stats
│   ├── i18n/                # Multi-language
│   └── scheduler/           # Cron tasks
├── examples/
│   ├── starter/             # Simple starter bot
│   ├── ecommerce/           # E-Commerce bot example
│   └── ai-bot/              # AI-powered bot example
└── go.mod
```

## Core Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/rs/zerolog` - Structured logging
- `github.com/stretchr/testify` - Testing assertions
- `github.com/joho/godotenv` - Environment variable loading
- `github.com/redis/go-redis/v9` - Redis client

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

## License

TeleRiHa is released under the MIT License. See [LICENSE](LICENSE) for details.
