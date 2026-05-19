// Package bot provides the core TeleRiHa bot framework for Telegram.
// It handles direct communication with the Telegram Bot API using net/http.
package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/Yashwanth-Kumar-26/teleriha/pkg/store"
)

// Bot represents the core Telegram bot instance.
// It manages the connection to Telegram's Bot API and routes updates.
type Bot struct {
	// Token is the bot token provided by @BotFather
	Token string

	// BaseURL is the Telegram Bot API base URL (defaults to api.telegram.org)
	BaseURL string

	// Client is the HTTP client used for requests
	Client *http.Client

	// Logger is used for structured logging
	Logger zerolog.Logger

	// Router handles incoming updates
	Router *Router

	// updateOffset tracks the last processed update ID for polling
	updateOffset int64

	// mu protects the updateOffset
	mu sync.Mutex

	// isRunning indicates if the bot is currently running
	isRunning bool

	// stopChan allows graceful shutdown; chan bool lets tests signal
	// without closing, avoiding races with stop-during-SetupWebhook.
	stopChan chan bool

	// plugins stores loaded plugins
	plugins map[string]Plugin

	// PluginRegistry is the plugin registry for this bot
	PluginRegistry *PluginRegistry

	// ConversationManager manages conversations
	ConversationManager *ConversationManager

	// Scoped data storage
	userData map[int64]map[string]interface{}
	chatData map[int64]map[string]interface{}
	botData  map[string]interface{}
	dataMu   sync.RWMutex
}

// BotOption is a function that configures a Bot.
type BotOption func(*Bot)

// WithBaseURL sets a custom base URL for the Telegram API.
func WithBaseURL(url string) BotOption {
	return func(b *Bot) {
		b.BaseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) BotOption {
	return func(b *Bot) {
		b.Client = client
	}
}

// WithLogger sets a custom logger.
func WithLogger(logger zerolog.Logger) BotOption {
	return func(b *Bot) {
		b.Logger = logger
	}
}

// WithStore attaches a store-backed persistent ConversationManager.
// Pass a *store.RedisStore for Redis-backed persistence or a
// *store.MemoryStore for in-process testing.
func WithStore(s store.Store) BotOption {
	return func(b *Bot) {
		b.ConversationManager = NewPersistedConversationManager(b.ConversationManager, s).ConversationManager
	}
}

// New creates a new Bot instance with the given token and options.
func New(token string, opts ...BotOption) *Bot {
	b := &Bot{
		Token:             token,
		BaseURL:           "https://api.telegram.org",
		Client:            &http.Client{Timeout: 30 * time.Second},
		Logger:            zerolog.Nop(),
		Router:            NewRouter(),
		updateOffset:      0,
		isRunning:         false,
		stopChan:          make(chan bool, 1),
		plugins:           make(map[string]Plugin),
		PluginRegistry:    NewPluginRegistry(),
		ConversationManager: NewConversationManager(),
		userData:          make(map[int64]map[string]interface{}),
		chatData:          make(map[int64]map[string]interface{}),
		botData:           make(map[string]interface{}),
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// Start begins processing updates using webhook mode.
// The handler is registered at the provided path.
func (b *Bot) StartWebhook(path string, port int) error {
	b.mu.Lock()
	if b.isRunning {
		b.mu.Unlock()
		return fmt.Errorf("bot is already running")
	}
	b.isRunning = true
	b.mu.Unlock()

	mux := http.NewServeMux()
	mux.HandleFunc(path, b.handleWebhook)

	addr := fmt.Sprintf(":%d", port)
	b.Logger.Info().Str("address", addr).Str("path", path).Msg("Starting webhook server")

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			b.Logger.Err(err).Msg("Webhook server error")
		}
	}()

	// Wait for stop signal
	<-b.stopChan

	b.Logger.Info().Msg("Shutting down webhook server")
	return server.Close()
}

// StartPolling begins processing updates using getUpdates (long polling).
func (b *Bot) StartPolling(pollInterval, timeout time.Duration) error {
	b.mu.Lock()
	if b.isRunning {
		b.mu.Unlock()
		return fmt.Errorf("bot is already running")
	}
	b.isRunning = true
	b.mu.Unlock()

	b.Logger.Info().Dur("interval", pollInterval).Dur("timeout", timeout).Msg("Starting polling mode")

	go func() {
		ticker := time.NewTicker(pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-b.stopChan:
				b.Logger.Info().Msg("Stopping polling mode")
				return
			case <-ticker.C:
				if err := b.pollUpdates(timeout); err != nil {
					b.Logger.Err(err).Msg("Error polling updates")
				}
			}
		}
	}()

	return nil
}

// Stop gracefully shuts down the bot.
func (b *Bot) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.isRunning {
		return
	}

	b.isRunning = false
	close(b.stopChan)
}

