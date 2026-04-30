package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

func TestRateLimitPlugin(t *testing.T) {
	b := bot.New("test-token")
	
	// Create plugin
	p := New(Config{
		MaxRequests: 2,
		Interval:    1 * time.Second,
	})
	
	// Init and Start
	err := p.Init(b)
	assert.NoError(t, err)
	
	err = p.Start(context.Background())
	assert.NoError(t, err)
	
	// Register a test command
	called := 0
	b.Router.On("/test", func(ctx *bot.Context) error {
		called++
		return nil
	})
	
	// Mock context for user 1
	ctx := &bot.Context{
		Sender: &bot.User{ID: 1},
		Chat:   &bot.Chat{ID: 100},
	}
	
	// Mock message
	msg := bot.Message{
		Text: "/test",
		From: ctx.Sender,
		Chat: ctx.Chat,
		Entities: []bot.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: 5},
		},
	}

	// First 2 requests should pass
	b.Router.HandleMessage(ctx, msg)
	b.Router.HandleMessage(ctx, msg)
	
	// Check if called 2 times
	assert.Equal(t, 2, called)

	// Third request should be blocked by middleware
	b.Router.HandleMessage(ctx, msg)
	assert.Equal(t, 2, called)
}
