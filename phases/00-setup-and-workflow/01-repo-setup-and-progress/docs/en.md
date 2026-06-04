# Repo Setup and Progress Tracking

> Your prep system should reduce friction, not create another project to maintain.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Set up a lightweight practice environment that records progress, preserves artifacts, and makes it obvious what lesson to do next.
**Prerequisites:** none
**Estimated time:** ~45 min
**Primary artifact:** progress checklist + repository validator

## The Problem

Many strong engineers sabotage their prep before the first design prompt. Notes are scattered, practice artifacts disappear, and there is no reliable signal for whether they are improving or just staying busy.

This lesson gives you a minimal prep control plane. The objective is not a fancy dashboard. The objective is a workflow you will still use after twenty lessons.

## Clarify

- What signals will prove you are improving: finished lessons, quiz scores, mock scores, or recurring mistake categories?
- How often will you study each week, and how much of that time is solo practice versus mock interviews?
- If you skip manual updates under pressure, what is the smallest metadata set worth preserving anyway?

## Requirements

### Functional

- Track lesson completion status and last touched date.
- Preserve links to notes, artifacts, and mock outcomes.
- Make the next recommended lesson obvious.

### Non-functional

- Update cost must stay under two minutes per session.
- Data model must survive partial or inconsistent updates.
- The workflow should be simple enough to keep using for months.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Lessons tracked | 180 planned lessons | defines progress file shape and status granularity |
| Session writes | 1 to 3 updates per study session | favors append-friendly or simple overwrite flows |
| Artifact references | tens to low hundreds | argues for storing paths, not duplicating content |
| Peak factor | short bursts after mocks | update flow should tolerate multiple edits in one evening |
| Rough cost | effectively zero beyond local files | do not overbuild this as an app |

## Architecture

Use a tiny three-part setup:

1. A canonical roadmap for lesson order.
2. A machine-readable progress file for status and timestamps.
3. Reusable artifacts per lesson stored next to the lesson itself.

The main idea is separation of concerns:

- `ROADMAP.md` defines what exists.
- `.progress.json` defines your current state.
- lesson folders hold durable outputs you may revisit later.

The code artifact for this lesson validates a progress snapshot so your prep system fails loudly when required metadata is missing.

## Data Model & APIs

Suggested record fields:

- `schema_version`
- `lesson`
- `status` such as `not_started`, `in_progress`, `done`, `assisted`
- `last_updated`
- `quiz_score`
- `quiz_history`
- `confidence` such as `low`, `medium`, `high`
- `mistake_tags`
- `notes_path`
- `artifact_paths`
- `last_reviewed_at`
- `next_review_at`
- `review_interval_days`
- `review_ease`
- `lapse_count`
- `mode_history`
- `feedback_history`

Minimal API surface:

- `validate(progress)` checks required fields and legal statuses.
- `recommendNext(progress, roadmap)` returns the next unfinished lesson.
- `syncReviewSchedule(progress, today)` derives `last_reviewed_at`, `next_review_at`, `review_interval_days`, `review_ease`, and `lapse_count` from the latest lesson evidence.
- `recommendReviews(progress, today)` returns due and overdue review lessons.

Suggested `feedback_history` shape per session:

- `session_type` such as `lesson`, `check_understanding`, `design_review`, `mock_interview`
- `completed_at`
- `summary`
- `strengths`
- `gaps`
- `highest_leverage_improvement`
- `dimensions[]` where each item includes:
  - `dimension` from `clarification`, `requirements`, `sizing`, `architecture`, `deep_dive`, `failure_modes`, `observability`, `trade_offs`, `communication`
  - `score` from 1 to 4
  - `evidence`
  - `next_action`

Keep the compatibility rule simple:

- `schema_version` defaults to the legacy shape when absent.
- newer versions should fail loudly with a migration hint instead of silently dropping fields.
- review and quiz history fields may start optional so the workflow stays low-friction early on.

Suggested scheduling rules:

