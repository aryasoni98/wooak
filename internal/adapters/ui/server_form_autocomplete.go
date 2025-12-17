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
	"strings"
)

// Autocomplete functions for ServerForm
// These functions provide intelligent autocomplete suggestions for various form fields.

// createSSHKeyAutocomplete creates an autocomplete function for SSH key file paths
func (sf *ServerForm) createSSHKeyAutocomplete() func(string) []string {
	return func(currentText string) []string {
		if currentText == "" {
			// Show available keys when field is empty
			availableKeys := GetAvailableSSHKeys()
			if len(availableKeys) == 0 {
				return nil
			}
			return availableKeys
		}

		// Split by comma to handle multiple keys
		keys := strings.Split(currentText, ",")
		lastKey := strings.TrimSpace(keys[len(keys)-1])

		// If the last key is empty (after a comma and space), show all available keys
		if lastKey == "" {
			availableKeys := GetAvailableSSHKeys()
			if len(availableKeys) == 0 {
				return nil
			}
			// Build suggestions with existing keys
			var suggestions []string
			prefix := ""
			if len(keys) > 1 {
				// Join all keys except the last empty one
				existingKeys := keys[:len(keys)-1]
				for i := range existingKeys {
					existingKeys[i] = strings.TrimSpace(existingKeys[i])
				}
				prefix = strings.Join(existingKeys, ", ") + ", "
			}
			for _, key := range availableKeys {
				suggestions = append(suggestions, prefix+key)
			}
			return suggestions
		}

		// Get available keys and filter based on what's being typed
		availableKeys := GetAvailableSSHKeys()
		if len(availableKeys) == 0 {
			return nil
		}

		// Convert to lowercase for case-insensitive matching
		searchTerm := strings.ToLower(lastKey)

		// Filter available keys
		var filtered []string
		prefix := ""
		if len(keys) > 1 {
			// Join all keys except the last one being typed
			existingKeys := keys[:len(keys)-1]
			for i := range existingKeys {
				existingKeys[i] = strings.TrimSpace(existingKeys[i])
			}
			prefix = strings.Join(existingKeys, ", ") + ", "
		}

		for _, key := range availableKeys {
			lowerKey := strings.ToLower(key)
			// Check if the key matches the search term
			if strings.Contains(lowerKey, searchTerm) || matchesSequence(lowerKey, searchTerm) {
				filtered = append(filtered, prefix+key)
			}
		}

		// If no matches found, return nil to allow Tab navigation
		if len(filtered) == 0 {
			return nil
		}

		return filtered
	}
}

// createKnownHostsAutocomplete creates an autocomplete function for known_hosts file paths
func (sf *ServerForm) createKnownHostsAutocomplete() func(string) []string {
	return func(currentText string) []string {
		if currentText == "" {
			// Show available known_hosts files when field is empty
			availableFiles := GetAvailableKnownHostsFiles()
			if len(availableFiles) == 0 {
				return nil
			}
			return availableFiles
		}

		// Split by space to handle multiple files
		files := strings.Split(currentText, " ")
		lastFile := strings.TrimSpace(files[len(files)-1])

		// If the last file is empty (after a space), show all available files
		if lastFile == "" {
			availableFiles := GetAvailableKnownHostsFiles()
			if len(availableFiles) == 0 {
				return nil
			}
			// Build suggestions with existing files
			var suggestions []string
			prefix := ""
			if len(files) > 1 {
				// Join all files except the last empty one
				existingFiles := files[:len(files)-1]
				for i := range existingFiles {
					existingFiles[i] = strings.TrimSpace(existingFiles[i])
				}
				prefix = strings.Join(existingFiles, " ") + " "
			}
			for _, file := range availableFiles {
				suggestions = append(suggestions, prefix+file)
			}
			return suggestions
		}

		// Get available files and filter based on what's being typed
		availableFiles := GetAvailableKnownHostsFiles()
		if len(availableFiles) == 0 {
			return nil
		}

		// Convert to lowercase for case-insensitive matching
		searchTerm := strings.ToLower(lastFile)

		// Filter available files
		var filtered []string
		prefix := ""
		if len(files) > 1 {
			// Join all files except the last one being typed
			existingFiles := files[:len(files)-1]
			for i := range existingFiles {
				existingFiles[i] = strings.TrimSpace(existingFiles[i])
			}
			prefix = strings.Join(existingFiles, " ") + " "
		}

		for _, file := range availableFiles {
			lowerFile := strings.ToLower(file)
			// Check if the file matches the search term
			if strings.Contains(lowerFile, searchTerm) || matchesSequence(lowerFile, searchTerm) {
				filtered = append(filtered, prefix+file)
			}
		}

		// If no matches found, return nil to allow Tab navigation
		if len(filtered) == 0 {
			return nil
		}

		return filtered
	}
}

// createAlgorithmAutocomplete creates an autocomplete function for algorithm input fields
func (sf *ServerForm) createAlgorithmAutocomplete(suggestions []string) func(string) []string {
	return func(currentText string) []string {
		if currentText == "" {
			// Return nil when empty to disable autocomplete, allowing Tab to navigate
			return nil
		}

		// Find the current word being typed
		words := strings.Split(currentText, ",")
		lastWord := strings.TrimSpace(words[len(words)-1])

		// If the last word is empty (after a comma), return nil to allow Tab navigation
		if lastWord == "" {
			return nil
		}

		// Handle prefix characters
		prefix := ""
		searchTerm := lastWord
		if lastWord != "" {
			if lastWord[0] == '+' || lastWord[0] == '-' || lastWord[0] == '^' {
				prefix = string(lastWord[0])
				if len(lastWord) > 1 {
					searchTerm = lastWord[1:]
				} else {
					// Just a prefix character, show all suggestions
					searchTerm = ""
				}
			}
		}

		// Filter suggestions - check if all characters appear in sequence
		var filtered []string
		for _, s := range suggestions {
			if searchTerm == "" || matchesSequence(strings.ToLower(s), strings.ToLower(searchTerm)) {
				// Build the complete text with the suggestion
				newWords := make([]string, len(words)-1)
				copy(newWords, words[:len(words)-1])
				newWords = append(newWords, prefix+s)
				filtered = append(filtered, strings.Join(newWords, ","))
			}
		}

		// If no matches found, return nil to allow Tab navigation
		if len(filtered) == 0 {
			return nil
		}

		return filtered
	}
}
