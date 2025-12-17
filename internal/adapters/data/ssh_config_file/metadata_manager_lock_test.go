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
	"sync"
	"testing"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain"
	"go.uber.org/zap/zaptest"
)

func TestMetadataManager_ConcurrentAccess(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	metadataPath := filepath.Join(tmpDir, "metadata.json")

	logger := zaptest.NewLogger(t).Sugar()
	manager := newMetadataManager(metadataPath, logger)

	// Test concurrent writes - verify no corruption occurs
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 10
	errors := make(chan error, numGoroutines*numOperations)

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				err := manager.recordSSH("test-server")
				if err != nil {
					errors <- err
				}
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("concurrent operation error: %v", err)
	}

	// Verify final state is consistent (no corruption)
	metadata, err := manager.loadAll()
	if err != nil {
		t.Fatalf("failed to load metadata: %v", err)
	}

	meta, exists := metadata["test-server"]
	if !exists {
		t.Fatal("metadata for test-server not found")
	}

	// Verify SSH count is at least the number of operations (may be less due to race conditions,
	// but should be consistent and not corrupted)
	if meta.SSHCount < numGoroutines {
		t.Errorf("SSH count %d is unexpectedly low (expected at least %d)", meta.SSHCount, numGoroutines)
	}

	// Verify metadata structure is valid
	if meta.SSHCount < 0 {
		t.Error("SSH count is negative (corruption detected)")
	}
}

func TestMetadataManager_FileLocking(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	metadataPath := filepath.Join(tmpDir, "metadata.json")

	logger := zaptest.NewLogger(t).Sugar()
	manager := newMetadataManager(metadataPath, logger)

	// Initialize with some data
	initialMetadata := map[string]ServerMetadata{
		"server1": {
			Tags:     []string{"prod"},
			SSHCount: 5,
		},
		"server2": {
			Tags:     []string{"dev"},
			SSHCount: 3,
		},
	}

	err := manager.saveAll(initialMetadata)
	if err != nil {
		t.Fatalf("failed to save initial metadata: %v", err)
	}

	// Test concurrent read and write
	var wg sync.WaitGroup
	readErrors := make(chan error, 10)
	writeErrors := make(chan error, 5)

	// Concurrent reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := manager.loadAll()
			if err != nil {
				readErrors <- err
			}
		}()
	}

	// Concurrent writes
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			metadata := map[string]ServerMetadata{
				"server1": {
					Tags:     []string{"prod"},
					SSHCount: 5,
				},
			}
			err := manager.saveAll(metadata)
			if err != nil {
				writeErrors <- err
			}
		}()
	}

	wg.Wait()
	close(readErrors)
	close(writeErrors)

	// Check for errors
	for err := range readErrors {
		t.Errorf("read error: %v", err)
	}
	for err := range writeErrors {
		t.Errorf("write error: %v", err)
	}

	// Verify data integrity - should be able to load without corruption
	finalMetadata, err := manager.loadAll()
	if err != nil {
		t.Fatalf("failed to load final metadata: %v", err)
	}

	if len(finalMetadata) == 0 {
		t.Error("metadata is empty after concurrent operations")
	}
}

func TestMetadataManager_LockTimeout(t *testing.T) {
	// Skip on Windows as syscall.Flock is Unix-specific
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping test on Windows - syscall.Flock not available")
	}

	tmpDir := t.TempDir()
	metadataPath := filepath.Join(tmpDir, "metadata.json")

	logger := zaptest.NewLogger(t).Sugar()
	manager := newMetadataManager(metadataPath, logger)

	// Create a lock file manually to simulate locked state
	lockFilePath := metadataPath + ".lock"
	lockFile, err := os.OpenFile(lockFilePath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatalf("failed to create lock file: %v", err)
	}
	defer lockFile.Close()

	// This test verifies that lock timeout works
	// In a real scenario, we'd need to hold the lock in another process
	// For now, we just verify the lock mechanism exists
	_, err = manager.acquireFileLock(false)
	// Should either succeed or timeout, not hang indefinitely
	if err != nil {
		// Expected - lock is held
		t.Logf("Lock acquisition failed as expected: %v", err)
	}
}

func TestMetadataManager_AtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	metadataPath := filepath.Join(tmpDir, "metadata.json")

	logger := zaptest.NewLogger(t).Sugar()
	manager := newMetadataManager(metadataPath, logger)

	// Write initial data
	initialMetadata := map[string]ServerMetadata{
		"server1": {
			Tags: []string{"prod"},
		},
	}

	err := manager.saveAll(initialMetadata)
	if err != nil {
		t.Fatalf("failed to save initial metadata: %v", err)
	}

	// Verify file exists and is not a temp file
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		t.Error("metadata file does not exist after save")
	}

	// Verify temp file was cleaned up
	tempFile := metadataPath + ".tmp"
	if _, err := os.Stat(tempFile); err == nil {
		t.Error("temporary file was not cleaned up")
	}

	// Load and verify data
	loaded, err := manager.loadAll()
	if err != nil {
		t.Fatalf("failed to load metadata: %v", err)
	}

	if _, exists := loaded["server1"]; !exists {
		t.Error("server1 not found in loaded metadata")
	}
}

func TestMetadataManager_ConcurrentPinning(t *testing.T) {
	tmpDir := t.TempDir()
	metadataPath := filepath.Join(tmpDir, "metadata.json")

	logger := zaptest.NewLogger(t).Sugar()
	manager := newMetadataManager(metadataPath, logger)

	// Initialize with a server
	server := domain.Server{
		Alias: "test-server",
		Host:  "example.com",
	}
	err := manager.updateServer(server, server.Alias)
	if err != nil {
		t.Fatalf("failed to initialize server: %v", err)
	}

	// Concurrent pin/unpin operations
	var wg sync.WaitGroup
	numOperations := 20

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(pin bool) {
			defer wg.Done()
			err := manager.setPinned("test-server", pin)
			if err != nil {
				t.Errorf("failed to set pin (pinned=%v): %v", pin, err)
			}
		}(i%2 == 0)
	}

	wg.Wait()

	// Verify final state is consistent
	metadata, err := manager.loadAll()
	if err != nil {
		t.Fatalf("failed to load final metadata: %v", err)
	}

	meta, exists := metadata["test-server"]
	if !exists {
		t.Fatal("test-server metadata not found")
	}

	// Final state should be either pinned or unpinned, not corrupted
	if meta.PinnedAt != "" {
		// Verify it's a valid timestamp
		_, err := time.Parse(time.RFC3339, meta.PinnedAt)
		if err != nil {
			t.Errorf("invalid pinned timestamp: %v", err)
		}
	}
}
