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

	aiDomain "github.com/aryasoni98/wooak/internal/core/domain/ai"
)

func TestAIService_NewAIService(t *testing.T) {
	config := aiDomain.DefaultAIConfig()
	service := NewAIService(config)

	if service == nil {
		t.Fatal("Expected AIService to be created, got nil")
	}

	if service.config != config {
		t.Error("Expected config to be set correctly")
	}

	if service.httpClient == nil {
		t.Error("Expected httpClient to be initialized")
	}

	if service.cache == nil {
		t.Error("Expected cache to be initialized")
	}
}

func TestAIService_GetConfig(t *testing.T) {
	config := aiDomain.DefaultAIConfig()
	service := NewAIService(config)

	retrievedConfig := service.GetConfig()
	if retrievedConfig != config {
		t.Error("Expected GetConfig to return the same config")
	}
}

func TestAIService_UpdateConfig(t *testing.T) {
	config := aiDomain.DefaultAIConfig()
	service := NewAIService(config)

	newConfig := aiDomain.DefaultAIConfig()
	newConfig.Model = "llama3.2:1b"
	newConfig.Temperature = 0.5

	err := service.UpdateConfig(newConfig)
	if err != nil {
		t.Errorf("Expected UpdateConfig to succeed, got error: %v", err)
	}

	if service.config.Model != "llama3.2:1b" {
		t.Error("Expected config to be updated")
	}
}

func TestAIService_GenerateRecommendation_Disabled(t *testing.T) {
	config := aiDomain.DefaultAIConfig()
	config.Enabled = false
	service := NewAIService(config)

	ctx := context.Background()
	configMap := map[string]string{
		"Host": "test-server",
		"User": "testuser",
	}

	_, err := service.GenerateRecommendation(ctx, configMap, "test context")
	if err == nil {
		t.Error("Expected error when AI service is disabled")
	}

	expectedError := "AI service is disabled"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestAIService_NaturalLanguageSearch_Disabled(t *testing.T) {
	config := aiDomain.DefaultAIConfig()
	config.Enabled = false
	service := NewAIService(config)

	ctx := context.Background()
	servers := []map[string]string{
		{"Host": "server1", "User": "user1"},
		{"Host": "server2", "User": "user2"},
	}

	_, err := service.NaturalLanguageSearch(ctx, "find production servers", servers)
	if err == nil {
		t.Error("Expected error when AI service is disabled")
	}

	expectedError := "AI service is disabled"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestAIService_AnalyzeSecurity_Disabled(t *testing.T) {
	config := aiDomain.DefaultAIConfig()
	config.Enabled = false
	service := NewAIService(config)

	ctx := context.Background()
	configMap := map[string]string{
		"Host": "test-server",
		"User": "testuser",
	}

	_, err := service.AnalyzeSecurity(ctx, configMap)
	if err == nil {
		t.Error("Expected error when AI service is disabled")
	}

	expectedError := "AI service is disabled"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestAIService_TestConnection_Disabled(t *testing.T) {
	config := aiDomain.DefaultAIConfig()
	config.Enabled = false
	service := NewAIService(config)

	ctx := context.Background()
	err := service.TestConnection(ctx)
	if err == nil {
		t.Error("Expected error when AI service is disabled")
	}

	expectedError := "AI service is disabled"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestAIService_TestConnection_Enabled(t *testing.T) {
	config := aiDomain.DefaultAIConfig()
	config.Enabled = true
	config.BaseURL = "http://localhost:11434"
	service := NewAIService(config)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// This test will fail if Ollama is not running, which is expected
	err := service.TestConnection(ctx)
	if err != nil {
		// Expected if Ollama is not running
		t.Logf("TestConnection failed as expected (Ollama not running): %v", err)
	}
}

func TestAIService_CacheIntegration(t *testing.T) {
	config := aiDomain.DefaultAIConfig()
	config.Enabled = true
	config.CacheEnabled = true
	config.CacheTTL = 1 * time.Minute
	service := NewAIService(config)

	// Test cache initialization
	if service.cache == nil {
		t.Error("Expected cache to be initialized")
	}

	// Test cache size
	initialSize := service.cache.Size()
	if initialSize != 0 {
		t.Errorf("Expected initial cache size to be 0, got %d", initialSize)
	}
}

func TestAIService_ModelValidation(t *testing.T) {
	// Test valid models
	validModels := []string{
		"llama3.2:3b",
		"llama3.2:1b",
		"llama3.1:8b",
		"codellama:7b",
		"mistral:7b",
	}

	for _, model := range validModels {
		if !ValidateModel(model) {
			t.Errorf("Expected model '%s' to be valid", model)
		}
	}

	// Test invalid models
	invalidModels := []string{
		"invalid-model",
		"llama3.2:999b",
		"",
		"unknown:model",
	}

	for _, model := range invalidModels {
		if ValidateModel(model) {
			t.Errorf("Expected model '%s' to be invalid", model)
		}
	}
}

func TestAIService_ModelInfo(t *testing.T) {
	modelInfo := GetModelInfo("llama3.2:3b")

	if modelInfo["size"] != "2.0GB" {
		t.Errorf("Expected size '2.0GB', got '%s'", modelInfo["size"])
	}

	if modelInfo["parameters"] != "3B" {
		t.Errorf("Expected parameters '3B', got '%s'", modelInfo["parameters"])
	}

	// Test unknown model
	unknownInfo := GetModelInfo("unknown-model")
	if unknownInfo["size"] != "Unknown" {
		t.Error("Expected unknown model to return 'Unknown' size")
	}
}

func TestAIService_AvailableModels(t *testing.T) {
	models := GetAvailableModels()

	if len(models) == 0 {
		t.Error("Expected at least one available model")
	}

	// Check for expected models
	expectedModels := []string{"llama3.2:3b", "llama3.2:1b", "llama3.1:8b"}
	for _, expected := range expectedModels {
		found := false
		for _, model := range models {
			if model == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected model '%s' to be in available models", expected)
		}
	}
}

func TestAIService_OllamaRunning(t *testing.T) {
	// Test with default URL
	running := IsOllamaRunning("http://localhost:11434")
	// This will be true or false depending on whether Ollama is actually running
	t.Logf("Ollama running status: %t", running)
}

func TestAIService_RequestIDGeneration(t *testing.T) {
	id1 := generateRequestID()
	id2 := generateRequestID()

	if id1 == id2 {
		t.Error("Expected request IDs to be unique")
	}

	if id1 == "" {
		t.Error("Expected request ID to be non-empty")
	}
}

func TestAIService_ResponseIDGeneration(t *testing.T) {
	id1 := generateResponseID()
	id2 := generateResponseID()

	if id1 == id2 {
		t.Error("Expected response IDs to be unique")
	}

	if id1 == "" {
		t.Error("Expected response ID to be non-empty")
	}
}

func TestAIService_RecommendationIDGeneration(t *testing.T) {
	id1 := generateRecommendationID()
	id2 := generateRecommendationID()

	if id1 == id2 {
		t.Error("Expected recommendation IDs to be unique")
	}

	if id1 == "" {
		t.Error("Expected recommendation ID to be non-empty")
	}
}

func TestAIService_HashGeneration(t *testing.T) {
	data1 := map[string]string{"key1": "value1"}
	data2 := map[string]string{"key1": "value2"}

	hash1 := generateHash(data1)
	hash2 := generateHash(data2)

	if hash1 == hash2 {
		t.Error("Expected different data to generate different hashes")
	}

	if len(hash1) != 16 {
		t.Errorf("Expected hash length to be 16, got %d", len(hash1))
	}
}
