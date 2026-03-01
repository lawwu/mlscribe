package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lawwu/mlscribe/markdown"
)

// Extract emits CLI commands to recreate the document from scratch.
func Extract(args []string) {
	fs := flag.NewFlagSet("extract", flag.ExitOnError)
	fs.Parse(args)

	requireArgs(fs, 1, "mlscribe extract <file>")

	path := fs.Arg(0)

	blocks, err := readDocument(path)
	if err != nil {
		fatalf("reading document: %v", err)
	}

	fmt.Printf("#!/bin/bash\n# Extracted from %s\nset -e\n\n", path)

	for _, block := range blocks {
		switch b := block.(type) {
		case *markdown.TitleBlock:
			fmt.Printf("mlscribe init %s %q\n", path, b.Title)

		case *markdown.CommentaryBlock:
			// Use heredoc for multi-line notes
			if strings.Contains(b.Content, "\n") {
				fmt.Printf("mlscribe note %s <<'EOF'\n%s\nEOF\n", path, b.Content)
			} else {
				fmt.Printf("mlscribe note %s %q\n", path, b.Content)
			}

		case *markdown.CodeBlock:
			// Use heredoc for code
			fmt.Printf("mlscribe exec %s %s <<'EOF'\n%s\nEOF\n", path, b.Language, b.Code)

		case *markdown.OutputBlock:
			// Output blocks are recreated by exec, skip

		case *markdown.ImageOutputBlock:
			fmt.Printf("mlscribe image %s %s --alt %q\n", path, b.Filename, b.AltText)

		case *markdown.MetricsBlock:
			// Reconstruct metrics JSON
			pairs := make([]string, 0, len(b.Metrics))
			for _, m := range b.Metrics {
				pairs = append(pairs, fmt.Sprintf("%q: %q", m.Name, m.Value))
			}
			jsonStr := "{" + strings.Join(pairs, ", ") + "}"
			if b.Metrics[0].Baseline != "" {
				fmt.Printf("# Note: baseline values not re-emitted in extract\n")
			}
			fmt.Printf("mlscribe metrics %s --json '%s'\n", path, jsonStr)

		case *markdown.GateBlock:
			cmd := fmt.Sprintf("mlscribe gate %s %q %s", path, b.GateName, b.Status)
			if b.Decision != "" {
				cmd += fmt.Sprintf(" --decision %s", b.Decision)
			}
			if b.Confidence != "" {
				cmd += fmt.Sprintf(" --confidence %s", b.Confidence)
			}
			if b.Findings != "" {
				cmd += fmt.Sprintf(" --findings %q", b.Findings)
			}
			fmt.Println(cmd)

		case *markdown.ArtifactRefBlock:
			cmd := fmt.Sprintf("mlscribe artifact %s --skill %s --out-dir %s",
				path, b.SkillName, extractDir(b.ArtifactPath))
			if b.RunID != "" {
				cmd += fmt.Sprintf(" --run-id %s", b.RunID)
			}
			fmt.Println(cmd)
		}
	}
}

func extractDir(p string) string {
	idx := strings.LastIndex(p, "/")
	if idx < 0 {
		return "."
	}
	return p[:idx]
}
