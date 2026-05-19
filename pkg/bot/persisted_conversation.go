package bot

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/store"
)

// PersistedConversationManager wraps a ConversationManager and adds
// transparent store-backed persistence for conversations.
type PersistedConversationManager struct {
	*ConversationManager
	store store.Store
	mu    sync.Mutex // serialises store writes per-manager
}

// NewPersistedConversationManager wraps cm with the given store.
// A nil store is replaced with a noop guard so callers never need a nil-check.
func NewPersistedConversationManager(cm *ConversationManager, s store.Store) *PersistedConversationManager {
	if s == nil {
		s = &noopStore{}
	}
	return &PersistedConversationManager{
		ConversationManager: cm,
		store:               s,
	}
}

// Start begins a new conversation and persists it.
func (pcm *PersistedConversationManager) Start(userID, chatID int64, id string) *Conversation {
	conv := pcm.ConversationManager.Start(userID, chatID, id)
	pcm.persist(userID, conv)
	return conv
}

// End ends the conversation and removes the persisted copy.
func (pcm *PersistedConversationManager) End(userID int64) {
	pcm.ConversationManager.End(userID)
	_ = pcm.store.Delete(convKey(userID))
}

// Get returns the conversation for userID, hydrating from the store on miss.
func (pcm *PersistedConversationManager) Get(userID int64) *Conversation {
	conv := pcm.ConversationManager.Get(userID)
	if conv != nil {
		return conv
	}
	return pcm.hydrate(userID)
}

// Handle processes a message in the context of a conversation, consulting
// the store for an in-memory miss before delegating upstream.
func (pcm *PersistedConversationManager) Handle(ctx *Context) (bool, error) {
	userID := ctx.SenderID()
	if userID == 0 {
		return false, nil
	}
	// Hydrate from store if not already in memory.
	pcm.Get(userID)
	return pcm.ConversationManager.Handle(ctx)
}

// UpdateState updates the conversation state and persists the change.
func (pcm *PersistedConversationManager) UpdateState(userID int64, state string) {
	pcm.ConversationManager.UpdateState(userID, state)
	pcm.persistIfExists(userID)
}

// SetData stores a key-value pair and persists the conversation.
func (pcm *PersistedConversationManager) SetData(userID int64, key string, value interface{}) {
	pcm.ConversationManager.SetData(userID, key, value)
	pcm.persistIfExists(userID)
}

// GetData delegates to the embedded manager; the full conversation is
// already hydrated by Get / Handle.
func (pcm *PersistedConversationManager) GetData(userID int64, key string) interface{} {
	return pcm.ConversationManager.GetData(userID, key)
}

// ClearData clears all conversation data and persists the change.
func (pcm *PersistedConversationManager) ClearData(userID int64) {
	pcm.ConversationManager.ClearData(userID)
	pcm.persistIfExists(userID)
}

// ---- internal helpers ----

func (pcm *PersistedConversationManager) persist(userID int64, conv *Conversation) {
	pcm.mu.Lock()
	defer pcm.mu.Unlock()

	data, err := json.Marshal(conv)
	if err != nil {
		return
	}
	_ = pcm.store.Set(convKey(userID), data, 24*time.Hour)
}

// persistIfExists serialises the current in-memory conversation for userID,
// if one exists, and writes it to the store.
func (pcm *PersistedConversationManager) persistIfExists(userID int64) {
	pcm.mu.Lock()
	defer pcm.mu.Unlock()

	conv := pcm.ConversationManager.Get(userID)
	if conv == nil {
		return
	}
	data, err := json.Marshal(conv)
	if err != nil {
		return
	}
	_ = pcm.store.Set(convKey(userID), data, 24*time.Hour)
}

func (pcm *PersistedConversationManager) hydrate(userID int64) *Conversation {
	pcm.mu.Lock()
	defer pcm.mu.Unlock()

	data, err := pcm.store.Get(convKey(userID))
	if err != nil || data == nil {
		return nil
	}

	var conv Conversation
	if err := json.Unmarshal(data, &conv); err != nil {
		return nil
	}

	pcm.ConversationManager.mu.Lock()
	pcm.conversations[userID] = &conv
	pcm.ConversationManager.mu.Unlock()

	return &conv
}

// convKey returns the store key for a given user ID. Format: "conv:{userID}".
func convKey(userID int64) string {
	return "conv:" + strconv.FormatInt(userID, 10)
}

// noopStore silently discards all operations; used as a nil-guard so callers
// never need a nil-check before calling store methods.
type noopStore struct{}

func (*noopStore) Get(string) ([]byte, error)              { return nil, errNotSet }
func (*noopStore) Set(string, []byte, time.Duration) error  { return errNotSet }
func (*noopStore) Delete(string) error                      { return errNotSet }
func (*noopStore) Exists(string) bool                       { return false }
func (*noopStore) Clear() error                             { return errNotSet }

var errNotSet = &storeNotSetError{}

type storeNotSetError struct{}

func (*storeNotSetError) Error() string { return "store not set" }
