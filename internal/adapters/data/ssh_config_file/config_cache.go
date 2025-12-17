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

package ssh_config_file

import (
	"sync"
	"time"

	"github.com/kevinburke/ssh_config"
)

// configCache holds a cached SSH config with its modification time
type configCache struct {
	mu          sync.RWMutex
	config      *ssh_config.Config
	modTime     time.Time
	fileSize    int64
	lastChecked time.Time
}

// newConfigCache creates a new config cache
func newConfigCache() *configCache {
	return &configCache{}
}

// get retrieves the cached config if it's still valid
func (c *configCache) get(fileSystem FileSystem, configPath string) (*ssh_config.Config, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// If cache is empty, it's not valid
	if c.config == nil {
		return nil, false
	}

	// Check file modification time
	info, err := fileSystem.Stat(configPath)
	if err != nil {
		// File doesn't exist or can't be accessed - cache is invalid
		return nil, false
	}

	// Check if file has been modified
	if !info.ModTime().Equal(c.modTime) || info.Size() != c.fileSize {
		return nil, false
	}

	return c.config, true
}

// set stores a config in the cache with its modification time
func (c *configCache) set(fileSystem FileSystem, configPath string, config *ssh_config.Config) {
	c.mu.Lock()
	defer c.mu.Unlock()

	info, err := fileSystem.Stat(configPath)
	if err != nil {
		// Can't get file info, don't cache
		return
	}

	c.config = config
	c.modTime = info.ModTime()
	c.fileSize = info.Size()
	c.lastChecked = time.Now()
}

// invalidate clears the cache
func (c *configCache) invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.config = nil
	c.modTime = time.Time{}
	c.fileSize = 0
	c.lastChecked = time.Time{}
}

// isCached returns whether there's a valid cached config
func (c *configCache) isCached() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config != nil
}
