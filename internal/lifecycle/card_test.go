package lifecycle

import "testing"

func TestCardTransitionRejectsSkippingVerification(t *testing.T) {
	card := Card{ID: "card-1", State: CardDraft}
	if err := card.Transition(CardMerged); err == nil {
		t.Fatal("expected invalid transition")
	}
}

func TestCardTransitionAllowsExecutionPath(t *testing.T) {
	card := Card{ID: "card-1", State: CardDraft}
	for _, state := range []CardState{CardReady, CardQueued, CardInProgress, CardVerifying, CardApproved, CardMerged} {
		if err := card.Transition(state); err != nil {
			t.Fatalf("transition to %s: %v", state, err)
		}
	}
}
