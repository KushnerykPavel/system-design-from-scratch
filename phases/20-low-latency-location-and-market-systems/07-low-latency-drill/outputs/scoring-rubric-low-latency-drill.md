# Scoring Rubric: Low-Latency Drill

Score each area from 1 to 4.

## Scope and Clarify

- 1: Jumps to components without a user-visible promise
- 2: Names a promise but leaves scope fuzzy
- 3: Clarifies the hot path and one key trade-off
- 4: Clarifies the hot path, freshness or exactness boundary, and likely failure pressure

## Sizing

- 1: Uses only vague scale words
- 2: Gives averages but no skew or burst estimate
- 3: Sizes main load and one hotspot
- 4: Sizes main load, hotspot amplification, and the first likely bottleneck

## Architecture

- 1: Lists components only
- 2: Has a diagram but no sync versus async boundary
- 3: Clear high-level design with one credible deep dive
- 4: Strong boundaries, hotspot story, and operationally realistic flow

## Operations

- 1: Says "add retries and replicas"
- 2: Mentions monitoring without specific signals
- 3: Names failure modes and detection metrics
- 4: Includes degraded mode, incident signals, and redesign logic
