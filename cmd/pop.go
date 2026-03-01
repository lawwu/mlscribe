package cmd

import (
	"flag"
	"fmt"
)

// Pop removes the last block from the document.
func Pop(args []string) {
	fs := flag.NewFlagSet("pop", flag.ExitOnError)
	fs.Parse(args)

	requireArgs(fs, 1, "mlscribe pop <file>")

	path := fs.Arg(0)

	blocks, err := readDocument(path)
	if err != nil {
		fatalf("reading document: %v", err)
	}

	if len(blocks) == 0 {
		fatalf("document is empty, nothing to pop")
	}

	removed := blocks[len(blocks)-1]
	blocks = blocks[:len(blocks)-1]

	if err := writeDocument(path, blocks); err != nil {
		fatalf("writing document: %v", err)
	}

	fmt.Printf("Removed %s block\n", removed.Type())
}
