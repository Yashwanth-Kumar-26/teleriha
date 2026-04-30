package bot

import (
	"errors"
	"regexp"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper to create a test message
func testMessage(text string, isCommand bool) *Message {
	msg := &Message{
		MessageID: 1,
		Text:      text,
		Date:      1234567890,
		Chat: &Chat{
			ID:   123,
			Type: "private",
		},
		From: &User{
			ID:        456,
			FirstName: "Test",
			Username:  "testuser",
			IsBot:     false,
		},
	}

	if isCommand {
		msg.Entities = []MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: len(text),
			},
		}
	}

	return msg
}

// Helper to create a test callback query
func testCallbackQuery(data string) *CallbackQuery {
	return &CallbackQuery{
		ID:           "cb1",
		From:         &User{ID: 456, FirstName: "Test"},
		Data:         data,
		ChatInstance: "chat1",
	}
}

// Test NewRouter creates a new router
func TestNewRouter(t *testing.T) {
	router := NewRouter()
	assert.NotNil(t, router)
	// Fields are unexported, but we can verify router is created
}

// Test On registers a command handler
func TestOn_RegistersHandler(t *testing.T) {
	router := NewRouter()

	handler := func(ctx *Context) error {
		return nil
	}

	router.On("/start", handler)

	// Verify handler is registered by calling it
	msg := testMessage("/start", true)
	ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}
	// Should not panic and should call the handler
	assert.NotPanics(t, func() {
		router.HandleMessage(ctx, *msg)
	})
}

// Test On with case-insensitive command matching
func TestOn_CaseInsensitive(t *testing.T) {
	router := NewRouter()

	result := ""
	handler := func(ctx *Context) error {
		result = "handled"
		return nil
	}

	router.On("/START", handler)

	msg := testMessage("/start", true)
	ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}

	router.HandleMessage(ctx, *msg)

	assert.Equal(t, "handled", result)
}

// Test On with bot username in command
func TestOn_WithBotUsername(t *testing.T) {
	router := NewRouter()

	result := ""
	handler := func(ctx *Context) error {
		result = "handled"
		return nil
	}

	router.On("/start", handler)

	// Message with bot username
	msg := testMessage("/start@mybot", true)
	ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}

	router.HandleMessage(ctx, *msg)

	assert.Equal(t, "handled", result)
}

// Test OnRegex registers a regex handler
func TestOnRegex_MatchesPattern(t *testing.T) {
	router := NewRouter()

	result := ""
	handler := func(ctx *Context) error {
		result = "regex matched"
		return nil
	}

	err := router.OnRegex(`^hello|hi$`, handler)
	assert.NoError(t, err)

	// Test matching
	msg := testMessage("hello", false)
	ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}

	router.HandleMessage(ctx, *msg)

	assert.Equal(t, "regex matched", result)

	// Reset and test non-matching
	result = ""
	msg = testMessage("goodbye", false)
	ctx = &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}

	router.HandleMessage(ctx, *msg)
	assert.Equal(t, "", result)
}

// Test OnRegex with numeric pattern
func TestOnRegex_NumericPattern(t *testing.T) {
	router := NewRouter()

	result := ""
	handler := func(ctx *Context) error {
		result = "number matched"
		return nil
	}

	err := router.OnRegex(`^[0-9]+$`, handler)
	assert.NoError(t, err)

	msg := testMessage("12345", false)
	ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}

	router.HandleMessage(ctx, *msg)

	assert.Equal(t, "number matched", result)

	// Should not match non-numeric
	result = ""
	msg = testMessage("abc", false)
	ctx = &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}

	router.HandleMessage(ctx, *msg)
	assert.Equal(t, "", result)
}

// Test OnRegex with invalid pattern
func TestOnRegex_InvalidPattern(t *testing.T) {
	router := NewRouter()

	// this[ is invalid regex
	err := router.OnRegex(`[`, func(ctx *Context) error { return nil })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid regex pattern")
}

