package bot

// SendMediaGroup sends a group of photos, videos, documents or audios as an album.
func (b *Bot) SendMediaGroup(chatID int64, media []interface{}, opts ...MessageOption) ([]Message, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
		"media": media,
	}
	for _, opt := range opts { opt(params) }
	var msgs []Message
	_, err := b.callMethod("sendMediaGroup", params, &msgs)
	return msgs, err
}
