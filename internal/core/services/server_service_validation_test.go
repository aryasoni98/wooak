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

	"github.com/aryasoni98/wooak/internal/core/domain"
	"go.uber.org/zap"
)

// mockServerRepository is a mock implementation for testing
type mockServerRepository struct {
	servers []domain.Server
	err     error
}

func (m *mockServerRepository) ListServers(query string) ([]domain.Server, error) {
	return m.servers, m.err
}

func (m *mockServerRepository) UpdateServer(server domain.Server, newServer domain.Server) error {
	return m.err
}

func (m *mockServerRepository) AddServer(server domain.Server) error {
	return m.err
}

func (m *mockServerRepository) DeleteServer(server domain.Server) error {
	return m.err
}

func (m *mockServerRepository) SetPinned(alias string, pinned bool) error {
	return m.err
}

func (m *mockServerRepository) RecordSSH(alias string) error {
	return m.err
}

func TestIsValidAlias(t *testing.T) {
	tests := []struct {
		name     string
		alias    string
		expected bool
	}{
		{
			name:     "valid simple alias",
			alias:    "server1",
			expected: true,
		},
		{
			name:     "valid alias with dots",
			alias:    "server.prod",
			expected: true,
		},
		{
			name:     "valid alias with dashes",
			alias:    "server-prod",
			expected: true,
		},
		{
			name:     "valid alias with underscores",
			alias:    "server_prod",
			expected: true,
		},
		{
			name:     "valid alias with numbers",
			alias:    "server123",
			expected: true,
		},
		{
			name:     "valid complex alias",
			alias:    "server-prod-123.aws",
			expected: true,
		},
		{
			name:     "empty alias",
			alias:    "",
			expected: false,
		},
		{
			name:     "alias too long",
			alias:    "a" + string(make([]byte, 101)),
			expected: false,
		},
		{
			name:     "alias with path traversal",
			alias:    "server/../etc",
			expected: false,
		},
		{
			name:     "alias with backslash",
			alias:    "server\\etc",
			expected: false,
		},
		{
			name:     "alias with semicolon",
			alias:    "server;rm -rf /",
			expected: false,
		},
		{
			name:     "alias with ampersand",
			alias:    "server&rm -rf /",
			expected: false,
		},
		{
			name:     "alias with pipe",
			alias:    "server|rm -rf /",
			expected: false,
		},
		{
			name:     "alias with backtick",
			alias:    "server`rm -rf /`",
			expected: false,
		},
		{
			name:     "alias with dollar sign",
			alias:    "server$HOME",
			expected: false,
		},
		{
			name:     "alias with parentheses",
			alias:    "server(rm -rf /)",
			expected: false,
		},
		{
			name:     "alias with angle brackets",
			alias:    "server<file",
			expected: false,
		},
		{
			name:     "alias with quotes",
			alias:    "server\"rm -rf /\"",
			expected: false,
		},
		{
			name:     "alias with single quotes",
			alias:    "server'rm -rf /'",
			expected: false,
		},
		{
			name:     "alias with newline",
			alias:    "server\nrm -rf /",
			expected: false,
		},
		{
			name:     "alias with carriage return",
			alias:    "server\rm -rf /",
			expected: false,
		},
		{
			name:     "alias with tab",
			alias:    "server\trm -rf /",
			expected: false,
		},
		{
			name:     "alias with double dots",
			alias:    "server..etc",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidAlias(tt.alias)
			if result != tt.expected {
				t.Errorf("isValidAlias(%q) = %v, expected %v", tt.alias, result, tt.expected)
			}
		})
	}
}

