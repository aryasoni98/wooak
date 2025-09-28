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

package ai

import (
	"context"
	"testing"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain/ai"
)

func TestAIService_GenerateRecommendationAsync(t *testing.T) {
	config := &ai.AIConfig{
		Enabled:  false, // Disabled for testing
		BaseURL:  "http://localhost:11434",
		Model:    "llama3.2:3b",
		Timeout:  5 * time.Second,
		CacheTTL: 1 * time.Hour,
	}

	service := NewAIService(config)
	ctx := context.Background()

	// Test async recommendation generation
	resultChan := service.GenerateRecommendationAsync(ctx, map[string]string{"host": "example.com"}, "test context")

	// Wait for result with timeout
	select {
	case result := <-resultChan:
		if result.Error == nil {
			t.Error("Expected error for disabled AI service")
		}
	case <-time.After(2 * time.Second):
		t.Error("Async operation timed out")
	}
}

func TestAIService_NaturalLanguageSearchAsync(t *testing.T) {
	config := &ai.AIConfig{
		Enabled:  false, // Disabled for testing
		BaseURL:  "http://localhost:11434",
		Model:    "llama3.2:3b",
		Timeout:  5 * time.Second,
		CacheTTL: 1 * time.Hour,
	}

	service := NewAIService(config)
	ctx := context.Background()

	// Test async natural language search
	servers := []map[string]string{
		{"host": "server1.example.com", "user": "admin"},
		{"host": "server2.example.com", "user": "root"},
	}

	resultChan := service.NaturalLanguageSearchAsync(ctx, "find production servers", servers)

	// Wait for result with timeout
	select {
	case result := <-resultChan:
		if result.Error == nil {
			t.Error("Expected error for disabled AI service")
		}
	case <-time.After(2 * time.Second):
		t.Error("Async operation timed out")
	}
}

func TestAIService_AnalyzeSecurityAsync(t *testing.T) {
	config := &ai.AIConfig{
		Enabled:  false, // Disabled for testing
		BaseURL:  "http://localhost:11434",
		Model:    "llama3.2:3b",
		Timeout:  5 * time.Second,
		CacheTTL: 1 * time.Hour,
	}

	service := NewAIService(config)
	ctx := context.Background()

	// Test async security analysis
	configMap := map[string]string{
		"host": "example.com",
		"user": "admin",
		"port": "22",
	}

	resultChan := service.AnalyzeSecurityAsync(ctx, configMap)

	// Wait for result with timeout
	select {
	case result := <-resultChan:
		if result.Error == nil {
			t.Error("Expected error for disabled AI service")
		}
	case <-time.After(2 * time.Second):
		t.Error("Async operation timed out")
	}
}

func TestAIService_TestConnectionAsync(t *testing.T) {
	config := &ai.AIConfig{
		Enabled:  false, // Disabled for testing
		BaseURL:  "http://localhost:11434",
		Model:    "llama3.2:3b",
		Timeout:  5 * time.Second,
		CacheTTL: 1 * time.Hour,
	}

	service := NewAIService(config)
	ctx := context.Background()

	// Test async connection test
	resultChan := service.TestConnectionAsync(ctx)

	// Wait for result with timeout
	select {
	case result := <-resultChan:
		if result.Error == nil {
			t.Error("Expected error for disabled AI service")
		}
		if result.Result {
			t.Error("Expected connection test to fail for disabled service")
		}
	case <-time.After(2 * time.Second):
		t.Error("Async operation timed out")
	}
}

func TestAIService_BatchProcessAsync(t *testing.T) {
	config := &ai.AIConfig{
		Enabled:  false, // Disabled for testing
		BaseURL:  "http://localhost:11434",
		Model:    "llama3.2:3b",
		Timeout:  5 * time.Second,
		CacheTTL: 1 * time.Hour,
	}

	service := NewAIService(config)
	ctx := context.Background()

	// Create test requests
	requests := []ai.AIRequest{
		{
			ID:        "req1",
			Type:      ai.RequestTypeGeneral,
			Prompt:    "Test prompt 1",
			Model:     "llama3.2:3b",
			MaxTokens: 100,
			Timestamp: time.Now(),
		},
		{
			ID:        "req2",
			Type:      ai.RequestTypeGeneral,
			Prompt:    "Test prompt 2",
			Model:     "llama3.2:3b",
			MaxTokens: 100,
			Timestamp: time.Now(),
		},
	}

	// Test async batch processing
	resultChan := service.BatchProcessAsync(ctx, requests)

	// Wait for result with timeout
	select {
	case results := <-resultChan:
		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}

		for i, result := range results {
			if result.Error == nil {
				t.Errorf("Expected error for request %d", i)
			}
		}
	case <-time.After(3 * time.Second):
		t.Error("Async batch operation timed out")
	}
}

func TestAIService_AsyncConcurrency(t *testing.T) {
	config := &ai.AIConfig{
		Enabled:  false, // Disabled for testing
		BaseURL:  "http://localhost:11434",
		Model:    "llama3.2:3b",
		Timeout:  5 * time.Second,
		CacheTTL: 1 * time.Hour,
	}

	service := NewAIService(config)
	ctx := context.Background()

	// Start multiple async operations concurrently
	const numOperations = 5
	resultChans := make([]<-chan AsyncResult[*ai.AIRecommendation], numOperations)

	for i := 0; i < numOperations; i++ {
		resultChans[i] = service.GenerateRecommendationAsync(ctx,
			map[string]string{"host": "example.com"}, "test context")
	}

	// Wait for all operations to complete
	completed := 0
	for _, resultChan := range resultChans {
		select {
		case result := <-resultChan:
			if result.Error == nil {
				t.Error("Expected error for disabled AI service")
			}
			completed++
		case <-time.After(3 * time.Second):
			t.Error("Async operation timed out")
		}
	}

	if completed != numOperations {
		t.Errorf("Expected %d operations to complete, got %d", numOperations, completed)
	}
}

func TestAsyncResult_GenericType(t *testing.T) {
	// Test that AsyncResult works with different types
	stringResult := AsyncResult[string]{
		Result: "test",
		Error:  nil,
	}

	if stringResult.Result != "test" {
		t.Error("Expected string result to be 'test'")
	}

	intResult := AsyncResult[int]{
		Result: 42,
		Error:  nil,
	}

	if intResult.Result != 42 {
		t.Error("Expected int result to be 42")
	}

	boolResult := AsyncResult[bool]{
		Result: true,
		Error:  nil,
	}

	if !boolResult.Result {
		t.Error("Expected bool result to be true")
	}
}
