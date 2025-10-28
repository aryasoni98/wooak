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
	"net"
	"net/http"
	"sync"
	"time"
)

// ConnectionPool manages a pool of HTTP connections for AI requests
type ConnectionPool struct {
	clients []*http.Client
	current int
	mutex   sync.Mutex
	config  *PoolConfig
}

// PoolConfig defines configuration for the connection pool
type PoolConfig struct {
	MaxConnections        int           // Maximum number of connections in pool
	MaxIdleConns          int           // Maximum idle connections per host
	MaxConnsPerHost       int           // Maximum connections per host
	IdleConnTimeout       time.Duration // Idle connection timeout
	ResponseHeaderTimeout time.Duration // Response header timeout
	RequestTimeout        time.Duration // Request timeout
}

// DefaultPoolConfig returns a default configuration for the connection pool
func DefaultPoolConfig() *PoolConfig {
	return &PoolConfig{
		MaxConnections:        DefaultMaxConnections * 2, // 10
		MaxIdleConns:          DefaultMaxIdleConns * 2,   // 10
		MaxConnsPerHost:       DefaultMaxConnsPerHost,
		IdleConnTimeout:       DefaultIdleConnTimeout,
		ResponseHeaderTimeout: DefaultResponseTimeout,
		RequestTimeout:        30 * time.Second,
	}
}

// NewConnectionPool creates a new HTTP connection pool
func NewConnectionPool(config *PoolConfig) *ConnectionPool {
	if config == nil {
		config = DefaultPoolConfig()
	}

	pool := &ConnectionPool{
		clients: make([]*http.Client, config.MaxConnections),
		config:  config,
	}

	// Initialize HTTP clients with optimized settings
	for i := 0; i < config.MaxConnections; i++ {
		pool.clients[i] = &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:          config.MaxIdleConns,
				MaxConnsPerHost:       config.MaxConnsPerHost,
				IdleConnTimeout:       config.IdleConnTimeout,
				ResponseHeaderTimeout: config.ResponseHeaderTimeout,
				DisableKeepAlives:     false,
				DisableCompression:    false,
				DialContext: (&net.Dialer{
					Timeout:   DefaultDialTimeout,
					KeepAlive: DefaultKeepAliveTimeout,
				}).DialContext,
				TLSHandshakeTimeout: DefaultTLSHandshakeTimeout,
			},
			Timeout: config.RequestTimeout,
		}
	}

	return pool
}

// GetClient returns an HTTP client from the pool
func (p *ConnectionPool) GetClient() *http.Client {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	client := p.clients[p.current]
	p.current = (p.current + 1) % len(p.clients)
	return client
}

// GetClientWithContext returns an HTTP client with a custom context
func (p *ConnectionPool) GetClientWithContext(ctx context.Context) *http.Client {
	client := p.GetClient()
	// Create a new client with the provided context
	return &http.Client{
		Transport: client.Transport,
		Timeout:   client.Timeout,
	}
}

// Stats returns connection pool statistics
func (p *ConnectionPool) Stats() map[string]interface{} {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return map[string]interface{}{
		"max_connections":    p.config.MaxConnections,
		"current_index":      p.current,
		"max_idle_conns":     p.config.MaxIdleConns,
		"max_conns_per_host": p.config.MaxConnsPerHost,
		"idle_conn_timeout":  p.config.IdleConnTimeout.Seconds(),
		"request_timeout":    p.config.RequestTimeout.Seconds(),
	}
}

// Close closes all connections in the pool
func (p *ConnectionPool) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, client := range p.clients {
		if transport, ok := client.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}
}

// Global connection pool instance
var (
	globalConnectionPool *ConnectionPool
	poolOnce             sync.Once
)

// GetGlobalConnectionPool returns the global connection pool instance
func GetGlobalConnectionPool() *ConnectionPool {
	poolOnce.Do(func() {
		globalConnectionPool = NewConnectionPool(DefaultPoolConfig())
	})
	return globalConnectionPool
}

// CloseGlobalConnectionPool closes the global connection pool
func CloseGlobalConnectionPool() {
	if globalConnectionPool != nil {
		globalConnectionPool.Close()
	}
}
