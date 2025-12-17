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

package tracing

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type traceIDKey struct{}

// TraceID represents a unique request trace identifier
type TraceID string

// GenerateTraceID generates a new unique trace ID
func GenerateTraceID() TraceID {
	// Generate a trace ID using timestamp and random number
	timestamp := time.Now().UnixNano()
	random := rand.Int63()
	return TraceID(fmt.Sprintf("%016x-%016x", timestamp, random))
}

// WithTraceID adds a trace ID to the context
func WithTraceID(ctx context.Context, traceID TraceID) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// GetTraceID retrieves the trace ID from the context
func GetTraceID(ctx context.Context) (TraceID, bool) {
	traceID, ok := ctx.Value(traceIDKey{}).(TraceID)
	return traceID, ok
}

// GetTraceIDOrNew retrieves the trace ID from context or generates a new one
func GetTraceIDOrNew(ctx context.Context) TraceID {
	if traceID, ok := GetTraceID(ctx); ok {
		return traceID
	}
	return GenerateTraceID()
}

// String returns the string representation of the trace ID
func (t TraceID) String() string {
	return string(t)
}
