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
	"testing"
	"time"
)

func TestAIProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider AIProvider
		expected AIProvider
	}{
		{"ollama provider", ProviderOllama, "ollama"},
		{"openai provider", ProviderOpenAI, "openai"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.provider != tt.expected {
				t.Errorf("Provider = %v, want %v", tt.provider, tt.expected)
			}
		})
	}
}

func TestAIModel(t *testing.T) {
	model := AIModel{
		Name:        "llama3.2:3b",
		Provider:    ProviderOllama,
		MaxTokens:   512,
		Temperature: 0.7,
		Enabled:     true,
	}

	if model.Name != "llama3.2:3b" {
		t.Errorf("Name = %v, want llama3.2:3b", model.Name)
	}
	if model.Provider != ProviderOllama {
		t.Errorf("Provider = %v, want ollama", model.Provider)
	}
	if model.MaxTokens != 512 {
		t.Errorf("MaxTokens = %v, want 512", model.MaxTokens)
	}
	if model.Temperature != 0.7 {
		t.Errorf("Temperature = %v, want 0.7", model.Temperature)
	}
	if !model.Enabled {
		t.Error("Enabled should be true")
	}
}

func TestAIRequestType(t *testing.T) {
	tests := []struct {
		name        string
		requestType AIRequestType
		expected    AIRequestType
	}{
		{"recommendation type", RequestTypeRecommendation, "recommendation"},
		{"search type", RequestTypeSearch, "search"},
		{"optimization type", RequestTypeOptimization, "optimization"},
		{"security type", RequestTypeSecurity, "security"},
		{"general type", RequestTypeGeneral, "general"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.requestType != tt.expected {
				t.Errorf("RequestType = %v, want %v", tt.requestType, tt.expected)
			}
		})
	}
}

func TestAIRequest(t *testing.T) {
	now := time.Now()
	request := AIRequest{
		ID:          "req-123",
		Type:        RequestTypeSearch,
		Prompt:      "Find production servers",
		Context:     map[string]interface{}{"environment": "production"},
		Model:       "llama3.2:3b",
		MaxTokens:   512,
		Temperature: 0.7,
		Timestamp:   now,
	}

	if request.ID != "req-123" {
		t.Errorf("ID = %v, want req-123", request.ID)
	}
	if request.Type != RequestTypeSearch {
		t.Errorf("Type = %v, want search", request.Type)
	}
	if request.Prompt != "Find production servers" {
		t.Errorf("Prompt = %v, want Find production servers", request.Prompt)
	}
	if request.Model != "llama3.2:3b" {
		t.Errorf("Model = %v, want llama3.2:3b", request.Model)
	}
	if request.Context["environment"] != "production" {
		t.Errorf("Context[environment] = %v, want production", request.Context["environment"])
	}
}

func TestAIResponse(t *testing.T) {
	now := time.Now()
	response := AIResponse{
		ID:             "resp-123",
		RequestID:      "req-123",
		Content:        "Here are the production servers...",
		Type:           RequestTypeSearch,
		Confidence:     0.95,
		Metadata:       map[string]interface{}{"source": "ollama"},
		Timestamp:      now,
		Model:          "llama3.2:3b",
		TokensUsed:     150,
		ProcessingTime: 2 * time.Second,
	}

	if response.ID != "resp-123" {
		t.Errorf("ID = %v, want resp-123", response.ID)
	}
	if response.RequestID != "req-123" {
		t.Errorf("RequestID = %v, want req-123", response.RequestID)
	}
	if response.Confidence != 0.95 {
		t.Errorf("Confidence = %v, want 0.95", response.Confidence)
	}
	if response.TokensUsed != 150 {
		t.Errorf("TokensUsed = %v, want 150", response.TokensUsed)
	}
	if response.ProcessingTime != 2*time.Second {
		t.Errorf("ProcessingTime = %v, want 2s", response.ProcessingTime)
	}
}

