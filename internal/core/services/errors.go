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
	"fmt"
)

// ErrorContext provides context information for error wrapping
type ErrorContext struct {
	Operation string
	TraceID   interface{} // Can be string or tracing.TraceID
	Fields    map[string]interface{}
}

// WrapError wraps an error with contextual information
func WrapError(err error, ctx ErrorContext) error {
	if err == nil {
		return nil
	}

	// Build context string
	contextStr := ctx.Operation
	if ctx.TraceID != nil {
		contextStr += fmt.Sprintf(" (trace_id: %v)", ctx.TraceID)
	}

	// Add fields to context
	if len(ctx.Fields) > 0 {
		for k, v := range ctx.Fields {
			contextStr += fmt.Sprintf(", %s: %v", k, v)
		}
	}

	return fmt.Errorf("failed to %s: %w", contextStr, err)
}

// WrapErrorf wraps an error with a formatted message and context
func WrapErrorf(err error, ctx ErrorContext, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	// Build context string
	contextStr := ctx.Operation
	if ctx.TraceID != nil {
		contextStr += fmt.Sprintf(" (trace_id: %v)", ctx.TraceID)
	}

	// Add fields to context
	if len(ctx.Fields) > 0 {
		for k, v := range ctx.Fields {
			contextStr += fmt.Sprintf(", %s: %v", k, v)
		}
	}

	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("failed to %s, %s: %w", contextStr, msg, err)
}

// NewError creates a new error with context
func NewError(ctx ErrorContext, format string, args ...interface{}) error {
	// Build context string
	contextStr := ctx.Operation
	if ctx.TraceID != nil {
		contextStr += fmt.Sprintf(" (trace_id: %v)", ctx.TraceID)
	}

	// Add fields to context
	if len(ctx.Fields) > 0 {
		for k, v := range ctx.Fields {
			contextStr += fmt.Sprintf(", %s: %v", k, v)
		}
	}

	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s: %s", contextStr, msg)
}

// WithTraceID adds a trace ID to an error context
func (ctx ErrorContext) WithTraceID(traceID interface{}) ErrorContext {
	ctx.TraceID = traceID
	return ctx
}

// WithField adds a field to the error context
func (ctx ErrorContext) WithField(key string, value interface{}) ErrorContext {
	if ctx.Fields == nil {
		ctx.Fields = make(map[string]interface{})
	}
	ctx.Fields[key] = value
	return ctx
}

// WithFields adds multiple fields to the error context
func (ctx ErrorContext) WithFields(fields map[string]interface{}) ErrorContext {
	if ctx.Fields == nil {
		ctx.Fields = make(map[string]interface{})
	}
	for k, v := range fields {
		ctx.Fields[k] = v
	}
	return ctx
}

// NewErrorContext creates a new error context
func NewErrorContext(operation string) ErrorContext {
	return ErrorContext{
		Operation: operation,
		Fields:    make(map[string]interface{}),
	}
}

// WrapValidationError wraps a validation error with context
func WrapValidationError(err error, ctx ErrorContext, field string, value interface{}) error {
	if err == nil {
		return nil
	}
	ctx = ctx.WithField("field", field).WithField("value", value)
	return WrapErrorf(err, ctx, "validation failed")
}

// WrapSecurityError wraps a security-related error with context
func WrapSecurityError(err error, ctx ErrorContext, reason string) error {
	if err == nil {
		return nil
	}
	ctx = ctx.WithField("security_reason", reason)
	return WrapErrorf(err, ctx, "security check failed")
}

// WrapOperationError wraps an operation error with context
func WrapOperationError(err error, ctx ErrorContext, operation string) error {
	if err == nil {
		return nil
	}
	ctx.Operation = operation
	return WrapError(err, ctx)
}

// NewValidationError creates a new validation error
func NewValidationError(ctx ErrorContext, field string, reason string) error {
	ctx = ctx.WithField("field", field)
	return NewError(ctx, "validation failed: %s", reason)
}

// NewSecurityError creates a new security error
func NewSecurityError(ctx ErrorContext, reason string) error {
	ctx = ctx.WithField("security_reason", reason)
	return NewError(ctx, "security violation: %s", reason)
}
