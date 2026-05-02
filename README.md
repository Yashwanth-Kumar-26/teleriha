# TeleRiHa - Telegram Bot Framework for Go

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License">
  <img src="https://img.shields.io/badge/Telegram%20API-v7.8-blue?style=flat-square" alt="Telegram API">
</p>

Build production-grade Telegram bots with a clean, modular, and developer-friendly framework. TeleRiHa provides a Gin-like API with zero external Telegram dependencies.

## Features

| Feature | Description |
|---------|-------------|
| **Pure `net/http`** | Direct communication with Telegram Bot API using Go's built-in HTTP package |
| **Zero Dependencies** | No wrappers or abstractions - just pure Telegram API calls |
| **Gin-like API** | Familiar router and middleware patterns |
| **Plugin System** | Extensible plugins: rate limiting, analytics, i18n, scheduling |
| **State Management** | In-memory or Redis-backed storage with pluggable interface |
| **Conversation Flows** | Multi-step conversation builder with state management |
| **Inline Keyboards** | Fluent builder API for inline/reply keyboards |

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/Yashwanth-Kumar-26/teleriha
cd TeleRiHa

# Build the CLI
cd cmd/riha && go build -o riha .

# Or install globally
go install github.com/Yashwanth-Kumar-26/teleriha/cmd/riha@latest
```

### Create Your Bot

1. Talk to [@BotFather](https://t.me/BotFather) on Telegram
2. Create a new bot with `/newbot`
3. Copy the bot token

### Run the Starter Bot

```bash
cd examples/starter
echo "BOT_TOKEN=your_bot_token_here" > .env
go run main.go
```

### Basic Example

```go
package main

import (
    "os"
    "time"

    "github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
    "github.com/joho/godotenv"
    "github.com/rs/zerolog/log"
)

func main() {
    godotenv.Load()

    b := bot.New(os.Getenv("BOT_TOKEN"))

    // Add middleware
    b.Router.Use(bot.Logger(log.Logger))
    b.Router.Use(bot.Recover(log.Logger))

    // Register handlers
    b.Router.On("/start", handleStart)
    b.Router.On("/help", handleHelp)
    b.Router.Default(handleDefault)

    // Start polling
    b.StartPolling(1*time.Second, 30*time.Second)
}

func handleStart(ctx *bot.Context) error {
    return ctx.Reply("Welcome! Type /help for commands.")
}

func handleHelp(ctx *bot.Context) error {
    return ctx.Reply("Available commands: /start, /help")
}

func handleDefault(ctx *bot.Context) error {
    return ctx.Reply("Unknown command. Try /help")
}
```

## Architecture

```
+-------------------------------------------------------------+
|                         TeleRiHa Bot                        |
+-------------------------------------------------------------+
|  +-------------+  +-------------+  +---------------------+   |
|  |   Router    |  | Middleware  |  |  Plugin Registry    |   |
|  |  On/OnRegex |  |  Logger/    |  |  RateLimit/Analytics|   |
|  |  OnCallback |  |  Recover/   |  |  I18n/Scheduler     |   |
|  |  OnInline   |  |  RateLimit  |  |                     |   |
|  +-------------+  +-------------+  +---------------------+   |
|         |                |                    |             |
|         +----------------+--------------------+             |
|                          |                                   |
|                          v                                   |
|  +-----------------------------------------------------+    |
|  |                     Context                          |    |
|  |  Reply/Send/Edit/Delete/AnswerCallback              |    |
|  +-----------------------------------------------------+    |
|                          |                                   |
|                          v                                   |
|  +-----------------------------------------------------+    |
|  |              Telegram Bot API (net/http)            |    |
|  +-----------------------------------------------------+    |
+-------------------------------------------------------------+
|  +-----------------------------------------------------+    |
|  |                    Store Layer                      |    |
|  |          MemoryStore        /    RedisStore          |    |
|  +-----------------------------------------------------+    |
+-------------------------------------------------------------+
```

## Core Concepts

### Bot

The Bot is the entry point for all TeleRiHa applications:

```go
b := bot.New(token,
    bot.WithLogger(log.Logger),
    bot.WithBaseURL("https://api.telegram.org"),
    bot.WithHTTPClient(httpClient),
)
```

### Router

Routes incoming updates to handlers:

```go
// Command routing
b.Router.On("/start", handleStart)
b.Router.On("/help", handleHelp)

// Regex matching
b.Router.OnRegex(`^hello$`, handleHello)

// Callback query routing
b.Router.OnCallback("confirm_", handleConfirm)

// Inline query
b.Router.OnInlineQuery(handleInlineQuery)

// Default fallback
b.Router.Default(handleDefault)
```

### Middleware

Apply cross-cutting concerns:

```go
b.Router.Use(bot.Logger(zerolog.Logger))
b.Router.Use(bot.Recover(zerolog.Logger))
b.Router.Use(bot.OnlyPrivate())  // Restrict to private chats
b.Router.Use(bot.OnlyGroups())   // Restrict to groups
```

### Context

Access update data and reply:

```go
func handler(ctx *bot.Context) error {
    // Access data
    text := ctx.Text()
    chatID := ctx.ChatID()
    senderID := ctx.SenderID()

    // Check update type
    if ctx.IsCommand() {
        cmd := ctx.Command()      // e.g., "start"
        args := ctx.CommandArgs() // e.g., "arg1 arg2"
    }

    // Reply
    ctx.Reply("Hello!")

    // Answer callback
    ctx.AnswerCallback("Clicked!", true)

    return nil
}
```

### Keyboards

Fluent keyboard builders:

```go
// Inline keyboard
kb := bot.NewInlineKeyboardBuilder().
    AddRow(bot.InlineButton("Yes", "confirm:yes"), bot.InlineButton("No", "confirm:no")).
    Build()
