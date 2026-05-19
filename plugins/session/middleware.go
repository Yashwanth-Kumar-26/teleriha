package session

import (
	"context"
	"fmt"
	"time"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
	"github.com/Yashwanth-Kumar-26/teleriha/pkg/store"
)

// SessionData represents the data for a user's session.
type SessionData struct {
	store     *store.SessionStore
	sessionID string
}

// Get retrieves a value from the session.
func (s *SessionData) Get(key string) ([]byte, error) {
	return s.store.GetState(s.sessionID, key)
}

// Set stores a value in the session.
func (s *SessionData) Set(key string, value []byte) error {
	return s.store.SetState(s.sessionID, key, value)
}

type sessionKeyType string
const sessionContextKey sessionKeyType = "bot_session"

// Middleware creates a session middleware using the given store.
func Middleware(s store.Store, ttl time.Duration) bot.Middleware {
	sessionStore := store.NewSessionStore(s, ttl)

	return func(next bot.Handler) bot.Handler {
		return func(ctx *bot.Context) error {
			userID := ctx.SenderID()
			if userID == 0 {
				return next(ctx)
			}

			// Format session ID as user:12345
			sessionID := "user:" + fmt.Sprintf("%d", userID)

			sessionData := &SessionData{
				store:     sessionStore,
				sessionID: sessionID,
			}

			// Store the session in the context's internal Go context
			goCtx := context.WithValue(ctx.Context(), sessionContextKey, sessionData)
			ctxWithSession := ctx.WithContext(goCtx)

			return next(ctxWithSession)
		}
	}
}

// GetSession retrieves the session data from a context.
// Returns nil if no session is available (e.g., middleware not used).
func GetSession(ctx *bot.Context) *SessionData {
	val := ctx.Context().Value(sessionContextKey)
	if sessionData, ok := val.(*SessionData); ok {
		return sessionData
	}
	return nil
}
