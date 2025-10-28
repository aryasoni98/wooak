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

// CacheEntry represents a cached entry with expiration
type CacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// AICache provides caching functionality for AI responses
type AICache struct {
	entries  map[string]*CacheEntry
	mutex    sync.RWMutex
	ttl      time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
}

// NewAICache creates a new AI cache
func NewAICache(ttl time.Duration) *AICache {
	cache := &AICache{
		entries:  make(map[string]*CacheEntry),
		ttl:      ttl,
		stopChan: make(chan struct{}),
	}

	// Start cleanup goroutine
	cache.wg.Add(1)
	go cache.cleanup()

	return cache
}

// Get retrieves a value from the cache
func (c *AICache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Value, true
}

// Set stores a value in the cache
func (c *AICache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = &CacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

// Delete removes a value from the cache
func (c *AICache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.entries, key)
}

// Clear removes all entries from the cache
func (c *AICache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries = make(map[string]*CacheEntry)
}

// Size returns the number of entries in the cache
func (c *AICache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.entries)
}

// cleanup removes expired entries from the cache
func (c *AICache) cleanup() {
	defer c.wg.Done()

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mutex.Lock()
			now := time.Now()
			for key, entry := range c.entries {
				if now.After(entry.ExpiresAt) {
					delete(c.entries, key)
				}
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
