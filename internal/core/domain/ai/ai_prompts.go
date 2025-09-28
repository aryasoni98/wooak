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
	"fmt"
	"strings"
)

// PromptTemplate represents a template for AI prompts
type PromptTemplate struct {
	Name        string        `json:"name"`
	Type        AIRequestType `json:"type"`
	Template    string        `json:"template"`
	Variables   []string      `json:"variables"`
	Description string        `json:"description"`
}

// GetPromptTemplates returns all available prompt templates
func GetPromptTemplates() map[string]*PromptTemplate {
	return map[string]*PromptTemplate{
		"server_recommendation": {
			Name:        "Server Recommendation",
			Type:        RequestTypeRecommendation,
			Template:    "Based on the SSH server configuration below, provide security and performance recommendations. Focus on best practices and potential improvements.\n\nServer Config:\n{{.config}}\n\nUser Context: {{.context}}\n\nProvide specific, actionable recommendations with priority levels.",
			Variables:   []string{"config", "context"},
			Description: "Generates recommendations for SSH server configurations",
		},
		"natural_language_search": {
			Name:        "Natural Language Search",
			Type:        RequestTypeSearch,
			Template:    "Find SSH servers that match this natural language description: \"{{.query}}\"\n\nAvailable servers:\n{{.servers}}\n\nReturn a ranked list of matching servers with explanations.",
			Variables:   []string{"query", "servers"},
			Description: "Searches servers using natural language queries",
		},
		"security_analysis": {
			Name:        "Security Analysis",
			Type:        RequestTypeSecurity,
			Template:    "Analyze the following SSH configuration for security vulnerabilities and best practices:\n\n{{.config}}\n\nProvide a security assessment with:\n1. Identified vulnerabilities\n2. Risk levels\n3. Recommended fixes\n4. Best practice suggestions",
			Variables:   []string{"config"},
			Description: "Performs security analysis of SSH configurations",
		},
		"connection_optimization": {
			Name:        "Connection Optimization",
			Type:        RequestTypeOptimization,
			Template:    "Optimize this SSH connection configuration for better performance and reliability:\n\n{{.config}}\n\nCurrent issues: {{.issues}}\n\nProvide optimization suggestions with expected improvements.",
			Variables:   []string{"config", "issues"},
			Description: "Optimizes SSH connection configurations",
		},
		"intelligent_suggestions": {
			Name:        "Intelligent Suggestions",
			Type:        RequestTypeGeneral,
			Template:    "Based on the user's SSH usage patterns and current configuration, suggest improvements and new features that would be helpful.\n\nUsage patterns: {{.patterns}}\nCurrent config: {{.config}}\n\nProvide personalized suggestions.",
			Variables:   []string{"patterns", "config"},
			Description: "Provides intelligent suggestions based on usage patterns",
		},
	}
}

// BuildPrompt builds a prompt from a template with variables
func (pt *PromptTemplate) BuildPrompt(variables map[string]string) string {
	prompt := pt.Template

	for _, variable := range pt.Variables {
		if value, exists := variables[variable]; exists {
			placeholder := fmt.Sprintf("{{.%s}}", variable)
			prompt = strings.ReplaceAll(prompt, placeholder, value)
		}
	}

	return prompt
}

// GetPromptForType returns the appropriate prompt template for a request type
func GetPromptForType(requestType AIRequestType) *PromptTemplate {
	templates := GetPromptTemplates()

	switch requestType {
	case RequestTypeRecommendation:
		return templates["server_recommendation"]
	case RequestTypeSearch:
		return templates["natural_language_search"]
	case RequestTypeSecurity:
		return templates["security_analysis"]
	case RequestTypeOptimization:
		return templates["connection_optimization"]
	case RequestTypeGeneral:
		return templates["intelligent_suggestions"]
	default:
		return templates["intelligent_suggestions"]
	}
}

// FormatServerConfig formats server configuration for AI prompts
func FormatServerConfig(config map[string]string) string {
	var parts []string
	for key, value := range config {
		if value != "" {
			parts = append(parts, fmt.Sprintf("%s: %s", key, value))
		}
	}
	return strings.Join(parts, "\n")
}

// FormatServerList formats a list of servers for AI prompts
func FormatServerList(servers []map[string]string) string {
	parts := make([]string, 0, len(servers))
	for i, server := range servers {
		serverInfo := fmt.Sprintf("Server %d:", i+1)
		for key, value := range server {
			if value != "" {
				serverInfo += fmt.Sprintf("\n  %s: %s", key, value)
			}
		}
		parts = append(parts, serverInfo)
	}
	return strings.Join(parts, "\n\n")
}
