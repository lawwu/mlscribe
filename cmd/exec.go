package cmd

import (
	"flag"
	"fmt"

	"github.com/lawwu/mlscribe/exec"
	"github.com/lawwu/mlscribe/markdown"
)

// Exec runs a code snippet, captures output, and appends both to the document.
func Exec(args []string) {
	fs := flag.NewFlagSet("exec", flag.ExitOnError)
	workDir := fs.String("workdir", "", "working directory for execution")
	fs.Parse(reorderArgs(args))

	requireArgs(fs, 2, "mlscribe exec <file> <lang> [code]")

	path := fs.Arg(0)
	lang := fs.Arg(1)
	var code string
	if fs.NArg() >= 3 {
		code = fs.Arg(2)
	}

	code, err := readTextArg(code)
	if err != nil {
		fatalf("reading code: %v", err)
	}
	if code == "" {
		fatalf("no code provided")
	}

	blocks, err := readDocument(path)
	if err != nil {
		fatalf("reading document: %v", err)
	}

	output, runErr := exec.RunCode(lang, code, *workDir)

	// Always append code and output blocks (even on error)
	blocks = append(blocks, &markdown.CodeBlock{Language: lang, Code: code})
	blocks = append(blocks, &markdown.OutputBlock{Output: output})

	if err := writeDocument(path, blocks); err != nil {
		fatalf("writing document: %v", err)
	}

	// Print output to stdout
	if output != "" {
		fmt.Println(output)
	}

	if runErr != nil {
		fatalf("execution error: %v", runErr)
	}
}
