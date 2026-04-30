package store

import (
	"sync"
	"time"
)

// MemoryStore is an in-memory key-value store for bot state.
// It's useful for development and testing, but for production
// you should use a persistent store like Redis.
type MemoryStore struct {
	// data stores all key-value pairs
	data map[string][]byte

	// expires stores expiration times for keys
	expires map[string]time.Time

	// mu protects the data and expires maps
	mu sync.RWMutex
}

// NewMemoryStore creates a new MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data:    make(map[string][]byte),
		expires: make(map[string]time.Time),
	}
}

// Get retrieves a value from the store.
func (s *MemoryStore) Get(key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if key is expired
	if exp, ok := s.expires[key]; ok {
		if time.Now().After(exp) {
			// Key is expired, but we can't delete it here
			// because we have a read lock. Just return nil.
			// The cleanup will happen on the next operation.
			return nil, ErrKeyNotFound
		}
	}

	val, ok := s.data[key]
	if !ok {
		return nil, ErrKeyNotFound
	}

	return val, nil
}

// Set stores a value with an optional expiration time.
func (s *MemoryStore) Set(key string, value []byte, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	if ttl > 0 {
		s.expires[key] = time.Now().Add(ttl)
	} else {
		delete(s.expires, key)
	}

	return nil
}

// Delete removes a key from the store.
func (s *MemoryStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	delete(s.expires, key)
	return nil
}

// Exists checks if a key exists and is not expired.
func (s *MemoryStore) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if key is expired
	if exp, ok := s.expires[key]; ok {
		if time.Now().After(exp) {
			return false
		}
	}

	_, ok := s.data[key]
	return ok
}

// Keys returns all keys in the store (including expired ones).
func (s *MemoryStore) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	return keys
}

// Clear removes all keys from the store.
func (s *MemoryStore) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string][]byte)
	s.expires = make(map[string]time.Time)
	return nil
}

// Cleanup removes all expired keys.
func (s *MemoryStore) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for key, exp := range s.expires {
		if now.After(exp) {
			delete(s.data, key)
			delete(s.expires, key)
		}
	}
}

// Size returns the number of keys in the store.
func (s *MemoryStore) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// (Moved to top)

// PrefixedMemoryStore is a MemoryStore that automatically adds a prefix to all keys.
type PrefixedMemoryStore struct {
	store  *MemoryStore
	prefix string
}

// NewPrefixedMemoryStore creates a new PrefixedMemoryStore.
func NewPrefixedMemoryStore(prefix string) *PrefixedMemoryStore {
	return &PrefixedMemoryStore{
		store:  NewMemoryStore(),
		prefix: prefix,
	}
}

// NewPrefixedMemoryStoreWithStore creates a new PrefixedMemoryStore with an existing store.
func NewPrefixedMemoryStoreWithStore(store *MemoryStore, prefix string) *PrefixedMemoryStore {
	return &PrefixedMemoryStore{
		store:  store,
		prefix: prefix,
	}
}

// Get retrieves a value from the store.
func (s *PrefixedMemoryStore) Get(key string) ([]byte, error) {
	return s.store.Get(s.fullKey(key))
}

// Set stores a value with an optional expiration time.
func (s *PrefixedMemoryStore) Set(key string, value []byte, ttl time.Duration) error {
	return s.store.Set(s.fullKey(key), value, ttl)
}

// Delete removes a key from the store.
func (s *PrefixedMemoryStore) Delete(key string) error {
	return s.store.Delete(s.fullKey(key))
}

// Exists checks if a key exists and is not expired.
func (s *PrefixedMemoryStore) Exists(key string) bool {
	return s.store.Exists(s.fullKey(key))
}

// Clear removes all keys with this prefix from the store.
func (s *PrefixedMemoryStore) Clear() error {
	// Get all keys and delete those with our prefix
	keys := s.store.Keys()
	for _, key := range keys {
		if len(key) >= len(s.prefix) && key[:len(s.prefix)] == s.prefix {
			s.store.Delete(key)
		}
	}
	return nil
}

// fullKey returns the full key with prefix.
func (s *PrefixedMemoryStore) fullKey(key string) string {
	return s.prefix + key
}
