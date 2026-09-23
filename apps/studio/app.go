package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/agentic-app-studio/studio/internal/agentadapter"
	"github.com/agentic-app-studio/studio/internal/compiler"
	"github.com/agentic-app-studio/studio/internal/execution"
	"github.com/agentic-app-studio/studio/internal/productgraph"
	"github.com/agentic-app-studio/studio/internal/scheduler"
	"github.com/agentic-app-studio/studio/internal/storage"
	"github.com/agentic-app-studio/studio/internal/terminal"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var ErrStoreUnavailable = errors.New("studio project store is unavailable")
var ErrInvalidSnapshot = errors.New("revision snapshot must be valid JSON")
var ErrSnapshotProjectMismatch = errors.New("revision snapshot project does not match target project")

type App struct {
	ctx        context.Context
	store      *storage.SQLiteProjectStore
	storeErr   error
	terminalMu sync.RWMutex
	terminals  map[string]*terminalRecord
}

type terminalRecord struct {
	info    terminal.Info
	session terminal.Session
	events  []terminal.Event
}

func NewApp() *App {
	return &App{terminals: make(map[string]*terminalRecord)}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	dataDir, err := os.UserConfigDir()
	if err != nil {
		return
	}
	databaseDir := filepath.Join(dataDir, "AgenticAppStudio")
	if err := os.MkdirAll(databaseDir, 0o700); err != nil {
		return
	}
	a.store, a.storeErr = storage.OpenSQLiteProjectStore(ctx, filepath.Join(databaseDir, "studio.db"))
}

func (a *App) shutdown(_ context.Context) {
	if a.store != nil {
		_ = a.store.Close()
	}
	a.terminalMu.Lock()
	defer a.terminalMu.Unlock()
	for id, record := range a.terminals {
		_ = record.session.Close()
		delete(a.terminals, id)
	}
}

func (a *App) Health() map[string]string {
	status := "ready"
	if a.store == nil || a.storeErr != nil {
		status = "store_unavailable"
	}
	return map[string]string{
		"status":  status,
		"runtime": "wails-go",
	}
}

func (a *App) DefaultWorkspaceDirectory() string {
	current, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(current, "go.mod")); err == nil {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			return current
		}
		current = parent
	}
}

func (a *App) SelectWorkspaceDirectory() (string, error) {
	if a.ctx == nil {
		return "", errors.New("native studio runtime is not ready")
	}
	return wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title:                "Choose a workspace folder",
		DefaultDirectory:     a.DefaultWorkspaceDirectory(),
		ShowHiddenFiles:      false,
		CanCreateDirectories: false,
	})
}

