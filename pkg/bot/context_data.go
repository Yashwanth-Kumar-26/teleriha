package bot

// UserData returns a thread-safe map for storing user-specific data.
func (c *Context) UserData() map[string]interface{} {
	if c.SenderID() == 0 {
		return nil
	}
	c.Bot.dataMu.Lock()
	defer c.Bot.dataMu.Unlock()
	
	if c.Bot.userData[c.SenderID()] == nil {
		c.Bot.userData[c.SenderID()] = make(map[string]interface{})
	}
	return c.Bot.userData[c.SenderID()]
}

// ChatData returns a thread-safe map for storing chat-specific data.
func (c *Context) ChatData() map[string]interface{} {
	if c.ChatID() == 0 {
		return nil
	}
	c.Bot.dataMu.Lock()
	defer c.Bot.dataMu.Unlock()
	
	if c.Bot.chatData[c.ChatID()] == nil {
		c.Bot.chatData[c.ChatID()] = make(map[string]interface{})
	}
	return c.Bot.chatData[c.ChatID()]
}

// BotData returns a thread-safe map for storing bot-global data.
func (c *Context) BotData() map[string]interface{} {
	c.Bot.dataMu.Lock()
	defer c.Bot.dataMu.Unlock()
	return c.Bot.botData
}
