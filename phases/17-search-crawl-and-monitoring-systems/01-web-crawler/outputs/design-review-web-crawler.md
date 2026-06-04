# Design Review Prompt — Web Crawler

Use this when reviewing a crawler design or interview answer.

## Clarify
- What kinds of pages deserve the freshest recrawl target?
- Do we need to respect robots.txt and per-host politeness?
- Is the product optimizing for search coverage, monitoring freshness, or archival completeness?

## Core checks
- The frontier is partitioned in a way that preserves per-host fairness.
- Work claims are leased so failed fetchers do not strand URLs.
- URL dedup and content dedup are treated as separate concerns.
- Recrawl cadence is budgeted by value and change rate, not by one global timer.

## Failure probes
- What happens when one host produces 20% of discovered URLs?
- How are stale robots policies detected and corrected?
- What prevents parser bugs from exploding the frontier with duplicates?
- How do you explain why a specific URL was delayed or skipped?
