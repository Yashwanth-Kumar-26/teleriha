package bot

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/store"
	"github.com/stretchr/testify/assert"
)

// seedStore returns a fresh MemoryStore pre-loaded with a serialised
// Conversation for user 42.  Returns the store and the raw JSON so
// tests can assert on exact bytes.
func seedStore(t *testing.T, conv *Conversation) *store.MemoryStore {
	t.Helper()
	s := store.NewMemoryStore()
	data, err := json.Marshal(conv)
	assert.NoError(t, err)
	err = s.Set("conv:42", data, 24*time.Hour)
	assert.NoError(t, err)
	return s
}

// freshCM returns a new ConversationManager and a MemoryStore pre-seeded
// with no conversation for user 42.
func freshCM(t *testing.T) (*ConversationManager, *store.MemoryStore) {
	t.Helper()
	return NewConversationManager(), store.NewMemoryStore()
}

func TestPersistedConversationManager_StartPersistsToStore(t *testing.T) {
	cm, s := freshCM(t)
	pcm := NewPersistedConversationManager(cm, s)

	pcm.Start(42, 42, "my_flow")

	data, err := s.Get("conv:42")
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// The stored blob should decode back to a Conversation.
	var decoded Conversation
	assert.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, int64(42), decoded.UserID)
	assert.Equal(t, int64(42), decoded.ChatID)
	assert.Equal(t, "my_flow", decoded.ID)
	assert.Equal(t, "start", decoded.State)
	assert.NotNil(t, decoded.Data)
}

func TestPersistedConversationManager_EndDeletesFromStore(t *testing.T) {
	cm, s := freshCM(t)
	pcm := NewPersistedConversationManager(cm, s)

	pcm.Start(42, 42, "my_flow")

	// Confirm it exists.
	_, err := s.Get("conv:42")
	assert.NoError(t, err)

	pcm.End(42)

	_, err = s.Get("conv:42")
	assert.ErrorIs(t, err, store.ErrKeyNotFound)
}

func TestPersistedConversationManager_GetHydratesFromStore(t *testing.T) {
	seeded := &Conversation{
		ID:     "my_flow",
		UserID: 42,
		ChatID: 99,
		State:  "step_b",
		Data:   map[string]interface{}{"name": "alice"},
	}
	s := seedStore(t, seeded)
	cm := NewConversationManager()
	pcm := NewPersistedConversationManager(cm, s)

	// In-memory map is empty; Get should hydrate from store.
	conv := pcm.Get(42)
	assert.NotNil(t, conv)
	assert.Equal(t, "my_flow", conv.ID)
	assert.Equal(t, int64(99), conv.ChatID)
	assert.Equal(t, "step_b", conv.State)
	assert.Equal(t, "alice", conv.Data["name"])

	// A second Get should come from in-memory now.
	conv2 := pcm.Get(42)
	assert.NotNil(t, conv2)
	assert.Equal(t, "step_b", conv2.State)
}

func TestPersistedConversationManager_HandleHydratesAndDelegates(t *testing.T) {
	seeded := &Conversation{
		ID:     "my_flow",
		UserID: 42,
		ChatID: 99,
		State:  "start",
		Data:   map[string]interface{}{},
	}
	s := seedStore(t, seeded)
	cm := NewConversationManager()
	pcm := NewPersistedConversationManager(cm, s)
	cm.Register("my_flow", func(ctx *Context, conv *Conversation) error {
		return nil
	})

	ctx := &Context{
		Update: Update{
			Message: &Message{
				From: &User{ID: 42},
				Chat: &Chat{ID: 99},
			},
		},
		Sender: &User{ID: 42},
		Chat:   &Chat{ID: 99},
	}

	handled, err := pcm.Handle(ctx)
	assert.NoError(t, err)
	assert.True(t, handled)
}

func TestPersistedConversationManager_GetReturnsNilWhenStoreMiss(t *testing.T) {
	_, s := freshCM(t)
	cm := NewConversationManager()
	pcm := NewPersistedConversationManager(cm, s)

	assert.Nil(t, pcm.Get(42))
}

