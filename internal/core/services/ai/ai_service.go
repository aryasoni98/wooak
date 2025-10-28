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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain/ai"
	"github.com/aryasoni98/wooak/internal/core/services/monitoring"
)

// AIService provides AI-powered functionality for Wooak
type AIService struct {
	config     *ai.AIConfig
	httpClient *http.Client
	cache      *AICache
	pool       *ConnectionPool
	monitoring *monitoring.MonitoringService
}

// NewAIService creates a new AI service
func NewAIService(config *ai.AIConfig) *AIService {
	// Create connection pool configuration
	poolConfig := &PoolConfig{
		MaxConnections:        5,
		MaxIdleConns:          5,
		MaxConnsPerHost:       3,
		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		RequestTimeout:        config.Timeout,
	}

	return &AIService{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		cache:      NewAICache(config.CacheTTL),
		pool:       NewConnectionPool(poolConfig),
		monitoring: nil, // Will be set via SetMonitoring
	}
}

// SetMonitoring sets the monitoring service for the AI service
func (s *AIService) SetMonitoring(mon *monitoring.MonitoringService) {
	s.monitoring = mon
}

// GenerateRecommendation generates AI recommendations for SSH configurations
func (s *AIService) GenerateRecommendation(ctx context.Context, config map[string]string, context string) (*ai.AIRecommendation, error) {
	start := time.Now()

	if !s.config.Enabled {
		return nil, fmt.Errorf("AI service is disabled")
	}

	// Check cache first
	cacheKey := fmt.Sprintf("recommendation:%s", generateHash(config))
	if cached, exists := s.cache.Get(cacheKey); exists {
		if rec, ok := cached.(*ai.AIRecommendation); ok {
			if s.monitoring != nil {
				s.monitoring.RecordCacheHit("ai_recommendation")
				s.monitoring.RecordOperation("ai_recommendation", time.Since(start), true)
			}
			return rec, nil
		}
	}

	if s.monitoring != nil {
		s.monitoring.RecordCacheMiss("ai_recommendation")
	}

	// Build prompt
	template := ai.GetPromptForType(ai.RequestTypeRecommendation)
	variables := map[string]string{
		"config":  ai.FormatServerConfig(config),
		"context": context,
	}
	prompt := template.BuildPrompt(variables)

	// Make AI request
	request := &ai.AIRequest{
		ID:          generateRequestID(),
		Type:        ai.RequestTypeRecommendation,
		Prompt:      prompt,
		Model:       s.config.Model,
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
		Timestamp:   time.Now(),
	}

	response, err := s.makeAIRequest(ctx, request)
	if err != nil {
		if s.monitoring != nil {
			s.monitoring.RecordOperation("ai_recommendation", time.Since(start), false)
		}
		return nil, fmt.Errorf("failed to generate recommendation: %w", err)
	}

	// Parse response into recommendation
	recommendation := s.parseRecommendation(response, config)

	// Cache the result
	if s.config.CacheEnabled {
		s.cache.Set(cacheKey, recommendation)
	}

	if s.monitoring != nil {
		s.monitoring.RecordOperation("ai_recommendation", time.Since(start), true)
		s.monitoring.GetMetrics().IncrementCounter("ai_request_total", map[string]string{
			"type":   "recommendation",
			"status": "success",
		})
	}

	return recommendation, nil
}

// NaturalLanguageSearch performs AI-powered natural language search
func (s *AIService) NaturalLanguageSearch(ctx context.Context, query string, servers []map[string]string) ([]map[string]interface{}, error) {
	if !s.config.Enabled {
		return nil, fmt.Errorf("AI service is disabled")
	}

	// Check cache first
	cacheKey := fmt.Sprintf("search:%s", generateHash(query))
	if cached, exists := s.cache.Get(cacheKey); exists {
		if results, ok := cached.([]map[string]interface{}); ok {
			return results, nil
		}
	}

	// Build prompt
	template := ai.GetPromptForType(ai.RequestTypeSearch)
	variables := map[string]string{
		"query":   query,
		"servers": ai.FormatServerList(servers),
	}
	prompt := template.BuildPrompt(variables)

	// Make AI request
	request := &ai.AIRequest{
		ID:          generateRequestID(),
		Type:        ai.RequestTypeSearch,
		Prompt:      prompt,
		Model:       s.config.Model,
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
		Timestamp:   time.Now(),
	}

	response, err := s.makeAIRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to perform natural language search: %w", err)
	}

	// Parse response into search results
	results := s.parseSearchResults(response)

	// Cache the result
	if s.config.CacheEnabled {
		s.cache.Set(cacheKey, results)
	}

	return results, nil
}

