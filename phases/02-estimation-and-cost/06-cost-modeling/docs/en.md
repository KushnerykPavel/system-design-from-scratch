# Rough Cost Modeling Without Getting Lost

> Cost modeling in interviews is not about perfect billing. It is about proving you notice which choice gets expensive first.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Build quick cost models that compare architectural options without drowning in cloud-pricing trivia.  
**Prerequisites:** `02-estimation-and-cost/02-storage-growth`, `02-estimation-and-cost/03-bandwidth-and-egress`  
**Estimated time:** ~75 min  
**Primary artifact:** cost worksheet + trade-off matrix  

## The Problem

Candidates often mention cost only at the end, and only in generic terms. Senior answers should identify the main spending levers early: compute, storage, bandwidth, replication, or managed-service premiums.

This lesson teaches quick component-based cost reasoning that supports decisions instead of stalling them.

## Clarify

- Are we comparing two architectures or estimating one design’s rough order of magnitude?
- Which line item is likely to dominate: compute, storage, or network?
- Does the interviewer care about developer velocity or only raw infrastructure spend?
- Are there multi-region or premium-SLA requirements that intentionally raise cost?

## Requirements

### Functional

- Break monthly spend into a few major buckets.
- Compare at least two architecture variants.
- Show which assumption dominates total cost sensitivity.

### Non-functional

- Keep the model rough and explainable.
- Avoid false precision from provider-specific pennies.
- Tie cost back to product or reliability value.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Stateless compute | 180 instances | major recurring line item |
| Monthly stored TB | 120 TB | storage baseline |
| Monthly egress | 8 PB | can dominate internet-facing services |
| Managed-service premium | +25% | changes buy-vs-build reasoning |
| Sensitivity driver | cache hit rate | can move both compute and egress cost |

## Architecture

A useful pattern:

1. List 3-5 main cost buckets.
2. Use round monthly units.
3. Compare a baseline and one optimization.
4. Name the assumption that swings the total most.

Example buckets:

- compute
- storage
- egress
- replicated cross-region traffic
- managed control plane premiums

## Data Model & APIs

The code artifact models monthly cost items:

```text
CostItem {
  Name
  UnitCost
  Units
}
```

Outputs:

- total monthly cost
- most expensive item

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| false precision | long pricing digression with low insight | round heavily and stay directional |
| hidden dominant line item | one expensive tier omitted | list major buckets explicitly |
| cost disconnected from value | cheapest design chosen blindly | tie savings to latency, reliability, or ops cost |
| sensitivity ignored | one assumption changes total massively | identify the swing factor clearly |

## Observability

- metric: spend by service and cost bucket
- metric: egress cost per request or per tenant
- metric: storage growth vs forecast
- metric: idle vs utilized compute capacity
- SLO: no single hidden cost bucket should exceed forecast without alerting

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| use round numbers | faster, clearer reasoning | less exact | provider-price deep dive |
| compare variants | supports decision-making | more setup | one-cost-number with no alternative |
| include managed-service premium | honest about build-vs-buy | rougher estimation | pretending control planes are free |

## Interview It

**Google framing:** "Design a log analytics system but keep cost sane." The signal is whether you know which lever to optimize first instead of listing all cloud services.

**Cloudflare framing:** "Design an edge product where origin egress and regional replication matter." The signal is whether you connect traffic shape to cost structure.

**Follow-ups:**
1. What if egress is 60% of spend?
2. What if managed storage reduces ops headcount but costs 30% more?
3. What if cache hit-rate improvement saves both bandwidth and database load?
4. What if multi-region is optional for v1?

## Ship It

- `outputs/cost-worksheet-rough-cost-modeling.md`
- `outputs/tradeoff-matrix-rough-cost-modeling.md`

## Exercises

1. **Easy** — Compare one-region vs two-region cost for a read-heavy API.  
2. **Medium** — Estimate whether better cache hit rate or storage compression saves more money.  
3. **Hard** — Justify paying more for a managed service because of operational risk reduction.  

## Further Reading

- [Google SRE book](https://sre.google/books/) — good context on capacity and operational trade-offs  
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline interview framework  
