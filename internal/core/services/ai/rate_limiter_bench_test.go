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
)

// BenchmarkRateLimiterAllow benchmarks rate limiter Allow operations
func BenchmarkRateLimiterAllow(b *testing.B) {
	limiter := NewRateLimiter(RateLimiterConfig{
		MaxTokens:      10.0,
		RefillRate:     2.0,
		BlockOnExhaust: false,
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = limiter.Allow()
		}
	})
}

// BenchmarkRateLimiterAllowN benchmarks rate limiter AllowN operations
func BenchmarkRateLimiterAllowN(b *testing.B) {
	limiter := NewRateLimiter(RateLimiterConfig{
		MaxTokens:      10.0,
		RefillRate:     2.0,
		BlockOnExhaust: false,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = limiter.AllowN(1.0)
	}
}

// BenchmarkRateLimiterStats benchmarks rate limiter Stats operations
func BenchmarkRateLimiterStats(b *testing.B) {
	limiter := NewRateLimiter(RateLimiterConfig{
		MaxTokens:      10.0,
		RefillRate:     2.0,
		BlockOnExhaust: false,
	})

	// Prime the limiter
	for i := 0; i < 5; i++ {
		_ = limiter.Allow()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = limiter.Stats()
	}
}
