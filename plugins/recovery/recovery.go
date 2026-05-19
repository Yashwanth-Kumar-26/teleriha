package recovery

import (
	"fmt"
	"runtime/debug"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

// Middleware returns a middleware that recovers from any panics
// and routes them to the bot's global ErrorHandler to prevent crashes.
func Middleware() bot.Middleware {
	return func(next bot.Handler) bot.Handler {
		return func(ctx *bot.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					// Convert panic to error
					err = fmt.Errorf("panic recovered: %v\n%s", r, string(debug.Stack()))
					ctx.Bot.Logger.Error().Msgf("Recovered from panic: %v", err)
				}
			}()
			
			return next(ctx)
		}
	}
}
