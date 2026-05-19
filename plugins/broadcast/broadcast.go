package broadcast

import (
	"context"
	"sync"
	"time"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

// Broadcaster manages mass message broadcasting safely honoring Telegram rate limits.
type Broadcaster struct {
	bot        *bot.Bot
	rateLimiter *time.Ticker
	mu         sync.Mutex
	status     BroadcastStatus
}

// BroadcastStatus tracks the progress of a broadcast.
type BroadcastStatus struct {
	Total      int
	Successful int
	Failed     int
	IsRunning  bool
}

// New creates a new Broadcaster instance.
func New(b *bot.Bot) *Broadcaster {
	return &Broadcaster{
		bot: b,
		// Standard Telegram rate limit for broadcast is 30 msg/sec
		rateLimiter: time.NewTicker(time.Second / 30),
	}
}

// Broadcast sends a message to a list of users.
func (b *Broadcaster) Broadcast(ctx context.Context, userIDs []int64, text string, opts ...bot.MessageOption) BroadcastStatus {
	b.mu.Lock()
	b.status = BroadcastStatus{Total: len(userIDs), IsRunning: true}
	b.mu.Unlock()

	go func() {
		defer func() {
			b.mu.Lock()
			b.status.IsRunning = false
			b.mu.Unlock()
		}()

		for _, userID := range userIDs {
			select {
			case <-ctx.Done():
				return
			case <-b.rateLimiter.C:
				_, err := b.bot.SendMessage(userID, text, opts...)
				
				b.mu.Lock()
				if err != nil {
					b.status.Failed++
				} else {
					b.status.Successful++
				}
				b.mu.Unlock()
			}
		}
	}()

	b.mu.Lock()
	defer b.mu.Unlock()
	return b.status
}

// Status returns the current status of the broadcast.
func (b *Broadcaster) Status() BroadcastStatus {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.status
}
