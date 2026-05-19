package bot

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApplicationBuilder(t *testing.T) {
	app := NewApplicationBuilder().
		Token("test-token").
		Build()

	assert.NotNil(t, app)
	assert.Equal(t, "test-token", app.Bot.Token)
}

func TestApplication_AddHandler(t *testing.T) {
	app := NewApplicationBuilder().
		Token("test-token").
		Build()

	commandCalled := false
	messageFilterCalled := false
	callbackCalled := false
	pollCalled := false
	pollAnswerCalled := false
	chatMemberCalled := false
	myChatMemberCalled := false
	chatJoinRequestCalled := false

	app.AddHandler(CommandHandler("/start", func(ctx *Context) error {
		commandCalled = true
		return nil
	}))

	app.AddHandler(MessageHandler(FilterPhoto(), func(ctx *Context) error {
		messageFilterCalled = true
		return nil
	}))

	app.AddHandler(CallbackQueryHandler("click_", func(ctx *Context) error {
		callbackCalled = true
		return nil
	}))

	app.AddHandler(PollHandler(func(ctx *Context) error {
		pollCalled = true
		return nil
	}))

	app.AddHandler(PollAnswerHandler(func(ctx *Context) error {
		pollAnswerCalled = true
		return nil
	}))

	app.AddHandler(ChatMemberHandler(func(ctx *Context) error {
		chatMemberCalled = true
		return nil
	}))

	app.AddHandler(MyChatMemberHandler(func(ctx *Context) error {
		myChatMemberCalled = true
		return nil
	}))

	app.AddHandler(ChatJoinRequestHandler(func(ctx *Context) error {
		chatJoinRequestCalled = true
		return nil
	}))

	// Trigger command
	msgCmd := testMessage("/start", true)
	ctxCmd := NewContext(app.Bot, Update{Message: msgCmd})
	app.Bot.Router.HandleMessage(ctxCmd, *msgCmd)
	assert.True(t, commandCalled)

	// Trigger filter handler
	msgPhoto := testMessage("", false)
	msgPhoto.Photo = []PhotoSize{{FileID: "photo123"}}
	ctxPhoto := NewContext(app.Bot, Update{Message: msgPhoto})
	app.Bot.Router.HandleMessage(ctxPhoto, *msgPhoto)
	assert.True(t, messageFilterCalled)

	// Trigger callback query
	cb := testCallbackQuery("click_button")
	ctxCb := NewContext(app.Bot, Update{CallbackQuery: cb})
	app.Bot.Router.HandleCallbackQuery(ctxCb, *cb)
	assert.True(t, callbackCalled)

	// Trigger Poll
	poll := Poll{ID: "poll123", Question: "Yes or No?"}
	ctxPoll := NewContext(app.Bot, Update{Poll: &poll})
	app.Bot.Router.HandlePoll(ctxPoll, poll)
	assert.True(t, pollCalled)

	// Trigger PollAnswer
	pollAnswer := PollAnswer{PollID: "poll123", OptionIDs: []int{0}}
	ctxPollAnswer := NewContext(app.Bot, Update{PollAnswer: &pollAnswer})
	app.Bot.Router.HandlePollAnswer(ctxPollAnswer, pollAnswer)
	assert.True(t, pollAnswerCalled)

	// Trigger ChatMember
	chatMember := ChatMemberUpdated{Chat: &Chat{ID: 111}, From: &User{ID: 222}}
	ctxChatMember := NewContext(app.Bot, Update{ChatMember: &chatMember})
	app.Bot.Router.HandleChatMember(ctxChatMember, chatMember)
	assert.True(t, chatMemberCalled)

	// Trigger MyChatMember
	myChatMember := ChatMemberUpdated{Chat: &Chat{ID: 111}, From: &User{ID: 222}}
	ctxMyChatMember := NewContext(app.Bot, Update{MyChatMember: &myChatMember})
	app.Bot.Router.HandleMyChatMember(ctxMyChatMember, myChatMember)
	assert.True(t, myChatMemberCalled)

	// Trigger ChatJoinRequest
	chatJoinRequest := ChatJoinRequest{Chat: &Chat{ID: 111}, From: &User{ID: 222}}
	ctxChatJoinRequest := NewContext(app.Bot, Update{ChatJoinRequest: &chatJoinRequest})
	app.Bot.Router.HandleChatJoinRequest(ctxChatJoinRequest, chatJoinRequest)
	assert.True(t, chatJoinRequestCalled)
}

func TestContextData(t *testing.T) {
	app := NewApplicationBuilder().
		Token("test-token").
		Build()

	msg := testMessage("hello", false)
	ctx := NewContext(app.Bot, Update{Message: msg})

	// Test BotData
	botData := ctx.BotData()
	assert.NotNil(t, botData)
	botData["global_config"] = "active"

	// Verify it persists in another context
	ctx2 := NewContext(app.Bot, Update{Message: msg})
	assert.Equal(t, "active", ctx2.BotData()["global_config"])

	// Test UserData
	userData := ctx.UserData()
	assert.NotNil(t, userData)
	userData["name"] = "Alice"

	// Verify user data is linked to user ID
	assert.Equal(t, "Alice", ctx2.UserData()["name"])

	// User data with different user ID should be empty
	msgDiffUser := testMessage("hello", false)
	msgDiffUser.From.ID = 999
	ctxDiffUser := NewContext(app.Bot, Update{Message: msgDiffUser})
	assert.Nil(t, ctxDiffUser.UserData()["name"])

	// Test ChatData
	chatData := ctx.ChatData()
	assert.NotNil(t, chatData)
	chatData["topic"] = "general"

	// Verify chat data persists
	assert.Equal(t, "general", ctx2.ChatData()["topic"])

	// Chat data with different chat ID should be empty
	msgDiffChat := testMessage("hello", false)
	msgDiffChat.Chat.ID = 888
	ctxDiffChat := NewContext(app.Bot, Update{Message: msgDiffChat})
	assert.Nil(t, ctxDiffChat.ChatData()["topic"])
}

func TestApplicationStop(t *testing.T) {
	app := NewApplicationBuilder().
		Token("test-token").
		Build()

	// Run polling and stop it immediately
	go func() {
		time.Sleep(100 * time.Millisecond)
		app.Stop()
	}()

	err := app.RunPolling()
	assert.NoError(t, err)
}

func TestApplicationRunWebhook(t *testing.T) {
	app := NewApplicationBuilder().
		Token("test-token").
		Build()

	// Simple context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(100 * time.Millisecond)
		app.Stop()
	}()

	// SetupWebhook will hit the API URL, which might fail or block in tests, so we just check it returns
	err := app.RunWebhook(ctx, "https://invalid-domain.local/webhook")
	assert.Error(t, err) // Should error because the domain doesn't exist or API fails
}
