package bot

import (
	"context"
	"time"
)

// Application provides a high-level interface for building and running bots.
type Application struct {
	Bot      *Bot
	JobQueue *JobQueue
}

// ApplicationBuilder helps construct an Application.
type ApplicationBuilder struct {
	token string
	opts  []BotOption
}

// NewApplicationBuilder creates a new ApplicationBuilder.
func NewApplicationBuilder() *ApplicationBuilder {
	return &ApplicationBuilder{}
}

// Token sets the bot token.
func (b *ApplicationBuilder) Token(token string) *ApplicationBuilder {
	b.token = token
	return b
}

// Opts adds bot options.
func (b *ApplicationBuilder) Opts(opts ...BotOption) *ApplicationBuilder {
	b.opts = append(b.opts, opts...)
	return b
}

// Build creates the Application.
func (b *ApplicationBuilder) Build() *Application {
	bot := New(b.token, b.opts...)
	return &Application{
		Bot:      bot,
		JobQueue: NewJobQueue(bot),
	}
}

// HandlerRegistrar is a function that registers a handler on a Router.
type HandlerRegistrar func(*Router)

// AddHandler adds a handler to the application's router.
func (a *Application) AddHandler(registrar HandlerRegistrar) {
	registrar(a.Bot.Router)
}

// RunPolling starts the bot using long polling.
// It blocks until the bot is stopped.
func (a *Application) RunPolling() error {
	a.Bot.Logger.Info().Msg("Application starting polling...")
	a.JobQueue.Start()
	defer a.JobQueue.Stop()
	return a.Bot.StartPolling(1*time.Second, 30*time.Second)
}

// RunWebhook starts the bot using webhooks.
func (a *Application) RunWebhook(ctx context.Context, url string) error {
	a.Bot.Logger.Info().Msg("Application starting webhook...")
	a.JobQueue.Start()
	defer a.JobQueue.Stop()
	return a.Bot.SetupWebhook(ctx, url)
}

// Stop gracefully stops the application.
func (a *Application) Stop() {
	a.JobQueue.Stop()
	a.Bot.Stop()
}

// --- Handler registrars ---

// CommandHandler creates a command handler registrar.
func CommandHandler(command string, callback Handler) HandlerRegistrar {
	return func(r *Router) {
		r.On(command, callback)
	}
}

// MessageHandler creates a message handler registrar.
func MessageHandler(filter Filter, callback Handler) HandlerRegistrar {
	return func(r *Router) {
		r.OnFilter(filter, callback)
	}
}

// CallbackQueryHandler creates a callback query handler registrar.
func CallbackQueryHandler(prefix string, callback Handler) HandlerRegistrar {
	return func(r *Router) {
		r.OnCallback(prefix, callback)
	}
}

// InlineQueryHandler creates an inline query handler registrar.
func InlineQueryHandler(callback Handler) HandlerRegistrar {
	return func(r *Router) {
		r.OnInlineQuery(callback)
	}
}

// PreCheckoutQueryHandler creates a pre-checkout query handler registrar.
func PreCheckoutQueryHandler(callback Handler) HandlerRegistrar {
	return func(r *Router) {
		r.OnPreCheckoutQuery(callback)
	}
}

// ShippingQueryHandler creates a shipping query handler registrar.
func ShippingQueryHandler(callback Handler) HandlerRegistrar {
	return func(r *Router) {
		r.OnShippingQuery(callback)
	}
}

// PollHandler creates a poll handler registrar.
func PollHandler(callback Handler) HandlerRegistrar {
	return func(r *Router) {
		r.OnPoll(callback)
	}
}

// PollAnswerHandler creates a poll answer handler registrar.
func PollAnswerHandler(callback Handler) HandlerRegistrar {
	return func(r *Router) {
		r.OnPollAnswer(callback)
	}
}

// ChatMemberHandler creates a chat member handler registrar.
func ChatMemberHandler(callback Handler) HandlerRegistrar {
	return func(r *Router) {
		r.OnChatMember(callback)
	}
}

// MyChatMemberHandler creates a my chat member handler registrar.
func MyChatMemberHandler(callback Handler) HandlerRegistrar {
	return func(r *Router) {
		r.OnMyChatMember(callback)
	}
}

// ChatJoinRequestHandler creates a chat join request handler registrar.
func ChatJoinRequestHandler(callback Handler) HandlerRegistrar {
	return func(r *Router) {
		r.OnChatJoinRequest(callback)
	}
}
