package bot

import "regexp"

// Filter represents a function that decides whether an update should be handled.
type Filter func(ctx *Context) bool

// FilterAnd returns true if all filters return true.
func FilterAnd(filters ...Filter) Filter {
	return func(ctx *Context) bool {
		for _, f := range filters {
			if !f(ctx) {
				return false
			}
		}
		return true
	}
}

// FilterOr returns true if any filter returns true.
func FilterOr(filters ...Filter) Filter {
	return func(ctx *Context) bool {
		for _, f := range filters {
			if f(ctx) {
				return true
			}
		}
		return false
	}
}

// FilterNot negates a filter.
func FilterNot(f Filter) Filter {
	return func(ctx *Context) bool {
		return !f(ctx)
	}
}

// FilterMessage ensures the update is a message.
func FilterMessage() Filter {
	return func(ctx *Context) bool {
		return ctx.Message != nil
	}
}

// FilterText ensures the message has text.
func FilterText() Filter {
	return func(ctx *Context) bool {
		return ctx.Message != nil && ctx.Message.Text != ""
	}
}

// FilterPhoto ensures the message has a photo.
func FilterPhoto() Filter {
	return func(ctx *Context) bool {
		return ctx.Message != nil && ctx.Message.Photo != nil
	}
}

// FilterVideo ensures the message has a video.
func FilterVideo() Filter {
	return func(ctx *Context) bool {
		return ctx.Message != nil && ctx.Message.Video != nil
	}
}

// FilterDocument ensures the message has a document.
func FilterDocument() Filter {
	return func(ctx *Context) bool {
		return ctx.Message != nil && ctx.Message.Document != nil
	}
}

// FilterSticker ensures the message has a sticker.
func FilterSticker() Filter {
	return func(ctx *Context) bool {
		return ctx.Message != nil && ctx.Message.Sticker != nil
	}
}

// FilterVoice ensures the message has a voice note.
func FilterVoice() Filter {
	return func(ctx *Context) bool {
		return ctx.Message != nil && ctx.Message.Voice != nil
	}
}

// FilterRegex returns true if the message text matches the pattern.
func FilterRegex(pattern string) Filter {
	re := regexp.MustCompile(pattern)
	return func(ctx *Context) bool {
		return ctx.Message != nil && re.MatchString(ctx.Message.Text)
	}
}

// FilterChatType filters messages based on chat type (e.g., "private", "group", "supergroup", "channel").
func FilterChatType(chatType string) Filter {
	return func(ctx *Context) bool {
		if ctx.Chat != nil {
			return ctx.Chat.Type == chatType
		}
		return false
	}
}

// FilterPrivate filters for private chats.
func FilterPrivate() Filter {
	return FilterChatType("private")
}

// FilterGroup filters for group or supergroup chats.
func FilterGroup() Filter {
	return FilterOr(FilterChatType("group"), FilterChatType("supergroup"))
}

// FilterForwarded filters for forwarded messages.
func FilterForwarded() Filter {
	return func(ctx *Context) bool {
		return ctx.Message != nil && ctx.Message.ForwardDate != 0
	}
}

// FilterReply filters for messages that are replying to another message.
func FilterReply() Filter {
	return func(ctx *Context) bool {
		return ctx.Message != nil && ctx.Message.ReplyToMessage != nil
	}
}
