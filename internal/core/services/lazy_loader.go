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

package services

import (
	"context"
	"sync"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain"
)

// LazyLoader provides lazy loading functionality for server details
type LazyLoader struct {
	loadedDetails map[string]*domain.Server
	loadingMutex  sync.RWMutex
	loadFunc      func(string) (*domain.Server, error)
	cache         map[string]*domain.Server
	cacheMutex    sync.RWMutex
	cacheTTL      time.Duration
}

// LazyLoaderConfig defines configuration for the lazy loader
type LazyLoaderConfig struct {
	CacheTTL time.Duration // Cache time-to-live
}

// DefaultLazyLoaderConfig returns a default configuration
func DefaultLazyLoaderConfig() *LazyLoaderConfig {
	return &LazyLoaderConfig{
		CacheTTL: 5 * time.Minute, // Cache for 5 minutes
	}
}

// NewLazyLoader creates a new lazy loader
func NewLazyLoader(loadFunc func(string) (*domain.Server, error), config *LazyLoaderConfig) *LazyLoader {
	if config == nil {
		config = DefaultLazyLoaderConfig()
	}

	return &LazyLoader{
		loadedDetails: make(map[string]*domain.Server),
		loadFunc:      loadFunc,
		cache:         make(map[string]*domain.Server),
		cacheTTL:      config.CacheTTL,
	}
}

// LoadServerDetails loads server details lazily
func (ll *LazyLoader) LoadServerDetails(serverID string) (*domain.Server, error) {
	// Check cache first
	ll.cacheMutex.RLock()
	if cached, exists := ll.cache[serverID]; exists {
		ll.cacheMutex.RUnlock()
		return cached, nil
	}
	ll.cacheMutex.RUnlock()

	// Check if already loading
	ll.loadingMutex.Lock()
	if server, exists := ll.loadedDetails[serverID]; exists {
		ll.loadingMutex.Unlock()
		return server, nil
	}
	ll.loadingMutex.Unlock()

	// Load the server details
	server, err := ll.loadFunc(serverID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	ll.cacheMutex.Lock()
	ll.cache[serverID] = server
	ll.cacheMutex.Unlock()

	// Store in loaded details
	ll.loadingMutex.Lock()
	ll.loadedDetails[serverID] = server
	ll.loadingMutex.Unlock()

	return server, nil
}

// AsyncResult represents the result of an async operation
type AsyncResult[T any] struct {
	Result T
	Error  error
}

// LoadServerDetailsAsync loads server details asynchronously
func (ll *LazyLoader) LoadServerDetailsAsync(ctx context.Context, serverID string) <-chan AsyncResult[*domain.Server] {
	result := make(chan AsyncResult[*domain.Server], 1)

	go func() {
		defer close(result)
		server, err := ll.LoadServerDetails(serverID)
		result <- AsyncResult[*domain.Server]{
			Result: server,
			Error:  err,
		}
	}()

	return result
}

// PreloadServerDetails preloads server details in the background
func (ll *LazyLoader) PreloadServerDetails(ctx context.Context, serverIDs []string) {
	go func() {
		for _, serverID := range serverIDs {
			select {
			case <-ctx.Done():
				return
			default:
				// Load server details in background
				_, _ = ll.LoadServerDetails(serverID) // Ignore errors in background loading
			}
		}
	}()
}

// BatchLoadServerDetails loads multiple server details concurrently
func (ll *LazyLoader) BatchLoadServerDetails(ctx context.Context, serverIDs []string) <-chan []AsyncResult[*domain.Server] {
	results := make(chan []AsyncResult[*domain.Server], 1)

	go func() {
		defer close(results)

		var wg sync.WaitGroup
		responses := make([]AsyncResult[*domain.Server], len(serverIDs))

		for i, serverID := range serverIDs {
			wg.Add(1)
			go func(index int, id string) {
				defer wg.Done()
				server, err := ll.LoadServerDetails(id)
				responses[index] = AsyncResult[*domain.Server]{
					Result: server,
					Error:  err,
				}
			}(i, serverID)
		}

		wg.Wait()
		results <- responses
	}()

	return results
}

// ClearCache clears the lazy loader cache
func (ll *LazyLoader) ClearCache() {
	ll.cacheMutex.Lock()
	defer ll.cacheMutex.Unlock()

	ll.cache = make(map[string]*domain.Server)
}

// GetCacheStats returns cache statistics
func (ll *LazyLoader) GetCacheStats() map[string]interface{} {
	ll.cacheMutex.RLock()
	defer ll.cacheMutex.RUnlock()

	ll.loadingMutex.RLock()
	defer ll.loadingMutex.RUnlock()

	return map[string]interface{}{
		"cached_servers":    len(ll.cache),
		"loaded_servers":    len(ll.loadedDetails),
		"cache_ttl_seconds": ll.cacheTTL.Seconds(),
	}
}

// IsLoaded checks if server details are already loaded
func (ll *LazyLoader) IsLoaded(serverID string) bool {
	ll.cacheMutex.RLock()
	defer ll.cacheMutex.RUnlock()

	_, exists := ll.cache[serverID]
	return exists
}

// GetLoadedServers returns a list of loaded server IDs
func (ll *LazyLoader) GetLoadedServers() []string {
	ll.cacheMutex.RLock()
	defer ll.cacheMutex.RUnlock()

	servers := make([]string, 0, len(ll.cache))
	for serverID := range ll.cache {
		servers = append(servers, serverID)
	}

	return servers
}
