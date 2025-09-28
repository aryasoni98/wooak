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
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"
)

// generateRequestID generates a unique request ID
func generateRequestID() string {
	now := time.Now().UnixNano()
	randNum, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return fmt.Sprintf("req_%d_%d", now, randNum.Int64())
}

// generateResponseID generates a unique response ID
func generateResponseID() string {
	now := time.Now().UnixNano()
	randNum, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return fmt.Sprintf("resp_%d_%d", now, randNum.Int64())
}

// generateRecommendationID generates a unique recommendation ID
func generateRecommendationID() string {
	now := time.Now().UnixNano()
	randNum, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return fmt.Sprintf("rec_%d_%d", now, randNum.Int64())
}

// generateHash generates a hash for caching purposes
func generateHash(data interface{}) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%v", data)))
	return fmt.Sprintf("%x", hash)[:16] // Use first 16 characters
}

// IsOllamaRunning checks if Ollama is running on the default port
func IsOllamaRunning(baseURL string) bool {
	// This would typically make a simple HTTP request to check if Ollama is running
	// For now, we'll return true as a placeholder
	// In a real implementation, you'd make a HEAD request to the base URL
	return true
}

// GetAvailableModels returns a list of available AI models
func GetAvailableModels() []string {
	return []string{
		"llama3.2:3b",
		"llama3.2:1b",
		"llama3.1:8b",
		"llama3.1:70b",
		"codellama:7b",
		"codellama:13b",
		"mistral:7b",
		"gemma:2b",
		"gemma:7b",
	}
}

// ValidateModel checks if a model name is valid
func ValidateModel(model string) bool {
	availableModels := GetAvailableModels()
	for _, availableModel := range availableModels {
		if model == availableModel {
			return true
		}
	}
	return false
}

// GetModelInfo returns information about a specific model
func GetModelInfo(model string) map[string]interface{} {
	modelInfo := map[string]map[string]interface{}{
		"llama3.2:3b": {
			"size":        "2.0GB",
			"parameters":  "3B",
			"context":     "128K",
			"description": "Fast, efficient model for general tasks",
		},
		"llama3.2:1b": {
			"size":        "1.3GB",
			"parameters":  "1B",
			"context":     "128K",
			"description": "Ultra-lightweight model for simple tasks",
		},
		"llama3.1:8b": {
			"size":        "4.7GB",
			"parameters":  "8B",
			"context":     "128K",
			"description": "Balanced model for complex reasoning",
		},
		"llama3.1:70b": {
			"size":        "40GB",
			"parameters":  "70B",
			"context":     "128K",
			"description": "High-performance model for advanced tasks",
		},
		"codellama:7b": {
			"size":        "3.8GB",
			"parameters":  "7B",
			"context":     "100K",
			"description": "Specialized for code generation and analysis",
		},
		"codellama:13b": {
			"size":        "7.3GB",
			"parameters":  "13B",
			"context":     "100K",
			"description": "Advanced code model for complex programming tasks",
		},
		"mistral:7b": {
			"size":        "4.1GB",
			"parameters":  "7B",
			"context":     "32K",
			"description": "Efficient model with strong reasoning capabilities",
		},
		"gemma:2b": {
			"size":        "1.6GB",
			"parameters":  "2B",
			"context":     "8K",
			"description": "Google's lightweight, efficient model",
		},
		"gemma:7b": {
			"size":        "5.4GB",
			"parameters":  "7B",
			"context":     "8K",
			"description": "Google's balanced model for various tasks",
		},
	}

	if info, exists := modelInfo[model]; exists {
		return info
	}

	return map[string]interface{}{
		"size":        "Unknown",
		"parameters":  "Unknown",
		"context":     "Unknown",
		"description": "Custom or unknown model",
	}
}
