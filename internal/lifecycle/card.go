package lifecycle

import (
	"errors"
	"fmt"
)

type CardState string

const (
	CardDraft         CardState = "draft"
	CardReady         CardState = "ready"
	CardQueued        CardState = "queued"
	CardInProgress    CardState = "in_progress"
	CardVerifying     CardState = "verifying"
	CardApproved      CardState = "approved"
	CardMerged        CardState = "merged"
	CardBlocked       CardState = "blocked"
	CardStale         CardState = "stale"
	CardRepair        CardState = "repair"
	CardFailed        CardState = "failed"
	CardCancelled     CardState = "cancelled"
	CardNeedsDecision CardState = "needs_decision"
)

type Card struct {
	ID                        string
	Goal                      string
	ProductGraphRevision      int
	ScopePaths                []string
	ForbiddenPaths            []string
	Dependencies              []string
	SemanticLocks             []string
	SharedResourceLocks       []string
	AcceptanceCriteria        []string
	VerificationCommands      []string
	RequiredAdapterCapability []string
	State                     CardState
}

var ErrInvalidTransition = errors.New("invalid card transition")

var allowedTransitions = map[CardState][]CardState{
	CardDraft:         {CardReady, CardCancelled},
	CardReady:         {CardQueued, CardBlocked, CardCancelled},
	CardQueued:        {CardInProgress, CardBlocked, CardCancelled},
	CardInProgress:    {CardVerifying, CardFailed, CardRepair, CardCancelled},
	CardVerifying:     {CardApproved, CardRepair, CardFailed},
	CardRepair:        {CardInProgress, CardFailed, CardCancelled},
	CardApproved:      {CardMerged, CardStale},
	CardMerged:        {CardStale},
	CardBlocked:       {CardReady, CardCancelled},
	CardStale:         {CardReady, CardCancelled},
	CardFailed:        {CardRepair, CardCancelled},
	CardNeedsDecision: {CardReady, CardCancelled},
}

func (c *Card) Transition(next CardState) error {
	for _, candidate := range allowedTransitions[c.State] {
		if candidate == next {
			c.State = next
			return nil
		}
	}
	return fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, c.State, next)
}
