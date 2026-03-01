// Package exec provides code execution for mlscribe.
package exec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCode executes a code snippet for the given language and returns combined stdout+stderr.
// workDir is the working directory for the subprocess; empty means current dir.
func RunCode(lang, code, workDir string) (string, error) {
	var cmd *exec.Cmd

	switch strings.ToLower(lang) {
	case "python", "python3":
		cmd = exec.Command("python3", "-c", code)
	case "bash", "sh":
		cmd = exec.Command("bash", "-c", code)
	case "zsh":
		cmd = exec.Command("zsh", "-c", code)
	case "ruby":
		cmd = exec.Command("ruby", "-e", code)
	case "node", "javascript", "js":
		cmd = exec.Command("node", "-e", code)
	case "r":
		cmd = exec.Command("Rscript", "-e", code)
	default:
		return "", fmt.Errorf("unsupported language: %s", lang)
	}

	if workDir != "" {
		cmd.Dir = workDir
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = os.Environ()

	err := cmd.Run()
	output := strings.TrimRight(out.String(), "\n")
	if err != nil {
		// Return output even on non-zero exit so the user sees error messages
		return output, fmt.Errorf("command exited with error: %w", err)
	}
	return output, nil
}

// RunFile executes a script file for the given language and returns combined stdout+stderr.
func RunFile(lang, filepath, workDir string) (string, error) {
	var cmd *exec.Cmd

	switch strings.ToLower(lang) {
	case "python", "python3":
		cmd = exec.Command("python3", filepath)
	case "bash", "sh":
		cmd = exec.Command("bash", filepath)
	case "zsh":
		cmd = exec.Command("zsh", filepath)
	case "ruby":
		cmd = exec.Command("ruby", filepath)
	case "node", "javascript", "js":
		cmd = exec.Command("node", filepath)
	case "r":
		cmd = exec.Command("Rscript", filepath)
	default:
		return "", fmt.Errorf("unsupported language: %s", lang)
	}

	if workDir != "" {
		cmd.Dir = workDir
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = os.Environ()

	err := cmd.Run()
	output := strings.TrimRight(out.String(), "\n")
	if err != nil {
		return output, fmt.Errorf("command exited with error: %w", err)
	}
	return output, nil
}
