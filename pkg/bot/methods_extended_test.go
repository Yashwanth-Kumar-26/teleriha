package bot

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBot_EditMessageCaption(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"caption":"new caption"}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.EditMessageCaption(123, 456, "new caption")
	assert.NoError(t, err)
	assert.Equal(t, "new caption", msg.Caption)
}

func TestBot_SendVideo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"video":{"file_id":"v1"}}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.SendVideo(123, "file_id", "caption")
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)
}

func TestBot_SendAudio(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"audio":{"file_id":"a1"}}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.SendAudio(123, "file_id", "caption")
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)
}

func TestBot_SendVoice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"voice":{"file_id":"vo1"}}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.SendVoice(123, "file_id", "caption")
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)
}

func TestBot_SendContact(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"contact":{"phone_number":"123","first_name":"n"}}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.SendContact(123, "123", "n")
	assert.NoError(t, err)
	assert.Equal(t, "123", msg.Contact.PhoneNumber)
}

func TestBot_SendDice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"dice":{"emoji":"🎲","value":6}}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	msg, err := b.SendDice(123, "🎲")
	assert.NoError(t, err)
	assert.Equal(t, 6, msg.Dice.Value)
}

func TestBot_GetChatMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":{"status":"member","user":{"id":123,"first_name":"n"}}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	member, err := b.GetChatMember(123, 456)
	assert.NoError(t, err)
	assert.Equal(t, "member", member.Status)
}

func TestBot_DeleteMyCommands(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	err := b.DeleteMyCommands()
	assert.NoError(t, err)
}

func TestBot_SendDocument_File(t *testing.T) {
	f, err := os.CreateTemp("", "test.txt")
	assert.NoError(t, err)
	defer os.Remove(f.Name())
	f.Write([]byte("test data"))
	f.Close()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":{"message_id":456,"document":{"file_id":"d1"}}}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	f2, _ := os.Open(f.Name())
	defer f2.Close()
	msg, err := b.SendDocument(123, f2, "caption")
	assert.NoError(t, err)
	assert.Equal(t, int64(456), msg.MessageID)
}

func TestBot_GetChatAdministrators(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ok":true,"result":[{"status":"creator","user":{"id":123,"first_name":"n"}}]}`)
	}))
	defer server.Close()

	b := New("token", WithBaseURL(server.URL))
	admins, err := b.GetChatAdministrators(123)
	assert.NoError(t, err)
	assert.Len(t, admins, 1)
	assert.Equal(t, "creator", admins[0].Status)
}
