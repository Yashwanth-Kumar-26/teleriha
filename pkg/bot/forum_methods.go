package bot

// CreateForumTopic creates a topic in a forum supergroup chat.
func (b *Bot) CreateForumTopic(chatID int64, name string, opts ...MessageOption) (*ForumTopic, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
		"name": name,
	}
	for _, opt := range opts { opt(params) }
	var topic ForumTopic
	_, err := b.callMethod("createForumTopic", params, &topic)
	return &topic, err
}

// EditForumTopic edits name and icon of a topic in a forum supergroup chat.
func (b *Bot) EditForumTopic(chatID int64, messageThreadID int, opts ...MessageOption) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"message_thread_id": messageThreadID,
	}
	for _, opt := range opts { opt(params) }
	_, err := b.callMethod("editForumTopic", params, nil)
	return err
}

// CloseForumTopic closes an open topic in a forum supergroup chat.
func (b *Bot) CloseForumTopic(chatID int64, messageThreadID int) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"message_thread_id": messageThreadID,
	}
	_, err := b.callMethod("closeForumTopic", params, nil)
	return err
}

// ReopenForumTopic reopens a closed topic in a forum supergroup chat.
func (b *Bot) ReopenForumTopic(chatID int64, messageThreadID int) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"message_thread_id": messageThreadID,
	}
	_, err := b.callMethod("reopenForumTopic", params, nil)
	return err
}

// DeleteForumTopic deletes a forum topic along with all its messages in a forum supergroup chat.
func (b *Bot) DeleteForumTopic(chatID int64, messageThreadID int) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"message_thread_id": messageThreadID,
	}
	_, err := b.callMethod("deleteForumTopic", params, nil)
	return err
}

// UnpinAllForumTopicMessages clears the list of pinned messages in a forum topic.
func (b *Bot) UnpinAllForumTopicMessages(chatID int64, messageThreadID int) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"message_thread_id": messageThreadID,
	}
	_, err := b.callMethod("unpinAllForumTopicMessages", params, nil)
	return err
}
