package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/agentic-app-studio/studio/internal/productgraph"
	"github.com/agentic-app-studio/studio/internal/storage"
)

func main() {
	store, err := storage.OpenSQLiteProjectStore(context.Background(), ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	project := productgraph.Project{
		ID:        "project.bootstrap",
		Name:      "Agentic App Studio",
		Revision:  1,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := store.CreateProject(context.Background(), project); err != nil {
		log.Fatal(err)
	}
	projects, err := store.ListProjects(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("studio bootstrap ready: %d project(s)\n", len(projects))
}
