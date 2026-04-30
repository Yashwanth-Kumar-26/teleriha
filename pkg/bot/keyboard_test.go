package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test NewInlineKeyboard creates a keyboard
func TestNewInlineKeyboard(t *testing.T) {
	rows := [][]InlineKeyboardButton{
		{ {Text: "A", CallbackData: "a"}, {Text: "B", CallbackData: "b"} },
	}
	keyboard := NewInlineKeyboard(rows...)

	assert.NotNil(t, keyboard)
	assert.Equal(t, rows, keyboard.InlineKeyboard)
}

// Test NewInlineKeyboardRow creates a row
func TestNewInlineKeyboardRow(t *testing.T) {
	buttons := []InlineKeyboardButton{
		{Text: "A", CallbackData: "a"},
		{Text: "B", CallbackData: "b"},
	}
	row := NewInlineKeyboardRow(buttons...)

	assert.Equal(t, buttons, row)
}

// Test InlineButton creates a button with callback data
func TestInlineButton(t *testing.T) {
	button := InlineButton("Click me", "action")

	assert.Equal(t, "Click me", button.Text)
	assert.Equal(t, "action", button.CallbackData)
}

// Test InlineButtonURL creates a button with URL
func TestInlineButtonURL(t *testing.T) {
	button := InlineButtonURL("Visit", "https://example.com")

	assert.Equal(t, "Visit", button.Text)
	assert.Equal(t, "https://example.com", button.URL)
}

// Test InlineButtonSwitch creates a switch button
func TestInlineButtonSwitch(t *testing.T) {
	button := InlineButtonSwitch("Inline", "query")

	assert.Equal(t, "Inline", button.Text)
	assert.Equal(t, "query", button.SwitchInlineQuery)
	assert.Equal(t, "query", button.SwitchInlineQueryCurrentChat)
}

// Test InlineButtonPay creates a pay button
func TestInlineButtonPay(t *testing.T) {
	button := InlineButtonPay("Pay Now")

	assert.Equal(t, "Pay Now", button.Text)
	assert.True(t, button.Pay)
}

// Test NewReplyKeyboard creates a reply keyboard
func TestNewReplyKeyboard(t *testing.T) {
	rows := [][]KeyboardButton{
		{ {Text: "A"}, {Text: "B"} },
	}
	keyboard := NewReplyKeyboard(rows...)

	assert.NotNil(t, keyboard)
	assert.Equal(t, rows, keyboard.Keyboard)
}

// Test NewReplyKeyboardRow creates a row of reply buttons
func TestNewReplyKeyboardRow(t *testing.T) {
	buttons := []KeyboardButton{
		{Text: "A"},
		{Text: "B"},
	}
	row := NewReplyKeyboardRow(buttons...)

	assert.Equal(t, buttons, row)
}

// Test ReplyButton creates a simple reply button
func TestReplyButton(t *testing.T) {
	button := ReplyButton("Click")

	assert.Equal(t, "Click", button.Text)
}

// Test ReplyButtonContact creates a contact request button
func TestReplyButtonContact(t *testing.T) {
	button := ReplyButtonContact("Share Contact")

	assert.Equal(t, "Share Contact", button.Text)
	assert.True(t, button.RequestContact)
}

// Test ReplyButtonLocation creates a location request button
func TestReplyButtonLocation(t *testing.T) {
	button := ReplyButtonLocation("Share Location")

	assert.Equal(t, "Share Location", button.Text)
	assert.True(t, button.RequestLocation)
}

// Test ReplyButtonPoll creates a poll request button
func TestReplyButtonPoll(t *testing.T) {
	button := ReplyButtonPoll("Create Poll", "quiz")

	assert.Equal(t, "Create Poll", button.Text)
	assert.NotNil(t, button.RequestPoll)
	assert.Equal(t, "quiz", button.RequestPoll.Type)
}

