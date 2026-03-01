package markdown

import (
	"strings"
	"testing"
)

func TestTitleBlockRender(t *testing.T) {
	b := &TitleBlock{
		Title:     "My Experiment",
		Timestamp: "2025-01-01T00:00:00Z",
		Version:   "1.0",
		DocID:     "abc-123",
	}
	got := b.Render()
	if !strings.Contains(got, "# My Experiment") {
		t.Errorf("missing title: %s", got)
	}
	if !strings.Contains(got, "timestamp: 2025-01-01T00:00:00Z") {
		t.Errorf("missing timestamp: %s", got)
	}
}

func TestCommentaryBlockRender(t *testing.T) {
	b := &CommentaryBlock{Content: "This is a note."}
	got := b.Render()
	if !strings.Contains(got, "This is a note.") {
		t.Errorf("unexpected render: %s", got)
	}
}

func TestCodeBlockRender(t *testing.T) {
	b := &CodeBlock{Language: "python", Code: "print('hello')"}
	got := b.Render()
	if !strings.HasPrefix(got, "```python\n") {
		t.Errorf("unexpected render: %s", got)
	}
	if !strings.Contains(got, "print('hello')") {
		t.Errorf("missing code: %s", got)
	}
}

func TestOutputBlockRender(t *testing.T) {
	b := &OutputBlock{Output: "hello"}
	got := b.Render()
	if !strings.Contains(got, "hello") {
		t.Errorf("missing output: %s", got)
	}
}

func TestMetricsBlockRender(t *testing.T) {
	b := &MetricsBlock{
		Metrics: []MetricEntry{
			{Name: "accuracy", Value: "0.94"},
			{Name: "loss", Value: "0.12"},
		},
	}
	got := b.Render()
	if !strings.Contains(got, "<!-- mlscribe:metrics -->") {
		t.Errorf("missing sentinel: %s", got)
	}
	if !strings.Contains(got, "accuracy") {
		t.Errorf("missing metric name: %s", got)
	}
	if !strings.Contains(got, "0.94") {
		t.Errorf("missing metric value: %s", got)
	}
}

func TestMetricsBlockWithBaselineRender(t *testing.T) {
	b := &MetricsBlock{
		Metrics: []MetricEntry{
			{Name: "accuracy", Value: "0.94", Baseline: "0.90", Delta: "+0.04"},
		},
	}
	got := b.Render()
	if !strings.Contains(got, "Baseline") {
		t.Errorf("missing Baseline column: %s", got)
	}
	if !strings.Contains(got, "+0.04") {
		t.Errorf("missing delta: %s", got)
	}
}

func TestGateBlockRender(t *testing.T) {
	b := &GateBlock{
		GateName:   "data-quality",
		Status:     "PASS",
		Decision:   "GO",
		Confidence: "high",
		Findings:   "All checks passed",
	}
	got := b.Render()
	if !strings.Contains(got, "<!-- mlscribe:gate -->") {
		t.Errorf("missing sentinel: %s", got)
	}
	if !strings.Contains(got, "data-quality") {
		t.Errorf("missing gate name: %s", got)
	}
	if !strings.Contains(got, "PASS") {
		t.Errorf("missing status: %s", got)
	}
}

func TestArtifactRefBlockRender(t *testing.T) {
	b := &ArtifactRefBlock{
		SkillName:    "eda",
		ArtifactPath: "/tmp/artifacts/eda.json",
		RunID:        "run-001",
	}
	got := b.Render()
	if !strings.Contains(got, "<!-- mlscribe:artifact_ref -->") {
		t.Errorf("missing sentinel: %s", got)
	}
	if !strings.Contains(got, "eda") {
		t.Errorf("missing skill name: %s", got)
	}
	if !strings.Contains(got, "run-001") {
		t.Errorf("missing run ID: %s", got)
	}
}
