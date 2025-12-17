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

	"github.com/aryasoni98/wooak/internal/core/domain"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestRepository_ListServers_EmptyConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")
	metaDataPath := filepath.Join(tmpDir, "metadata.json")

	// Create empty config file
	if err := os.WriteFile(configPath, []byte(""), 0o600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	logger := zap.NewNop().Sugar()
	repo := NewRepository(logger, configPath, metaDataPath).(*Repository)

	servers, err := repo.ListServers("")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(servers) != 0 {
		t.Errorf("Expected 0 servers, got %d", len(servers))
	}
}

func TestRepository_ListServers_WithQuery(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")
	metaDataPath := filepath.Join(tmpDir, "metadata.json")

	configContent := `
Host server1
    HostName example.com
    User user1

Host server2
    HostName example2.com
    User user2
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	logger := zap.NewNop().Sugar()
	repo := NewRepository(logger, configPath, metaDataPath).(*Repository)

	// Test listing all servers
	servers, err := repo.ListServers("")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(servers) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(servers))
	}

	// Test filtering by query
	filtered, err := repo.ListServers("server1")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(filtered) != 1 {
		t.Errorf("Expected 1 server after filtering, got %d", len(filtered))
	}
	if len(filtered) > 0 && filtered[0].Alias != "server1" {
		t.Errorf("Expected server1, got %s", filtered[0].Alias)
	}
}

func TestRepository_AddServer(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")
	metaDataPath := filepath.Join(tmpDir, "metadata.json")

	// Create empty config file
	if err := os.WriteFile(configPath, []byte(""), 0o600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	logger := zap.NewNop().Sugar()
	repo := NewRepository(logger, configPath, metaDataPath).(*Repository)

	server := domain.Server{
		Alias: "test-server",
		Host:  "example.com",
		User:  "user",
		Port:  22,
	}

	err := repo.AddServer(server)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify server was added
	servers, err := repo.ListServers("")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(servers) != 1 {
		t.Fatalf("Expected 1 server, got %d", len(servers))
	}
	if servers[0].Alias != "test-server" {
		t.Errorf("Expected alias 'test-server', got %s", servers[0].Alias)
	}
}

func TestRepository_AddServer_Duplicate(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")
	metaDataPath := filepath.Join(tmpDir, "metadata.json")

	configContent := `
Host test-server
    HostName example.com
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	logger := zap.NewNop().Sugar()
	repo := NewRepository(logger, configPath, metaDataPath).(*Repository)

	server := domain.Server{
		Alias: "test-server",
		Host:  "example2.com",
	}

	err := repo.AddServer(server)
	if err == nil {
		t.Error("Expected error when adding duplicate server, got nil")
	}
}

func TestRepository_UpdateServer(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")
	metaDataPath := filepath.Join(tmpDir, "metadata.json")

	configContent := `
Host test-server
    HostName example.com
    User olduser
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	logger := zap.NewNop().Sugar()
	repo := NewRepository(logger, configPath, metaDataPath).(*Repository)

	oldServer := domain.Server{
		Alias: "test-server",
		Host:  "example.com",
		User:  "olduser",
	}

	newServer := domain.Server{
		Alias: "test-server",
		Host:  "example.com",
		User:  "newuser",
	}

	err := repo.UpdateServer(oldServer, newServer)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify server was updated
	servers, err := repo.ListServers("")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(servers) != 1 {
		t.Fatalf("Expected 1 server, got %d", len(servers))
	}
	if servers[0].User != "newuser" {
		t.Errorf("Expected user 'newuser', got %s", servers[0].User)
	}
}

func TestRepository_DeleteServer(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")
	metaDataPath := filepath.Join(tmpDir, "metadata.json")

	configContent := `
Host server1
    HostName example.com

Host server2
    HostName example2.com
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	logger := zap.NewNop().Sugar()
	repo := NewRepository(logger, configPath, metaDataPath).(*Repository)

	server := domain.Server{
		Alias: "server1",
	}

	err := repo.DeleteServer(server)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify server was deleted
	servers, err := repo.ListServers("")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(servers) != 1 {
		t.Fatalf("Expected 1 server, got %d", len(servers))
	}
	if servers[0].Alias != "server2" {
		t.Errorf("Expected remaining server to be 'server2', got %s", servers[0].Alias)
	}
}

func TestRepository_DeleteServer_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")
	metaDataPath := filepath.Join(tmpDir, "metadata.json")

	configContent := `
Host server1
    HostName example.com
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	logger := zap.NewNop().Sugar()
	repo := NewRepository(logger, configPath, metaDataPath).(*Repository)

	server := domain.Server{
		Alias: "nonexistent",
	}

	err := repo.DeleteServer(server)
	if err == nil {
		t.Error("Expected error when deleting non-existent server, got nil")
	}
}

func TestRepository_SetPinned(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")
	metaDataPath := filepath.Join(tmpDir, "metadata.json")

	configContent := `
Host test-server
    HostName example.com
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	logger := zap.NewNop().Sugar()
	repo := NewRepository(logger, configPath, metaDataPath).(*Repository)

	// Pin server
	err := repo.SetPinned("test-server", true)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Unpin server
	err = repo.SetPinned("test-server", false)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestRepository_RecordSSH(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")
	metaDataPath := filepath.Join(tmpDir, "metadata.json")

	configContent := `
Host test-server
    HostName example.com
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	logger := zap.NewNop().Sugar()
	repo := NewRepository(logger, configPath, metaDataPath).(*Repository)

	// Record SSH access
	err := repo.RecordSSH("test-server")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify metadata was updated
	servers, err := repo.ListServers("")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(servers) > 0 {
		// SSH count should be incremented
		if servers[0].SSHCount < 1 {
			t.Errorf("Expected SSHCount >= 1, got %d", servers[0].SSHCount)
		}
	}
}

func TestRepository_LoadConfig_Cache(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")
	metaDataPath := filepath.Join(tmpDir, "metadata.json")

	configContent := `
Host test-server
    HostName example.com
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	logger, _ := config.Build()
	repo := NewRepository(logger.Sugar(), configPath, metaDataPath).(*Repository)

	// First load - should read from file
	cfg1, err := repo.loadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Second load - should use cache
	cfg2, err := repo.loadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify both configs are the same
	if len(cfg1.Hosts) != len(cfg2.Hosts) {
		t.Errorf("Expected same number of hosts, got %d and %d", len(cfg1.Hosts), len(cfg2.Hosts))
	}
}
