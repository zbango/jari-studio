package storage

import (
	"context"
	"errors"

	"github.com/agentic-app-studio/studio/internal/productgraph"
)

var ErrRevisionNotFound = errors.New("revision not found")

type ProjectStore interface {
	CreateProject(context.Context, productgraph.Project) error
	GetProject(context.Context, string) (productgraph.Project, error)
	ListProjects(context.Context) ([]productgraph.Project, error)
	CreateRevision(context.Context, productgraph.Revision) error
	ListRevisions(context.Context, string) ([]productgraph.Revision, error)
}
