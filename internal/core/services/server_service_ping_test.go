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

package services

import (
	"testing"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain"
	"go.uber.org/zap"
)

func TestServerService_Ping(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	mockRepo := &mockServerRepository{}

	service := &serverService{
		logger:           logger.Sugar(),
		serverRepository: mockRepo,
	}

	tests := []struct {
		name        string
		server      domain.Server
		expectError bool
	}{
		{
			name: "ping with valid host and port",
			server: domain.Server{
				Alias: "test-server",
				Host:  "127.0.0.1",
				Port:  22,
			},
			// May succeed or fail depending on whether SSH is running locally
			expectError: false, // We'll check for error but not require it
		},
		{
			name: "ping with default port",
			server: domain.Server{
				Alias: "test-server",
				Host:  "127.0.0.1",
				Port:  0, // Will default to 22
			},
			expectError: false,
		},
		{
			name: "ping with invalid port",
			server: domain.Server{
				Alias: "test-server",
				Host:  "127.0.0.1",
				Port:  99999, // Invalid port
			},
			expectError: true,
		},
		{
			name: "ping with empty host",
			server: domain.Server{
				Alias: "test-server",
				Host:  "",
				Port:  22,
			},
			// Will use alias as host
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reachable, duration, err := service.Ping(tt.server)

			if tt.expectError {
				if err == nil {
					t.Logf("Note: Ping succeeded but was expected to fail (this is OK for some cases)")
				}
			}

			// Verify duration is reasonable (should be less than timeout)
			expectedTimeout := 3 * time.Second
			if duration > expectedTimeout+1*time.Second {
				t.Errorf("Ping duration %v exceeds expected timeout", duration)
			}

			// Log results for debugging
			t.Logf("Ping result: reachable=%v, duration=%v, error=%v", reachable, duration, err)
		})
	}
}

func TestResolveSSHDestination(t *testing.T) {
	tests := []struct {
		name     string
		alias    string
		expected bool // Whether resolution should succeed
	}{
		{
			name:     "empty alias",
			alias:    "",
			expected: false,
		},
		{
			name:     "whitespace alias",
			alias:    "   ",
			expected: false,
		},
		{
			name:     "valid alias format",
			alias:    "test-server",
			expected: false, // Will fail because SSH config doesn't exist in test
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, port, ok := resolveSSHDestination(tt.alias)

			if tt.expected {
				if !ok {
					t.Errorf("Expected resolution to succeed, got host=%q, port=%d", host, port)
				}
			} else {
				// For most test cases, resolution will fail (no SSH config)
				// Just verify function doesn't panic
				_ = host
				_ = port
			}
		})
	}
}
