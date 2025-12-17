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
	"os"
	"sync"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain/ai"
	"github.com/aryasoni98/wooak/internal/core/services/monitoring"
	"go.uber.org/zap"
)

// AIService provides AI-powered functionality for Wooak
type AIService struct {
	config      *ai.AIConfig
	httpClient  *http.Client
	cache       *AICache
	pool        *ConnectionPool
	monitoring  *monitoring.MonitoringService
	retryConfig *RetryConfig
	logger      *zap.SugaredLogger
	rateLimiter *RateLimiter
}

// NewAIService creates a new AI service
func NewAIService(config *ai.AIConfig) *AIService {
	return NewAIServiceWithLogger(config, nil)
}

// NewAIServiceWithLogger creates a new AI service with a logger instance
func NewAIServiceWithLogger(config *ai.AIConfig, logger *zap.SugaredLogger) *AIService {
	// Create connection pool configuration
	poolConfig := &PoolConfig{
		MaxConnections:        DefaultMaxConnections,
		MaxIdleConns:          DefaultMaxIdleConns,
		MaxConnsPerHost:       DefaultMaxConnsPerHost,
		IdleConnTimeout:       DefaultIdleConnTimeout,
		ResponseHeaderTimeout: DefaultResponseTimeout,
		RequestTimeout:        config.Timeout,
	}

	// Create rate limiter with default configuration
	rateLimiter := NewRateLimiter(RateLimiterConfig{
		MaxTokens:      DefaultRateLimitMaxTokens,
		RefillRate:     DefaultRateLimitRefillRate,
		BlockOnExhaust: DefaultRateLimitBlock,
	})

	return &AIService{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		cache:       NewAICache(config.CacheTTL),
		pool:        NewConnectionPool(poolConfig),
		monitoring:  nil, // Will be set via SetMonitoring
		retryConfig: DefaultRetryConfig(),
		logger:      logger,
		rateLimiter: rateLimiter,
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

// makeAIRequest makes a request to the AI provider with retry logic
func (s *AIService) makeAIRequest(ctx context.Context, request *ai.AIRequest) (*ai.AIResponse, error) {
	startTime := time.Now()

	// Check rate limit before making request
	if !s.rateLimiter.Allow() {
		if s.monitoring != nil {
			s.monitoring.GetMetrics().IncrementCounter("ai_rate_limit_exceeded_total", map[string]string{
				"provider": string(s.config.Provider),
				"type":     string(request.Type),
			})
		}
		if s.logger != nil {
			s.logger.Warnw("AI request rate limited", "provider", s.config.Provider, "type", request.Type)
		}
		return nil, fmt.Errorf("rate limit exceeded for AI provider %s", s.config.Provider)
	}

	var response *ai.AIResponse
	var err error

	// Wrap the AI request in retry logic
	retryErr := RetryWithBackoff(ctx, s.retryConfig, func() error {
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
			return fmt.Errorf("unsupported AI provider: %s", s.config.Provider)
		}
		return err
	})

	if retryErr != nil {
		if s.monitoring != nil {
			s.monitoring.GetMetrics().IncrementCounter("ai_request_error_total", map[string]string{
				"provider": string(s.config.Provider),
				"type":     string(request.Type),
			})
		}
		return nil, fmt.Errorf("AI request failed after retries: %w", retryErr)
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
		return nil, fmt.Errorf("failed to marshal Ollama request (request_id: %s, model: %s): %w", request.ID, request.Model, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request (request_id: %s, url: %s): %w", request.ID, url, err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Use connection pool for better performance
	client := s.pool.GetClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make Ollama HTTP request (request_id: %s, url: %s): %w", request.ID, url, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error but don't fail the request
			if s.logger != nil {
				s.logger.Warnw("Failed to close response body", "error", closeErr, "request_id", request.ID)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "Warning: Failed to close response body: %v\n", closeErr)
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama API request failed (request_id: %s, status_code: %d, url: %s)", request.ID, resp.StatusCode, url)
	}

	var ollamaResponse struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ollamaResponse); err != nil {
		return nil, fmt.Errorf("failed to decode Ollama response (request_id: %s, status_code: %d): %w", request.ID, resp.StatusCode, err)
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

// makeOpenAIRequest makes a request to OpenAI Chat Completions API
func (s *AIService) makeOpenAIRequest(ctx context.Context, request *ai.AIRequest) (*ai.AIResponse, error) {
	// Validate API key is configured
	if s.config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key not configured (request_id: %s). Please set APIKey in AIConfig", request.ID)
	}

	// Determine base URL - use config if set, otherwise default
	baseURL := s.config.BaseURL
	if baseURL == "" {
		baseURL = OpenAIBaseURL
	}
	url := fmt.Sprintf("%s%s", baseURL, OpenAIEndpoint)

	// Determine model - use request model if set, otherwise config model, otherwise default
	model := request.Model
	if model == "" {
		model = s.config.Model
	}
	if model == "" {
		model = OpenAIDefaultModel
	}

	// Build OpenAI Chat Completions request payload
	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": request.Prompt,
			},
		},
		"max_tokens":  request.MaxTokens,
		"temperature": request.Temperature,
		"stream":      false,
	}

	// Use max_tokens from config if request doesn't specify
	if request.MaxTokens == 0 {
		payload["max_tokens"] = s.config.MaxTokens
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OpenAI request (request_id: %s, model: %s): %w", request.ID, model, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request (request_id: %s, url: %s): %w", request.ID, url, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.APIKey))

	// Use connection pool for better performance
	client := s.pool.GetClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make OpenAI HTTP request (request_id: %s, url: %s): %w", request.ID, url, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error but don't fail the request
			if s.logger != nil {
				s.logger.Warnw("Failed to close response body", "error", closeErr, "request_id", request.ID)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "Warning: Failed to close response body: %v\n", closeErr)
			}
		}
	}()

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		// Try to read error response body for better error messages
		var errorResponse map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if s.logger != nil {
			s.logger.Errorw("OpenAI API request failed", "request_id", request.ID, "status_code", resp.StatusCode, "url", url, "response_body", errorResponse)
		}
		return nil, fmt.Errorf("OpenAI API request failed (request_id: %s, status_code: %d, url: %s)", request.ID, resp.StatusCode, url)
	}

	// Parse OpenAI response
	var openAIResponse struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		Model   string `json:"model"`
		Choices []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&openAIResponse); err != nil {
		return nil, fmt.Errorf("failed to decode OpenAI response (request_id: %s, status_code: %d): %w", request.ID, resp.StatusCode, err)
	}

	// Extract content from first choice
	if len(openAIResponse.Choices) == 0 {
		return nil, fmt.Errorf("OpenAI API returned no choices (request_id: %s)", request.ID)
	}

	content := openAIResponse.Choices[0].Message.Content
	tokensUsed := openAIResponse.Usage.TotalTokens

	return &ai.AIResponse{
		ID:         generateResponseID(),
		RequestID:  request.ID,
		Content:    content,
		Type:       request.Type,
		Model:      openAIResponse.Model,
		TokensUsed: tokensUsed,
		Timestamp:  time.Now(),
	}, nil
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

// Stop stops the AI service and cleans up resources with graceful shutdown
func (s *AIService) Stop() {
	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()

	// Create a channel to signal when cleanup is done
	done := make(chan struct{})

	go func() {
		// Stop cache cleanup goroutines
		s.cache.Stop()

		// Close connection pool
		s.pool.Close()

		close(done)
	}()

	// Wait for cleanup to complete or timeout
	select {
	case <-done:
		// Cleanup completed successfully
	case <-ctx.Done():
		// Timeout occurred, log but don't fail
		if s.logger != nil {
			s.logger.Warnw("AI service shutdown timed out", "timeout", DefaultShutdownTimeout)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Warning: AI service shutdown timed out after %v\n", DefaultShutdownTimeout)
		}
	}
}
