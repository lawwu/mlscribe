// Package artifact defines the schema and vocabulary for mlscribe JSON artifacts.
package artifact

// Decision values for gate blocks and artifacts (matches agentic-ml-plugin vocabulary).
const (
	DecisionGO          = "GO"
	DecisionNOGO        = "NO-GO"
	DecisionCONDITIONAL = "CONDITIONAL"
)

// Confidence levels.
const (
	ConfidenceHigh   = "high"
	ConfidenceMedium = "medium"
	ConfidenceLow    = "low"
)

// Gate status values.
const (
	GateStatusPASS    = "PASS"
	GateStatusFAIL    = "FAIL"
	GateStatusSKIPPED = "SKIPPED"
)

// Severity levels for findings.
const (
	SeverityCritical = "critical"
	SeverityHigh     = "high"
	SeverityMedium   = "medium"
	SeverityLow      = "low"
	SeverityInfo     = "info"
)

// IsValidDecision returns true if the string is a valid decision value.
func IsValidDecision(d string) bool {
	switch d {
	case DecisionGO, DecisionNOGO, DecisionCONDITIONAL:
		return true
	}
	return false
}

// IsValidConfidence returns true if the string is a valid confidence level.
func IsValidConfidence(c string) bool {
	switch c {
	case ConfidenceHigh, ConfidenceMedium, ConfidenceLow:
		return true
	}
	return false
}

// IsValidGateStatus returns true if the string is a valid gate status.
func IsValidGateStatus(s string) bool {
	switch s {
	case GateStatusPASS, GateStatusFAIL, GateStatusSKIPPED:
		return true
	}
	return false
}