// AnalyzeSecurity performs AI-powered security analysis
func (s *AIService) AnalyzeSecurity(ctx context.Context, config map[string]string) (*ai.AIRecommendation, error) {
	if !s.config.Enabled {
		return nil, fmt.Errorf("AI service is disabled")
	}

	// Check cache first
	cacheKey := fmt.Sprintf("security:%s", generateHash(config))
	if cached, exists := s.cache.Get(cacheKey); exists {
		if rec, ok := cached.(*ai.AIRecommendation); ok {
			return rec, nil
		}
	}

	// Build prompt
	template := ai.GetPromptForType(ai.RequestTypeSecurity)
	variables := map[string]string{
		"config": ai.FormatServerConfig(config),
	}
	prompt := template.BuildPrompt(variables)

	// Make AI request
	request := &ai.AIRequest{
		ID:          generateRequestID(),
		Type:        ai.RequestTypeSecurity,
		Prompt:      prompt,
		Model:       s.config.Model,
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
		Timestamp:   time.Now(),
	}

	response, err := s.makeAIRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze security: %w", err)
	}

	// Parse response into security recommendation
	recommendation := s.parseSecurityAnalysis(response, config)

	// Cache the result
	if s.config.CacheEnabled {
		s.cache.Set(cacheKey, recommendation)
	}

	return recommendation, nil
}

// makeAIRequest makes a request to the AI provider
func (s *AIService) makeAIRequest(ctx context.Context, request *ai.AIRequest) (*ai.AIResponse, error) {
	startTime := time.Now()

	var response *ai.AIResponse
	var err error

	switch s.config.Provider {
	case ai.ProviderOllama:
		response, err = s.makeOllamaRequest(ctx, request)
	case ai.ProviderOpenAI:
		response, err = s.makeOpenAIRequest(ctx, request)
	default:
		if s.monitoring != nil {
			s.monitoring.GetMetrics().IncrementCounter("ai_request_error_total", map[string]string{
				"provider": string(s.config.Provider),
				"reason":   "unsupported_provider",
			})
		}
		return nil, fmt.Errorf("unsupported AI provider: %s", s.config.Provider)
	}

	if err != nil {
		if s.monitoring != nil {
			s.monitoring.GetMetrics().IncrementCounter("ai_request_error_total", map[string]string{
				"provider": string(s.config.Provider),
				"type":     string(request.Type),
			})
		}
		return nil, err
	}

	response.ProcessingTime = time.Since(startTime)

	if s.monitoring != nil {
		s.monitoring.GetMetrics().RecordTimer("ai_request_duration_seconds", response.ProcessingTime, map[string]string{
			"provider": string(s.config.Provider),
			"type":     string(request.Type),
		})

		// Record token usage if available (accumulated total)
		if response.TokensUsed > 0 {
			s.monitoring.GetMetrics().AddToCounter("ai_tokens_used_total", float64(response.TokensUsed), map[string]string{
				"provider": string(s.config.Provider),
				"type":     string(request.Type),
			})
		}
	}

	return response, nil
}

// makeOllamaRequest makes a request to Ollama
func (s *AIService) makeOllamaRequest(ctx context.Context, request *ai.AIRequest) (*ai.AIResponse, error) {
	url := fmt.Sprintf("%s/api/generate", s.config.BaseURL)

	payload := map[string]interface{}{
		"model":  request.Model,
		"prompt": request.Prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": request.Temperature,
			"num_predict": request.MaxTokens,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Use connection pool for better performance
	client := s.pool.GetClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: Failed to close response body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI request failed with status: %d", resp.StatusCode)
	}

	var ollamaResponse struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ollamaResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ai.AIResponse{
		ID:        generateResponseID(),
		RequestID: request.ID,
		Content:   ollamaResponse.Response,
		Type:      request.Type,
		Model:     request.Model,
		Timestamp: time.Now(),
	}, nil
}

// makeOpenAIRequest makes a request to OpenAI
// NOTE: OpenAI integration is not yet implemented. This is a placeholder for future functionality.
// Currently, only Ollama provider is supported. To use AI features, please configure Ollama.
// See: https://ollama.ai for installation instructions.
func (s *AIService) makeOpenAIRequest(ctx context.Context, request *ai.AIRequest) (*ai.AIResponse, error) {
	return nil, fmt.Errorf("OpenAI integration not yet implemented - please use Ollama provider instead")
}

// parseRecommendation parses AI response into a recommendation
func (s *AIService) parseRecommendation(response *ai.AIResponse, config map[string]string) *ai.AIRecommendation {
	return &ai.AIRecommendation{
		ID:          generateRecommendationID(),
		Type:        ai.RecTypeOptimization,
		Title:       "AI-Generated Recommendation",
		Description: response.Content,
		Confidence:  0.8, // Default confidence
		Priority:    ai.PriorityMedium,
		Timestamp:   time.Now(),
		Metadata: map[string]interface{}{
			"model":       response.Model,
			"request_id":  response.RequestID,
			"config_hash": generateHash(config),
		},
	}
}

// parseSearchResults parses AI response into search results
func (s *AIService) parseSearchResults(response *ai.AIResponse) []map[string]interface{} {
	// Simple parsing - in a real implementation, this would be more sophisticated
	return []map[string]interface{}{
		{
			"content": response.Content,
			"score":   0.8,
			"model":   response.Model,
		},
	}
}

