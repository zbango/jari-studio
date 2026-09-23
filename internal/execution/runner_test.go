package execution

import (
	"context"
	"testing"

	"github.com/agentic-app-studio/studio/internal/agentadapter"
	"github.com/agentic-app-studio/studio/internal/productgraph"
)

func TestRequireCapabilities(t *testing.T) {
	if err := requireCapabilities([]string{"structured_events"}, []string{"interactive"}); err == nil {
		t.Fatal("expected missing capability")
	}
	if err := requireCapabilities([]string{"structured_events"}, []string{"structured_events", "diff_events"}); err != nil {
		t.Fatalf("expected capabilities to pass: %v", err)
	}
}

func TestSessionRequestForCardContainsExecutionContract(t *testing.T) {
	request := SessionRequestForCard(productgraph.CardSpec{ID: "CARD-1", Goal: "Build module", AcceptanceCriteria: []string{"AC-1"}, VerificationCommands: []string{"go test ./..."}}, "/tmp/worktree")
	if request.WorkingDirectory != "/tmp/worktree" || request.Prompt == "" {
		t.Fatalf("unexpected request: %+v", request)
	}
}

func TestCLIAdapterRejectsMissingRequestFields(t *testing.T) {
	adapter, err := agentadapter.NewCLIAdapter(agentadapter.Manifest{ID: "test"}, "echo", func(prompt string) []string { return []string{prompt} })
	if err != nil {
		t.Fatalf("create adapter: %v", err)
	}
	if _, err := adapter.Start(context.Background(), agentadapter.SessionRequest{}); err != agentadapter.ErrInvalidRequest {
		t.Fatalf("expected invalid request, got %v", err)
	}
}
