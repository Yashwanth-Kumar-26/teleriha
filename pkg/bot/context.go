package bot

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

// Context represents the context for handling a Telegram update.
// It provides access to the bot, the update, and convenience methods.
type Context struct {
	// Bot is the TeleRiHa bot instance
	Bot *Bot

	// Update is the raw update from Telegram
	Update Update

	// Message is the message from the update (may be nil)
	Message *Message

	// CallbackQuery is the callback query from the update (may be nil)
	CallbackQuery *CallbackQuery

	// InlineQuery is the inline query from the update (may be nil)
	InlineQuery *InlineQuery

	// ChosenInlineResult is the chosen inline result (may be nil)
	ChosenInlineResult *ChosenInlineResult

	// Chat is the chat from the update (cached for convenience)
	Chat *Chat

	// Sender is the sender user (cached for convenience)
	Sender *User

	// ctx is the underlying Go context for cancellation
	ctx context.Context
}

// NewContext creates a new Context from an update.
func NewContext(b *Bot, update Update) *Context {
	ctx := &Context{
		Bot:    b,
		Update: update,
		ctx:    context.Background(),
	}

	// Set convenience fields
	switch {
	case update.Message != nil:
		ctx.Message = update.Message
		ctx.Chat = update.Message.Chat
		ctx.Sender = update.Message.From
	case update.CallbackQuery != nil:
		ctx.CallbackQuery = update.CallbackQuery
		if update.CallbackQuery.Message != nil {
			ctx.Message = update.CallbackQuery.Message
			ctx.Chat = update.CallbackQuery.Message.Chat
		}
		ctx.Sender = update.CallbackQuery.From
	case update.InlineQuery != nil:
		ctx.InlineQuery = update.InlineQuery
		ctx.Sender = update.InlineQuery.From
	case update.ChosenInlineResult != nil:
		ctx.ChosenInlineResult = update.ChosenInlineResult
		ctx.Sender = update.ChosenInlineResult.From
	}

	return ctx
}

// WithContext wrapper for Go's context.Context.
func (c *Context) WithContext(ctx context.Context) *Context {
	return &Context{
		Bot:              c.Bot,
		Update:           c.Update,
		Message:          c.Message,
		CallbackQuery:    c.CallbackQuery,
		InlineQuery:      c.InlineQuery,
		ChosenInlineResult: c.ChosenInlineResult,
		Chat:             c.Chat,
		Sender:           c.Sender,
		ctx:              ctx,
	}
}

// Context returns the underlying Go context.
func (c *Context) Context() context.Context {
	return c.ctx
}

// ChatID returns the chat ID, or 0 if not available.
func (c *Context) ChatID() int64 {
	if c.Chat != nil {
		return c.Chat.ID
	}
	return 0
}

// SenderID returns the sender user ID, or 0 if not available.
func (c *Context) SenderID() int64 {
	if c.Sender != nil {
		return c.Sender.ID
	}
	return 0
}

// Text returns the message text, or empty string if not available.
func (c *Context) Text() string {
	if c.Message != nil {
		return c.Message.Text
	}
	return ""
}

// CallbackData returns the callback query data, or empty string if not available.
func (c *Context) CallbackData() string {
	if c.CallbackQuery != nil {
		return c.CallbackQuery.Data
	}
	return ""
}

// InlineQueryText returns the inline query text, or empty string if not available.
func (c *Context) InlineQueryText() string {
	if c.InlineQuery != nil {
		return c.InlineQuery.Query
	}
	return ""
}

// Reply sends a text message to the chat this context came from.
func (c *Context) Reply(text string, opts ...MessageOption) (*Message, error) {
	chatID := c.ChatID()
	if chatID == 0 {
		return nil, fmt.Errorf("no chat ID in context")
	}
	return c.Bot.SendMessage(chatID, text, opts...)
}

// AnswerCallback answers a callback query.
func (c *Context) AnswerCallback(text string, showAlert bool) error {
	if c.CallbackQuery == nil {
		return fmt.Errorf("no callback query in context")
	}
	return c.Bot.AnswerCallbackQuery(c.CallbackQuery.ID, text, showAlert)
}

