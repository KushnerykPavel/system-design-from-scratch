# Freshness Decision Checklist

- Name the data class before naming the TTL.
- State the acceptable stale window in user-visible terms.
- Decide whether TTL alone is enough or whether explicit invalidation is required.
- Keep a safety backstop when relying on invalidation events.
- Include version or timestamp metadata when stale behavior must be debugged.
- Ask whether all fields in the object really need the same freshness target.
- Measure invalidation latency and version skew, not only hit ratio.
