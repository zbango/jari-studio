package productgraph

import "time"

type Status string

const (
	StatusObserved   Status = "observed"
	StatusInferred   Status = "inferred"
	StatusProposed   Status = "proposed"
	StatusConfirmed  Status = "confirmed"
	StatusRejected   Status = "rejected"
	StatusSuperseded Status = "superseded"
)

type SourceRef struct {
	Kind string `json:"kind"`
	Ref  string `json:"ref"`
}

type Provenance struct {
	DerivedFrom  []string  `json:"derivedFrom,omitempty"`
	GeneratedBy  []string  `json:"generatedBy,omitempty"`
	AttributedTo []string  `json:"attributedTo,omitempty"`
	Confidence   float64   `json:"confidence"`
	GeneratedAt  time.Time `json:"generatedAt"`
}

type NodeMetadata struct {
	Sources    []SourceRef `json:"source,omitempty"`
	Status     Status      `json:"status"`
	Impact     string      `json:"impact,omitempty"`
	Provenance Provenance  `json:"provenance,omitempty"`
	Revision   int         `json:"revision"`
}

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Brief     string    `json:"brief,omitempty"`
	Revision  int       `json:"revision"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Module struct {
	ID          string       `json:"id"`
	ProjectID   string       `json:"projectId"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Metadata    NodeMetadata `json:"metadata"`
}

type Entity struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Fields   []Field      `json:"fields,omitempty"`
	Metadata NodeMetadata `json:"metadata"`
}

type Field struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Type     string       `json:"type"`
	Required bool         `json:"required,omitempty"`
	Metadata NodeMetadata `json:"metadata"`
}

type Relation struct {
	ID         string       `json:"id"`
	FromEntity string       `json:"fromEntity"`
	ToEntity   string       `json:"toEntity"`
	Kind       string       `json:"kind"`
	Metadata   NodeMetadata `json:"metadata"`
}

type Validation struct {
	ID       string       `json:"id"`
	TargetID string       `json:"targetId"`
	Rule     string       `json:"rule"`
	Message  string       `json:"message,omitempty"`
	Metadata NodeMetadata `json:"metadata"`
}

type Gap struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Severity string       `json:"severity"`
	Impact   string       `json:"impact"`
	Status   Status       `json:"status"`
	Metadata NodeMetadata `json:"metadata"`
}

type Requirement struct {
	ID        string       `json:"id"`
	Statement string       `json:"statement"`
	Metadata  NodeMetadata `json:"metadata"`
}

type AcceptanceCriterion struct {
	ID            string       `json:"id"`
	RequirementID string       `json:"requirementId"`
	Statement     string       `json:"statement"`
	Mandatory     bool         `json:"mandatory"`
	Verification  string       `json:"verificationType"`
	Metadata      NodeMetadata `json:"metadata"`
}

type CardSpec struct {
	ID                        string   `json:"id"`
	Goal                      string   `json:"goal"`
	ScopePaths                []string `json:"scopePaths"`
	ForbiddenPaths            []string `json:"forbiddenPaths"`
	Dependencies              []string `json:"dependencies,omitempty"`
	SemanticLocks              []string `json:"semanticLocks,omitempty"`
	SharedResourceLocks        []string `json:"sharedResourceLocks,omitempty"`
	AcceptanceCriteria        []string `json:"acceptanceCriteria"`
	VerificationCommands      []string `json:"verificationCommands"`
	RequiredAdapterCapability []string `json:"requiredAdapterCapabilities"`
	State                     string   `json:"state"`
}

type AdapterManifest struct {
	ID           string   `json:"id"`
	Version      string   `json:"version"`
	Transport    string   `json:"transport"`
	Capabilities []string `json:"capabilities"`
}

type CardScheduleDecision struct {
	CardID  string   `json:"cardId"`
	State   string   `json:"state"`
	Reasons []string `json:"reasons,omitempty"`
}

type ScheduleResult struct {
	Adapter   AdapterManifest        `json:"adapter"`
	Decisions []CardScheduleDecision `json:"decisions"`
}

type CompiledPlan struct {
	RevisionID         string                `json:"revisionId"`
	Requirements       []Requirement         `json:"requirements"`
	AcceptanceCriteria []AcceptanceCriterion `json:"acceptanceCriteria"`
	Cards              []CardSpec            `json:"cards"`
	BlockingGaps       []Gap                 `json:"blockingGaps,omitempty"`
}

type Graph struct {
	Project     Project      `json:"project"`
	Modules     []Module     `json:"modules,omitempty"`
	Entities    []Entity     `json:"entities,omitempty"`
	Relations   []Relation   `json:"relations,omitempty"`
	Validations []Validation `json:"validations,omitempty"`
	Gaps        []Gap        `json:"gaps,omitempty"`
}

type Revision struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	Number    int       `json:"number"`
	Status    Status    `json:"status"`
	Snapshot  string    `json:"snapshot"`
	CreatedAt time.Time `json:"createdAt"`
}
