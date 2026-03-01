// Package markdown defines block types for mlscribe documents.
package markdown

import (
	"fmt"
	"strings"
)

// Block is the interface implemented by all document block types.
type Block interface {
	// Render returns the markdown representation of the block.
	Render() string
	// Type returns a string identifier for the block type.
	Type() string
}

// --- Inherited block types ---

// TitleBlock represents the document title and metadata header.
type TitleBlock struct {
	Title     string
	Timestamp string // RFC3339
	Version   string
	DocID     string // UUID
}

func (b *TitleBlock) Type() string { return "title" }
func (b *TitleBlock) Render() string {
	var sb strings.Builder
	sb.WriteString("# " + b.Title + "\n\n")
	sb.WriteString("<!-- mlscribe:meta\n")
	sb.WriteString("timestamp: " + b.Timestamp + "\n")
	sb.WriteString("version: " + b.Version + "\n")
	sb.WriteString("doc_id: " + b.DocID + "\n")
	sb.WriteString("-->\n")
	return sb.String()
}

// CommentaryBlock is free-form markdown prose.
type CommentaryBlock struct {
	Content string
}

func (b *CommentaryBlock) Type() string { return "commentary" }
func (b *CommentaryBlock) Render() string {
	return b.Content + "\n"
}

// CodeBlock is a fenced code block with a language tag.
type CodeBlock struct {
	Language string
	Code     string
}

func (b *CodeBlock) Type() string { return "code" }
func (b *CodeBlock) Render() string {
	return fmt.Sprintf("```%s\n%s\n```\n", b.Language, b.Code)
}

// OutputBlock is captured stdout/stderr from an exec command.
type OutputBlock struct {
	Output string
}

func (b *OutputBlock) Type() string { return "output" }
func (b *OutputBlock) Render() string {
	if b.Output == "" {
		return "```\n\n```\n"
	}
	return fmt.Sprintf("```\n%s\n```\n", b.Output)
}

// ImageOutputBlock is an image reference embedded in the document.
type ImageOutputBlock struct {
	AltText  string
	Filename string
}

func (b *ImageOutputBlock) Type() string { return "image" }
func (b *ImageOutputBlock) Render() string {
	return fmt.Sprintf("![%s](%s)\n", b.AltText, b.Filename)
}

// --- New ML block types ---

// MetricEntry is a single row in a metrics table.
type MetricEntry struct {
	Name     string
	Value    string
	Baseline string // optional
	Delta    string // optional, computed
}

// MetricsBlock is a structured key-value metrics table.
type MetricsBlock struct {
	Metrics []MetricEntry
}

func (b *MetricsBlock) Type() string { return "metrics" }
func (b *MetricsBlock) Render() string {
	var sb strings.Builder
	sb.WriteString("<!-- mlscribe:metrics -->\n")

	hasBaseline := false
	for _, m := range b.Metrics {
		if m.Baseline != "" {
			hasBaseline = true
			break
		}
	}

	if hasBaseline {
		sb.WriteString("| Metric | Value | Baseline | Delta |\n")
		sb.WriteString("|--------|-------|----------|-------|\n")
		for _, m := range b.Metrics {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", m.Name, m.Value, m.Baseline, m.Delta))
		}
	} else {
		sb.WriteString("| Metric | Value |\n")
		sb.WriteString("|--------|-------|\n")
		for _, m := range b.Metrics {
			sb.WriteString(fmt.Sprintf("| %s | %s |\n", m.Name, m.Value))
		}
	}
	sb.WriteString("<!-- /mlscribe:metrics -->\n")
	return sb.String()
}

// GateBlock records a lifecycle gate decision.
type GateBlock struct {
	GateName   string
	Status     string // PASS | FAIL | SKIPPED
	Decision   string // GO | NO-GO | CONDITIONAL
	Confidence string // high | medium | low
	Findings   string // JSON string or plain text
}

func (b *GateBlock) Type() string { return "gate" }
func (b *GateBlock) Render() string {
	var sb strings.Builder
	sb.WriteString("<!-- mlscribe:gate -->\n")
	sb.WriteString(fmt.Sprintf("**Gate**: %s\n\n", b.GateName))
	sb.WriteString(fmt.Sprintf("| Field | Value |\n"))
	sb.WriteString(fmt.Sprintf("|-------|-------|\n"))
	sb.WriteString(fmt.Sprintf("| Status | %s |\n", b.Status))
	if b.Decision != "" {
		sb.WriteString(fmt.Sprintf("| Decision | %s |\n", b.Decision))
	}
	if b.Confidence != "" {
		sb.WriteString(fmt.Sprintf("| Confidence | %s |\n", b.Confidence))
	}
	if b.Findings != "" {
		sb.WriteString(fmt.Sprintf("\n**Findings**: %s\n", b.Findings))
	}
	sb.WriteString("<!-- /mlscribe:gate -->\n")
	return sb.String()
}

// ArtifactRefBlock references an emitted JSON artifact file.
type ArtifactRefBlock struct {
	SkillName    string
	ArtifactPath string
	RunID        string
}

func (b *ArtifactRefBlock) Type() string { return "artifact_ref" }
func (b *ArtifactRefBlock) Render() string {
	var sb strings.Builder
	sb.WriteString("<!-- mlscribe:artifact_ref -->\n")
	sb.WriteString(fmt.Sprintf("**Artifact**: [%s](%s)", b.SkillName, b.ArtifactPath))
	if b.RunID != "" {
		sb.WriteString(fmt.Sprintf(" (run: %s)", b.RunID))
	}
	sb.WriteString("\n<!-- /mlscribe:artifact_ref -->\n")
	return sb.String()
}
