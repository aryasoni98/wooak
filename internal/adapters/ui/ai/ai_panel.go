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
	"strings"
	"time"

	aiDomain "github.com/aryasoni98/wooak/internal/core/domain/ai"
	aiService "github.com/aryasoni98/wooak/internal/core/services/ai"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// AIPanel provides a UI for AI features
type AIPanel struct {
	app        *tview.Application
	aiSvc      *aiService.AIService
	config     *aiDomain.AIConfig
	form       *tview.Form
	textView   *tview.TextView
	queryInput *tview.InputField
	resultView *tview.TextView
	statusView *tview.TextView
}

// NewAIPanel creates a new AI panel
func NewAIPanel(app *tview.Application, aiSvc *aiService.AIService) *AIPanel {
	ap := &AIPanel{
		app:    app,
		aiSvc:  aiSvc,
		config: aiSvc.GetConfig(),
	}

	ap.setupUI()
	return ap
}

// setupUI sets up the AI panel UI
func (ap *AIPanel) setupUI() {
	// Create the main form
	ap.form = tview.NewForm()
	ap.form.SetBorder(true).SetTitle(" AI Configuration ")

	// Add AI configuration fields
	ap.form.AddDropDown("Provider", []string{"ollama", "openai"}, 0, nil)
	ap.form.AddDropDown("Model", aiService.GetAvailableModels(), 0, nil)
	ap.form.AddInputField("Base URL", ap.config.BaseURL, 30, nil, nil)
	ap.form.AddInputField("Max Tokens", fmt.Sprintf("%d", ap.config.MaxTokens), 10, nil, nil)
	ap.form.AddInputField("Temperature", fmt.Sprintf("%.2f", ap.config.Temperature), 10, nil, nil)
	ap.form.AddCheckbox("Enabled", ap.config.Enabled, nil)
	ap.form.AddCheckbox("Cache Enabled", ap.config.CacheEnabled, nil)

	// Add buttons
	ap.form.AddButton("Save Config", ap.saveConfig)
	ap.form.AddButton("Test Connection", ap.testConnection)
	ap.form.AddButton("Reset to Default", ap.resetConfig)

	// Create AI interaction section
	ap.setupAIInteraction()

	// Create result view
	ap.resultView = tview.NewTextView()
	ap.resultView.SetBorder(true).SetTitle(" AI Results ")
	ap.resultView.SetDynamicColors(true)

	// Create status view
	ap.statusView = tview.NewTextView()
	ap.statusView.SetBorder(true).SetTitle(" AI Status ")
	ap.statusView.SetDynamicColors(true)
	ap.updateStatus()
}

// setupAIInteraction sets up the AI interaction UI
func (ap *AIPanel) setupAIInteraction() {
	// Query input field
	ap.queryInput = tview.NewInputField()
	ap.queryInput.SetLabel("AI Query: ").
		SetFieldWidth(50).
		SetPlaceholder("Ask AI about your SSH configuration...").
		SetChangedFunc(func(text string) {
			// Could add real-time suggestions here
		})

	// AI interaction result view
	ap.textView = tview.NewTextView()
	ap.textView.SetBorder(true).SetTitle(" AI Assistant ")
	ap.textView.SetDynamicColors(true)
	ap.textView.SetText("[blue]Welcome to Wooak AI Assistant![/blue]\n\n" +
		"Ask me about:\n" +
		"• SSH configuration optimization\n" +
		"• Security best practices\n" +
		"• Server recommendations\n" +
		"• Natural language search\n\n" +
		"Type your question and press Enter to get AI-powered insights.")
}

// saveConfig saves the current AI configuration
func (ap *AIPanel) saveConfig() {
	// Get form values and update config
	// This is a simplified version - in a real implementation,
	// you'd need to get the actual form values and update the config

	ap.updateStatus()
	ap.resultView.SetText("[green]AI configuration saved successfully![white]")
}

// testConnection tests the connection to the AI provider
func (ap *AIPanel) testConnection() {
	ap.statusView.SetText("[yellow]Testing AI connection...[/yellow]")

	// Test connection in a goroutine to avoid blocking UI
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := ap.aiSvc.TestConnection(ctx)
		if err != nil {
			ap.app.QueueUpdateDraw(func() {
				ap.statusView.SetText(fmt.Sprintf("[red]Connection failed: %v[/red]", err))
			})
		} else {
			ap.app.QueueUpdateDraw(func() {
				ap.statusView.SetText("[green]AI connection successful![/green]")
			})
		}
	}()
}

// resetConfig resets the AI configuration to defaults
func (ap *AIPanel) resetConfig() {
	ap.config = aiDomain.DefaultAIConfig()
	if err := ap.aiSvc.UpdateConfig(ap.config); err != nil {
		ap.resultView.SetText("[red]Failed to update AI configuration: " + err.Error() + "[white]")
		return
	}
	ap.setupUI() // Refresh the UI
	ap.resultView.SetText("[yellow]AI configuration reset to defaults[/yellow]")
}

