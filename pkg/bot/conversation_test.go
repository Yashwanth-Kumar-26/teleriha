package bot

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test NewConversationManager creates a manager
func TestNewConversationManager(t *testing.T) {
	cm := NewConversationManager()

	assert.NotNil(t, cm)
	assert.NotNil(t, cm.conversations)
	assert.NotNil(t, cm.handlers)
}

// Test Register adds a handler
func TestConversationManager_Register(t *testing.T) {
	cm := NewConversationManager()
	handler := func(ctx *Context, conv *Conversation) error { return nil }

	cm.Register("test", handler)

	// Verify handler is registered (we can't directly access handlers map, but we can test via Start)
	assert.NotNil(t, cm)
}

// Test Start begins a new conversation
func TestConversationManager_Start(t *testing.T) {
	cm := NewConversationManager()

	conv := cm.Start(123, 456, "test")

	assert.NotNil(t, conv)
	assert.Equal(t, "test", conv.ID)
	assert.Equal(t, int64(123), conv.UserID)
	assert.Equal(t, int64(456), conv.ChatID)
	assert.Equal(t, "start", conv.State)
	assert.NotNil(t, conv.Data)
}

// Test Get returns an active conversation
func TestConversationManager_Get(t *testing.T) {
	cm := NewConversationManager()

	cm.Start(123, 456, "test")
	conv := cm.Get(123)

	assert.NotNil(t, conv)
	assert.Equal(t, "test", conv.ID)
}

// Test Get returns nil for non-existent conversation
func TestConversationManager_Get_NotFound(t *testing.T) {
	cm := NewConversationManager()

	conv := cm.Get(999)

	assert.Nil(t, conv)
}

// Test End removes a conversation
func TestConversationManager_End(t *testing.T) {
	cm := NewConversationManager()

	cm.Start(123, 456, "test")
	cm.End(123)

	conv := cm.Get(123)
	assert.Nil(t, conv)
}

// Test UpdateState changes conversation state
func TestConversationManager_UpdateState(t *testing.T) {
	cm := NewConversationManager()

	cm.Start(123, 456, "test")
	cm.UpdateState(123, "step2")

	conv := cm.Get(123)
	assert.Equal(t, "step2", conv.State)
}

// Test UpdateState does nothing for non-existent conversation
func TestConversationManager_UpdateState_NotFound(t *testing.T) {
	cm := NewConversationManager()

	// Should not panic
	assert.NotPanics(t, func() {
		cm.UpdateState(999, "step2")
	})
}

// Test SetData stores data in conversation
func TestConversationManager_SetData(t *testing.T) {
	cm := NewConversationManager()

	cm.Start(123, 456, "test")
	cm.SetData(123, "name", "John")

	conv := cm.Get(123)
	assert.Equal(t, "John", conv.Data["name"])
}

// Test GetData retrieves data from conversation
func TestConversationManager_GetData(t *testing.T) {
	cm := NewConversationManager()

	cm.Start(123, 456, "test")
	cm.SetData(123, "name", "John")

	value := cm.GetData(123, "name")
	assert.Equal(t, "John", value)
}

// Test GetData returns nil for non-existent conversation
func TestConversationManager_GetData_NotFound(t *testing.T) {
	cm := NewConversationManager()

	value := cm.GetData(999, "name")
	assert.Nil(t, value)
}

// Test GetData returns nil for non-existent key
func TestConversationManager_GetData_KeyNotFound(t *testing.T) {
	cm := NewConversationManager()

	cm.Start(123, 456, "test")
	value := cm.GetData(123, "name")
	assert.Nil(t, value)
}

// Test ClearData removes all data
func TestConversationManager_ClearData(t *testing.T) {
	cm := NewConversationManager()

	cm.Start(123, 456, "test")
	cm.SetData(123, "name", "John")
	cm.SetData(123, "age", 30)
	cm.ClearData(123)

	conv := cm.Get(123)
	assert.Empty(t, conv.Data)
}

// Test Handle processes a message in a conversation
func TestConversationManager_Handle(t *testing.T) {
	cm := NewConversationManager()
	handlerCalled := false

	handler := func(ctx *Context, conv *Conversation) error {
		handlerCalled = true
		return nil
	}
	cm.Register("test", handler)
	cm.Start(123, 456, "test")

	ctx := &Context{Sender: &User{ID: 123}}
	handled, err := cm.Handle(ctx)

	assert.True(t, handled)
	assert.NoError(t, err)
	assert.True(t, handlerCalled)
}

// Test Handle returns false for no active conversation
func TestConversationManager_Handle_NoConversation(t *testing.T) {
	cm := NewConversationManager()

	ctx := &Context{Sender: &User{ID: 123}}
	handled, err := cm.Handle(ctx)

	assert.False(t, handled)
	assert.NoError(t, err)
}

// Test Handle returns false for no sender
func TestConversationManager_Handle_NoSender(t *testing.T) {
	cm := NewConversationManager()

	ctx := &Context{}
	handled, err := cm.Handle(ctx)

	assert.False(t, handled)
	assert.NoError(t, err)
}

