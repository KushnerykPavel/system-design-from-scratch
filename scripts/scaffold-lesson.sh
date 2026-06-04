#!/usr/bin/env bash
# Scaffold a system design lesson with docs, Go code stub, tests, outputs, and quiz.
set -euo pipefail

if [[ $# -lt 2 ]]; then
  cat <<'USAGE' >&2
Usage: scripts/scaffold-lesson.sh <phase-dir> <lesson-slug> [title]

Examples:
  scripts/scaffold-lesson.sh 03-design-framework-and-timing 01-four-step-interview-loop "Four-Step Interview Loop"
  scripts/scaffold-lesson.sh 14-rate-limiters-ids-and-hashing 02-distributed-rate-limiter
USAGE
  exit 2
fi

PHASE="$1"
LESSON="$2"
TITLE="${3:-}"

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PHASE_DIR="$REPO_ROOT/phases/$PHASE"
LESSON_DIR="$PHASE_DIR/$LESSON"

[[ -d "$PHASE_DIR" ]] || { echo "error: missing phase dir phases/$PHASE" >&2; exit 1; }
[[ ! -e "$LESSON_DIR" ]] || { echo "error: lesson exists phases/$PHASE/$LESSON" >&2; exit 1; }
[[ "$LESSON" =~ ^[0-9]{2}-[a-z0-9-]+$ ]] || { echo "error: lesson slug must match NN-kebab-case" >&2; exit 1; }

mkdir -p "$LESSON_DIR/code" "$LESSON_DIR/docs" "$LESSON_DIR/outputs"

PRETTY_TITLE="$TITLE"
if [[ -z "$PRETTY_TITLE" ]]; then
  PRETTY_TITLE="$(echo "${LESSON#[0-9][0-9]-}" | tr '-' ' ' | awk '{for (i=1; i<=NF; i++) $i=toupper(substr($i,1,1)) substr($i,2);}1')"
fi

cat >"$LESSON_DIR/docs/en.md" <<EOF
# $PRETTY_TITLE

> [One-line motto.]

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** [one sentence]
**Prerequisites:** [prior lessons]
**Estimated time:** ~75 min
**Primary artifact:** [interview card / capacity sheet / simulator]

## The Problem

[Prompt and context.]

## Clarify

- [question]
- [question]
- [assumption]

## Requirements

### Functional

- [requirement]

### Non-functional

- [requirement]

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| QPS | | |
| Storage | | |
| Bandwidth | | |
| Peak factor | | |
| Rough cost | | |

## Architecture

[Diagram or component narrative.]

## Data Model & APIs

[Entities, boundaries, interfaces.]

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
|         |           |            |

## Observability

- [metric]
- [log]
- [trace]
- [SLO]

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
|          |         |      |                      |

## Interview It

**Google framing:** [prompt]

**Cloudflare framing:** [prompt]

## Ship It

[Artifact list.]

## Exercises

1. **Easy** — [variation]
2. **Medium** — [variation]
3. **Hard** — [variation]

## Further Reading

- [link] — [why]
EOF

cat >"$LESSON_DIR/code/main.go" <<'EOF'
package main

func main() {}
EOF

cat >"$LESSON_DIR/code/main_test.go" <<'EOF'
package main

import "testing"

func TestPlaceholder(t *testing.T) {}
EOF

cat >"$LESSON_DIR/quiz.json" <<'EOF'
{
  "questions": []
}
EOF

touch "$LESSON_DIR/outputs/.gitkeep"

echo "created phases/$PHASE/$LESSON/"
