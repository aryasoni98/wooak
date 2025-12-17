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

package security

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain/security"
	"go.uber.org/zap"
)

// AuditLogger handles security event logging
type AuditLogger struct {
	policy     *security.SecurityPolicy
	logPath    string
	eventQueue chan *security.SecurityEvent
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	started    bool
	mutex      sync.RWMutex
	startOnce  sync.Once
	stopOnce   sync.Once
	logger     *zap.SugaredLogger
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(policy *security.SecurityPolicy) *AuditLogger {
	return NewAuditLoggerWithLogger(policy, nil)
}

// NewAuditLoggerWithLogger creates a new audit logger with a logger instance
func NewAuditLoggerWithLogger(policy *security.SecurityPolicy, logger *zap.SugaredLogger) *AuditLogger {
	home, _ := os.UserHomeDir()
	logDir := filepath.Join(home, ".wooak", "logs")

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, LogDirectoryPermissions); err != nil {
		// Log error but continue - this is not critical
		if logger != nil {
			logger.Warnw("Failed to create log directory", "error", err, "directory", logDir)
		} else {
			// Fallback to stderr if no logger provided
			_, _ = fmt.Fprintf(os.Stderr, "Warning: Failed to create log directory: %v\n", err)
		}
	}

	logPath := filepath.Join(logDir, "security-audit.log")
	ctx, cancel := context.WithCancel(context.Background())

	return &AuditLogger{
		policy:     policy,
		logPath:    logPath,
		eventQueue: make(chan *security.SecurityEvent, DefaultAuditLogQueueSize),
		ctx:        ctx,
		cancel:     cancel,
		started:    false,
		logger:     logger,
	}
}

// LogEvent logs a security event
func (al *AuditLogger) LogEvent(event *security.SecurityEvent) {
	if !al.policy.EnableAuditLog {
		return
	}

	// Check if we should log this event based on severity
	if !al.shouldLogEvent(event.Severity) {
		return
	}

	// Format the log entry
	logEntry := map[string]interface{}{
		"timestamp": event.Timestamp.Format(time.RFC3339),
		"id":        event.ID,
		"type":      event.Type,
		"severity":  event.Severity,
		"message":   event.Message,
		"source":    event.Source,
		"action":    event.Action,
		"result":    event.Result,
	}

	if event.User != "" {
		logEntry["user"] = event.User
	}

	if event.Host != "" {
		logEntry["host"] = event.Host
	}

	if len(event.Details) > 0 {
		logEntry["details"] = event.Details
	}

	// Convert to JSON
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		if al.logger != nil {
			al.logger.Errorw("Failed to marshal audit log entry", "error", err, "event_id", event.ID)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to marshal audit log entry: %v\n", err)
		}
		return
	}

	// Write to log file
	al.writeToLogFile(string(jsonData))
}

// writeToLogFile writes a log entry to the audit log file
func (al *AuditLogger) writeToLogFile(entry string) {
	file, err := os.OpenFile(al.logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, LogFilePermissions)
	if err != nil {
		if al.logger != nil {
			al.logger.Errorw("Failed to open audit log file", "error", err, "path", al.logPath)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to open audit log file: %v\n", err)
		}
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			if al.logger != nil {
				al.logger.Warnw("Failed to close audit log file", "error", closeErr, "path", al.logPath)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to close audit log file: %v\n", closeErr)
			}
		}
	}()

	_, err = file.WriteString(entry + "\n")
	if err != nil {
		if al.logger != nil {
			al.logger.Errorw("Failed to write to audit log file", "error", err, "path", al.logPath)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to write to audit log file: %v\n", err)
		}
	}
}

// shouldLogEvent determines if an event should be logged based on severity
func (al *AuditLogger) shouldLogEvent(severity security.SecurityEventSeverity) bool {
	switch al.policy.AuditLogLevel {
	case "error":
		return severity == security.SeverityError || severity == security.SeverityCritical
	case "warn":
		return severity == security.SeverityWarning || severity == security.SeverityError || severity == security.SeverityCritical
	case "info":
		return true
	default:
		return true
	}
}

// GetAuditLogPath returns the path to the audit log file
func (al *AuditLogger) GetAuditLogPath() string {
	return al.logPath
}

// CleanupOldLogs removes audit log entries older than the retention period
func (al *AuditLogger) CleanupOldLogs() error {
	if al.policy.RetentionDays <= 0 {
		return nil
	}

	cutoffTime := time.Now().AddDate(0, 0, -al.policy.RetentionDays)

	// Read all log entries
	file, err := os.Open(al.logPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			if al.logger != nil {
				al.logger.Warnw("Failed to close audit log file during cleanup", "error", closeErr, "path", al.logPath)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to close audit log file: %v\n", closeErr)
			}
		}
	}()

	var validEntries []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse timestamp from log entry
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			continue
		}

		if timestampStr, ok := logEntry["timestamp"].(string); ok {
			if timestamp, err := time.Parse(time.RFC3339, timestampStr); err == nil {
				if timestamp.After(cutoffTime) {
					validEntries = append(validEntries, line)
				}
			}
		}
	}

	// Write back only valid entries
	return al.writeValidEntries(validEntries)
}

// writeValidEntries writes the valid log entries back to the file
func (al *AuditLogger) writeValidEntries(entries []string) error {
	file, err := os.OpenFile(al.logPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			if al.logger != nil {
				al.logger.Warnw("Failed to close audit log file during write", "error", closeErr, "path", al.logPath)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to close audit log file: %v\n", closeErr)
			}
		}
	}()

	for _, entry := range entries {
		_, err := file.WriteString(entry + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// StartBackgroundLogging starts the background logging worker
func (al *AuditLogger) StartBackgroundLogging() {
	al.startOnce.Do(func() {
		al.started = true
		al.wg.Add(1)
		go al.backgroundWorker()
	})
}

// StopBackgroundLogging stops the background logging worker
func (al *AuditLogger) StopBackgroundLogging() {
	al.stopOnce.Do(func() {
		al.mutex.Lock()
		defer al.mutex.Unlock()

		if !al.started {
			return
		}

		al.cancel()
		al.wg.Wait()
		al.started = false
	})
}

// LogEventAsync logs a security event asynchronously
func (al *AuditLogger) LogEventAsync(event *security.SecurityEvent) {
	if !al.policy.EnableAuditLog {
		return
	}

	// Ensure background logging is started
	al.mutex.RLock()
	started := al.started
	al.mutex.RUnlock()

	if !started {
		al.StartBackgroundLogging()
	}

	// Non-blocking send to queue
	select {
	case al.eventQueue <- event:
		// Event queued successfully
	default:
		// Queue is full, fall back to synchronous logging
		al.LogEvent(event)
	}
}

// backgroundWorker processes events from the queue
func (al *AuditLogger) backgroundWorker() {
	defer al.wg.Done()

	for {
		select {
		case event := <-al.eventQueue:
			al.LogEvent(event)
		case <-al.ctx.Done():
			// Process remaining events before shutdown
			for {
				select {
				case event := <-al.eventQueue:
					al.LogEvent(event)
				default:
					return
				}
			}
		}
	}
}

// GetQueueStats returns statistics about the event queue
func (al *AuditLogger) GetQueueStats() map[string]interface{} {
	al.mutex.RLock()
	defer al.mutex.RUnlock()

	return map[string]interface{}{
		"queue_length":              len(al.eventQueue),
		"queue_capacity":            cap(al.eventQueue),
		"background_logging_active": al.started,
	}
}
