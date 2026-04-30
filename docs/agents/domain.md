# TeleRiHa Domain Language Glossary

This document provides a detailed explanation of the TeleRiHa domain language, including code examples and relationships between concepts.

## Core Terms

### Bot
The main entry point of the framework. It holds the configuration, the HTTP client, and the router.
```go
b := bot.New("token")
b.StartPolling(time.Second, time.Minute)
```

### Context
A wrapper around the Telegram update and the bot instance. It provides scoped methods for replying and accessing sender info.
```go
func handler(ctx *bot.Context) error {
    return ctx.Reply("Hello!")
}
```

### Router
The component responsible for matching updates to handlers. It supports commands, regex, and callback queries.
```go
b.Router.On("/start", startHandler)
b.Router.OnRegex(`^/help`, helpHandler)
```

### Plugin
A modular unit that can be registered to add functionality. Plugins have a lifecycle: `Init`, `Start`, and `Stop`.
```go
type MyPlugin struct { bot.BasePlugin }
b.PluginRegistry.RegisterPlugin(&MyPlugin{})
```

### Store
The abstraction layer for state management.
```go
s := store.NewMemoryStore()
s.Set("key", []byte("value"), time.Hour)
```

### Handler
A function signature `func(*Context) error`. It is the building block of bot logic.

### Middleware
A function that wraps a handler: `func(Handler) Handler`. Used for logging, recovery, and auth.

## Concept Relationships

- A **Bot** has one **Router** and one **PluginRegistry**.
- The **Router** contains many **Handlers** (and optionally **Middlewares**).
- Every **Handler** receives a **Context** when triggered.
- **Plugins** often interact with the **Store** to persist state.
- **Middlewares** are applied to **Handlers** to inject logic before/after execution.
