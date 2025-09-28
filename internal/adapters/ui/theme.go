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
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ModernTheme defines a clean, modern color palette for the Wooak TUI
type ModernTheme struct {
	// Background colors
	Background     tcell.Color // Main background
	Surface        tcell.Color // Panel backgrounds
	SurfaceVariant tcell.Color // Secondary surfaces
	Overlay        tcell.Color // Overlay backgrounds

	// Text colors
	OnBackground     tcell.Color // Primary text on background
	OnSurface        tcell.Color // Primary text on surface
	OnSurfaceVariant tcell.Color // Secondary text on surface
	Muted            tcell.Color // Muted/secondary text

	// Accent colors
	Primary          tcell.Color // Primary accent (blue)
	PrimaryVariant   tcell.Color // Primary variant
	Secondary        tcell.Color // Secondary accent (green)
	SecondaryVariant tcell.Color // Secondary variant

	// Status colors
	Success tcell.Color // Success states
	Warning tcell.Color // Warning states
	Error   tcell.Color // Error states
	Info    tcell.Color // Info states

	// Interactive colors
	Selected     tcell.Color // Selected item background
	SelectedText tcell.Color // Selected item text
	Hover        tcell.Color // Hover state
	Focus        tcell.Color // Focus state

	// Border colors
	Border      tcell.Color // Default borders
	BorderFocus tcell.Color // Focused borders
	BorderMuted tcell.Color // Muted borders

	// Special colors
	Header    tcell.Color // Header background
	Footer    tcell.Color // Footer background
	Separator tcell.Color // Separator lines
}

// GetModernTheme returns a modern, clean theme for Wooak
func GetModernTheme() *ModernTheme {
	return &ModernTheme{
		// Background colors - Dark theme with subtle variations
		Background:     tcell.Color16,  // Deep black
		Surface:        tcell.Color235, // Dark gray
		SurfaceVariant: tcell.Color236, // Slightly lighter gray
		Overlay:        tcell.Color237, // Overlay gray

		// Text colors - High contrast for readability
		OnBackground:     tcell.Color255, // Pure white
		OnSurface:        tcell.Color252, // Off-white
		OnSurfaceVariant: tcell.Color248, // Light gray
		Muted:            tcell.Color244, // Medium gray

		// Accent colors - Modern blue and green palette
		Primary:          tcell.Color39,  // Bright blue
		PrimaryVariant:   tcell.Color75,  // Lighter blue
		Secondary:        tcell.Color82,  // Bright green
		SecondaryVariant: tcell.Color118, // Lighter green

		// Status colors - Clear status indicators
		Success: tcell.Color46,  // Bright green
		Warning: tcell.Color214, // Orange
		Error:   tcell.Color196, // Red
		Info:    tcell.Color51,  // Cyan

		// Interactive colors - Clear selection states
		Selected:     tcell.Color24,  // Dark blue
		SelectedText: tcell.Color255, // White
		Hover:        tcell.Color238, // Light gray
		Focus:        tcell.Color39,  // Blue

		// Border colors - Subtle but visible
		Border:      tcell.Color238, // Medium gray
		BorderFocus: tcell.Color39,  // Blue
		BorderMuted: tcell.Color236, // Dark gray

		// Special colors
		Header:    tcell.Color234, // Dark header
		Footer:    tcell.Color235, // Dark footer
		Separator: tcell.Color238, // Separator gray
	}
}

// ApplyTheme applies the modern theme to tview styles
func (t *ModernTheme) ApplyTheme() {
	// Set global tview styles
	tview.Styles.PrimitiveBackgroundColor = t.Background
	tview.Styles.ContrastBackgroundColor = t.Surface
	tview.Styles.BorderColor = t.Border
	tview.Styles.TitleColor = t.OnSurface
	tview.Styles.PrimaryTextColor = t.OnSurface
	tview.Styles.TertiaryTextColor = t.Muted
	tview.Styles.SecondaryTextColor = t.OnSurfaceVariant
	tview.Styles.GraphicsColor = t.Border
}

// GetColorName returns a human-readable name for a color (for debugging)
func (t *ModernTheme) GetColorName(color tcell.Color) string {
	//nolint:exhaustive // Only check theme colors, not all possible tcell colors
	switch color {
	case t.Background:
		return "Background"
	case t.Surface:
		return "Surface"
	case t.SurfaceVariant:
		return "SurfaceVariant"
	case t.Overlay:
		return "Overlay"
	case t.OnBackground:
		return "OnBackground"
	case t.OnSurface:
		return "OnSurface"
	case t.OnSurfaceVariant:
		return "OnSurfaceVariant"
	case t.Muted:
		return "Muted"
	case t.Primary:
		return "Primary"
	case t.PrimaryVariant:
		return "PrimaryVariant"
	case t.Secondary:
		return "Secondary"
	case t.SecondaryVariant:
		return "SecondaryVariant"
	case t.Success:
		return "Success"
	case t.Warning:
		return "Warning"
	case t.Error:
		return "Error"
	case t.Info:
		return "Info"
	case t.Selected:
		return "Selected"
	case t.SelectedText:
		return "SelectedText"
	case t.Hover:
		return "Hover"
	case t.Focus:
		return "Focus"
	case t.Border:
		return "Border"
	case t.BorderFocus:
		return "BorderFocus"
	case t.BorderMuted:
		return "BorderMuted"
	case t.Header:
		return "Header"
	case t.Footer:
		return "Footer"
	case t.Separator:
		return "Separator"
	default:
		return "Unknown"
	}
}
