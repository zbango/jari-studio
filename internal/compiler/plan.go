package compiler

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/agentic-app-studio/studio/internal/productgraph"
)

var (
	ErrInvalidGraph      = errors.New("revision snapshot is not a valid product graph")
	ErrBlockingGaps      = errors.New("revision contains unresolved blocking gaps")
	ErrRevisionProjectID = errors.New("revision snapshot project does not match revision")
)

type snapshot struct {
	Project     productgraph.Project      `json:"project"`
	Modules     []productgraph.Module     `json:"modules"`
	Entities    []productgraph.Entity     `json:"entities"`
	Relations   []productgraph.Relation   `json:"relations"`
	Validations []productgraph.Validation `json:"validations"`
	Gaps        []productgraph.Gap        `json:"gaps"`
}

func Compile(revision productgraph.Revision) (productgraph.CompiledPlan, error) {
	var graph snapshot
	if err := json.Unmarshal([]byte(revision.Snapshot), &graph); err != nil || graph.Project.ID == "" {
		return productgraph.CompiledPlan{}, ErrInvalidGraph
	}
	if graph.Project.ID != revision.ProjectID {
		return productgraph.CompiledPlan{}, ErrRevisionProjectID
	}
	blocking := make([]productgraph.Gap, 0)
	for _, gap := range graph.Gaps {
		if strings.EqualFold(gap.Severity, "blocking") && gap.Status != productgraph.StatusRejected && gap.Status != productgraph.StatusSuperseded {
			blocking = append(blocking, gap)
		}
	}
	if len(blocking) > 0 {
		return productgraph.CompiledPlan{RevisionID: revision.ID, BlockingGaps: blocking}, fmt.Errorf("%w: %d", ErrBlockingGaps, len(blocking))
	}

	plan := productgraph.CompiledPlan{RevisionID: revision.ID, BlockingGaps: blocking}
	foundationID := cardID(revision.ID, "foundation")
	plan.Requirements = append(plan.Requirements, productgraph.Requirement{
		ID:        requirementID(revision.ID, "foundation"),
		Statement: "Generated application bundle has a runnable shared foundation",
		Metadata:  proposedMetadata(revision.Number),
	})
	plan.AcceptanceCriteria = append(plan.AcceptanceCriteria, productgraph.AcceptanceCriterion{
		ID:            acceptanceID(revision.ID, "foundation"),
		RequirementID: plan.Requirements[0].ID,
		Statement:     "Web, desktop and mobile bundles share the approved project contracts",
		Mandatory:     true,
		Verification:  "integration",
		Metadata:      proposedMetadata(revision.Number),
	})
	plan.Cards = append(plan.Cards, productgraph.CardSpec{
		ID:                        foundationID,
		Goal:                      "Prepare the generated application foundation",
		ScopePaths:                []string{"packages/shared/**", "services/api/**"},
		ForbiddenPaths:            []string{".env*", "secrets/**"},
		AcceptanceCriteria:        []string{acceptanceID(revision.ID, "foundation")},
		VerificationCommands:      []string{"go test ./...", "npm run check"},
		RequiredAdapterCapability: []string{"structured_events", "working_directory", "diff_events"},
		State:                     "ready",
	})

	for _, module := range graph.Modules {
		slug := slugify(module.Name)
		reqID := requirementID(revision.ID, "module-"+slug)
		acID := acceptanceID(revision.ID, "module-"+slug)
		plan.Requirements = append(plan.Requirements, productgraph.Requirement{ID: reqID, Statement: fmt.Sprintf("User can access the %s module", module.Name), Metadata: proposedMetadata(revision.Number)})
		plan.AcceptanceCriteria = append(plan.AcceptanceCriteria, productgraph.AcceptanceCriterion{ID: acID, RequirementID: reqID, Statement: fmt.Sprintf("The %s module renders its approved navigation entry and initial state", module.Name), Mandatory: true, Verification: "e2e", Metadata: proposedMetadata(revision.Number)})
		plan.Cards = append(plan.Cards, productgraph.CardSpec{
			ID:                        cardID(revision.ID, "module-"+slug),
			Goal:                      fmt.Sprintf("Build the %s module across application surfaces", module.Name),
			ScopePaths:                []string{"apps/desktop/**", "apps/web/**", "apps/mobile/**", "packages/shared/**"},
			ForbiddenPaths:            []string{"auth/**", "infra/secrets/**"},
			Dependencies:              []string{foundationID},
			AcceptanceCriteria:        []string{acID},
			VerificationCommands:      []string{"npm run check", "npm run build"},
			RequiredAdapterCapability: []string{"structured_events", "working_directory", "changed_files"},
			State:                     "blocked",
		})
	}

	for _, entity := range graph.Entities {
		slug := slugify(entity.Name)
		reqID := requirementID(revision.ID, "entity-"+slug)
		acID := acceptanceID(revision.ID, "entity-"+slug)
		plan.Requirements = append(plan.Requirements, productgraph.Requirement{ID: reqID, Statement: fmt.Sprintf("System stores %s according to the approved domain proposal", entity.Name), Metadata: proposedMetadata(revision.Number)})
		plan.AcceptanceCriteria = append(plan.AcceptanceCriteria, productgraph.AcceptanceCriterion{ID: acID, RequirementID: reqID, Statement: fmt.Sprintf("The %s entity validates and persists its proposed fields", entity.Name), Mandatory: true, Verification: "integration", Metadata: proposedMetadata(revision.Number)})
		plan.Cards = append(plan.Cards, productgraph.CardSpec{
			ID:                        cardID(revision.ID, "entity-"+slug),
			Goal:                      fmt.Sprintf("Implement the %s domain contract", entity.Name),
			ScopePaths:                []string{"services/api/**", "apps/desktop/**", "packages/shared/**"},
			ForbiddenPaths:            []string{"auth/**", "infra/secrets/**"},
			Dependencies:              []string{foundationID},
			AcceptanceCriteria:        []string{acID},
			VerificationCommands:      []string{"go test ./...", "npm run check"},
			RequiredAdapterCapability: []string{"structured_events", "working_directory", "changed_files"},
			State:                     "blocked",
		})
	}

	if len(plan.Cards) > 1 {
		verifyID := cardID(revision.ID, "verification")
		dependencies := make([]string, 0, len(plan.Cards))
		for _, card := range plan.Cards[1:] {
			dependencies = append(dependencies, card.ID)
		}
		plan.Cards = append(plan.Cards, productgraph.CardSpec{ID: verifyID, Goal: "Run integration verification for the compiled bundle", ScopePaths: []string{"tests/**", "packages/testkit/**"}, Dependencies: dependencies, AcceptanceCriteria: []string{acceptanceID(revision.ID, "foundation")}, VerificationCommands: []string{"npm run build", "go test ./..."}, RequiredAdapterCapability: []string{"structured_events", "diff_events"}, State: "blocked"})
	}
	return plan, nil
}

func proposedMetadata(revision int) productgraph.NodeMetadata {
	return productgraph.NodeMetadata{Status: productgraph.StatusProposed, Revision: revision}
}

var nonWord = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func slugify(value string) string {
	clean := strings.Trim(nonWord.ReplaceAllString(strings.ToLower(value), "-"), "-")
	if clean == "" {
		return "unnamed"
	}
	return clean
}

func requirementID(revisionID, name string) string {
	return "REQ." + strings.ToUpper(slugify(revisionID+"-"+name))
}
func acceptanceID(revisionID, name string) string {
	return "AC." + strings.ToUpper(slugify(revisionID+"-"+name))
}
func cardID(revisionID, name string) string {
	return "CARD." + strings.ToUpper(slugify(revisionID+"-"+name))
}
