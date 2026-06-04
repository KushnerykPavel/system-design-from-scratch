# Architecture Note-Taking System

> Good notes compress future thinking; bad notes only preserve old confusion.

**Type:** Learn
**Company focus:** Balanced
**Learning goal:** Build a note-taking structure that preserves design decisions, open questions, and trade-offs in a format you can reuse during mocks and reviews.
**Prerequisites:** `00-setup-and-workflow/01-repo-setup-and-progress`
**Estimated time:** ~45 min
**Primary artifact:** architecture note template

## The Problem

Candidates often take notes as a linear transcript of what they said. That makes review painfully slow. Senior prep notes should preserve reasoning, not just chronology.

The note system in this course is deliberately opinionated: capture scope, sizing, architecture, deep dives, failure modes, and what changed your mind.

## Clarify

- Are the notes for active practice, later retrieval, or mock feedback handoff?
- What parts of a design answer do you regularly forget a week later: assumptions, numbers, or trade-offs?
- If a mock partner read the notes cold, what would they need in order to continue the conversation?

## Requirements

### Functional

- Capture the prompt, clarified scope, and assumptions.
- Record the capacity model before the architecture section.
- Preserve trade-offs, failure modes, and unresolved questions.

### Non-functional

- Notes must be skimmable in under three minutes.
- The template must work for both product backends and infrastructure prompts.
- The system should avoid turning note-taking into a second full-time task.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Notes per lesson | 1 primary note plus revisions | favors a stable, reusable template |
| Read frequency | high during reviews and mocks | structure matters more than visual flourish |
| Storage | tiny text files | markdown is enough |
| Peak factor | heavy use before interviews | retrieval speed matters more than formatting |
| Rough cost | near zero | optimize for consistency, not tooling |

## Architecture

Organize each note around the same arc as the lesson:

1. Prompt and clarified scope
2. Prioritized requirements
3. Capacity sheet summary
4. High-level architecture
5. One or two deep dives
6. Failure modes and observability
7. Trade-offs and redesign triggers
8. Debrief notes

This architecture turns every note into a portable design artifact instead of a diary.

## Data Model & APIs

Recommended headings:

- `Prompt`
- `Clarifications`
- `Requirements`
- `Capacity model`
- `Architecture`
- `Deep dives`
- `Failure modes`
- `Observability`
- `Trade-offs`
- `What I missed`

Suggested metadata:

- lesson slug
- date
- company focus
- mock or solo mode
- confidence level

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| notes become transcripts | long paragraphs with no structure | force fixed headings and short bullets |
| numbers get omitted | architecture section appears before capacity model | keep a mandatory sizing block near the top |
| trade-offs disappear | notes read like a final answer | add a dedicated trade-off table every time |
| notes are never revisited | no links from mocks or progress tracking | attach notes to lesson status and debriefs |

## Observability

- metric: percentage of notes containing a filled capacity section
- metric: percentage of notes with at least one explicit trade-off
- metric: average time to skim a prior note before a mock
- SLO: a prior design note should be understandable by your future self within three minutes

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| fixed template headings | faster review and comparison | less expressive for unusual prompts | free-form journaling |
| short bullets over long prose | high scan speed | may hide nuance unless you rewrite carefully | narrative notes |
| separate debrief block | captures learning loop | one more section to maintain | mixing feedback into architecture body |

## Interview It

**Google framing:** "How would you capture reusable notes from system design practice so future sessions improve rather than restart from zero?"

**Cloudflare framing:** "How would you structure notes for platform and edge designs where operational details and failure assumptions matter as much as the core topology?"

**Follow-ups:**
1. How would the note template change for a low-latency prompt?
2. What note sections are optional during a 20-minute practice drill?
3. How do you mark uncertain claims without slowing yourself down?
4. What should a mock interviewer add that the candidate should not overwrite?
5. How do you avoid copying polished final answers into your notes?

## Ship It

- `outputs/architecture-note-template.md`
- `outputs/design-review-note-quality.md`

## Exercises

1. **Easy** — Convert one old unstructured design note into the new template.
2. **Medium** — Add a section that highlights redesign triggers at 10x scale.
3. **Hard** — Create one note template variant for solo practice and one for live mocks, then justify the differences.

## Further Reading

- [Markdown Guide](https://www.markdownguide.org/basic-syntax/) — enough structure without tool lock-in
- [liquidslr/system-design-notes](https://github.com/liquidslr/system-design-notes) — useful reference for repeatable interview sections
