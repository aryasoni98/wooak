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

	securityDomain "github.com/aryasoni98/wooak/internal/core/domain/security"
)

func TestSecurityService_NewSecurityService(t *testing.T) {
	policy := securityDomain.DefaultSecurityPolicy()
	service := NewSecurityService(policy)

	if service == nil {
		t.Fatal("Expected SecurityService to be created, got nil")
	}

	if service.policy != policy {
		t.Error("Expected policy to be set correctly")
	}

	if service.validator == nil {
		t.Error("Expected validator to be initialized")
	}
}

func TestSecurityService_GetSecurityPolicy(t *testing.T) {
	policy := securityDomain.DefaultSecurityPolicy()
	service := NewSecurityService(policy)

	retrievedPolicy := service.GetSecurityPolicy()
	if retrievedPolicy != policy {
		t.Error("Expected GetSecurityPolicy to return the same policy")
	}
}

func TestSecurityService_UpdateSecurityPolicy(t *testing.T) {
	policy := securityDomain.DefaultSecurityPolicy()
	service := NewSecurityService(policy)

	newPolicy := securityDomain.DefaultSecurityPolicy()
	newPolicy.MinKeySize = 4096
	newPolicy.RequireHostKeyCheck = true

	err := service.UpdateSecurityPolicy(newPolicy)
	if err != nil {
		t.Errorf("Expected UpdateSecurityPolicy to succeed, got error: %v", err)
	}

	if service.policy.MinKeySize != 4096 {
		t.Error("Expected policy to be updated")
	}

	if service.policy.RequireHostKeyCheck != true {
		t.Error("Expected RequireHostKeyCheck to be updated")
	}
}

func TestSecurityService_UpdateSecurityPolicy_Valid(t *testing.T) {
	policy := securityDomain.DefaultSecurityPolicy()
	service := NewSecurityService(policy)

	// Test updating with a valid policy
	newPolicy := securityDomain.DefaultSecurityPolicy()
	newPolicy.MinKeySize = 3072
	newPolicy.RequireHostKeyCheck = true

	err := service.UpdateSecurityPolicy(newPolicy)
	if err != nil {
		t.Errorf("Expected UpdateSecurityPolicy to succeed, got error: %v", err)
	}

	// Verify the policy was updated
	if service.policy.MinKeySize != 3072 {
		t.Error("Expected MinKeySize to be updated to 3072")
	}

	if service.policy.RequireHostKeyCheck != true {
		t.Error("Expected RequireHostKeyCheck to be updated to true")
	}
}

func TestSecurityService_ValidateSSHKey(t *testing.T) {
	policy := securityDomain.DefaultSecurityPolicy()
	service := NewSecurityService(policy)

	// Test with invalid key data
	invalidKey := "invalid-key-data"
	result := service.ValidateSSHKey(invalidKey)

	if result == nil {
		t.Fatal("Expected validation result, got nil")
	}

	if result.IsValid {
		t.Error("Expected invalid key to be marked as invalid")
	}

	if len(result.Issues) == 0 {
		t.Error("Expected issues to be reported for invalid key")
	}
}

func TestSecurityService_ValidateSSHKey_InvalidKey(t *testing.T) {
	policy := securityDomain.DefaultSecurityPolicy()
	service := NewSecurityService(policy)

	// Test with completely invalid key data
	invalidKey := "not-a-key-at-all"
	result := service.ValidateSSHKey(invalidKey)

	if result == nil {
		t.Fatal("Expected validation result, got nil")
	}

	if result.IsValid {
		t.Error("Expected invalid key to be marked as invalid")
	}

	if len(result.Issues) == 0 {
		t.Error("Expected issues to be reported for invalid key")
	}
}
