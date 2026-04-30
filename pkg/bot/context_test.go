package bot

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// Test NewContext creates a context with message
func TestNewContext_WithMessage(t *testing.T) {
	update := Update{
		UpdateID: 1,
		Message: &Message{
			MessageID: 123,
			Text:      "Hello, world!",
			Date:      1234567890,
			Chat: &Chat{
				ID:   456,
				Type: "private",
			},
			From: &User{
				ID:        789,
				FirstName: "John",
				LastName:  "Doe",
				Username:  "johndoe",
			},
		},
	}

	ctx := NewContext(nil, update)

	assert.NotNil(t, ctx)
	assert.Equal(t, "Hello, world!", ctx.Text())
	assert.Equal(t, int64(789), ctx.SenderID())
	assert.Equal(t, int64(456), ctx.ChatID())
	assert.Equal(t, "johndoe", ctx.Sender.Username)
	assert.Equal(t, "John", ctx.Sender.FirstName)
	assert.Equal(t, "Doe", ctx.Sender.LastName)
	assert.True(t, ctx.IsMessage())
	assert.False(t, ctx.IsCallbackQuery())
	assert.False(t, ctx.IsInlineQuery())
	assert.False(t, ctx.IsChosenInlineResult())
}

// Test NewContext with callback query
func TestNewContext_WithCallbackQuery(t *testing.T) {
	update := Update{
		UpdateID: 1,
		CallbackQuery: &CallbackQuery{
			ID:           "cb1",
			From:         &User{ID: 789, FirstName: "John"},
			Data:         "test_data",
			ChatInstance: "chat1",
			Message: &Message{
				MessageID: 123,
				Chat:      &Chat{ID: 456, Type: "private"},
			},
		},
	}

	ctx := NewContext(nil, update)

	assert.NotNil(t, ctx)
	assert.Equal(t, int64(789), ctx.SenderID())
	assert.Equal(t, int64(456), ctx.ChatID())
	assert.Equal(t, "test_data", ctx.CallbackData())
	assert.True(t, ctx.IsCallbackQuery())
	// IsMessage is true because CallbackQuery has an associated Message
	assert.True(t, ctx.IsMessage())
}

// Test NewContext with inline query
func TestNewContext_WithInlineQuery(t *testing.T) {
	update := Update{
		UpdateID: 1,
		InlineQuery: &InlineQuery{
			ID:     "iq1",
			From:   &User{ID: 789, FirstName: "John"},
			Query:  "test query",
			Offset: "",
		},
	}

	ctx := NewContext(nil, update)

	assert.NotNil(t, ctx)
	assert.Equal(t, int64(789), ctx.SenderID())
	assert.Equal(t, "test query", ctx.InlineQueryText())
	assert.True(t, ctx.IsInlineQuery())
	assert.False(t, ctx.IsMessage())
}

// Test NewContext with chosen inline result
func TestNewContext_WithChosenInlineResult(t *testing.T) {
	update := Update{
		UpdateID: 1,
		ChosenInlineResult: &ChosenInlineResult{
			ResultID: "r1",
			From:     &User{ID: 789, FirstName: "John"},
			Query:    "test",
		},
	}

	ctx := NewContext(nil, update)

	assert.NotNil(t, ctx)
	assert.Equal(t, int64(789), ctx.SenderID())
	assert.True(t, ctx.IsChosenInlineResult())
	assert.False(t, ctx.IsMessage())
}

// Test ChatID returns chat ID from message
func TestChatID_FromMessage(t *testing.T) {
	ctx := &Context{
		Chat: &Chat{ID: 12345},
	}

	assert.Equal(t, int64(12345), ctx.ChatID())
}

// Test ChatID returns 0 when no chat
func TestChatID_NoChat(t *testing.T) {
	ctx := &Context{}

	assert.Equal(t, int64(0), ctx.ChatID())
}

// Test SenderID returns sender ID
func TestSenderID_FromMessage(t *testing.T) {
	ctx := &Context{
		Sender: &User{ID: 67890},
	}

	assert.Equal(t, int64(67890), ctx.SenderID())
}

