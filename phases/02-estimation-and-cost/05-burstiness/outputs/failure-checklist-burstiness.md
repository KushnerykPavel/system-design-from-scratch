# Failure Checklist — Burstiness

- Peak factor stated explicitly, not hidden in averages.
- Arrival and service rates compared on the same time basis.
- Recovery time computed after the burst, not assumed away.
- Retry amplification considered.
- Hot partition risk called out separately from fleet average.
- Queue expiry or useless stale work discussed.
