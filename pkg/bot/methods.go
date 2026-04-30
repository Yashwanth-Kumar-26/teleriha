package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// EditMessageOption configures message editing.
type EditMessageOption func(map[string]interface{})

// WithNewParseMode sets the parse mode for edited message.
func WithNewParseMode(mode string) EditMessageOption {
	return func(params map[string]interface{}) {
		params["parse_mode"] = mode
	}
}

// WithNewReplyMarkup sets the reply markup for edited message.
func WithNewReplyMarkup(markup interface{}) EditMessageOption {
	return func(params map[string]interface{}) {
		params["reply_markup"] = markup
	}
}

// EditMessageText edits the text of a message.
func (b *Bot) EditMessageText(chatID int64, messageID int64, text string, opts ...EditMessageOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
		"text":       text,
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("editMessageText", params, &msg)
	return &msg, err
}

// EditMessageCaption edits the caption of a message.
func (b *Bot) EditMessageCaption(chatID int64, messageID int64, caption string, opts ...EditMessageOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
		"caption":    caption,
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("editMessageCaption", params, &msg)
	return &msg, err
}

// DeleteMessage deletes a message.
func (b *Bot) DeleteMessage(chatID int64, messageID int64) error {
	params := map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
	}

	_, err := b.callMethod("deleteMessage", params, nil)
	return err
}

// AnswerCallbackQuery sends an answer to a callback query.
func (b *Bot) AnswerCallbackQuery(callbackQueryID string, text string, showAlert bool) error {
	params := map[string]interface{}{
		"callback_query_id": callbackQueryID,
	}

	if text != "" {
		params["text"] = text
	}

	if showAlert {
		params["show_alert"] = true
	}

	_, err := b.callMethod("answerCallbackQuery", params, nil)
	return err
}

// SendPhotoOption configures photo sending.
type SendPhotoOption func(map[string]interface{})

// WithPhotoCaption sets the caption for the photo.
func WithPhotoCaption(caption string) SendPhotoOption {
	return func(params map[string]interface{}) {
		params["caption"] = caption
	}
}

// SendPhoto sends a photo to a chat.
func (b *Bot) SendPhoto(chatID int64, photo interface{}, caption string, opts ...SendPhotoOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}

	// Handle different types of photo input
	switch p := photo.(type) {
	case string:
		// Could be a file ID or URL
		if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
			params["photo"] = p
		} else {
			params["photo"] = p // Assume file ID
		}
	case *os.File, *File:
		// For file uploads, we need multipart
		return b.sendPhotoFile(chatID, p, caption, opts...)
	default:
		return nil, fmt.Errorf("unsupported photo type: %T", photo)
	}

	if caption != "" {
		params["caption"] = caption
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("sendPhoto", params, &msg)
	return &msg, err
}

