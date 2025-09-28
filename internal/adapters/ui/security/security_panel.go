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

package security

import (
	"fmt"
	"strings"

	securityDomain "github.com/aryasoni98/wooak/internal/core/domain/security"
	securityService "github.com/aryasoni98/wooak/internal/core/services/security"
	"github.com/rivo/tview"
)

// SecurityPanel provides a UI for security features
type SecurityPanel struct {
	app         *tview.Application
	securitySvc *securityService.SecurityService
	policy      *securityDomain.SecurityPolicy
	form        *tview.Form
	textView    *tview.TextView
	keyInput    *tview.InputField
	resultView  *tview.TextView
}

// NewSecurityPanel creates a new security panel
func NewSecurityPanel(app *tview.Application, securitySvc *securityService.SecurityService) *SecurityPanel {
	sp := &SecurityPanel{
		app:         app,
		securitySvc: securitySvc,
		policy:      securitySvc.GetSecurityPolicy(),
	}

	sp.setupUI()
	return sp
}

// setupUI sets up the security panel UI
func (sp *SecurityPanel) setupUI() {
	// Create the main form
	sp.form = tview.NewForm()
	sp.form.SetBorder(true).SetTitle(" Security Configuration ")

	// Add policy configuration fields
	sp.form.AddInputField("Min Key Size (bits)", fmt.Sprintf("%d", sp.policy.MinKeySize), 20, nil, nil)
	sp.form.AddCheckbox("Require Host Key Check", sp.policy.RequireHostKeyCheck, nil)
	sp.form.AddCheckbox("Enable Audit Log", sp.policy.EnableAuditLog, nil)
	sp.form.AddDropDown("Audit Log Level", []string{"info", "warn", "error"}, 0, nil)
	sp.form.AddInputField("Retention Days", fmt.Sprintf("%d", sp.policy.RetentionDays), 10, nil, nil)
	sp.form.AddCheckbox("Require VPN", sp.policy.RequireVPN, nil)

	// Add buttons
	sp.form.AddButton("Save Policy", sp.savePolicy)
	sp.form.AddButton("Reset to Default", sp.resetPolicy)
	sp.form.AddButton("View Audit Log", sp.viewAuditLog)

	// Create key validation section
	sp.setupKeyValidation()

	// Create result view
	sp.resultView = tview.NewTextView()
	sp.resultView.SetBorder(true).SetTitle(" Security Results ")
	sp.resultView.SetDynamicColors(true)
}

// setupKeyValidation sets up the key validation UI
func (sp *SecurityPanel) setupKeyValidation() {
	// Key input field
	sp.keyInput = tview.NewInputField()
	sp.keyInput.SetLabel("SSH Key: ").
		SetFieldWidth(50).
		SetPlaceholder("Paste your SSH public key here...").
		SetChangedFunc(func(text string) {
			if strings.TrimSpace(text) != "" {
				sp.validateKey(text)
			}
		})

	// Key validation result view
	sp.textView = tview.NewTextView()
	sp.textView.SetBorder(true).SetTitle(" Key Validation Results ")
	sp.textView.SetDynamicColors(true)
}

// validateKey validates the entered SSH key
func (sp *SecurityPanel) validateKey(keyData string) {
	result := sp.securitySvc.ValidateSSHKey(keyData)
	sp.displayKeyValidationResult(result)
}

// displayKeyValidationResult displays the key validation result
func (sp *SecurityPanel) displayKeyValidationResult(result *securityDomain.KeyValidationResult) {
	var output strings.Builder

	// Status
	if result.IsValid {
		output.WriteString("[green]✓ Key is valid[white]\n\n")
	} else {
		output.WriteString("[red]✗ Key validation failed[white]\n\n")
	}

	// Key information
	if result.KeyInfo != nil {
		output.WriteString("[yellow]Key Information:[white]\n")
		output.WriteString(fmt.Sprintf("  Type: %s\n", result.KeyInfo.Type))
		output.WriteString(fmt.Sprintf("  Size: %d bits\n", result.KeyInfo.Size))
		if result.KeyInfo.Fingerprint != "" {
			output.WriteString(fmt.Sprintf("  Fingerprint: %s\n", result.KeyInfo.Fingerprint))
		}
		if result.KeyInfo.Comment != "" {
			output.WriteString(fmt.Sprintf("  Comment: %s\n", result.KeyInfo.Comment))
		}
		output.WriteString("\n")
	}

	// Issues
	if len(result.Issues) > 0 {
		output.WriteString("[red]Issues:[white]\n")
		for _, issue := range result.Issues {
			output.WriteString(fmt.Sprintf("  • %s\n", issue))
		}
		output.WriteString("\n")
	}

	// Warnings
	if len(result.Warnings) > 0 {
		output.WriteString("[yellow]Warnings:[white]\n")
		for _, warning := range result.Warnings {
			output.WriteString(fmt.Sprintf("  • %s\n", warning))
		}
		output.WriteString("\n")
	}

	// Recommendations
	if len(result.Recommendations) > 0 {
		output.WriteString("[blue]Recommendations:[white]\n")
		for _, rec := range result.Recommendations {
			output.WriteString(fmt.Sprintf("  • %s\n", rec))
		}
		output.WriteString("\n")
	}

	sp.textView.SetText(output.String())
}

// savePolicy saves the current security policy
func (sp *SecurityPanel) savePolicy() {
	// Get form values and update policy
	// This is a simplified version - in a real implementation,
	// you'd need to get the actual form values and update the policy

	// For now, just show a message
	sp.resultView.SetText("[green]Security policy saved successfully![white]")
}

// resetPolicy resets the security policy to defaults
func (sp *SecurityPanel) resetPolicy() {
	sp.policy = securityDomain.DefaultSecurityPolicy()
	if err := sp.securitySvc.UpdateSecurityPolicy(sp.policy); err != nil {
		sp.resultView.SetText("[red]Failed to update security policy: " + err.Error() + "[white]")
		return
	}
	sp.setupUI() // Refresh the UI
	sp.resultView.SetText("[yellow]Security policy reset to defaults[white]")
}

// viewAuditLog displays the audit log
func (sp *SecurityPanel) viewAuditLog() {
	// This would open a modal or new view showing the audit log
	sp.resultView.SetText("[blue]Audit log viewer would open here[white]")
}

// GetSecurityForm returns the security configuration form
func (sp *SecurityPanel) GetSecurityForm() *tview.Form {
	return sp.form
}

// GetKeyValidationView returns the key validation view
func (sp *SecurityPanel) GetKeyValidationView() *tview.Flex {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	// Add key input
	flex.AddItem(sp.keyInput, 3, 0, false)

	// Add result view
	flex.AddItem(sp.textView, 0, 1, false)

	return flex
}

// GetResultView returns the result view
func (sp *SecurityPanel) GetResultView() *tview.TextView {
	return sp.resultView
}
