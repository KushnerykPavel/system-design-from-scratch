# Proximity Service Design Review Sheet

Use this when reviewing a nearby-search answer.

## Must-Have Decisions

- What objects are moving, and what freshness target matters?
- What is the approximate geo partitioning strategy?
- Where does exact distance filtering happen?
- How are availability and business filters applied?
- What is the fallback when local supply is sparse?

## Hotspot Review

- What happens in a downtown or stadium hotspot?
- How many candidates can one query expand to?
- Is there a hot-cell split or replication plan?
- Which metrics reveal skew before latency collapses?

## Failure Review

- How are stale or missing updates detected?
- What happens if updates arrive out of order?
- How does the system degrade when the fresh-update pipeline lags?
- What operator knob exists for widening-radius safety?

## Senior-Level Follow-Up

- What changes at 10x read load?
- What changes at 10x write load?
- Which guarantee would you relax first if cost pressure rises?
