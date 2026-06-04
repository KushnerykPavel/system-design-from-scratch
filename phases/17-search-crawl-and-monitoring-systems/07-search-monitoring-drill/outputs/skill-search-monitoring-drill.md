# Skill Sheet — Search and Monitoring Drill

Use this for a 45-minute self-run mock.

## Opening structure
1. Clarify the product, user, and top success metric.
2. State 3 to 5 assumptions and move forward.
3. Do rough QPS, storage, freshness, and retention math.
4. Propose a high-level design and get buy-in.
5. Pick one deep dive deliberately.

## Good deep-dive choices
- crawler frontier fairness
- autocomplete freshness and ranking snapshots
- metrics cardinality control
- alert routing quality and fallback
- index rollout safety

## Must-cover close-out points
- one important failure mode with detection and mitigation
- one observability plan tied to an SLO or freshness target
- one explicit trade-off with benefit and cost
- one redesign after a changed assumption

## Scoring prompts
- Did I clarify something that changed the design?
- Did I size before architecture?
- Did I choose one deep dive instead of many shallow ones?
- Did I explain what breaks and how I would notice?
- Did I adapt when the constraint changed?
