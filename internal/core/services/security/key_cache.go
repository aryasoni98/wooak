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

package security

import (
	"sync"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain/security"
)

// KeyValidationResult represents a cached key validation result
type KeyValidationResult struct {
	Result    *security.KeyValidationResult
	Timestamp time.Time
	ExpiresAt time.Time
}

// KeyValidationCache provides LRU caching for SSH key validation results
type KeyValidationCache struct {
	cache       map[string]*KeyValidationResult
	mutex       sync.RWMutex
	maxSize     int
	ttl         time.Duration
	accessOrder []string // For LRU eviction
}

// NewKeyValidationCache creates a new key validation cache
func NewKeyValidationCache(maxSize int, ttl time.Duration) *KeyValidationCache {
	return &KeyValidationCache{
		cache:       make(map[string]*KeyValidationResult),
		maxSize:     maxSize,
		ttl:         ttl,
		accessOrder: make([]string, 0, maxSize),
	}
}

// Get retrieves a cached key validation result
func (c *KeyValidationCache) Get(keyPath string) (*security.KeyValidationResult, bool) {
	c.mutex.RLock()
	result, exists := c.cache[keyPath]
	if !exists {
		c.mutex.RUnlock()
		return nil, false
	}

	// Check if expired
	if time.Now().After(result.ExpiresAt) {
		c.mutex.RUnlock()
		return nil, false
	}

	// Need to upgrade to write lock for updating access order
	c.mutex.RUnlock()
	c.mutex.Lock()

	// Double-check the entry still exists and is not expired after acquiring write lock
	if result, exists := c.cache[keyPath]; exists && !time.Now().After(result.ExpiresAt) {
		// Update access order for LRU
		c.updateAccessOrder(keyPath)
		c.mutex.Unlock()
		return result.Result, true
	}

	c.mutex.Unlock()
	return nil, false
}

// Set stores a key validation result in the cache
func (c *KeyValidationCache) Set(keyPath string, result *security.KeyValidationResult) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	cacheEntry := &KeyValidationResult{
		Result:    result,
		Timestamp: now,
		ExpiresAt: now.Add(c.ttl),
	}

	// If key already exists, update it
	if _, exists := c.cache[keyPath]; exists {
		c.cache[keyPath] = cacheEntry
		c.updateAccessOrder(keyPath)
		return
	}

	// If cache is full, evict least recently used item
	if len(c.cache) >= c.maxSize {
		c.evictLRU()
	}

	// Add new entry
	c.cache[keyPath] = cacheEntry
	c.accessOrder = append(c.accessOrder, keyPath)
}

// Clear removes all entries from the cache
func (c *KeyValidationCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache = make(map[string]*KeyValidationResult)
	c.accessOrder = make([]string, 0, c.maxSize)
}

// Size returns the current number of cached entries
func (c *KeyValidationCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.cache)
}

// Cleanup removes expired entries from the cache
func (c *KeyValidationCache) Cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, result := range c.cache {
		if now.After(result.ExpiresAt) {
			delete(c.cache, key)
			c.removeFromAccessOrder(key)
		}
	}
}

// Stats returns cache statistics
func (c *KeyValidationCache) Stats() map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	now := time.Now()
	expiredCount := 0
	for _, result := range c.cache {
		if now.After(result.ExpiresAt) {
			expiredCount++
		}
	}

	return map[string]interface{}{
		"total_entries":   len(c.cache),
		"expired_entries": expiredCount,
		"max_size":        c.maxSize,
		"ttl_seconds":     c.ttl.Seconds(),
	}
}

// updateAccessOrder moves the key to the end of the access order (most recently used)
// Note: This method assumes the caller holds the write lock
func (c *KeyValidationCache) updateAccessOrder(key string) {
	// Remove from current position
	c.removeFromAccessOrder(key)
	// Add to end (most recently used)
	c.accessOrder = append(c.accessOrder, key)
}

// removeFromAccessOrder removes a key from the access order
// Note: This method assumes the caller holds the write lock
func (c *KeyValidationCache) removeFromAccessOrder(key string) {
	for i, k := range c.accessOrder {
		if k == key {
			c.accessOrder = append(c.accessOrder[:i], c.accessOrder[i+1:]...)
			break
		}
	}
}

// evictLRU removes the least recently used entry
func (c *KeyValidationCache) evictLRU() {
	if len(c.accessOrder) == 0 {
		return
	}

	// Remove the first (least recently used) entry
	lruKey := c.accessOrder[0]
	delete(c.cache, lruKey)
	c.accessOrder = c.accessOrder[1:]
}
