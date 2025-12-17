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
	"fmt"
	"testing"

	"github.com/aryasoni98/wooak/internal/core/domain"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// mockRepository is a simple mock for benchmarking
type mockRepository struct {
	servers []domain.Server
}

func (m *mockRepository) ListServers(query string) ([]domain.Server, error) {
	return m.servers, nil
}

func (m *mockRepository) UpdateServer(server, newServer domain.Server) error {
	return nil
}

func (m *mockRepository) AddServer(server domain.Server) error {
	return nil
}

func (m *mockRepository) DeleteServer(server domain.Server) error {
	return nil
}

func (m *mockRepository) SetPinned(alias string, pinned bool) error {
	return nil
}

func (m *mockRepository) RecordSSH(alias string) error {
	return nil
}

// BenchmarkListServers benchmarks server listing
func BenchmarkListServers(b *testing.B) {
	// Create mock repository with sample servers
	servers := make([]domain.Server, 100)
	for i := 0; i < 100; i++ {
		servers[i] = domain.Server{
			Alias: fmt.Sprintf("server%d", i%10),
			Host:  "example.com",
			User:  "user",
			Port:  22,
			Tags:  []string{"production", "web"},
		}
	}

	mockRepo := &mockRepository{servers: servers}

	// Create minimal logger
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	logger, _ := config.Build()
	sugar := logger.Sugar()

	service := &serverService{
		serverRepository: mockRepo,
		logger:           sugar,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.ListServers("")
		if err != nil {
			b.Fatalf("Failed to list servers: %v", err)
		}
	}
}

// BenchmarkListServersWithQuery benchmarks server listing with query
func BenchmarkListServersWithQuery(b *testing.B) {
	servers := make([]domain.Server, 100)
	for i := 0; i < 100; i++ {
		servers[i] = domain.Server{
			Alias: "server" + string(rune('0'+i%10)),
			Host:  "example.com",
			User:  "user",
			Port:  22,
			Tags:  []string{"production", "web"},
		}
	}

	mockRepo := &mockRepository{servers: servers}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	logger, _ := config.Build()
	sugar := logger.Sugar()

	service := &serverService{
		serverRepository: mockRepo,
		logger:           sugar,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.ListServers("server1")
		if err != nil {
			b.Fatalf("Failed to list servers: %v", err)
		}
	}
}

// BenchmarkValidateServer benchmarks server validation
func BenchmarkValidateServer(b *testing.B) {
	server := domain.Server{
		Alias: "test-server",
		Host:  "example.com",
		User:  "user",
		Port:  22,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := validateServer(server)
		if err != nil {
			b.Fatalf("Validation failed: %v", err)
		}
	}
}