// sendPhotoFile sends a photo file to a chat using multipart form.
func (b *Bot) sendPhotoFile(chatID int64, photo interface{}, caption string, opts ...SendPhotoOption) (*Message, error) {
	var fileBytes []byte
	var fileName string

	switch p := photo.(type) {
	case *os.File:
		var err error
		fileBytes, err = os.ReadFile(p.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		fileName = filepath.Base(p.Name())
	case *File:
		// Fork now, we assume File contains the file ID
		params := map[string]interface{}{
			"chat_id": chatID,
			"photo":   p.FileID,
		}
		if caption != "" {
			params["caption"] = caption
		}
		for _, opt := range opts {
			opt(params)
		}
		var msg Message
		_, err := b.callMethod("sendPhoto", params, &msg)
		return &msg, err
	default:
		return nil, fmt.Errorf("unsupported photo type for file upload: %T", photo)
	}

	// Create multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add photo field
	part, err := writer.CreateFormFile("photo", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(fileBytes); err != nil {
		return nil, fmt.Errorf("failed to write file data: %w", err)
	}

	// Add chat_id field
	if err := writer.WriteField("chat_id", fmt.Sprintf("%d", chatID)); err != nil {
		return nil, fmt.Errorf("failed to write chat_id field: %w", err)
	}

	// Add caption field
	if caption != "" {
		if err := writer.WriteField("caption", caption); err != nil {
			return nil, fmt.Errorf("failed to write caption field: %w", err)
		}
	}

	// Add fields from options
	for _, opt := range opts {
		opt(map[string]interface{}{})
		// Options are not applied here as we're using multipart
		// This is a simplification; in reality, we'd need to handle options differently
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	url := b.buildURL("sendPhoto")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url.String(), &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse struct {
		Ok     bool   `json:"ok"`
		Result Message `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram API returned ok=false")
	}

	return &apiResponse.Result, nil
}

// SendDocumentOption configures document sending.
type SendDocumentOption func(map[string]interface{})

// WithDocumentCaption sets the caption for the document.
func WithDocumentCaption(caption string) SendDocumentOption {
	return func(params map[string]interface{}) {
		params["caption"] = caption
	}
}

// SendDocument sends a document to a chat.
func (b *Bot) SendDocument(chatID int64, document interface{}, caption string, opts ...SendDocumentOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}

	switch d := document.(type) {
	case string:
		// Could be a file ID or URL
		params["document"] = d
	case *os.File:
		return b.sendDocumentFile(chatID, d, caption, opts...)
	default:
		return nil, fmt.Errorf("unsupported document type: %T", document)
	}

	if caption != "" {
		params["caption"] = caption
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("sendDocument", params, &msg)
	return &msg, err
}

// sendDocumentFile sends a document file to a chat using multipart form.
func (b *Bot) sendDocumentFile(chatID int64, document *os.File, caption string, opts ...SendDocumentOption) (*Message, error) {
	fileBytes, err := os.ReadFile(document.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	fileName := filepath.Base(document.Name())

	// Create multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add document field
	part, err := writer.CreateFormFile("document", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(fileBytes); err != nil {
		return nil, fmt.Errorf("failed to write file data: %w", err)
	}

	// Add chat_id field
	if err := writer.WriteField("chat_id", fmt.Sprintf("%d", chatID)); err != nil {
		return nil, fmt.Errorf("failed to write chat_id field: %w", err)
	}

	// Add caption field
	if caption != "" {
		if err := writer.WriteField("caption", caption); err != nil {
			return nil, fmt.Errorf("failed to write caption field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	url := b.buildURL("sendDocument")
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url.String(), &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse struct {
		Ok     bool   `json:"ok"`
		Result Message `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResponse.Ok {
		return nil, fmt.Errorf("telegram API returned ok=false")
	}

	return &apiResponse.Result, nil
}

// SendSticker sends a sticker to a chat.
func (b *Bot) SendSticker(chatID int64, sticker interface{}) (*Message, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}

	switch s := sticker.(type) {
	case string:
		params["sticker"] = s
	default:
		return nil, fmt.Errorf("unsupported sticker type: %T", sticker)
	}

	var msg Message
	_, err := b.callMethod("sendSticker", params, &msg)
	return &msg, err
}

// SendVideo sends a video to a chat.
func (b *Bot) SendVideo(chatID int64, video interface{}, caption string, opts ...MessageOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}

	switch v := video.(type) {
	case string:
		params["video"] = v
	default:
		return nil, fmt.Errorf("unsupported video type: %T", video)
	}

	if caption != "" {
		params["caption"] = caption
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("sendVideo", params, &msg)
	return &msg, err
}

// SendAudio sends an audio file to a chat.
func (b *Bot) SendAudio(chatID int64, audio interface{}, caption string, opts ...MessageOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}

	switch a := audio.(type) {
	case string:
		params["audio"] = a
	default:
		return nil, fmt.Errorf("unsupported audio type: %T", audio)
	}

	if caption != "" {
		params["caption"] = caption
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("sendAudio", params, &msg)
	return &msg, err
}

// SendVoice sends a voice message to a chat.
func (b *Bot) SendVoice(chatID int64, voice interface{}, caption string, opts ...MessageOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}

	switch v := voice.(type) {
	case string:
		params["voice"] = v
	default:
		return nil, fmt.Errorf("unsupported voice type: %T", voice)
	}

	if caption != "" {
		params["caption"] = caption
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("sendVoice", params, &msg)
	return &msg, err
}

// SendLocation sends a location to a chat.
func (b *Bot) SendLocation(chatID int64, latitude, longitude float64, opts ...MessageOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id":    chatID,
		"latitude":   latitude,
		"longitude":  longitude,
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("sendLocation", params, &msg)
	return &msg, err
}

// SendContact sends a contact to a chat.
func (b *Bot) SendContact(chatID int64, phoneNumber, firstName string, opts ...MessageOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id":      chatID,
		"phone_number": phoneNumber,
		"first_name":   firstName,
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("sendContact", params, &msg)
	return &msg, err
}

// SendPoll sends a poll to a chat.
func (b *Bot) SendPoll(chatID int64, question string, options []string, opts ...MessageOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id":   chatID,
		"question":  question,
		"options":   options,
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("sendPoll", params, &msg)
	return &msg, err
}

// SendDice sends a dice message to a chat.
func (b *Bot) SendDice(chatID int64, emoji string, opts ...MessageOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}

	if emoji != "" {
		params["emoji"] = emoji
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("sendDice", params, &msg)
	return &msg, err
}

// ForwardMessage forwards a message to a chat.
func (b *Bot) ForwardMessage(chatID int64, fromChatID int64, messageID int64, opts ...MessageOption) (*Message, error) {
	params := map[string]interface{}{
		"chat_id":      chatID,
		"from_chat_id": fromChatID,
		"message_id":   messageID,
	}

	for _, opt := range opts {
		opt(params)
	}

	var msg Message
	_, err := b.callMethod("forwardMessage", params, &msg)
	return &msg, err
}

// GetUpdates returns recent updates (for polling mode).
func (b *Bot) GetUpdates(offset, limit int, timeout int) ([]Update, error) {
	params := map[string]interface{}{
		"offset":  offset,
		"limit":   limit,
		"timeout": timeout,
	}

	var updates []Update
	_, err := b.callMethod("getUpdates", params, &updates)
	if err != nil {
		return nil, err
	}

	return updates, nil
}

// GetChat returns information about a chat.
func (b *Bot) GetChat(chatID int64) (*Chat, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}

	var chat Chat
	_, err := b.callMethod("getChat", params, &chat)
	return &chat, err
}

// GetChatAdministrators returns the administrators of a chat.
func (b *Bot) GetChatAdministrators(chatID int64) ([]ChatMember, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}

	var members []ChatMember
	_, err := b.callMethod("getChatAdministrators", params, &members)
	if err != nil {
		return nil, err
	}

	return members, nil
}

// GetChatMember returns a chat member.
func (b *Bot) GetChatMember(chatID int64, userID int64) (*ChatMember, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	}

	var member ChatMember
	_, err := b.callMethod("getChatMember", params, &member)
	return &member, err
}

