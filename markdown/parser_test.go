package markdown

import (
	"testing"
)

func TestRoundTrip(t *testing.T) {
	blocks := []Block{
		&TitleBlock{
			Title:     "Test Document",
			Timestamp: "2025-01-01T00:00:00Z",
			Version:   "1.0",
			DocID:     "abc-123",
		},
		&CommentaryBlock{Content: "This is a note."},
		&CodeBlock{Language: "python", Code: "print('hello')"},
		&OutputBlock{Output: "hello"},
		&MetricsBlock{
			Metrics: []MetricEntry{
				{Name: "accuracy", Value: "0.94"},
				{Name: "f1", Value: "0.91"},
			},
		},
		&GateBlock{
			GateName:   "model-quality",
			Status:     "PASS",
			Decision:   "GO",
			Confidence: "high",
			Findings:   "Metrics above threshold",
		},
		&ArtifactRefBlock{
			SkillName:    "eda",
			ArtifactPath: "/tmp/eda.json",
			RunID:        "run-001",
		},
	}

	rendered := Render(blocks)
	parsed, err := Parse(rendered)
	if err != nil {
		t.Fatalf("parsing rendered document: %v", err)
	}

	if len(parsed) != len(blocks) {
		t.Fatalf("expected %d blocks, got %d\nparsed:\n%s", len(blocks), len(parsed), rendered)
	}

	for i, expected := range blocks {
		if parsed[i].Type() != expected.Type() {
			t.Errorf("block %d: expected type %q, got %q", i, expected.Type(), parsed[i].Type())
		}
	}
}

func TestParseTitleBlock(t *testing.T) {
	input := `# My Experiment

<!-- mlscribe:meta
timestamp: 2025-01-01T00:00:00Z
version: 1.0
doc_id: abc-123
-->
`
	blocks, err := Parse(input)
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}
	if len(blocks) == 0 {
		t.Fatal("no blocks parsed")
	}

	tb, ok := blocks[0].(*TitleBlock)
	if !ok {
		t.Fatalf("expected TitleBlock, got %T", blocks[0])
	}

	if tb.Title != "My Experiment" {
		t.Errorf("title: want %q, got %q", "My Experiment", tb.Title)
	}
	if tb.Timestamp != "2025-01-01T00:00:00Z" {
		t.Errorf("timestamp: want %q, got %q", "2025-01-01T00:00:00Z", tb.Timestamp)
	}
	if tb.DocID != "abc-123" {
		t.Errorf("doc_id: want %q, got %q", "abc-123", tb.DocID)
	}
}

func TestParseMetricsBlock(t *testing.T) {
	input := `<!-- mlscribe:metrics -->
| Metric | Value |
|--------|-------|
| accuracy | 0.94 |
| loss | 0.12 |
<!-- /mlscribe:metrics -->
`
	blocks, err := Parse(input)
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}
	if len(blocks) == 0 {
		t.Fatal("no blocks parsed")
	}

	mb, ok := blocks[0].(*MetricsBlock)
	if !ok {
		t.Fatalf("expected MetricsBlock, got %T", blocks[0])
	}

	if len(mb.Metrics) != 2 {
		t.Errorf("expected 2 metrics, got %d", len(mb.Metrics))
	}
}

func TestParseGateBlock(t *testing.T) {
	input := `<!-- mlscribe:gate -->
**Gate**: model-quality

| Field | Value |
|-------|-------|
| Status | PASS |
| Decision | GO |
| Confidence | high |

**Findings**: All good
<!-- /mlscribe:gate -->
`
	blocks, err := Parse(input)
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}
	if len(blocks) == 0 {
		t.Fatal("no blocks parsed")
	}

	gb, ok := blocks[0].(*GateBlock)
	if !ok {
		t.Fatalf("expected GateBlock, got %T", blocks[0])
	}

	if gb.GateName != "model-quality" {
		t.Errorf("gate name: want %q, got %q", "model-quality", gb.GateName)
	}
	if gb.Status != "PASS" {
		t.Errorf("status: want %q, got %q", "PASS", gb.Status)
	}
	if gb.Decision != "GO" {
		t.Errorf("decision: want %q, got %q", "GO", gb.Decision)
	}
}

func TestParseCodeAndOutputBlocks(t *testing.T) {
	input := "```python\nprint('hello')\n```\n\n```\nhello\n```\n"

	blocks, err := Parse(input)
	if err != nil {
		t.Fatalf("parsing: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}

	cb, ok := blocks[0].(*CodeBlock)
	if !ok {
		t.Fatalf("expected CodeBlock, got %T", blocks[0])
	}
	if cb.Language != "python" {
		t.Errorf("language: want %q, got %q", "python", cb.Language)
	}

	ob, ok := blocks[1].(*OutputBlock)
	if !ok {
		t.Fatalf("expected OutputBlock, got %T", blocks[1])
	}
	if ob.Output != "hello" {
		t.Errorf("output: want %q, got %q", "hello", ob.Output)
	}
}
