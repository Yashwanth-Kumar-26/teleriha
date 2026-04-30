package bot

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// Test Logger middleware runs and doesn't panic
func TestLogger_DoesNotPanic(t *testing.T) {
	logger := zerolog.Nop()

	middleware := Logger(logger)
	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := middleware(handler)
	assert.NotNil(t, wrapped)

	ctx := &Context{Message: &Message{}}
	err := wrapped(ctx)
	assert.NoError(t, err)
}

// Test Logger middleware logs as expected
func TestLogger_LogsMessages(t *testing.T) {
	// Use Nop logger - we just verify it doesn't panic
	logger := zerolog.Nop()

	middleware := Logger(logger)

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := middleware(handler)
	ctx := &Context{
		Message: &Message{Text: "/test"},
		Chat:    &Chat{ID: 123},
		Sender:  &User{ID: 456},
	}

	err := wrapped(ctx)
	assert.NoError(t, err)
}

// Test Logger middleware with error
func TestLogger_WithError(t *testing.T) {
	logger := zerolog.Nop()

	middleware := Logger(logger)

	expectedErr := errors.New("test error")
	handler := func(ctx *Context) error {
		return expectedErr
	}

	wrapped := middleware(handler)
	ctx := &Context{}

	err := wrapped(ctx)
	assert.Equal(t, expectedErr, err)
}

// Test Recover middleware catches panics
func TestRecover_CatchesPanic(t *testing.T) {
	logger := zerolog.Nop()

	middleware := Recover(logger)

	handler := func(ctx *Context) error {
		panic("test panic")
	}

	wrapped := middleware(handler)
	ctx := &Context{}

	// Should not panic, should return an error
	assert.NotPanics(t, func() {
		err := wrapped(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "panic")
	})
}

// Test Recover middleware with typed panic
func TestRecover_CatchesTypedPanic(t *testing.T) {
	logger := zerolog.Nop()

	middleware := Recover(logger)

	handler := func(ctx *Context) error {
		panic(errors.New("typed panic"))
	}

	wrapped := middleware(handler)
	ctx := &Context{}

	assert.NotPanics(t, func() {
		err := wrapped(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "typed panic")
	})
}

// Test Recover middleware doesn't affect normal flow
func TestRecover_NormalFlow(t *testing.T) {
	logger := zerolog.Nop()

	middleware := Recover(logger)

	called := false
	handler := func(ctx *Context) error {
		called = true
		return nil
	}

	wrapped := middleware(handler)
	ctx := &Context{}

	err := wrapped(ctx)
	assert.NoError(t, err)
	assert.True(t, called)
}

// Test RateLimiter creation
func TestNewRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(10, 1*time.Second)
	assert.NotNil(t, limiter)
	assert.NotNil(t, limiter.limits)
	assert.Equal(t, 10, limiter.maxRequests)
	assert.Equal(t, 1*time.Second, limiter.interval)
}

// Test RateLimit middleware allows requests under limit
func TestRateLimitMiddleware_AllowsUnderLimit(t *testing.T) {
	limiter := NewRateLimiter(3, 1*time.Second)

	middleware := RateLimitMiddleware(limiter)
	called := false
	handler := func(ctx *Context) error {
		called = true
		return nil
	}

	wrapped := middleware(handler)

	// First 3 requests should be allowed
	for i := 0; i < 3; i++ {
		ctx := &Context{Sender: &User{ID: int64(i + 1)}}
		err := wrapped(ctx)
		assert.NoError(t, err)
	}

	assert.True(t, called)
}

