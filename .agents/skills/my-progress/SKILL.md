---
name: my-progress
version: 1.0.0
description: Personal dashboard for System Design from Scratch. Shows placement, quiz results, likely weak spots, and the next lesson to study.
---

# My Progress

Use `.progress.json` and `ROADMAP.md` to show:
- placement phase
- lessons completed
- assisted vs unassisted completion ratio
- hours done vs remaining
- recent quiz performance
- due reviews
- overdue reviews
- one weak area trend
- one strongest area trend
- the next lesson to study
- one improving dimension trend from `feedback_history`
- one explicit recommendation with a reason, preferring overdue reviews before new lessons
- the top recurring `mistake_tag` this week and the suggested corrective drill
- when possible, recommend a short drill mode such as `capacity_drill`, `tradeoff_drill`, `failure_drill`, `clarification_drill`, or `communication_drill` instead of a full lesson replay
- consistency signals such as current streak and active days in the last week
- a compact weekly summary of sessions, lessons touched, and review backlog

When `feedback_history` exists, summarize weak and improving patterns using the shared dimensions:
- `clarification`
- `requirements`
- `sizing`
- `architecture`
- `deep_dive`
- `failure_modes`
- `observability`
- `trade_offs`
- `communication`
