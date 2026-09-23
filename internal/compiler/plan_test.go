package compiler

import (
	"testing"
	"time"

	"github.com/agentic-app-studio/studio/internal/productgraph"
)

func TestCompileRevisionProducesDeterministicPlan(t *testing.T) {
	revision := productgraph.Revision{
		ID: "project.demo.revision.1", ProjectID: "project.demo", Number: 1,
		Snapshot:  `{"project":{"id":"project.demo"},"modules":[{"name":"Members"}],"entities":[{"name":"Member"}]}`,
		CreatedAt: time.Now().UTC(),
	}
	first, err := Compile(revision)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	second, err := Compile(revision)
	if err != nil {
		t.Fatalf("compile second: %v", err)
	}
	if len(first.Requirements) != 3 || len(first.AcceptanceCriteria) != 3 || len(first.Cards) != 4 {
		t.Fatalf("unexpected plan sizes: requirements=%d ac=%d cards=%d", len(first.Requirements), len(first.AcceptanceCriteria), len(first.Cards))
	}
	if first.Cards[1].ID != second.Cards[1].ID || first.Cards[1].Goal != second.Cards[1].Goal {
		t.Fatalf("plan is not deterministic: %+v vs %+v", first.Cards[1], second.Cards[1])
	}
}

func TestCompileRevisionBlocksUnresolvedBlockingGap(t *testing.T) {
	revision := productgraph.Revision{
		ID: "project.demo.revision.2", ProjectID: "project.demo", Number: 2,
		Snapshot: `{"project":{"id":"project.demo"},"gaps":[{"id":"gap-1","severity":"blocking","status":"proposed"}]}`,
	}
	_, err := Compile(revision)
	if err == nil || err.Error() == "" {
		t.Fatal("expected blocking gap error")
	}
}
