package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	b := New("token")
	assert.NotNil(t, b)
	assert.Equal(t, "token", b.Token)
	assert.Equal(t, "https://api.telegram.org", b.BaseURL)
}

func TestBotOptions(t *testing.T) {
	client := &http.Client{}
	b := New("token",
		WithBaseURL("https://example.com"),
		WithHTTPClient(client),
	)
	assert.Equal(t, "https://example.com", b.BaseURL)
	assert.Equal(t, client, b.Client)
}

func TestBot_GetMe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/getMe", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":{"id":123,"is_bot":true,"first_name":"TestBot"}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	user, err := b.GetMe()
	assert.NoError(t, err)
	assert.Equal(t, int64(123), user.ID)
	assert.Equal(t, "TestBot", user.FirstName)
}

func TestBot_SendMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/sendMessage", r.URL.Path)
		body, _ := io.ReadAll(r.Body)
		var params map[string]interface{}
		json.Unmarshal(body, &params)
		assert.Equal(t, float64(123), params["chat_id"])
		assert.Equal(t, "hello", params["text"])
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"text":"hello"}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.SendMessage(123, "hello")
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)
}

func TestBot_SetWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/setWebhook", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	err := b.SetWebhook("https://example.com/hook", 100, []string{"message"})
	assert.NoError(t, err)
}

func TestBot_DeleteWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/deleteWebhook", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	err := b.DeleteWebhook()
	assert.NoError(t, err)
}

func TestBot_Polling(t *testing.T) {
	handlerCalled := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/bottoken/getUpdates" {
			fmt.Fprint(w, `{"ok":true,"result":[{"update_id":1,"message":{"message_id":100,"text":"/ping"}}]}`)
		}
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	
	b.Router.On("/ping", func(ctx *Context) error {
		close(handlerCalled)
		return nil
	})

	err := b.StartPolling(10*time.Millisecond, 1*time.Second)
	assert.NoError(t, err)
	
	// Wait for handler to be called
	select {
	case <-handlerCalled:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for handler to be called")
	}

	b.Stop()
}

func TestBot_WebhookHandler(t *testing.T) {
	b := New("token")
	
	called := false
	b.Router.On("/ping", func(ctx *Context) error {
		called = true
		return nil
	})

	req := httptest.NewRequest("POST", "/hook", io.NopCloser(bytes.NewBufferString(`{"update_id":1,"message":{"message_id":100,"text":"/ping"}}`)))
	w := httptest.NewRecorder()
	
	b.handleWebhook(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, called)
}

func TestBot_BuildURL(t *testing.T) {
	b := New("token")
	u := b.buildURL("getMe")
	assert.Equal(t, "https://api.telegram.org/bottoken/getMe", u.String())
}

func TestBot_ProcessUpdate_CallbackQuery(t *testing.T) {
	b := New("token")
	called := false
	b.Router.OnCallback("click", func(ctx *Context) error {
		called = true
		return nil
	})

	b.processUpdate(Update{
		CallbackQuery: &CallbackQuery{
			Data: "click_here",
		},
	})
	assert.True(t, called)
}

func TestBot_ProcessUpdate_InlineQuery(t *testing.T) {
	b := New("token")
	called := false
	b.Router.OnInlineQuery(func(ctx *Context) error {
		called = true
		return nil
	})

	b.processUpdate(Update{
		InlineQuery: &InlineQuery{
			Query: "test",
		},
	})
	assert.True(t, called)
}

func TestBot_ProcessUpdate_ChosenInlineResult(t *testing.T) {
	b := New("token")
	called := false
	b.Router.OnChosenInlineResult(func(ctx *Context) error {
		called = true
		return nil
	})

	b.processUpdate(Update{
		ChosenInlineResult: &ChosenInlineResult{
			ResultID: "1",
		},
	})
	assert.True(t, called)
}

func TestBot_Builder(t *testing.T) {
	b := New("token")
	logger := b.Builder()
	assert.NotNil(t, logger)
}

func TestBot_CallMethod_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"ok":false,"description":"bad request","error_code":400}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	_, err := b.callMethod("test", nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code 400")
}

func TestBot_CallMethod_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":false,"description":"api error","error_code":500}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	_, err := b.callMethod("test", nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "telegram API error: api error")
}
