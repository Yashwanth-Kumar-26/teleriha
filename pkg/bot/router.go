package bot

import (
	"fmt"
	"regexp"
	"strings"
)

// Handler is a function that handles a Telegram update.
type Handler func(*Context) error

// Middleware is a function that wraps a Handler.
type Middleware func(Handler) Handler

// Router handles routing of Telegram updates to the appropriate handlers.
type Router struct {
	// messageHandlers maps command names to handlers
	messageHandlers map[string]Handler

	// regexHandlers maps regex patterns to handlers
	regexHandlers []struct {
		pattern *regexp.Regexp
		handler Handler
	}

	// callbackHandlers maps callback data prefixes to handlers
	callbackHandlers map[string]Handler

	// inlineQueryHandlers handles inline queries
	inlineQueryHandler Handler

	// chosenInlineResultHandler handles chosen inline results
	chosenInlineResultHandler Handler

	// preCheckoutQueryHandler handles pre-checkout queries
	preCheckoutQueryHandler Handler

	// shippingQueryHandler handles shipping queries
	shippingQueryHandler Handler

	// successfulPaymentHandler handles successful payments
	successfulPaymentHandler Handler

	// pollHandler handles poll updates
	pollHandler Handler

	// pollAnswerHandler handles poll answer updates
	pollAnswerHandler Handler

	// chatMemberHandler handles chat member updates
	chatMemberHandler Handler

	// myChatMemberHandler handles my chat member updates
	myChatMemberHandler Handler

	// chatJoinRequestHandler handles chat join requests
	chatJoinRequestHandler Handler

	// filterHandlers handles advanced filter routing
	filterHandlers []struct {
		filter  Filter
		handler Handler
	}

	// errorHandler handles any errors returned by handlers
	errorHandler func(err error, ctx *Context)

	// middlewares is the list of global middlewares
	middlewares []Middleware

	// groupMiddlewares maps group names to their middlewares
	groupMiddlewares map[string][]Middleware
}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	return &Router{
		messageHandlers:    make(map[string]Handler),
		callbackHandlers:   make(map[string]Handler),
		groupMiddlewares:   make(map[string][]Middleware),
	}
}

// Use adds global middleware to the router.
func (r *Router) Use(middleware Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

// UseGroup adds middleware to a specific group.
func (r *Router) UseGroup(group string, middleware Middleware) {
	r.groupMiddlewares[group] = append(r.groupMiddlewares[group], middleware)
}

// Group creates a new group with optional middlewares.
func (r *Router) Group(group string, middlewares ...Middleware) *GroupRouter {
	return &GroupRouter{
		router:      r,
		group:       group,
		middlewares: middlewares,
	}
}

// GroupRouter is a router for a specific group.
type GroupRouter struct {
	router      *Router
	group       string
	middlewares []Middleware
}

// Use adds middleware to the group.
func (gr *GroupRouter) Use(middleware Middleware) {
	gr.middlewares = append(gr.middlewares, middleware)
}

// On registers a command handler for this group.
func (gr *GroupRouter) On(command string, handler Handler) {
	// Get all middlewares for this group
	var mws []Middleware
	if len(gr.middlewares) > 0 {
		mws = append(mws, gr.middlewares...)
	}
	if len(gr.router.groupMiddlewares[gr.group]) > 0 {
		mws = append(mws, gr.router.groupMiddlewares[gr.group]...)
	}

	// Apply middlewares
	finalHandler := handler
	for i := len(mws) - 1; i >= 0; i-- {
		finalHandler = mws[i](finalHandler)
	}

	// Add to router with group prefix
	fullCommand := gr.group + ":" + command
	gr.router.messageHandlers[fullCommand] = finalHandler
}

// On registers a handler for a specific command.
// Commands are case-insensitive and can include bot username.
// Example: /start, /start@mybot
func (r *Router) On(command string, handler Handler) {
	r.messageHandlers[strings.ToLower(command)] = r.applyMiddlewares(handler)
}

// OnRegex registers a handler for messages matching a regex pattern.
func (r *Router) OnRegex(pattern string, handler Handler) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}
	r.regexHandlers = append(r.regexHandlers, struct {
		pattern *regexp.Regexp
		handler Handler
	}{
		pattern: re,
		handler: r.applyMiddlewares(handler),
	})
	return nil
}

// OnCallback registers a handler for callback queries with a specific prefix.
func (r *Router) OnCallback(prefix string, handler Handler) {
	r.callbackHandlers[prefix] = r.applyMiddlewares(handler)
}

// OnInlineQuery registers a handler for inline queries.
func (r *Router) OnInlineQuery(handler Handler) {
	r.inlineQueryHandler = r.applyMiddlewares(handler)
}

// OnChosenInlineResult registers a handler for chosen inline results.
func (r *Router) OnChosenInlineResult(handler Handler) {
	r.chosenInlineResultHandler = r.applyMiddlewares(handler)
}

