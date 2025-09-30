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
	"fmt"
	"os"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain/security"
)

// SecurityService provides security-related functionality
type SecurityService struct {
	policy    *security.SecurityPolicy
	validator *security.KeyValidator
	auditLog  *AuditLogger
	keyCache  *KeyValidationCache
}

// NewSecurityService creates a new security service
func NewSecurityService(policy *security.SecurityPolicy) *SecurityService {
	validator := security.NewKeyValidator(policy)
	auditLog := NewAuditLogger(policy)
	keyCache := NewKeyValidationCache(100, 1*time.Hour) // Cache 100 keys for 1 hour

	return &SecurityService{
		policy:    policy,
		validator: validator,
		auditLog:  auditLog,
		keyCache:  keyCache,
	}
}

// ValidateSSHKey validates an SSH key against the security policy
func (s *SecurityService) ValidateSSHKey(keyData string) *security.KeyValidationResult {
	// Check cache first
	if cached, exists := s.keyCache.Get(keyData); exists {
		return cached
	}

	// Validate and set the key in the cache
	result := s.validator.ValidateKey(keyData)
	s.keyCache.Set(keyData, result)

	if !result.IsValid {
		s.auditLog.LogEventAsync(security.NewSecurityEvent(
			security.EventTypeKeyValidation,
			security.SeverityWarning,
			"SSH key validation failed",
		).WithSource("security_service").WithDetails("is_valid", result.IsValid))
	}

	return result
}

// ValidateSSHConfig validates the SSH configuration for security issues
func (s *SecurityService) ValidateSSHConfig(configPath string) *SecurityConfigValidationResult {
	result := &SecurityConfigValidationResult{
		IsValid:         true,
		Issues:          []string{},
		Warnings:        []string{},
		Recommendations: []string{},
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		result.Warnings = append(result.Warnings, "SSH config file does not exist")
		return result
	}

	// Check file permissions
	info, err := os.Stat(configPath)
	if err != nil {
		result.Issues = append(result.Issues, fmt.Sprintf("Cannot read SSH config file: %v", err))
		result.IsValid = false
		return result
	}

	// Check if file permissions are too permissive
	mode := info.Mode()
	if mode&0o077 != 0 {
		result.Issues = append(result.Issues, "SSH config file has overly permissive permissions")
		result.IsValid = false
	}

	// Log the validation
	s.auditLog.LogEvent(security.NewSecurityEvent(
		security.EventTypeConfigChange,
		security.SeverityInfo,
		"SSH config validation completed",
	).WithSource("security_service").WithDetails("config_path", configPath))

	return result
}

// CheckHostSecurity checks if a host meets security requirements
func (s *SecurityService) CheckHostSecurity(host string) *HostSecurityResult {
	result := &HostSecurityResult{
		Host:            host,
		IsSecure:        true,
		Issues:          []string{},
		Warnings:        []string{},
		Recommendations: []string{},
	}

	// Check if host is in blocked list
	for _, blockedHost := range s.policy.BlockedHosts {
		if host == blockedHost {
			result.IsSecure = false
			result.Issues = append(result.Issues, "Host is in blocked hosts list")
			break
		}
	}

	// Check if host is in allowed list (if allowed list is not empty)
	if len(s.policy.AllowedHosts) > 0 {
		allowed := false
		for _, allowedHost := range s.policy.AllowedHosts {
			if host == allowedHost {
				allowed = true
				break
			}
		}
		if !allowed {
			result.IsSecure = false
			result.Issues = append(result.Issues, "Host is not in allowed hosts list")
		}
	}

	// Log the security check
	severity := security.SeverityInfo
	if !result.IsSecure {
		severity = security.SeverityWarning
	}

	s.auditLog.LogEvent(security.NewSecurityEvent(
		security.EventTypeConnection,
		severity,
		fmt.Sprintf("Host security check completed for %s", host),
	).WithSource("security_service").WithHost(host).WithDetails("is_secure", result.IsSecure))

	return result
}

// GetSecurityPolicy returns the current security policy
func (s *SecurityService) GetSecurityPolicy() *security.SecurityPolicy {
	return s.policy
}

// UpdateSecurityPolicy updates the security policy
func (s *SecurityService) UpdateSecurityPolicy(newPolicy *security.SecurityPolicy) error {
	// Log the policy update
	s.auditLog.LogEvent(security.NewSecurityEvent(
		security.EventTypeConfigChange,
		security.SeverityInfo,
		"Security policy updated",
	).WithSource("security_service"))

	s.policy = newPolicy
	s.validator = security.NewKeyValidator(newPolicy)
	s.auditLog = NewAuditLogger(newPolicy)
	// Clear cache when policy changes as validation results may change
	s.keyCache.Clear()

	return nil
}

// GetCacheStats returns cache statistics
func (s *SecurityService) GetCacheStats() map[string]interface{} {
	return s.keyCache.Stats()
}

// ClearCache clears the key validation cache
func (s *SecurityService) ClearCache() {
	s.keyCache.Clear()
}

// CleanupCache removes expired entries from the cache
func (s *SecurityService) CleanupCache() {
	s.keyCache.Cleanup()
}

// StartBackgroundLogging starts background audit logging
func (s *SecurityService) StartBackgroundLogging() {
	s.auditLog.StartBackgroundLogging()
}

// StopBackgroundLogging stops background audit logging
func (s *SecurityService) StopBackgroundLogging() {
	s.auditLog.StopBackgroundLogging()
}

// GetAuditLogStats returns audit logging statistics
func (s *SecurityService) GetAuditLogStats() map[string]interface{} {
	return s.auditLog.GetQueueStats()
}

// SecurityConfigValidationResult represents the result of SSH config validation
type SecurityConfigValidationResult struct {
	IsValid         bool     `json:"is_valid"`
	Issues          []string `json:"issues"`
	Warnings        []string `json:"warnings"`
	Recommendations []string `json:"recommendations"`
}

// HostSecurityResult represents the result of host security check
type HostSecurityResult struct {
	Host            string   `json:"host"`
	IsSecure        bool     `json:"is_secure"`
	Issues          []string `json:"issues"`
	Warnings        []string `json:"warnings"`
	Recommendations []string `json:"recommendations"`
}
