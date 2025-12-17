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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain"
	"go.uber.org/zap"
)

type ServerMetadata struct {
	Tags     []string `json:"tags,omitempty"`
	LastSeen string   `json:"last_seen,omitempty"`
	PinnedAt string   `json:"pinned_at,omitempty"`
	SSHCount int      `json:"ssh_count,omitempty"`
}

type metadataManager struct {
	filePath string
	logger   *zap.SugaredLogger
	mu       sync.RWMutex // Protects in-process concurrent access
}

const (
	// Default lock timeout for file locking operations
	defaultLockTimeout = 5 * time.Second
	// Default lock retry interval
	defaultLockRetryInterval = 100 * time.Millisecond
)

func newMetadataManager(filePath string, logger *zap.SugaredLogger) *metadataManager {
	return &metadataManager{filePath: filePath, logger: logger}
}

// loadAll loads all metadata from the file with file locking protection.
func (m *metadataManager) loadAll() (map[string]ServerMetadata, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.loadAllWithLock()
}

// loadAllWithLock performs the actual load operation (called with lock held).
func (m *metadataManager) loadAllWithLock() (map[string]ServerMetadata, error) {
	metadata := make(map[string]ServerMetadata)

	if _, err := os.Stat(m.filePath); os.IsNotExist(err) {
		return metadata, nil
	}

	// Acquire file lock for reading
	lockFile, err := m.acquireFileLock(true)
	if err != nil {
		return nil, fmt.Errorf("acquire read lock for '%s': %w", m.filePath, err)
	}
	defer func() {
		if lockFile != nil {
			_ = m.releaseFileLock(lockFile)
		}
	}()

	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return nil, fmt.Errorf("read metadata '%s': %w", m.filePath, err)
	}

	if len(data) == 0 {
		return metadata, nil
	}

	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("parse metadata JSON '%s': %w", m.filePath, err)
	}

	return metadata, nil
}

// saveAll saves all metadata to the file with file locking protection.
func (m *metadataManager) saveAll(metadata map[string]ServerMetadata) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.saveAllWithLock(metadata)
}

// saveAllWithLock performs the actual save operation (called with lock held).
//
// This function implements atomic file writes using a write-temp-then-rename pattern:
//  1. Write JSON data to temporary file (metadata.json.tmp)
//  2. Acquire exclusive file lock
//  3. Rename temp file to final location (atomic on most filesystems)
//  4. Release lock
//
// This ensures:
//   - Atomicity: readers never see partially written files
//   - Consistency: file is either old or new, never corrupted
//   - Durability: rename is atomic, preventing data loss on crashes
//
// The temporary file is cleaned up on error to prevent disk space leaks.
func (m *metadataManager) saveAllWithLock(metadata map[string]ServerMetadata) error {
	if err := m.ensureDirectory(); err != nil {
		m.logger.Errorw("failed to ensure metadata directory", "path", m.filePath, "error", err)
		return fmt.Errorf("ensure metadata directory for '%s': %w", m.filePath, err)
	}

	// Acquire exclusive file lock for writing
	lockFile, err := m.acquireFileLock(false)
	if err != nil {
		return fmt.Errorf("acquire write lock for '%s': %w", m.filePath, err)
	}
	defer func() {
		if lockFile != nil {
			_ = m.releaseFileLock(lockFile)
		}
	}()

	// Write to temporary file first for atomic operation
	tempFile := m.filePath + ".tmp"
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		m.logger.Errorw("failed to marshal metadata", "path", m.filePath, "error", err)
		return fmt.Errorf("marshal metadata for '%s': %w", m.filePath, err)
	}

	if err := os.WriteFile(tempFile, data, 0o600); err != nil {
		m.logger.Errorw("failed to write temporary metadata file", "path", tempFile, "error", err)
		return fmt.Errorf("write temporary metadata '%s': %w", tempFile, err)
	}

	// Atomic rename
	if err := os.Rename(tempFile, m.filePath); err != nil {
		m.logger.Errorw("failed to rename temporary metadata file", "temp", tempFile, "target", m.filePath, "error", err)
		// Clean up temp file on error
		_ = os.Remove(tempFile)
		return fmt.Errorf("rename temporary metadata '%s' to '%s': %w", tempFile, m.filePath, err)
	}

	return nil
}

