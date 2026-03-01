package artifact

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEmitBasic(t *testing.T) {
	dir := t.TempDir()

	opts := EmitOptions{
		SkillName: "eda",
		RunID:     "run-001",
		OutDir:    dir,
	}

	path, err := Emit(opts)
	if err != nil {
		t.Fatalf("Emit: %v", err)
	}

	if path != filepath.Join(dir, "eda.json") {
		t.Errorf("unexpected path: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading artifact: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("parsing artifact JSON: %v", err)
	}

	if result["schema_version"] != SchemaVersion {
		t.Errorf("schema_version: want %q, got %q", SchemaVersion, result["schema_version"])
	}
	if result["skill_name"] != "eda" {
		t.Errorf("skill_name: want %q, got %v", "eda", result["skill_name"])
	}
	if result["run_id"] != "run-001" {
		t.Errorf("run_id: want %q, got %v", "run-001", result["run_id"])
	}
}

func TestEmitWithExtraJSON(t *testing.T) {
	dir := t.TempDir()

	opts := EmitOptions{
		SkillName: "model-eval",
		RunID:     "run-002",
		OutDir:    dir,
		ExtraJSON: []byte(`{"decision": "GO", "confidence": "high", "summary": "All good", "findings": []}`),
	}

	path, err := Emit(opts)
	if err != nil {
		t.Fatalf("Emit: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading artifact: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("parsing artifact JSON: %v", err)
	}

	if result["decision"] != "GO" {
		t.Errorf("decision: want %q, got %v", "GO", result["decision"])
	}
	if result["confidence"] != "high" {
		t.Errorf("confidence: want %q, got %v", "high", result["confidence"])
	}
	if result["summary"] != "All good" {
		t.Errorf("summary: want %q, got %v", "All good", result["summary"])
	}
	// schema_version should still be set
	if result["schema_version"] != SchemaVersion {
		t.Errorf("schema_version: want %q, got %q", SchemaVersion, result["schema_version"])
	}
}

func TestEmitRequiresSkillName(t *testing.T) {
	dir := t.TempDir()
	_, err := Emit(EmitOptions{OutDir: dir})
	if err == nil {
		t.Error("expected error for empty skill name")
	}
}

func TestEmitRequiresOutDir(t *testing.T) {
	_, err := Emit(EmitOptions{SkillName: "test"})
	if err == nil {
		t.Error("expected error for empty out dir")
	}
}

func TestEmitCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "nested", "artifacts")

	_, err := Emit(EmitOptions{
		SkillName: "test",
		OutDir:    subdir,
	})
	if err != nil {
		t.Fatalf("Emit: %v", err)
	}

	if _, err := os.Stat(subdir); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}
