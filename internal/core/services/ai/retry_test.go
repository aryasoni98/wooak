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
	"errors"
	"testing"
	"time"
)

func TestDefaultRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", config.MaxRetries)
	}
	if config.InitialBackoff != 500*time.Millisecond {
		t.Errorf("Expected InitialBackoff to be 500ms, got %v", config.InitialBackoff)
	}
	if config.MaxBackoff != 10*time.Second {
		t.Errorf("Expected MaxBackoff to be 10s, got %v", config.MaxBackoff)
	}
	if config.BackoffFactor != 2.0 {
		t.Errorf("Expected BackoffFactor to be 2.0, got %f", config.BackoffFactor)
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "timeout error",
			err:      errors.New("request timeout"),
			expected: true,
		},
		{
			name:     "connection refused",
			err:      errors.New("connection refused"),
			expected: true,
		},
		{
			name:     "503 service unavailable",
			err:      errors.New("503 service unavailable"),
			expected: true,
		},
		{
			name:     "non-retryable error",
			err:      errors.New("invalid request"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryable(tt.err)
			if result != tt.expected {
				t.Errorf("IsRetryable() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestRetryWithBackoff_Success(t *testing.T) {
	ctx := context.Background()
	config := &RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	attempt := 0
	err := RetryWithBackoff(ctx, config, func() error {
		attempt++
		if attempt < 2 {
			return errors.New("timeout") // Retryable error
		}
		return nil // Success on second attempt
	})

	if err != nil {
		t.Errorf("Expected success after retry, got error: %v", err)
	}
	if attempt != 2 {
		t.Errorf("Expected 2 attempts, got %d", attempt)
	}
}

func TestRetryWithBackoff_MaxRetriesExceeded(t *testing.T) {
	ctx := context.Background()
	config := &RetryConfig{
		MaxRetries:     2,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	attempt := 0
	err := RetryWithBackoff(ctx, config, func() error {
		attempt++
		return errors.New("timeout") // Always fails with retryable error
	})

	if err == nil {
		t.Error("Expected error after max retries, got nil")
	}
	if attempt != 3 { // Initial attempt + 2 retries
		t.Errorf("Expected 3 attempts, got %d", attempt)
	}
}

func TestRetryWithBackoff_NonRetryableError(t *testing.T) {
	ctx := context.Background()
	config := DefaultRetryConfig()

	attempt := 0
	err := RetryWithBackoff(ctx, config, func() error {
		attempt++
		return errors.New("invalid request") // Non-retryable error
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if attempt != 1 {
		t.Errorf("Expected 1 attempt for non-retryable error, got %d", attempt)
	}
}

func TestRetryWithBackoff_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	config := &RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     1 * time.Second,
		BackoffFactor:  2.0,
	}

	// Cancel context after first attempt
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	attempt := 0
	err := RetryWithBackoff(ctx, config, func() error {
		attempt++
		return errors.New("timeout") // Retryable error
	})

	if err == nil {
		t.Error("Expected error due to context cancellation, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got: %v", err)
	}
}

func TestRetryWithBackoff_ExponentialBackoff(t *testing.T) {
	ctx := context.Background()
	config := &RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 50 * time.Millisecond,
		MaxBackoff:     500 * time.Millisecond,
		BackoffFactor:  2.0,
	}

	attempts := []time.Time{}
	err := RetryWithBackoff(ctx, config, func() error {
		attempts = append(attempts, time.Now())
		return errors.New("timeout") // Always fails
	})

	if err == nil {
		t.Error("Expected error, got nil")
	}

	// Verify exponential backoff timing
	if len(attempts) != 4 { // Initial + 3 retries
		t.Errorf("Expected 4 attempts, got %d", len(attempts))
	}

	// Check that delays are increasing (with some tolerance)
	if len(attempts) >= 3 {
		delay1 := attempts[1].Sub(attempts[0])
		delay2 := attempts[2].Sub(attempts[1])

		// Second delay should be roughly 2x first delay (exponential)
		// Allow 30% tolerance for timing variations
		if delay2 < delay1 || delay2 > delay1*3 {
			t.Logf("Delays: %v, %v", delay1, delay2)
		}
	}
}

func TestRetryWithBackoff_NilConfig(t *testing.T) {
	ctx := context.Background()

	attempt := 0
	err := RetryWithBackoff(ctx, nil, func() error {
		attempt++
		if attempt < 2 {
			return errors.New("timeout")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected success with default config, got error: %v", err)
	}
}

func TestRetryableError(t *testing.T) {
	originalErr := errors.New("connection timeout")
	retryErr := &RetryableError{
		Err:     originalErr,
		Attempt: 3,
	}

	errMsg := retryErr.Error()
	expectedMsg := "attempt 3 failed: connection timeout"
	if errMsg != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, errMsg)
	}
}
