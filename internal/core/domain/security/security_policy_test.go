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
	"testing"
	"time"
)

func TestDefaultSecurityPolicy(t *testing.T) {
	policy := DefaultSecurityPolicy()

	if policy == nil {
		t.Fatal("DefaultSecurityPolicy returned nil")
	}

	// Test key validation settings
	if policy.MinKeySize != 2048 {
		t.Errorf("MinKeySize = %v, want 2048", policy.MinKeySize)
	}

	expectedKeyTypes := []string{"rsa", "ed25519", "ecdsa"}
	if len(policy.AllowedKeyTypes) != len(expectedKeyTypes) {
		t.Errorf("AllowedKeyTypes length = %v, want %v", len(policy.AllowedKeyTypes), len(expectedKeyTypes))
	}

	for i, kt := range expectedKeyTypes {
		if policy.AllowedKeyTypes[i] != kt {
			t.Errorf("AllowedKeyTypes[%v] = %v, want %v", i, policy.AllowedKeyTypes[i], kt)
		}
	}

	if policy.KeyExpiryWarning != 30*24*time.Hour {
		t.Errorf("KeyExpiryWarning = %v, want 30 days", policy.KeyExpiryWarning)
	}

	// Test connection security
	if !policy.RequireHostKeyCheck {
		t.Error("RequireHostKeyCheck should be true")
	}

	if policy.MaxConnectionTime != 60 {
		t.Errorf("MaxConnectionTime = %v, want 60", policy.MaxConnectionTime)
	}

	// Test audit settings
	if !policy.EnableAuditLog {
		t.Error("EnableAuditLog should be true")
	}

	if policy.AuditLogLevel != "info" {
		t.Errorf("AuditLogLevel = %v, want info", policy.AuditLogLevel)
	}

	if policy.RetentionDays != 90 {
		t.Errorf("RetentionDays = %v, want 90", policy.RetentionDays)
	}

	// Test password policy
	if policy.RequirePasswordAuth {
		t.Error("RequirePasswordAuth should be false")
	}

	if policy.MinPasswordLength != 8 {
		t.Errorf("MinPasswordLength = %v, want 8", policy.MinPasswordLength)
	}

	// Test network security
	if len(policy.AllowedHosts) != 0 {
		t.Errorf("AllowedHosts should be empty, got %v", len(policy.AllowedHosts))
	}

	if len(policy.BlockedHosts) != 0 {
		t.Errorf("BlockedHosts should be empty, got %v", len(policy.BlockedHosts))
	}

	if policy.RequireVPN {
		t.Error("RequireVPN should be false")
	}
}

func TestSecurityPolicyCustomization(t *testing.T) {
	policy := SecurityPolicy{
		MinKeySize:          4096,
		AllowedKeyTypes:     []string{"ed25519"},
		KeyExpiryWarning:    60 * 24 * time.Hour,
		RequireHostKeyCheck: true,
		MaxConnectionTime:   30,
		EnableAuditLog:      true,
		AuditLogLevel:       "warn",
		RetentionDays:       180,
		RequirePasswordAuth: true,
		MinPasswordLength:   12,
		AllowedHosts:        []string{"*.example.com", "192.168.1.*"},
		BlockedHosts:        []string{"malicious.com"},
		RequireVPN:          true,
	}

	if policy.MinKeySize != 4096 {
		t.Errorf("MinKeySize = %v, want 4096", policy.MinKeySize)
	}

	if len(policy.AllowedKeyTypes) != 1 || policy.AllowedKeyTypes[0] != "ed25519" {
		t.Errorf("AllowedKeyTypes = %v, want [ed25519]", policy.AllowedKeyTypes)
	}

	if policy.KeyExpiryWarning != 60*24*time.Hour {
		t.Errorf("KeyExpiryWarning = %v, want 60 days", policy.KeyExpiryWarning)
	}

	if policy.MaxConnectionTime != 30 {
		t.Errorf("MaxConnectionTime = %v, want 30", policy.MaxConnectionTime)
	}

	if policy.AuditLogLevel != "warn" {
		t.Errorf("AuditLogLevel = %v, want warn", policy.AuditLogLevel)
	}

	if policy.RetentionDays != 180 {
		t.Errorf("RetentionDays = %v, want 180", policy.RetentionDays)
	}

	if !policy.RequirePasswordAuth {
		t.Error("RequirePasswordAuth should be true")
	}

	if policy.MinPasswordLength != 12 {
		t.Errorf("MinPasswordLength = %v, want 12", policy.MinPasswordLength)
	}

	if len(policy.AllowedHosts) != 2 {
		t.Errorf("AllowedHosts length = %v, want 2", len(policy.AllowedHosts))
	}

	if len(policy.BlockedHosts) != 1 {
		t.Errorf("BlockedHosts length = %v, want 1", len(policy.BlockedHosts))
	}

	if !policy.RequireVPN {
		t.Error("RequireVPN should be true")
	}
}