func TestServerService_ValidateSSHAccess(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name        string
		alias       string
		servers     []domain.Server
		repoErr     error
		expectError bool
	}{
		{
			name:  "valid alias exists in servers",
			alias: "server1",
			servers: []domain.Server{
				{Alias: "server1", Host: "192.168.1.1"},
				{Alias: "server2", Host: "192.168.1.2"},
			},
			expectError: false,
		},
		{
			name:  "alias not found in servers",
			alias: "nonexistent",
			servers: []domain.Server{
				{Alias: "server1", Host: "192.168.1.1"},
				{Alias: "server2", Host: "192.168.1.2"},
			},
			expectError: true,
		},
		{
			name:        "repository error",
			alias:       "server1",
			servers:     []domain.Server{},
			repoErr:     &MockError{message: "repository error"},
			expectError: true,
		},
		{
			name:        "empty servers list",
			alias:       "server1",
			servers:     []domain.Server{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockServerRepository{
				servers: tt.servers,
				err:     tt.repoErr,
			}

			service := &serverService{
				logger:           logger.Sugar(),
				serverRepository: mockRepo,
			}

			err := service.validateSSHAccess(tt.alias)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

func TestServerService_SSH_Validation(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name        string
		alias       string
		servers     []domain.Server
		expectError bool
		errorMsg    string
	}{
		{
			name:  "valid alias with existing server",
			alias: "server1",
			servers: []domain.Server{
				{Alias: "server1", Host: "192.168.1.1"},
			},
			expectError: true, // Will fail due to actual SSH execution
		},
		{
			name:        "invalid alias format",
			alias:       "server;rm -rf /",
			servers:     []domain.Server{},
			expectError: true,
			errorMsg:    "invalid alias format",
		},
		{
			name:  "valid alias format but server not found",
			alias: "nonexistent",
			servers: []domain.Server{
				{Alias: "server1", Host: "192.168.1.1"},
			},
			expectError: true,
			errorMsg:    "alias not found in known servers list", // Updated to match new error format
		},
		{
			name:        "empty alias",
			alias:       "",
			servers:     []domain.Server{},
			expectError: true,
			errorMsg:    "invalid alias format",
		},
		{
			name:        "alias with command injection attempt",
			alias:       "server`whoami`",
			servers:     []domain.Server{},
			expectError: true,
			errorMsg:    "invalid alias format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockServerRepository{
				servers: tt.servers,
			}

			service := &serverService{
				logger:           logger.Sugar(),
				serverRepository: mockRepo,
			}

			err := service.SSH(tt.alias)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

// Helper functions
type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || substr == "" ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestServerService_ListServers tests the ListServers business logic
func TestServerService_ListServers(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name     string
		query    string
		servers  []domain.Server
		expected int
	}{
		{
			name:  "list all servers",
			query: "",
			servers: []domain.Server{
				{Alias: "server1", Host: "192.168.1.1", Tags: []string{"prod"}},
				{Alias: "server2", Host: "192.168.1.2", Tags: []string{"dev"}},
				{Alias: "server3", Host: "192.168.1.3", Tags: []string{"prod"}},
			},
			expected: 3,
		},
		{
			name:  "filter by query",
			query: "server1",
			servers: []domain.Server{
				{Alias: "server1", Host: "192.168.1.1", Tags: []string{"prod"}},
			},
			expected: 1,
		},
		{
			name:  "filter by host",
			query: "192.168.1.1",
			servers: []domain.Server{
				{Alias: "server1", Host: "192.168.1.1", Tags: []string{"prod"}},
			},
			expected: 1,
		},
		{
			name:  "filter by tag",
			query: "prod",
			servers: []domain.Server{
				{Alias: "server1", Host: "192.168.1.1", Tags: []string{"prod"}},
				{Alias: "server3", Host: "192.168.1.3", Tags: []string{"prod"}},
			},
			expected: 2,
		},
		{
			name:     "empty servers list",
			query:    "",
			servers:  []domain.Server{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockServerRepository{
				servers: tt.servers,
			}

			service := &serverService{
				logger:           logger.Sugar(),
				serverRepository: mockRepo,
			}

			result, err := service.ListServers(tt.query)
			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if len(result) != tt.expected {
				t.Errorf("Expected %d servers, got %d", tt.expected, len(result))
			}
		})
	}
}

// TestServerService_AddServer tests the AddServer business logic
func TestServerService_AddServer(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name        string
		server      domain.Server
		repoErr     error
		expectError bool
	}{
		{
			name: "valid server",
			server: domain.Server{
				Alias: "server1",
				Host:  "192.168.1.1",
				User:  "admin",
				Port:  22,
			},
			expectError: false,
		},
		{
			name: "repository error",
			server: domain.Server{
				Alias: "server1",
				Host:  "192.168.1.1",
			},
			repoErr:     &MockError{message: "repository error"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockServerRepository{
				err: tt.repoErr,
			}

			service := &serverService{
				logger:           logger.Sugar(),
				serverRepository: mockRepo,
			}

			err := service.AddServer(tt.server)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

// TestServerService_UpdateServer tests the UpdateServer business logic
func TestServerService_UpdateServer(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name        string
		oldServer   domain.Server
		newServer   domain.Server
		repoErr     error
		expectError bool
	}{
		{
			name: "valid update",
			oldServer: domain.Server{
				Alias: "server1",
				Host:  "192.168.1.1",
			},
			newServer: domain.Server{
				Alias: "server1",
				Host:  "192.168.1.2",
			},
			expectError: false,
		},
		{
			name: "repository error",
			oldServer: domain.Server{
				Alias: "server1",
				Host:  "192.168.1.1",
			},
			newServer: domain.Server{
				Alias: "server1",
				Host:  "192.168.1.2",
			},
			repoErr:     &MockError{message: "repository error"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockServerRepository{
				err: tt.repoErr,
			}

			service := &serverService{
				logger:           logger.Sugar(),
				serverRepository: mockRepo,
			}

			err := service.UpdateServer(tt.oldServer, tt.newServer)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

// TestServerService_DeleteServer tests the DeleteServer business logic
func TestServerService_DeleteServer(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name        string
		server      domain.Server
		repoErr     error
		expectError bool
	}{
		{
			name: "valid deletion",
			server: domain.Server{
				Alias: "server1",
				Host:  "192.168.1.1",
			},
			expectError: false,
		},
		{
			name: "repository error",
			server: domain.Server{
				Alias: "server1",
				Host:  "192.168.1.1",
			},
			repoErr:     &MockError{message: "repository error"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockServerRepository{
				err: tt.repoErr,
			}

			service := &serverService{
				logger:           logger.Sugar(),
				serverRepository: mockRepo,
			}

			err := service.DeleteServer(tt.server)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

// TestServerService_SetPinned tests the SetPinned business logic
func TestServerService_SetPinned(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name        string
		alias       string
		pinned      bool
		repoErr     error
		expectError bool
	}{
		{
			name:        "pin server",
			alias:       "server1",
			pinned:      true,
			expectError: false,
		},
		{
			name:        "unpin server",
			alias:       "server1",
			pinned:      false,
			expectError: false,
		},
		{
			name:        "repository error",
			alias:       "server1",
			pinned:      true,
			repoErr:     &MockError{message: "repository error"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockServerRepository{
				err: tt.repoErr,
			}

			service := &serverService{
				logger:           logger.Sugar(),
				serverRepository: mockRepo,
			}

			err := service.SetPinned(tt.alias, tt.pinned)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}
