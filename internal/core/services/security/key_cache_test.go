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
	"fmt"
	"testing"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain/security"
)

func TestKeyValidationCache_NewCache(t *testing.T) {
	cache := NewKeyValidationCache(10, 1*time.Hour)

	if cache == nil {
		t.Fatal("Expected cache to be created")
	}

	if cache.maxSize != 10 {
		t.Errorf("Expected maxSize to be 10, got %d", cache.maxSize)
	}

	if cache.ttl != 1*time.Hour {
		t.Errorf("Expected ttl to be 1 hour, got %v", cache.ttl)
	}

	if len(cache.cache) != 0 {
		t.Errorf("Expected empty cache, got %d entries", len(cache.cache))
	}
}

func TestKeyValidationCache_SetAndGet(t *testing.T) {
	cache := NewKeyValidationCache(10, 1*time.Hour)

	keyPath := "/path/to/key"
	result := &security.KeyValidationResult{
		IsValid:         true,
		Issues:          []string{},
		Warnings:        []string{},
		Recommendations: []string{},
	}

	// Set value
	cache.Set(keyPath, result)

	// Get value
	retrieved, exists := cache.Get(keyPath)
	if !exists {
		t.Fatal("Expected key to exist in cache")
	}

	if retrieved != result {
		t.Error("Expected retrieved result to match original")
	}

	if retrieved.IsValid != result.IsValid {
		t.Error("Expected IsValid to match")
	}
}

func TestKeyValidationCache_GetNonExistent(t *testing.T) {
	cache := NewKeyValidationCache(10, 1*time.Hour)

	_, exists := cache.Get("/non/existent/key")
	if exists {
		t.Error("Expected non-existent key to not exist in cache")
	}
}

func TestKeyValidationCache_Expiration(t *testing.T) {
	cache := NewKeyValidationCache(10, 100*time.Millisecond)

	keyPath := "/path/to/key"
	result := &security.KeyValidationResult{
		IsValid:         true,
		Issues:          []string{},
		Warnings:        []string{},
		Recommendations: []string{},
	}

	// Set value
	cache.Set(keyPath, result)

	// Should exist immediately
	_, exists := cache.Get(keyPath)
	if !exists {
		t.Error("Expected key to exist immediately after setting")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should not exist after expiration
	_, exists = cache.Get(keyPath)
	if exists {
		t.Error("Expected key to be expired")
	}
}

func TestKeyValidationCache_LRU(t *testing.T) {
	cache := NewKeyValidationCache(3, 1*time.Hour)

	// Add 3 items
	for i := 0; i < 3; i++ {
		keyPath := fmt.Sprintf("/key/%d", i)
		result := &security.KeyValidationResult{
			IsValid:         true,
			Issues:          []string{},
			Warnings:        []string{},
			Recommendations: []string{},
		}
		cache.Set(keyPath, result)
	}

	// Cache should be full
	if cache.Size() != 3 {
		t.Errorf("Expected cache size to be 3, got %d", cache.Size())
	}

	// Access first key to make it recently used
	cache.Get("/key/0")

	// Add 4th item, should evict least recently used (key/1)
	cache.Set("/key/3", &security.KeyValidationResult{
		IsValid:         true,
		Issues:          []string{},
		Warnings:        []string{},
		Recommendations: []string{},
	})

	// Cache should still be size 3
	if cache.Size() != 3 {
		t.Errorf("Expected cache size to be 3, got %d", cache.Size())
	}

	// key/0 should still exist (was accessed)
	_, exists := cache.Get("/key/0")
	if !exists {
		t.Error("Expected key/0 to still exist (was recently accessed)")
	}

	// key/1 should be evicted
	_, exists = cache.Get("/key/1")
	if exists {
		t.Error("Expected key/1 to be evicted (least recently used)")
	}

	// key/2 should still exist
	_, exists = cache.Get("/key/2")
	if !exists {
		t.Error("Expected key/2 to still exist")
	}

	// key/3 should exist (newly added)
	_, exists = cache.Get("/key/3")
	if !exists {
		t.Error("Expected key/3 to exist (newly added)")
	}
}

func TestKeyValidationCache_Clear(t *testing.T) {
	cache := NewKeyValidationCache(10, 1*time.Hour)

	// Add some items
	for i := 0; i < 5; i++ {
		keyPath := fmt.Sprintf("/key/%d", i)
		result := &security.KeyValidationResult{
			IsValid:         true,
			Issues:          []string{},
			Warnings:        []string{},
			Recommendations: []string{},
		}
		cache.Set(keyPath, result)
	}

	if cache.Size() != 5 {
		t.Errorf("Expected cache size to be 5, got %d", cache.Size())
	}

	// Clear cache
	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("Expected cache size to be 0 after clear, got %d", cache.Size())
	}

	// All keys should be gone
	for i := 0; i < 5; i++ {
		keyPath := fmt.Sprintf("/key/%d", i)
		_, exists := cache.Get(keyPath)
		if exists {
			t.Errorf("Expected key %s to be cleared", keyPath)
		}
	}
}

