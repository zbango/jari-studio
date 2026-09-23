package scheduler

import (
	"fmt"

	"github.com/agentic-app-studio/studio/internal/productgraph"
)

func Evaluate(cards []productgraph.CardSpec, adapter productgraph.AdapterManifest) productgraph.ScheduleResult {
	result := productgraph.ScheduleResult{Adapter: adapter, Decisions: make([]productgraph.CardScheduleDecision, 0, len(cards))}
	known := make(map[string]productgraph.CardSpec, len(cards))
	for _, card := range cards {
		known[card.ID] = card
	}
	for _, card := range cards {
		decision := productgraph.CardScheduleDecision{CardID: card.ID, State: "ready"}
		for _, dependencyID := range card.Dependencies {
			dependency, ok := known[dependencyID]
			if !ok {
				decision.State = "blocked"
				decision.Reasons = append(decision.Reasons, fmt.Sprintf("missing dependency: %s", dependencyID))
				continue
			}
			if dependency.State != "merged" {
				decision.State = "blocked"
				decision.Reasons = append(decision.Reasons, fmt.Sprintf("dependency not merged: %s", dependencyID))
			}
		}
		capabilities := make(map[string]struct{}, len(adapter.Capabilities))
		for _, capability := range adapter.Capabilities {
			capabilities[capability] = struct{}{}
		}
		for _, required := range card.RequiredAdapterCapability {
			if _, ok := capabilities[required]; !ok {
				decision.State = "blocked"
				decision.Reasons = append(decision.Reasons, fmt.Sprintf("adapter missing capability: %s", required))
			}
		}
		if decision.State == "ready" && card.State != "ready" {
			decision.State = card.State
			if card.State == "blocked" {
				decision.Reasons = append(decision.Reasons, "card is blocked by its current lifecycle state")
			}
		}
		result.Decisions = append(result.Decisions, decision)
	}
	return result
}

