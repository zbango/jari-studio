package storage

import (
	"context"
	"errors"
	"sort"
	"sync"

	"github.com/agentic-app-studio/studio/internal/productgraph"
)

var ErrProjectNotFound = errors.New("project not found")

type MemoryProjectStore struct {
	mu        sync.RWMutex
	projects  map[string]productgraph.Project
	revisions map[string][]productgraph.Revision
}

func NewMemoryProjectStore() *MemoryProjectStore {
	return &MemoryProjectStore{
		projects:  make(map[string]productgraph.Project),
		revisions: make(map[string][]productgraph.Revision),
	}
}

func (s *MemoryProjectStore) CreateProject(_ context.Context, project productgraph.Project) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.projects[project.ID]; exists {
		return errors.New("project already exists")
	}
	s.projects[project.ID] = project
	return nil
}

func (s *MemoryProjectStore) GetProject(_ context.Context, id string) (productgraph.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	project, ok := s.projects[id]
	if !ok {
		return productgraph.Project{}, ErrProjectNotFound
	}
	return project, nil
}

func (s *MemoryProjectStore) ListProjects(_ context.Context) ([]productgraph.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	projects := make([]productgraph.Project, 0, len(s.projects))
	for _, project := range s.projects {
		projects = append(projects, project)
	}
	sort.Slice(projects, func(i, j int) bool { return projects[i].ID < projects[j].ID })
	return projects, nil
}

func (s *MemoryProjectStore) CreateRevision(_ context.Context, revision productgraph.Revision) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.projects[revision.ProjectID]; !exists {
		return ErrProjectNotFound
	}
	for _, existing := range s.revisions[revision.ProjectID] {
		if existing.ID == revision.ID || existing.Number == revision.Number {
			return errors.New("revision already exists")
		}
	}
	s.revisions[revision.ProjectID] = append(s.revisions[revision.ProjectID], revision)
	return nil
}

func (s *MemoryProjectStore) ListRevisions(_ context.Context, projectID string) ([]productgraph.Revision, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, exists := s.projects[projectID]; !exists {
		return nil, ErrProjectNotFound
	}
	revisions := append([]productgraph.Revision(nil), s.revisions[projectID]...)
	sort.Slice(revisions, func(i, j int) bool { return revisions[i].Number < revisions[j].Number })
	return revisions, nil
}
