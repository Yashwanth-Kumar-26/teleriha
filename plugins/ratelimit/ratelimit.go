package ratelimit

import (
	"context"
	"time"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

// Config holds the configuration for the ratelimit plugin.
type Config struct {
	MaxRequests int
	Interval    time.Duration
}

// RateLimitPlugin implements the bot.Plugin interface.
type RateLimitPlugin struct {
	bot.BasePlugin
	config  Config
	limiter *bot.RateLimiter
}

// New creates a new RateLimitPlugin.
func New(config Config) *RateLimitPlugin {
	return &RateLimitPlugin{
		BasePlugin: *bot.NewBasePlugin("ratelimit"),
		config:     config,
		limiter:    bot.NewRateLimiter(config.MaxRequests, config.Interval),
	}
}

// Init initializes the plugin.
func (p *RateLimitPlugin) Init(b *bot.Bot) error {
	if err := p.BasePlugin.Init(b); err != nil {
		return err
	}

	// Add the ratelimit middleware to the bot's router
	b.Router.Use(bot.RateLimitMiddleware(p.limiter))
	return nil
}

// Start starts the plugin.
func (p *RateLimitPlugin) Start(ctx context.Context) error {
	// Start a cleanup goroutine
	go func() {
		ticker := time.NewTicker(p.config.Interval * 2)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				p.limiter.Cleanup()
			}
		}
	}()
	return p.BasePlugin.Start(ctx)
}
