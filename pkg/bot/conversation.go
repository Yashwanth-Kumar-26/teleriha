package bot

import (
	"fmt"
	"sync"
)

// ConversationManager manages multi-step conversations for the bot.
// It tracks the state of each user's conversation.
type ConversationManager struct {
	// conversations maps user IDs to their current conversation
	conversations map[int64]*Conversation

	// handlers maps conversation IDs to their handlers
	handlers map[string]ConversationHandler

	// mu protects the conversations and handlers maps
	mu sync.RWMutex
}

// NewConversationManager creates a new ConversationManager.
func NewConversationManager() *ConversationManager {
	return &ConversationManager{
		conversations: make(map[int64]*Conversation),
		handlers:       make(map[string]ConversationHandler),
	}
}

// Conversation represents a user's current conversation state.
type Conversation struct {
	// ID is the conversation identifier
	ID string

	// UserID is the user's Telegram ID
	UserID int64

	// ChatID is the chat ID where the conversation is happening
	ChatID int64

	// State is the current state of the conversation
	State string

	// Data contains conversation-specific data
	Data map[string]interface{}
}

// ConversationHandler is a function that handles a conversation step.
type ConversationHandler func(*Context, *Conversation) error

// Register registers a conversation handler with the given ID.
func (cm *ConversationManager) Register(id string, handler ConversationHandler) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.handlers[id] = handler
}

// Start starts a new conversation for the given user.
func (cm *ConversationManager) Start(userID, chatID int64, id string) *Conversation {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conv := &Conversation{
		ID:     id,
		UserID: userID,
		ChatID: chatID,
		State:  "start",
		Data:   make(map[string]interface{}),
	}

	cm.conversations[userID] = conv
	return conv
}

// Get returns the current conversation for the given user.
func (cm *ConversationManager) Get(userID int64) *Conversation {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.conversations[userID]
}

// End ends the conversation for the given user.
func (cm *ConversationManager) End(userID int64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.conversations, userID)
}

// UpdateState updates the state of a conversation.
func (cm *ConversationManager) UpdateState(userID int64, state string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if conv, ok := cm.conversations[userID]; ok {
		conv.State = state
	}
}

// SetData sets a key-value pair in the conversation data.
func (cm *ConversationManager) SetData(userID int64, key string, value interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if conv, ok := cm.conversations[userID]; ok {
		conv.Data[key] = value
	}
}

// GetData returns a value from the conversation data.
func (cm *ConversationManager) GetData(userID int64, key string) interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if conv, ok := cm.conversations[userID]; ok {
		return conv.Data[key]
	}
	return nil
}

// ClearData clears all data from a conversation.
func (cm *ConversationManager) ClearData(userID int64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if conv, ok := cm.conversations[userID]; ok {
		conv.Data = make(map[string]interface{})
	}
}

// Handle processes a message in the context of a conversation.
func (cm *ConversationManager) Handle(ctx *Context) (bool, error) {
	userID := ctx.SenderID()
	if userID == 0 {
		return false, nil
	}

	conv := cm.Get(userID)
	if conv == nil {
		// No active conversation
		return false, nil
	}

	handler, ok := cm.handlers[conv.ID]
	if !ok {
		// No handler for this conversation
		cm.End(userID)
		return false, nil
	}

	// Call the handler
	err := handler(ctx, conv)

	// If the handler wants to end the conversation, it should do so explicitly
	// by calling cm.End(userID)

	return true, err
}

// ConversationMiddleware creates a middleware that checks for active conversations.
func (cm *ConversationManager) ConversationMiddleware() Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) error {
			handled, err := cm.Handle(ctx)
			if err != nil {
				return err
			}
			if handled {
				return nil
			}
			return next(ctx)
		}
	}
}

// Simple Conversation Builder for common use cases

// ConversationBuilder provides a fluent interface for building conversations.
type ConversationBuilder struct {
	cm          *ConversationManager
	id          string
	steps       map[string]ConversationHandler
	startHandler ConversationHandler
}

// NewConversationBuilder creates a new ConversationBuilder.
func NewConversationBuilder(cm *ConversationManager, id string) *ConversationBuilder {
	return &ConversationBuilder{
		cm:    cm,
		id:    id,
		steps: make(map[string]ConversationHandler),
	}
}

// Start sets the starting handler for the conversation.
func (cb *ConversationBuilder) Start(handler ConversationHandler) *ConversationBuilder {
	cb.startHandler = handler
	return cb
}

// Step adds a step with the given state and handler.
func (cb *ConversationBuilder) Step(state string, handler ConversationHandler) *ConversationBuilder {
	cb.steps[state] = handler
	return cb
}

// Next creates a handler that transitions to the next state.
func (cb *ConversationBuilder) Next(nextState string) ConversationHandler {
	return func(ctx *Context, conv *Conversation) error {
		cb.cm.UpdateState(conv.UserID, nextState)
		return nil
	}
}

// End creates a handler that ends the conversation.
func (cb *ConversationBuilder) End() ConversationHandler {
	return func(ctx *Context, conv *Conversation) error {
		cb.cm.End(conv.UserID)
		return nil
	}
}

// WaitForText returns a handler that waits for a text message and stores it.
func (cb *ConversationBuilder) WaitForText(key string, nextState string) ConversationHandler {
	return func(ctx *Context, conv *Conversation) error {
		text := ctx.Text()
		if text == "" {
			return fmt.Errorf("please send a text message")
		}

		cb.cm.SetData(conv.UserID, key, text)
		cb.cm.UpdateState(conv.UserID, nextState)
		return nil
	}
}

// WaitForCallback returns a handler that waits for a callback query.
func (cb *ConversationBuilder) WaitForCallback(expectedData string, nextState string) ConversationHandler {
	return func(ctx *Context, conv *Conversation) error {
		if !ctx.IsCallbackQuery() {
			return fmt.Errorf("please use the button")
		}

		callbackData := ctx.CallbackData()
		if callbackData != expectedData {
			return fmt.Errorf("unexpected button")
		}

		cb.cm.UpdateState(conv.UserID, nextState)

		// Answer the callback
		if err := ctx.AnswerCallback("", false); err != nil {
			return fmt.Errorf("failed to answer callback: %w", err)
		}

		// Delete the original message
		if err := ctx.Delete(); err != nil {
			// Ignore error
			_ = err
		}

		return nil
	}
}

// StoreValue returns a handler that stores a value and transitions.
func (cb *ConversationBuilder) StoreValue(key string, value interface{}, nextState string) ConversationHandler {
	return func(ctx *Context, conv *Conversation) error {
		cb.cm.SetData(conv.UserID, key, value)
		cb.cm.UpdateState(conv.UserID, nextState)
		return nil
	}
}

// Build registers the conversation with the manager.
func (cb *ConversationBuilder) Build() {
	cb.cm.Register(cb.id, func(ctx *Context, conv *Conversation) error {
		handler, ok := cb.steps[conv.State]
		if !ok {
			if conv.State == "start" && cb.startHandler != nil {
				return cb.startHandler(ctx, conv)
			}
			// End conversation if no handler found
			cb.cm.End(conv.UserID)
			return fmt.Errorf("conversation ended")
		}
		return handler(ctx, conv)
	})
}