func TestSecurityPolicyMinKeySize(t *testing.T) {
	tests := []struct {
		name       string
		minKeySize int
	}{
		{"1024 bits", 1024},
		{"2048 bits", 2048},
		{"4096 bits", 4096},
		{"8192 bits", 8192},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := SecurityPolicy{
				MinKeySize: tt.minKeySize,
			}
			if policy.MinKeySize != tt.minKeySize {
				t.Errorf("MinKeySize = %v, want %v", policy.MinKeySize, tt.minKeySize)
			}
		})
	}
}

func TestSecurityPolicyAuditLevels(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"info level", "info"},
		{"warn level", "warn"},
		{"error level", "error"},
		{"debug level", "debug"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := SecurityPolicy{
				AuditLogLevel: tt.level,
			}
			if policy.AuditLogLevel != tt.level {
				t.Errorf("AuditLogLevel = %v, want %v", policy.AuditLogLevel, tt.level)
			}
		})
	}
}

func TestSecurityPolicyHostFiltering(t *testing.T) {
	policy := SecurityPolicy{
		AllowedHosts: []string{"prod.example.com", "staging.example.com"},
		BlockedHosts: []string{"bad.example.com", "malicious.com"},
	}

	if len(policy.AllowedHosts) != 2 {
		t.Errorf("AllowedHosts length = %v, want 2", len(policy.AllowedHosts))
	}

	if len(policy.BlockedHosts) != 2 {
		t.Errorf("BlockedHosts length = %v, want 2", len(policy.BlockedHosts))
	}

	// Verify allowed hosts
	expectedAllowed := map[string]bool{
		"prod.example.com":    true,
		"staging.example.com": true,
	}

	for _, host := range policy.AllowedHosts {
		if !expectedAllowed[host] {
			t.Errorf("Unexpected allowed host: %v", host)
		}
	}

	// Verify blocked hosts
	expectedBlocked := map[string]bool{
		"bad.example.com": true,
		"malicious.com":   true,
	}

	for _, host := range policy.BlockedHosts {
		if !expectedBlocked[host] {
			t.Errorf("Unexpected blocked host: %v", host)
		}
	}
}

func TestSecurityPolicyConnectionTimeout(t *testing.T) {
	tests := []struct {
		name              string
		maxConnectionTime int
		expected          int
	}{
		{"15 minutes", 15, 15},
		{"30 minutes", 30, 30},
		{"60 minutes", 60, 60},
		{"120 minutes", 120, 120},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := SecurityPolicy{
				MaxConnectionTime: tt.maxConnectionTime,
			}
			if policy.MaxConnectionTime != tt.expected {
				t.Errorf("MaxConnectionTime = %v, want %v", policy.MaxConnectionTime, tt.expected)
			}
		})
	}
}

func TestSecurityPolicyKeyExpiryWarning(t *testing.T) {
	tests := []struct {
		name             string
		keyExpiryWarning time.Duration
		expected         time.Duration
	}{
		{"7 days", 7 * 24 * time.Hour, 7 * 24 * time.Hour},
		{"30 days", 30 * 24 * time.Hour, 30 * 24 * time.Hour},
		{"90 days", 90 * 24 * time.Hour, 90 * 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy := SecurityPolicy{
				KeyExpiryWarning: tt.keyExpiryWarning,
			}
			if policy.KeyExpiryWarning != tt.expected {
				t.Errorf("KeyExpiryWarning = %v, want %v", policy.KeyExpiryWarning, tt.expected)
			}
		})
	}
}