func TestPersistedConversationManager_UpdateStatePersists(t *testing.T) {
	pcm := NewPersistedConversationManager(NewConversationManager(), store.NewMemoryStore())
	pcm.Start(42, 99, "flow")

	pcm.UpdateState(42, "step_two")

	data, err := pcm.store.Get("conv:42")
	assert.NoError(t, err)
	var conv Conversation
	assert.NoError(t, json.Unmarshal(data, &conv))
	assert.Equal(t, "step_two", conv.State)
}

func TestPersistedConversationManager_SetDataPersists(t *testing.T) {
	pcm := NewPersistedConversationManager(NewConversationManager(), store.NewMemoryStore())
	pcm.Start(42, 99, "flow")

	pcm.SetData(42, "answer", 42)

	data, err := pcm.store.Get("conv:42")
	assert.NoError(t, err)
	var conv Conversation
	assert.NoError(t, json.Unmarshal(data, &conv))
	assert.Equal(t, float64(42), conv.Data["answer"])
}

func TestPersistedConversationManager_ClearDataPersists(t *testing.T) {
	pcm := NewPersistedConversationManager(NewConversationManager(), store.NewMemoryStore())
	pcm.Start(42, 99, "flow")
	pcm.SetData(42, "answer", 42)

	pcm.ClearData(42)

	data, err := pcm.store.Get("conv:42")
	assert.NoError(t, err)
	var conv Conversation
	assert.NoError(t, json.Unmarshal(data, &conv))
	assert.Empty(t, conv.Data)
}

func TestPersistedConversationManager_NilStoreIsNoop(t *testing.T) {
	cm := NewConversationManager()
	pcm := NewPersistedConversationManager(cm, nil)

	// Start should not panic with nil store.
	pcm.Start(42, 99, "flow")
	// Get should not panic.
	assert.NotNil(t, pcm.Get(42))
}

func TestNewPersistedConversationManager_WrapsCM(t *testing.T) {
	cm := NewConversationManager()
	s := store.NewMemoryStore()
	pcm := NewPersistedConversationManager(cm, s)

	// The embedded CM should be the same instance we passed in.
	assert.Equal(t, cm, pcm.ConversationManager)
}

func TestStoreRegistry_StoresAndRetrievesConvos(t *testing.T) {
	// Each user gets its own persisted manager sharing one store —
	// simulates a real bot serving multiple users via Redis.
	s := store.NewMemoryStore()
	cm1 := NewConversationManager()
	cm2 := NewConversationManager()
	pcm1 := NewPersistedConversationManager(cm1, s)
	pcm2 := NewPersistedConversationManager(cm2, s)

	pcm1.Start(1, 1, "flow_a")
	pcm1.Start(2, 2, "flow_b")
	pcm2.Start(3, 3, "flow_c")

	c1 := pcm1.Get(1)
	assert.Equal(t, "flow_a", c1.ID)

	c2 := pcm2.Get(2)
	assert.Equal(t, "flow_b", c2.ID)

	c3 := pcm1.Get(3)
	assert.Equal(t, "flow_c", c3.ID)

	// End only one user; others remain.
	pcm1.End(1)
	assert.Nil(t, pcm1.Get(1))
	assert.NotNil(t, pcm2.Get(2))
	assert.NotNil(t, pcm1.Get(3))
}

func TestPersistedConversationManager_RoundTripEndToEnd(t *testing.T) {
	s := store.NewMemoryStore()
	cm := NewConversationManager()
	pcm := NewPersistedConversationManager(cm, s)

	// 1. Start → should be persisted
	pcm.Start(100, 200, "survey")
	c := pcm.Get(100)
	assert.NotNil(t, c)
	assert.Equal(t, "survey", c.ID)
	assert.Equal(t, int64(200), c.ChatID)
	assert.Equal(t, "start", c.State)

	// 2. mutate in-memory and persist
	c.State = "name"
	pcm.persistIfExists(100)

	// 3. Verify store reflects mutation
	data, err := s.Get("conv:100")
	assert.NoError(t, err)
	var fromStore Conversation
	assert.NoError(t, json.Unmarshal(data, &fromStore))
	assert.Equal(t, "name", fromStore.State)

	// 4. End → key should be gone
	pcm.End(100)
	_, err = s.Get("conv:100")
	assert.ErrorIs(t, err, store.ErrKeyNotFound)
}

func TestWithStore_BotOption(t *testing.T) {
	s := store.NewMemoryStore()
	bot := New("test-token", WithStore(s))
	assert.NotNil(t, bot.ConversationManager)
}