// handleWebhook processes incoming webhook updates.
func (b *Bot) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		b.Logger.Err(err).Msg("Failed to read request body")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	var update Update
	if err := json.Unmarshal(body, &update); err != nil {
		b.Logger.Err(err).Msg("Failed to unmarshal update")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	b.Logger.Debug().Interface("update", update).Msg("Received update")

	// Process the update
	b.processUpdate(update)

	w.WriteHeader(http.StatusOK)
}

// pollUpdates fetches updates from Telegram using long polling.
func (b *Bot) pollUpdates(timeout time.Duration) error {
	url := b.buildURL("getUpdates")

	params := url.Query()
	b.mu.Lock()
	params.Set("offset", fmt.Sprintf("%d", b.updateOffset+1))
	b.mu.Unlock()
	params.Set("timeout", fmt.Sprintf("%d", int(timeout.Seconds())))
	params.Set("allowed_updates", "message,callback_query,inline_query,chosen_inline_result,channel_post,edited_message,edited_channel_post,pre_checkout_query,shipping_query,poll,poll_answer,chat_member,my_chat_member,chat_join_request")
	url.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := b.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Ok {
		return fmt.Errorf("telegram API returned ok=false")
	}

	for _, update := range result.Result {
		b.processUpdate(update)

		b.mu.Lock()
		if update.UpdateID >= b.updateOffset {
			b.updateOffset = update.UpdateID
		}
		b.mu.Unlock()
	}

	return nil
}

// processUpdate routes the update to the appropriate handler.
func (b *Bot) processUpdate(update Update) {
	ctx := NewContext(b, update)

	// Handle different update types
	switch {
	case update.Message != nil:
		b.Router.HandleMessage(ctx, *update.Message)
	case update.CallbackQuery != nil:
		b.Router.HandleCallbackQuery(ctx, *update.CallbackQuery)
	case update.InlineQuery != nil:
		b.Router.HandleInlineQuery(ctx, *update.InlineQuery)
	case update.ChosenInlineResult != nil:
		b.Router.HandleChosenInlineResult(ctx, *update.ChosenInlineResult)
	case update.PreCheckoutQuery != nil:
		b.Router.HandlePreCheckoutQuery(ctx, *update.PreCheckoutQuery)
	case update.ShippingQuery != nil:
		b.Router.HandleShippingQuery(ctx, *update.ShippingQuery)
	case update.Poll != nil:
		b.Router.HandlePoll(ctx, *update.Poll)
	case update.PollAnswer != nil:
		b.Router.HandlePollAnswer(ctx, *update.PollAnswer)
	case update.ChatMember != nil:
		b.Router.HandleChatMember(ctx, *update.ChatMember)
	case update.MyChatMember != nil:
		b.Router.HandleMyChatMember(ctx, *update.MyChatMember)
	case update.ChatJoinRequest != nil:
		b.Router.HandleChatJoinRequest(ctx, *update.ChatJoinRequest)
	default:
		b.Logger.Warn().Interface("update", update).Msg("Received unknown update type")
	}
}

// SetWebhook sets up a webhook for the bot.
func (b *Bot) SetWebhook(url string, maxConnections int, allowedUpdates []string) error {
	params := map[string]interface{}{
		"url": url,
	}

	if maxConnections > 0 {
		params["max_connections"] = maxConnections
	}

	if len(allowedUpdates) > 0 {
		params["allowed_updates"] = allowedUpdates
	}

	_, err := b.callMethod("setWebhook", params, nil)
	return err
}

// DeleteWebhook removes the webhook.
func (b *Bot) DeleteWebhook() error {
	_, err := b.callMethod("deleteWebhook", nil, nil)
	return err
}

// SetupWebhook registers the given URL as a Telegram webhook and then starts the
// local HTTP server to receive updates. It combines SetWebhook and StartWebhook
// into a single call so callers can configure and launch a webhook-bound bot in
// one step.
// Returns ctx.Err() if the context is cancelled while the server is running.
func (b *Bot) SetupWebhook(ctx context.Context, url string) error {
	if err := b.SetWebhook(url, 100, nil); err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	// StartWebhook blocks until Stop() is called; run it in a goroutine so that
	// the context select below can detect cancellation.
	go func() {
		_ = b.StartWebhook("/webhook", 8080)
	}()

	// Wait for context cancellation and propagate it as the return error.
	<-ctx.Done()
	return ctx.Err()
}

// GetMe returns information about the bot.
func (b *Bot) GetMe() (*User, error) {
	var user User
	_, err := b.callMethod("getMe", nil, &user)
	return &user, err
}

// callMethod makes a generic Telegram API call with retry logic.
func (b *Bot) callMethod(method string, params map[string]interface{}, result interface{}) ([]byte, error) {
	return b.callMethodWithRetry(method, params, result, 3)
}

