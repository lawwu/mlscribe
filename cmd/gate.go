package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/lawwu/mlscribe/artifact"
	"github.com/lawwu/mlscribe/markdown"
)

// Gate appends a gate decision block to the document.
func Gate(args []string) {
	fs := flag.NewFlagSet("gate", flag.ExitOnError)
	decision := fs.String("decision", "", "GO | NO-GO | CONDITIONAL")
	confidence := fs.String("confidence", "", "high | medium | low")
	findings := fs.String("findings", "", "findings text or JSON")
	fs.Parse(reorderArgs(args))

	requireArgs(fs, 3, "mlscribe gate <file> <gate-name> <PASS|FAIL|SKIPPED>")

	path := fs.Arg(0)
	gateName := fs.Arg(1)
	status := strings.ToUpper(fs.Arg(2))

	if !artifact.IsValidGateStatus(status) {
		fatalf("invalid gate status %q (must be PASS, FAIL, or SKIPPED)", status)
	}

	if *decision != "" && !artifact.IsValidDecision(*decision) {
		fatalf("invalid decision %q (must be GO, NO-GO, or CONDITIONAL)", *decision)
	}

	if *confidence != "" && !artifact.IsValidConfidence(*confidence) {
		fatalf("invalid confidence %q (must be high, medium, or low)", *confidence)
	}

	blocks, err := readDocument(path)
	if err != nil {
		fatalf("reading document: %v", err)
	}

	block := &markdown.GateBlock{
		GateName:   gateName,
		Status:     status,
		Decision:   *decision,
		Confidence: *confidence,
		Findings:   *findings,
	}

	blocks = append(blocks, block)
	if err := writeDocument(path, blocks); err != nil {
		fatalf("writing document: %v", err)
	}

	fmt.Printf("Appended gate block: %s [%s]\n", gateName, status)
}