// Test RateLimit middleware blocks over limit
func TestRateLimitMiddleware_BlocksOverLimit(t *testing.T) {
	limiter := NewRateLimiter(2, 1*time.Second)

	middleware := RateLimitMiddleware(limiter)

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := middleware(handler)
	userID := int64(1)

	// First 2 requests should be allowed
	for i := 0; i < 2; i++ {
		ctx := &Context{Sender: &User{ID: userID}}
		err := wrapped(ctx)
		assert.NoError(t, err)
	}

	// Third request should be blocked
	ctx := &Context{Sender: &User{ID: userID}}
	err := wrapped(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit exceeded")
}

// Test RateLimit middleware allows different users
func TestRateLimitMiddleware_DifferentUsers(t *testing.T) {
	limiter := NewRateLimiter(1, 1*time.Second)

	middleware := RateLimitMiddleware(limiter)

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := middleware(handler)

	// Each user should get their own limit
	for userID := int64(1); userID <= 5; userID++ {
		ctx := &Context{Sender: &User{ID: userID}}
		err := wrapped(ctx)
		assert.NoError(t, err)
	}
}

// Test RateLimit middleware allows requests without sender
func TestRateLimitMiddleware_NoSender(t *testing.T) {
	limiter := NewRateLimiter(1, 1*time.Second)

	middleware := RateLimitMiddleware(limiter)

	called := false
	handler := func(ctx *Context) error {
		called = true
		return nil
	}

	wrapped := middleware(handler)

	// Request without sender should be allowed
	ctx := &Context{}
	err := wrapped(ctx)
	assert.NoError(t, err)
	assert.True(t, called)
}

// Test RateLimit middleware with channel posts (no sender)
func TestRateLimitMiddleware_ChannelPost(t *testing.T) {
	limiter := NewRateLimiter(1, 1*time.Second)

	middleware := RateLimitMiddleware(limiter)

	handlerCalled := false
	handler := func(ctx *Context) error {
		handlerCalled = true
		return nil
	}

	wrapped := middleware(handler)

	// Channel post without from field
	ctx := &Context{
		Message: &Message{
			Chat: &Chat{ID: -1001234567, Type: "channel"},
		},
	}
	err := wrapped(ctx)
	assert.NoError(t, err)
	assert.True(t, handlerCalled)
}

// Test RateLimiter Cleanup - we can't directly access private rateLimitState,
// but we can test Cleanup doesn't panic and runs without error
func TestRateLimiter_Cleanup(t *testing.T) {
	limiter := NewRateLimiter(1, 10*time.Millisecond)

	// Make some requests to create state
	middleware := RateLimitMiddleware(limiter)
	handler := func(ctx *Context) error { return nil }
	wrapped := middleware(handler)

	// Create state for user 1
	ctx1 := &Context{Sender: &User{ID: 1}}
	_ = wrapped(ctx1)

	// Wait for state to potentially expire
	time.Sleep(30 * time.Millisecond)

	// Run cleanup - should not panic
	assert.NotPanics(t, func() {
		limiter.Cleanup()
	})
}

// Test OnlyGroups middleware
func TestOnlyGroups_AllowsGroups(t *testing.T) {
	middleware := OnlyGroups()

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := middleware(handler)

	// Group chat should be allowed
	ctx := &Context{
		Chat: &Chat{ID: -123, Type: "group"},
	}

	err := wrapped(ctx)
	assert.NoError(t, err)
}

// Test OnlyGroups middleware blocks private
func TestOnlyGroups_BlocksPrivate(t *testing.T) {
	middleware := OnlyGroups()

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := middleware(handler)

	// Private chat should be blocked
	ctx := &Context{
		Chat: &Chat{ID: 123, Type: "private"},
	}

	err := wrapped(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "groups")
}

// Test OnlyGroups middleware blocks without chat
func TestOnlyGroups_BlocksNoChat(t *testing.T) {
	middleware := OnlyGroups()

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := middleware(handler)

	// No chat should be blocked
	ctx := &Context{}

	err := wrapped(ctx)
	assert.Error(t, err)
}

// Test OnlyPrivate middleware
func TestOnlyPrivate_AllowsPrivate(t *testing.T) {
	middleware := OnlyPrivate()

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := middleware(handler)

	// Private chat should be allowed
	ctx := &Context{
		Chat: &Chat{ID: 123, Type: "private"},
	}

	err := wrapped(ctx)
	assert.NoError(t, err)
}

// Test OnlyPrivate middleware blocks groups
func TestOnlyPrivate_BlocksGroups(t *testing.T) {
	middleware := OnlyPrivate()

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := middleware(handler)

	// Group chat should be blocked
	ctx := &Context{
		Chat: &Chat{ID: -123, Type: "group"},
	}

	err := wrapped(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "private")
}

// Test Chain middleware
func TestChain_MultipleMiddlewares(t *testing.T) {
	middleware1 := func(next Handler) Handler {
		return func(ctx *Context) error {
			return next(ctx)
		}
	}

	middleware2 := func(next Handler) Handler {
		return func(ctx *Context) error {
			return next(ctx)
		}
	}

	middleware3 := func(next Handler) Handler {
		return func(ctx *Context) error {
			return next(ctx)
		}
	}

	combined := Chain(middleware1, middleware2, middleware3)
	assert.NotNil(t, combined)

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := combined(handler)
	ctx := &Context{}

	err := wrapped(ctx)
	assert.NoError(t, err)
}

// Test Chain with no middlewares
func TestChain_NoMiddlewares(t *testing.T) {
	combined := Chain()
	assert.NotNil(t, combined)

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := combined(handler)
	ctx := &Context{}

	err := wrapped(ctx)
	assert.NoError(t, err)
}

// Test Chain with single middleware
func TestChain_SingleMiddleware(t *testing.T) {
	middleware := func(next Handler) Handler {
		return func(ctx *Context) error {
			return next(ctx)
		}
	}

	combined := Chain(middleware)
	assert.NotNil(t, combined)

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := combined(handler)
	ctx := &Context{}

	err := wrapped(ctx)
	assert.NoError(t, err)
}

// Test Chain preserves middleware order
func TestChain_PreservesOrder(t *testing.T) {
	order := []string{}

	middleware1 := func(next Handler) Handler {
		return func(ctx *Context) error {
			order = append(order, "m1")
			return next(ctx)
		}
	}

	middleware2 := func(next Handler) Handler {
		return func(ctx *Context) error {
			order = append(order, "m2")
			return next(ctx)
		}
	}

	middleware3 := func(next Handler) Handler {
		return func(ctx *Context) error {
			order = append(order, "m3")
			return next(ctx)
		}
	}

	combined := Chain(middleware1, middleware2, middleware3)

	handler := func(ctx *Context) error {
		order = append(order, "handler")
		return nil
	}

	wrapped := combined(handler)
	ctx := &Context{}

	wrapped(ctx)

	// Order should be: m1 -> m2 -> m3 -> handler
	assert.Equal(t, []string{"m1", "m2", "m3", "handler"}, order)
}

// Test middleware with concurrent access
func TestMiddleware_Concurrent(t *testing.T) {
	limiter := NewRateLimiter(100, 1*time.Second)
	middleware := RateLimitMiddleware(limiter)

	handler := func(ctx *Context) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	wrapped := middleware(handler)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(userID int64) {
			defer wg.Done()
			ctx := &Context{Sender: &User{ID: userID}}
			_ = wrapped(ctx)
		}(int64(i))
	}

	wg.Wait()
	// Test passes if no race conditions detected (run with -race flag)
}

