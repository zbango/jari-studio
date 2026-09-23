package scheduler

import (
	"testing"

	"github.com/agentic-app-studio/studio/internal/productgraph"
)

func TestEvaluateBlocksDependencyAndMissingCapability(t *testing.T) {
	cards := []productgraph.CardSpec{
		{ID: "CARD.FOUNDATION", State: "ready", RequiredAdapterCapability: []string{"structured_events"}},
		{ID: "CARD.MODULE", State: "blocked", Dependencies: []string{"CARD.FOUNDATION"}, RequiredAdapterCapability: []string{"changed_files"}},
	}
	result := Evaluate(cards, productgraph.AdapterManifest{ID: "generic-pty", Capabilities: []string{"structured_events"}})
	if result.Decisions[0].State != "ready" {
		t.Fatalf("foundation should be ready: %+v", result.Decisions[0])
	}
	if result.Decisions[1].State != "blocked" || len(result.Decisions[1].Reasons) != 2 {
		t.Fatalf("module should be blocked for two reasons: %+v", result.Decisions[1])
	}
}

func TestEvaluateAllowsMergedDependencies(t *testing.T) {
	cards := []productgraph.CardSpec{
		{ID: "CARD.FOUNDATION", State: "merged", RequiredAdapterCapability: []string{"structured_events"}},
		{ID: "CARD.MODULE", State: "ready", Dependencies: []string{"CARD.FOUNDATION"}, RequiredAdapterCapability: []string{"structured_events"}},
	}
	result := Evaluate(cards, productgraph.AdapterManifest{ID: "codex", Capabilities: []string{"structured_events"}})
	if result.Decisions[1].State != "ready" || len(result.Decisions[1].Reasons) != 0 {
		t.Fatalf("module should be ready: %+v", result.Decisions[1])
	}
}

