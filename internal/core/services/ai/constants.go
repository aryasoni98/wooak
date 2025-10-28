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

import "time"

const (
	// Cache configuration
	DefaultCacheCleanupInterval = 5 * time.Minute

	// Connection pool configuration
	DefaultMaxConnections      = 5
	DefaultMaxIdleConns        = 5
	DefaultMaxConnsPerHost     = 3
	DefaultIdleConnTimeout     = 90 * time.Second
	DefaultResponseTimeout     = 10 * time.Second
	DefaultDialTimeout         = 5 * time.Second
	DefaultKeepAliveTimeout    = 30 * time.Second
	DefaultTLSHandshakeTimeout = 5 * time.Second

	// Ollama service configuration
	OllamaConnectionTimeout = 2 * time.Second
	OllamaHealthCheckURL    = "/api/tags"
)
