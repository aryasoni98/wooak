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
	"os"
	"testing"
	"time"

	"github.com/kevinburke/ssh_config"
)

// mockFileSystemForCache is a mock filesystem for cache testing
type mockFileSystemForCache struct {
	DefaultFileSystem
	modTime  time.Time
	fileSize int64
	statErr  error
}

func (m *mockFileSystemForCache) Stat(name string) (os.FileInfo, error) {
	if m.statErr != nil {
		return nil, m.statErr
	}
	return &mockFileInfo{
		name:    name,
		size:    m.fileSize,
		modTime: m.modTime,
	}, nil
}

type mockFileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return m.size }
func (m *mockFileInfo) Mode() os.FileMode  { return 0o600 }
func (m *mockFileInfo) ModTime() time.Time { return m.modTime }
func (m *mockFileInfo) IsDir() bool        { return false }
func (m *mockFileInfo) Sys() interface{}   { return nil }

func TestConfigCache_Get_Set(t *testing.T) {
	cache := newConfigCache()
	fs := &mockFileSystemForCache{
		modTime:  time.Now(),
		fileSize: 1024,
	}
	configPath := "/test/config"

	// Test empty cache
	cfg, valid := cache.get(fs, configPath)
	if valid {
		t.Error("Expected cache to be invalid when empty")
	}
	if cfg != nil {
		t.Error("Expected nil config from empty cache")
	}

	// Set cache
	testConfig := &ssh_config.Config{
		Hosts: []*ssh_config.Host{
			{Patterns: []*ssh_config.Pattern{}},
		},
	}
	cache.set(fs, configPath, testConfig)

	// Test getting from cache
	cfg, valid = cache.get(fs, configPath)
	if !valid {
		t.Error("Expected cache to be valid after setting")
	}
	if cfg == nil {
		t.Error("Expected non-nil config from cache")
	}
}

func TestConfigCache_Invalidate(t *testing.T) {
	cache := newConfigCache()
	fs := &mockFileSystemForCache{
		modTime:  time.Now(),
		fileSize: 1024,
	}
	configPath := "/test/config"

	testConfig := &ssh_config.Config{
		Hosts: []*ssh_config.Host{},
	}
	cache.set(fs, configPath, testConfig)

	// Verify cache is set
	if !cache.isCached() {
		t.Error("Expected cache to be set")
	}

	// Invalidate
	cache.invalidate()

	// Verify cache is cleared
	if cache.isCached() {
		t.Error("Expected cache to be invalidated")
	}

	cfg, valid := cache.get(fs, configPath)
	if valid {
		t.Error("Expected cache to be invalid after invalidation")
	}
	if cfg != nil {
		t.Error("Expected nil config after invalidation")
	}
}

func TestConfigCache_FileModified(t *testing.T) {
	cache := newConfigCache()
	originalTime := time.Now()
	fs := &mockFileSystemForCache{
		modTime:  originalTime,
		fileSize: 1024,
	}
	configPath := "/test/config"

	testConfig := &ssh_config.Config{
		Hosts: []*ssh_config.Host{},
	}
	cache.set(fs, configPath, testConfig)

	// Verify cache is valid
	_, valid := cache.get(fs, configPath)
	if !valid {
		t.Error("Expected cache to be valid")
	}

	// Modify file time
	fs.modTime = originalTime.Add(1 * time.Second)

	// Cache should be invalid now
	_, valid = cache.get(fs, configPath)
	if valid {
		t.Error("Expected cache to be invalid after file modification")
	}
}

func TestConfigCache_FileSizeChanged(t *testing.T) {
	cache := newConfigCache()
	fs := &mockFileSystemForCache{
		modTime:  time.Now(),
		fileSize: 1024,
	}
	configPath := "/test/config"

	testConfig := &ssh_config.Config{
		Hosts: []*ssh_config.Host{},
	}
	cache.set(fs, configPath, testConfig)

	// Verify cache is valid
	_, valid := cache.get(fs, configPath)
	if !valid {
		t.Error("Expected cache to be valid")
	}

	// Change file size
	fs.fileSize = 2048

	// Cache should be invalid now
	_, valid = cache.get(fs, configPath)
	if valid {
		t.Error("Expected cache to be invalid after file size change")
	}
}

func TestConfigCache_FileNotExists(t *testing.T) {
	cache := newConfigCache()
	fs := &mockFileSystemForCache{
		modTime:  time.Now(),
		fileSize: 1024,
		statErr:  os.ErrNotExist,
	}
	configPath := "/test/config"

	testConfig := &ssh_config.Config{
		Hosts: []*ssh_config.Host{},
	}
	cache.set(fs, configPath, testConfig)

	// Cache should be invalid when file doesn't exist
	cfg, valid := cache.get(fs, configPath)
	if valid {
		t.Error("Expected cache to be invalid when file doesn't exist")
	}
	if cfg != nil {
		t.Error("Expected nil config when file doesn't exist")
	}
}