// parseSecurityAnalysis parses AI response into security analysis
func (s *AIService) parseSecurityAnalysis(response *ai.AIResponse, config map[string]string) *ai.AIRecommendation {
	return &ai.AIRecommendation{
		ID:          generateRecommendationID(),
		Type:        ai.RecTypeSecurity,
		Title:       "Security Analysis",
		Description: response.Content,
		Confidence:  0.9, // Higher confidence for security analysis
		Priority:    ai.PriorityHigh,
		Timestamp:   time.Now(),
		Metadata: map[string]interface{}{
			"model":       response.Model,
			"request_id":  response.RequestID,
			"config_hash": generateHash(config),
		},
	}
}

// GetConfig returns the current AI configuration
func (s *AIService) GetConfig() *ai.AIConfig {
	return s.config
}

// UpdateConfig updates the AI configuration
func (s *AIService) UpdateConfig(newConfig *ai.AIConfig) error {
	s.config = newConfig
	s.httpClient.Timeout = newConfig.Timeout
	s.cache = NewAICache(newConfig.CacheTTL)
	return nil
}

// TestConnection tests the connection to the AI provider
func (s *AIService) TestConnection(ctx context.Context) error {
	if !s.config.Enabled {
		return fmt.Errorf("AI service is disabled")
	}

	testRequest := &ai.AIRequest{
		ID:        generateRequestID(),
		Type:      ai.RequestTypeGeneral,
		Prompt:    "Hello, this is a test message. Please respond with 'OK'.",
		Model:     s.config.Model,
		MaxTokens: 10,
		Timestamp: time.Now(),
	}

	_, err := s.makeAIRequest(ctx, testRequest)
	return err
}

// AsyncResult represents the result of an async operation
type AsyncResult[T any] struct {
	Result T
	Error  error
}

// GenerateRecommendationAsync generates AI recommendations asynchronously
func (s *AIService) GenerateRecommendationAsync(ctx context.Context, config map[string]string, context string) <-chan AsyncResult[*ai.AIRecommendation] {
	result := make(chan AsyncResult[*ai.AIRecommendation], 1)

	go func() {
		defer close(result)
		recommendation, err := s.GenerateRecommendation(ctx, config, context)
		result <- AsyncResult[*ai.AIRecommendation]{
			Result: recommendation,
			Error:  err,
		}
	}()

	return result
}

// NaturalLanguageSearchAsync performs natural language search asynchronously
func (s *AIService) NaturalLanguageSearchAsync(ctx context.Context, query string, servers []map[string]string) <-chan AsyncResult[[]map[string]interface{}] {
	result := make(chan AsyncResult[[]map[string]interface{}], 1)

	go func() {
		defer close(result)
		response, err := s.NaturalLanguageSearch(ctx, query, servers)
		result <- AsyncResult[[]map[string]interface{}]{
			Result: response,
			Error:  err,
		}
	}()

	return result
}

// AnalyzeSecurityAsync analyzes security configurations asynchronously
func (s *AIService) AnalyzeSecurityAsync(ctx context.Context, config map[string]string) <-chan AsyncResult[*ai.AIRecommendation] {
	result := make(chan AsyncResult[*ai.AIRecommendation], 1)

	go func() {
		defer close(result)
		response, err := s.AnalyzeSecurity(ctx, config)
		result <- AsyncResult[*ai.AIRecommendation]{
			Result: response,
			Error:  err,
		}
	}()

	return result
}

// TestConnectionAsync tests AI connection asynchronously
func (s *AIService) TestConnectionAsync(ctx context.Context) <-chan AsyncResult[bool] {
	result := make(chan AsyncResult[bool], 1)

	go func() {
		defer close(result)
		err := s.TestConnection(ctx)
		result <- AsyncResult[bool]{
			Result: err == nil,
			Error:  err,
		}
	}()

	return result
}

// BatchProcessAsync processes multiple AI requests concurrently
func (s *AIService) BatchProcessAsync(ctx context.Context, requests []ai.AIRequest) <-chan []AsyncResult[*ai.AIResponse] {
	results := make(chan []AsyncResult[*ai.AIResponse], 1)

	go func() {
		defer close(results)

		var wg sync.WaitGroup
		responses := make([]AsyncResult[*ai.AIResponse], len(requests))

		for i, req := range requests {
			wg.Add(1)
			go func(index int, request ai.AIRequest) {
				defer wg.Done()
				response, err := s.makeAIRequest(ctx, &request)
				responses[index] = AsyncResult[*ai.AIResponse]{
					Result: response,
					Error:  err,
				}
			}(i, req)
		}

		wg.Wait()
		results <- responses
	}()

	return results
}

// GetConnectionPoolStats returns connection pool statistics
func (s *AIService) GetConnectionPoolStats() map[string]interface{} {
	return s.pool.Stats()
}

// CloseConnectionPool closes the connection pool
func (s *AIService) CloseConnectionPool() {
	s.pool.Close()
}

// Stop stops the AI service and cleans up resources
func (s *AIService) Stop() {
	s.cache.Stop()
	s.pool.Close()
}