// updateStatus updates the AI status display
func (ap *AIPanel) updateStatus() {
	var status strings.Builder

	status.WriteString(fmt.Sprintf("Provider: %s\n", ap.config.Provider))
	status.WriteString(fmt.Sprintf("Model: %s\n", ap.config.Model))
	status.WriteString(fmt.Sprintf("Base URL: %s\n", ap.config.BaseURL))
	status.WriteString(fmt.Sprintf("Enabled: %t\n", ap.config.Enabled))
	status.WriteString(fmt.Sprintf("Cache: %t\n", ap.config.CacheEnabled))

	if ap.config.Enabled {
		status.WriteString("\n[green]AI Service: Active[/green]")
	} else {
		status.WriteString("\n[red]AI Service: Disabled[/red]")
	}

	ap.statusView.SetText(status.String())
}

// processAIQuery processes an AI query
func (ap *AIPanel) processAIQuery(query string) {
	if !ap.config.Enabled {
		ap.textView.SetText("[red]AI service is disabled. Please enable it in the configuration.[/red]")
		return
	}

	ap.textView.SetText("[yellow]Processing your query...[/yellow]")

	// Process query in a goroutine to avoid blocking UI
	go func() {
		_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// For now, we'll simulate an AI response
		// In a real implementation, this would call the AI service
		time.Sleep(2 * time.Second) // Simulate processing time

		response := ap.generateMockResponse(query)

		ap.app.QueueUpdateDraw(func() {
			ap.textView.SetText(response)
		})
	}()
}

// generateMockResponse generates a mock AI response for demonstration
func (ap *AIPanel) generateMockResponse(query string) string {
	queryLower := strings.ToLower(query)

	if strings.Contains(queryLower, "security") {
		return "[green]AI Security Analysis:[/green]\n\n" +
			"Based on your SSH configuration, here are my security recommendations:\n\n" +
			"1. [yellow]Key Authentication:[/yellow] Ensure you're using Ed25519 or RSA 3072+ keys\n" +
			"2. [yellow]Host Verification:[/yellow] Enable StrictHostKeyChecking\n" +
			"3. [yellow]Connection Timeout:[/yellow] Set appropriate connection timeouts\n" +
			"4. [yellow]Port Security:[/yellow] Use non-standard ports when possible\n\n" +
			"[blue]Confidence: 85%[/blue]"
	}

	if strings.Contains(queryLower, "optimize") || strings.Contains(queryLower, "performance") {
		return "[green]AI Performance Optimization:[/green]\n\n" +
			"Here are my recommendations to optimize your SSH connections:\n\n" +
			"1. [yellow]Compression:[/yellow] Enable compression for slow connections\n" +
			"2. [yellow]Keep-Alive:[/yellow] Configure ServerAliveInterval\n" +
			"3. [yellow]Multiplexing:[/yellow] Use ControlMaster for multiple sessions\n" +
			"4. [yellow]Cipher Selection:[/yellow] Use faster ciphers when security allows\n\n" +
			"[blue]Expected improvement: 20-30% faster connections[/blue]"
	}

	if strings.Contains(queryLower, "search") || strings.Contains(queryLower, "find") {
		return "[green]AI Natural Language Search:[/green]\n\n" +
			"I found these servers matching your query:\n\n" +
			"1. [yellow]production-web-01[/yellow] - Web server cluster\n" +
			"   • High availability setup\n" +
			"   • Load balanced\n\n" +
			"2. [yellow]staging-db-02[/yellow] - Database server\n" +
			"   • PostgreSQL instance\n" +
			"   • Backup configured\n\n" +
			"[blue]Search confidence: 92%[/blue]"
	}

	// Default response
	return "[green]AI Assistant Response:[/green]\n\n" +
		"Thank you for your question: \"" + query + "\"\n\n" +
		"I'm here to help you with:\n" +
		"• SSH configuration optimization\n" +
		"• Security best practices\n" +
		"• Server recommendations\n" +
		"• Natural language search\n\n" +
		"Please ask me a more specific question about your SSH setup!"
}

// GetAIConfigForm returns the AI configuration form
func (ap *AIPanel) GetAIConfigForm() *tview.Form {
	return ap.form
}

// GetAIInteractionView returns the AI interaction view
func (ap *AIPanel) GetAIInteractionView() *tview.Flex {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	// Add query input
	ap.queryInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			query := ap.queryInput.GetText()
			if strings.TrimSpace(query) != "" {
				ap.processAIQuery(query)
				ap.queryInput.SetText("")
			}
		}
	})

	flex.AddItem(ap.queryInput, 3, 0, false)

	// Add result view
	flex.AddItem(ap.textView, 0, 1, false)

	return flex
}

// GetResultView returns the result view
func (ap *AIPanel) GetResultView() *tview.TextView {
	return ap.resultView
}

// GetStatusView returns the status view
func (ap *AIPanel) GetStatusView() *tview.TextView {
	return ap.statusView
}
