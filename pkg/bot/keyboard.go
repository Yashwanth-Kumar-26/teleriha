package bot

import (
	"fmt"
)

// Keyboard provides helper functions for creating various types of keyboards.

// NewInlineKeyboard creates a new inline keyboard.
func NewInlineKeyboard(rows ...[]InlineKeyboardButton) *InlineKeyboardMarkup {
	return &InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

// NewInlineKeyboardRow creates a new row of inline buttons.
func NewInlineKeyboardRow(buttons ...InlineKeyboardButton) []InlineKeyboardButton {
	return buttons
}

// InlineButton creates a new inline keyboard button with text and callback data.
func InlineButton(text, callbackData string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:         text,
		CallbackData: callbackData,
	}
}

// InlineButtonURL creates a new inline keyboard button with text and URL.
func InlineButtonURL(text, url string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		URL:  url,
	}
}

// InlineButtonSwitch creates a new inline keyboard button that switches to inline mode.
func InlineButtonSwitch(text, query string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:                 text,
		SwitchInlineQuery:    query,
		SwitchInlineQueryCurrentChat: query,
	}
}

// InlineButtonPay creates a new inline keyboard button for payments.
func InlineButtonPay(text string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text: text,
		Pay:  true,
	}
}

// NewReplyKeyboard creates a new reply keyboard.
func NewReplyKeyboard(rows ...[]KeyboardButton) *ReplyKeyboardMarkup {
	return &ReplyKeyboardMarkup{
		Keyboard: rows,
	}
}

// NewReplyKeyboardRow creates a new row of reply buttons.
func NewReplyKeyboardRow(buttons ...KeyboardButton) []KeyboardButton {
	return buttons
}

// ReplyButton creates a new reply keyboard button with text.
func ReplyButton(text string) KeyboardButton {
	return KeyboardButton{
		Text: text,
	}
}

// ReplyButtonContact creates a new reply keyboard button that requests contact.
func ReplyButtonContact(text string) KeyboardButton {
	return KeyboardButton{
		Text:           text,
		RequestContact: true,
	}
}

// ReplyButtonLocation creates a new reply keyboard button that requests location.
func ReplyButtonLocation(text string) KeyboardButton {
	return KeyboardButton{
		Text:            text,
		RequestLocation: true,
	}
}

// ReplyButtonPoll creates a new reply keyboard button that requests a poll.
func ReplyButtonPoll(text string, pollType string) KeyboardButton {
	return KeyboardButton{
		Text:   text,
		RequestPoll: &KeyboardButtonPollType{
			Type: pollType,
		},
	}
}

// NewForceReply creates a new force reply keyboard.
func NewForceReply(placeholder string) *ForceReply {
	return &ForceReply{
		InputFieldPlaceholder: placeholder,
		Selective:            true,
	}
}

// RemoveKeyboard creates a keyboard that removes the current keyboard.
func RemoveKeyboard() *ReplyKeyboardMarkup {
	return &ReplyKeyboardMarkup{
		Keyboard:        [][]KeyboardButton{},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
		Selective:       true,
	}
}

// HideKeyboard creates a keyboard that hides the current keyboard.
func HideKeyboard() *ReplyKeyboardMarkup {
	return &ReplyKeyboardMarkup{
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
		Selective:       true,
	}
}

// InlineKeyboardBuilder provides a fluent interface for building inline keyboards.
type InlineKeyboardBuilder struct {
	rows [][]InlineKeyboardButton
}

// NewInlineKeyboardBuilder creates a new InlineKeyboardBuilder.
func NewInlineKeyboardBuilder() *InlineKeyboardBuilder {
	return &InlineKeyboardBuilder{}
}

// AddRow adds a new row with the given buttons.
func (b *InlineKeyboardBuilder) AddRow(buttons ...InlineKeyboardButton) *InlineKeyboardBuilder {
	b.rows = append(b.rows, buttons)
	return b
}

// AddButton adds a single button in a new row.
func (b *InlineKeyboardBuilder) AddButton(text, callbackData string) *InlineKeyboardBuilder {
	b.rows = append(b.rows, []InlineKeyboardButton{
		{Text: text, CallbackData: callbackData},
	})
	return b
}

// AddButtons adds multiple buttons in a single row.
func (b *InlineKeyboardBuilder) AddButtons(buttons ...InlineKeyboardButton) *InlineKeyboardBuilder {
	b.rows = append(b.rows, buttons)
	return b
}

// AddURLButton adds a URL button in a new row.
func (b *InlineKeyboardBuilder) AddURLButton(text, url string) *InlineKeyboardBuilder {
	b.rows = append(b.rows, []InlineKeyboardButton{
		{Text: text, URL: url},
	})
	return b
}

// Add Pay button
func (b *InlineKeyboardBuilder) AddPayButton(text string) *InlineKeyboardBuilder {
	b.rows = append(b.rows, []InlineKeyboardButton{
		{Text: text, Pay: true},
	})
	return b
}

