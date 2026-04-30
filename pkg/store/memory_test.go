package store

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test NewMemoryStore creates a store
func TestNewMemoryStore(t *testing.T) {
	store := NewMemoryStore()

	assert.NotNil(t, store)
	assert.NotNil(t, store.data)
	assert.NotNil(t, store.expires)
}

// Test Set and Get
func TestMemoryStore_SetGet(t *testing.T) {
	store := NewMemoryStore()

	value := []byte("test value")
	err := store.Set("key1", value, 0)
	assert.NoError(t, err)

	got, err := store.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, value, got)
}

// Test Get returns ErrKeyNotFound for non-existent key
func TestMemoryStore_Get_NotFound(t *testing.T) {
	store := NewMemoryStore()

	_, err := store.Get("nonexistent")
	assert.Equal(t, ErrKeyNotFound, err)
}

// Test Delete removes a key
func TestMemoryStore_Delete(t *testing.T) {
	store := NewMemoryStore()

	store.Set("key1", []byte("value"), 0)
	assert.True(t, store.Exists("key1"))

	err := store.Delete("key1")
	assert.NoError(t, err)

	assert.False(t, store.Exists("key1"))
}

// Test Delete non-existent key
func TestMemoryStore_Delete_NotFound(t *testing.T) {
	store := NewMemoryStore()

	// Should not error
	err := store.Delete("nonexistent")
	assert.NoError(t, err)
}

// Test Exists returns true for existing key
func TestMemoryStore_Exists(t *testing.T) {
	store := NewMemoryStore()

	store.Set("key1", []byte("value"), 0)

	assert.True(t, store.Exists("key1"))
}

// Test Exists returns false for non-existent key
func TestMemoryStore_Exists_NotFound(t *testing.T) {
	store := NewMemoryStore()

	assert.False(t, store.Exists("nonexistent"))
}

// Test Clear removes all keys
func TestMemoryStore_Clear(t *testing.T) {
	store := NewMemoryStore()

	store.Set("key1", []byte("value1"), 0)
	store.Set("key2", []byte("value2"), 0)
	store.Set("key3", []byte("value3"), 0)

	err := store.Clear()
	assert.NoError(t, err)

	assert.False(t, store.Exists("key1"))
	assert.False(t, store.Exists("key2"))
	assert.False(t, store.Exists("key3"))
}

// Test Keys returns all keys
func TestMemoryStore_Keys(t *testing.T) {
	store := NewMemoryStore()

	store.Set("key1", []byte("value1"), 0)
	store.Set("key2", []byte("value2"), 0)
	store.Set("key3", []byte("value3"), 0)

	keys := store.Keys()
	assert.Len(t, keys, 3)
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
	assert.Contains(t, keys, "key3")
}

// Test Size returns correct count
func TestMemoryStore_Size(t *testing.T) {
	store := NewMemoryStore()

	assert.Equal(t, 0, store.Size())

	store.Set("key1", []byte("value1"), 0)
	assert.Equal(t, 1, store.Size())

	store.Set("key2", []byte("value2"), 0)
	assert.Equal(t, 2, store.Size())

	store.Delete("key1")
	assert.Equal(t, 1, store.Size())
}

// Test Set with TTL
func TestMemoryStore_Set_WithTTL(t *testing.T) {
	store := NewMemoryStore()

	store.Set("key1", []byte("value"), 100*time.Millisecond)

	// Should exist immediately
	assert.True(t, store.Exists("key1"))

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	assert.False(t, store.Exists("key1"))
}

// Test Get with expired key
func TestMemoryStore_Get_Expired(t *testing.T) {
	store := NewMemoryStore()

	store.Set("key1", []byte("value"), 50*time.Millisecond)

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	_, err := store.Get("key1")
	assert.Equal(t, ErrKeyNotFound, err)
}

