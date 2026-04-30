package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

func TestBot_Integration_Start_Help(t *testing.T) {
	// Create a mock Telegram API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		// Handle getUpdates
		if r.URL.Path == "/bot-token/getUpdates" {
			resp := struct {
				Ok     bool         `json:"ok"`
				Result []bot.Update `json:"result"`
			}{
				Ok: true,
				Result: []bot.Update{
					{
						UpdateID: 1,
						Message: &bot.Message{
							MessageID: 101,
							Text:      "/start",
							Chat:      &bot.Chat{ID: 123, Type: "private"},
							From:      &bot.User{ID: 456, FirstName: "Tester"},
						},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Handle sendMessage
		if r.URL.Path == "/bot-token/sendMessage" {
			var params map[string]interface{}
			json.NewDecoder(r.Body).Decode(&params)
			
			assert.Equal(t, float64(123), params["chat_id"])
			assert.Contains(t, params["text"], "Hello Tester")
			
			resp := struct {
				Ok     bool         `json:"ok"`
				Result bot.Message `json:"result"`
			}{
				Ok: true,
				Result: bot.Message{
					MessageID: 102,
					Text:      params["text"].(string),
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		
		fmt.Fprintf(w, `{"ok":true,"result":[]}`)
	}))
	defer server.Close()

	// Create bot with mock server URL
	b := bot.New("-token", bot.WithBaseURL(server.URL))
	
	// Register handler
	b.Router.On("/start", func(ctx *bot.Context) error {
		_, err := ctx.Reply("Hello " + ctx.Sender.FirstName)
		return err
	})

	// Run polling for a short time
	err := b.StartPolling(10*time.Millisecond, 10*time.Millisecond)
	assert.NoError(t, err)

	// Wait for processing
	time.Sleep(100 * time.Millisecond)
	
	b.Stop()
}