func (m *metadataManager) updateServer(server domain.Server, oldAlias string) error {
	metadata, err := m.loadAll()
	if err != nil {
		m.logger.Errorw("failed to load metadata in updateServer", "path", m.filePath, "alias", server.Alias, "old_alias", oldAlias, "error", err)
		return fmt.Errorf("load metadata: %w", err)
	}

	if oldAlias != server.Alias {
		oldMeta, ok := metadata[oldAlias]
		if ok {
			metadata[server.Alias] = oldMeta
		}
		delete(metadata, oldAlias)
	}

	existing := metadata[server.Alias]
	merged := existing

	merged.Tags = server.Tags

	if !server.LastSeen.IsZero() {
		merged.LastSeen = server.LastSeen.Format(time.RFC3339)
	}

	if !server.PinnedAt.IsZero() {
		merged.PinnedAt = server.PinnedAt.Format(time.RFC3339)
	}

	if server.SSHCount > 0 {
		merged.SSHCount = server.SSHCount
	}

	metadata[server.Alias] = merged
	return m.saveAll(metadata)
}

func (m *metadataManager) deleteServer(alias string) error {
	metadata, err := m.loadAll()
	if err != nil {
		m.logger.Errorw("failed to load metadata in deleteServer", "path", m.filePath, "alias", alias, "error", err)
		return fmt.Errorf("load metadata: %w", err)
	}

	delete(metadata, alias)
	return m.saveAll(metadata)
}

func (m *metadataManager) setPinned(alias string, pinned bool) error {
	metadata, err := m.loadAll()
	if err != nil {
		m.logger.Errorw("failed to load metadata in setPinned", "path", m.filePath, "alias", alias, "pinned", pinned, "error", err)
		return fmt.Errorf("load metadata: %w", err)
	}

	meta := metadata[alias]
	if pinned {
		meta.PinnedAt = time.Now().Format(time.RFC3339)
	} else {
		meta.PinnedAt = ""
	}

	metadata[alias] = meta
	return m.saveAll(metadata)
}

func (m *metadataManager) recordSSH(alias string) error {
	metadata, err := m.loadAll()
	if err != nil {
		m.logger.Errorw("failed to load metadata in recordSSH", "path", m.filePath, "alias", alias, "error", err)
		return fmt.Errorf("load metadata: %w", err)
	}

	meta := metadata[alias]
	meta.LastSeen = time.Now().Format(time.RFC3339)
	meta.SSHCount++

	metadata[alias] = meta
	return m.saveAll(metadata)
}

func (m *metadataManager) ensureDirectory() error {
	dir := filepath.Dir(m.filePath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("mkdir '%s': %w", dir, err)
	}
	return nil
}

// acquireFileLock acquires a file lock (shared for read, exclusive for write).
//
// This function implements file-level locking using syscall.Flock to prevent
// concurrent access from multiple processes. It uses a separate lock file
// (.lock extension) to coordinate access to the metadata file.
//
// Locking strategy:
//   - Shared lock (LOCK_SH): allows multiple readers simultaneously
//   - Exclusive lock (LOCK_EX): prevents all other access (readers and writers)
//   - Non-blocking (LOCK_NB): returns immediately if lock unavailable
//   - Timeout-based retry: retries every 100ms up to 5 seconds
//
// The lock file is created if it doesn't exist and has restrictive permissions (0600).
// This ensures only the owner can access the lock file.
//
// Returns the lock file handle (must be closed via releaseFileLock) and any error.
func (m *metadataManager) acquireFileLock(shared bool) (*os.File, error) {
	lockFilePath := m.filePath + ".lock"
	lockDir := filepath.Dir(lockFilePath)

	// Ensure lock directory exists
	if err := os.MkdirAll(lockDir, 0o750); err != nil {
		return nil, fmt.Errorf("create lock directory: %w", err)
	}

	// Open lock file
	lockFile, err := os.OpenFile(lockFilePath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open lock file: %w", err)
	}

	// Determine lock type
	lockType := syscall.LOCK_EX // Exclusive lock for writes
	if shared {
		lockType = syscall.LOCK_SH // Shared lock for reads
	}

	// Try to acquire lock with timeout
	deadline := time.Now().Add(defaultLockTimeout)
	for time.Now().Before(deadline) {
		err = syscall.Flock(int(lockFile.Fd()), lockType|syscall.LOCK_NB)
		if err == nil {
			return lockFile, nil
		}

		// If error is not EWOULDBLOCK, it's a real error
		if err != syscall.EWOULDBLOCK {
			_ = lockFile.Close()
			return nil, fmt.Errorf("flock error: %w", err)
		}

		// Wait before retrying
		time.Sleep(defaultLockRetryInterval)
	}

	_ = lockFile.Close()
	return nil, fmt.Errorf("timeout acquiring lock after %v", defaultLockTimeout)
}

// releaseFileLock releases a file lock.
func (m *metadataManager) releaseFileLock(lockFile *os.File) error {
	if lockFile == nil {
		return nil
	}

	err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)
	if closeErr := lockFile.Close(); closeErr != nil && err == nil {
		err = closeErr
	}

	return err
}