// Build returns the InlineKeyboardMarkup.
func (b *InlineKeyboardBuilder) Build() *InlineKeyboardMarkup {
	return &InlineKeyboardMarkup{
		InlineKeyboard: b.rows,
	}
}

// ReplyKeyboardBuilder provides a fluent interface for building reply keyboards.
type ReplyKeyboardBuilder struct {
	resizeKeyboard bool
	onetime        bool
	selective      bool
	placeholder    string
	rows          [][]KeyboardButton
}

// NewReplyKeyboardBuilder creates a new ReplyKeyboardBuilder.
func NewReplyKeyboardBuilder() *ReplyKeyboardBuilder {
	return &ReplyKeyboardBuilder{}
}

// Resize sets the resize keyboard option.
func (b *ReplyKeyboardBuilder) Resize() *ReplyKeyboardBuilder {
	b.resizeKeyboard = true
	return b
}

// OneTime sets the one-time keyboard option.
func (b *ReplyKeyboardBuilder) OneTime() *ReplyKeyboardBuilder {
	b.onetime = true
	return b
}

// Selective sets the selective option.
func (b *ReplyKeyboardBuilder) Selective() *ReplyKeyboardBuilder {
	b.selective = true
	return b
}

// Placeholder sets the input field placeholder.
func (b *ReplyKeyboardBuilder) Placeholder(text string) *ReplyKeyboardBuilder {
	b.placeholder = text
	return b
}

// AddRow adds a new row with the given buttons.
func (b *ReplyKeyboardBuilder) AddRow(buttons ...KeyboardButton) *ReplyKeyboardBuilder {
	b.rows = append(b.rows, buttons)
	return b
}

// AddButton adds a single button in a new row.
func (b *ReplyKeyboardBuilder) AddButton(text string) *ReplyKeyboardBuilder {
	b.rows = append(b.rows, []KeyboardButton{
		{Text: text},
	})
	return b
}

// AddButtons adds multiple buttons in a single row.
func (b *ReplyKeyboardBuilder) AddButtons(buttons ...KeyboardButton) *ReplyKeyboardBuilder {
	b.rows = append(b.rows, buttons)
	return b
}

// Build returns the ReplyKeyboardMarkup.
func (b *ReplyKeyboardBuilder) Build() *ReplyKeyboardMarkup {
	return &ReplyKeyboardMarkup{
		Keyboard:             b.rows,
		ResizeKeyboard:       b.resizeKeyboard,
		OneTimeKeyboard:      b.onetime,
		Selective:            b.selective,
		InputFieldPlaceholder: b.placeholder,
	}
}

// ForceReply represents a force reply keyboard (users will have to reply to the message).
type ForceReply struct {
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
	Selective             bool   `json:"selective,omitempty"`
}

// Simple utility functions for common keyboard patterns

// YesNoKeyboard creates a simple Yes/No inline keyboard.
func YesNoKeyboard(yesCallback, noCallback string) *InlineKeyboardMarkup {
	return NewInlineKeyboard(
		NewInlineKeyboardRow(
			InlineButton("Yes", yesCallback),
			InlineButton("No", noCallback),
		),
	)
}

// ConfirmCancelKeyboard creates a Confirm/Cancel inline keyboard.
func ConfirmCancelKeyboard(confirmCallback, cancelCallback string) *InlineKeyboardMarkup {
	return NewInlineKeyboard(
		NewInlineKeyboardRow(
			InlineButton("Confirm", confirmCallback),
			InlineButton("Cancel", cancelCallback),
		),
	)
}

// NumberKeyboard creates a numeric keyboard with buttons 1-9 and 0.
func NumberKeyboard() *ReplyKeyboardMarkup {
	return NewReplyKeyboard(
		NewReplyKeyboardRow(
			ReplyButton("1"),
			ReplyButton("2"),
			ReplyButton("3"),
		),
		NewReplyKeyboardRow(
			ReplyButton("4"),
			ReplyButton("5"),
			ReplyButton("6"),
		),
		NewReplyKeyboardRow(
			ReplyButton("7"),
			ReplyButton("8"),
			ReplyButton("9"),
		),
		NewReplyKeyboardRow(
			ReplyButton("0"),
		),
	)
}

// PaginationKeyboard creates a pagination keyboard with Previous and Next buttons.
func PaginationKeyboard(prevCallback, nextCallback string, currentPage, totalPages int) *InlineKeyboardMarkup {
	var buttons []InlineKeyboardButton

	// Previous button
	if currentPage > 1 {
		buttons = append(buttons, InlineButton("⬅️ Previous", prevCallback))
	}

	// Page indicator (as text, not a button)
	pageText := fmt.Sprintf("Page %d/%d", currentPage, totalPages)
	buttons = append(buttons, InlineKeyboardButton{
		Text:         pageText,
		CallbackData: "no_op",
	})

	// Next button
	if currentPage < totalPages {
		buttons = append(buttons, InlineButton("Next ➡️", nextCallback))
	}

	return NewInlineKeyboard(buttons)
}
