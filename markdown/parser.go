// Package markdown provides parsing of mlscribe documents.
package markdown

import (
	"bufio"
	"strings"
)

// Parse reads a mlscribe markdown document and returns its blocks.
func Parse(content string) ([]Block, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var blocks []Block
	var lines []string

	// Collect all lines first
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Title block: starts with "# "
		if strings.HasPrefix(line, "# ") {
			block, consumed := parseTitleBlock(lines, i)
			blocks = append(blocks, block)
			i += consumed
			continue
		}

		// ML metrics block
		if line == "<!-- mlscribe:metrics -->" {
			block, consumed := parseMetricsBlock(lines, i)
			blocks = append(blocks, block)
			i += consumed
			continue
		}

		// ML gate block
		if line == "<!-- mlscribe:gate -->" {
			block, consumed := parseGateBlock(lines, i)
			blocks = append(blocks, block)
			i += consumed
			continue
		}

		// ML artifact_ref block
		if line == "<!-- mlscribe:artifact_ref -->" {
			block, consumed := parseArtifactRefBlock(lines, i)
			blocks = append(blocks, block)
			i += consumed
			continue
		}

		// Fenced code block
		if strings.HasPrefix(line, "```") {
			block, consumed := parseCodeOrOutputBlock(lines, i)
			blocks = append(blocks, block)
			i += consumed
			continue
		}

		// Image block: starts with "!["
		if strings.HasPrefix(line, "![") {
			block := parseImageBlock(line)
			blocks = append(blocks, block)
			i++
			continue
		}

		// Skip blank lines between blocks
		if strings.TrimSpace(line) == "" {
			i++
			continue
		}

		// Commentary block: collect consecutive non-special lines
		block, consumed := parseCommentaryBlock(lines, i)
		blocks = append(blocks, block)
		i += consumed
	}

	return blocks, nil
}

func parseTitleBlock(lines []string, start int) (*TitleBlock, int) {
	title := strings.TrimPrefix(lines[start], "# ")
	block := &TitleBlock{Title: title}
	i := start + 1

	// Skip blank line after title
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}

	// Check for meta comment block
	if i < len(lines) && lines[i] == "<!-- mlscribe:meta" {
		i++
		for i < len(lines) && lines[i] != "-->" {
			kv := lines[i]
			if after, ok := strings.CutPrefix(kv, "timestamp: "); ok {
				block.Timestamp = after
			} else if after, ok := strings.CutPrefix(kv, "version: "); ok {
				block.Version = after
			} else if after, ok := strings.CutPrefix(kv, "doc_id: "); ok {
				block.DocID = after
			}
			i++
		}
		if i < len(lines) && lines[i] == "-->" {
			i++ // consume "-->"
		}
	}

	return block, i - start
}

func parseCodeOrOutputBlock(lines []string, start int) (Block, int) {
	fence := lines[start]
	lang := strings.TrimPrefix(fence, "```")
	lang = strings.TrimSpace(lang)

	var body strings.Builder
	i := start + 1
	for i < len(lines) && lines[i] != "```" {
		if body.Len() > 0 {
			body.WriteByte('\n')
		}
		body.WriteString(lines[i])
		i++
	}
	if i < len(lines) {
		i++ // consume closing ```
	}

	if lang == "" {
		return &OutputBlock{Output: body.String()}, i - start
	}
	return &CodeBlock{Language: lang, Code: body.String()}, i - start
}

func parseImageBlock(line string) *ImageOutputBlock {
	// Format: ![alt](filename)
	line = strings.TrimPrefix(line, "![")
	parts := strings.SplitN(line, "](", 2)
	if len(parts) != 2 {
		return &ImageOutputBlock{AltText: line}
	}
	alt := parts[0]
	filename := strings.TrimSuffix(parts[1], ")")
	return &ImageOutputBlock{AltText: alt, Filename: filename}
}

