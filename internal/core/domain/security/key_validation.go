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
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"
)

// KeyValidationResult represents the result of key validation
type KeyValidationResult struct {
	IsValid         bool     `json:"is_valid"`
	Issues          []string `json:"issues"`
	Warnings        []string `json:"warnings"`
	KeyInfo         *KeyInfo `json:"key_info,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`
}

// KeyInfo contains information about a validated key
type KeyInfo struct {
	Type        string    `json:"type"`        // rsa, ed25519, ecdsa
	Size        int       `json:"size"`        // Key size in bits
	Fingerprint string    `json:"fingerprint"` // Key fingerprint
	Created     time.Time `json:"created,omitempty"`
	Expires     time.Time `json:"expires,omitempty"`
	Comment     string    `json:"comment,omitempty"`
}

// KeyValidator provides key validation functionality
type KeyValidator struct {
	policy *SecurityPolicy
}

// NewKeyValidator creates a new key validator with the given security policy
func NewKeyValidator(policy *SecurityPolicy) *KeyValidator {
	return &KeyValidator{
		policy: policy,
	}
}

// ValidateKey validates an SSH key against the security policy
func (kv *KeyValidator) ValidateKey(keyData string) *KeyValidationResult {
	result := &KeyValidationResult{
		IsValid:         true,
		Issues:          []string{},
		Warnings:        []string{},
		Recommendations: []string{},
	}

	// Parse the key
	keyInfo, err := kv.parseKey(keyData)
	if err != nil {
		result.IsValid = false
		result.Issues = append(result.Issues, fmt.Sprintf("Failed to parse key: %v", err))
		return result
	}

	result.KeyInfo = keyInfo

	// Validate key type
	if !kv.isAllowedKeyType(keyInfo.Type) {
		result.IsValid = false
		result.Issues = append(result.Issues, fmt.Sprintf("Key type '%s' is not allowed. Allowed types: %s",
			keyInfo.Type, strings.Join(kv.policy.AllowedKeyTypes, ", ")))
	}

	// Validate key size
	if keyInfo.Size < kv.policy.MinKeySize {
		result.IsValid = false
		result.Issues = append(result.Issues, fmt.Sprintf("Key size %d bits is below minimum required size of %d bits",
			keyInfo.Size, kv.policy.MinKeySize))
	}

	// Check for weak key sizes
	const rsaKeyType = "rsa"
	if keyInfo.Type == rsaKeyType && keyInfo.Size < 3072 {
		result.Warnings = append(result.Warnings, "RSA key size is below recommended 3072 bits for new keys")
		result.Recommendations = append(result.Recommendations, "Consider using a 3072-bit or larger RSA key, or switch to Ed25519")
	}

	// Check key expiry
	if !keyInfo.Expires.IsZero() && time.Until(keyInfo.Expires) < kv.policy.KeyExpiryWarning {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Key expires in %v", time.Until(keyInfo.Expires)))
		result.Recommendations = append(result.Recommendations, "Consider renewing the key before it expires")
	}

	// Add general recommendations
	if keyInfo.Type == rsaKeyType {
		result.Recommendations = append(result.Recommendations, "Consider using Ed25519 keys for better security and performance")
	}

	return result
}

// parseKey parses an SSH key and extracts information
func (kv *KeyValidator) parseKey(keyData string) (*KeyInfo, error) {
	// Remove any whitespace and newlines
	keyData = strings.TrimSpace(keyData)

	// Parse PEM block
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	keyInfo := &KeyInfo{}

	// Determine key type and size based on PEM type
	switch block.Type {
	case "RSA PRIVATE KEY", "RSA PUBLIC KEY":
		keyInfo.Type = "rsa"
		if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
			keyInfo.Size = key.N.BitLen()
		} else if key, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
			keyInfo.Size = key.N.BitLen()
		} else {
			return nil, fmt.Errorf("failed to parse RSA key")
		}
	case "PRIVATE KEY", "PUBLIC KEY":
		// Try to parse as generic private key
		if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			if rsaKey, ok := key.(*rsa.PrivateKey); ok {
				keyInfo.Type = "rsa"
				keyInfo.Size = rsaKey.N.BitLen()
			} else {
				keyInfo.Type = "unknown"
			}
		} else {
			return nil, fmt.Errorf("failed to parse private key")
		}
	case "OPENSSH PRIVATE KEY":
		// This is more complex and would require additional parsing
		// For now, we'll assume it's a valid OpenSSH key
		keyInfo.Type = "openssh"
		keyInfo.Size = 0 // Would need to parse the key to get size
	default:
		return nil, fmt.Errorf("unsupported key type: %s", block.Type)
	}

	// Extract comment if present
	if comment := strings.TrimSpace(block.Headers["Comment"]); comment != "" {
		keyInfo.Comment = comment
	}

	return keyInfo, nil
}

// isAllowedKeyType checks if the key type is allowed by the policy
func (kv *KeyValidator) isAllowedKeyType(keyType string) bool {
	for _, allowedType := range kv.policy.AllowedKeyTypes {
		if keyType == allowedType {
			return true
		}
	}
	return false
}
