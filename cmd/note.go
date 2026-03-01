package cmd

import (
	"flag"

	"github.com/lawwu/mlscribe/markdown"
)

// Note appends a commentary block to the document.
func Note(args []string) {
	fs := flag.NewFlagSet("note", flag.ExitOnError)
	fs.Parse(args)

	requireArgs(fs, 1, "mlscribe note <file> [text]")

	path := fs.Arg(0)
	var text string
	if fs.NArg() >= 2 {
		text = fs.Arg(1)
	}

	text, err := readTextArg(text)
	if err != nil {
		fatalf("reading text: %v", err)
	}
	if text == "" {
		fatalf("no text provided")
	}

	blocks, err := readDocument(path)
	if err != nil {
		fatalf("reading document: %v", err)
	}

	blocks = append(blocks, &markdown.CommentaryBlock{Content: text})
	if err := writeDocument(path, blocks); err != nil {
		fatalf("writing document: %v", err)
	}
}
