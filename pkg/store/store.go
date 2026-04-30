// Package store provides interfaces and implementations for storing bot state.
package store

import (
	"errors"
	"fmt"
	"time"
)

// ErrKeyNotFound is returned when a key is not found in the store.
var ErrKeyNotFound = errors.New("key not found")

// Store is the interface that all stores must implement.
type Store interface {
	// Get retrieves a value from the store.
	// Returns ErrKeyNotFound if the key doesn't exist or is expired.
	Get(key string) ([]byte, error)

	// Set stores a value with an optional expiration time.
	// If ttl is 0, the key never expires.
	Set(key string, value []byte, ttl time.Duration) error

	// Delete removes a key from the store.
	Delete(key string) error

	// Exists checks if a key exists and is not expired.
	Exists(key string) bool

	// Clear removes all keys from the store.
	Clear() error
}

// StringStore is a helper interface for storing strings.
type StringStore interface {
	Store

	// GetString retrieves a string value from the store.
	GetString(key string) (string, error)

	// SetString stores a string value with an optional expiration time.
	SetString(key string, value string, ttl time.Duration) error
}

// stringStoreWrapper wraps a Store to provide string operations.
type stringStoreWrapper struct {
	Store
}

// NewStringStore wraps a Store to provide string operations.
func NewStringStore(store Store) StringStore {
	return &stringStoreWrapper{Store: store}
}

// GetString retrieves a string value from the store.
func (s *stringStoreWrapper) GetString(key string) (string, error) {
	data, err := s.Store.Get(key)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SetString stores a string value with an optional expiration time.
func (s *stringStoreWrapper) SetString(key string, value string, ttl time.Duration) error {
	return s.Store.Set(key, []byte(value), ttl)
}

// NamespaceStore provides a namespaced view of a store.
type NamespaceStore struct {
	store     Store
	namespace string
}

// NewNamespaceStore creates a new NamespaceStore.
func NewNamespaceStore(store Store, namespace string) *NamespaceStore {
	return &NamespaceStore{
		store:     store,
		namespace: namespace + ":",
	}
}

// Get retrieves a value from the namespaced store.
func (s *NamespaceStore) Get(key string) ([]byte, error) {
	return s.store.Get(s.fullKey(key))
}

// Set stores a value in the namespaced store.
func (s *NamespaceStore) Set(key string, value []byte, ttl time.Duration) error {
	return s.store.Set(s.fullKey(key), value, ttl)
}

// Delete removes a key from the namespaced store.
func (s *NamespaceStore) Delete(key string) error {
	return s.store.Delete(s.fullKey(key))
}

// Exists checks if a key exists in the namespaced store.
func (s *NamespaceStore) Exists(key string) bool {
	return s.store.Exists(s.fullKey(key))
}

// Clear removes all keys in this namespace from the store.
func (s *NamespaceStore) Clear() error {
	// This is a limitation: we can't efficiently clear all keys in a namespace
	// without scanning all keys. The implementation depends on the underlying store.
	// For now, we'll just return nil and let the caller handle it.
	// A proper implementation would need to track keys by namespace.
	return nil
}

// fullKey returns the full key with namespace.
func (s *NamespaceStore) fullKey(key string) string {
	return s.namespace + key
}

// UserStore provides methods for storing user-specific data.
type UserStore struct {
	store Store
}

// NewUserStore creates a new UserStore.
func NewUserStore(store Store) *UserStore {
	return &UserStore{store: store}
}

// Prefix for user keys.
const userPrefix = "user:"

// Get retrieves user data.
func (s *UserStore) Get(userID int64, key string) ([]byte, error) {
	return s.store.Get(userPrefix + keyForInt(userID) + ":" + key)
}

// Set stores user data.
func (s *UserStore) Set(userID int64, key string, value []byte, ttl time.Duration) error {
	return s.store.Set(userPrefix + keyForInt(userID) + ":" + key, value, ttl)
}

// Delete removes user data.
func (s *UserStore) Delete(userID int64, key string) error {
	return s.store.Delete(userPrefix + keyForInt(userID) + ":" + key)
}

// KeyForInt converts an int64 to a string key.
func keyForInt(i int64) string {
	return fmt.Sprintf("%d", i)
}

// ChatStore provides methods for storing chat-specific data.
type ChatStore struct {
	store Store
}

// NewChatStore creates a new ChatStore.
func NewChatStore(store Store) *ChatStore {
	return &ChatStore{store: store}
}

// Prefix for chat keys.
const chatPrefix = "chat:"

// Get retrieves chat data.
func (s *ChatStore) Get(chatID int64, key string) ([]byte, error) {
	return s.store.Get(chatPrefix + keyForInt(chatID) + ":" + key)
}

// Set stores chat data.
func (s *ChatStore) Set(chatID int64, key string, value []byte, ttl time.Duration) error {
	return s.store.Set(chatPrefix + keyForInt(chatID) + ":" + key, value, ttl)
}

// Delete removes chat data.
func (s *ChatStore) Delete(chatID int64, key string) error {
	return s.store.Delete(chatPrefix + keyForInt(chatID) + ":" + key)
}

// sessionStore provides methods for storing session data.
type SessionStore struct {
	store Store
	ttl   time.Duration
}

// NewSessionStore creates a new SessionStore with a default TTL.
func NewSessionStore(store Store, ttl time.Duration) *SessionStore {
	return &SessionStore{
		store: store,
		ttl:   ttl,
	}
}

// Prefix for session keys.
const sessionPrefix = "session:"

// Get retrieves session data.
func (s *SessionStore) Get(sessionID string) ([]byte, error) {
	return s.store.Get(sessionPrefix + sessionID)
}

// Set stores session data.
func (s *SessionStore) Set(sessionID string, value []byte) error {
	return s.store.Set(sessionPrefix+sessionID, value, s.ttl)
}

// Delete removes session data.
func (s *SessionStore) Delete(sessionID string) error {
	return s.store.Delete(sessionPrefix + sessionID)
}

// GetState retrieves a specific value from session data.
func (s *SessionStore) GetState(sessionID, key string) ([]byte, error) {
	return s.store.Get(sessionPrefix + sessionID + ":" + key)
}

// SetState stores a specific value in session data.
func (s *SessionStore) SetState(sessionID, key string, value []byte) error {
	return s.store.Set(sessionPrefix+sessionID+":"+key, value, s.ttl)
}