// Test Cleanup removes expired keys
func TestMemoryStore_Cleanup(t *testing.T) {
	store := NewMemoryStore()

	store.Set("key1", []byte("value1"), 50*time.Millisecond)
	store.Set("key2", []byte("value2"), 0) // No TTL

	// Wait for key1 to expire
	time.Sleep(100 * time.Millisecond)

	// key1 should still exist (expired but not cleaned)
	// This tests that Expires doesn't auto-clean, Cleanup does
	store.Cleanup()

	// Now key1 should be gone
	assert.False(t, store.Exists("key1"))
	// key2 should still exist
	assert.True(t, store.Exists("key2"))
}

// Test Set with zero TTL (never expires)
func TestMemoryStore_Set_ZeroTTL(t *testing.T) {
	store := NewMemoryStore()

	store.Set("key1", []byte("value"), 0)

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Should still exist
	assert.True(t, store.Exists("key1"))
}

// Test overwriting existing key
func TestMemoryStore_Overwrite(t *testing.T) {
	store := NewMemoryStore()

	store.Set("key1", []byte("value1"), 0)
	store.Set("key1", []byte("value2"), 0)

	got, err := store.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, []byte("value2"), got)
}

// Test concurrent access
func TestMemoryStore_Concurrent(t *testing.T) {
	store := NewMemoryStore()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			key := "key" + string(rune(idx))
			store.Set(key, []byte("value"), 0)
			store.Get(key)
			store.Set(key, []byte("value2"), 0)
			store.Delete(key)
		}(i)
	}

	wg.Wait()
	// Test passes if no race conditions detected (run with -race flag)
}

// Test concurrent read and write
func TestMemoryStore_ConcurrentRW(t *testing.T) {
	store := NewMemoryStore()

	// Pre-populate some keys
	for i := 0; i < 10; i++ {
		store.Set("key"+string(rune(i)), []byte("value"), 0)
	}

	var wg sync.WaitGroup

	// Readers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				store.Get("key" + string(rune(j%10)))
				store.Exists("key" + string(rune(j%10)))
			}
		}()
	}

	// Writers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				key := "key" + string(rune((idx+j)%10))
				store.Set(key, []byte("newvalue"), 0)
			}
		}(i)
	}

	wg.Wait()
}

// Test PrefixedMemoryStore
func TestNewPrefixedMemoryStore(t *testing.T) {
	store := NewPrefixedMemoryStore("test:")

	assert.NotNil(t, store)
	assert.NotNil(t, store.store)
	assert.Equal(t, "test:", store.prefix)
}

// Test PrefixedMemoryStore Get
func TestPrefixedMemoryStore_Get(t *testing.T) {
	store := NewPrefixedMemoryStore("test:")

	store.Set("key1", []byte("value"), 0)

	got, err := store.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, []byte("value"), got)
}

// Test PrefixedMemoryStore internal key structure
func TestPrefixedMemoryStore_KeyStructure(t *testing.T) {
	store := NewPrefixedMemoryStore("myprefix:")

	store.Set("key1", []byte("value"), 0)

	// The internal store should have the prefixed key
	// We can verify by checking if the key exists with prefix
	// Since we can't access internal store directly, we test via Get
	got, err := store.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, []byte("value"), got)

	// Setting same key without prefix in the underlying store should not conflict
	// (we create a new prefixed store to test this)
	underlying := NewMemoryStore()
	prefixed := NewPrefixedMemoryStoreWithStore(underlying, "prefix:")
	underlying.Set("key1", []byte("direct"), 0)
	prefixed.Set("key1", []byte("prefixed"), 0)

	// Getting from prefixed store should return prefixed value
	got, _ = prefixed.Get("key1")
	assert.Equal(t, []byte("prefixed"), got)

	// Getting from underlying store should return direct value
	got, _ = underlying.Get("key1")
	assert.Equal(t, []byte("direct"), got)

	// The prefixed key should be "prefix:key1" in underlying store
	got, _ = underlying.Get("prefix:key1")
	assert.Equal(t, []byte("prefixed"), got)
}

// Test PrefixedMemoryStore Delete
func TestPrefixedMemoryStore_Delete(t *testing.T) {
	store := NewPrefixedMemoryStore("test:")

	store.Set("key1", []byte("value"), 0)
	assert.True(t, store.Exists("key1"))

	store.Delete("key1")
	assert.False(t, store.Exists("key1"))
}

