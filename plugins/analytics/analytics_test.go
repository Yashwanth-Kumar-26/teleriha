package analytics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

func TestAnalyticsPlugin(t *testing.T) {
	b := bot.New("test-token")
	
	// Create plugin
	p := New()
	
	// Init
	err := p.Init(b)
	assert.NoError(t, err)
	
	// Register a test command
	b.Router.On("/test", func(ctx *bot.Context) error {
		return nil
	})
	
	// Mock message
	ctx := &bot.Context{
		Sender: &bot.User{ID: 1},
		Chat:   &bot.Chat{ID: 100},
	}
	msg := bot.Message{
		Text: "/test",
		From: ctx.Sender,
		Chat: ctx.Chat,
		Entities: []bot.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: 5},
		},
	}

	// Handle a few messages
	b.Router.HandleMessage(ctx, msg)
	b.Router.HandleMessage(ctx, msg)
	
	// Check stats
	stats := p.GetStats()
	assert.Equal(t, int64(2), stats.TotalMessages)
	assert.Equal(t, 1, stats.TotalUsers)
}
