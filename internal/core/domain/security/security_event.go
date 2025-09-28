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
	"time"
)

// SecurityEventType represents the type of security event
type SecurityEventType string

const (
	EventTypeConnection         SecurityEventType = "connection"
	EventTypeKeyValidation      SecurityEventType = "key_validation"
	EventTypeConfigChange       SecurityEventType = "config_change"
	EventTypeAccessDenied       SecurityEventType = "access_denied"
	EventTypePolicyViolation    SecurityEventType = "policy_violation"
	EventTypeSuspiciousActivity SecurityEventType = "suspicious_activity"
)

// SecurityEventSeverity represents the severity level of a security event
type SecurityEventSeverity string

const (
	SeverityInfo     SecurityEventSeverity = "info"
	SeverityWarning  SecurityEventSeverity = "warning"
	SeverityError    SecurityEventSeverity = "error"
	SeverityCritical SecurityEventSeverity = "critical"
)

// SecurityEvent represents a security-related event in the system
type SecurityEvent struct {
	ID        string                 `json:"id"`
	Type      SecurityEventType      `json:"type"`
	Severity  SecurityEventSeverity  `json:"severity"`
	Timestamp time.Time              `json:"timestamp"`
	User      string                 `json:"user,omitempty"`
	Host      string                 `json:"host,omitempty"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Source    string                 `json:"source"` // Component that generated the event
	Action    string                 `json:"action"` // Action taken
	Result    string                 `json:"result"` // Result of the action
}

// NewSecurityEvent creates a new security event
func NewSecurityEvent(eventType SecurityEventType, severity SecurityEventSeverity, message string) *SecurityEvent {
	return &SecurityEvent{
		ID:        generateEventID(),
		Type:      eventType,
		Severity:  severity,
		Timestamp: time.Now(),
		Message:   message,
		Details:   make(map[string]interface{}),
	}
}

// WithUser sets the user for the security event
func (e *SecurityEvent) WithUser(user string) *SecurityEvent {
	e.User = user
	return e
}

// WithHost sets the host for the security event
func (e *SecurityEvent) WithHost(host string) *SecurityEvent {
	e.Host = host
	return e
}

// WithDetails adds details to the security event
func (e *SecurityEvent) WithDetails(key string, value interface{}) *SecurityEvent {
	e.Details[key] = value
	return e
}

// WithSource sets the source component for the security event
func (e *SecurityEvent) WithSource(source string) *SecurityEvent {
	e.Source = source
	return e
}

// WithAction sets the action for the security event
func (e *SecurityEvent) WithAction(action string) *SecurityEvent {
	e.Action = action
	return e
}

// WithResult sets the result for the security event
func (e *SecurityEvent) WithResult(result string) *SecurityEvent {
	e.Result = result
	return e
}

// generateEventID generates a unique event ID
func generateEventID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