// Test NewForceReply creates a force reply keyboard
func TestNewForceReply(t *testing.T) {
	keyboard := NewForceReply("Type something...")

	assert.NotNil(t, keyboard)
	assert.Equal(t, "Type something...", keyboard.InputFieldPlaceholder)
	assert.True(t, keyboard.Selective)
}

// Test RemoveKeyboard creates a remove keyboard
func TestRemoveKeyboard(t *testing.T) {
	keyboard := RemoveKeyboard()

	assert.NotNil(t, keyboard)
	assert.Empty(t, keyboard.Keyboard)
	assert.True(t, keyboard.ResizeKeyboard)
	assert.True(t, keyboard.OneTimeKeyboard)
	assert.True(t, keyboard.Selective)
}

// Test HideKeyboard creates a hide keyboard
func TestHideKeyboard(t *testing.T) {
	keyboard := HideKeyboard()

	assert.NotNil(t, keyboard)
	assert.True(t, keyboard.ResizeKeyboard)
	assert.True(t, keyboard.OneTimeKeyboard)
	assert.True(t, keyboard.Selective)
}

// Test InlineKeyboardBuilder AddRow
func TestInlineKeyboardBuilder_AddRow(t *testing.T) {
	builder := NewInlineKeyboardBuilder()
	button := InlineKeyboardButton{Text: "A", CallbackData: "a"}

	builder.AddRow(button)
	result := builder.Build()

	assert.Len(t, result.InlineKeyboard, 1)
	assert.Len(t, result.InlineKeyboard[0], 1)
	assert.Equal(t, "A", result.InlineKeyboard[0][0].Text)
}

// Test InlineKeyboardBuilder AddButton
func TestInlineKeyboardBuilder_AddButton(t *testing.T) {
	builder := NewInlineKeyboardBuilder()

	builder.AddButton("Click", "action")
	result := builder.Build()

	assert.Len(t, result.InlineKeyboard, 1)
	assert.Len(t, result.InlineKeyboard[0], 1)
	assert.Equal(t, "Click", result.InlineKeyboard[0][0].Text)
	assert.Equal(t, "action", result.InlineKeyboard[0][0].CallbackData)
}

// Test InlineKeyboardBuilder AddButtons
func TestInlineKeyboardBuilder_AddButtons(t *testing.T) {
	builder := NewInlineKeyboardBuilder()
	button1 := InlineKeyboardButton{Text: "A", CallbackData: "a"}
	button2 := InlineKeyboardButton{Text: "B", CallbackData: "b"}

	builder.AddButtons(button1, button2)
	result := builder.Build()

	assert.Len(t, result.InlineKeyboard, 1)
	assert.Len(t, result.InlineKeyboard[0], 2)
}

// Test InlineKeyboardBuilder AddURLButton
func TestInlineKeyboardBuilder_AddURLButton(t *testing.T) {
	builder := NewInlineKeyboardBuilder()

	builder.AddURLButton("Visit", "https://example.com")
	result := builder.Build()

	assert.Len(t, result.InlineKeyboard, 1)
	assert.Equal(t, "Visit", result.InlineKeyboard[0][0].Text)
	assert.Equal(t, "https://example.com", result.InlineKeyboard[0][0].URL)
}

// Test InlineKeyboardBuilder AddPayButton
func TestInlineKeyboardBuilder_AddPayButton(t *testing.T) {
	builder := NewInlineKeyboardBuilder()

	builder.AddPayButton("Pay")
	result := builder.Build()

	assert.Len(t, result.InlineKeyboard, 1)
	assert.True(t, result.InlineKeyboard[0][0].Pay)
}

// Test InlineKeyboardBuilder multiple rows
func TestInlineKeyboardBuilder_MultipleRows(t *testing.T) {
	builder := NewInlineKeyboardBuilder()

	builder.AddButton("A", "a").AddRow(InlineKeyboardButton{Text: "B", CallbackData: "b"})
	result := builder.Build()

	assert.Len(t, result.InlineKeyboard, 2)
}

