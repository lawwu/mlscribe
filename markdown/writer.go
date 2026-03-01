// Package markdown provides writing of mlscribe documents.
package markdown

import (
	"os"
	"strings"
)

// Render serializes a slice of blocks to a markdown string.
func Render(blocks []Block) string {
	var sb strings.Builder
	for i, b := range blocks {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(b.Render())
	}
	return sb.String()
}

// ReadFile reads a document file and parses it into blocks.
func ReadFile(path string) ([]Block, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(string(data))
}

// WriteFile serializes blocks and writes them to a file.
func WriteFile(path string, blocks []Block) error {
	content := Render(blocks)
	return os.WriteFile(path, []byte(content), 0644)
}
