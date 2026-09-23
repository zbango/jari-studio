package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/agentic-app-studio/studio/internal/productgraph"
	_ "modernc.org/sqlite"
)

type SQLiteProjectStore struct {
	db *sql.DB
}

func OpenSQLiteProjectStore(ctx context.Context, path string) (*SQLiteProjectStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	store := &SQLiteProjectStore{db: db}
	if err := store.initialize(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *SQLiteProjectStore) initialize(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    brief TEXT NOT NULL DEFAULT '',
    revision INTEGER NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS revisions (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id),
    number INTEGER NOT NULL,
    status TEXT NOT NULL,
    snapshot TEXT NOT NULL,
    created_at TEXT NOT NULL,
    UNIQUE(project_id, number)
);`)
	return err
}

func (s *SQLiteProjectStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteProjectStore) CreateProject(ctx context.Context, project productgraph.Project) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO projects (id, name, brief, revision, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		project.ID,
		project.Name,
		project.Brief,
		project.Revision,
		project.CreatedAt.UTC().Format(time.RFC3339Nano),
		project.UpdatedAt.UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (s *SQLiteProjectStore) GetProject(ctx context.Context, id string) (productgraph.Project, error) {
	var project productgraph.Project
	var createdAt, updatedAt string
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, brief, revision, created_at, updated_at FROM projects WHERE id = ?`, id,
	).Scan(&project.ID, &project.Name, &project.Brief, &project.Revision, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return productgraph.Project{}, ErrProjectNotFound
	}
	if err != nil {
		return productgraph.Project{}, err
	}
	project.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return productgraph.Project{}, err
	}
	project.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return productgraph.Project{}, err
	}
	return project, nil
}

func (s *SQLiteProjectStore) ListProjects(ctx context.Context) ([]productgraph.Project, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, brief, revision, created_at, updated_at FROM projects ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]productgraph.Project, 0)
	for rows.Next() {
		var project productgraph.Project
		var createdAt, updatedAt string
		if err := rows.Scan(&project.ID, &project.Name, &project.Brief, &project.Revision, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		project.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
		if err != nil {
			return nil, err
		}
		project.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *SQLiteProjectStore) CreateRevision(ctx context.Context, revision productgraph.Revision) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO revisions (id, project_id, number, status, snapshot, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		revision.ID,
		revision.ProjectID,
		revision.Number,
		revision.Status,
		revision.Snapshot,
		revision.CreatedAt.UTC().Format(time.RFC3339Nano),
	)
	return err
}

func (s *SQLiteProjectStore) ListRevisions(ctx context.Context, projectID string) ([]productgraph.Revision, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, project_id, number, status, snapshot, created_at FROM revisions WHERE project_id = ? ORDER BY number`,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	revisions := make([]productgraph.Revision, 0)
	for rows.Next() {
		var revision productgraph.Revision
		var createdAt string
		if err := rows.Scan(&revision.ID, &revision.ProjectID, &revision.Number, &revision.Status, &revision.Snapshot, &createdAt); err != nil {
			return nil, err
		}
		revision.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
		if err != nil {
			return nil, err
		}
		revisions = append(revisions, revision)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return revisions, nil
}