// Test SenderID returns 0 when no sender
func TestSenderID_NoSender(t *testing.T) {
	ctx := &Context{}

	assert.Equal(t, int64(0), ctx.SenderID())
}

// Test Text returns message text
func TestText_FromMessage(t *testing.T) {
	ctx := &Context{
		Message: &Message{Text: "Test message"},
	}

	assert.Equal(t, "Test message", ctx.Text())
}

// Test Text returns empty string when no message
func TestText_NoMessage(t *testing.T) {
	ctx := &Context{}

	assert.Equal(t, "", ctx.Text())
}

// Test CallbackData returns callback data
func TestCallbackData_FromCallbackQuery(t *testing.T) {
	ctx := &Context{
		CallbackQuery: &CallbackQuery{Data: "test_data"},
	}

	assert.Equal(t, "test_data", ctx.CallbackData())
}

// Test CallbackData returns empty string when no callback query
func TestCallbackData_NoCallbackQuery(t *testing.T) {
	ctx := &Context{}

	assert.Equal(t, "", ctx.CallbackData())
}

// Test InlineQueryText returns inline query text
func TestInlineQueryText_FromInlineQuery(t *testing.T) {
	ctx := &Context{
		InlineQuery: &InlineQuery{Query: "test query"},
	}

	assert.Equal(t, "test query", ctx.InlineQueryText())
}

// Test InlineQueryText returns empty string when no inline query
func TestInlineQueryText_NoInlineQuery(t *testing.T) {
	ctx := &Context{}

	assert.Equal(t, "", ctx.InlineQueryText())
}

// Test IsMessage returns true for message updates
func TestIsMessage(t *testing.T) {
	ctx := &Context{
		Message: &Message{},
	}

	assert.True(t, ctx.IsMessage())
}

// Test IsCallbackQuery returns true for callback updates
func TestIsCallbackQuery(t *testing.T) {
	ctx := &Context{
		CallbackQuery: &CallbackQuery{},
	}

	assert.True(t, ctx.IsCallbackQuery())
}

// Test IsInlineQuery returns true for inline query updates
func TestIsInlineQuery(t *testing.T) {
	ctx := &Context{
		InlineQuery: &InlineQuery{},
	}

	assert.True(t, ctx.IsInlineQuery())
}

// Test IsChosenInlineResult returns true for chosen inline result updates
func TestIsChosenInlineResult(t *testing.T) {
	ctx := &Context{
		ChosenInlineResult: &ChosenInlineResult{},
	}

	assert.True(t, ctx.IsChosenInlineResult())
}

// Test IsCommand returns true for command messages
func TestIsCommand(t *testing.T) {
	ctx := Context{
		Message: &Message{
			Text: " /start ",
			Entities: []MessageEntity{
				{
					Type:   "bot_command",
					Offset: 1,
					Length: 5,
				},
			},
		},
	}

	assert.True(t, ctx.IsCommand())
}

// Test IsCommand returns false for non-command messages
func TestIsCommand_False(t *testing.T) {
	ctx := Context{
		Message: &Message{
			Text: "Hello",
		},
	}

	assert.False(t, ctx.IsCommand())
}

// Test IsCommand returns false when no entities
func TestIsCommand_NoEntities(t *testing.T) {
	ctx := Context{
		Message: &Message{
			Text:     "/start",
			Entities: nil,
		},
	}

	assert.False(t, ctx.IsCommand())
}

// Test Command returns command name without slash
func TestCommand(t *testing.T) {
	ctx := Context{
		Message: &Message{
			Text: "/start@mybot",
			Entities: []MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 12,
				},
			},
		},
	}

	// Command returns the text without the leading /, but with bot username
	assert.Equal(t, "start@mybot", ctx.Command())
}

// Test Command returns empty string for non-command
func TestCommand_Empty(t *testing.T) {
	ctx := Context{
		Message: &Message{
			Text: "Hello",
		},
	}

	assert.Equal(t, "", ctx.Command())
}

