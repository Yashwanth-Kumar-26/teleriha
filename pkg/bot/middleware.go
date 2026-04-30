package bot

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Logger is a middleware that logs all handled updates.
func Logger(logger zerolog.Logger) Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) error {
			start := time.Now()

			// Build a logger with context
			event := logger.Info().
				Int64("user_id", ctx.SenderID()).
				Int64("chat_id", ctx.ChatID()).
				Str("type", ctx.updateType())

			if ctx.IsCommand() {
				event = event.Str("command", ctx.Command())
			}

			event.Msg("Handling update")

			// Call the next handler
			err := next(ctx)

			// Log completion
			logger.Info().
				Int64("user_id", ctx.SenderID()).
				Int64("chat_id", ctx.ChatID()).
				Str("type", ctx.updateType()).
				Dur("duration", time.Since(start)).
				Err(err).
				Msg("Update handled")

			return err
		}
	}
}

//updateType returns the type of update as a string.
func (c *Context) updateType() string {
	switch {
	case c.IsMessage():
		return "message"
	case c.IsCallbackQuery():
		return "callback_query"
	case c.IsInlineQuery():
		return "inline_query"
	case c.IsChosenInlineResult():
		return "chosen_inline_result"
	default:
		return "unknown"
	}
}

// Recover is a middleware that recovers from panics.
func Recover(logger zerolog.Logger) Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					// Convert panic to error
					if e, ok := r.(error); ok {
						err = e
					} else {
						err = fmt.Errorf("panic: %v", r)
					}

					// Log the panic
					logger.Err(err).
						Str("type", ctx.updateType()).
						Int64("user_id", ctx.SenderID()).
						Int64("chat_id", ctx.ChatID()).
						Msg("Recovered from panic")
				}
			}()

			return next(ctx)
		}
	}
}

// RateLimit is a middleware that implements rate limiting per user.
// It uses a simple token bucket algorithm.
// This is a basic in-memory implementation; for production, use a distributed store.
type RateLimiter struct {
	// limits maps user IDs to their rate limit state
	limits map[int64]*rateLimitState
	// maxRequests is the maximum number of requests per interval
	maxRequests int
	// interval is the time window for rate limiting
	interval time.Duration
	// mu protects the limits map
	mu sync.Mutex
}

type rateLimitState struct {
	count     int
	windowStart time.Time
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(maxRequests int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		limits:     make(map[int64]*rateLimitState),
		maxRequests: maxRequests,
		interval:   interval,
	}
}

// RateLimitMiddleware creates a rate limiting middleware.
func RateLimitMiddleware(rateLimiter *RateLimiter) Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) error {
			userID := ctx.SenderID()
			if userID == 0 {
				// Allow messages without a sender (e.g., channel posts)
				return next(ctx)
			}

			rateLimiter.mu.Lock()
			defer rateLimiter.mu.Unlock()

			state := rateLimiter.limits[userID]
			now := time.Now()

			// Check if we need to reset the window
			if state == nil || now.Sub(state.windowStart) >= rateLimiter.interval {
				state = &rateLimitState{
					count:     0,
					windowStart: now,
				}
				rateLimiter.limits[userID] = state
			}

			// Check if rate limit is exceeded
			if state.count >= rateLimiter.maxRequests {
				return fmt.Errorf("rate limit exceeded: %d requests per %v", rateLimiter.maxRequests, rateLimiter.interval)
			}

			// Increment counter
			state.count++

			return next(ctx)
		}
	}
}

// Cleanup removes old rate limit states.
func (r *RateLimiter) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for userID, state := range r.limits {
		if now.Sub(state.windowStart) >= r.interval {
			delete(r.limits, userID)
		}
	}
}

// OnlyGroups is a middleware that only allows messages from groups.
func OnlyGroups() Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) error {
			if ctx.Chat == nil || ctx.Chat.Type == "private" {
				return fmt.Errorf("this command can only be used in groups")
			}
			return next(ctx)
		}
	}
}

// OnlyPrivate is a middleware that only allows messages from private chats.
func OnlyPrivate() Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) error {
			if ctx.Chat == nil || ctx.Chat.Type != "private" {
				return fmt.Errorf("this command can only be used in private chats")
			}
			return next(ctx)
		}
	}
}

// IsAdmin is a middleware that checks if the sender is an admin in the chat.
func IsAdmin() Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) error {
			if ctx.Message == nil || ctx.Message.From == nil {
				return fmt.Errorf("cannot check admin status")
			}

			// For now, just check if the user is the bot owner
			// In a real implementation, you'd need to get chat members
			// and check their status
			// This is a placeholder implementation
			if !ctx.isBotOwner() {
				return fmt.Errorf("you must be an admin to use this command")
			}

			return next(ctx)
		}
	}
}

// isBotOwner checks if the sender is the bot owner (placeholder).
// In a real implementation, you'd configure the owner ID.
func (c *Context) isBotOwner() bool {
	// Placeholder: replace with actual owner ID check
	return false
}

// Chain combines multiple middlewares into one.
func Chain(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		finalHandler := next
		for i := len(middlewares) - 1; i >= 0; i-- {
			finalHandler = middlewares[i](finalHandler)
		}
		return finalHandler
	}
}
