[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

# TeleRiHa - Telegram Bot Framework for Go

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-Apache%202.0-blue?style=flat-square" alt="License">
  <img src="https://img.shields.io/badge/Telegram%20API-v7.8-blue?style=flat-square" alt="Telegram API">
</p>

Build production-grade Telegram bots with a clean, modular, and developer-friendly framework. TeleRiHa provides a Gin-like API with zero external Telegram dependencies.

## Features

| Feature | Description |
|---------|-------------|
| **Pure `net/http`** | Direct communication with Telegram Bot API using Go's built-in HTTP package |
| **Zero Dependencies** | No wrappers or abstractions - just pure Telegram API calls |
| **Fluent Application Builder** | Clean `ApplicationBuilder` & handler registry architecture |
| **Stateful FSM Scene Manager** | Complex multi-step user conversations layered on top of Session Middleware |
| **Composable Filters** | Flexible message filters supporting logical operators (`FilterAnd`, `FilterOr`, `FilterNot`) |
| **Robust Job Queue** | Native background task executor supporting one-off and repeating jobs |
| **Enterprise Payments** | Full support for standard payments and Telegram Stars (`XTR`), pre-checkout query routing |
| **Telegram Gateway API** | Decoupled client for low-cost, direct verification/OTP messaging |
| **Multipart Media Uploads** | Dynamic media detection and local file uploads |
| **Inline Keyboards** | Fluent builder API for inline/reply keyboards and paginated interfaces |

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

    // 1. Build the Application
    app := bot.NewApplicationBuilder().
        Token(os.Getenv("BOT_TOKEN")).
        Opts(bot.WithLogger(log.Logger)).
        Build()

    // 2. Add Middlewares
    app.Bot.Router.Use(bot.Logger(log.Logger))
    app.Bot.Router.Use(bot.Recover(log.Logger))

    // 3. Register Handlers
    app.AddHandler(bot.CommandHandler("/start", handleStart))
    app.AddHandler(bot.CommandHandler("/help", handleHelp))
    
    // Add filtering handler
    app.AddHandler(bot.MessageHandler(bot.FilterPhoto(), func(ctx *bot.Context) error {
        return ctx.Reply("Thanks for the photo!")
    }))

    // 4. Start background jobs
    app.JobQueue.RunRepeating(func(ctx *bot.JobContext) error {
        log.Info().Msg("Performing periodic maintenance...")
        return nil
    }, 1*time.Hour, time.Now())

    // 5. Start long polling
    if err := app.RunPolling(); err != nil {
        log.Fatal().Err(err).Msg("Failed to run bot")
    }
}

func handleStart(ctx *bot.Context) error {
    // Store user session/state via Context Data
    userData := ctx.UserData()
    userData["visits"] = 1

    return ctx.Reply("Welcome! Type /help for commands.")
}

func handleHelp(ctx *bot.Context) error {
    return ctx.Reply("Available commands: /start, /help")
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

### Composable Filters

TeleRiHa provides a robust message filtering system. You can compose complex logic using logical operators:

```go
// Match only photos from private chats
photoPrivateFilter := bot.FilterAnd(bot.FilterPhoto(), bot.FilterPrivate())

app.AddHandler(bot.MessageHandler(photoPrivateFilter, func(ctx *bot.Context) error {
    return ctx.Reply("Received a private photo!")
}))

// Match text messages OR videos
textOrVideoFilter := bot.FilterOr(bot.FilterText(), bot.FilterVideo())
```

### Thread-Safe Context Data

Context data provides easy state management inside individual user, chat, or global scope:

```go
func handleStart(ctx *bot.Context) error {
    // User-specific data (persists per user ID)
    userData := ctx.UserData()
    userData["name"] = ctx.Sender.FirstName

    // Chat-specific data (persists per chat ID)
    chatData := ctx.ChatData()
    chatData["active"] = true

    // Global application data
    botData := ctx.BotData()
    botData["system_status"] = "online"
    
    return ctx.Reply("Context data stored!")
}
```

### Job Queue

The Job Queue handles scheduling background tasks:

```go
// Run once after 5 seconds
app.JobQueue.RunOnce(func(ctx *bot.JobContext) error {
    log.Info().Msg("One-off job running!")
    return nil
}, time.Now().Add(5*time.Second))

// Run repeating every 30 seconds
app.JobQueue.RunRepeating(func(ctx *bot.JobContext) error {
    log.Info().Msg("Repeating task executing...")
    return nil
}, 30*time.Second, time.Now())
```

### FSM Scene Manager

For multi-step conversation flows (like ordering or onboarding), you can register scenes:

```go
import "github.com/Yashwanth-Kumar-26/teleriha/plugins/scenes"

// 1. Create a scene manager
sm := scenes.NewSceneManager()

// 2. Define scene steps
onboardingScene := scenes.NewScene("onboarding",
    func(ctx *bot.Context) error {
        return ctx.Reply("Welcome! What is your name?")
    },
    func(ctx *bot.Context) error {
        name := ctx.Text()
        ctx.UserData()["name"] = name
        return ctx.Reply("Great! How old are you?")
    },
    func(ctx *bot.Context) error {
        age := ctx.Text()
        ctx.UserData()["age"] = age
        
        // Scene automatically completes after last step
        return ctx.Reply("Registration complete!")
    },
)
sm.AddScene(onboardingScene)

// 3. Register scenes middleware
app.Bot.Router.Use(sm.Middleware())

// 4. Enter user into the scene from a command
app.AddHandler(bot.CommandHandler("/register", func(ctx *bot.Context) error {
    return sm.Enter(ctx, "onboarding")
}))
```

### Stateful Session Storage

Pluggable storage backends backing persistent state and scenes:

```go
// SQLite backend for persistent sessions
sqliteStore, _ := store.NewSQLiteStore("bot.db")

app := bot.NewApplicationBuilder().
    Token(token).
    Opts(bot.WithStore(sqliteStore)). // Persists scenes state
    Build()

// Attach session middleware
app.Bot.Router.Use(session.Middleware(sqliteStore, 24*time.Hour))
```

### Enterprise Payments & Telegram Stars (`XTR`)

Fully supports standard currency billing and Telegram Stars:

```go
app.AddHandler(bot.CommandHandler("/invoice", func(ctx *bot.Context) error {
    prices := []bot.LabeledPrice{
        {Label: "Premium access", Amount: 100}, // 100 Stars
    }
    
    // Send invoice using Telegram Stars (XTR)
    _, err := app.Bot.SendInvoice(
        ctx.ChatID(),
        "Premium Upgrade",
        "Unlock premium enterprise capabilities.",
        "payload_upgrade",
        "", // Empty provider token for Telegram Stars
        "XTR",
        prices,
    )
    return err
}))

// Answer Pre-Checkout query
app.AddHandler(bot.PreCheckoutQueryHandler(func(ctx *bot.Context) error {
    return app.Bot.AnswerPreCheckoutQuery(ctx.PreCheckoutQuery.ID, true, "")
}))
```

### Telegram Gateway API (OTP/Verification)

Decoupled client to send direct OTP/verification text messages at cost:

```go
import "github.com/Yashwanth-Kumar-26/teleriha/pkg/gateway"

client := gateway.NewClient("GATEWAY_TOKEN")

// Verification check request
result, err := client.SendVerificationMessage("+1234567890", gateway.VerificationOptions{
    CodeLength: 6,
})
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

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

TeleRiHa is released under the Apache License 2.0. See [LICENSE](LICENSE) for details.