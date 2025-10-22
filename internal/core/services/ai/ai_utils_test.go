// Copyright 2025.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ai

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGenerateRequestID(t *testing.T) {
	// Test that ID is generated successfully
	id1 := generateRequestID()
	if id1 == "" {
		t.Error("Expected non-empty request ID")
	}

	// Test that IDs are unique
	id2 := generateRequestID()
	if id1 == id2 {
		t.Error("Expected unique request IDs")
	}

	// Test that ID has expected format
	if len(id1) < 10 {
		t.Error("Expected request ID to have reasonable length")
	}
}

func TestGenerateResponseID(t *testing.T) {
	// Test that ID is generated successfully
	id1 := generateResponseID()
	if id1 == "" {
		t.Error("Expected non-empty response ID")
	}

	// Test that IDs are unique
	id2 := generateResponseID()
	if id1 == id2 {
		t.Error("Expected unique response IDs")
	}

	// Test that ID has expected format
	if len(id1) < 10 {
		t.Error("Expected response ID to have reasonable length")
	}
}

func TestGenerateRecommendationID(t *testing.T) {
	// Test that ID is generated successfully
	id1 := generateRecommendationID()
	if id1 == "" {
		t.Error("Expected non-empty recommendation ID")
	}

	// Test that IDs are unique
	id2 := generateRecommendationID()
	if id1 == id2 {
		t.Error("Expected unique recommendation IDs")
	}

	// Test that ID has expected format
	if len(id1) < 10 {
		t.Error("Expected recommendation ID to have reasonable length")
	}
}

func TestGenerateHash(t *testing.T) {
	// Test with string input
	hash1 := generateHash("test string")
	if hash1 == "" {
		t.Error("Expected non-empty hash")
	}

	// Test that same input produces same hash
	hash2 := generateHash("test string")
	if hash1 != hash2 {
		t.Error("Expected same input to produce same hash")
	}

	// Test that different input produces different hash
	hash3 := generateHash("different string")
	if hash1 == hash3 {
		t.Error("Expected different input to produce different hash")
	}

	// Test with different data types
	hash4 := generateHash(123)
	hash5 := generateHash(456)
	if hash4 == hash5 {
		t.Error("Expected different numeric inputs to produce different hashes")
	}
}

func TestIDGenerationWithTimeSeparation(t *testing.T) {
	// Test that IDs generated with time separation are different
	id1 := generateRequestID()
	time.Sleep(1 * time.Millisecond) // Small delay to ensure different timestamps
	id2 := generateRequestID()

	if id1 == id2 {
		t.Error("Expected IDs with time separation to be different")
	}
}

func TestIDGenerationConcurrency(t *testing.T) {
	// Test concurrent ID generation
	const numGoroutines = 10
	const numIDs = 100

	ids := make(chan string, numGoroutines*numIDs)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numIDs; j++ {
				ids <- generateRequestID()
			}
		}()
	}

	// Collect all IDs
	allIDs := make(map[string]bool)
	for i := 0; i < numGoroutines*numIDs; i++ {
		id := <-ids
		if allIDs[id] {
			t.Errorf("Found duplicate ID: %s", id)
		}
		allIDs[id] = true
	}
}

func TestIsOllamaRunning(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse int
		expectedResult bool
	}{
		{
			name:           "ollama running - 200 OK",
			serverResponse: http.StatusOK,
			expectedResult: true,
		},
		{
			name:           "ollama not running - 404",
			serverResponse: http.StatusNotFound,
			expectedResult: false,
		},
		{
			name:           "ollama error - 500",
			serverResponse: http.StatusInternalServerError,
			expectedResult: false,
		},
		{
			name:           "ollama unauthorized - 401",
			serverResponse: http.StatusUnauthorized,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/api/tags" {
					w.WriteHeader(tt.serverResponse)
					w.Write([]byte(`{"models":[]}`))
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			// Test IsOllamaRunning with the test server
			result := IsOllamaRunning(server.URL)
			if result != tt.expectedResult {
				t.Errorf("IsOllamaRunning() = %v, expected %v", result, tt.expectedResult)
			}
		})
	}
}

func TestIsOllamaRunning_InvalidURL(t *testing.T) {
	// Test with invalid URL
	result := IsOllamaRunning("http://invalid-url-that-does-not-exist")
	if result {
		t.Error("Expected IsOllamaRunning to return false for invalid URL")
	}
}

func TestIsOllamaRunning_Timeout(t *testing.T) {
	// Create a server that takes too long to respond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // Longer than our 2-second timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Test that timeout works
	result := IsOllamaRunning(server.URL)
	if result {
		t.Error("Expected IsOllamaRunning to return false due to timeout")
	}
}

func TestIsOllamaRunning_WrongEndpoint(t *testing.T) {
	// Create a server that doesn't have /api/tags endpoint
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Test that wrong endpoint returns false
	result := IsOllamaRunning(server.URL)
	if result {
		t.Error("Expected IsOllamaRunning to return false for wrong endpoint")
	}
}
