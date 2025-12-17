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
	"path/filepath"
	"testing"

	"github.com/kevinburke/ssh_config"
)

// BenchmarkLoadConfig benchmarks config loading from file
func BenchmarkLoadConfig(b *testing.B) {
	// Create a temporary config file
	tmpDir := b.TempDir()
	configPath := filepath.Join(tmpDir, "config")

	// Write a sample config
	configContent := `
Host server1
    HostName example.com
    User user1
    Port 22

Host server2
    HostName example2.com
    User user2
    Port 2222
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		b.Fatalf("Failed to create test config: %v", err)
	}

	repo := &Repository{
		configPath: configPath,
		fileSystem: DefaultFileSystem{},
		cache:      newConfigCache(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.loadConfig()
		if err != nil {
			b.Fatalf("Failed to load config: %v", err)
		}
	}
}

// BenchmarkLoadConfigCached benchmarks cached config loading
func BenchmarkLoadConfigCached(b *testing.B) {
	tmpDir := b.TempDir()
	configPath := filepath.Join(tmpDir, "config")

	configContent := `
Host server1
    HostName example.com
    User user1
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		b.Fatalf("Failed to create test config: %v", err)
	}

	repo := &Repository{
		configPath: configPath,
		fileSystem: DefaultFileSystem{},
		cache:      newConfigCache(),
	}

	// Prime the cache
	_, err := repo.loadConfig()
	if err != nil {
		b.Fatalf("Failed to prime cache: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := repo.loadConfig()
		if err != nil {
			b.Fatalf("Failed to load cached config: %v", err)
		}
	}
}

// BenchmarkDecodeConfig benchmarks SSH config decoding
func BenchmarkDecodeConfig(b *testing.B) {
	configContent := `
Host server1
    HostName example.com
    User user1
    Port 22
    IdentityFile ~/.ssh/id_rsa

Host server2
    HostName example2.com
    User user2
    Port 2222
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ssh_config.DecodeBytes([]byte(configContent))
		if err != nil {
			b.Fatalf("Failed to decode config: %v", err)
		}
	}
}
