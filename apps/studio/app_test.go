package main

import (
	"context"
	"testing"

	"github.com/agentic-app-studio/studio/internal/storage"
)

func TestAppRevisionBinding(t *testing.T) {
	ctx := context.Background()
	store, err := storage.OpenSQLiteProjectStore(ctx, ":memory:")
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer store.Close()

	app := &App{ctx: ctx, store: store}
	project, err := app.CreateProject("Reference", "A business brief")
	if err != nil {
		t.Fatalf("create project: %v", err)
	}
	if _, err := app.CreateRevision(project.ID, "not-json"); err != ErrInvalidSnapshot {
		t.Fatalf("expected invalid snapshot error, got %v", err)
	}
	if _, err := app.CreateRevision(project.ID, `{"project":{"id":"other"}}`); err != ErrSnapshotProjectMismatch {
		t.Fatalf("expected project mismatch error, got %v", err)
	}

	revision, err := app.CreateRevision(project.ID, `{"project":{"id":"`+project.ID+`"}}`)
	if err != nil {
		t.Fatalf("create revision: %v", err)
	}
	if revision.Number != 1 || revision.ProjectID != project.ID {
		t.Fatalf("unexpected revision: %+v", revision)
	}
}