// callMethodWithRetry attempts to call the API with retries for transient errors.
func (b *Bot) callMethodWithRetry(method string, params map[string]interface{}, result interface{}, maxRetries int) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		body, err := b.doCallMethod(method, params, result)
		if err == nil {
			return body, nil
		}

		lastErr = err

		// Check if error is retryable (429 Too Many Requests, 5xx errors)
		if !isRetryableError(err) {
			return nil, err
		}

		// Exponential backoff: 1s, 2s, 4s
		if attempt < maxRetries-1 {
			backoff := time.Duration(1<<attempt) * time.Second
			b.Logger.Debug().
				Err(err).
				Int("attempt", attempt+1).
				Dur("backoff", backoff).
				Str("method", method).
				Msg("Retrying API call")
			time.Sleep(backoff)
		}
	}

	return nil, lastErr
}

// isRetryableError determines if an error should trigger a retry.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// Retry on 429 (rate limit) or 5xx (server errors)
	return strings.Contains(errStr, "status code 429") ||
		strings.Contains(errStr, "status code 5") ||
		strings.Contains(errStr, "status code 503")
}

// doCallMethod performs the actual HTTP call.
func (b *Bot) doCallMethod(method string, params map[string]interface{}, result interface{}) ([]byte, error) {
	url := b.buildURL(method)

	// Add params to URL for GET requests, or to body for POST
	var bodyReader io.Reader
	contentType := "application/json"

	if len(params) > 0 {
		isMultipart := false
		for _, v := range params {
			if _, ok := v.(io.Reader); ok {
				isMultipart = true
				break
			}
		}

		if isMultipart {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			for k, v := range params {
				if r, ok := v.(io.Reader); ok {
					filename := "file"
					if f, ok := v.(*os.File); ok {
						filename = filepath.Base(f.Name())
					}
					part, err := writer.CreateFormFile(k, filename)
					if err != nil {
						return nil, fmt.Errorf("failed to create form file: %w", err)
					}
					if _, err := io.Copy(part, r); err != nil {
						return nil, fmt.Errorf("failed to copy file data: %w", err)
					}
				} else {
					var valStr string
					switch val := v.(type) {
					case string:
						valStr = val
					case int, int64, float64, bool:
						valStr = fmt.Sprintf("%v", val)
					default:
						b, _ := json.Marshal(val)
						valStr = string(b)
					}
					if err := writer.WriteField(k, valStr); err != nil {
						return nil, fmt.Errorf("failed to write form field: %w", err)
					}
				}
			}
			if err := writer.Close(); err != nil {
				return nil, fmt.Errorf("failed to close multipart writer: %w", err)
			}
			bodyReader = body
			contentType = writer.FormDataContentType()
		} else {
			jsonBody, err := json.Marshal(params)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal params: %w", err)
			}
			bodyReader = bytes.NewReader(jsonBody)
		}
	}

	var req *http.Request
	var err error

	if bodyReader != nil {
		req, err = http.NewRequestWithContext(context.Background(), http.MethodPost, url.String(), bodyReader)
		if err == nil {
			req.Header.Set("Content-Type", contentType)
		}
	} else {
		req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, url.String(), nil)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	// Parse the response
	var apiResponse struct {
		Ok          bool            `json:"ok"`
		Result      json.RawMessage `json:"result"`
		Description string          `json:"description"`
		ErrorCode   int             `json:"error_code"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return body, fmt.Errorf("failed to parse API response: %w", err)
	}

	if !apiResponse.Ok {
		if apiResponse.Description != "" {
			return body, fmt.Errorf("telegram API error: %s (code: %d)", apiResponse.Description, apiResponse.ErrorCode)
		}
		return body, fmt.Errorf("telegram API returned ok=false")
	}

	// Unmarshal result if provided
	if result != nil && len(apiResponse.Result) > 0 {
		if err := json.Unmarshal(apiResponse.Result, result); err != nil {
			return body, fmt.Errorf("failed to unmarshal result: %w", err)
		}
	}

	return body, nil
}

// buildURL constructs the full URL for a Telegram API method.
func (b *Bot) buildURL(method string) *url.URL {
	base, _ := url.Parse(b.BaseURL)
	path := fmt.Sprintf("bot%s/%s", b.Token, method)
	return base.JoinPath(path)
}

// SendMessage sends a text message to a chat.
func (b *Bot) SendMessage(chatID int64, text string, opts ...MessageOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("sendMessage", params, &msg)
	return &msg, err
}

// MessageOption configures a message.
type MessageOption func(map[string]interface{})

// WithParseMode sets the parse mode for the message.
func WithParseMode(mode string) MessageOption {
	return func(params map[string]interface{}) {
		params["parse_mode"] = mode
	}
}

// WithReplyMarkup sets the reply markup for the message.
func WithReplyMarkup(markup interface{}) MessageOption {
	return func(params map[string]interface{}) {
		params["reply_markup"] = markup
	}
}

// WithDisableNotification disables notification for the message.
func WithDisableNotification() MessageOption {
	return func(params map[string]interface{}) {
		params["disable_notification"] = true
	}
}

// Bot returns the logger for the bot.
func (b *Bot) Builder() zerolog.Logger {
	return b.Logger
}
