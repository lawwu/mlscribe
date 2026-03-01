package cmd

import (
	"flag"
	"fmt"

	"github.com/lawwu/mlscribe/exec"
	"github.com/lawwu/mlscribe/markdown"
)

// Image copies an image into the document directory and appends an image block.
func Image(args []string) {
	fs := flag.NewFlagSet("image", flag.ExitOnError)
	alt := fs.String("alt", "", "alt text for the image")
	fs.Parse(reorderArgs(args))

	requireArgs(fs, 2, "mlscribe image <file> <image-path>")

	path := fs.Arg(0)
	imagePath := fs.Arg(1)

	blocks, err := readDocument(path)
	if err != nil {
		fatalf("reading document: %v", err)
	}

	filename, err := exec.CopyImageToDocDir(imagePath, path)
	if err != nil {
		fatalf("copying image: %v", err)
	}

	altText := *alt
	if altText == "" {
		altText = filename
	}

	blocks = append(blocks, &markdown.ImageOutputBlock{
		AltText:  altText,
		Filename: filename,
	})

	if err := writeDocument(path, blocks); err != nil {
		fatalf("writing document: %v", err)
	}

	fmt.Printf("Embedded image: %s\n", filename)
}