// Test ReplyKeyboardBuilder Resize
func TestReplyKeyboardBuilder_Resize(t *testing.T) {
	builder := NewReplyKeyboardBuilder()
	builder.Resize().AddButton("A")
	result := builder.Build()

	assert.True(t, result.ResizeKeyboard)
}

// Test ReplyKeyboardBuilder OneTime
func TestReplyKeyboardBuilder_OneTime(t *testing.T) {
	builder := NewReplyKeyboardBuilder()
	builder.OneTime().AddButton("A")
	result := builder.Build()

	assert.True(t, result.OneTimeKeyboard)
}

// Test ReplyKeyboardBuilder Selective
func TestReplyKeyboardBuilder_Selective(t *testing.T) {
	builder := NewReplyKeyboardBuilder()
	builder.Selective().AddButton("A")
	result := builder.Build()

	assert.True(t, result.Selective)
}

// Test ReplyKeyboardBuilder Placeholder
func TestReplyKeyboardBuilder_Placeholder(t *testing.T) {
	builder := NewReplyKeyboardBuilder()
	builder.Placeholder("Type...").AddButton("A")
	result := builder.Build()

	assert.Equal(t, "Type...", result.InputFieldPlaceholder)
}

// Test ReplyKeyboardBuilder AddRow
func TestReplyKeyboardBuilder_AddRow(t *testing.T) {
	builder := NewReplyKeyboardBuilder()
	button := KeyboardButton{Text: "A"}

	builder.AddRow(button)
	result := builder.Build()

	assert.Len(t, result.Keyboard, 1)
	assert.Len(t, result.Keyboard[0], 1)
	assert.Equal(t, "A", result.Keyboard[0][0].Text)
}

// Test ReplyKeyboardBuilder AddButton
func TestReplyKeyboardBuilder_AddButton(t *testing.T) {
	builder := NewReplyKeyboardBuilder()

	builder.AddButton("Click")
	result := builder.Build()

	assert.Len(t, result.Keyboard, 1)
	assert.Len(t, result.Keyboard[0], 1)
	assert.Equal(t, "Click", result.Keyboard[0][0].Text)
}

// Test ReplyKeyboardBuilder AddButtons
func TestReplyKeyboardBuilder_AddButtons(t *testing.T) {
	builder := NewReplyKeyboardBuilder()
	button1 := KeyboardButton{Text: "A"}
	button2 := KeyboardButton{Text: "B"}

	builder.AddButtons(button1, button2)
	result := builder.Build()

	assert.Len(t, result.Keyboard, 1)
	assert.Len(t, result.Keyboard[0], 2)
}

// Test YesNoKeyboard creates yes/no keyboard
func TestYesNoKeyboard(t *testing.T) {
	keyboard := YesNoKeyboard("yes", "no")

	assert.NotNil(t, keyboard)
	assert.Len(t, keyboard.InlineKeyboard, 1)
	assert.Len(t, keyboard.InlineKeyboard[0], 2)
	assert.Equal(t, "Yes", keyboard.InlineKeyboard[0][0].Text)
	assert.Equal(t, "No", keyboard.InlineKeyboard[0][1].Text)
	assert.Equal(t, "yes", keyboard.InlineKeyboard[0][0].CallbackData)
	assert.Equal(t, "no", keyboard.InlineKeyboard[0][1].CallbackData)
}

// Test ConfirmCancelKeyboard creates confirm/cancel keyboard
func TestConfirmCancelKeyboard(t *testing.T) {
	keyboard := ConfirmCancelKeyboard("confirm", "cancel")

	assert.NotNil(t, keyboard)
	assert.Len(t, keyboard.InlineKeyboard, 1)
	assert.Len(t, keyboard.InlineKeyboard[0], 2)
	assert.Equal(t, "Confirm", keyboard.InlineKeyboard[0][0].Text)
	assert.Equal(t, "Cancel", keyboard.InlineKeyboard[0][1].Text)
}

