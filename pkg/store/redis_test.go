package store

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisStore(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)
	defer mr.Close()

	s := NewRedisStore(&redis.Options{
		Addr: mr.Addr(),
	})

	// Test Set and Get
	err = s.Set("key1", []byte("value1"), 1*time.Second)
	assert.NoError(t, err)

	val, err := s.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), val)

	// Test Exists
	assert.True(t, s.Exists("key1"))

	// Test Delete
	err = s.Delete("key1")
	assert.NoError(t, err)
	assert.False(t, s.Exists("key1"))

	// Test Expiration
	err = s.Set("key2", []byte("value2"), 10*time.Millisecond)
	assert.NoError(t, err)
	time.Sleep(20 * time.Millisecond)
	mr.FastForward(1 * time.Second) // Force expiration in miniredis

	_, err = s.Get("key2")
	assert.Error(t, err)

	// Test Clear
	s.Set("key3", []byte("value3"), 0)
	err = s.Clear()
	assert.NoError(t, err)
	assert.False(t, s.Exists("key3"))
}
