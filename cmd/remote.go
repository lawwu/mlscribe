package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/lawwu/mlscribe/markdown"
)

// PostToRemote sends the document content to the MLSCRIBE_REMOTE_URL if set.
// This is called automatically after writes when the env var is present.
func PostToRemote(path string, blocks []markdown.Block) {
	remoteURL := os.Getenv("MLSCRIBE_REMOTE_URL")
	if remoteURL == "" {
		return
	}

	content := markdown.Render(blocks)
	req, err := http.NewRequest("POST", remoteURL, bytes.NewBufferString(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlscribe: remote: building request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "text/markdown")
	req.Header.Set("X-Mlscribe-Document", path)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlscribe: remote: posting document: %v\n", err)
		return
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 400 {
		fmt.Fprintf(os.Stderr, "mlscribe: remote: server returned %d\n", resp.StatusCode)
	}
}