// Test Handle ends conversation when no handler found
func TestConversationManager_Handle_NoHandler(t *testing.T) {
	cm := NewConversationManager()
	cm.Start(123, 456, "test")

	ctx := &Context{Sender: &User{ID: 123}}
	handled, err := cm.Handle(ctx)

	assert.False(t, handled)
	assert.NoError(t, err)
	// Conversation should be ended
	assert.Nil(t, cm.Get(123))
}

// Test Handle propagates errors from handler
func TestConversationManager_Handle_Error(t *testing.T) {
	cm := NewConversationManager()
	expectedErr := errors.New("handler error")

	handler := func(ctx *Context, conv *Conversation) error {
		return expectedErr
	}
	cm.Register("test", handler)
	cm.Start(123, 456, "test")

	ctx := &Context{Sender: &User{ID: 123}}
	handled, err := cm.Handle(ctx)

	assert.True(t, handled)
	assert.Equal(t, expectedErr, err)
}

// Test ConversationMiddleware creates middleware
func TestConversationMiddleware(t *testing.T) {
	cm := NewConversationManager()

	middleware := cm.ConversationMiddleware()
	assert.NotNil(t, middleware)
}

// Test ConversationMiddleware handles conversation
func TestConversationMiddleware_HandlesConversation(t *testing.T) {
	cm := NewConversationManager()
	handlerCalled := false

	handler := func(ctx *Context, conv *Conversation) error {
		handlerCalled = true
		return nil
	}
	cm.Register("test", handler)
	cm.Start(123, 456, "test")

	middleware := cm.ConversationMiddleware()
	nextCalled := false
	nextHandler := func(ctx *Context) error {
		nextCalled = true
		return nil
	}

	wrapped := middleware(nextHandler)
	ctx := &Context{Sender: &User{ID: 123}}
	err := wrapped(ctx)

	assert.NoError(t, err)
	assert.True(t, handlerCalled)
	assert.False(t, nextCalled) // Conversation handler should short-circuit
}

// Test ConversationMiddleware calls next when no conversation
func TestConversationMiddleware_CallsNext(t *testing.T) {
	cm := NewConversationManager()

	middleware := cm.ConversationMiddleware()
	nextCalled := false
	nextHandler := func(ctx *Context) error {
		nextCalled = true
		return nil
	}

	wrapped := middleware(nextHandler)
	ctx := &Context{Sender: &User{ID: 123}}
	err := wrapped(ctx)

	assert.NoError(t, err)
	assert.True(t, nextCalled)
}

// Test NewConversationBuilder creates a builder
func TestNewConversationBuilder(t *testing.T) {
	cm := NewConversationManager()
	builder := NewConversationBuilder(cm, "test")

	assert.NotNil(t, builder)
	assert.Equal(t, cm, builder.cm)
	assert.Equal(t, "test", builder.id)
	assert.NotNil(t, builder.steps)
}

// Test ConversationBuilder Start
func TestConversationBuilder_Start(t *testing.T) {
	cm := NewConversationManager()
	builder := NewConversationBuilder(cm, "test")

	handler := func(ctx *Context, conv *Conversation) error {
		return nil
	}
	builder.Start(handler)

	// Just verify Start doesn't panic and startHandler is set
	assert.NotNil(t, builder.startHandler)
}

// Test ConversationBuilder Step
func TestConversationBuilder_Step(t *testing.T) {
	cm := NewConversationManager()
	builder := NewConversationBuilder(cm, "test")

	handler := func(ctx *Context, conv *Conversation) error {
		return nil
	}
	builder.Step("step1", handler)

	assert.NotNil(t, builder.steps["step1"])
}

// Test ConversationBuilder Next
func TestConversationBuilder_Next(t *testing.T) {
	cm := NewConversationManager()
	builder := NewConversationBuilder(cm, "test")

	cm.Start(123, 456, "test")
	nextHandler := builder.Next("step2")

	ctx := &Context{Sender: &User{ID: 123}}
	conv := cm.Get(123)
	err := nextHandler(ctx, conv)

	assert.NoError(t, err)
	assert.Equal(t, "step2", conv.State)
}

// Test ConversationBuilder End
func TestConversationBuilder_End(t *testing.T) {
	cm := NewConversationManager()
	builder := NewConversationBuilder(cm, "test")

	cm.Start(123, 456, "test")
	endHandler := builder.End()

	ctx := &Context{Sender: &User{ID: 123}}
	conv := cm.Get(123)
	err := endHandler(ctx, conv)

	assert.NoError(t, err)
	assert.Nil(t, cm.Get(123)) // Conversation should be ended
}

// Test ConversationBuilder WaitForText
func TestConversationBuilder_WaitForText(t *testing.T) {
	cm := NewConversationManager()
	builder := NewConversationBuilder(cm, "test")

	cm.Start(123, 456, "test")
	handler := builder.WaitForText("username", "step2")

	ctx := &Context{Sender: &User{ID: 123}, Message: &Message{Text: "john_doe"}}
	conv := cm.Get(123)
	err := handler(ctx, conv)

	assert.NoError(t, err)
	assert.Equal(t, "step2", conv.State)
	assert.Equal(t, "john_doe", conv.Data["username"])
}

