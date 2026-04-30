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
func (r *Router) HandleMessage(ctx *Context, message Message) {
	if message.Text == "" {
		// Handle non-text messages
		if r.messageHandlers[""] != nil {
			ctx.Message = &message
			r.messageHandlers[""](ctx)
		}
		return
	}

	text := message.Text

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
			if handler, ok := r.messageHandlers[cmd]; ok {
				ctx.Message = &message
				handler(ctx)
				return
			}
		}
	}

	// Try regex handlers
	for _, rh := range r.regexHandlers {
		if rh.pattern.MatchString(text) {
			ctx.Message = &message
			rh.handler(ctx)
			return
		}
	}

	// Try default handler
	if r.messageHandlers[""] != nil {
		ctx.Message = &message
		r.messageHandlers[""](ctx)
	}
}

// HandleCallbackQuery routes a callback query to the appropriate handler.
func (r *Router) HandleCallbackQuery(ctx *Context, callback CallbackQuery) {
	ctx.CallbackQuery = &callback

	// Try to find a prefix match
	for prefix, handler := range r.callbackHandlers {
		if strings.HasPrefix(callback.Data, prefix) {
			handler(ctx)
			return
		}
	}

	// try default callback handler if exists
	if handler, ok := r.callbackHandlers[""]; ok {
		handler(ctx)
	}
}

// HandleInlineQuery routes an inline query to the handler.
func (r *Router) HandleInlineQuery(ctx *Context, query InlineQuery) {
	ctx.InlineQuery = &query
	if r.inlineQueryHandler != nil {
		r.inlineQueryHandler(ctx)
	}
}

// HandleChosenInlineResult routes a chosen inline result to the handler.
func (r *Router) HandleChosenInlineResult(ctx *Context, result ChosenInlineResult) {
	ctx.ChosenInlineResult = &result
	if r.chosenInlineResultHandler != nil {
		r.chosenInlineResultHandler(ctx)
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
		middlewares:       r.middlewares,
		groupMiddlewares:  r.groupMiddlewares,
	}
}
