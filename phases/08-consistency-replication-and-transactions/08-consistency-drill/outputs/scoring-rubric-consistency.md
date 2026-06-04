# Scoring Rubric - Consistency Drill

## Strong answer signals

- Clarifies which flows actually need stronger guarantees.
- Maps guarantees by entity instead of giving one blanket rule.
- Chooses a replication model that matches the guarantee.
- Keeps transactions narrow and explains where sagas begin.
- Names failure behavior under lag, failover, and contention.
- Uses observability to prove the contract is working.

## Weak answer signals

- Says "eventual consistency" or "strong consistency" without product meaning.
- Treats replication choice as a vendor preference.
- Uses transactions everywhere without discussing hotspots.
- Mentions sagas without compensation or status visibility.
- Ignores clock, failover, or stale-read risk.
