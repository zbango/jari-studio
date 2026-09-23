package storage

import (
	"context"
	"testing"
	"time"

	"github.com/agentic-app-studio/studio/internal/productgraph"
)

func TestMemoryProjectStoreRoundTrip(t *testing.T) {
	ctx := context.Background()
	store := NewMemoryProjectStore()
	project := productgraph.Project{ID: "project-1", Name: "Reference", Revision: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	if err := store.CreateProject(ctx, project); err != nil {
		t.Fatalf("create project: %v", err)
	}
	got, err := store.GetProject(ctx, project.ID)
	if err != nil {
		t.Fatalf("get project: %v", err)
	}
	if got != project {
		t.Fatalf("got %+v, want %+v", got, project)
	}
}
