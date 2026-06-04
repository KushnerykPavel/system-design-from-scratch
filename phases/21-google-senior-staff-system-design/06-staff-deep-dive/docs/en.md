# Staff-Level Deep-Dive Selection

> Staff-level depth is not "more detail everywhere." It is picking the one subsystem where better judgment most changes the design.

**Type:** Learn
**Company focus:** Google
**Learning goal:** Learn how to choose and justify a deep dive that reveals staff-level judgment instead of accidental implementation trivia.
**Prerequisites:** `03-design-framework-and-timing/04-deep-dive-selection`, `11-observability-slos-and-debugging/08-observability-drill`, `19-payments-wallets-and-ordering/07-financial-consistency-drill`
**Estimated time:** ~60 min
**Primary artifact:** deep-dive selector worksheet

## The Problem

By the time a candidate is operating at senior or staff level, the interviewer often expects more than a clean high-level design. They want to see whether you can choose the highest-leverage subsystem and unpack guarantees, bottlenecks, failure behavior, and operational consequences there.

Weak answers deep dive into the most familiar area. Strong answers deep dive into the area that makes the rest of the architecture credible.

## Clarify

- Which subsystem carries the hardest guarantee or biggest scaling risk?
- Where could one wrong assumption invalidate the whole design?
- Which part is most likely to trigger interviewer pushback?
- Which deep dive would distinguish staff-level reasoning from senior-level completeness?

If multiple deep dives are possible, pick the one with the largest combination of product risk, operational complexity, and architectural leverage.

## Requirements

### Functional

- Identify 2 to 3 plausible deep-dive targets.
- Choose one and justify it explicitly.
- Explain the chosen subsystem's guarantees, interfaces, bottlenecks, and failure modes.
- Connect the deep dive back to system-wide trade-offs.
- Defend why other deep dives were not first priority.

### Non-functional

- Avoid deep-diving into a component only because it is personally familiar.
- Avoid low-value implementation detail with no design consequence.
- Keep the deep dive connected to product and operational outcomes.
- Preserve time for wrap-up and redesign.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Candidate deep-dive areas | 2 to 4 | too many means the architecture is still fuzzy |
| Time available for depth | 8 to 12 min | enough for one strong subsystem, not three |
| Major failure modes in the dive | 2 to 4 | depth should include degraded behavior, not only the happy path |
| System-wide consequences | at least 2 | proves the deep dive matters beyond itself |

## Architecture

Use a selection rubric:

```text
deep_dive_score =
  design leverage
  + correctness risk
  + scaling risk
  + operational risk
  - implementation trivia
```

Common good deep dives:

1. Write-path correctness and idempotency.
2. Partitioning and hotspot control.
3. Cache invalidation and freshness boundaries.
4. Multi-region failover and recovery semantics.
5. Control-plane rollout safety.

Common weak deep dives:

1. API field naming unless it changes the architecture materially.
2. Database internals unrelated to the prompt's main tension.
3. Generic auth boilerplate in prompts where data movement is the real risk.

## Data Model & APIs

Useful structure:

```text
deep_dive_candidate(
  component,
  reason,
  primary_risk,
  system_impact
)
```

Helpful verbal APIs:

- `ListCandidates()`
- `PickDeepDive(component, reason)`
- `ExplainGuarantee(component)`
- `ConnectBackToArchitecture()`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| deepest discussion is not the most important one | deep dive never changes trust in the design | rank candidates by leverage before choosing |
| deep dive becomes implementation trivia | no system-level consequence is named | tie every detail back to a requirement or risk |
| candidate cannot explain why other dives were deferred | choice seems accidental | name one rejected dive and why it matters less now |
| deep dive consumes the whole interview | no wrap-up or redesign remains | time-box the depth explicitly |

## Observability

- metric: does the chosen deep dive target the system's hardest risk?
- metric: were guarantees, bottlenecks, and failure modes all covered?
- metric: were trade-offs from the deep dive reflected back into the full design?
- log: rejected deep-dive candidates and reasons
- trace: candidate list -> chosen subsystem -> consequences -> redesign
- SLO: choose one deep dive that most increases confidence in the architecture

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| one high-leverage deep dive | shows judgment and depth | less breadth elsewhere | many shallow mini-dives |
| explicit candidate ranking | clearer choice rationale | a little more upfront time | accidental selection |
| connect depth back to system | preserves coherence | requires discipline | isolated subsystem monologue |

## Interview It

**Google framing:** "What part would you dive into next?" Staff-level answers sound intentional here. They say why that subsystem matters more than the others, and what risk it reveals.

**Follow-ups:**
1. Why not deep dive the storage layer instead?
2. What if the interviewer pushes into a different subsystem?
3. How do you know your chosen area is actually the highest risk?
4. What changes if scale grows 10x but the product remains the same?

## Ship It

- `outputs/deep-dive-selector-google-staff.md`

## Exercises

1. **Easy** - Pick the best deep dive for a notification system.
2. **Medium** - Compare deep-diving partitioning versus deep-diving consistency for a metadata service.
3. **Hard** - Justify a staff-level deep dive for a globally distributed API gateway where both rollout safety and origin protection matter.

## Further Reading

- [Google SRE book](https://sre.google/sre-book/table-of-contents/) - useful for identifying subsystems where operational risk dominates
- [Designing Data-Intensive Applications](https://dataintensive.net/) - helpful background for choosing depth around correctness, replication, and scale
