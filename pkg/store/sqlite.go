package store

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore implements the Store interface using SQLite.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates a new SQLiteStore with the given database path.
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS store (
			key TEXT PRIMARY KEY,
			value BLOB,
			expires_at DATETIME
		)
	`)
	if err != nil {
		return nil, err
	}

	// Start a background goroutine to clean up expired keys
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			_, _ = db.Exec(`DELETE FROM store WHERE expires_at IS NOT NULL AND expires_at < ?`, time.Now())
		}
	}()

	return &SQLiteStore{db: db}, nil
}

// Get retrieves a value from the store.
func (s *SQLiteStore) Get(key string) ([]byte, error) {
	var value []byte
	var expiresAt sql.NullTime

	err := s.db.QueryRow(`SELECT value, expires_at FROM store WHERE key = ?`, key).Scan(&value, &expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}

	if expiresAt.Valid && expiresAt.Time.Before(time.Now()) {
		_ = s.Delete(key)
		return nil, ErrKeyNotFound
	}

	return value, nil
}

// Set stores a value with an optional expiration time.
func (s *SQLiteStore) Set(key string, value []byte, ttl time.Duration) error {
	var expiresAt *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		expiresAt = &t
	}

	_, err := s.db.Exec(`
		INSERT INTO store (key, value, expires_at) 
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET 
			value = excluded.value, 
			expires_at = excluded.expires_at
	`, key, value, expiresAt)
	return err
}

// Delete removes a key from the store.
func (s *SQLiteStore) Delete(key string) error {
	_, err := s.db.Exec(`DELETE FROM store WHERE key = ?`, key)
	return err
}

// Exists checks if a key exists and is not expired.
func (s *SQLiteStore) Exists(key string) bool {
	var expiresAt sql.NullTime
	err := s.db.QueryRow(`SELECT expires_at FROM store WHERE key = ?`, key).Scan(&expiresAt)
	if err != nil {
		return false
	}
	if expiresAt.Valid && expiresAt.Time.Before(time.Now()) {
		_ = s.Delete(key)
		return false
	}
	return true
}

// Clear removes all keys from the store.
func (s *SQLiteStore) Clear() error {
	_, err := s.db.Exec(`DELETE FROM store`)
	return err
}

// Close closes the database connection.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