// OnPreCheckoutQuery registers a handler for pre-checkout queries.
func (r *Router) OnPreCheckoutQuery(handler Handler) {
	r.preCheckoutQueryHandler = handler
}

// OnFilter registers a message filter handler.
func (r *Router) OnFilter(filter Filter, handler Handler) {
	r.filterHandlers = append(r.filterHandlers, struct {
		filter  Filter
		handler Handler
	}{filter, handler})
}

// OnError registers a global error handler for the router.
func (r *Router) OnError(handler func(err error, ctx *Context)) {
	r.errorHandler = handler
}

// handleError safely executes the global error handler.
func (r *Router) handleError(err error, ctx *Context) {
	if r.errorHandler != nil {
		r.errorHandler(err, ctx)
	} else if ctx != nil && ctx.Bot != nil {
		ctx.Bot.Logger.Err(err).Msg("Unhandled handler error")
	} else {
		fmt.Println("Unhandled handler error:", err)
	}
}

// OnShippingQuery registers a handler for shipping queries.
func (r *Router) OnShippingQuery(handler Handler) {
	r.shippingQueryHandler = handler
}

// OnSuccessfulPayment registers a handler for successful payments.
func (r *Router) OnSuccessfulPayment(handler Handler) {
	r.successfulPaymentHandler = handler
}

// Default registers a default handler for messages that don't match any command.
func (r *Router) Default(handler Handler) {
	r.messageHandlers[""] = r.applyMiddlewares(handler)
}

// applyMiddlewares applies all global middlewares to a handler.
func (r *Router) applyMiddlewares(handler Handler) Handler {
	finalHandler := handler
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		finalHandler = r.middlewares[i](finalHandler)
	}
	return finalHandler
}

// HandleMessage routes a message to the appropriate handler.
func (r *Router) HandleMessage(ctx *Context, msg Message) {
	// Check for successful payment first
	if msg.SuccessfulPayment != nil && r.successfulPaymentHandler != nil {
		if err := r.successfulPaymentHandler(ctx); err != nil {
			r.handleError(err, ctx)
		}
		return
	}

	// Try message filter handlers first
	for _, fh := range r.filterHandlers {
		ctx.Message = &msg // temporarily set for filter
		if fh.filter(ctx) {
			if err := fh.handler(ctx); err != nil {
				r.handleError(err, ctx)
			}
			return
		}
	}

	if msg.Text == "" {
		// Handle non-text messages
		if r.messageHandlers[""] != nil {
			ctx.Message = &msg
			if err := r.messageHandlers[""](ctx); err != nil {
				r.handleError(err, ctx)
			}
		}
		return
	}

	text := msg.Text

	// Try command handlers
	if strings.HasPrefix(text, "/") {
		// Extract command name
		parts := strings.Fields(text)
		if len(parts) > 0 {
			cmd := strings.ToLower(parts[0])
			// Remove bot username if present
			if idx := strings.Index(cmd, "@"); idx != -1 {
				cmd = cmd[:idx]
			}

			// Try exact match
			if h, ok := r.messageHandlers[cmd]; ok {
				ctx.Message = &msg
				if err := h(ctx); err != nil {
					r.handleError(err, ctx)
				}
				return
			}
		}
	}

	// Try regex handlers
	for _, rh := range r.regexHandlers {
		if rh.pattern.MatchString(text) {
			ctx.Message = &msg
			if err := rh.handler(ctx); err != nil {
				r.handleError(err, ctx)
			}
			return
		}
	}

	// Try default handler
	if h, ok := r.messageHandlers[""]; ok {
		ctx.Message = &msg
		if err := h(ctx); err != nil {
			r.handleError(err, ctx)
		}
	}
}

// HandleCallbackQuery routes a callback query to the appropriate handler.
func (r *Router) HandleCallbackQuery(ctx *Context, callback CallbackQuery) {
	ctx.CallbackQuery = &callback

	// Try to find a prefix match
	for prefix, h := range r.callbackHandlers {
		if strings.HasPrefix(callback.Data, prefix) {
			if err := h(ctx); err != nil {
				r.handleError(err, ctx)
			}
			return
		}
	}

	// try default callback handler if exists
	if h, ok := r.callbackHandlers[""]; ok {
		if err := h(ctx); err != nil {
			r.handleError(err, ctx)
		}
	}
}

// HandleInlineQuery routes an inline query to the handler.
func (r *Router) HandleInlineQuery(ctx *Context, query InlineQuery) {
	ctx.InlineQuery = &query
	if r.inlineQueryHandler != nil {
		if err := r.inlineQueryHandler(ctx); err != nil {
			r.handleError(err, ctx)
		}
	}
}

// HandleChosenInlineResult routes a chosen inline result to the handler.
func (r *Router) HandleChosenInlineResult(ctx *Context, result ChosenInlineResult) {
	ctx.ChosenInlineResult = &result
	if r.chosenInlineResultHandler != nil {
		if err := r.chosenInlineResultHandler(ctx); err != nil {
			r.handleError(err, ctx)
		}
	}
}

