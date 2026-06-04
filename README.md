# System Design from Scratch

> *From ambiguous prompts to senior-level architecture decisions. Design resilient systems, defend trade-offs, and practice Google and Cloudflare interview loops.*

![status](https://img.shields.io/badge/status-bootstrap-blue) ![license](https://img.shields.io/badge/license-MIT-blue) ![phases](https://img.shields.io/badge/phases-24-green) ![lessons](https://img.shields.io/badge/lessons-180-green) ![hours](https://img.shields.io/badge/hours-~248-green)

**180 lessons. 24 phases. ~248 hours.**

This course is for engineers who already know the usual system design buzzwords and want something sharper: better requirement handling, better trade-off language, better failure-mode depth, and better interview performance for senior-level loops.

You do not just memorize "design Twitter" answers. You practice:
- clarifying ambiguous prompts
- estimating scale fast
- choosing architecture under constraints
- defending consistency and cost trade-offs
- planning observability, rollout, and incident response
- redesigning when the interviewer changes the rules

> ⚠️ **Educational use only.** Designs, diagrams, and code are optimized for learning and interview reasoning, not direct production deployment.

## Why this course

Most system design material falls into two weak extremes:

- **Notes-only prep**: broad coverage, but shallow reasoning and no reusable artifacts.
- **Production architecture blogs**: rich detail, but not structured for interview training.

This course bridges both. It borrows the breadth checklist and 4-step interview loop popularized in repositories like [`liquidslr/system-design-notes`](https://github.com/liquidslr/system-design-notes), then goes further with:

- explicit capacity worksheets
- failure-mode tables
- observability and SLO planning
- migration and rollout strategy
- Google-specific and Cloudflare-specific interview framing
- mock interview rubrics and scenario banks

## Target roles

- **Google Senior / Staff SWE** — backend, infra, storage, serving, SRE-adjacent loops
- **Google Senior / Staff SRE** — architecture, reliability, scaling, operational trade-offs
- **Cloudflare Senior SWE** — edge, traffic, cache, gateway, abuse, origin protection
- **Cloudflare Senior Platform / Infra** — edge control planes, routing, resilience, observability

## Course shape

| Block | Phases | Hours |
|-------|--------|-------|
| Workflow, scoping, estimation | 00-03 | 34h |
| Core architecture building blocks | 04-13 | 102h |
| Canonical interview systems | 14-20 | 78h |
| Company specialization + mocks | 21-23 | 34h |

See [ROADMAP.md](ROADMAP.md) for the full lesson inventory.

## Languages

- **Go** — primary language for small simulators, validators, and reference implementations
- **Markdown** — primary medium for architecture notes, answer templates, and review artifacts

## Lesson structure

Every lesson follows the same system-design arc:

1. **The Problem** — the prompt, user, workload, and business pressure
2. **Clarify** — scope, assumptions, and success criteria
3. **Requirements** — functional + non-functional priorities
4. **Capacity Model** — QPS, storage, bandwidth, latency, cost rough sizing
5. **Architecture** — high-level design and main components
6. **Data Model & APIs** — interfaces, contracts, ownership boundaries
7. **Failure Modes** — what breaks, how you detect it, how you degrade safely
8. **Observability** — metrics, logs, traces, dashboards, SLOs
9. **Trade-offs** — alternatives and why this design wins here
10. **Interview It** — Google / Cloudflare framing and follow-ups
11. **Ship It** — reusable output in `outputs/`
12. **Exercises** — redesign, extension, and stress-constraint drills

See [LESSON_TEMPLATE.md](LESSON_TEMPLATE.md).

## Reusable artifacts

Each major lesson ships at least one high-signal artifact:

- **Design review prompts** — `outputs/design-review-*.md`
- **Interview cards** — `outputs/interview-card-*.md`
- **Capacity worksheets** — `outputs/capacity-sheet-*.md`
- **Trade-off matrices** — `outputs/tradeoff-matrix-*.md`
- **Failure-mode checklists** — `outputs/failure-checklist-*.md`
- **Observability checklists** — `outputs/observability-checklist-*.md`
- **Skills** — `outputs/skill-*.md`
- **Small Go tools** — `outputs/tool-*/` or lesson `code/`

## Quick start

```bash
git clone https://github.com/KushnerykPavel/system-design-from-scratch
cd system-design-from-scratch
# pick a lesson
cd phases/03-design-framework-and-timing/01-four-step-interview-loop
go test ./...
go run code/main.go --prompt="Design a rate limiter"
```

## Built-in agent commands

Slash commands ship with the repo and are backed by `.claude/skills/`:

| Command | What it does |
|---------|--------------|
| `/find-your-level` | Placement quiz and recommended starting phase |
| `/check-understanding <phase>` | Per-phase quiz and review recommendations |
| `/lesson <phase> <lesson>` | Guided lesson tutoring loop |
| `/design-review <phase> <lesson>` | Architecture review against the lesson rubric |
| `/capacity-drill <topic>` | Fast estimation round with QPS, storage, bandwidth, and cost |
| `/tradeoff-drill <topic>` | Trade-off sparring: consistency, latency, cost, complexity |
| `/mock-interview google-system-design` | Senior/staff Google-style mock loop |
| `/mock-interview cloudflare-system-design` | Cloudflare edge/platform mock loop |
| `/my-progress` | Personal dashboard and next-step recommendation |

Progress persists to `.progress.json` (git-ignored).
Use [.progress.example.json](/Users/pavelkushneryk/Documents/vsprojects/linkedin_blog/cryptograpy_education/system-design-from-scratch/.progress.example.json) as the reference shape for richer quiz, review, and confidence tracking.
Session feedback should use the shared dimensions `clarification`, `requirements`, `sizing`, `architecture`, `deep_dive`, `failure_modes`, `observability`, `trade_offs`, and `communication`.
The progress validator can also sync review dates from lesson results with `go run phases/00-setup-and-workflow/01-repo-setup-and-progress/code/main.go -progress .progress.json -sync-reviews`.
Recurring `mistake_tags` now feed review priority and weekly drill recommendations.
Adaptive recommendations can now prefer short practice modes like `capacity_drill` and `tradeoff_drill` over a full lesson replay when the weakness is narrow.
The progress dashboard now also surfaces completion breakdown, assisted ratio, streaks, and a compact weekly summary.

## Initial sample lessons

This bootstrap includes three representative lessons:

- [Four-Step Interview Loop](phases/03-design-framework-and-timing/01-four-step-interview-loop/docs/en.md)
- [Distributed Rate Limiter](phases/14-rate-limiters-ids-and-hashing/02-distributed-rate-limiter/docs/en.md)
- [Global API Edge Gateway](phases/22-cloudflare-edge-platform-design/01-global-api-edge-gateway/docs/en.md)

These are here to prove out the teaching model before the rest of the roadmap is authored.

## Acknowledgments

- Coverage inspiration and baseline interview framing from [`liquidslr/system-design-notes`](https://github.com/liquidslr/system-design-notes)
- Course scaffolding from [template-course-from-scratch](https://github.com/pavelkushneryk/template-course-from-scratch)
- Structural inspiration from the other `*-from-scratch` courses in this workspace

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — see [LICENSE](LICENSE).
