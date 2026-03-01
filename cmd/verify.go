package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lawwu/mlscribe/exec"
	"github.com/lawwu/mlscribe/markdown"
)

// Verify re-runs all code blocks and compares their output to the stored output.
// Exits with code 1 if any output differs.
func Verify(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	workDir := fs.String("workdir", "", "working directory for execution")
	fs.Parse(args)

	requireArgs(fs, 1, "mlscribe verify <file>")

	path := fs.Arg(0)

	blocks, err := readDocument(path)
	if err != nil {
		fatalf("reading document: %v", err)
	}

	failures := 0
	i := 0
	for i < len(blocks) {
		block := blocks[i]
		codeBlock, ok := block.(*markdown.CodeBlock)
		if !ok {
			i++
			continue
		}

		// Look ahead for an OutputBlock
		var storedOutput string
		outputIdx := -1
		if i+1 < len(blocks) {
			if ob, ok := blocks[i+1].(*markdown.OutputBlock); ok {
				storedOutput = ob.Output
				outputIdx = i + 1
			}
		}

		// Re-run the code
		newOutput, runErr := exec.RunCode(codeBlock.Language, codeBlock.Code, *workDir)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "WARN: execution error for %s block: %v\n", codeBlock.Language, runErr)
		}

		if outputIdx >= 0 {
			if normalize(newOutput) != normalize(storedOutput) {
				fmt.Fprintf(os.Stderr, "MISMATCH for %s block:\n", codeBlock.Language)
				fmt.Fprintf(os.Stderr, "  stored:  %q\n", storedOutput)
				fmt.Fprintf(os.Stderr, "  current: %q\n", newOutput)
				failures++
			} else {
				fmt.Printf("OK: %s block\n", codeBlock.Language)
			}
			i = outputIdx + 1
		} else {
			fmt.Printf("OK: %s block (no stored output)\n", codeBlock.Language)
			i++
		}
	}

	if failures > 0 {
		fmt.Fprintf(os.Stderr, "%d block(s) failed verification\n", failures)
		os.Exit(1)
	}
	fmt.Println("All blocks verified successfully")
}

// normalize trims whitespace for comparison.
func normalize(s string) string {
	return strings.TrimSpace(s)
}