// HandlePreCheckoutQuery routes a pre-checkout query to the handler.
func (r *Router) HandlePreCheckoutQuery(ctx *Context, q PreCheckoutQuery) {
	ctx.PreCheckoutQuery = &q
	if r.preCheckoutQueryHandler != nil {
		if err := r.preCheckoutQueryHandler(ctx); err != nil {
			r.handleError(err, ctx)
		}
	}
}

// HandleShippingQuery executes the registered handler for a shipping query.
func (r *Router) HandleShippingQuery(ctx *Context, q ShippingQuery) {
	ctx.ShippingQuery = &q
	if r.shippingQueryHandler != nil {
		if err := r.shippingQueryHandler(ctx); err != nil {
			r.handleError(err, ctx)
		}
	}
}

// OnPoll registers a handler for poll updates.
func (r *Router) OnPoll(handler Handler) {
	r.pollHandler = r.applyMiddlewares(handler)
}

// HandlePoll routes a poll update to the handler.
func (r *Router) HandlePoll(ctx *Context, poll Poll) {
	ctx.Update.Poll = &poll // Make sure it's in Update
	if r.pollHandler != nil {
		if err := r.pollHandler(ctx); err != nil {
			r.handleError(err, ctx)
		}
	}
}

// OnPollAnswer registers a handler for poll answer updates.
func (r *Router) OnPollAnswer(handler Handler) {
	r.pollAnswerHandler = r.applyMiddlewares(handler)
}

// HandlePollAnswer routes a poll answer update to the handler.
func (r *Router) HandlePollAnswer(ctx *Context, answer PollAnswer) {
	ctx.Update.PollAnswer = &answer
	if r.pollAnswerHandler != nil {
		if err := r.pollAnswerHandler(ctx); err != nil {
			r.handleError(err, ctx)
		}
	}
}

// OnChatMember registers a handler for chat member updates.
func (r *Router) OnChatMember(handler Handler) {
	r.chatMemberHandler = r.applyMiddlewares(handler)
}

// HandleChatMember routes a chat member update to the handler.
func (r *Router) HandleChatMember(ctx *Context, chatMember ChatMemberUpdated) {
	ctx.Update.ChatMember = &chatMember
	ctx.Chat = chatMember.Chat
	ctx.Sender = chatMember.From
	if r.chatMemberHandler != nil {
		if err := r.chatMemberHandler(ctx); err != nil {
			r.handleError(err, ctx)
		}
	}
}

// OnMyChatMember registers a handler for my chat member updates.
func (r *Router) OnMyChatMember(handler Handler) {
	r.myChatMemberHandler = r.applyMiddlewares(handler)
}

// HandleMyChatMember routes a my chat member update to the handler.
func (r *Router) HandleMyChatMember(ctx *Context, myChatMember ChatMemberUpdated) {
	ctx.Update.MyChatMember = &myChatMember
	ctx.Chat = myChatMember.Chat
	ctx.Sender = myChatMember.From
	if r.myChatMemberHandler != nil {
		if err := r.myChatMemberHandler(ctx); err != nil {
			r.handleError(err, ctx)
		}
	}
}

// OnChatJoinRequest registers a handler for chat join requests.
func (r *Router) OnChatJoinRequest(handler Handler) {
	r.chatJoinRequestHandler = r.applyMiddlewares(handler)
}

// HandleChatJoinRequest routes a chat join request to the handler.
func (r *Router) HandleChatJoinRequest(ctx *Context, req ChatJoinRequest) {
	ctx.Update.ChatJoinRequest = &req
	ctx.Chat = req.Chat
	ctx.Sender = req.From
	if r.chatJoinRequestHandler != nil {
		if err := r.chatJoinRequestHandler(ctx); err != nil {
			r.handleError(err, ctx)
		}
	}
}


// Clone creates a new router with the same state.
func (r *Router) Clone() *Router {
	return &Router{
		messageHandlers:   r.messageHandlers,
		regexHandlers:     r.regexHandlers,
		callbackHandlers:  r.callbackHandlers,
		inlineQueryHandler: r.inlineQueryHandler,
		chosenInlineResultHandler: r.chosenInlineResultHandler,
		preCheckoutQueryHandler:   r.preCheckoutQueryHandler,
		shippingQueryHandler:      r.shippingQueryHandler,
		successfulPaymentHandler:  r.successfulPaymentHandler,
		pollHandler:               r.pollHandler,
		pollAnswerHandler:         r.pollAnswerHandler,
		chatMemberHandler:         r.chatMemberHandler,
		myChatMemberHandler:       r.myChatMemberHandler,
		chatJoinRequestHandler:    r.chatJoinRequestHandler,
		filterHandlers:            r.filterHandlers,
		errorHandler:              r.errorHandler,
		middlewares:               r.middlewares,
		groupMiddlewares:          r.groupMiddlewares,
	}
}
