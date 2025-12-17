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

package services

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorContext_WithTraceID(t *testing.T) {
	ctx := NewErrorContext("test operation")
	ctx = ctx.WithTraceID("trace-123")

	if ctx.TraceID != "trace-123" {
		t.Errorf("Expected trace ID 'trace-123', got %v", ctx.TraceID)
	}
}

func TestErrorContext_WithField(t *testing.T) {
	ctx := NewErrorContext("test operation")
	ctx = ctx.WithField("key1", "value1")
	ctx = ctx.WithField("key2", 42)

	if ctx.Fields["key1"] != "value1" {
		t.Errorf("Expected field 'key1' to be 'value1', got %v", ctx.Fields["key1"])
	}
	if ctx.Fields["key2"] != 42 {
		t.Errorf("Expected field 'key2' to be 42, got %v", ctx.Fields["key2"])
	}
}

func TestErrorContext_WithFields(t *testing.T) {
	ctx := NewErrorContext("test operation")
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}
	ctx = ctx.WithFields(fields)

	if ctx.Fields["key1"] != "value1" {
		t.Errorf("Expected field 'key1' to be 'value1', got %v", ctx.Fields["key1"])
	}
	if ctx.Fields["key2"] != 42 {
		t.Errorf("Expected field 'key2' to be 42, got %v", ctx.Fields["key2"])
	}
	if ctx.Fields["key3"] != true {
		t.Errorf("Expected field 'key3' to be true, got %v", ctx.Fields["key3"])
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")
	ctx := NewErrorContext("test operation").
		WithTraceID("trace-123").
		WithField("key", "value")

	wrappedErr := WrapError(originalErr, ctx)

	if wrappedErr == nil {
		t.Fatal("Expected wrapped error, got nil")
	}

	errStr := wrappedErr.Error()
	if !strings.Contains(errStr, "test operation") {
		t.Errorf("Expected error to contain 'test operation', got: %s", errStr)
	}
	if !strings.Contains(errStr, "trace-123") {
		t.Errorf("Expected error to contain 'trace-123', got: %s", errStr)
	}
	if !strings.Contains(errStr, "original error") {
		t.Errorf("Expected error to contain 'original error', got: %s", errStr)
	}

	// Verify error wrapping
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should wrap original error")
	}
}

func TestWrapError_NilError(t *testing.T) {
	ctx := NewErrorContext("test operation")
	wrappedErr := WrapError(nil, ctx)

	if wrappedErr != nil {
		t.Errorf("Expected nil error, got: %v", wrappedErr)
	}
}

func TestWrapErrorf(t *testing.T) {
	originalErr := errors.New("original error")
	ctx := NewErrorContext("test operation").
		WithTraceID("trace-123")

	wrappedErr := WrapErrorf(originalErr, ctx, "additional context: %s", "extra info")

	if wrappedErr == nil {
		t.Fatal("Expected wrapped error, got nil")
	}

	errStr := wrappedErr.Error()
	if !strings.Contains(errStr, "test operation") {
		t.Errorf("Expected error to contain 'test operation', got: %s", errStr)
	}
	if !strings.Contains(errStr, "additional context: extra info") {
		t.Errorf("Expected error to contain formatted message, got: %s", errStr)
	}
	if !strings.Contains(errStr, "original error") {
		t.Errorf("Expected error to contain 'original error', got: %s", errStr)
	}
}

func TestNewError(t *testing.T) {
	ctx := NewErrorContext("test operation").
		WithTraceID("trace-123").
		WithField("key", "value")

	err := NewError(ctx, "something went wrong: %s", "details")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "test operation") {
		t.Errorf("Expected error to contain 'test operation', got: %s", errStr)
	}
	if !strings.Contains(errStr, "something went wrong: details") {
		t.Errorf("Expected error to contain formatted message, got: %s", errStr)
	}
}

func TestWrapValidationError(t *testing.T) {
	originalErr := errors.New("validation failed")
	ctx := NewErrorContext("validate server").
		WithTraceID("trace-123")

	wrappedErr := WrapValidationError(originalErr, ctx, "Alias", "invalid-value")

	if wrappedErr == nil {
		t.Fatal("Expected wrapped error, got nil")
	}

	errStr := wrappedErr.Error()
	if !strings.Contains(errStr, "validation failed") {
		t.Errorf("Expected error to contain 'validation failed', got: %s", errStr)
	}
	if !strings.Contains(errStr, "Alias") {
		t.Errorf("Expected error to contain field name 'Alias', got: %s", errStr)
	}
}

func TestWrapSecurityError(t *testing.T) {
	originalErr := errors.New("security check failed")
	ctx := NewErrorContext("SSH connection").
		WithTraceID("trace-123")

	wrappedErr := WrapSecurityError(originalErr, ctx, "invalid alias format")

	if wrappedErr == nil {
		t.Fatal("Expected wrapped error, got nil")
	}

	errStr := wrappedErr.Error()
	if !strings.Contains(errStr, "security check failed") {
		t.Errorf("Expected error to contain 'security check failed', got: %s", errStr)
	}
	if !strings.Contains(errStr, "invalid alias format") {
		t.Errorf("Expected error to contain security reason, got: %s", errStr)
	}
}

func TestNewValidationError(t *testing.T) {
	ctx := NewErrorContext("validate server").
		WithTraceID("trace-123")

	err := NewValidationError(ctx, "Alias", "cannot be empty")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "validation failed") {
		t.Errorf("Expected error to contain 'validation failed', got: %s", errStr)
	}
	if !strings.Contains(errStr, "cannot be empty") {
		t.Errorf("Expected error to contain reason, got: %s", errStr)
	}
}

func TestNewSecurityError(t *testing.T) {
	ctx := NewErrorContext("SSH connection").
		WithTraceID("trace-123").
		WithField("alias", "test-server")

	err := NewSecurityError(ctx, "alias not found in known servers list")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	errStr := err.Error()
	if !strings.Contains(errStr, "security violation") {
		t.Errorf("Expected error to contain 'security violation', got: %s", errStr)
	}
	if !strings.Contains(errStr, "alias not found in known servers list") {
		t.Errorf("Expected error to contain reason, got: %s", errStr)
	}
}

func TestErrorContext_Chaining(t *testing.T) {
	ctx := NewErrorContext("test operation").
		WithTraceID("trace-123").
		WithField("key1", "value1").
		WithField("key2", "value2").
		WithFields(map[string]interface{}{
			"key3": "value3",
			"key4": 42,
		})

	if ctx.Operation != "test operation" {
		t.Errorf("Expected operation 'test operation', got: %s", ctx.Operation)
	}
	if ctx.TraceID != "trace-123" {
		t.Errorf("Expected trace ID 'trace-123', got: %v", ctx.TraceID)
	}
	if len(ctx.Fields) != 4 {
		t.Errorf("Expected 4 fields, got: %d", len(ctx.Fields))
	}
}

func TestWrapOperationError(t *testing.T) {
	originalErr := errors.New("operation failed")
	ctx := NewErrorContext("original operation").
		WithTraceID("trace-123")

	wrappedErr := WrapOperationError(originalErr, ctx, "new operation")

	if wrappedErr == nil {
		t.Fatal("Expected wrapped error, got nil")
	}

	errStr := wrappedErr.Error()
	if !strings.Contains(errStr, "new operation") {
		t.Errorf("Expected error to contain 'new operation', got: %s", errStr)
	}
	if !strings.Contains(errStr, "operation failed") {
		t.Errorf("Expected error to contain original error, got: %s", errStr)
	}
}