func TestKeyValidationCache_Size(t *testing.T) {
	cache := NewKeyValidationCache(10, 1*time.Hour)

	if cache.Size() != 0 {
		t.Errorf("Expected initial size to be 0, got %d", cache.Size())
	}

	// Add items
	for i := 0; i < 3; i++ {
		keyPath := fmt.Sprintf("/key/%d", i)
		result := &security.KeyValidationResult{
			IsValid:         true,
			Issues:          []string{},
			Warnings:        []string{},
			Recommendations: []string{},
		}
		cache.Set(keyPath, result)
	}

	if cache.Size() != 3 {
		t.Errorf("Expected size to be 3, got %d", cache.Size())
	}
}

func TestKeyValidationCache_Cleanup(t *testing.T) {
	cache := NewKeyValidationCache(10, 50*time.Millisecond)

	// Add items with different expiration times
	for i := 0; i < 3; i++ {
		keyPath := fmt.Sprintf("/key/%d", i)
		result := &security.KeyValidationResult{
			IsValid:         true,
			Issues:          []string{},
			Warnings:        []string{},
			Recommendations: []string{},
		}
		cache.Set(keyPath, result)
	}

	if cache.Size() != 3 {
		t.Errorf("Expected cache size to be 3, got %d", cache.Size())
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Cleanup should remove expired entries
	cache.Cleanup()

	if cache.Size() != 0 {
		t.Errorf("Expected cache size to be 0 after cleanup, got %d", cache.Size())
	}
}

func TestKeyValidationCache_Stats(t *testing.T) {
	cache := NewKeyValidationCache(10, 1*time.Hour)

	stats := cache.Stats()

	if stats["total_entries"] != 0 {
		t.Errorf("Expected total_entries to be 0, got %v", stats["total_entries"])
	}

	if stats["max_size"] != 10 {
		t.Errorf("Expected max_size to be 10, got %v", stats["max_size"])
	}

	if stats["ttl_seconds"] != 3600.0 {
		t.Errorf("Expected ttl_seconds to be 3600, got %v", stats["ttl_seconds"])
	}

	// Add some items
	for i := 0; i < 3; i++ {
		keyPath := fmt.Sprintf("/key/%d", i)
		result := &security.KeyValidationResult{
			IsValid:         true,
			Issues:          []string{},
			Warnings:        []string{},
			Recommendations: []string{},
		}
		cache.Set(keyPath, result)
	}

	stats = cache.Stats()

	if stats["total_entries"] != 3 {
		t.Errorf("Expected total_entries to be 3, got %v", stats["total_entries"])
	}

	if stats["expired_entries"] != 0 {
		t.Errorf("Expected expired_entries to be 0, got %v", stats["expired_entries"])
	}
}

func TestKeyValidationCache_ConcurrentAccess(t *testing.T) {
	cache := NewKeyValidationCache(100, 1*time.Hour)

	// Test concurrent writes
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				keyPath := fmt.Sprintf("/key/%d/%d", id, j)
				result := &security.KeyValidationResult{
					IsValid:         true,
					Issues:          []string{},
					Warnings:        []string{},
					Recommendations: []string{},
				}
				cache.Set(keyPath, result)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Cache should have some entries (may be less than 100 due to LRU eviction)
	if cache.Size() == 0 {
		t.Error("Expected cache to have some entries after concurrent writes")
	}

	// Test concurrent reads
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				keyPath := fmt.Sprintf("/key/%d/%d", id, j)
				cache.Get(keyPath)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// No panics should occur
}
