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
	"time"
)

// AIProvider represents different AI providers
type AIProvider string

const (
	ProviderOllama AIProvider = "ollama"
	ProviderOpenAI AIProvider = "openai"
)

// AIModel represents an AI model configuration
type AIModel struct {
	Name        string     `json:"name"`        // Model name (e.g., "llama3.2:3b")
	Provider    AIProvider `json:"provider"`    // AI provider
	MaxTokens   int        `json:"max_tokens"`  // Maximum tokens for responses
	Temperature float64    `json:"temperature"` // Response creativity (0.0-1.0)
	Enabled     bool       `json:"enabled"`     // Whether model is enabled
}

// AIRequest represents a request to the AI service
type AIRequest struct {
	ID          string                 `json:"id"`
	Type        AIRequestType          `json:"type"`
	Prompt      string                 `json:"prompt"`
	Context     map[string]interface{} `json:"context,omitempty"`
	Model       string                 `json:"model,omitempty"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// AIRequestType represents the type of AI request
type AIRequestType string

const (
	RequestTypeRecommendation AIRequestType = "recommendation"
	RequestTypeSearch         AIRequestType = "search"
	RequestTypeOptimization   AIRequestType = "optimization"
	RequestTypeSecurity       AIRequestType = "security"
	RequestTypeGeneral        AIRequestType = "general"
)

// AIResponse represents a response from the AI service
type AIResponse struct {
	ID             string                 `json:"id"`
	RequestID      string                 `json:"request_id"`
	Content        string                 `json:"content"`
	Type           AIRequestType          `json:"type"`
	Confidence     float64                `json:"confidence"` // 0.0-1.0
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	Timestamp      time.Time              `json:"timestamp"`
	Model          string                 `json:"model"`
	TokensUsed     int                    `json:"tokens_used,omitempty"`
	ProcessingTime time.Duration          `json:"processing_time"`
}

// AIRecommendation represents an AI-generated recommendation
type AIRecommendation struct {
	ID          string                 `json:"id"`
	Type        RecommendationType     `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Confidence  float64                `json:"confidence"`
	Priority    Priority               `json:"priority"`
	Actions     []string               `json:"actions,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// RecommendationType represents the type of recommendation
type RecommendationType string

const (
	RecTypeSecurity     RecommendationType = "security"
	RecTypePerformance  RecommendationType = "performance"
	RecTypeOptimization RecommendationType = "optimization"
	RecTypeBestPractice RecommendationType = "best_practice"
	RecTypeConnection   RecommendationType = "connection"
)

// Priority represents the priority level of a recommendation
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// AIConfig represents the AI service configuration
type AIConfig struct {
	Provider     AIProvider    `json:"provider"`
	Model        string        `json:"model"`
	BaseURL      string        `json:"base_url,omitempty"`
	APIKey       string        `json:"api_key,omitempty"`
	MaxTokens    int           `json:"max_tokens"`
	Temperature  float64       `json:"temperature"`
	Timeout      time.Duration `json:"timeout"`
	Enabled      bool          `json:"enabled"`
	CacheEnabled bool          `json:"cache_enabled"`
	CacheTTL     time.Duration `json:"cache_ttl"`
}

// DefaultAIConfig returns the default AI configuration
func DefaultAIConfig() *AIConfig {
	return &AIConfig{
		Provider:     ProviderOllama,
		Model:        "llama3.2:3b",
		BaseURL:      "http://localhost:11434",
		MaxTokens:    512,
		Temperature:  0.7,
		Timeout:      30 * time.Second,
		Enabled:      true,
		CacheEnabled: true,
		CacheTTL:     1 * time.Hour,
	}
}
