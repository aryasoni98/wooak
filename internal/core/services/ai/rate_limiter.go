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
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter.
// The token bucket algorithm allows bursts of traffic up to maxTokens while
// maintaining an average rate of refillRate tokens per second.
//
// Algorithm:
//   - Tokens are added continuously at refillRate per second
//   - Tokens cannot exceed maxTokens (bucket capacity)
//   - Requests consume tokens; if insufficient tokens, request is denied
//   - If blockOnExhaust is true, requests wait until tokens are available
//
// This implementation is thread-safe and uses a mutex to protect shared state.
type RateLimiter struct {
	mu             sync.Mutex
	tokens         float64   // Current number of available tokens
	maxTokens      float64   // Maximum bucket capacity
	refillRate     float64   // Tokens added per second
	lastRefill     time.Time // Last time tokens were refilled
	blockOnExhaust bool      // Whether to block when tokens exhausted
}

// RateLimiterConfig configures a rate limiter
type RateLimiterConfig struct {
	MaxTokens      float64 // Maximum number of tokens
	RefillRate     float64 // Tokens per second
	BlockOnExhaust bool    // Whether to block when tokens are exhausted
	InitialTokens  float64 // Initial token count (defaults to MaxTokens)
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	initialTokens := config.InitialTokens
	if initialTokens == 0 {
		initialTokens = config.MaxTokens
	}

	return &RateLimiter{
		tokens:         initialTokens,
		maxTokens:      config.MaxTokens,
		refillRate:     config.RefillRate,
		lastRefill:     time.Now(),
		blockOnExhaust: config.BlockOnExhaust,
	}
}

// Allow checks if a request is allowed and consumes a token if available
func (rl *RateLimiter) Allow() bool {
	return rl.AllowN(1)
}

// AllowN checks if N requests are allowed and consumes tokens if available
func (rl *RateLimiter) AllowN(n float64) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refill()

	if rl.tokens >= n {
		rl.tokens -= n
		return true
	}

	return false
}

// Wait waits until tokens are available and then consumes them
func (rl *RateLimiter) Wait(ctx context.Context, n float64) error {
	if !rl.blockOnExhaust {
		return fmt.Errorf("rate limiter is not configured to block")
	}

	for {
		rl.mu.Lock()
		rl.refill()

		if rl.tokens >= n {
			rl.tokens -= n
			rl.mu.Unlock()
			return nil
		}

		// Calculate how long to wait
		tokensNeeded := n - rl.tokens
		waitTime := time.Duration(float64(time.Second) * tokensNeeded / rl.refillRate)
		rl.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			// Continue loop to check again
		}
	}
}

// refill adds tokens based on elapsed time (must be called with lock held).
//
// The refill calculation:
//   - Calculates time elapsed since last refill
//   - Multiplies elapsed time by refillRate to get tokens to add
//   - Caps tokens at maxTokens to prevent bucket overflow
//   - Updates lastRefill timestamp for next calculation
//
// This ensures tokens accumulate continuously even when no requests are made,
// allowing bursts of traffic up to maxTokens.
func (rl *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	tokensToAdd := elapsed * rl.refillRate

	if tokensToAdd > 0 {
		rl.tokens = minFloat64(rl.tokens+tokensToAdd, rl.maxTokens)
		rl.lastRefill = now
	}
}

// Stats returns current rate limiter statistics
func (rl *RateLimiter) Stats() map[string]interface{} {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refill()

	return map[string]interface{}{
		"tokens_available": rl.tokens,
		"max_tokens":       rl.maxTokens,
		"refill_rate":      rl.refillRate,
		"block_on_exhaust": rl.blockOnExhaust,
	}
}

// minFloat64 returns the minimum of two float64 values
func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
