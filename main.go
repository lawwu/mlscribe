// mlscribe: ML/DS demo document CLI tool.
// Inspired by simonw/showboat, extended with metrics, gate decisions,
// and JSON artifact emission compatible with the agentic-ml-plugin schema.
package main

import (
	"fmt"
	"os"

	_ "embed"

	"github.com/lawwu/mlscribe/cmd"
)

//go:embed help.txt
var helpText string

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		fmt.Print(helpText)
		os.Exit(1)
	}

	subcmd := os.Args[1]
	args := os.Args[2:]

	switch subcmd {
	case "init":
		cmd.Init(args)
	case "note":
		cmd.Note(args)
	case "exec":
		cmd.Exec(args)
	case "image":
		cmd.Image(args)
	case "pop":
		cmd.Pop(args)
	case "verify":
		cmd.Verify(args)
	case "extract":
		cmd.Extract(args)
	case "metrics":
		cmd.Metrics(args)
	case "gate":
		cmd.Gate(args)
	case "artifact":
		cmd.Artifact(args)
	case "--version", "version":
		fmt.Printf("mlscribe version %s\n", version)
	case "--help", "help", "-h":
		fmt.Print(helpText)
	default:
		fmt.Fprintf(os.Stderr, "mlscribe: unknown subcommand %q\n\n", subcmd)
		fmt.Print(helpText)
		os.Exit(1)
	}
}
