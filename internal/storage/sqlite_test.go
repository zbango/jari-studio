package storage

import (
	"context"
	"testing"
	"time"

	"github.com/agentic-app-studio/studio/internal/productgraph"
)

func TestSQLiteProjectStoreRoundTrip(t *testing.T) {
	ctx := context.Background()
	store, err := OpenSQLiteProjectStore(ctx, ":memory:")
	if err != nil {
		t.Fatalf("open sqlite store: %v", err)
	}
	defer store.Close()

	created := time.Date(2026, 9, 2, 12, 0, 0, 0, time.UTC)
	project := productgraph.Project{
		ID: "project-sqlite", Name: "SQLite Project", Brief: "Offline-first", Revision: 2,
		CreatedAt: created, UpdatedAt: created,
	}
	if err := store.CreateProject(ctx, project); err != nil {
		t.Fatalf("create project: %v", err)
	}
	got, err := store.GetProject(ctx, project.ID)
	if err != nil {
		t.Fatalf("get project: %v", err)
	}
	if !got.CreatedAt.Equal(project.CreatedAt) || got.ID != project.ID || got.Brief != project.Brief {
		t.Fatalf("got %+v, want %+v", got, project)
	}

	revision := productgraph.Revision{
		ID: "revision-1", ProjectID: project.ID, Number: 1,
		Status: productgraph.StatusProposed, Snapshot: `{"project":{"id":"project-sqlite"}}`, CreatedAt: created,
	}
	if err := store.CreateRevision(ctx, revision); err != nil {
		t.Fatalf("create revision: %v", err)
	}
	revisions, err := store.ListRevisions(ctx, project.ID)
	if err != nil {
		t.Fatalf("list revisions: %v", err)
	}
	if len(revisions) != 1 || revisions[0].Snapshot != revision.Snapshot {
		t.Fatalf("got revisions %+v, want one matching revision", revisions)
	}
}