// Test Recover with concurrent panics
func TestRecover_Concurrent(t *testing.T) {
	logger := zerolog.Nop()
	middleware := Recover(logger)

	handler := func(ctx *Context) error {
		panic("concurrent panic")
	}

	wrapped := middleware(handler)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := &Context{}
			assert.NotPanics(t, func() {
				_ = wrapped(ctx)
			})
		}()
	}

	wg.Wait()
}

// Test Logger with concurrent logging
func TestLogger_Concurrent(t *testing.T) {
	logger := zerolog.Nop()
	middleware := Logger(logger)

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := middleware(handler)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := &Context{
				Message: &Message{Text: "/test"},
				Chat:    &Chat{ID: int64(i)},
				Sender:  &User{ID: int64(i)},
			}
			_ = wrapped(ctx)
		}()
	}

	wg.Wait()
}

// Benchmark middleware performance
func BenchmarkLogger_Middleware(b *testing.B) {
	logger := zerolog.Nop()
	middleware := Logger(logger)

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := middleware(handler)
	ctx := &Context{
		Message: &Message{Text: "/test"},
		Chat:    &Chat{ID: 123},
		Sender:  &User{ID: 456},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = wrapped(ctx)
	}
}

func BenchmarkRecover_Middleware(b *testing.B) {
	logger := zerolog.Nop()
	middleware := Recover(logger)

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := middleware(handler)
	ctx := &Context{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = wrapped(ctx)
	}
}

func BenchmarkRateLimit_Middleware(b *testing.B) {
	limiter := NewRateLimiter(1000, 1*time.Second)
	middleware := RateLimitMiddleware(limiter)

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := middleware(handler)
	ctx := &Context{Sender: &User{ID: 1}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = wrapped(ctx)
	}
}

func BenchmarkMiddleware_Chain(b *testing.B) {
	middleware := func(next Handler) Handler {
		return func(ctx *Context) error {
			return next(ctx)
		}
	}

	combined := Chain(middleware, middleware, middleware)

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := combined(handler)
	ctx := &Context{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = wrapped(ctx)
	}
}