- `assisted` or quiz score `0-5`: schedule the next review in `1` day and increment `lapse_count`
- quiz score `6-7`: schedule the next review in `3-7` days depending on prior interval and weak dimensions
- quiz score `8`: schedule the next review in `7-14` days, longer when confidence is high and weak dimensions are absent
- repeated success should multiply the interval by `review_ease`
- repeated misses should reset the interval to `1`

Controlled `mistake_tags`:

- `skips_clarification`
- `weak_requirements`
- `no_sizing`
- `component_soup`
- `bad_deep_dive_choice`
- `weak_consistency_reasoning`
- `shallow_failure_modes`
- `missing_observability`
- `weak_tradeoffs`
- `weak_communication`
- `rushed_architecture`
- `over_indexes_storage`
- `no_operational_story`
- `does_not_tie_back_to_requirements`

Use those tags to answer:

- which recurring mistakes should raise review priority?
- which short drill best addresses the top repeated mistake this week?
- whether a lesson should be revisited because the same mistake appears across prompts

Suggested short drill modes:

- `capacity_drill` for `no_sizing` or weak `sizing`
- `tradeoff_drill` for `weak_tradeoffs` or `weak_consistency_reasoning`
- `failure_drill` for `missing_observability`, `shallow_failure_modes`, or `no_operational_story`
- `clarification_drill` for `skips_clarification` or `rushed_architecture`
- `communication_drill` for `weak_communication`
- `redesign_drill` for `component_soup`, `bad_deep_dive_choice`, or `over_indexes_storage`

Useful dashboard signals:

- completion breakdown: `done`, `assisted`, `in_progress`, `not_started`
- assisted ratio among completed lessons
- current streak and active days in the last 7 days
- due today vs overdue review backlog
- strongest and weakest dimensions by recent feedback averages
- weekly summary of sessions, lessons touched, and review pressure

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| progress file becomes stale | lessons marked complete but no recent timestamps | require `last_updated` on every write |
| prep artifacts are lost | referenced note path does not exist | validate artifact and notes paths regularly |
| too much metadata to maintain | skipped updates after sessions | keep required fields minimal and default the rest |
| false sense of progress | many completed lessons but no mock evidence | pair lesson status with quiz and mock signals |

## Observability

- metric: lessons completed per week
- metric: average days between touching the same weak area
- metric: percentage of lessons with notes and quiz data attached
- SLO: after any session, progress state should be updated within 24 hours

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| local JSON progress file | trivial to edit and validate | no built-in UI or history analysis | custom app or spreadsheet |
| lesson-local artifacts | keeps context near source material | more files in the repo | centralized notes bucket |
| minimal metadata | low maintenance | less rich analytics | detailed journaling every session |

## Interview It

**Google framing:** "How would you design your own interview prep workflow so you can detect whether you are improving rather than just repeating exercises?"

**Cloudflare framing:** "How would you structure a reusable practice system for edge and platform interviews where lesson artifacts, debriefs, and failure drills need to stay searchable?"

**Follow-ups:**
1. What changes if two people share the same repo for pair practice?
2. How do you distinguish `done` from `assisted` without shaming the learner?
3. What if the progress file is edited manually and drifts from reality?
4. How would you surface stale topics that have not been revisited in 30 days?
5. What is the minimum rollout if you were starting tonight?

## Ship It

- `outputs/interview-card-repo-setup-and-progress.md`
- `outputs/progress-checklist.md`

## Exercises

1. **Easy** — Add a new status that marks a lesson as blocked on a missing prerequisite.
2. **Medium** — Design a weekly review summary derived from the same progress file.
3. **Hard** — Extend the workflow so mock interview scores influence the next recommended lesson.

## Further Reading

- [GitHub Engineering Career Growth](https://github.blog/engineering/) — useful mindset for tracking improvement through evidence rather than vibes
- [liquidslr/system-design-notes](https://github.com/liquidslr/system-design-notes) — reference for keeping system design practice structured and repeatable