func TestRecommendationType(t *testing.T) {
	tests := []struct {
		name     string
		recType  RecommendationType
		expected RecommendationType
	}{
		{"security recommendation", RecTypeSecurity, "security"},
		{"performance recommendation", RecTypePerformance, "performance"},
		{"optimization recommendation", RecTypeOptimization, "optimization"},
		{"best practice recommendation", RecTypeBestPractice, "best_practice"},
		{"connection recommendation", RecTypeConnection, "connection"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.recType != tt.expected {
				t.Errorf("RecommendationType = %v, want %v", tt.recType, tt.expected)
			}
		})
	}
}

func TestPriority(t *testing.T) {
	tests := []struct {
		name     string
		priority Priority
		expected Priority
	}{
		{"low priority", PriorityLow, "low"},
		{"medium priority", PriorityMedium, "medium"},
		{"high priority", PriorityHigh, "high"},
		{"critical priority", PriorityCritical, "critical"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.priority != tt.expected {
				t.Errorf("Priority = %v, want %v", tt.priority, tt.expected)
			}
		})
	}
}

func TestAIRecommendation(t *testing.T) {
	now := time.Now()
	recommendation := AIRecommendation{
		ID:          "rec-123",
		Type:        RecTypeSecurity,
		Title:       "Update SSH key",
		Description: "Your SSH key is using an outdated algorithm",
		Confidence:  0.85,
		Priority:    PriorityHigh,
		Actions:     []string{"Generate new key", "Update authorized_keys"},
		Metadata:    map[string]interface{}{"key_type": "rsa-1024"},
		Timestamp:   now,
	}

	if recommendation.ID != "rec-123" {
		t.Errorf("ID = %v, want rec-123", recommendation.ID)
	}
	if recommendation.Type != RecTypeSecurity {
		t.Errorf("Type = %v, want security", recommendation.Type)
	}
	if recommendation.Priority != PriorityHigh {
		t.Errorf("Priority = %v, want high", recommendation.Priority)
	}
	if len(recommendation.Actions) != 2 {
		t.Errorf("Actions length = %v, want 2", len(recommendation.Actions))
	}
}

func TestDefaultAIConfig(t *testing.T) {
	config := DefaultAIConfig()

	if config == nil {
		t.Fatal("DefaultAIConfig returned nil")
	}

	if config.Provider != ProviderOllama {
		t.Errorf("Provider = %v, want ollama", config.Provider)
	}
	if config.Model != "llama3.2:3b" {
		t.Errorf("Model = %v, want llama3.2:3b", config.Model)
	}
	if config.BaseURL != "http://localhost:11434" {
		t.Errorf("BaseURL = %v, want http://localhost:11434", config.BaseURL)
	}
	if config.MaxTokens != 512 {
		t.Errorf("MaxTokens = %v, want 512", config.MaxTokens)
	}
	if config.Temperature != 0.7 {
		t.Errorf("Temperature = %v, want 0.7", config.Temperature)
	}
	if config.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", config.Timeout)
	}
	if !config.Enabled {
		t.Error("Enabled should be true")
	}
	if !config.CacheEnabled {
		t.Error("CacheEnabled should be true")
	}
	if config.CacheTTL != 1*time.Hour {
		t.Errorf("CacheTTL = %v, want 1h", config.CacheTTL)
	}
}

func TestAIConfigCustomization(t *testing.T) {
	config := AIConfig{
		Provider:     ProviderOpenAI,
		Model:        "gpt-4",
		BaseURL:      "https://api.openai.com/v1",
		APIKey:       "sk-test-key",
		MaxTokens:    2048,
		Temperature:  0.5,
		Timeout:      60 * time.Second,
		Enabled:      true,
		CacheEnabled: false,
		CacheTTL:     30 * time.Minute,
	}

	if config.Provider != ProviderOpenAI {
		t.Errorf("Provider = %v, want openai", config.Provider)
	}
	if config.Model != "gpt-4" {
		t.Errorf("Model = %v, want gpt-4", config.Model)
	}
	if config.APIKey != "sk-test-key" {
		t.Errorf("APIKey = %v, want sk-test-key", config.APIKey)
	}
	if config.CacheEnabled {
		t.Error("CacheEnabled should be false")
	}
}
