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
    "context"
    "fmt"
    "time"

    "go.uber.org/zap"

    "github.com/aryasoni98/wooak/internal/core/ports"
    aiService "github.com/aryasoni98/wooak/internal/core/services/ai"
    securityService "github.com/aryasoni98/wooak/internal/core/services/security"
    "github.com/rivo/tview"
)

type App interface {
    Run() error
}

type tui struct {
    logger *zap.SugaredLogger

    version string
    commit  string

    app           *tview.Application
    serverService ports.ServerService
    securitySvc   *securityService.SecurityService
    aiSvc         *aiService.AIService

    header     *AppHeader
    searchBar  *SearchBar
    hintBar    *tview.TextView
    serverList *ServerList
    details    *ServerDetails
    statusBar  *tview.TextView

    root    *tview.Flex
    left    *tview.Flex
    content *tview.Flex

    sortMode      SortMode
    searchVisible bool
}

func NewTUI(logger *zap.SugaredLogger, ss ports.ServerService, securitySvc *securityService.SecurityService, aiSvc *aiService.AIService, version, commit string) App {
    return &tui{
        logger:        logger,
        app:           tview.NewApplication(),
        serverService: ss,
        securitySvc:   securitySvc,
        aiSvc:         aiSvc,
        version:       version,
        commit:        commit,
    }
}

func (t *tui) Run() error {
    defer func() {
        if r := recover(); r != nil {
            t.logger.Errorw("panic recovered", "error", r)
            t.displayError(fmt.Errorf("unexpected error occurred: %v", r))
        }
    }()
    t.app.EnableMouse(true)
    t.initializeTheme().buildComponents().buildLayout().bindEvents().loadInitialData()
    t.app.SetRoot(t.root, true)
    t.logger.Infow("starting TUI application", "version", t.version, "commit", t.commit)
    if err := t.app.Run(); err != nil {
        t.logger.Errorw("application run error", "error", err)
        t.displayError(err)
        return err
    }
    return nil
}

func (t *tui) initializeTheme() *tui {
    // Apply modern theme
    theme := GetModernTheme()
    theme.ApplyTheme()
    return t
}

func (t *tui) buildComponents() *tui {
    t.header = NewAppHeader(t.version, t.commit, RepoURL)
    t.searchBar = NewSearchBar().
        OnSearch(t.handleSearchInput).
        OnEscape(t.hideSearchBar)
    t.hintBar = NewHintBar()
    t.serverList = NewServerList().
        OnSelectionChange(t.handleServerSelectionChange)
    t.details = NewServerDetails()
    t.statusBar = NewStatusBar()

    // default sort mode
    t.sortMode = SortByAliasAsc

    return t
}

func (t *tui) buildLayout() *tui {
    t.left = tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(t.hintBar, 1, 0, false).
        AddItem(t.serverList, 0, 1, true)

    right := tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(t.details, 0, 1, false)

    t.content = tview.NewFlex().SetDirection(tview.FlexColumn).
        AddItem(t.left, 0, 3, true).
        AddItem(right, 0, 2, false)

    t.root = tview.NewFlex().SetDirection(tview.FlexRow).
        AddItem(t.header, 2, 0, false).
        AddItem(t.content, 0, 1, true).
        AddItem(t.statusBar, 1, 0, false)
    return t
}

func (t *tui) bindEvents() *tui {
    t.root.SetInputCapture(t.handleGlobalKeys)
    return t
}

func (t *tui) loadInitialData() *tui {
    // Use context with timeout for network operation
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    t.showLoading("Loading servers...")

    servers, err := t.serverService.ListServersWithContext(ctx, "")
    t.hideLoading()

    if err != nil {
        t.displayError(fmt.Errorf("failed to load servers: %w", err))
        return t
    }

    sortServersForUI(servers, t.sortMode)
    t.updateListTitle()
    t.serverList.UpdateServers(servers)

    return t
}

func (t *tui) updateListTitle() {
    if t.serverList != nil {
        t.serverList.SetTitle(" Servers — Sort: " + t.sortMode.String() + " ")
    }
}

// showLoading displays a loading message in the status bar.
func (t *tui) showLoading(message string) {
    t.app.QueueUpdateDraw(func() {
        t.statusBar.SetText("[yellow]" + message + " [::b](press Esc to cancel)[-:-:-]")
        t.statusBar.Show()
    })
}

// hideLoading clears the status bar.
func (t *tui) hideLoading() {
    t.app.QueueUpdateDraw(func() {
        t.statusBar.SetText("")
        t.statusBar.Hide()
    })
}

// displayError centralizes error display in status bar and logs it.
func (t *tui) displayError(err error) {
    if err == nil {
        return
    }
    t.logger.Errorw("UI error", "error", err)
    t.app.QueueUpdateDraw(func() {
        t.statusBar.SetText("[red]Error: " + err.Error() + " [-]")
        t.statusBar.Show()
    })
}
