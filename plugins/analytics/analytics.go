package analytics

import (
	"sync"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

// Stats holds the analytics data.
type Stats struct {
	TotalMessages int64
	TotalUsers    int
}

// AnalyticsPlugin implements the bot.Plugin interface.
type AnalyticsPlugin struct {
	bot.BasePlugin
	mu            sync.RWMutex
	totalMessages int64
	users         map[int64]struct{}
}

// New creates a new AnalyticsPlugin.
func New() *AnalyticsPlugin {
	return &AnalyticsPlugin{
		BasePlugin: *bot.NewBasePlugin("analytics"),
		users:      make(map[int64]struct{}),
	}
}

// Init initializes the plugin.
func (p *AnalyticsPlugin) Init(b *bot.Bot) error {
	if err := p.BasePlugin.Init(b); err != nil {
		return err
	}

	// Add analytics middleware
	b.Router.Use(func(next bot.Handler) bot.Handler {
		return func(ctx *bot.Context) error {
			p.track(ctx)
			return next(ctx)
		}
	})
	return nil
}

// track records the update in stats.
func (p *AnalyticsPlugin) track(ctx *bot.Context) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.totalMessages++
	if ctx.Sender != nil {
		p.users[ctx.Sender.ID] = struct{}{}
	}
}

// GetStats returns the current stats.
func (p *AnalyticsPlugin) GetStats() Stats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return Stats{
		TotalMessages: p.totalMessages,
		TotalUsers:    len(p.users),
	}
}
