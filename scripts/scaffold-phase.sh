#!/usr/bin/env bash
# Add a new phase directory and a stub README.
# Usage: scripts/scaffold-phase.sh NN phase-slug "Phase Name"
set -euo pipefail

if [[ $# -lt 3 ]]; then
  echo "usage: scripts/scaffold-phase.sh NN phase-slug \"Phase Name\"" >&2
  exit 2
fi

NUM="$1"
SLUG="$2"
NAME="$3"

if [[ ! "$NUM" =~ ^[0-9]{2}$ ]]; then
  echo "error: NN must be two digits (e.g. 03)" >&2
  exit 1
fi
if [[ ! "$SLUG" =~ ^[a-z0-9-]+$ ]]; then
  echo "error: slug must be kebab-case" >&2
  exit 1
fi

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DIR="$REPO_ROOT/phases/$NUM-$SLUG"

if [[ -e "$DIR" ]]; then
  echo "error: phase already exists: $DIR" >&2
  exit 1
fi

mkdir -p "$DIR"
cat >"$DIR/README.md" <<EOF
# Phase $NUM — $NAME

Lessons live in \`NN-lesson-slug/\` subdirectories.

Scaffold a lesson:

    scripts/scaffold-lesson.sh $NUM-$SLUG 01-first-lesson "First Lesson"
EOF

echo "created phases/$NUM-$SLUG/"
echo "next: add a row in ROADMAP.md and update course.yml phases list."
