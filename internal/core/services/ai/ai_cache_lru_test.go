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
	"testing"
	"time"
)

func TestAICache_LRU_Eviction(t *testing.T) {
	// Create cache with small max size for testing
	cache := NewAICacheWithConfig(CacheConfig{
		TTL:     1 * time.Hour,
		MaxSize: 3, // Small size to test eviction
	})
	defer cache.Stop()

	// Fill cache to capacity
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	if cache.Size() != 3 {
		t.Errorf("expected cache size 3, got %d", cache.Size())
	}

	// Access key1 to make it most recently used
	_, _ = cache.Get("key1")

	// Add a new entry - should evict least recently used (key2)
	cache.Set("key4", "value4")

	if cache.Size() != 3 {
		t.Errorf("expected cache size 3 after eviction, got %d", cache.Size())
	}

	// key2 should be evicted (least recently used)
	if _, exists := cache.Get("key2"); exists {
		t.Error("key2 should have been evicted")
	}

	// key1, key3, and key4 should still exist
	if _, exists := cache.Get("key1"); !exists {
		t.Error("key1 should still exist")
	}
	if _, exists := cache.Get("key3"); !exists {
		t.Error("key3 should still exist")
	}
	if _, exists := cache.Get("key4"); !exists {
		t.Error("key4 should still exist")
	}
}

func TestAICache_LRU_AccessOrder(t *testing.T) {
	cache := NewAICacheWithConfig(CacheConfig{
		TTL:     1 * time.Hour,
		MaxSize: 3,
	})
	defer cache.Stop()

	// Add entries
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Access key2 to make it most recently used
	_, _ = cache.Get("key2")

	// Add new entry - should evict key1 (least recently used)
	cache.Set("key4", "value4")

	// key1 should be evicted
	if _, exists := cache.Get("key1"); exists {
		t.Error("key1 should have been evicted (least recently used)")
	}

	// key2, key3, key4 should exist
	if _, exists := cache.Get("key2"); !exists {
		t.Error("key2 should exist (most recently used)")
	}
}

func TestAICache_UnlimitedSize(t *testing.T) {
	cache := NewAICacheWithConfig(CacheConfig{
		TTL:     1 * time.Hour,
		MaxSize: 0, // Unlimited
	})
	defer cache.Stop()

	// Add many entries
	for i := 0; i < 100; i++ {
		cache.Set(string(rune(i)), i)
	}

	if cache.Size() != 100 {
		t.Errorf("expected cache size 100, got %d", cache.Size())
	}

	// No eviction should occur
	if cache.MaxSize() != 0 {
		t.Error("cache should have unlimited size")
	}
}

func TestAICache_UpdateExistingEntry(t *testing.T) {
	cache := NewAICacheWithConfig(CacheConfig{
		TTL:     1 * time.Hour,
		MaxSize: 3,
	})
	defer cache.Stop()

	// Add entry
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Update existing entry - should not evict
	cache.Set("key1", "updated_value1")

	if cache.Size() != 3 {
		t.Errorf("expected cache size 3 after update, got %d", cache.Size())
	}

	// Verify updated value
	val, exists := cache.Get("key1")
	if !exists {
		t.Error("key1 should exist")
	}
	if val != "updated_value1" {
		t.Errorf("expected 'updated_value1', got %v", val)
	}
}

func TestAICache_Stats(t *testing.T) {
	cache := NewAICacheWithConfig(CacheConfig{
		TTL:     30 * time.Minute,
		MaxSize: 100,
	})
	defer cache.Stop()

	// Add some entries
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	stats := cache.Stats()
	if stats["size"].(int) != 2 {
		t.Errorf("expected size 2, got %d", stats["size"])
	}
	if stats["max_size"].(int) != 100 {
		t.Errorf("expected max_size 100, got %d", stats["max_size"])
	}
}

func TestAICache_LRU_MultipleEvictions(t *testing.T) {
	cache := NewAICacheWithConfig(CacheConfig{
		TTL:     1 * time.Hour,
		MaxSize: 2,
	})
	defer cache.Stop()

	// Fill cache
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// Access key1
	_, _ = cache.Get("key1")

	// Add key3 - should evict key2
	cache.Set("key3", "value3")

	if _, exists := cache.Get("key2"); exists {
		t.Error("key2 should have been evicted")
	}

	// Access key3
	_, _ = cache.Get("key3")

	// Add key4 - should evict key1 (least recently used now)
	cache.Set("key4", "value4")

	if _, exists := cache.Get("key1"); exists {
		t.Error("key1 should have been evicted")
	}

	// key3 and key4 should exist
	if _, exists := cache.Get("key3"); !exists {
		t.Error("key3 should exist")
	}
	if _, exists := cache.Get("key4"); !exists {
		t.Error("key4 should exist")
	}
}