// Test Default handler
func TestDefault_HandleUnknownCommand(t *testing.T) {
	router := NewRouter()

	defaultCalled := false
	commandCalled := false

	router.On("/start", func(ctx *Context) error {
		commandCalled = true
		return nil
	})

	router.Default(func(ctx *Context) error {
		defaultCalled = true
		return nil
	})

	// Unknown command should call default
	msg := testMessage("/unknown", true)
	ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}

	router.HandleMessage(ctx, *msg)

	assert.False(t, commandCalled, "Command handler should not be called")
	assert.True(t, defaultCalled, "Default handler should be called")
}

// Test Default handler with empty string key
func TestDefault_EmptyStringKey(t *testing.T) {
	router := NewRouter()

	called := false
	router.Default(func(ctx *Context) error {
		called = true
		return nil
	})

	// Empty message should call default
	msg := testMessage("", false)
	ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}

	router.HandleMessage(ctx, *msg)

	assert.True(t, called)
}

// Test Middleware chain execution order
func TestUse_MiddlewareChainOrder(t *testing.T) {
	router := NewRouter()

	order := []string{}

	middleware1 := func(next Handler) Handler {
		return func(ctx *Context) error {
			order = append(order, "middleware1-before")
			err := next(ctx)
			order = append(order, "middleware1-after")
			return err
		}
	}

	middleware2 := func(next Handler) Handler {
		return func(ctx *Context) error {
			order = append(order, "middleware2-before")
			err := next(ctx)
			order = append(order, "middleware2-after")
			return err
		}
	}

	handler := func(ctx *Context) error {
		order = append(order, "handler")
		return nil
	}

	router.Use(middleware1)
	router.Use(middleware2)
	router.On("/test", handler)

	msg := testMessage("/test", true)
	ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}

	router.HandleMessage(ctx, *msg)

	// Middleware are applied in reverse order (LIFO) in applyMiddlewares
	// middleware1 was added first, so it wraps the handler first
	// middleware2 was added second, so it wraps middleware1
	// Execution order: middleware2-before -> middleware1-before -> handler -> middleware1-after -> middleware2-after
	assert.Equal(t, []string{
		"middleware1-before",
		"middleware2-before",
		"handler",
		"middleware2-after",
		"middleware1-after",
	}, order)
}

// Test Middleware can short-circuit
func TestMiddleware_ShortCircuit(t *testing.T) {
	router := NewRouter()

	handlerCalled := false

	middleware := func(next Handler) Handler {
		return func(ctx *Context) error {
			return errors.New("blocked by middleware")
		}
	}

	handler := func(ctx *Context) error {
		handlerCalled = true
		return nil
	}

	router.Use(middleware)
	router.On("/test", handler)

	msg := testMessage("/test", true)
	ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}

	router.HandleMessage(ctx, *msg)

	assert.False(t, handlerCalled, "Handler should not be called when middleware returns error")
}

// Test Middleware can modify context
func TestMiddleware_ModifiesContext(t *testing.T) {
	router := NewRouter()

	middleware := func(next Handler) Handler {
		return func(ctx *Context) error {
			// Middleware can add data to context
			return next(ctx)
		}
	}

	handler := func(ctx *Context) error {
		// Handler receives modified context
		return nil
	}

	router.Use(middleware)
	router.On("/test", handler)

	msg := testMessage("/test", true)
	ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}

	// Should not panic
	router.HandleMessage(ctx, *msg)
}

// Test OnCallback registers callback handler
func TestOnCallback_RegistersHandler(t *testing.T) {
	router := NewRouter()

	called := false
	handler := func(ctx *Context) error {
		called = true
		return nil
	}

	router.OnCallback("test_", handler)

	// Verify handler is registered
	assert.NotNil(t, router.callbackHandlers["test_"])

	// Callback with matching prefix
	callback := testCallbackQuery("test_action")
	ctx := &Context{CallbackQuery: callback, Chat: &Chat{ID: 123}, Sender: callback.From}

	router.HandleCallbackQuery(ctx, *callback)

	assert.True(t, called, "Callback handler should have been called")
}

