# Redesign After a 10x Scale Change

> Senior answers get stronger when scale changes force them to revisit assumptions instead of merely enlarging the same diagram.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Practice redesigning a solid baseline system after a 10x workload shift changes hotspots, economics, and operational risk.  
**Prerequisites:** `02-estimation-and-cost/08-uncertain-inputs`, `03-design-framework-and-timing/06-redesign-prompts`, `20-low-latency-location-and-market-systems/07-low-latency-drill`  
**Estimated time:** ~75 min  
**Primary artifact:** ten-x redesign checklist + trade-off matrix  

## The Problem

You already have a reasonable design. The interviewer now says one important dimension grew by 10x:

- user count
- write QPS
- peak burst factor
- object size
- regional spread
- retention period

This lesson is about redesign discipline. The right response is rarely "add more servers." A 10x shift usually changes at least one of the following:

- the dominant bottleneck
- the most important cost term
- the best consistency boundary
- the safest work placement
- the operational failure mode

## Clarify

- Which exact dimension grew by 10x, and did anything else stay fixed?
- Is the new requirement steady-state, burst-only, or premium-path only?
- Which guarantee is still sacred after the scale jump?
- Is the redesign allowed to be incremental, or can it assume a greenfield rewrite?

## Requirements

### Functional

- Identify which original assumptions are now invalid.
- Re-estimate the hot path with enough numbers to justify change.
- Propose one or two architecture changes that target the new bottleneck.
- Explain rollout, migration, and observability for the redesign.

### Non-functional

- Avoid rewriting the whole system unless necessary.
- Keep cost and operational complexity visible.
- Preserve the most important user-visible guarantee deliberately.
- Show which trade-offs worsened as scale increased.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Scale jump | 10x on one chosen axis | determines what actually changed |
| Headroom before redesign | often less than 2x | explains why the old design is no longer safe |
| New bottleneck | CPU, fanout, network, storage, or coordination | focuses the redesign |
| Migration duration | days to months depending on data movement | changes rollout and dual-run cost |
| Cost multiplier | may grow more or less than 10x | exposes hidden economic risk |

## Architecture

Recommended redesign loop:

```text
state original design
  -> identify the broken assumption
  -> re-size the new bottleneck
  -> choose the smallest meaningful architecture change
  -> explain migration and rollback
  -> name the new failure mode and metric
```

Good redesign moves include:

1. repartitioning or reshaping keys
2. pushing work async or precomputing more state
3. splitting hot and cold paths
4. introducing regionalization or tenancy tiers
5. tightening or relaxing consistency on the right path

## Data Model & APIs

Useful redesign template:

```text
ten_x_redesign(
  changed_dimension,
  broken_assumption,
  new_bottleneck,
  architecture_change,
  migration_plan,
  new_metrics
)
```

Key questions:

- which API or state boundary now feels too expensive?
- which data shape is creating hotspots?
- which client-visible contract must stay stable through migration?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| redesign is only "add more boxes" | no bottleneck or cost term changed | force one broken assumption and one bottleneck estimate |
| redesign replaces everything | migration and rollback are absent | prefer smallest high-leverage change first |
| cost is ignored | new capacity is feasible technically but not economically | estimate at least one major cost multiplier |
| new failure mode is missed | the scaled design creates different incidents | add one new metric and one new degraded mode |

## Observability

- metric: the new bottleneck indicator that justified redesign
- metric: migration progress and dual-path correctness checks
- metric: cost or resource efficiency before and after the change
- log: assumptions invalidated by the scale jump
- trace: old path versus new path during shadow or canary rollout
- SLO: preserve one user-visible guarantee while scaling the chosen dimension

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| smallest meaningful redesign | lower migration risk | may leave future headroom limited | full rewrite on first pressure signal |
| explicit dual-run or shadowing | safer migration | temporary cost increase | one-shot cutover |
| bottleneck-specific optimization | targeted gains | more heterogeneous architecture | uniform scaling everywhere |

## Interview It

**Google framing:** Expect the interviewer to test whether your redesign follows from new math rather than from habit.

**Cloudflare framing:** Expect pressure on how the redesign affects routing, cacheability, propagation, or regional blast radius as scale changes.

**Follow-ups:**
1. Which old assumption broke first?
2. Which cost term now dominates?
3. What migration risk worries you most?
4. What new metric would you page on?
5. What gets worse because of the redesign?

## Ship It

- `outputs/redesign-checklist-ten-x.md`
- `outputs/tradeoff-matrix-ten-x-redesign.md`

## Exercises

1. **Easy** — Take one earlier lesson and name the first bottleneck that breaks at 10x.
2. **Medium** — Propose one incremental redesign and one migration guardrail.
3. **Hard** — Redesign the same system after the scale jump happens only in one region with a 20x peak factor.

## Further Reading

- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) — good for thinking about growth, change management, and operational safety  
- [Designing Data-Intensive Applications](https://dataintensive.net/) — strong reference for changing bottlenecks under scale  
