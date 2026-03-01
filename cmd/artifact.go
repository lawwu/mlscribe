package cmd

import (
	"flag"
	"fmt"

	"github.com/lawwu/mlscribe/artifact"
	"github.com/lawwu/mlscribe/markdown"
)

// Artifact writes a JSON artifact file and appends an ArtifactRefBlock to the document.
func Artifact(args []string) {
	fs := flag.NewFlagSet("artifact", flag.ExitOnError)
	skillName := fs.String("skill", "", "skill name (required)")
	outDir := fs.String("out-dir", "", "output directory for the artifact file (required)")
	runID := fs.String("run-id", "", "run ID for artifact correlation")
	jsonStr := fs.String("json", "", "artifact JSON payload")
	fs.Parse(reorderArgs(args))

	requireArgs(fs, 1, "mlscribe artifact <file> --skill <name> --out-dir <dir> [--run-id <id>] [--json '{...}']")

	path := fs.Arg(0)

	if *skillName == "" {
		fatalf("--skill is required")
	}
	if *outDir == "" {
		fatalf("--out-dir is required")
	}

	// Read JSON from flag or stdin
	var rawJSON string
	if *jsonStr != "" {
		rawJSON = *jsonStr
	} else {
		var err error
		rawJSON, err = readTextArg("")
		if err != nil {
			fatalf("reading artifact JSON: %v", err)
		}
	}

	opts := artifact.EmitOptions{
		SkillName: *skillName,
		RunID:     *runID,
		OutDir:    *outDir,
	}
	if rawJSON != "" {
		opts.ExtraJSON = []byte(rawJSON)
	}

	artifactPath, err := artifact.Emit(opts)
	if err != nil {
		fatalf("emitting artifact: %v", err)
	}

	// Append ArtifactRefBlock to document
	blocks, err := readDocument(path)
	if err != nil {
		fatalf("reading document: %v", err)
	}

	blocks = append(blocks, &markdown.ArtifactRefBlock{
		SkillName:    *skillName,
		ArtifactPath: artifactPath,
		RunID:        *runID,
	})

	if err := writeDocument(path, blocks); err != nil {
		fatalf("writing document: %v", err)
	}

	fmt.Printf("Artifact written: %s\n", artifactPath)
}
