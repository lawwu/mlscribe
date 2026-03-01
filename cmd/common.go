// Package cmd contains subcommand implementations for mlscribe.
package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lawwu/mlscribe/markdown"
)

// readDocument reads blocks from a file, or returns an error.
func readDocument(path string) ([]markdown.Block, error) {
	return markdown.ReadFile(path)
}

// writeDocument serializes blocks and writes to a file.
func writeDocument(path string, blocks []markdown.Block) error {
	return markdown.WriteFile(path, blocks)
}

// readTextArg reads text from either the argument (if non-empty) or stdin.
func readTextArg(arg string) (string, error) {
	if arg != "" {
		return arg, nil
	}
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("reading stdin: %w", err)
	}
	return strings.Join(lines, "\n"), nil
}

// fatalf prints an error message and exits with code 1.
func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "mlscribe: "+format+"\n", args...)
	os.Exit(1)
}

// requireArgs checks that at least n positional args are present.
func requireArgs(fs *flag.FlagSet, n int, usage string) {
	if fs.NArg() < n {
		fmt.Fprintf(os.Stderr, "usage: %s\n", usage)
		os.Exit(1)
	}
}

// reorderArgs moves flag arguments (--flag value or --flag=value) before
// positional arguments so that flag.FlagSet.Parse works regardless of order.
// This is needed because Go's flag package stops parsing at the first
// non-flag argument.
func reorderArgs(args []string) []string {
	var flags []string
	var positional []string
	i := 0
	for i < len(args) {
		arg := args[i]
		if strings.HasPrefix(arg, "--") || (strings.HasPrefix(arg, "-") && len(arg) > 1 && arg[1] != '-') {
			// It's a flag; check if it has an inline value (--flag=value)
			if strings.Contains(arg, "=") {
				flags = append(flags, arg)
				i++
			} else {
				// The next arg is the value, if it doesn't start with '-'
				flags = append(flags, arg)
				i++
				if i < len(args) && !strings.HasPrefix(args[i], "-") {
					flags = append(flags, args[i])
					i++
				}
			}
		} else {
			positional = append(positional, arg)
			i++
		}
	}
	return append(flags, positional...)
}
