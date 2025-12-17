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
	"fmt"
	"time"
)

// RetryConfig configures retry behavior for AI requests
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	BackoffFactor  float64
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 500 * time.Millisecond,
		MaxBackoff:     10 * time.Second,
		BackoffFactor:  2.0,
	}
}

// RetryableError represents an error that can be retried
type RetryableError struct {
	Err     error
	Attempt int
}

func (e *RetryableError) Error() string {
	return fmt.Sprintf("attempt %d failed: %v", e.Attempt, e.Err)
}

// IsRetryable determines if an error should trigger a retry
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Network errors, timeouts, and temporary errors are retryable
	errStr := err.Error()
	retryablePatterns := []string{
		"timeout",
		"connection refused",
		"connection reset",
		"temporary failure",
		"service unavailable",
		"502",
		"503",
		"504",
	}

	for _, pattern := range retryablePatterns {
		if contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || substr == "" ||
		len(s) >= len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// RetryWithBackoff executes a function with exponential backoff retry logic.
//
// This implements a robust retry strategy for transient failures:
//   - Exponential backoff: wait time doubles after each failed attempt
//   - Configurable retries: MaxRetries controls total attempts (initial + retries)
//   - Context-aware: respects context cancellation for graceful shutdown
//   - Retryable error detection: only retries on transient errors (network, timeouts)
//
// Algorithm:
//  1. Execute function
//  2. If success, return immediately
//  3. If non-retryable error, return immediately
//  4. If retryable error and attempts remaining:
//     - Calculate backoff: min(InitialBackoff * (2^attempt), MaxBackoff)
//     - Wait for backoff duration (respecting context)
//     - Retry
//  5. If all attempts exhausted, return RetryableError
//
// Example with defaults (InitialBackoff=500ms, MaxBackoff=10s, MaxRetries=3):
//
//	Attempt 1: immediate
//	Attempt 2: wait 500ms
//	Attempt 3: wait 1s
//	Attempt 4: wait 2s
//	Total: ~3.5s maximum
func RetryWithBackoff(ctx context.Context, config *RetryConfig, fn func() error) error {
	if config == nil {
		config = DefaultRetryConfig()
	}

	var lastErr error
	backoff := config.InitialBackoff

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Check if context is canceled
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry canceled: %w", ctx.Err())
		default:
		}

		// Execute the function
		err := fn()
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Check if error is retryable
		if !IsRetryable(err) {
			return fmt.Errorf("non-retryable error: %w", err)
		}

		// Don't sleep after the last attempt
		if attempt >= config.MaxRetries {
			break
		}

		// Calculate backoff duration
		sleepDuration := backoff
		if sleepDuration > config.MaxBackoff {
			sleepDuration = config.MaxBackoff
		}

		// Wait with context awareness
		select {
		case <-time.After(sleepDuration):
			// Continue to next attempt
		case <-ctx.Done():
			return fmt.Errorf("retry canceled during backoff: %w", ctx.Err())
		}

		// Increase backoff for next attempt
		backoff = time.Duration(float64(backoff) * config.BackoffFactor)
	}

	return &RetryableError{
		Err:     lastErr,
		Attempt: config.MaxRetries + 1,
	}
}