func parseMetricsBlock(lines []string, start int) (*MetricsBlock, int) {
	block := &MetricsBlock{}
	i := start + 1

	// Find the header row to determine columns
	headerParsed := false
	hasBaseline := false

	for i < len(lines) && lines[i] != "<!-- /mlscribe:metrics -->" {
		line := lines[i]
		i++

		if !strings.HasPrefix(line, "|") {
			continue
		}
		// Skip separator row
		if strings.Contains(line, "---") {
			continue
		}

		cols := parseTableRow(line)
		if !headerParsed {
			// This is the header row
			for _, c := range cols {
				if strings.EqualFold(strings.TrimSpace(c), "baseline") {
					hasBaseline = true
					break
				}
			}
			headerParsed = true
			continue
		}

		// Data row
		if hasBaseline && len(cols) >= 4 {
			block.Metrics = append(block.Metrics, MetricEntry{
				Name:     strings.TrimSpace(cols[0]),
				Value:    strings.TrimSpace(cols[1]),
				Baseline: strings.TrimSpace(cols[2]),
				Delta:    strings.TrimSpace(cols[3]),
			})
		} else if len(cols) >= 2 {
			block.Metrics = append(block.Metrics, MetricEntry{
				Name:  strings.TrimSpace(cols[0]),
				Value: strings.TrimSpace(cols[1]),
			})
		}
	}

	if i < len(lines) {
		i++ // consume closing sentinel
	}

	return block, i - start
}

func parseGateBlock(lines []string, start int) (*GateBlock, int) {
	block := &GateBlock{}
	i := start + 1
	headerParsed := false

	for i < len(lines) && lines[i] != "<!-- /mlscribe:gate -->" {
		line := lines[i]
		i++

		// Gate name line: **Gate**: <name>
		if strings.HasPrefix(line, "**Gate**: ") {
			block.GateName = strings.TrimPrefix(line, "**Gate**: ")
			continue
		}

		// Findings line: **Findings**: <text>
		if strings.HasPrefix(line, "**Findings**: ") {
			block.Findings = strings.TrimPrefix(line, "**Findings**: ")
			continue
		}

		if !strings.HasPrefix(line, "|") {
			continue
		}
		if strings.Contains(line, "---") {
			headerParsed = true
			continue
		}
		if !headerParsed {
			continue
		}

		cols := parseTableRow(line)
		if len(cols) < 2 {
			continue
		}
		field := strings.TrimSpace(cols[0])
		value := strings.TrimSpace(cols[1])
		switch field {
		case "Status":
			block.Status = value
		case "Decision":
			block.Decision = value
		case "Confidence":
			block.Confidence = value
		}
	}

	if i < len(lines) {
		i++ // consume closing sentinel
	}

	return block, i - start
}

func parseArtifactRefBlock(lines []string, start int) (*ArtifactRefBlock, int) {
	block := &ArtifactRefBlock{}
	i := start + 1

	for i < len(lines) && lines[i] != "<!-- /mlscribe:artifact_ref -->" {
		line := lines[i]
		i++

		// **Artifact**: [skill](path) (run: runID)
		if strings.HasPrefix(line, "**Artifact**: [") {
			inner := strings.TrimPrefix(line, "**Artifact**: [")
			// Parse skill name
			parts := strings.SplitN(inner, "](", 2)
			if len(parts) == 2 {
				block.SkillName = parts[0]
				rest := parts[1]
				// Parse path and optional run ID
				if idx := strings.Index(rest, ") (run: "); idx >= 0 {
					block.ArtifactPath = rest[:idx]
					runPart := rest[idx+len(") (run: "):]
					block.RunID = strings.TrimSuffix(runPart, ")")
				} else {
					block.ArtifactPath = strings.TrimSuffix(rest, ")")
				}
			}
		}
	}

	if i < len(lines) {
		i++ // consume closing sentinel
	}

	return block, i - start
}

func parseCommentaryBlock(lines []string, start int) (*CommentaryBlock, int) {
	var sb strings.Builder
	i := start
	for i < len(lines) {
		line := lines[i]
		// Stop at special block starters
		if strings.HasPrefix(line, "# ") ||
			strings.HasPrefix(line, "```") ||
			strings.HasPrefix(line, "![") ||
			strings.HasPrefix(line, "<!-- mlscribe:") {
			break
		}
		if sb.Len() > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(line)
		i++
	}
	return &CommentaryBlock{Content: sb.String()}, i - start
}

// parseTableRow splits a markdown table row into columns.
func parseTableRow(line string) []string {
	line = strings.TrimPrefix(line, "|")
	line = strings.TrimSuffix(line, "|")
	parts := strings.Split(line, "|")
	return parts
}
