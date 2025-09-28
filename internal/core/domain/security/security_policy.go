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

// SecurityPolicy defines the security configuration for Wooak
type SecurityPolicy struct {
	// Key validation settings
	MinKeySize       int           `json:"min_key_size"`       // Minimum RSA key size (bits)
	AllowedKeyTypes  []string      `json:"allowed_key_types"`  // Allowed key types (rsa, ed25519, ecdsa)
	KeyExpiryWarning time.Duration `json:"key_expiry_warning"` // Warning before key expires

	// Connection security
	RequireHostKeyCheck bool `json:"require_host_key_check"` // Require host key verification
	MaxConnectionTime   int  `json:"max_connection_time"`    // Max connection time in minutes

	// Audit settings
	EnableAuditLog bool   `json:"enable_audit_log"` // Enable audit logging
	AuditLogLevel  string `json:"audit_log_level"`  // Audit log level (info, warn, error)
	RetentionDays  int    `json:"retention_days"`   // Log retention in days

	// Password policy
	RequirePasswordAuth bool `json:"require_password_auth"` // Require password authentication
	MinPasswordLength   int  `json:"min_password_length"`   // Minimum password length

	// Network security
	AllowedHosts []string `json:"allowed_hosts"` // Whitelist of allowed hosts
	BlockedHosts []string `json:"blocked_hosts"` // Blacklist of blocked hosts
	RequireVPN   bool     `json:"require_vpn"`   // Require VPN connection
}

// DefaultSecurityPolicy returns the default security policy
func DefaultSecurityPolicy() *SecurityPolicy {
	return &SecurityPolicy{
		MinKeySize:          2048,
		AllowedKeyTypes:     []string{"rsa", "ed25519", "ecdsa"},
		KeyExpiryWarning:    30 * 24 * time.Hour, // 30 days
		RequireHostKeyCheck: true,
		MaxConnectionTime:   60, // 60 minutes
		EnableAuditLog:      true,
		AuditLogLevel:       "info",
		RetentionDays:       90,
		RequirePasswordAuth: false,
		MinPasswordLength:   8,
		AllowedHosts:        []string{},
		BlockedHosts:        []string{},
		RequireVPN:          false,
	}
}
