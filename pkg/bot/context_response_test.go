package bot

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext_ResponseMethods(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"text":"response"}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	ctx := NewContext(b, Update{
		Message: &Message{
			MessageID: 123,
			Chat:      &Chat{ID: 111},
			From:      &User{ID: 222},
		},
	})

	// Test Reply
	msg, err := ctx.Reply("hello")
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)

	// Test Send
	msg, err = ctx.Send(333, "hi")
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)

	// Test DeleteMessage
	err = ctx.DeleteMessage(456)
	assert.NoError(t, err)

	// Test Delete
	err = ctx.Delete()
	assert.NoError(t, err)

	// Test EditMessageText
	msg, err = ctx.EditMessageText("new text")
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)
}

func TestContext_AnswerCallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	ctx := NewContext(b, Update{
		CallbackQuery: &CallbackQuery{
			ID: "cb_id",
		},
	})

	err := ctx.AnswerCallback("done", false)
	assert.NoError(t, err)
}

func TestContext_Response_NoChatID(t *testing.T) {
	ctx := &Context{}
	_, err := ctx.Reply("fail")
	assert.Error(t, err)
	assert.Equal(t, "no chat ID in context", err.Error())
}

func TestContext_AnswerCallback_NoQuery(t *testing.T) {
	ctx := &Context{}
	err := ctx.AnswerCallback("fail", false)
	assert.Error(t, err)
	assert.Equal(t, "no callback query in context", err.Error())
}

func TestContext_Edit_NoMessage(t *testing.T) {
	ctx := &Context{}
	_, err := ctx.EditMessageText("fail")
	assert.Error(t, err)
	assert.Equal(t, "no message to edit in context", err.Error())
}

func TestContext_Delete_NoMessage(t *testing.T) {
	ctx := &Context{}
	err := ctx.Delete()
	assert.Error(t, err)
	assert.Equal(t, "no message to delete in context", err.Error())
}
