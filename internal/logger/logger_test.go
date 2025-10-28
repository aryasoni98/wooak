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

package logger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		appName string
		wantErr bool
	}{
		{
			name:    "valid app name",
			appName: "TestApp",
			wantErr: false,
		},
		{
			name:    "empty app name",
			appName: "",
			wantErr: false, // Should still work
		},
		{
			name:    "app name with special characters",
			appName: "Test-App_123",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.appName)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && logger == nil {
				t.Error("New() returned nil logger without error")
			}
			if logger != nil {
				// Verify we can log without panic
				logger.Info("test message")
				logger.Infow("test message with fields", "key", "value")
			}
		})
	}
}

func TestLoggerLevels(t *testing.T) {
	logger, err := New("TestApp")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test that all log levels work without panic
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	logger.Debugw("debug with fields", "key", "value")
	logger.Infow("info with fields", "key", "value")
	logger.Warnw("warn with fields", "key", "value")
	logger.Errorw("error with fields", "key", "value")
}

func TestLoggerSync(t *testing.T) {
	logger, err := New("TestApp")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test sync doesn't panic
	err = logger.Sync()
	// Sync may return error on some platforms (like stdout/stderr), which is acceptable
	if err != nil {
		t.Logf("Sync returned error (expected on some platforms): %v", err)
	}
}

func TestLoggerConcurrency(t *testing.T) {
	logger, err := New("TestApp")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test concurrent logging doesn't cause race conditions
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.Infow("concurrent log", "id", id)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestLoggerWithFields(t *testing.T) {
	logger, err := New("TestApp")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test logging with various field types
	logger.Infow("message with fields",
		"string", "value",
		"int", 42,
		"bool", true,
		"float", 3.14,
	)

	// Test with nested logger
	childLogger := logger.With("component", "test")
	childLogger.Info("child logger message")
}

func TestLoggerErrorHandling(t *testing.T) {
	logger, err := New("TestApp")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test error logging doesn't panic
	testErr := zap.Error(nil)
	if testErr.Type == zapcore.ErrorType {
		// This is expected - zap wraps nil errors
		t.Log("zap.Error correctly handles nil")
	}

	// Test with actual error
	logger.Errorw("error occurred", "error", err)
}

func TestLoggerNamed(t *testing.T) {
	logger, err := New("TestApp")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test named logger
	namedLogger := logger.Named("subsystem")
	namedLogger.Info("named logger message")

	// Verify we can create multiple named loggers
	logger1 := logger.Named("component1")
	logger2 := logger.Named("component2")

	logger1.Info("component1 message")
	logger2.Info("component2 message")
}
