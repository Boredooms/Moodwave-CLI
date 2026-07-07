// Package cache provides a simple, thread-safe, disk-backed LRU cache
// for storing music metadata, API responses, and mood inference hints.
//
// The cache is designed to be very lightweight:
//   - In-memory LRU index with a configurable max entry count.
//   - Entries persisted as individual JSON files on disk.
//   - TTL-based expiry: stale entries are ignored and lazily cleaned.
//   - No background goroutine required.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Entry is a single cached item.
type Entry struct {
	Key       string          `json:"key"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
	ExpiresAt time.Time       `json:"expires_at"`
}

// IsExpired returns true if the entry's TTL has passed.
func (e *Entry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache is a thread-safe, LRU, disk-backed cache.
type Cache struct {
	mu         sync.Mutex
	dir        string
	maxEntries int
	defaultTTL time.Duration
	index      []string // LRU order: index[0] = oldest
}

// New creates a Cache backed by the given directory.
// maxEntries is the maximum number of entries to keep (LRU eviction).
// defaultTTL is the default entry lifetime.
func New(dir string, maxEntries int, defaultTTL time.Duration) (*Cache, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("cache: creating directory %s: %w", dir, err)
	}
	c := &Cache{
		dir:        dir,
		maxEntries: maxEntries,
		defaultTTL: defaultTTL,
	}
	c.buildIndex()
	return c, nil
}

// Get retrieves a cached value by key, unmarshalling it into dest.
// Returns (false, nil) if the key is not found or is expired.
func (c *Cache) Get(key string, dest interface{}) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	path := c.entryPath(key)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("cache: reading %s: %w", key, err)
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		// Corrupted entry — treat as miss and remove.
		_ = os.Remove(path)
		return false, nil
	}

	if entry.IsExpired() {
		_ = os.Remove(path)
		c.removeFromIndex(key)
		return false, nil
	}

	if err := json.Unmarshal(entry.Data, dest); err != nil {
		return false, fmt.Errorf("cache: unmarshalling %s: %w", key, err)
	}

	// Move to end of LRU index (most recently used).
	c.touchInIndex(key)

	return true, nil
}

// Set stores a value in the cache with the default TTL.
func (c *Cache) Set(key string, value interface{}) error {
	return c.SetTTL(key, value, c.defaultTTL)
}

// SetTTL stores a value in the cache with an explicit TTL.
func (c *Cache) SetTTL(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache: marshalling %s: %w", key, err)
	}

	now := time.Now()
	entry := Entry{
		Key:       key,
		Data:      data,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
	}

	entryData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("cache: marshalling entry %s: %w", key, err)
	}

	path := c.entryPath(key)
	if err := os.WriteFile(path, entryData, 0644); err != nil {
		return fmt.Errorf("cache: writing %s: %w", key, err)
	}

	c.addToIndex(key)
	c.evictIfNeeded()

	return nil
}

// Delete removes a key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_ = os.Remove(c.entryPath(key))
	c.removeFromIndex(key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return fmt.Errorf("cache: reading dir: %w", err)
	}
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			_ = os.Remove(filepath.Join(c.dir, e.Name()))
		}
	}
	c.index = nil
	return nil
}

// Stats returns basic cache statistics.
func (c *Cache) Stats() (int, int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.index), c.maxEntries
}

// entryPath returns the file path for a cache key.
// Keys are hashed to avoid filesystem path issues.
func (c *Cache) entryPath(key string) string {
	h := sha256.Sum256([]byte(key))
	name := hex.EncodeToString(h[:16]) + ".json"
	return filepath.Join(c.dir, name)
}

// buildIndex rebuilds the LRU index from files on disk.
func (c *Cache) buildIndex() {
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return
	}
	type entry struct {
		name      string
		createdAt time.Time
	}
	var valid []entry
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		valid = append(valid, entry{name: e.Name(), createdAt: info.ModTime()})
	}
	// Use filenames as index (already hashed keys).
	c.index = make([]string, 0, len(valid))
	for _, e := range valid {
		c.index = append(c.index, strings.TrimSuffix(e.name, ".json"))
	}
}

// addToIndex adds a key to the LRU index.
func (c *Cache) addToIndex(key string) {
	h := sha256.Sum256([]byte(key))
	id := hex.EncodeToString(h[:16])
	c.removeFromIndexByID(id)
	c.index = append(c.index, id)
}

// touchInIndex moves a key to the end of the LRU index.
func (c *Cache) touchInIndex(key string) {
	h := sha256.Sum256([]byte(key))
	id := hex.EncodeToString(h[:16])
	c.removeFromIndexByID(id)
	c.index = append(c.index, id)
}

// removeFromIndex removes a key from the LRU index.
func (c *Cache) removeFromIndex(key string) {
	h := sha256.Sum256([]byte(key))
	id := hex.EncodeToString(h[:16])
	c.removeFromIndexByID(id)
}

func (c *Cache) removeFromIndexByID(id string) {
	for i, k := range c.index {
		if k == id {
			c.index = append(c.index[:i], c.index[i+1:]...)
			return
		}
	}
}

// evictIfNeeded removes the oldest entries if we exceed maxEntries.
func (c *Cache) evictIfNeeded() {
	for len(c.index) > c.maxEntries {
		oldest := c.index[0]
		c.index = c.index[1:]
		_ = os.Remove(filepath.Join(c.dir, oldest+".json"))
	}
}
