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
	"context"
	"net/http"
	"testing"
	"time"
)

func TestConnectionPool_NewPool(t *testing.T) {
	config := &PoolConfig{
		MaxConnections:        5,
		MaxIdleConns:          3,
		MaxConnsPerHost:       2,
		IdleConnTimeout:       60 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		RequestTimeout:        10 * time.Second,
	}

	pool := NewConnectionPool(config)

	if pool == nil {
		t.Fatal("Expected pool to be created")
	}

	if len(pool.clients) != config.MaxConnections {
		t.Errorf("Expected %d clients, got %d", config.MaxConnections, len(pool.clients))
	}

	if pool.config != config {
		t.Error("Expected config to be set")
	}
}

func TestConnectionPool_DefaultConfig(t *testing.T) {
	pool := NewConnectionPool(nil)

	if pool == nil {
		t.Fatal("Expected pool to be created with default config")
	}

	expectedMaxConnections := 10
	if len(pool.clients) != expectedMaxConnections {
		t.Errorf("Expected %d clients with default config, got %d", expectedMaxConnections, len(pool.clients))
	}
}

func TestConnectionPool_GetClient(t *testing.T) {
	config := &PoolConfig{
		MaxConnections: 3,
		RequestTimeout: 5 * time.Second,
	}

	pool := NewConnectionPool(config)

	// Get multiple clients and verify they're different instances
	clients := make([]*http.Client, 5)
	for i := 0; i < 5; i++ {
		clients[i] = pool.GetClient()
	}

	// Should cycle through the pool
	if clients[0] == clients[1] {
		t.Error("Expected different clients from pool")
	}

	if clients[1] == clients[2] {
		t.Error("Expected different clients from pool")
	}

	// Should start cycling after MaxConnections
	if clients[0] != clients[3] {
		t.Error("Expected client to cycle back to first after MaxConnections")
	}
}

func TestConnectionPool_GetClientWithContext(t *testing.T) {
	config := &PoolConfig{
		MaxConnections: 2,
		RequestTimeout: 5 * time.Second,
	}

	pool := NewConnectionPool(config)
	ctx := context.Background()

	client := pool.GetClientWithContext(ctx)

	if client == nil {
		t.Fatal("Expected client to be returned")
	}

	if client.Timeout != config.RequestTimeout {
		t.Errorf("Expected timeout to be %v, got %v", config.RequestTimeout, client.Timeout)
	}
}

func TestConnectionPool_Stats(t *testing.T) {
	config := &PoolConfig{
		MaxConnections:        5,
		MaxIdleConns:          3,
		MaxConnsPerHost:       2,
		IdleConnTimeout:       60 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		RequestTimeout:        10 * time.Second,
	}

	pool := NewConnectionPool(config)
	stats := pool.Stats()

	expectedStats := map[string]interface{}{
		"max_connections":    config.MaxConnections,
		"current_index":      0, // Should start at 0
		"max_idle_conns":     config.MaxIdleConns,
		"max_conns_per_host": config.MaxConnsPerHost,
		"idle_conn_timeout":  config.IdleConnTimeout.Seconds(),
		"request_timeout":    config.RequestTimeout.Seconds(),
	}

	for key, expectedValue := range expectedStats {
		if stats[key] != expectedValue {
			t.Errorf("Expected stats[%s] to be %v, got %v", key, expectedValue, stats[key])
		}
	}
}

func TestConnectionPool_ConcurrentAccess(t *testing.T) {
	config := &PoolConfig{
		MaxConnections: 3,
		RequestTimeout: 5 * time.Second,
	}

	pool := NewConnectionPool(config)

	// Test concurrent access
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			client := pool.GetClient()
			if client == nil {
				t.Error("Expected client to be returned")
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// No panics should occur
}

func TestConnectionPool_Close(t *testing.T) {
	config := &PoolConfig{
		MaxConnections: 2,
		RequestTimeout: 5 * time.Second,
	}

	pool := NewConnectionPool(config)

	// Close should not panic
	pool.Close()

	// Stats should still work after close
	stats := pool.Stats()
	if stats == nil {
		t.Error("Expected stats to be available after close")
	}
}

func TestGlobalConnectionPool(t *testing.T) {
	// Test that global pool is created
	pool1 := GetGlobalConnectionPool()
	pool2 := GetGlobalConnectionPool()

	if pool1 == nil || pool2 == nil {
		t.Fatal("Expected global pool to be created")
	}

	// Should be the same instance
	if pool1 != pool2 {
		t.Error("Expected global pool to be a singleton")
	}

	// Test close
	CloseGlobalConnectionPool()
	// Should not panic
}

func TestConnectionPool_ClientReuse(t *testing.T) {
	config := &PoolConfig{
		MaxConnections: 2,
		RequestTimeout: 5 * time.Second,
	}

	pool := NewConnectionPool(config)

	// Get clients and verify they're properly configured
	client1 := pool.GetClient()
	client2 := pool.GetClient()

	if client1.Timeout != config.RequestTimeout {
		t.Errorf("Expected client1 timeout to be %v, got %v", config.RequestTimeout, client1.Timeout)
	}

	if client2.Timeout != config.RequestTimeout {
		t.Errorf("Expected client2 timeout to be %v, got %v", config.RequestTimeout, client2.Timeout)
	}

	// Verify transport settings
	if transport1, ok := client1.Transport.(*http.Transport); ok {
		if transport1.MaxIdleConns != config.MaxIdleConns {
			t.Errorf("Expected MaxIdleConns to be %d, got %d", config.MaxIdleConns, transport1.MaxIdleConns)
		}
	} else {
		t.Error("Expected client to have http.Transport")
	}
}
