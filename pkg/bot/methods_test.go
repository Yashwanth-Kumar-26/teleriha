package bot

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBot_EditMessageText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/editMessageText", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"text":"edited"}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.EditMessageText(123, 456, "edited", WithNewParseMode("HTML"))
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)
	assert.Equal(t, "edited", msg.Text)
}

func TestBot_DeleteMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/deleteMessage", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	err := b.DeleteMessage(123, 456)
	assert.NoError(t, err)
}

func TestBot_AnswerCallbackQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/answerCallbackQuery", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	err := b.AnswerCallbackQuery("id", "answer", true)
	assert.NoError(t, err)
}

func TestBot_SendPhoto(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/sendPhoto", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"photo":[{"file_id":"p1"}]}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.SendPhoto(123, "file_id", "caption", WithPhotoCaption("new caption"))
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)
}

func TestBot_SendPhoto_File(t *testing.T) {
	// Create a temp file
	f, err := os.CreateTemp("", "test.png")
	assert.NoError(t, err)
	defer os.Remove(f.Name())
	f.Write([]byte("fake image data"))
	f.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/sendPhoto", r.URL.Path)
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"photo":[{"file_id":"p1"}]}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	
	f2, _ := os.Open(f.Name())
	defer f2.Close()
	
	msg, err := b.SendPhoto(123, f2, "caption")
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)
}

func TestBot_SendDocument(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/sendDocument", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"document":{"file_id":"d1"}}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.SendDocument(123, "file_id", "caption")
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)
}

func TestBot_SendSticker(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/sendSticker", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"sticker":{"file_id":"s1"}}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.SendSticker(123, "file_id")
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)
}

func TestBot_SendLocation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/sendLocation", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"location":{"latitude":1.0,"longitude":2.0}}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.SendLocation(123, 1.0, 2.0)
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)
}

func TestBot_SendPoll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/sendPoll", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"poll":{"id":"1","question":"q","options":[]}}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.SendPoll(123, "question", []string{"a", "b"})
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)
}

func TestBot_ForwardMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/forwardMessage", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":789}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.ForwardMessage(123, 456, 100)
	assert.NoError(t, err)
	assert.Equal(t, int64(789), msg.MessageID)
}

func TestBot_GetChat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/getChat", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":{"id":123,"type":"private"}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	chat, err := b.GetChat(123)
	assert.NoError(t, err)
	assert.Equal(t, int64(123), chat.ID)
}

func TestBot_SetMyCommands(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/setMyCommands", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	err := b.SetMyCommands([]BotCommand{{Command: "start", Description: "start"}})
	assert.NoError(t, err)
}

func TestBot_GetMyCommands(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/getMyCommands", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":[{"command":"start","description":"start"}]}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	cmds, err := b.GetMyCommands()
	assert.NoError(t, err)
	assert.Len(t, cmds, 1)
}

func TestBot_GetWebhookInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/getWebhookInfo", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":{"url":"https://example.com"}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	info, err := b.GetWebhookInfo()
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", info.URL)
}

func TestBot_AnswerInlineQuery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/bottoken/answerInlineQuery", r.URL.Path)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	err := b.AnswerInlineQuery("id", []InlineQueryResult{{Type: "article", ID: "1", Title: "t"}})
	assert.NoError(t, err)
}