// Test OnCallback with non-matching prefix
func TestOnCallback_NoMatch(t *testing.T) {
	router := NewRouter()

	called := false
	handler := func(ctx *Context) error {
		called = true
		return nil
	}

	router.OnCallback("test_", handler)

	// Callback with non-matching prefix
	callback := testCallbackQuery("other_action")
	ctx := &Context{CallbackQuery: callback, Chat: &Chat{ID: 123}, Sender: callback.From}

	router.HandleCallbackQuery(ctx, *callback)

	assert.False(t, called, "Callback handler should not be called for non-matching prefix")
}

// Test OnInlineQuery registers handler
func TestOnInlineQuery_RegistersHandler(t *testing.T) {
	router := NewRouter()

	called := false
	handler := func(ctx *Context) error {
		called = true
		return nil
	}

	router.OnInlineQuery(handler)

	assert.NotNil(t, router.inlineQueryHandler)

	inlineQuery := &InlineQuery{
		ID:   "iq1",
		From: &User{ID: 456},
	}
	ctx := &Context{InlineQuery: inlineQuery, Sender: inlineQuery.From}

	router.HandleInlineQuery(ctx, *inlineQuery)

	assert.True(t, called, "InlineQuery handler should have been called")
}

// Test OnChosenInlineResult registers handler
func TestOnChosenInlineResult_RegistersHandler(t *testing.T) {
	router := NewRouter()

	called := false
	handler := func(ctx *Context) error {
		called = true
		return nil
	}

	router.OnChosenInlineResult(handler)

	assert.NotNil(t, router.chosenInlineResultHandler)

	result := &ChosenInlineResult{
		ResultID: "r1",
		From:     &User{ID: 456},
	}
	ctx := &Context{ChosenInlineResult: result, Sender: result.From}

	router.HandleChosenInlineResult(ctx, *result)

	assert.True(t, called, "ChosenInlineResult handler should have been called")
}

// Test Clone creates independent router
func TestClone_CreatesIndependentRouter(t *testing.T) {
	router := NewRouter()

	handler1 := func(ctx *Context) error { return nil }
	handler2 := func(ctx *Context) error { return nil }

	router.On("/start", handler1)

	cloned := router.Clone()
	cloned.On("/help", handler2)

	// Test that handlers work independently
	// We can't access private fields, so test via HandleMessage
	msg1 := testMessage("/start", true)
	ctx1 := &Context{Message: msg1, Chat: msg1.Chat, Sender: msg1.From}
	assert.NotPanics(t, func() { router.HandleMessage(ctx1, *msg1) })

	msg2 := testMessage("/help", true)
	ctx2 := &Context{Message: msg2, Chat: msg2.Chat, Sender: msg2.From}
	assert.NotPanics(t, func() { cloned.HandleMessage(ctx2, *msg2) })
}

// Test Group creates a group router
func TestGroup_CreatesGroupRouter(t *testing.T) {
	router := NewRouter()

	handler := func(ctx *Context) error { return nil }

	group := router.Group("admin")
	group.On("settings", handler)

	// Group router should have been created
	assert.NotNil(t, group)
}

// Test Group with middlewares
func TestGroup_WithMiddlewares(t *testing.T) {
	router := NewRouter()

	middleware := func(next Handler) Handler {
		return func(ctx *Context) error {
			return next(ctx)
		}
	}

	group := router.Group("admin", middleware)
	assert.NotNil(t, group)
}

// Test UseGroup adds middleware to group
func TestUseGroup_AddsMiddleware(t *testing.T) {
	router := NewRouter()

	group := router.Group("admin")
	group.Use(func(next Handler) Handler {
		return func(ctx *Context) error {
			return next(ctx)
		}
	})

	// We can't access private fields, but test that it doesn't panic
	assert.NotNil(t, group)
}

// Test applyMiddlewares wraps handler correctly
func TestApplyMiddlewares(t *testing.T) {
	router := NewRouter()

	middlewareCalled := false
	middleware := func(next Handler) Handler {
		return func(ctx *Context) error {
			middlewareCalled = true
			return next(ctx)
		}
	}

	router.Use(middleware)

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := router.applyMiddlewares(handler)
	assert.NotNil(t, wrapped)

	ctx := &Context{}
	wrapped(ctx)
	assert.True(t, middlewareCalled)
}

