package execution

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/agentic-app-studio/studio/internal/agentadapter"
	"github.com/agentic-app-studio/studio/internal/productgraph"
)

type RunResult struct {
	CardID       string               `json:"cardId"`
	Worktree     string               `json:"worktree"`
	State        string               `json:"state"`
	Events       []agentadapter.Event `json:"events"`
	ChangedFiles []string             `json:"changedFiles"`
}

func ExecuteCard(ctx context.Context, card productgraph.CardSpec, repository, worktreeRoot, baseRef string, adapter agentadapter.Adapter) (RunResult, error) {
	if card.ID == "" || repository == "" || worktreeRoot == "" || baseRef == "" || adapter == nil {
		return RunResult{}, fmt.Errorf("card, repository, worktree root, base ref and adapter are required")
	}
	if err := requireCapabilities(card.RequiredAdapterCapability, adapter.Manifest().Capabilities); err != nil {
		return RunResult{}, err
	}
	if err := os.MkdirAll(worktreeRoot, 0o700); err != nil {
		return RunResult{}, err
	}
	worktree := filepath.Join(worktreeRoot, safeName(card.ID))
	branch := "studio/" + safeName(card.ID)
	manager := WorktreeManager{}
	if err := manager.Add(ctx, repository, worktree, branch, baseRef); err != nil {
		return RunResult{}, err
	}
	result := RunResult{CardID: card.ID, Worktree: worktree, State: "in_progress"}
	session, err := adapter.Start(ctx, SessionRequestForCard(card, worktree))
	if err != nil {
		_ = manager.Remove(context.Background(), repository, worktree)
		return RunResult{}, err
	}
	for event := range session.Events() {
		result.Events = append(result.Events, event)
	}
	if err := session.Wait(); err != nil {
		result.State = "failed"
		return result, err
	}
	changed, err := changedFiles(ctx, worktree, baseRef)
	if err != nil {
		result.State = "failed"
		return result, err
	}
	if err := ValidateChangedFiles(changed, card.ScopePaths, card.ForbiddenPaths); err != nil {
		result.State = "failed"
		return result, err
	}
	result.ChangedFiles = changed
	result.State = "verifying"
	return result, nil
}

func SessionRequestForCard(card productgraph.CardSpec, worktree string) agentadapter.SessionRequest {
	return agentadapter.SessionRequest{WorkingDirectory: worktree, Prompt: fmt.Sprintf(
		"Implement card %s. Goal: %s. Acceptance criteria: %s. Verification commands: %s. Do not modify files outside the allowed scope.",
		card.ID, card.Goal, strings.Join(card.AcceptanceCriteria, "; "), strings.Join(card.VerificationCommands, "; "))}
}

func requireCapabilities(required, available []string) error {
	set := make(map[string]struct{}, len(available))
	for _, capability := range available {
		set[capability] = struct{}{}
	}
	for _, capability := range required {
		if _, ok := set[capability]; !ok {
			return fmt.Errorf("adapter missing capability: %s", capability)
		}
	}
	return nil
}

func changedFiles(ctx context.Context, worktree, baseRef string) ([]string, error) {
	command := exec.CommandContext(ctx, "git", "-C", worktree, "diff", "--name-only", baseRef)
	output, err := command.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}
	return lines, nil
}

func safeName(value string) string {
	value = strings.NewReplacer("/", "-", "\\", "-", "..", "-").Replace(value)
	return strings.Trim(value, ".-")
}
