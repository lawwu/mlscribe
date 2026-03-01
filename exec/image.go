// Package exec provides image validation and copying for mlscribe.
package exec

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var validImageExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".svg":  true,
	".webp": true,
}

// ValidateImage checks that the file exists and has a supported extension.
func ValidateImage(path string) error {
	ext := strings.ToLower(filepath.Ext(path))
	if !validImageExtensions[ext] {
		return fmt.Errorf("unsupported image format %q (supported: png, jpg, jpeg, gif, svg, webp)", ext)
	}
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("image file not found: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("%s is a directory, not an image file", path)
	}
	return nil
}

// CopyImageToDocDir copies an image file into the document's directory.
// Returns the filename (without directory) that was copied.
func CopyImageToDocDir(srcPath, docPath string) (string, error) {
	if err := ValidateImage(srcPath); err != nil {
		return "", err
	}

	docDir := filepath.Dir(docPath)
	filename := filepath.Base(srcPath)
	dstPath := filepath.Join(docDir, filename)

	// If src and dst are the same file, nothing to do
	srcAbs, _ := filepath.Abs(srcPath)
	dstAbs, _ := filepath.Abs(dstPath)
	if srcAbs == dstAbs {
		return filename, nil
	}

	src, err := os.Open(srcPath)
	if err != nil {
		return "", fmt.Errorf("opening source image: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("creating destination image: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("copying image: %w", err)
	}

	return filename, nil
}
