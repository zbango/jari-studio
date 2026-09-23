package execution

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

type WorktreeManager struct {
	GitBinary string
}

func (m WorktreeManager) Add(ctx context.Context, repository, worktree, branch, baseRef string) error {
	if repository == "" || worktree == "" || branch == "" || baseRef == "" {
		return fmt.Errorf("repository, worktree, branch and base ref are required")
	}
	cleanWorktree := filepath.Clean(worktree)
	command := exec.CommandContext(ctx, m.binary(), "-C", repository, "worktree", "add", "-b", branch, cleanWorktree, baseRef)
	if output, err := command.CombinedOutput(); err != nil {
		return fmt.Errorf("git worktree add: %w: %s", err, output)
	}
	return nil
}

func (m WorktreeManager) Remove(ctx context.Context, repository, worktree string) error {
	if repository == "" || worktree == "" {
		return fmt.Errorf("repository and worktree are required")
	}
	command := exec.CommandContext(ctx, m.binary(), "-C", repository, "worktree", "remove", "--force", filepath.Clean(worktree))
	if output, err := command.CombinedOutput(); err != nil {
		return fmt.Errorf("git worktree remove: %w: %s", err, output)
	}
	return nil
}

func (m WorktreeManager) binary() string {
	if m.GitBinary != "" {
		return m.GitBinary
	}
	return "git"
}