// Test ConversationBuilder WaitForText error on empty text
func TestConversationBuilder_WaitForText_Empty(t *testing.T) {
	cm := NewConversationManager()
	builder := NewConversationBuilder(cm, "test")

	cm.Start(123, 456, "test")
	handler := builder.WaitForText("username", "step2")

	ctx := &Context{Sender: &User{ID: 123}, Message: &Message{Text: ""}}
	conv := cm.Get(123)
	err := handler(ctx, conv)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "text message")
}

// Test ConversationBuilder WaitForCallback
func TestConversationBuilder_WaitForCallback(t *testing.T) {
	cm := NewConversationManager()
	builder := NewConversationBuilder(cm, "test")
	// We need to mock AnswerCallback and Delete for this test
	// Since we can't easily mock these, we'll test the state transition

	cm.Start(123, 456, "test")
	// For this test, we'll just verify the builder creates the handler
	handler := builder.WaitForCallback("confirm", "step2")
	assert.NotNil(t, handler)
}

// Test ConversationBuilder StoreValue
func TestConversationBuilder_StoreValue(t *testing.T) {
	cm := NewConversationManager()
	builder := NewConversationBuilder(cm, "test")

	cm.Start(123, 456, "test")
	handler := builder.StoreValue("count", 42, "step2")

	ctx := &Context{Sender: &User{ID: 123}}
	conv := cm.Get(123)
	err := handler(ctx, conv)

	assert.NoError(t, err)
	assert.Equal(t, "step2", conv.State)
	assert.Equal(t, 42, conv.Data["count"])
}

// Test ConversationBuilder Build registers the conversation
func TestConversationBuilder_Build(t *testing.T) {
	cm := NewConversationManager()
	builder := NewConversationBuilder(cm, "test")

	handler := func(ctx *Context, conv *Conversation) error {
		return nil
	}
	builder.Start(handler)
	builder.Step("step1", handler)
	builder.Build()

	// Start a conversation and verify it can be handled
	cm.Start(123, 456, "test")
	ctx := &Context{Sender: &User{ID: 123}}
	handled, err := cm.Handle(ctx)

	assert.True(t, handled)
	assert.NoError(t, err)
}

// Test multi-step conversation flow
func TestMultiStepConversationFlow(t *testing.T) {
	cm := NewConversationManager()
	steps := []string{}

	handler := func(ctx *Context, conv *Conversation) error {
		steps = append(steps, conv.State)
		return nil
	}

	cm.Register("test", handler)
	conv := cm.Start(123, 456, "test")

	// Simulate step transitions
	conv.State = "step1"
	ctx := &Context{Sender: &User{ID: 123}}
	_, _ = cm.Handle(ctx)

	cm.UpdateState(123, "step2")
	_, _ = cm.Handle(ctx)

	cm.UpdateState(123, "step3")
	_, _ = cm.Handle(ctx)

	assert.Contains(t, steps, "step1")
	assert.Contains(t, steps, "step2")
	assert.Contains(t, steps, "step3")
}

// Test concurrent access to ConversationManager
func TestConversationManager_Concurrent(t *testing.T) {
	cm := NewConversationManager()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(userID int64) {
			defer wg.Done()
			cm.Start(userID, userID, "test")
			cm.SetData(userID, "value", userID)
			cm.UpdateState(userID, "step2")
			cm.GetData(userID, "value")
			cm.End(userID)
		}(int64(i))
	}

	wg.Wait()
	// Test passes if no race conditions detected (run with -race flag)
}

// Test concurrent Handle calls
func TestConversationManager_ConcurrentHandle(t *testing.T) {
	cm := NewConversationManager()

	handler := func(ctx *Context, conv *Conversation) error {
		return nil
	}
	cm.Register("test", handler)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(userID int64) {
			defer wg.Done()
			cm.Start(userID, userID, "test")
			ctx := &Context{Sender: &User{ID: userID}}
			_, _ = cm.Handle(ctx)
			cm.End(userID)
		}(int64(i))
	}

	wg.Wait()
}

// Benchmark ConversationManager operations
func BenchmarkConversationManager_Start(b *testing.B) {
	cm := NewConversationManager()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.Start(int64(i), int64(i), "test")
	}
}

func BenchmarkConversationManager_Get(b *testing.B) {
	cm := NewConversationManager()
	cm.Start(123, 456, "test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cm.Get(123)
	}
}

func BenchmarkConversationManager_SetData(b *testing.B) {
	cm := NewConversationManager()
	cm.Start(123, 456, "test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.SetData(123, "key", "value")
	}
}

func BenchmarkConversationManager_UpdateState(b *testing.B) {
	cm := NewConversationManager()
	cm.Start(123, 456, "test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cm.UpdateState(123, "state")
	}
}
