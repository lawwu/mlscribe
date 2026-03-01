# mlscribe

A CLI tool for creating and managing ML/DS demo documents, with metrics tracking, gate decisions, and JSON artifact emission compatible with the [agentic-ml-plugin](https://github.com/lawwu/agentic-ml-plugin) schema.

## Installation

You can run it without installing it first using uvx:

```bash
uvx mlscribe --help
```

Or install it like this, then run `mlscribe --help`:

```bash
uv tool install mlscribe
# or
uv pip install mlscribe
```

You can also install the Go binary directly:

```bash
go install github.com/lawwu/mlscribe@latest
```

Or run it without installation like this:

```bash
go run github.com/lawwu/mlscribe@latest --help
```

## Usage

```
mlscribe <command> [options]
```

### Core Commands

| Command | Description |
|---------|-------------|
| `init <file> <title>` | Create a new document |
| `note <file> [text]` | Append commentary (reads from stdin if no text arg) |
| `exec <file> <lang> [code]` | Run code, capture output, append both to document |
| `image <file> <image-path>` | Copy image into doc dir and embed a reference |
| `pop <file>` | Remove the last block |
| `verify <file>` | Re-run all code blocks and compare output (exits 1 on mismatch) |
| `extract <file>` | Emit CLI commands to recreate the document |

### ML/DS Commands

| Command | Description |
|---------|-------------|
| `metrics <file> --json '{...}'` | Append a metrics table block |
| `gate <file> <name> <PASS\|FAIL\|SKIPPED>` | Append a gate decision block |
| `artifact <file> --skill <name> --out-dir <dir>` | Write a JSON artifact file |

## Examples

```bash
# Create a new document
mlscribe init experiment.md "Model Training Experiment"

# Add a note
mlscribe note experiment.md "Training ResNet-50 on CIFAR-10"

# Run code and capture output
mlscribe exec experiment.md python "print('hello world')"

# Record metrics
mlscribe metrics experiment.md --json '{"accuracy": 0.94, "loss": 0.12}'

# Record metrics with baseline delta
mlscribe metrics experiment.md --json '{"accuracy": 0.94}' --baseline baseline.json

# Record a gate decision
mlscribe gate experiment.md data-quality PASS --decision GO --confidence high

# Emit a JSON artifact compatible with agentic-ml-plugin
mlscribe artifact experiment.md --skill eda --out-dir ./artifacts \
  --run-id run-001 --json '{"decision": "GO", "summary": "Data quality good"}'

# Verify all code blocks still produce the same output
mlscribe verify experiment.md

# Extract CLI commands to recreate the document
mlscribe extract experiment.md
```

## Options

**exec options:**
- `--workdir <dir>` — Working directory for code execution

**metrics options:**
- `--json <json>` — Metrics as JSON object (or pipe via stdin)
- `--baseline <file>` — Baseline metrics file for delta computation

**gate options:**
- `--decision GO|NO-GO|CONDITIONAL` — Gate decision
- `--confidence high|medium|low` — Confidence level
- `--findings <text>` — Findings text or JSON

**artifact options:**
- `--skill <name>` — Skill name (required)
- `--out-dir <dir>` — Output directory (required)
- `--run-id <id>` — Run ID for artifact correlation
- `--json <json>` — Artifact payload JSON (or pipe via stdin)

## Supported Languages (exec)

`python`, `python3`, `bash`, `sh`, `zsh`, `ruby`, `node`, `javascript`, `js`, `r`

## Artifact Schema

JSON artifacts match the agentic-ml-plugin base schema (`schema_version: "1.0"`):

```json
{
  "schema_version": "1.0",
  "skill_name": "...",
  "run_id": "...",
  "timestamp": "...",
  "decision": "...",
  "confidence": "...",
  "summary": "...",
  "findings": "...",
  "next_commands": []
}
```

## Environment

| Variable | Description |
|----------|-------------|
| `MLSCRIBE_REMOTE_URL` | If set, POST the document to this URL after each write |

## Acknowledgements

Inspired by [simonw/showboat](https://github.com/simonw/showboat) by Simon Willison.