// EditMessageText edits the text of a message.
func (c *Context) EditMessageText(text string, opts ...EditMessageOption) (*Message, error) {
	if c.CallbackQuery != nil && c.CallbackQuery.Message != nil {
		return c.Bot.EditMessageText(
			c.CallbackQuery.Message.Chat.ID,
			c.CallbackQuery.Message.MessageID,
			text,
			opts...,
		)
	}
	if c.Message != nil {
		return c.Bot.EditMessageText(
			c.Message.Chat.ID,
			c.Message.MessageID,
			text,
			opts...,
		)
	}
	return nil, fmt.Errorf("no message to edit in context")
}

// DeleteMessage deletes a message.
func (c *Context) DeleteMessage(messageID int64) error {
	return c.Bot.DeleteMessage(c.ChatID(), messageID)
}

// Delete deletes the message this context came from.
func (c *Context) Delete() error {
	if c.Message != nil {
		return c.Bot.DeleteMessage(c.Message.Chat.ID, c.Message.MessageID)
	}
	return fmt.Errorf("no message to delete in context")
}

// Send sends a message to a specific chat ID.
func (c *Context) Send(chatID int64, text string, opts ...MessageOption) (*Message, error) {
	return c.Bot.SendMessage(chatID, text, opts...)
}

// SendPhoto sends a photo to a chat.
func (c *Context) SendPhoto(chatID int64, photo interface{}, caption string, opts ...SendPhotoOption) (*Message, error) {
	return c.Bot.SendPhoto(chatID, photo, caption, opts...)
}

// SendDocument sends a document to a chat.
func (c *Context) SendDocument(chatID int64, document interface{}, caption string, opts ...SendDocumentOption) (*Message, error) {
	return c.Bot.SendDocument(chatID, document, caption, opts...)
}

// IsCallbackQuery returns true if this context came from a callback query.
func (c *Context) IsCallbackQuery() bool {
	return c.CallbackQuery != nil
}

// IsInlineQuery returns true if this context came from an inline query.
func (c *Context) IsInlineQuery() bool {
	return c.InlineQuery != nil
}

// IsChosenInlineResult returns true if this context came from a chosen inline result.
func (c *Context) IsChosenInlineResult() bool {
	return c.ChosenInlineResult != nil
}

// IsMessage returns true if this context came from a message.
func (c *Context) IsMessage() bool {
	return c.Message != nil
}

// IsCommand returns true if the message is a command.
func (c *Context) IsCommand() bool {
	if c.Message == nil || c.Message.Entities == nil {
		return false
	}
	for _, entity := range c.Message.Entities {
		if entity.Type == "bot_command" {
			return true
		}
	}
	return false
}

// Command returns the command name (without the bot username and @).
// Returns empty string if the message is not a command.
func (c *Context) Command() string {
	if !c.IsCommand() || c.Message == nil {
		return ""
	}

	text := c.Message.Text
	if len(text) == 0 {
		return ""
	}

	// Find the bot_command entity
	for _, entity := range c.Message.Entities {
		if entity.Type == "bot_command" {
			if entity.Offset >= 0 && entity.Offset < len(text) {
				end := entity.Offset + entity.Length
				if end > len(text) {
					end = len(text)
				}
				cmd := text[entity.Offset:end]

				// Remove the command from the text for arguments
				// Return just the command name without /
				if len(cmd) > 0 && cmd[0] == '/' {
					return cmd[1:]
				}
				return cmd
			}
		}
	}
	return ""
}

// CommandArgs returns the arguments after the command.
// Returns empty string if the message is not a command.
func (c *Context) CommandArgs() string {
	if !c.IsCommand() || c.Message == nil {
		return ""
	}

	text := c.Message.Text
	if len(text) == 0 {
		return ""
	}

	// Find the bot_command entity
	for _, entity := range c.Message.Entities {
		if entity.Type == "bot_command" {
			if entity.Offset >= 0 && entity.Offset < len(text) {
				start := entity.Offset + entity.Length
				if start < len(text) {
					// Skip whitespace after command
					for start < len(text) && (text[start] == ' ' || text[start] == '\t') {
						start++
					}
					return text[start:]
				}
			}
		}
	}
	return ""
}

// Log is a convenience method for logging.
func (c *Context) Log() zerolog.Logger {
	return c.Bot.Logger.With().
		Int64("chat_id", c.ChatID()).
		Int64("user_id", c.SenderID()).
		Logger()
}
