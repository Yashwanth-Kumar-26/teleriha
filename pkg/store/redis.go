package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)


// RedisStore implements the Store interface using Redis.
type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisStore creates a new RedisStore.
func NewRedisStore(opts *redis.Options) *RedisStore {
	return &RedisStore{
		client: redis.NewClient(opts),
		ctx:    context.Background(),
	}
}

// Get retrieves a value from the store.
func (s *RedisStore) Get(key string) ([]byte, error) {
	val, err := s.client.Get(s.ctx, key).Bytes()
	if err == redis.Nil {
		return nil, ErrKeyNotFound
	}
	return val, err
}

// Set stores a value with an optional expiration time.
func (s *RedisStore) Set(key string, value []byte, ttl time.Duration) error {
	return s.client.Set(s.ctx, key, value, ttl).Err()
}

// Delete removes a key from the store.
func (s *RedisStore) Delete(key string) error {
	return s.client.Del(s.ctx, key).Err()
}

// Exists checks if a key exists and is not expired.
func (s *RedisStore) Exists(key string) bool {
	n, err := s.client.Exists(s.ctx, key).Result()
	return err == nil && n > 0
}

// Clear removes all keys from the store.
func (s *RedisStore) Clear() error {
	return s.client.FlushDB(s.ctx).Err()
}
