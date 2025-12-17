// Copyright 2025.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ai

import (
	"sync"
	"time"
)

// CacheEntry represents a cached entry with expiration and access tracking
type CacheEntry struct {
	Value      interface{}
	ExpiresAt  time.Time
	LastAccess time.Time // For LRU eviction
}

// AICache provides caching functionality for AI responses with LRU eviction
type AICache struct {
	entries     map[string]*CacheEntry
	mutex       sync.RWMutex
	ttl         time.Duration
	maxSize     int // Maximum number of entries (0 = unlimited)
	stopChan    chan struct{}
	wg          sync.WaitGroup
	accessOrder []string // Track access order for LRU (most recent at end)
}

// CacheConfig configures an AI cache
type CacheConfig struct {
	TTL     time.Duration // Time-to-live for cache entries
	MaxSize int           // Maximum number of entries (0 = unlimited)
}

// DefaultCacheConfig returns a default cache configuration
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		TTL:     DefaultCacheTTL,
		MaxSize: DefaultCacheMaxSize, // Default to 1000 entries
	}
}

// NewAICache creates a new AI cache with default configuration
func NewAICache(ttl time.Duration) *AICache {
	return NewAICacheWithConfig(CacheConfig{
		TTL:     ttl,
		MaxSize: DefaultCacheMaxSize,
	})
}

// NewAICacheWithConfig creates a new AI cache with custom configuration
func NewAICacheWithConfig(config CacheConfig) *AICache {
	cache := &AICache{
		entries:     make(map[string]*CacheEntry),
		ttl:         config.TTL,
		maxSize:     config.MaxSize,
		stopChan:    make(chan struct{}),
		accessOrder: make([]string, 0),
	}

	// Start cleanup goroutine
	cache.wg.Add(1)
	go cache.cleanup()

	return cache
}

// Get retrieves a value from the cache and updates its access time for LRU
func (c *AICache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	now := time.Now()
	if now.After(entry.ExpiresAt) {
		// Remove expired entry
		c.removeEntryUnsafe(key)
		return nil, false
	}

	// Update access time and move to end of access order (most recently used)
	entry.LastAccess = now
	c.updateAccessOrderUnsafe(key)

	return entry.Value, true
}

// Set stores a value in the cache, evicting LRU entries if at capacity
func (c *AICache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()

	// If entry already exists, update it
	if _, exists := c.entries[key]; exists {
		c.entries[key].Value = value
		c.entries[key].ExpiresAt = now.Add(c.ttl)
		c.entries[key].LastAccess = now
		c.updateAccessOrderUnsafe(key)
		return
	}

	// Check if we need to evict entries
	if c.maxSize > 0 && len(c.entries) >= c.maxSize {
		c.evictLRUUnsafe()
	}

	// Add new entry
	c.entries[key] = &CacheEntry{
		Value:      value,
		ExpiresAt:  now.Add(c.ttl),
		LastAccess: now,
	}
	c.updateAccessOrderUnsafe(key)
}

// Delete removes a value from the cache
func (c *AICache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.removeEntryUnsafe(key)
}

// removeEntryUnsafe removes an entry without acquiring lock (caller must hold lock)
func (c *AICache) removeEntryUnsafe(key string) {
	delete(c.entries, key)
	// Remove from access order
	for i, k := range c.accessOrder {
		if k == key {
			c.accessOrder = append(c.accessOrder[:i], c.accessOrder[i+1:]...)
			break
		}
	}
}

// Clear removes all entries from the cache
func (c *AICache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries = make(map[string]*CacheEntry)
	c.accessOrder = make([]string, 0)
}

// Size returns the number of entries in the cache
func (c *AICache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.entries)
}

// MaxSize returns the maximum number of entries allowed in the cache
func (c *AICache) MaxSize() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.maxSize
}

// Stats returns cache statistics
func (c *AICache) Stats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return map[string]interface{}{
		"size":     len(c.entries),
		"max_size": c.maxSize,
		"ttl":      c.ttl.String(),
	}
}

// updateAccessOrderUnsafe updates the access order for LRU tracking (caller must hold lock)
func (c *AICache) updateAccessOrderUnsafe(key string) {
	// Remove key from current position if it exists
	for i, k := range c.accessOrder {
		if k == key {
			c.accessOrder = append(c.accessOrder[:i], c.accessOrder[i+1:]...)
			break
		}
	}
	// Add to end (most recently used)
	c.accessOrder = append(c.accessOrder, key)
}

// evictLRUUnsafe evicts the least recently used entry (caller must hold lock)
func (c *AICache) evictLRUUnsafe() {
	if len(c.accessOrder) == 0 {
		return
	}

	// Remove the first entry (least recently used)
	keyToEvict := c.accessOrder[0]
	c.removeEntryUnsafe(keyToEvict)
}

// cleanup removes expired entries from the cache
func (c *AICache) cleanup() {
	defer c.wg.Done()

	ticker := time.NewTicker(DefaultCacheCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mutex.Lock()
			now := time.Now()
			expiredKeys := make([]string, 0)
			for key, entry := range c.entries {
				if now.After(entry.ExpiresAt) {
					expiredKeys = append(expiredKeys, key)
				}
			}
			// Remove expired entries
			for _, key := range expiredKeys {
				c.removeEntryUnsafe(key)
			}
			c.mutex.Unlock()
		case <-c.stopChan:
			return // Proper exit when stopped
		}
	}
}

// Stop stops the cache cleanup goroutine
func (c *AICache) Stop() {
	// Check if already stopped without holding lock
	select {
	case <-c.stopChan:
		// Already stopped
		return
	default:
	}

	// Close stop channel while holding lock
	c.mutex.Lock()
	select {
	case <-c.stopChan:
		// Another goroutine already stopped it
		c.mutex.Unlock()
		return
	default:
		close(c.stopChan)
	}
	c.mutex.Unlock()

	// Wait for cleanup goroutine to finish (no lock needed here)
	c.wg.Wait()
}