// Test PrefixedMemoryStore Exists
func TestPrefixedMemoryStore_Exists(t *testing.T) {
	store := NewPrefixedMemoryStore("test:")

	assert.False(t, store.Exists("key1"))

	store.Set("key1", []byte("value"), 0)
	assert.True(t, store.Exists("key1"))
}

// Test PrefixedMemoryStore Clear
func TestPrefixedMemoryStore_Clear(t *testing.T) {
	store := NewPrefixedMemoryStore("test:")

	store.Set("key1", []byte("value1"), 0)
	store.Set("key2", []byte("value2"), 0)

	// Add a key with different prefix to underlying store
	// We can't directly access the underlying store, but we can test that Clear works
	store.Clear()

	// Our keys should be gone
	assert.False(t, store.Exists("key1"))
	assert.False(t, store.Exists("key2"))
}

// Test StringStoreWrapper GetString
func TestStringStore_GetString(t *testing.T) {
	store := NewMemoryStore()
	stringStore := NewStringStore(store)

	stringStore.SetString("key1", "test value", 0)

	got, err := stringStore.GetString("key1")
	assert.NoError(t, err)
	assert.Equal(t, "test value", got)
}

// Test StringStoreWrapper SetString
func TestStringStore_SetString(t *testing.T) {
	store := NewMemoryStore()
	stringStore := NewStringStore(store)

	err := stringStore.SetString("key1", "test value", 0)
	assert.NoError(t, err)

	got, err := store.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, []byte("test value"), got)
}

// Test NamespaceStore
func TestNewNamespaceStore(t *testing.T) {
	store := NewMemoryStore()
	nsStore := NewNamespaceStore(store, "myns")

	assert.NotNil(t, nsStore)
	assert.NotNil(t, nsStore.store)
}

// Test NamespaceStore key prefixing
func TestNamespaceStore_KeyPrefixing(t *testing.T) {
	store := NewMemoryStore()
	nsStore := NewNamespaceStore(store, "user")

	nsStore.Set("123:name", []byte("John"), 0)

	// The key should be "user:123:name" in the underlying store
	got, err := store.Get("user:123:name")
	assert.NoError(t, err)
	assert.Equal(t, []byte("John"), got)

	// Getting via namespace store
	got, err = nsStore.Get("123:name")
	assert.NoError(t, err)
	assert.Equal(t, []byte("John"), got)
}

// Test UserStore
func TestNewUserStore(t *testing.T) {
	store := NewMemoryStore()
	userStore := NewUserStore(store)

	assert.NotNil(t, userStore)
}

// Test UserStore Set and Get
func TestUserStore_SetGet(t *testing.T) {
	store := NewMemoryStore()
	userStore := NewUserStore(store)

	err := userStore.Set(123, "name", []byte("John"), 0)
	assert.NoError(t, err)

	got, err := userStore.Get(123, "name")
	assert.NoError(t, err)
	assert.Equal(t, []byte("John"), got)
}

// Test UserStore Delete
func TestUserStore_Delete(t *testing.T) {
	store := NewMemoryStore()
	userStore := NewUserStore(store)

	userStore.Set(123, "name", []byte("John"), 0)
	err := userStore.Delete(123, "name")
	assert.NoError(t, err)

	_, err = userStore.Get(123, "name")
	assert.Equal(t, ErrKeyNotFound, err)
}

// Test UserStore key structure
func TestUserStore_KeyStructure(t *testing.T) {
	store := NewMemoryStore()
	userStore := NewUserStore(store)

	userStore.Set(123, "name", []byte("John"), 0)

	// The key should be "user:123:name" in the underlying store
	got, err := store.Get("user:123:name")
	assert.NoError(t, err)
	assert.Equal(t, []byte("John"), got)
}

// Test ChatStore
func TestNewChatStore(t *testing.T) {
	store := NewMemoryStore()
	chatStore := NewChatStore(store)

	assert.NotNil(t, chatStore)
}

