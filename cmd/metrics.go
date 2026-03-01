package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/lawwu/mlscribe/markdown"
)

// Metrics appends a metrics table block to the document.
func Metrics(args []string) {
	fs := flag.NewFlagSet("metrics", flag.ExitOnError)
	jsonStr := fs.String("json", "", "metrics as JSON object (e.g. '{\"accuracy\": 0.94}')")
	baselineFile := fs.String("baseline", "", "baseline metrics JSON file for delta computation")
	fs.Parse(reorderArgs(args))

	requireArgs(fs, 1, "mlscribe metrics <file> --json '{...}'")

	path := fs.Arg(0)

	// Read metrics JSON from flag or stdin
	var rawJSON string
	if *jsonStr != "" {
		rawJSON = *jsonStr
	} else {
		var err error
		rawJSON, err = readTextArg("")
		if err != nil {
			fatalf("reading metrics JSON: %v", err)
		}
	}

	if rawJSON == "" {
		fatalf("no metrics JSON provided (use --json or stdin)")
	}

	var metricsMap map[string]interface{}
	if err := json.Unmarshal([]byte(rawJSON), &metricsMap); err != nil {
		fatalf("parsing metrics JSON: %v", err)
	}

	// Optionally read baseline
	var baselineMap map[string]interface{}
	if *baselineFile != "" {
		data, err := os.ReadFile(*baselineFile)
		if err != nil {
			fatalf("reading baseline file: %v", err)
		}
		if err := json.Unmarshal(data, &baselineMap); err != nil {
			fatalf("parsing baseline JSON: %v", err)
		}
	}

	// Build sorted metric entries
	names := make([]string, 0, len(metricsMap))
	for k := range metricsMap {
		names = append(names, k)
	}
	sort.Strings(names)

	var entries []markdown.MetricEntry
	for _, name := range names {
		val := fmt.Sprintf("%v", metricsMap[name])
		entry := markdown.MetricEntry{Name: name, Value: val}

		if baselineMap != nil {
			if bv, ok := baselineMap[name]; ok {
				entry.Baseline = fmt.Sprintf("%v", bv)
				// Compute delta if both are numeric
				entry.Delta = computeDelta(metricsMap[name], bv)
			}
		}
		entries = append(entries, entry)
	}

	blocks, err := readDocument(path)
	if err != nil {
		fatalf("reading document: %v", err)
	}

	blocks = append(blocks, &markdown.MetricsBlock{Metrics: entries})
	if err := writeDocument(path, blocks); err != nil {
		fatalf("writing document: %v", err)
	}

	fmt.Printf("Appended metrics block (%d metrics)\n", len(entries))
}

func computeDelta(current, baseline interface{}) string {
	cf, cok := toFloat(current)
	bf, bok := toFloat(baseline)
	if !cok || !bok {
		return ""
	}
	delta := cf - bf
	if delta >= 0 {
		return fmt.Sprintf("+%s", formatFloat(delta))
	}
	return formatFloat(delta)
}

func toFloat(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	case string:
		f, err := strconv.ParseFloat(x, 64)
		return f, err == nil
	}
	return 0, false
}

func formatFloat(f float64) string {
	// Use up to 4 decimal places, strip trailing zeros
	s := strconv.FormatFloat(f, 'f', 4, 64)
	// Strip trailing zeros after decimal point
	for len(s) > 1 && s[len(s)-1] == '0' && s[len(s)-2] != '.' {
		s = s[:len(s)-1]
	}
	return s
}
