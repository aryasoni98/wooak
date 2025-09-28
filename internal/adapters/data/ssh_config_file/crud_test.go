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

package ssh_config_file

import (
	"testing"
)

func TestConvertCLIForwardToConfigFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic local forward",
			input:    "8080:localhost:80",
			expected: "8080 localhost:80",
		},
		{
			name:     "local forward with bind address",
			input:    "127.0.0.1:8080:localhost:80",
			expected: "127.0.0.1:8080 localhost:80",
		},
		{
			name:     "local forward with wildcard bind",
			input:    "*:8080:localhost:80",
			expected: "*:8080 localhost:80",
		},
		{
			name:     "remote forward",
			input:    "8080:localhost:3000",
			expected: "8080 localhost:3000",
		},
		{
			name:     "remote forward with bind address",
			input:    "0.0.0.0:80:localhost:8080",
			expected: "0.0.0.0:80 localhost:8080",
		},
		{
			name:     "forward with IPv6 address",
			input:    "8080:[2001:db8::1]:80",
			expected: "8080 [2001:db8::1]:80",
		},
		{
			name:     "forward with domain",
			input:    "3306:db.example.com:3306",
			expected: "3306 db.example.com:3306",
		},
		{
			name:     "invalid format - only one colon",
			input:    "8080:localhost",
			expected: "8080:localhost", // returned as-is
		},
		{
			name:     "invalid format - no colons",
			input:    "8080",
			expected: "8080", // returned as-is
		},
	}

	r := &Repository{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.convertCLIForwardToConfigFormat(tt.input)
			if result != tt.expected {
				t.Errorf("convertCLIForwardToConfigFormat(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertConfigForwardToCLIFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic local forward",
			input:    "8080 localhost:80",
			expected: "8080:localhost:80",
		},
		{
			name:     "local forward with bind address",
			input:    "127.0.0.1:8080 localhost:80",
			expected: "127.0.0.1:8080:localhost:80",
		},
		{
			name:     "local forward with wildcard bind",
			input:    "*:8080 localhost:80",
			expected: "*:8080:localhost:80",
		},
		{
			name:     "remote forward",
			input:    "8080 localhost:3000",
			expected: "8080:localhost:3000",
		},
		{
			name:     "remote forward with bind address",
			input:    "0.0.0.0:80 localhost:8080",
			expected: "0.0.0.0:80:localhost:8080",
		},
		{
			name:     "forward with IPv6 address",
			input:    "8080 [2001:db8::1]:80",
			expected: "8080:[2001:db8::1]:80",
		},
		{
			name:     "forward with domain",
			input:    "3306 db.example.com:3306",
			expected: "3306:db.example.com:3306",
		},
		{
			name:     "already in CLI format",
			input:    "8080:localhost:80",
			expected: "8080:localhost:80", // returned as-is
		},
		{
			name:     "no space separator",
			input:    "8080",
			expected: "8080", // returned as-is
		},
	}

	r := &Repository{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.convertConfigForwardToCLIFormat(tt.input)
			if result != tt.expected {
				t.Errorf("convertConfigForwardToCLIFormat(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
