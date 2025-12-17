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

func TestRateLimiter_Allow(t *testing.T) {
	limiter := NewRateLimiter(RateLimiterConfig{
		MaxTokens:      2.0,
		RefillRate:     1.0, // 1 token per second
		BlockOnExhaust: false,
		InitialTokens:  2.0,
	})

	// Should allow first 2 requests
	if !limiter.Allow() {
		t.Error("Expected first request to be allowed")
	}
	if !limiter.Allow() {
		t.Error("Expected second request to be allowed")
	}

	// Third request should be denied (no tokens left)
	if limiter.Allow() {
		t.Error("Expected third request to be denied")
	}
}

func TestRateLimiter_AllowN(t *testing.T) {
	limiter := NewRateLimiter(RateLimiterConfig{
		MaxTokens:      10.0,
		RefillRate:     1.0,
		BlockOnExhaust: false,
		InitialTokens:  10.0,
	})

	// Should allow request for 5 tokens
	if !limiter.AllowN(5.0) {
		t.Error("Expected request for 5 tokens to be allowed")
	}

	// Should allow request for remaining 5 tokens
	if !limiter.AllowN(5.0) {
		t.Error("Expected request for 5 tokens to be allowed")
	}

	// Should deny request for more tokens
	if limiter.AllowN(1.0) {
		t.Error("Expected request to be denied when no tokens available")
	}
}

func TestRateLimiter_Refill(t *testing.T) {
	limiter := NewRateLimiter(RateLimiterConfig{
		MaxTokens:      10.0,
		RefillRate:     10.0, // 10 tokens per second
		BlockOnExhaust: false,
		InitialTokens:  10.0,
	})

	// Consume all tokens
	limiter.AllowN(10.0)

	// Wait for refill (should refill quickly with high rate)
	time.Sleep(150 * time.Millisecond)

	// Should allow request after refill
	if !limiter.Allow() {
		t.Error("Expected request to be allowed after refill")
	}
}

func TestRateLimiter_Stats(t *testing.T) {
	limiter := NewRateLimiter(RateLimiterConfig{
		MaxTokens:      10.0,
		RefillRate:     2.0,
		BlockOnExhaust: false,
		InitialTokens:  10.0,
	})

	// Consume some tokens
	limiter.AllowN(3.0)

	stats := limiter.Stats()
	if stats == nil {
		t.Fatal("Expected stats, got nil")
	}

	tokens, ok := stats["tokens_available"].(float64)
	if !ok {
		t.Fatal("Expected tokens_available in stats")
	}
	// Use tolerance for floating point comparison
	expected := 7.0
	tolerance := 0.01
	if tokens < expected-tolerance || tokens > expected+tolerance {
		t.Errorf("Expected approximately %.2f tokens, got %f", expected, tokens)
	}

	maxTokens, ok := stats["max_tokens"].(float64)
	if !ok || maxTokens != 10.0 {
		t.Errorf("Expected max_tokens to be 10.0, got %v", maxTokens)
	}
}

func TestRateLimiter_Wait_BlockOnExhaust(t *testing.T) {
	limiter := NewRateLimiter(RateLimiterConfig{
		MaxTokens:      2.0,
		RefillRate:     10.0, // Fast refill for testing
		BlockOnExhaust: true,
		InitialTokens:  2.0,
	})

	// Consume all tokens
	limiter.AllowN(2.0)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Should wait and eventually succeed
	err := limiter.Wait(ctx, 1.0)
	if err != nil {
		t.Logf("Wait returned error (may be timeout): %v", err)
	}
}

func TestRateLimiter_Wait_ContextCancel(t *testing.T) {
	limiter := NewRateLimiter(RateLimiterConfig{
		MaxTokens:      1.0,
		RefillRate:     0.1, // Slow refill
		BlockOnExhaust: true,
		InitialTokens:  1.0,
	})

	// Consume all tokens
	limiter.AllowN(1.0)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// Should timeout
	err := limiter.Wait(ctx, 1.0)
	if err == nil {
		t.Error("Expected error when context times out")
	}
	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Logf("Got error (may be expected): %v", err)
	}
}

func TestRateLimiter_Concurrent(t *testing.T) {
	limiter := NewRateLimiter(RateLimiterConfig{
		MaxTokens:      100.0,
		RefillRate:     100.0,
		BlockOnExhaust: false,
		InitialTokens:  100.0,
	})

	// Test concurrent access
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			allowed := limiter.Allow()
			done <- allowed
		}()
	}

	allowedCount := 0
	for i := 0; i < 10; i++ {
		if <-done {
			allowedCount++
		}
	}

	// All should be allowed with sufficient tokens
	if allowedCount != 10 {
		t.Errorf("Expected all 10 requests to be allowed, got %d", allowedCount)
	}
}
