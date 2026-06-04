#!/usr/bin/env bash
# Bootstrap a new course from this template.
# Reads course.yml, substitutes {{PLACEHOLDERS}} across template files,
# creates phase directories, and seeds ROADMAP.md.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

if [[ ! -f course.yml ]]; then
  echo "error: course.yml not found in $REPO_ROOT" >&2
  exit 1
fi

# Lightweight YAML reader using Python (no PyYAML dep — uses regex for our flat schema).
python3 - "$REPO_ROOT" <<'PY'
import os, re, sys, datetime, pathlib, json, shutil

root = pathlib.Path(sys.argv[1])
text = (root / "course.yml").read_text()

def grab(key, default=""):
    m = re.search(rf'^\s*{re.escape(key)}\s*:\s*"?(.*?)"?\s*$', text, re.M)
    return m.group(1) if m else default

title    = grab("title")
slug     = grab("slug")
tagline  = grab("tagline")
topic    = grab("topic")
website  = grab("website")
github   = grab("github")

# Languages list (simple parser).
langs = re.findall(r'^\s*-\s+([a-z][a-z0-9+]*)\s*$', text, re.M)
# Filter to language names by intersecting with a known set + entries appearing under `languages:` block.
lang_block = re.search(r'languages:\s*\n((?:\s*-\s*\w+\s*\n)+)', text)
if lang_block:
    langs = re.findall(r'-\s*(\w+)', lang_block.group(1))
else:
    langs = ["python"]

primary_lang = langs[0] if langs else "python"
run_cmd_map = {
    "python":     "python code/main.py",
    "rust":       "cargo run --manifest-path code/Cargo.toml",
    "typescript": "ts-node code/main.ts",
    "julia":      "julia code/main.jl",
}
run_cmd = run_cmd_map.get(primary_lang, f"# run code/main.<ext>")

# Phases.
phase_lines = re.findall(
    r'-\s*\{\s*num:\s*"?(\d+)"?\s*,\s*slug:\s*"?([a-z0-9-]+)"?\s*,\s*name:\s*"([^"]+)"\s*,\s*hours:\s*(\d+)\s*\}',
    text,
)
total_phases = len(phase_lines)
total_hours  = sum(int(h) for _,_,_,h in phase_lines)
total_lessons = 0  # filled after lessons are scaffolded

# Substitution map.
year = datetime.date.today().year
sub = {
    "COURSE_TITLE":   title,
    "COURSE_SLUG":    slug,
    "COURSE_TAGLINE": tagline,
    "TOPIC":          topic,
    "GITHUB_URL":     github,
    "WEBSITE_URL":    website,
    "RUN_CMD":        run_cmd,
    "PRIMARY_LANG":   primary_lang,
    "TOTAL_PHASES":   str(total_phases),
    "TOTAL_LESSONS":  str(total_lessons),
    "TOTAL_HOURS":    str(total_hours),
    "YEAR":           str(year),
    "AUTHOR":         os.environ.get("USER","you"),
    "MAINTAINER_EMAIL": os.environ.get("GIT_EMAIL","maintainer@example.com"),
    "LANGUAGES_BLOCK": "\n".join(f"- **{l.capitalize()}**" for l in langs),
    "PHASE_TABLE":    "\n".join(
        f"| {num} | [{name}](phases/{num}-{slg}) | ~{h}h |"
        for num,slg,name,h in phase_lines
    ) or "| _no phases yet — edit course.yml_ | | |",
}
sub["PHASE_TABLE"] = "| #  | Phase | Hours |\n|----|-------|-------|\n" + sub["PHASE_TABLE"]

targets = [
    "README.md","AGENTS.md","ROADMAP.md","CHANGELOG.md","LICENSE",
    "CODE_OF_CONDUCT.md","FORKING.md","CONTRIBUTING.md",
    ".claude/skills/find-your-level/SKILL.md",
    ".claude/skills/check-understanding/SKILL.md",
    ".claude/skills/my-progress/SKILL.md",
]

for rel in targets:
    p = root / rel
    if not p.exists():
        continue
    body = p.read_text()
    for k,v in sub.items():
        body = body.replace("{{"+k+"}}", v)
    p.write_text(body)
    print(f"filled: {rel}")

# Create phase directories.
for num,slg,name,h in phase_lines:
    d = root / "phases" / f"{num}-{slg}"
    d.mkdir(parents=True, exist_ok=True)
    readme = d / "README.md"
    if not readme.exists():
        readme.write_text(f"# Phase {num} — {name}\n\n~{h} hours.\n\nLessons live in `NN-lesson-slug/` subdirectories.\n")
    print(f"phase dir: phases/{num}-{slg}/")

# Update totals back into course.yml (totals block).
new_totals = f"totals:\n  phases: {total_phases}\n  lessons: {total_lessons}\n  hours: {total_hours}\n"
ymlp = root / "course.yml"
yml = ymlp.read_text()
yml = re.sub(r'totals:\s*\n(?:\s+\w+:\s*\d+\s*\n){1,3}', new_totals, yml)
ymlp.write_text(yml)

# Remove the template's own self-readme.
template_readme = root / "TEMPLATE_README.md"
if template_readme.exists():
    template_readme.unlink()
    print("removed: TEMPLATE_README.md")

print("\ndone. next:")
print("  1. review README.md, ROADMAP.md, AGENTS.md")
print("  2. add lessons: scripts/scaffold-lesson.sh <phase-slug> <NN-lesson-slug> \"Title\"")
print("  3. git init && git add -A && git commit -m \"feat: bootstrap\"")
PY