func (a *App) CreateProject(name, brief string) (productgraph.Project, error) {
	if a.store == nil || a.storeErr != nil {
		return productgraph.Project{}, ErrStoreUnavailable
	}
	now := time.Now().UTC()
	project := productgraph.Project{
		ID:        "project." + now.Format("20060102T150405.000000000Z"),
		Name:      name,
		Brief:     brief,
		Revision:  1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := a.store.CreateProject(a.ctx, project); err != nil {
		return productgraph.Project{}, err
	}
	return project, nil
}

func (a *App) ListProjects() ([]productgraph.Project, error) {
	if a.store == nil || a.storeErr != nil {
		return nil, ErrStoreUnavailable
	}
	return a.store.ListProjects(a.ctx)
}

func (a *App) CreateRevision(projectID, snapshot string) (productgraph.Revision, error) {
	if a.store == nil || a.storeErr != nil {
		return productgraph.Revision{}, ErrStoreUnavailable
	}
	if !json.Valid([]byte(snapshot)) {
		return productgraph.Revision{}, ErrInvalidSnapshot
	}
	var envelope struct {
		Project struct {
			ID string `json:"id"`
		} `json:"project"`
	}
	if err := json.Unmarshal([]byte(snapshot), &envelope); err != nil {
		return productgraph.Revision{}, ErrInvalidSnapshot
	}
	if envelope.Project.ID != "" && envelope.Project.ID != projectID {
		return productgraph.Revision{}, ErrSnapshotProjectMismatch
	}
	if _, err := a.store.GetProject(a.ctx, projectID); err != nil {
		return productgraph.Revision{}, err
	}
	revisions, err := a.store.ListRevisions(a.ctx, projectID)
	if err != nil {
		return productgraph.Revision{}, err
	}
	now := time.Now().UTC()
	revision := productgraph.Revision{
		ID:        fmt.Sprintf("%s.revision.%d", projectID, len(revisions)+1),
		ProjectID: projectID,
		Number:    len(revisions) + 1,
		Status:    productgraph.StatusProposed,
		Snapshot:  snapshot,
		CreatedAt: now,
	}
	if err := a.store.CreateRevision(a.ctx, revision); err != nil {
		return productgraph.Revision{}, err
	}
	return revision, nil
}

func (a *App) ListRevisions(projectID string) ([]productgraph.Revision, error) {
	if a.store == nil || a.storeErr != nil {
		return nil, ErrStoreUnavailable
	}
	return a.store.ListRevisions(a.ctx, projectID)
}

func (a *App) CompileRevision(projectID, revisionID string) (productgraph.CompiledPlan, error) {
	if a.store == nil || a.storeErr != nil {
		return productgraph.CompiledPlan{}, ErrStoreUnavailable
	}
	if _, err := a.store.GetProject(a.ctx, projectID); err != nil {
		return productgraph.CompiledPlan{}, err
	}
	revisions, err := a.store.ListRevisions(a.ctx, projectID)
	if err != nil {
		return productgraph.CompiledPlan{}, err
	}
	for _, revision := range revisions {
		if revision.ID == revisionID {
			return compiler.Compile(revision)
		}
	}
	return productgraph.CompiledPlan{}, storage.ErrRevisionNotFound
}

func (a *App) EvaluateSchedule(projectID, revisionID, adapterID string, capabilities []string) (productgraph.ScheduleResult, error) {
	plan, err := a.CompileRevision(projectID, revisionID)
	if err != nil {
		return productgraph.ScheduleResult{}, err
	}
	return scheduler.Evaluate(plan.Cards, productgraph.AdapterManifest{ID: adapterID, Version: "local", Transport: "structured", Capabilities: capabilities}), nil
}

func (a *App) AvailableAdapters() []agentadapter.Manifest {
	return agentadapter.BuiltInManifests()
}

func (a *App) ExecuteCard(cardJSON, repository, worktreeRoot, baseRef, adapterID string) (execution.RunResult, error) {
	var card productgraph.CardSpec
	if err := json.Unmarshal([]byte(cardJSON), &card); err != nil {
		return execution.RunResult{}, fmt.Errorf("invalid card JSON: %w", err)
	}
	adapter, err := agentadapter.NewBuiltInAdapter(adapterID)
	if err != nil {
		return execution.RunResult{}, err
	}
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	return execution.ExecuteCard(ctx, card, repository, worktreeRoot, baseRef, adapter)
}

func (a *App) StartTerminal(command string, args []string, directory string) (terminal.Info, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	session, err := (terminal.Manager{}).Start(ctx, command, args, directory, nil, 120, 40)
	if err != nil {
		return terminal.Info{}, err
	}
	id := fmt.Sprintf("terminal.%d", time.Now().UnixNano())
	record := &terminalRecord{info: terminal.Info{ID: id, Command: command, Directory: directory, State: terminal.StateRunning}, session: session}
	a.terminalMu.Lock()
	a.terminals[id] = record
	a.terminalMu.Unlock()
	go func() {
		for event := range session.Events() {
			a.terminalMu.Lock()
			if current, ok := a.terminals[id]; ok {
				if len(current.events) >= 500 {
					current.events = current.events[1:]
				}
				current.events = append(current.events, event)
			}
			a.terminalMu.Unlock()
		}
		err := session.Wait()
		a.terminalMu.Lock()
		if current, ok := a.terminals[id]; ok {
			if err != nil {
				current.info.State = terminal.StateFailed
			} else {
				current.info.State = terminal.StateExited
			}
		}
		a.terminalMu.Unlock()
	}()
	return record.info, nil
}

func (a *App) ListTerminals() []terminal.Info {
	a.terminalMu.RLock()
	defer a.terminalMu.RUnlock()
	terminals := make([]terminal.Info, 0, len(a.terminals))
	for _, record := range a.terminals {
		terminals = append(terminals, record.info)
	}
	return terminals
}

func (a *App) TerminalEvents(id string) []terminal.Event {
	a.terminalMu.RLock()
	defer a.terminalMu.RUnlock()
	record, ok := a.terminals[id]
	if !ok {
		return nil
	}
	return append([]terminal.Event(nil), record.events...)
}

func (a *App) WriteTerminalInput(id, input string) error {
	a.terminalMu.RLock()
	record, ok := a.terminals[id]
	a.terminalMu.RUnlock()
	if !ok {
		return errors.New("terminal session not found")
	}
	return record.session.WriteInput([]byte(input))
}

func (a *App) ResizeTerminal(id string, cols, rows uint16) error {
	a.terminalMu.RLock()
	record, ok := a.terminals[id]
	a.terminalMu.RUnlock()
	if !ok {
		return errors.New("terminal session not found")
	}
	return record.session.Resize(cols, rows)
}

func (a *App) InterruptTerminal(id string) error {
	a.terminalMu.RLock()
	record, ok := a.terminals[id]
	a.terminalMu.RUnlock()
	if !ok {
		return errors.New("terminal session not found")
	}
	return record.session.Interrupt()
}

func (a *App) CloseTerminal(id string) error {
	a.terminalMu.Lock()
	record, ok := a.terminals[id]
	if ok {
		delete(a.terminals, id)
	}
	a.terminalMu.Unlock()
	if !ok {
		return errors.New("terminal session not found")
	}
	return record.session.Close()
}
