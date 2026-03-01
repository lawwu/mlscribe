package cmd

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lawwu/mlscribe/markdown"
)

// Init creates a new mlscribe document.
func Init(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	version := fs.String("version", "1.0", "document version")
	fs.Parse(args)

	requireArgs(fs, 2, "mlscribe init <file> <title>")

	path := fs.Arg(0)
	title := fs.Arg(1)

	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		fatalf("file already exists: %s", path)
	}

	block := &markdown.TitleBlock{
		Title:     title,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   *version,
		DocID:     uuid.New().String(),
	}

	blocks := []markdown.Block{block}
	if err := writeDocument(path, blocks); err != nil {
		fatalf("writing document: %v", err)
	}

	fmt.Printf("Created %s\n", path)
}
