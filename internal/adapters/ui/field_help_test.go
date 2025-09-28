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

package ui

import (
	"testing"
)

func TestGetFieldHelp(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		wantNil   bool
		checkFunc func(*FieldHelp) bool
	}{
		{
			name:      "Host field should have help",
			fieldName: "Host",
			wantNil:   false,
			checkFunc: func(h *FieldHelp) bool {
				return h.Field == "Host" &&
					h.Description != "" &&
					h.Syntax != "" &&
					len(h.Examples) > 0
			},
		},
		{
			name:      "ProxyJump field should have help with version info",
			fieldName: "ProxyJump",
			wantNil:   false,
			checkFunc: func(h *FieldHelp) bool {
				return h.Field == "ProxyJump" &&
					h.Since != "" &&
					h.Category == "Connection"
			},
		},
		{
			name:      "LocalForward field should have correct category",
			fieldName: "LocalForward",
			wantNil:   false,
			checkFunc: func(h *FieldHelp) bool {
				return h.Category == "Forwarding"
			},
		},
		{
			name:      "Unknown field should return nil",
			fieldName: "NonExistentField",
			wantNil:   true,
			checkFunc: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			help := GetFieldHelp(tt.fieldName)

			if tt.wantNil {
				if help != nil {
					t.Errorf("GetFieldHelp(%s) = %v, want nil", tt.fieldName, help)
				}
				return
			}

			if help == nil {
				t.Errorf("GetFieldHelp(%s) = nil, want non-nil", tt.fieldName)
				return
			}

			if tt.checkFunc != nil && !tt.checkFunc(help) {
				t.Errorf("GetFieldHelp(%s) returned help that doesn't match expected criteria", tt.fieldName)
			}
		})
	}
}

func TestGetFieldsByCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		minCount int // Minimum expected fields in category
	}{
		{"Basic category", "Basic", 5},
		{"Connection category", "Connection", 5},
		{"Forwarding category", "Forwarding", 4},
		{"Authentication category", "Authentication", 5},
		{"Security category", "Security", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := GetFieldsByCategory(tt.category)
			if len(fields) < tt.minCount {
				t.Errorf("GetFieldsByCategory(%s) returned %d fields, want at least %d",
					tt.category, len(fields), tt.minCount)
			}
		})
	}
}

func TestGetAllCategories(t *testing.T) {
	categories := GetAllCategories()

	// Check we have at least the main categories
	expectedCategories := []string{
		"Basic", "Connection", "Forwarding", "Authentication", "Security",
	}

	for _, expected := range expectedCategories {
		found := false
		for _, cat := range categories {
			if cat == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetAllCategories() missing expected category: %s", expected)
		}
	}

	if len(categories) < len(expectedCategories) {
		t.Errorf("GetAllCategories() returned %d categories, want at least %d",
			len(categories), len(expectedCategories))
	}
}

func TestHelpContent(t *testing.T) {
	// Test that critical fields have comprehensive help
	criticalFields := []string{
		"Host", "Port", "User", "ProxyJump", "LocalForward",
		"ControlMaster", "StrictHostKeyChecking",
	}

	for _, field := range criticalFields {
		help := GetFieldHelp(field)
		if help == nil {
			t.Errorf("Critical field %s has no help", field)
			continue
		}

		if help.Description == "" {
			t.Errorf("Field %s has no description", field)
		}

		if help.Syntax == "" && len(help.Examples) == 0 {
			t.Errorf("Field %s has neither syntax nor examples", field)
		}

		if help.Default == "" {
			t.Logf("Warning: Field %s has no default value specified", field)
		}
	}
}
