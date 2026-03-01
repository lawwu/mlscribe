// Package artifact defines Go structs matching the agentic-ml-plugin base artifact schema.
package artifact

// SchemaVersion is the current schema version.
const SchemaVersion = "1.0"

// BaseArtifact matches the base schema from the agentic-ml-plugin (schemas.md).
// All skill-specific artifacts embed or extend this structure.
type BaseArtifact struct {
	SchemaVersion   string    `json:"schema_version"`
	SkillName       string    `json:"skill_name"`
	RunID           string    `json:"run_id"`
	Timestamp       string    `json:"timestamp"`
	DurationSeconds float64   `json:"duration_seconds,omitempty"`
	Decision        string    `json:"decision"`
	Confidence      string    `json:"confidence"`
	Summary         string    `json:"summary"`
	Findings        []Finding `json:"findings"`
	NextCommands    []string  `json:"next_commands,omitempty"`
}

// Finding is a single finding entry in an artifact.
type Finding struct {
	Severity    string `json:"severity"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description"`
	Remediation string `json:"remediation,omitempty"`
}

// MetricsArtifact is a skill artifact containing metrics data.
type MetricsArtifact struct {
	BaseArtifact
	Metrics map[string]interface{} `json:"metrics,omitempty"`
}