// Test ChatStore Set and Get
func TestChatStore_SetGet(t *testing.T) {
	store := NewMemoryStore()
	chatStore := NewChatStore(store)

	err := chatStore.Set(123, "title", []byte("General"), 0)
	assert.NoError(t, err)

	got, err := chatStore.Get(123, "title")
	assert.NoError(t, err)
	assert.Equal(t, []byte("General"), got)
}

// Test ChatStore key structure
func TestChatStore_KeyStructure(t *testing.T) {
	store := NewMemoryStore()
	chatStore := NewChatStore(store)

	chatStore.Set(-1001234567, "title", []byte("Group Chat"), 0)

	// The key should be "chat:-1001234567:title" in the underlying store
	got, err := store.Get("chat:-1001234567:title")
	assert.NoError(t, err)
	assert.Equal(t, []byte("Group Chat"), got)
}

// Test SessionStore
func TestNewSessionStore(t *testing.T) {
	store := NewMemoryStore()
	SessionStore := NewSessionStore(store, 30*time.Minute)

	assert.NotNil(t, SessionStore)
	assert.Equal(t, 30*time.Minute, SessionStore.ttl)
}

// Test SessionStore Set and Get
func TestSessionStore_SetGet(t *testing.T) {
	store := NewMemoryStore()
	SessionStore := NewSessionStore(store, 0)

	err := SessionStore.Set("session123", []byte("data"))
	assert.NoError(t, err)

	got, err := SessionStore.Get("session123")
	assert.NoError(t, err)
	assert.Equal(t, []byte("data"), got)
}

// Test SessionStore SetState and GetState
func TestSessionStore_SetStateGetState(t *testing.T) {
	store := NewMemoryStore()
	SessionStore := NewSessionStore(store, 0)

	err := SessionStore.SetState("session123", "step", []byte("value"))
	assert.NoError(t, err)

	got, err := SessionStore.GetState("session123", "step")
	assert.NoError(t, err)
	assert.Equal(t, []byte("value"), got)
}

// Test SessionStore Delete
func TestSessionStore_Delete(t *testing.T) {
	store := NewMemoryStore()
	SessionStore := NewSessionStore(store, 0)

	SessionStore.Set("session123", []byte("data"))
	SessionStore.Delete("session123")

	_, err := SessionStore.Get("session123")
	assert.Equal(t, ErrKeyNotFound, err)
}

// Benchmark MemoryStore operations
func BenchmarkMemoryStore_Set(b *testing.B) {
	store := NewMemoryStore()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Set("key", []byte("value"), 0)
	}
}

func BenchmarkMemoryStore_Get(b *testing.B) {
	store := NewMemoryStore()
	store.Set("key", []byte("value"), 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Get("key")
	}
}

func BenchmarkMemoryStore_Delete(b *testing.B) {
	store := NewMemoryStore()
	store.Set("key", []byte("value"), 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Delete("key")
		store.Set("key", []byte("value"), 0)
	}
}

func BenchmarkMemoryStore_Exists(b *testing.B) {
	store := NewMemoryStore()
	store.Set("key", []byte("value"), 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Exists("key")
	}
}

func BenchmarkMemoryStore_Clear(b *testing.B) {
	store := NewMemoryStore()
	for i := 0; i < 100; i++ {
		store.Set("key"+string(rune(i)), []byte("value"), 0)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.Clear()
		for j := 0; j < 100; j++ {
			store.Set("key"+string(rune(j)), []byte("value"), 0)
		}
	}
}

func BenchmarkUserStore_SetGet(b *testing.B) {
	store := NewMemoryStore()
	userStore := NewUserStore(store)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userStore.Set(123, "key", []byte("value"), 0)
		userStore.Get(123, "key")
	}
}

func BenchmarkSessionStore_SetGet(b *testing.B) {
	store := NewMemoryStore()
	SessionStore := NewSessionStore(store, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SessionStore.Set("session", []byte("data"))
		SessionStore.Get("session")
	}
}