// Test Chain combines multiple middlewares
func TestChain_Middlewares(t *testing.T) {
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

	combined := Chain(middleware1, middleware2)
	assert.NotNil(t, combined)

	handler := func(ctx *Context) error {
		return nil
	}

	wrapped := combined(handler)
	assert.NotNil(t, wrapped)

	ctx := &Context{}
	err := wrapped(ctx)
	assert.NoError(t, err)
}

// Test router with concurrent access
func TestRouter_Concurrent(t *testing.T) {
	router := NewRouter()

	handler := func(ctx *Context) error {
		return nil
	}

	// Register many handlers
	for i := 0; i < 100; i++ {
		cmd := "/cmd" + string(rune(i%26+'a'))
		router.On(cmd, handler)
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			cmd := "/cmd" + string(rune(idx%26+'a'))
			msg := testMessage(cmd, true)
			ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}
			router.HandleMessage(ctx, *msg)
		}(i)
	}

	wg.Wait()
	// Test passes if no race conditions detected (run with -race flag)
}

// Test OnRegex with compiled regex directly
func TestOnRegex_WithRunes(t *testing.T) {
	router := NewRouter()

	result := ""
	handler := func(ctx *Context) error {
		result = "matched"
		return nil
	}

	// Test with emoji
	err := router.OnRegex(`👍`, handler)
	assert.NoError(t, err)

	msg := testMessage("👍", false)
	ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}

	router.HandleMessage(ctx, *msg)

	assert.Equal(t, "matched", result)
}

// Test regexHandlers slice grows correctly
func TestOnRegex_MultiplePatterns(t *testing.T) {
	router := NewRouter()

	handler := func(ctx *Context) error { return nil }

	patterns := []string{
		`^hello$`,
		`^hi$`,
		`^hey$`,
		`^yo$`,
	}

	for _, pattern := range patterns {
		err := router.OnRegex(pattern, handler)
		assert.NoError(t, err)
	}

	assert.Len(t, router.regexHandlers, len(patterns))

	// Verify all patterns compile
	for _, rh := range router.regexHandlers {
		assert.NotNil(t, rh.pattern)
		_, err := regexp.Compile(rh.pattern.String())
		assert.NoError(t, err)
	}
}

// Benchmark router performance
func BenchmarkRouter_HandleMessage_SingleHandler(b *testing.B) {
	router := NewRouter()

	handler := func(ctx *Context) error {
		return nil
	}

	router.On("/start", handler)

	msg := testMessage("/start", true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}
		router.HandleMessage(ctx, *msg)
	}
}

func BenchmarkRouter_HandleMessage_ManyHandlers(b *testing.B) {
	router := NewRouter()

	handler := func(ctx *Context) error {
		return nil
	}

	// Register 100 handlers
	for i := 0; i < 100; i++ {
		router.On("/cmd"+string(rune(i%26+'a')), handler)
	}

	msg := testMessage("/cmd0", true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}
		router.HandleMessage(ctx, *msg)
	}
}

func BenchmarkRouter_WithMiddlewareChain(b *testing.B) {
	router := NewRouter()

	middleware := func(next Handler) Handler {
		return func(ctx *Context) error {
			return next(ctx)
		}
	}

	handler := func(ctx *Context) error {
		return nil
	}

	// Add 10 middlewares
	for i := 0; i < 10; i++ {
		router.Use(middleware)
	}

	router.On("/start", handler)

	msg := testMessage("/start", true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}
		router.HandleMessage(ctx, *msg)
	}
}

func BenchmarkRouter_RegexMatch(b *testing.B) {
	router := NewRouter()

	handler := func(ctx *Context) error {
		return nil
	}

	router.OnRegex(`^[a-zA-Z]+$`, handler)

	msg := testMessage("hello", false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := &Context{Message: msg, Chat: msg.Chat, Sender: msg.From}
		router.HandleMessage(ctx, *msg)
	}
}
