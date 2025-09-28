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
	"fmt"
	"testing"
	"time"
)

func TestAICache_NewAICache(t *testing.T) {
	ttl := 5 * time.Minute
	cache := NewAICache(ttl)

	if cache == nil {
		t.Fatal("Expected cache to be created, got nil")
	}

	if cache.ttl != ttl {
		t.Errorf("Expected TTL to be %v, got %v", ttl, cache.ttl)
	}

	if cache.entries == nil {
		t.Error("Expected entries map to be initialized")
	}

	if len(cache.entries) != 0 {
		t.Error("Expected initial entries to be empty")
	}
}

func TestAICache_SetAndGet(t *testing.T) {
	cache := NewAICache(1 * time.Minute)
	key := "test-key"
	value := "test-value"

	// Set value
	cache.Set(key, value)

	// Get value
	retrieved, exists := cache.Get(key)
	if !exists {
		t.Error("Expected value to exist in cache")
	}

	if retrieved != value {
		t.Errorf("Expected value '%s', got '%s'", value, retrieved)
	}
}

func TestAICache_GetNonExistent(t *testing.T) {
	cache := NewAICache(1 * time.Minute)
	key := "non-existent-key"

	_, exists := cache.Get(key)
	if exists {
		t.Error("Expected non-existent key to return false")
	}
}

func TestAICache_Delete(t *testing.T) {
	cache := NewAICache(1 * time.Minute)
	key := "test-key"
	value := "test-value"

	// Set value
	cache.Set(key, value)

	// Verify it exists
	_, exists := cache.Get(key)
	if !exists {
		t.Error("Expected value to exist before deletion")
	}

	// Delete value
	cache.Delete(key)

	// Verify it's gone
	_, exists = cache.Get(key)
	if exists {
		t.Error("Expected value to be deleted")
	}
}

func TestAICache_Clear(t *testing.T) {
	cache := NewAICache(1 * time.Minute)

	// Add multiple entries
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Verify entries exist
	if cache.Size() != 3 {
		t.Errorf("Expected cache size to be 3, got %d", cache.Size())
	}

	// Clear cache
	cache.Clear()

	// Verify cache is empty
	if cache.Size() != 0 {
		t.Errorf("Expected cache size to be 0 after clear, got %d", cache.Size())
	}

	// Verify individual entries are gone
	_, exists := cache.Get("key1")
	if exists {
		t.Error("Expected key1 to be cleared")
	}
}

func TestAICache_Size(t *testing.T) {
	cache := NewAICache(1 * time.Minute)

	// Initial size should be 0
	if cache.Size() != 0 {
		t.Errorf("Expected initial size to be 0, got %d", cache.Size())
	}

	// Add entries and check size
	cache.Set("key1", "value1")
	if cache.Size() != 1 {
		t.Errorf("Expected size to be 1, got %d", cache.Size())
	}

	cache.Set("key2", "value2")
	if cache.Size() != 2 {
		t.Errorf("Expected size to be 2, got %d", cache.Size())
	}

	// Delete entry and check size
	cache.Delete("key1")
	if cache.Size() != 1 {
		t.Errorf("Expected size to be 1 after deletion, got %d", cache.Size())
	}
}

func TestAICache_Expiration(t *testing.T) {
	// Use very short TTL for testing
	cache := NewAICache(100 * time.Millisecond)
	key := "test-key"
	value := "test-value"

	// Set value
	cache.Set(key, value)

	// Verify it exists immediately
	_, exists := cache.Get(key)
	if !exists {
		t.Error("Expected value to exist immediately after setting")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Verify it's expired
	_, exists = cache.Get(key)
	if exists {
		t.Error("Expected value to be expired")
	}
}

func TestAICache_ConcurrentAccess(t *testing.T) {
	cache := NewAICache(1 * time.Minute)
	done := make(chan bool, 10)

	// Start multiple goroutines
	for i := 0; i < 10; i++ {
		go func(i int) {
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)

			// Set value
			cache.Set(key, value)

			// Get value
			retrieved, exists := cache.Get(key)
			if !exists || retrieved != value {
				t.Errorf("Concurrent access failed for key %s", key)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify final state
	if cache.Size() != 10 {
		t.Errorf("Expected final cache size to be 10, got %d", cache.Size())
	}
}

func TestAICache_ComplexValues(t *testing.T) {
	cache := NewAICache(1 * time.Minute)

	// Test with complex data structures
	complexValue := map[string]interface{}{
		"string": "test",
		"number": 42,
		"array":  []string{"a", "b", "c"},
		"nested": map[string]string{"key": "value"},
	}

	key := "complex-key"
	cache.Set(key, complexValue)

	retrieved, exists := cache.Get(key)
	if !exists {
		t.Error("Expected complex value to exist in cache")
	}

	// Type assertion to verify structure
	if retrievedMap, ok := retrieved.(map[string]interface{}); ok {
		if retrievedMap["string"] != "test" {
			t.Error("Expected string value to be preserved")
		}
		if retrievedMap["number"] != 42 {
			t.Error("Expected number value to be preserved")
		}
	} else {
		t.Error("Expected retrieved value to be a map")
	}
}

func TestAICache_UpdateExistingKey(t *testing.T) {
	cache := NewAICache(1 * time.Minute)
	key := "test-key"

	// Set initial value
	cache.Set(key, "initial-value")

	// Verify initial value
	retrieved, exists := cache.Get(key)
	if !exists || retrieved != "initial-value" {
		t.Error("Expected initial value to be set correctly")
	}

	// Update value
	cache.Set(key, "updated-value")

	// Verify updated value
	retrieved, exists = cache.Get(key)
	if !exists || retrieved != "updated-value" {
		t.Error("Expected value to be updated")
	}

	// Verify size hasn't changed
	if cache.Size() != 1 {
		t.Errorf("Expected cache size to remain 1, got %d", cache.Size())
	}
}
