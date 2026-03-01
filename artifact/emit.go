// Package artifact provides JSON artifact emission for mlscribe.
package artifact

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// EmitOptions configures artifact emission.
type EmitOptions struct {
	SkillName string
	RunID     string
	OutDir    string
	// ExtraJSON is raw JSON that will be merged into the artifact.
	// Fields from ExtraJSON take precedence over defaults.
	ExtraJSON json.RawMessage
}

// Emit writes a JSON artifact file to <outDir>/<skillName>.json.
// It merges the provided JSON with required schema fields.
// Returns the path of the written file.
func Emit(opts EmitOptions) (string, error) {
	if opts.SkillName == "" {
		return "", fmt.Errorf("skill name is required")
	}
	if opts.OutDir == "" {
		return "", fmt.Errorf("output directory is required")
	}

	// Start with a base document
	base := map[string]interface{}{
		"schema_version": SchemaVersion,
		"skill_name":     opts.SkillName,
		"run_id":         opts.RunID,
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
		"decision":       DecisionGO,
		"confidence":     ConfidenceMedium,
		"summary":        "",
		"findings":       []interface{}{},
	}

	// Merge extra JSON on top (extra fields win)
	if len(opts.ExtraJSON) > 0 {
		var extra map[string]interface{}
		if err := json.Unmarshal(opts.ExtraJSON, &extra); err != nil {
			return "", fmt.Errorf("parsing artifact JSON: %w", err)
		}
		for k, v := range extra {
			base[k] = v
		}
	}

	// Always ensure these are set (may have been overwritten by extra)
	if _, ok := base["schema_version"]; !ok {
		base["schema_version"] = SchemaVersion
	}

	// Ensure output directory exists
	if err := os.MkdirAll(opts.OutDir, 0755); err != nil {
		return "", fmt.Errorf("creating output directory: %w", err)
	}

	outPath := filepath.Join(opts.OutDir, opts.SkillName+".json")
	data, err := json.MarshalIndent(base, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling artifact: %w", err)
	}

	if err := os.WriteFile(outPath, data, 0644); err != nil {
		return "", fmt.Errorf("writing artifact file: %w", err)
	}

	return outPath, nil
}
