---
lesson: 04-hot-and-cold-data
focus: balanced
---

## Clarify first
- Is the hotspot read-heavy, write-heavy, or both?
- What latency matters for recent data versus history?
- Can archive retrieval be asynchronous?

## Core design
- Keep the working set on the fast path
- Tier cold history deliberately
- Name the restore workflow, not just the archive tier

## Failure probes
- What if one tenant becomes half of all traffic?
- What if archive restore has never been tested?
- What if tiering moves data too early?