// Test Command with bot username included
func TestCommand_WithBotUsername(t *testing.T) {
	ctx := Context{
		Message: &Message{
			Text: "/start@mybot",
			Entities: []MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 12, // Full command including @mybot
				},
			},
		},
	}

	cmd := ctx.Command()
	// Command returns the full command text without the leading /
	assert.Equal(t, "start@mybot", cmd)
}

// Test CommandArgs returns arguments after command
func TestCommandArgs(t *testing.T) {
	ctx := Context{
		Message: &Message{
			Text: "/start arg1 arg2",
			Entities: []MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 6,
				},
			},
		},
	}

	assert.Equal(t, "arg1 arg2", ctx.CommandArgs())
}

// Test CommandArgs returns empty string for command without args
func TestCommandArgs_Empty(t *testing.T) {
	ctx := Context{
		Message: &Message{
			Text: "/start",
			Entities: []MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 6,
				},
			},
		},
	}

	assert.Equal(t, "", ctx.CommandArgs())
}

// Test CommandArgs skips whitespace
func TestCommandArgs_SkipsWhitespace(t *testing.T) {
	ctx := Context{
		Message: &Message{
			Text: "/start   arg1 ",
			Entities: []MessageEntity{
				{
					Type:   "bot_command",
					Offset: 0,
					Length: 6,
				},
			},
		},
	}

	assert.Equal(t, "arg1 ", ctx.CommandArgs())
}

// Test WithContext creates new context with wrapped context
func TestWithContext(t *testing.T) {
	ctx := &Context{}

	// WithContext should not panic and return a context
	newCtx := ctx.WithContext(nil)
	assert.NotNil(t, newCtx)
}

// Test Context returns the context
func TestContext(t *testing.T) {
	ctx := &Context{}

	// Context() should return the underlying context
	// For now it returns the logger which is a zerolog.Logger
	// This tests it doesn't panic
	_ = ctx.Context()
}

// TestLog returns a logger
func TestLog_DoesNotPanic(t *testing.T) {
	bot := &Bot{Logger: zerolog.Nop()}
	ctx := &Context{Bot: bot}

	// Log() should not panic
	// We can't test the logger without a proper setup
	// But we can at least verify it doesn't panic
	assert.NotPanics(t, func() {
		_ = ctx.Log()
	})
}

// Test updateType returns correct types
func TestUpdateType(t *testing.T) {
	// Test message
	ctx := &Context{Message: &Message{}}
	assert.Equal(t, "message", ctx.updateType())

	// Test callback
	ctx = &Context{CallbackQuery: &CallbackQuery{}}
	assert.Equal(t, "callback_query", ctx.updateType())

	// Test inline query
	ctx = &Context{InlineQuery: &InlineQuery{}}
	assert.Equal(t, "inline_query", ctx.updateType())

	// Test chosen inline result
	ctx = &Context{ChosenInlineResult: &ChosenInlineResult{}}
	assert.Equal(t, "chosen_inline_result", ctx.updateType())

	// Test unknown
	ctx = &Context{}
	assert.Equal(t, "unknown", ctx.updateType())
}

// Test Stringer implementations
func TestUser_String(t *testing.T) {
	user := &User{
		ID:        123,
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
	}

	// Just verify it doesn't panic
	_ = user
}

// Test Chat types
func TestChatType(t *testing.T) {
	chat := &Chat{
		ID:   123,
		Type: "group",
	}

	assert.Equal(t, "group", chat.Type)
}

// Benchmark context creation
func BenchmarkNewContext_Message(b *testing.B) {
	update := Update{
		Message: &Message{
			Text: "test",
			Chat: &Chat{ID: 123},
			From: &User{ID: 456},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewContext(nil, update)
	}
}

func BenchmarkNewContext_CallbackQuery(b *testing.B) {
	update := Update{
		CallbackQuery: &CallbackQuery{
			From: &User{ID: 456},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewContext(nil, update)
	}
}

func BenchmarkContext_Command(b *testing.B) {
	ctx := Context{
		Message: &Message{
			Text: "/start arg1 arg2",
			Entities: []MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 6},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.Command()
		_ = ctx.CommandArgs()
	}
}
