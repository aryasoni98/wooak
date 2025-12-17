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
	"fmt"
)

// Tab management functions for ServerForm
// These functions handle tab navigation, rendering, and display logic.

func (sf *ServerForm) titleForMode() string {
	if sf.mode == ServerFormEdit {
		return "Edit Server"
	}
	return "Add Server"
}

func (sf *ServerForm) getCurrentTabIndex() int {
	for i, tab := range sf.tabs {
		if tab == sf.currentTab {
			return i
		}
	}
	return 0
}

func (sf *ServerForm) calculateTabsWidth(useAbbrev bool) int {
	width := 0
	for i, tab := range sf.tabs {
		tabName := tab
		if useAbbrev {
			tabName = sf.tabAbbrev[tab]
		}
		width += len(tabName) + 2 // space + name + space
		if i < len(sf.tabs)-1 {
			width += 3 // " | " separator
		}
	}
	return width
}

func (sf *ServerForm) determineDisplayMode(width int) string {
	if width <= 20 { // Width unknown or too small
		return "full"
	}

	fullWidth := sf.calculateTabsWidth(false)
	if fullWidth <= width-10 {
		return "full"
	}

	abbrevWidth := sf.calculateTabsWidth(true)
	if abbrevWidth <= width-10 {
		return "abbrev"
	}

	return "scroll"
}

func (sf *ServerForm) renderTab(tab string, isCurrent bool, useAbbrev bool, index int) string {
	tabName := tab
	if useAbbrev {
		tabName = sf.tabAbbrev[tab]
	}
	regionID := fmt.Sprintf("tab_%d", index)
	if isCurrent {
		return fmt.Sprintf("[%q][black:white:b] %s [-:-:-][%q] ", regionID, tabName, "")
	}
	return fmt.Sprintf("[%q][gray::u] %s [-:-:-][%q] ", regionID, tabName, "")
}

func (sf *ServerForm) renderScrollableTabs(currentIdx, width int) string {
	var tabText string
	availableWidth := width - 8 // Reserve space for scroll indicators

	// Calculate visible count
	visibleCount := sf.calculateVisibleTabCount(availableWidth)
	if visibleCount < 2 {
		visibleCount = 2
	}

	// Add left scroll indicator
	if currentIdx > 0 {
		tabText = "[gray]◀ [-]"
	}

	// Calculate range
	start, end := sf.calculateVisibleRange(currentIdx, visibleCount, len(sf.tabs))

	// Render visible tabs
	for i := start; i < end && i < len(sf.tabs); i++ {
		tabText += sf.renderTab(sf.tabs[i], sf.tabs[i] == sf.currentTab, true, i)
		if i < end-1 && i < len(sf.tabs)-1 {
			tabText += tabSeparator
		}
	}

	// Add right scroll indicator
	if currentIdx < len(sf.tabs)-1 {
		tabText += " [gray]▶[-]"
	}

	return tabText
}

func (sf *ServerForm) calculateVisibleTabCount(availableWidth int) int {
	visibleCount := 0
	currentWidth := 0

	for i := 0; i < len(sf.tabs) && currentWidth < availableWidth; i++ {
		abbrev := sf.tabAbbrev[sf.tabs[i]]
		tabWidth := len(abbrev) + 2
		if i > 0 {
			tabWidth += 3 // separator
		}
		if currentWidth+tabWidth <= availableWidth {
			visibleCount++
			currentWidth += tabWidth
		} else {
			break
		}
	}

	return visibleCount
}

func (sf *ServerForm) calculateVisibleRange(currentIdx, visibleCount, totalTabs int) (int, int) {
	halfVisible := visibleCount / 2
	start := currentIdx - halfVisible + 1
	end := start + visibleCount

	// Adjust boundaries
	if start < 0 {
		start = 0
		end = visibleCount
	}
	if end > totalTabs {
		end = totalTabs
		start = end - visibleCount
		if start < 0 {
			start = 0
		}
	}

	return start, end
}

func (sf *ServerForm) updateTabBar() {
	currentIdx := sf.getCurrentTabIndex()

	// Build tab text with scroll indicator if needed
	var tabText string

	// Check if we need to show scroll indicators
	x, y, width, height := sf.tabBar.GetInnerRect()
	_ = x
	_ = y
	_ = height

	displayMode := sf.determineDisplayMode(width)

	switch displayMode {
	case "scroll":
		tabText = sf.renderScrollableTabs(currentIdx, width)
	case "abbrev":
		// Show all tabs with abbreviated names
		for i, tab := range sf.tabs {
			tabText += sf.renderTab(tab, tab == sf.currentTab, true, i)
			if i < len(sf.tabs)-1 {
				tabText += tabSeparator
			}
		}
	default: // "full"
		// Show all tabs with full names
		for i, tab := range sf.tabs {
			tabText += sf.renderTab(tab, tab == sf.currentTab, false, i)
			if i < len(sf.tabs)-1 {
				tabText += tabSeparator
			}
		}
	}

	sf.tabBar.SetText(tabText)

	// Set up mouse click handler using highlight regions
	sf.tabBar.SetHighlightedFunc(func(added, removed, remaining []string) {
		if len(added) > 0 {
			// Extract tab index from region ID (format: "tab_0", "tab_1", etc)
			for _, regionID := range added {
				if len(regionID) > 4 && regionID[:4] == "tab_" {
					idx := int(regionID[4] - '0')
					if idx < len(sf.tabs) {
						sf.switchToTab(sf.tabs[idx])
					}
				}
			}
		}
	})
}

func (sf *ServerForm) switchToTab(tabName string) {
	for _, tab := range sf.tabs {
		if tab != tabName {
			continue
		}

		sf.currentTab = tabName
		sf.pages.SwitchToPage(tabName)
		sf.updateTabBar()

		// Set focus to the form in the newly selected tab
		if form, exists := sf.forms[tabName]; exists && sf.app != nil {
			sf.app.SetFocus(form)
		}
		break
	}
}

func (sf *ServerForm) nextTab() {
	for i, tab := range sf.tabs {
		if tab == sf.currentTab {
			// Loop to first tab if at the last tab
			if i == len(sf.tabs)-1 {
				sf.switchToTab(sf.tabs[0])
			} else {
				sf.switchToTab(sf.tabs[i+1])
			}
			break
		}
	}
}

func (sf *ServerForm) prevTab() {
	for i, tab := range sf.tabs {
		if tab == sf.currentTab {
			// Loop to last tab if at the first tab
			if i == 0 {
				sf.switchToTab(sf.tabs[len(sf.tabs)-1])
			} else {
				sf.switchToTab(sf.tabs[i-1])
			}
			break
		}
	}
}
