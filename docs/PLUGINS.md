# TeleRiHa Plugin Development Guide

Learn how to create custom plugins for TeleRiHa.

## Table of Contents

- [Plugin Interface](#plugin-interface)
- [BasePlugin](#baseplugin)
- [Creating a Plugin](#creating-a-plugin)
- [Plugin Lifecycle](#plugin-lifecycle)
- [Examples](#examples)

---

## Plugin Interface

All TeleRiHa plugins must implement the Plugin interface:

```go
type Plugin interface {
    Name() string
    Init(bot *Bot) error
    Start(ctx context.Context) error
    Stop() error
}
```

| Method | Description |
|--------|-------------|
| Name() | Returns the plugin's unique name |
| Init(bot) | Called when plugin is loaded |
| Start(ctx) | Called when bot starts |
| Stop() | Called when bot stops |

---

## BasePlugin

TeleRiHa provides BasePlugin with default implementations:

```go
type BasePlugin struct {
    name   string
    bot    *Bot
    active bool
}

func NewBasePlugin(name string) *BasePlugin
```

Embed it to get sensible defaults:

```go
type MyPlugin struct {
    *bot.BasePlugin
    // Your fields
}
```

---

## Creating a Plugin

### Step 1: Define Structure

```go
package myplugin

import "context"

type MyPlugin struct {
    *bot.BasePlugin
    config Config
}

type Config struct {
    Setting1 string
    Setting2 int
}
```

### Step 2: Create Constructor

```go
func New(config Config) *MyPlugin {
    return &MyPlugin{
        BasePlugin: bot.NewBasePlugin("my-plugin"),
        config:     config,
    }
}
```

### Step 3: Implement Init

```go
func (p *MyPlugin) Init(b *bot.Bot) error {
    if err := p.BasePlugin.Init(b); err != nil {
        return err
    }
    b.Router.Use(p.MyMiddleware())
    return nil
}
```

### Step 4: Implement Start

```go
func (p *MyPlugin) Start(ctx context.Context) error {
    if err := p.BasePlugin.Start(ctx); err != nil {
        return err
    }
    go p.backgroundTask(ctx)
    return nil
}
```

### Step 5: Implement Stop

```go
func (p *MyPlugin) Stop() error {
    p.active = false
    return nil
}
```

---

## Plugin Lifecycle

```
bot.New()           -> Create bot instance
RegisterPlugin()    -> Register plugins
LoadAll()           -> Call Init() on each plugin
StartAll()          -> Call Start() on each plugin
StartPolling()      -> Begin receiving updates

--- On Shutdown ---

Stop()              -> Signal all goroutines to stop
StopAll()           -> Call Stop() on each plugin
```

---

## Examples

### Logging Plugin

```go
package logplugin

import (
    "context"
    "time"

    "github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
    "github.com/rs/zerolog"
)

type LogPlugin struct {
    *bot.BasePlugin
    logger zerolog.Logger
}

func New(logger zerolog.Logger) *LogPlugin {
    return &LogPlugin{
        BasePlugin: bot.NewBasePlugin("logging"),
        logger:     logger,
    }
}

func (p *LogPlugin) Init(b *bot.Bot) error {
    p.BasePlugin.Init(b)
    b.Router.Use(p.loggingMiddleware())
    return nil
}

func (p *LogPlugin) loggingMiddleware() bot.Middleware {
    return func(next bot.Handler) bot.Handler {
        return func(ctx *bot.Context) error {
            start := time.Now()
            p.logger.Info().
                Int64("user_id", ctx.SenderID()).
                Msg("Update received")

            err := next(ctx)

            p.logger.Info().
                Int64("user_id", ctx.SenderID()).
                Dur("duration", time.Since(start)).
                Err(err).
                Msg("Update processed")

            return err
        }
    }
}
```

### Reminder Plugin

```go
package reminder

import (
    "context"
    "fmt"
    "time"

    "github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

type ReminderPlugin struct {
    *bot.BasePlugin
    interval time.Duration
    message  string
}

func New(interval time.Duration, message string) *ReminderPlugin {
    return &ReminderPlugin{
        BasePlugin: bot.NewBasePlugin("reminder"),
        interval:   interval,
        message:    message,
    }
}

func (p *ReminderPlugin) Start(ctx context.Context) error {
    p.BasePlugin.Start(ctx)

    ticker := time.NewTicker(p.interval)
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                p.Bot().SendMessage(adminID, "Reminder: "+p.message)
            }
        }
    }()
    return nil
}
```

---

## Using Custom Plugins

```go
b := bot.New(token)

// Register plugin
b.PluginRegistry.RegisterPlugin(myplugin.New(config))
b.PluginRegistry.LoadAll(b)
```

---

## Best Practices

1. **Always Call BasePlugin Methods**
   ```go
   func (p *MyPlugin) Init(b *bot.Bot) error {
       if err := p.BasePlugin.Init(b); err != nil {
           return err
       }
       // ...
   }
   ```

2. **Handle Context Cancellation**
   ```go
   func (p *MyPlugin) Start(ctx context.Context) error {
       go func() {
           for {
               select {
               case <-ctx.Done():
                   return
               // ...
               }
           }
       }()
       return nil
   }
   ```

3. **Document Configuration**
   ```go
   // Config holds the configuration for MyPlugin.
   type Config struct {
       Setting1 string
       Setting2 int
   }
   ```