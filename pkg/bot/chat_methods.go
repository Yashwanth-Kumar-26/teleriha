package bot

// BanChatMember bans a user from a chat.
func (b *Bot) BanChatMember(chatID int64, userID int64, opts ...MessageOption) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	}
	for _, opt := range opts { opt(params) }
	_, err := b.callMethod("banChatMember", params, nil)
	return err
}

// UnbanChatMember unbans a user from a chat.
func (b *Bot) UnbanChatMember(chatID int64, userID int64, opts ...MessageOption) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	}
	for _, opt := range opts { opt(params) }
	_, err := b.callMethod("unbanChatMember", params, nil)
	return err
}

// RestrictChatMember restricts a user in a supergroup.
func (b *Bot) RestrictChatMember(chatID int64, userID int64, permissions ChatPermissions, opts ...MessageOption) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
		"permissions": permissions,
	}
	for _, opt := range opts { opt(params) }
	_, err := b.callMethod("restrictChatMember", params, nil)
	return err
}

// PromoteChatMember promotes or demotes a user in a supergroup or channel.
func (b *Bot) PromoteChatMember(chatID int64, userID int64, opts ...MessageOption) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	}
	for _, opt := range opts { opt(params) }
	_, err := b.callMethod("promoteChatMember", params, nil)
	return err
}

// SetChatAdministratorCustomTitle sets a custom title for an administrator.
func (b *Bot) SetChatAdministratorCustomTitle(chatID int64, userID int64, customTitle string) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
		"custom_title": customTitle,
	}
	_, err := b.callMethod("setChatAdministratorCustomTitle", params, nil)
	return err
}

// SetChatPermissions sets default chat permissions for all members.
func (b *Bot) SetChatPermissions(chatID int64, permissions ChatPermissions, opts ...MessageOption) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"permissions": permissions,
	}
	for _, opt := range opts { opt(params) }
	_, err := b.callMethod("setChatPermissions", params, nil)
	return err
}

// ExportChatInviteLink generates a new primary invite link for a chat.
func (b *Bot) ExportChatInviteLink(chatID int64) (string, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}
	var link string
	_, err := b.callMethod("exportChatInviteLink", params, &link)
	return link, err
}

// CreateChatInviteLink creates an additional invite link for a chat.
func (b *Bot) CreateChatInviteLink(chatID int64, opts ...MessageOption) (*ChatInviteLink, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}
	for _, opt := range opts { opt(params) }
	var link ChatInviteLink
	_, err := b.callMethod("createChatInviteLink", params, &link)
	return &link, err
}

// EditChatInviteLink edits a non-primary invite link created by the bot.
func (b *Bot) EditChatInviteLink(chatID int64, inviteLink string, opts ...MessageOption) (*ChatInviteLink, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
		"invite_link": inviteLink,
	}
	for _, opt := range opts { opt(params) }
	var link ChatInviteLink
	_, err := b.callMethod("editChatInviteLink", params, &link)
	return &link, err
}

// RevokeChatInviteLink revokes an invite link created by the bot.
func (b *Bot) RevokeChatInviteLink(chatID int64, inviteLink string) (*ChatInviteLink, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
		"invite_link": inviteLink,
	}
	var link ChatInviteLink
	_, err := b.callMethod("revokeChatInviteLink", params, &link)
	return &link, err
}

// ApproveChatJoinRequest approves a chat join request.
func (b *Bot) ApproveChatJoinRequest(chatID int64, userID int64) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	}
	_, err := b.callMethod("approveChatJoinRequest", params, nil)
	return err
}

// DeclineChatJoinRequest declines a chat join request.
func (b *Bot) DeclineChatJoinRequest(chatID int64, userID int64) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	}
	_, err := b.callMethod("declineChatJoinRequest", params, nil)
	return err
}

// SetChatPhoto sets a new profile photo for the chat.
func (b *Bot) SetChatPhoto(chatID int64, photo interface{}) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"photo": photo,
	}
	_, err := b.callMethod("setChatPhoto", params, nil)
	return err
}

// DeleteChatPhoto deletes a chat photo.
func (b *Bot) DeleteChatPhoto(chatID int64) error {
	params := map[string]interface{}{
		"chat_id": chatID,
	}
	_, err := b.callMethod("deleteChatPhoto", params, nil)
	return err
}

// SetChatTitle changes the title of a chat.
func (b *Bot) SetChatTitle(chatID int64, title string) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"title": title,
	}
	_, err := b.callMethod("setChatTitle", params, nil)
	return err
}

// SetChatDescription changes the description of a chat.
func (b *Bot) SetChatDescription(chatID int64, description string) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"description": description,
	}
	_, err := b.callMethod("setChatDescription", params, nil)
	return err
}

// PinChatMessage adds a message to the list of pinned messages in a chat.
func (b *Bot) PinChatMessage(chatID int64, messageID int64, opts ...MessageOption) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"message_id": messageID,
	}
	for _, opt := range opts { opt(params) }
	_, err := b.callMethod("pinChatMessage", params, nil)
	return err
}

// UnpinChatMessage removes a message from the list of pinned messages.
func (b *Bot) UnpinChatMessage(chatID int64, messageID int64) error {
	params := map[string]interface{}{
		"chat_id": chatID,
		"message_id": messageID,
	}
	_, err := b.callMethod("unpinChatMessage", params, nil)
	return err
}

// UnpinAllChatMessages clears the list of pinned messages.
func (b *Bot) UnpinAllChatMessages(chatID int64) error {
	params := map[string]interface{}{
		"chat_id": chatID,
	}
	_, err := b.callMethod("unpinAllChatMessages", params, nil)
	return err
}

// LeaveChat lets the bot leave a group, supergroup or channel.
func (b *Bot) LeaveChat(chatID int64) error {
	params := map[string]interface{}{
		"chat_id": chatID,
	}
	_, err := b.callMethod("leaveChat", params, nil)
	return err
}

// GetChatMemberCount returns the number of members in a chat.
func (b *Bot) GetChatMemberCount(chatID int64) (int, error) {
	params := map[string]interface{}{
		"chat_id": chatID,
	}
	var count int
	_, err := b.callMethod("getChatMemberCount", params, &count)
	return count, err
}
