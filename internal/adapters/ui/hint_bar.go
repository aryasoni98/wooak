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

func NewHintBar() *tview.TextView {
	hint := tview.NewTextView().SetDynamicColors(true)
	hint.SetBackgroundColor(tcell.Color236)
	hint.SetText("[#AAAAAA]Press [#39BFFF::b]/[-:-:b] to search…  •  [#39BFFF]↑↓[-] Navigate  •  [#39BFFF]Enter[-] SSH  •  [#39BFFF]c[-] Copy  •  [#39BFFF]g[-] Ping  •  [#39BFFF]r[-] Refresh  •  [#39BFFF]a[-] Add  •  [#39BFFF]e[-] Edit  •  [#39BFFF]t[-] Tags  •  [#39BFFF]d[-] Delete  •  [#39BFFF]p[-] Pin  •  [#39BFFF]s[-] Sort  •  [#39BFFF]z[-] Security  •  [#39BFFF]i[-] AI Assistant[-]")
	return hint
}