// SetMyCommands sets the bot's commands.
func (b *Bot) SetMyCommands(commands []BotCommand) error {
	params := map[string]interface{}{
		"commands": commands,
	}

	_, err := b.callMethod("setMyCommands", params, nil)
	return err
}

// DeleteMyCommands deletes the bot's commands.
func (b *Bot) DeleteMyCommands() error {
	_, err := b.callMethod("deleteMyCommands", nil, nil)
	return err
}

// GetMyCommands returns the bot's commands.
func (b *Bot) GetMyCommands() ([]BotCommand, error) {
	var commands []BotCommand
	_, err := b.callMethod("getMyCommands", nil, &commands)
	if err != nil {
		return nil, err
	}

	return commands, nil
}

// BotCommand represents a bot command.
type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// WebhookInfo represents information about the current webhook.
type WebhookInfo struct {
	URL          string `json:"url,omitempty"`
	HasCustomCertificate bool   `json:"has_custom_certificate,omitempty"`
	PendingUpdateCount int    `json:"pending_update_count,omitempty"`
	IPAddress    string `json:"ip_address,omitempty"`
	LastErrorDate int64  `json:"last_error_date,omitempty"`
	LastErrorMessage string `json:"last_error_message,omitempty"`
	LastSynchronizationErrorDate int64 `json:"last_synchronization_error_date,omitempty"`
	MaxConnections int `json:"max_connections,omitempty"`
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}

// GetWebhookInfo returns information about the current webhook.
func (b *Bot) GetWebhookInfo() (*WebhookInfo, error) {
	var info WebhookInfo
	_, err := b.callMethod("getWebhookInfo", nil, &info)
	return &info, err
}

// AnswerInlineQuery answers an inline query.
func (b *Bot) AnswerInlineQuery(inlineQueryID string, results []InlineQueryResult) error {
	params := map[string]interface{}{
		"inline_query_id": inlineQueryID,
		"results":        results,
	}

	_, err := b.callMethod("answerInlineQuery", params, nil)
	return err
}

// InlineQueryResult represents an inline query result.
// This is a simplified version; Telegram has many result types.
type InlineQueryResult struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	MessageText string `json:"message_text,omitempty"`
	PhotoURL    string `json:"photo_url,omitempty"`
	ThumbURL    string `json:"thumb_url,omitempty"`
	InputMessageContent *InputMessageContent `json:"input_message_content,omitempty"`
	ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// InputMessageContent represents the content of a message to be sent.
type InputMessageContent struct {
	MessageText string `json:"message_text,omitempty"`
	ParseMode   string `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool `json:"disable_web_page_preview,omitempty"`
}