// Test NumberKeyboard creates numeric keyboard
func TestNumberKeyboard(t *testing.T) {
	keyboard := NumberKeyboard()

	assert.NotNil(t, keyboard)
	assert.Len(t, keyboard.Keyboard, 4)
	// Check first row has 1, 2, 3
	assert.Equal(t, "1", keyboard.Keyboard[0][0].Text)
	assert.Equal(t, "2", keyboard.Keyboard[0][1].Text)
	assert.Equal(t, "3", keyboard.Keyboard[0][2].Text)
	// Check last row has 0
	assert.Equal(t, "0", keyboard.Keyboard[3][0].Text)
}

// Test PaginationKeyboard with middle page
func TestPaginationKeyboard_MiddlePage(t *testing.T) {
	keyboard := PaginationKeyboard("prev", "next", 2, 5)

	assert.NotNil(t, keyboard)
	assert.Len(t, keyboard.InlineKeyboard, 1)
	assert.Len(t, keyboard.InlineKeyboard[0], 3)
	// Should have Previous, Page 2/5, Next
	assert.Equal(t, "⬅️ Previous", keyboard.InlineKeyboard[0][0].Text)
	assert.Equal(t, "Page 2/5", keyboard.InlineKeyboard[0][1].Text)
	assert.Equal(t, "Next ➡️", keyboard.InlineKeyboard[0][2].Text)
}

// Test PaginationKeyboard first page
func TestPaginationKeyboard_FirstPage(t *testing.T) {
	keyboard := PaginationKeyboard("prev", "next", 1, 5)

	assert.NotNil(t, keyboard)
	assert.Len(t, keyboard.InlineKeyboard, 1)
	assert.Len(t, keyboard.InlineKeyboard[0], 2)
	// Should have Page 1/5, Next (no Previous)
	assert.Equal(t, "Page 1/5", keyboard.InlineKeyboard[0][0].Text)
	assert.Equal(t, "Next ➡️", keyboard.InlineKeyboard[0][1].Text)
}

// Test PaginationKeyboard last page
func TestPaginationKeyboard_LastPage(t *testing.T) {
	keyboard := PaginationKeyboard("prev", "next", 5, 5)

	assert.NotNil(t, keyboard)
	assert.Len(t, keyboard.InlineKeyboard, 1)
	assert.Len(t, keyboard.InlineKeyboard[0], 2)
	// Should have Previous, Page 5/5 (no Next)
	assert.Equal(t, "⬅️ Previous", keyboard.InlineKeyboard[0][0].Text)
	assert.Equal(t, "Page 5/5", keyboard.InlineKeyboard[0][1].Text)
}

// Test PaginationKeyboard single page
func TestPaginationKeyboard_SinglePage(t *testing.T) {
	keyboard := PaginationKeyboard("prev", "next", 1, 1)

	assert.NotNil(t, keyboard)
	assert.Len(t, keyboard.InlineKeyboard, 1)
	assert.Len(t, keyboard.InlineKeyboard[0], 1)
	// Should only have Page 1/1
	assert.Equal(t, "Page 1/1", keyboard.InlineKeyboard[0][0].Text)
}

// Benchmark InlineKeyboard creation
func BenchmarkNewInlineKeyboard(b *testing.B) {
	rows := [][]InlineKeyboardButton{
		{ {Text: "A", CallbackData: "a"} },
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewInlineKeyboard(rows...)
	}
}

func BenchmarkInlineKeyboardBuilder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder := NewInlineKeyboardBuilder()
		builder.AddButton("A", "a").AddButton("B", "b")
		_ = builder.Build()
	}
}

func BenchmarkNumberKeyboard(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NumberKeyboard()
	}
}

func BenchmarkYesNoKeyboard(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = YesNoKeyboard("yes", "no")
	}
}
