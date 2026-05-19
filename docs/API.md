# TeleRiHa API Reference

Complete API documentation for TeleRiHa v1.0.0.

## Table of Contents

- [Bot](#bot)
- [Router](#router)
- [Context](#context)
- [Middleware](#middleware)
- [Keyboards](#keyboards)
- [Plugin System](#plugin-system)
- [Conversation Manager](#conversation-manager)
- [State Storage](#state-storage)

---

## Bot

### bot.New

```go
func New(token string, opts ...BotOption) *Bot
```

Creates a new Bot instance with the given Telegram token.

**Parameters:**
- token (string) - Bot token from @BotFather

**Options:**
- WithLogger(zerolog.Logger) - Set custom logger
- WithBaseURL(string) - Custom API base URL
- WithHTTPClient(*http.Client) - Custom HTTP client

**Example:**
```go
b := bot.New("123456:ABC-DEF...")
b := bot.New(token, bot.WithLogger(log.Logger))
```

### bot.Bot.StartPolling

```go
func (b *Bot) StartPolling(pollInterval, timeout time.Duration) error
```

Starts long polling for updates.

**Parameters:**
- pollInterval (time.Duration) - Time between poll requests
- timeout (time.Duration) - Request timeout for Telegram

**Example:**
```go
b.StartPolling(1*time.Second, 30*time.Second)
```

### bot.Bot.StartWebhook

```go
func (b *Bot) StartWebhook(path string, port int) error
```

Starts webhook server for receiving updates.

**Parameters:**
- path (string) - Webhook endpoint path
- port (int) - Server port number

### bot.Bot.Stop

```go
func (b *Bot) Stop()
```

Gracefully stops the bot.

### bot.Bot.GetMe

```go
func (b *Bot) GetMe() (*User, error)
```

Returns information about the bot.

---

## Sending Messages

### bot.Bot.SendMessage

```go
func (b *Bot) SendMessage(chatID int64, text string, opts ...MessageOption) (*Message, error)
```

Sends a text message.

**Options:**
- WithParseMode("Markdown"|"HTML") - Message formatting
- WithReplyMarkup(InlineKeyboardMarkup) - Inline keyboard
- WithDisableNotification() - Mute notification

### bot.Bot.SendPhoto

```go
func (b *Bot) SendPhoto(chatID int64, photo interface{}, caption string, opts ...SendPhotoOption) (*Message, error)
```

Sends a photo. Accepts file ID, URL, or *os.File.

### bot.Bot.SendDocument

```go
func (b *Bot) SendDocument(chatID int64, document interface{}, caption string, opts ...SendDocumentOption) (*Message, error)
```

Sends a document.

### bot.Bot.SendLocation

```go
func (b *Bot) SendLocation(chatID int64, latitude, longitude float64, opts ...MessageOption) (*Message, error)
```

Sends a location.

### bot.Bot.SendContact

```go
func (b *Bot) SendContact(chatID int64, phoneNumber, firstName string, opts ...MessageOption) (*Message, error)
```

Sends a contact.

### bot.Bot.SendPoll

```go
func (b *Bot) SendPoll(chatID int64, question string, options []string, opts ...MessageOption) (*Message, error)
```

Sends a poll.

---

## Editing & Deleting

### bot.Bot.DeleteMessage

```go
func (b *Bot) DeleteMessage(chatID int64, messageID int64) error
```

Deletes a message.

### bot.Bot.AnswerCallbackQuery

```go
func (b *Bot) AnswerCallbackQuery(callbackQueryID string, text string, showAlert bool) error
```

Answers a callback query.

### bot.Bot.EditMessageText

```go
func (b *Bot) EditMessageText(chatID int64, messageID int64, text string, opts ...EditMessageOption) (*Message, error)
```

Edits message text.

---

## Router

### bot.Router.On

```go
func (r *Router) On(command string, handler Handler)
```

Registers a handler for a command. Commands are case-insensitive.

```go
b.Router.On("/start", handleStart)
b.Router.On("/help", handleHelp)
```

### bot.Router.OnRegex

```go
func (r *Router) OnRegex(pattern string, handler Handler) error
```

Registers a handler for messages matching a regex pattern.

```go
b.Router.OnRegex(`^hello$`, handleHello)
b.Router.OnRegex(`^\d+$`, handleNumber)  // Matches only numbers
```

### bot.Router.OnCallback

```go
func (r *Router) OnCallback(prefix string, handler Handler)
```

Registers a handler for callback queries with a prefix match.

```go
b.Router.OnCallback("confirm_", handleConfirm)
// Matches "confirm_yes", "confirm_no", etc.
```

### bot.Router.OnInlineQuery

```go
func (r *Router) OnInlineQuery(handler Handler)
```

Registers a handler for inline queries.

### bot.Router.Default

```go
func (r *Router) Default(handler Handler)
```

Registers a default handler for unhandled messages.

### bot.Router.Use

```go
func (r *Router) Use(middleware Middleware)
```

Adds global middleware to the router.

### bot.Router.Group

```go
func (r *Router) Group(group string, middlewares ...Middleware) *GroupRouter
```

Creates a route group with optional middlewares.

```go
admin := b.Router.Group("admin", bot.IsAdmin())
admin.On("/stats", handleStats)
```

---

## Context

Context provides access to the current update and convenience methods.

### Accessing Data

```go
ctx.ChatID()           // Get chat ID
ctx.SenderID()         // Get sender user ID
ctx.Text()             // Get message text
ctx.Command()          // Get command name (without /)
ctx.CommandArgs()      // Get command arguments
ctx.CallbackData()     // Get callback data
ctx.Sender             // Access sender User object
ctx.Chat               // Access Chat object
ctx.Message            // Access Message object
```

### Type Checking

```go
ctx.IsMessage()            // Check if update is a message
ctx.IsCommand()           // Check if message is a command
ctx.IsCallbackQuery()      // Check if update is a callback query
ctx.IsInlineQuery()        // Check if update is an inline query
```

### Sending Replies

```go
ctx.Reply(text string, opts ...MessageOption) (*Message, error)
ctx.Send(chatID int64, text string, opts ...MessageOption) (*Message, error)
ctx.AnswerCallback(text string, showAlert bool) error
ctx.Delete() error
```

### Logging

```go
ctx.Log()  // Returns logger with chat_id and user_id context
```

---

## Middleware

Middleware wraps handlers to add cross-cutting functionality.

### Built-in Middleware

#### bot.Logger

```go
func Logger(logger zerolog.Logger) Middleware
```

Logs all handled updates with duration.

#### bot.Recover

```go
func Recover(logger zerolog.Logger) Middleware
```

Recovers from panics and logs them.

#### bot.RateLimitMiddleware

```go
func RateLimitMiddleware(rateLimiter *RateLimiter) Middleware
```

Rate limits requests per user.

#### bot.OnlyGroups

```go
func OnlyGroups() Middleware
```

Only allows messages from groups.

#### bot.OnlyPrivate

```go
func OnlyPrivate() Middleware
```

Only allows messages from private chats.

#### bot.Chain

```go
func Chain(middlewares ...Middleware) Middleware
```

Combines multiple middlewares.

### Writing Custom Middleware

```go
func MyMiddleware() Middleware {
    return func(next Handler) Handler {
        return func(ctx *Context) error {
            // Before handler
            start := time.Now()

            // Call next handler
            err := next(ctx)

            // After handler
            log.Printf("Handler took %v", time.Since(start))

            return err
        }
    }
}
```

---

## Keyboards

### Inline Keyboards

#### bot.NewInlineKeyboardBuilder

```go
func NewInlineKeyboardBuilder() *InlineKeyboardBuilder
```

Creates a fluent inline keyboard builder.

```go
kb := bot.NewInlineKeyboardBuilder().
    AddRow(
        bot.InlineButton("Yes", "confirm:yes"),
        bot.InlineButton("No", "confirm:no"),
    ).
    AddURLButton("Open Site", "https://example.com").
    Build()
```

#### Inline Keyboard Button Helpers

```go
bot.InlineButton(text, callbackData string)           // Callback button
bot.InlineButtonURL(text, url string)                 // URL button
bot.InlineButtonSwitch(text, query string)           // Switch inline query
bot.InlineButtonPay(text string)                      // Payment button
```

#### Utility Keyboards

```go
bot.YesNoKeyboard(yesCallback, noCallback string)              // Yes/No buttons
bot.ConfirmCancelKeyboard(confirm, cancel string)              // Confirm/Cancel
bot.PaginationKeyboard(prev, next string, page, total int)    // Pagination
```

### Reply Keyboards

#### bot.NewReplyKeyboardBuilder

```go
func NewReplyKeyboardBuilder() *ReplyKeyboardBuilder
```

Creates a fluent reply keyboard builder.

```go
rk := bot.NewReplyKeyboardBuilder().
    Resize().OneTime().
    AddRow(
        bot.ReplyButton("Option 1"),
        bot.ReplyButton("Option 2"),
    ).
    Build()
```

#### Reply Keyboard Helpers

```go
bot.ReplyButton(text string)                    // Simple button
bot.ReplyButtonContact(text string)             // Request contact
bot.ReplyButtonLocation(text string)           // Request location
bot.RemoveKeyboard()                            // Remove keyboard
bot.HideKeyboard()                              // Hide keyboard
```

---

## Plugin System

### Plugin Interface

```go
type Plugin interface {
    Name() string
    Init(bot *Bot) error
    Start(ctx context.Context) error
    Stop() error
}
```

### BasePlugin

```go
type BasePlugin struct {
    name   string
    bot    *Bot
    active bool
}
```

Provides default implementations. Embed in your plugin:

```go
type MyPlugin struct {
    bot.BasePlugin
    config MyConfig
}

func New(config MyConfig) *MyPlugin {
    return &MyPlugin{
        BasePlugin: *bot.NewBasePlugin("my-plugin"),
        config:     config,
    }
}
```

### Plugin Registry

```go
b.PluginRegistry.Register(name string, fn PluginFunc)
b.PluginRegistry.RegisterPlugin(plugin Plugin)
b.PluginRegistry.LoadAll(bot *Bot) error
b.PluginRegistry.StartAll(ctx context.Context) error
b.PluginRegistry.StopAll() error
```

### Built-in Plugins

#### ratelimit

```go
import "github.com/Yashwanth-Kumar-26/teleriha/plugins/ratelimit"

b.PluginRegistry.RegisterPlugin(ratelimit.New(ratelimit.Config{
    MaxRequests: 5,
    Interval:    time.Minute,
}))
```

#### analytics

```go
import "github.com/Yashwanth-Kumar-26/teleriha/plugins/analytics"

b.PluginRegistry.RegisterPlugin(analytics.New())
```

---

## Conversation Manager

Manages multi-step conversation flows with state persistence.

### ConversationBuilder

```go
cm := b.ConversationManager

cb := bot.NewConversationBuilder(cm, "order-flow")
cb.Start(startHandler).
    Step("waiting_name", waitForName).
    Step("waiting_email", waitForEmail).
    Step("confirm", confirmOrder).
    Build()
```

### Conversation State

```go
type Conversation struct {
    ID     string                 // Conversation identifier
    UserID int64                 // User's Telegram ID
    ChatID int64                 // Chat ID
    State  string                 // Current state
    Data   map[string]interface{} // Conversation data
}
```

### Helper Methods

```go
cb.Next(nextState string)              // Transition to next state
cb.End()                               // End conversation
cb.WaitForText(key, nextState)         // Wait for text input
cb.WaitForCallback(data, nextState)    // Wait for callback button
```

---

## State Storage

### Store Interface

```go
type Store interface {
    Get(key string) ([]byte, error)
    Set(key string, value []byte, ttl time.Duration) error
    Delete(key string) error
    Exists(key string) bool
    Clear() error
}
```

### MemoryStore

```go
import "github.com/Yashwanth-Kumar-26/teleriha/pkg/store"

store := store.NewMemoryStore()
store.Set("key", []byte("value"), time.Hour)
data, _ := store.Get("key")
```

### RedisStore

```go
import (
    "github.com/Yashwanth-Kumar-26/teleriha/pkg/store"
    "github.com/redis/go-redis/v9"
)

store := store.NewRedisStore(&redis.Options{
    Addr: "localhost:6379",
    DB:   0,
})
```

### Store Wrappers

#### UserStore

```go
us := store.NewUserStore(baseStore)
us.Set(userID, "preferences", data, time.Hour)
```

#### ChatStore

```go
cs := store.NewChatStore(baseStore)
cs.Set(chatID, "settings", data, 0)
```

#### SessionStore

```go
ss := store.NewSessionStore(baseStore, 24*time.Hour)
ss.Set(sessionID, data)
```

### SQLiteStore

```go
import "github.com/Yashwanth-Kumar-26/teleriha/pkg/store"

// Creates or opens a persistent SQLite database file
sqliteStore, err := store.NewSQLiteStore("bot.db")
if err != nil {
    log.Fatal(err)
}
defer sqliteStore.Close()
```

---

## FSM Scene Manager

The State Machine Scene Manager simplifies multi-step conversations (like user onboarding or wizards).

### Scene Definition

```go
import "github.com/Yashwanth-Kumar-26/teleriha/plugins/scenes"

onboardingScene := scenes.NewScene("onboarding",
    // Step 1
    func(ctx *bot.Context) error {
        return ctx.Reply("Welcome! What is your name?")
    },
    // Step 2
    func(ctx *bot.Context) error {
        name := ctx.Message.Text
        // Save user name in session store
        sess := session.GetSession(ctx)
        sess.Set("name", []byte(name))
        return ctx.Reply("Got it! What is your age?")
    },
    // Step 3
    func(ctx *bot.Context) error {
        age := ctx.Message.Text
        // Retrieve name from session
        sess := session.GetSession(ctx)
        nameBytes, _ := sess.Get("name")
        
        ctx.Reply(fmt.Sprintf("Thanks %s! You are registered.", string(nameBytes)))
        return sceneManager.Leave(ctx) // Clean up scene state
    },
)
```

### Scene Setup & Middleware

```go
sceneManager := scenes.NewSceneManager()
sceneManager.AddScene(onboardingScene)

// Scene middleware must run after the session middleware
app.Bot.Router.Use(session.Middleware(sqliteStore, 24*time.Hour))
app.Bot.Router.Use(sceneManager.Middleware())
```

### Entering a Scene

```go
app.AddHandler(bot.CommandHandler("/register", func(ctx *bot.Context) error {
    return sceneManager.Enter(ctx, "onboarding")
}))
```

---

## Job Queue (Background Tasks)

The framework has an integrated Job Queue scheduler supporting delayed and periodic background tasks.

### Running a Periodic/Repeating Task

```go
app.JobQueue.RunRepeating(func(ctx *bot.JobContext) error {
    log.Println("Periodic cleanup task executed!")
    return nil
}, 1*time.Hour, time.Now())
```

### Running a One-Off Delayed Task

```go
app.JobQueue.RunOnce(func(ctx *bot.JobContext) error {
    log.Println("Delayed welcome task executed!")
    return nil
}, 5*time.Minute)
```

---

## Telegram Gateway API Client

Interacts with the Telegram Gateway API for high-deliverability phone verification.

```go
import "github.com/Yashwanth-Kumar-26/teleriha/pkg/gateway"

client := gateway.NewClient("YOUR_GATEWAY_TOKEN")

// Check if verification SMS can be sent
res, err := client.CheckSendAbility(context.Background(), "+1234567890")
if err == nil && res.Ok {
    fmt.Printf("OTP send ability verified. Request ID: %s\n", res.Result.RequestID)
}
```

---

## Payments & Telegram Stars

### Sending an Invoice (Telegram Stars)

For Telegram Stars payments, set the currency to `"XTR"` and omit the provider token:

```go
prices := []bot.LabeledPrice{
    {Label: "T-Shirt", Amount: 50}, // 50 Stars
}

_, err := app.Bot.SendInvoice(
    ctx.ChatID(),
    "Product Title",
    "Product Description",
    "payload_xyz",
    "", // Empty provider token for Stars
    "XTR",
    prices,
)
```

### Handling PreCheckout Queries

```go
app.AddHandler(bot.PreCheckoutQueryHandler(func(ctx *bot.Context) error {
    // Approve the checkout query
    return app.Bot.AnswerPreCheckoutQuery(ctx.PreCheckoutQuery.ID, true, "")
}))
```

---

## Pagination Keyboard Builder

Generates dynamic, interactive inline buttons for list pagination.

```go
pb := bot.NewPaginationBuilder(totalPages, currentPage, "callback_prefix_")
keyboard := pb.Build()

// Send paginated keyboard markup
_, err := ctx.Reply("Please choose page:", bot.WithReplyMarkup(keyboard))
```

---

## Message Filters

Filters are used to conditionally route messages to specific handlers.

```go
// Filter photo messages
app.AddHandler(bot.MessageHandler(bot.FilterPhoto(), func(ctx *bot.Context) error {
    return ctx.Reply("Photo received!")
}))

// Logical operators combination
filterOr := bot.FilterOr(bot.FilterPhoto(), bot.FilterDocument())
app.AddHandler(bot.MessageHandler(filterOr, func(ctx *bot.Context) error {
    return ctx.Reply("Media or Document received!")
}))
```