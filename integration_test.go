package main_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary
	dir, err := os.MkdirTemp("", "mlscribe-test-*")
	if err != nil {
		panic("creating temp dir: " + err.Error())
	}
	defer os.RemoveAll(dir)

	binaryPath = filepath.Join(dir, "mlscribe")
	out, err := exec.Command("go", "build", "-o", binaryPath, ".").CombinedOutput()
	if err != nil {
		panic("building binary: " + string(out))
	}

	os.Exit(m.Run())
}

func run(t *testing.T, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	return string(out), exitCode
}

func TestE2EFullWorkflow(t *testing.T) {
	dir := t.TempDir()
	doc := filepath.Join(dir, "experiment.md")

	// init
	out, code := run(t, "init", doc, "My ML Experiment")
	if code != 0 {
		t.Fatalf("init failed (code %d): %s", code, out)
	}
	if !strings.Contains(out, "Created") {
		t.Errorf("init: expected 'Created' in output: %s", out)
	}

	// note
	out, code = run(t, "note", doc, "Training a simple model")
	if code != 0 {
		t.Fatalf("note failed (code %d): %s", code, out)
	}

	// exec
	out, code = run(t, "exec", doc, "python", "print('hello from python')")
	if code != 0 {
		t.Fatalf("exec failed (code %d): %s", code, out)
	}
	if !strings.Contains(out, "hello from python") {
		t.Errorf("exec: expected output in stdout: %s", out)
	}

	// metrics
	out, code = run(t, "metrics", doc, "--json", `{"accuracy": 0.94, "loss": 0.12}`)
	if code != 0 {
		t.Fatalf("metrics failed (code %d): %s", code, out)
	}
	if !strings.Contains(out, "2 metrics") {
		t.Errorf("metrics: expected '2 metrics' in output: %s", out)
	}

	// gate
	out, code = run(t, "gate", doc, "model-quality", "PASS",
		"--decision", "GO", "--confidence", "high", "--findings", "Accuracy above threshold")
	if code != 0 {
		t.Fatalf("gate failed (code %d): %s", code, out)
	}
	if !strings.Contains(out, "model-quality") {
		t.Errorf("gate: expected gate name in output: %s", out)
	}

	// artifact
	artifactDir := filepath.Join(dir, "artifacts")
	out, code = run(t, "artifact", doc,
		"--skill", "model-eval",
		"--out-dir", artifactDir,
		"--run-id", "run-001",
		"--json", `{"decision": "GO", "confidence": "high", "summary": "Model passed all gates"}`)
	if code != 0 {
		t.Fatalf("artifact failed (code %d): %s", code, out)
	}

	// Verify artifact file exists and is valid JSON
	artifactFile := filepath.Join(artifactDir, "model-eval.json")
	data, err := os.ReadFile(artifactFile)
	if err != nil {
		t.Fatalf("reading artifact file: %v", err)
	}
	var artifactData map[string]interface{}
	if err := json.Unmarshal(data, &artifactData); err != nil {
		t.Fatalf("parsing artifact JSON: %v", err)
	}
	if artifactData["schema_version"] != "1.0" {
		t.Errorf("artifact: schema_version: want '1.0', got %v", artifactData["schema_version"])
	}
	if artifactData["skill_name"] != "model-eval" {
		t.Errorf("artifact: skill_name: want 'model-eval', got %v", artifactData["skill_name"])
	}
	if artifactData["run_id"] != "run-001" {
		t.Errorf("artifact: run_id: want 'run-001', got %v", artifactData["run_id"])
	}

	// verify
	out, code = run(t, "verify", doc)
	if code != 0 {
		t.Fatalf("verify failed (code %d): %s", code, out)
	}
	if !strings.Contains(out, "verified successfully") {
		t.Errorf("verify: expected 'verified successfully': %s", out)
	}

	// pop
	out, code = run(t, "pop", doc)
	if code != 0 {
		t.Fatalf("pop failed (code %d): %s", code, out)
	}
	if !strings.Contains(out, "Removed") {
		t.Errorf("pop: expected 'Removed' in output: %s", out)
	}

	// extract
	out, code = run(t, "extract", doc)
	if code != 0 {
		t.Fatalf("extract failed (code %d): %s", code, out)
	}
	if !strings.Contains(out, "mlscribe init") {
		t.Errorf("extract: expected 'mlscribe init' in output: %s", out)
	}
	if !strings.Contains(out, "mlscribe note") {
		t.Errorf("extract: expected 'mlscribe note' in output: %s", out)
	}
}

func TestInitFailsIfFileExists(t *testing.T) {
	dir := t.TempDir()
	doc := filepath.Join(dir, "experiment.md")

	run(t, "init", doc, "First")
	_, code := run(t, "init", doc, "Second")
	if code == 0 {
		t.Error("expected non-zero exit when file already exists")
	}
}

func TestPopOnEmptyDocument(t *testing.T) {
	// init creates one block, so pop removes it
	dir := t.TempDir()
	doc := filepath.Join(dir, "experiment.md")
	run(t, "init", doc, "Test")

	_, code := run(t, "pop", doc)
	if code != 0 {
		t.Error("popping the only block should succeed")
	}

	// Second pop should fail
	_, code = run(t, "pop", doc)
	if code == 0 {
		t.Error("expected non-zero exit when popping empty document")
	}
}

func TestVersion(t *testing.T) {
	out, code := run(t, "--version")
	if code != 0 {
		t.Fatalf("--version failed (code %d): %s", code, out)
	}
	if !strings.Contains(out, "mlscribe version") {
		t.Errorf("expected version string: %s", out)
	}
}

func TestGateInvalidStatus(t *testing.T) {
	dir := t.TempDir()
	doc := filepath.Join(dir, "experiment.md")
	run(t, "init", doc, "Test")

	_, code := run(t, "gate", doc, "my-gate", "INVALID")
	if code == 0 {
		t.Error("expected non-zero exit for invalid gate status")
	}
}

func TestMetricsWithBaselineDelta(t *testing.T) {
	dir := t.TempDir()
	doc := filepath.Join(dir, "experiment.md")
	run(t, "init", doc, "Test")

	// Write baseline file
	baselineFile := filepath.Join(dir, "baseline.json")
	if err := os.WriteFile(baselineFile, []byte(`{"accuracy": 0.90, "loss": 0.15}`), 0644); err != nil {
		t.Fatalf("writing baseline: %v", err)
	}

	out, code := run(t, "metrics", doc,
		"--json", `{"accuracy": 0.94, "loss": 0.12}`,
		"--baseline", baselineFile)
	if code != 0 {
		t.Fatalf("metrics with baseline failed (code %d): %s", code, out)
	}

	// Check document contains delta values
	content, err := os.ReadFile(doc)
	if err != nil {
		t.Fatalf("reading doc: %v", err)
	}
	if !strings.Contains(string(content), "Delta") {
		t.Errorf("expected Delta column in document: %s", string(content))
	}
}