ctx.Reply("Continue?", bot.WithReplyMarkup(kb))

// Reply keyboard
rk := bot.NewReplyKeyboardBuilder().
    Resize().OneTime().
    AddRow(bot.ReplyButton("Option 1"), bot.ReplyButton("Option 2")).
    Build()
ctx.Reply("Choose:", bot.WithReplyMarkup(rk))
```

### Conversations

Multi-step conversation flows:

```go
cm := b.ConversationManager

cb := bot.NewConversationBuilder(cm, "order")
cb.Start(startHandler).
    Step("waiting_name", waitForName).
    Step("waiting_address", waitForAddress).
    Step("confirm", confirmOrder).
    Build()

// In handler
func waitForName(ctx *bot.Context, conv *bot.Conversation) error {
    name := ctx.Text()
    cm.SetData(ctx.SenderID(), "name", name)
    cm.UpdateState(ctx.SenderID(), "waiting_address")
    return ctx.Reply("Got it! Now enter your address:")
}
```

### State Storage

Pluggable storage backends:

```go
// In-memory (default)
store := store.NewMemoryStore()

// Redis (production)
store := store.NewRedisStore(&redis.Options{
    Addr: "localhost:6379",
})

// Namespaced stores
userStore := store.NewUserStore(store)
userStore.Set(userID, "preferences", data, time.Hour)

chatStore := store.NewChatStore(store)
sessionStore := store.NewSessionStore(store, 24*time.Hour)
```

## Plugins

TeleRiHa includes production-ready plugins:

### Rate Limiting

```go
import "github.com/Yashwanth-Kumar-26/teleriha/plugins/ratelimit"

b.PluginRegistry.RegisterPlugin(ratelimit.New(ratelimit.Config{
    MaxRequests: 5,
    Interval:    time.Minute,
}))
```

### Analytics

```go
import "github.com/Yashwanth-Kumar-26/teleriha/plugins/analytics"

b.PluginRegistry.RegisterPlugin(analytics.New())
```

### Internationalization

```go
import "github.com/Yashwanth-Kumar-26/teleriha/plugins/i18n"

i18nPlugin := i18n.New("en")
i18nPlugin.AddTranslation("en", "hello", "Hello!")
i18nPlugin.AddTranslation("es", "hello", "Hola!")
b.PluginRegistry.RegisterPlugin(i18nPlugin)
```

### Scheduler

```go
import "github.com/Yashwanth-Kumar-26/teleriha/plugins/scheduler"

sched := scheduler.New()
sched.AddJob("daily", "0 0 * * *", func() error {
    return sendDailyDigest()
})
b.PluginRegistry.RegisterPlugin(sched)
```

## Comparison with Alternatives

| Feature | TeleRiHa | go-telegram-bot-api | grammY| telebot |
|---------|----------|---------------------|------|---------|
| Pure net/http | Yes | No | No | No |
| Zero deps | Yes | No | No | No |
| Gin-like API | Yes | No | No | Yes |
| Plugin system | Yes | No | Yes | No |
| Built-in store | Yes | No | No | No |
| Conversation builder | Yes | No | Yes | No |

## Project Structure

```
TeleRiHa/
|-- cmd/
|   +-- riha/                  # CLI tool
|-- pkg/
|   +-- bot/                   # Core framework
|       +-- teleriha.go        # Bot, polling, webhooks
|       +-- context.go         # Update context
|       +-- router.go          # Handler routing
|       +-- middleware.go      # Middleware chain
|       +-- keyboard.go        # Keyboard builders
|       +-- plugin.go          # Plugin interface
|       +-- conversation.go    # Conversation flows
|       +-- methods.go         # API methods
|       +-- types.go           # Telegram types
|   +-- store/                 # Storage backends
|       +-- store.go           # Store interface
|       +-- memory.go          # In-memory implementation
|       +-- redis.go           # Redis implementation
|-- plugins/                   # Built-in plugins
|   +-- ratelimit/            # Rate limiting
|   +-- analytics/             # Usage analytics
|   +-- i18n/                 # Translations
|   +-- scheduler/            # Cron jobs
|-- examples/                  # Example bots
|   +-- starter/              # Simple bot
|   +-- ecommerce/           # E-commerce bot
|   +-- ai-bot/               # AI-powered bot
|-- docs/                     # Documentation
|-- tests/                    # Test utilities
+-- go.mod
```

## Configuration

### Bot Options

```go
bot.New(token,
    bot.WithLogger(logger),           // Custom zerolog logger
    bot.WithBaseURL("https://api..."), // Custom API base URL
    bot.WithHTTPClient(client),        // Custom HTTP client
)
```

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| BOT_TOKEN | Yes | Telegram bot token from @BotFather |
| WEBHOOK_URL | No | Public URL for webhook mode |
| REDIS_URL | No | Redis connection for distributed storage |

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/bot/...
```

## Contributing

Contributions are welcome! See CONTRIBUTING.md for guidelines.

## License

TeleRiHa is released under the MIT License. See LICENSE for details.