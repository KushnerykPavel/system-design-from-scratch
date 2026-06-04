# Lightning Trade-off Rounds

> The fastest way to sound senior is not to avoid trade-offs, but to say the painful one out loud before the interviewer asks.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Practice short trade-off drills that force clear decision framing, rejected alternatives, and measurable consequences.  
**Prerequisites:** `03-design-framework-and-timing/05-wrap-up`, `21-google-senior-staff-system-design/05-risk-communication`, `23-mixed-mocks-and-redesign-drills/05-lightning-capacity-rounds`  
**Estimated time:** ~60 min  
**Primary artifact:** trade-off flashcards + wrap-up worksheet  

## The Problem

You have two to four minutes to answer prompts like:

- read-through cache or write-through cache?
- single-region primary or active-active?
- queue or workflow engine?
- strict consistency or bounded staleness?
- direct fanout or pull-based delivery?

The point is not to list pros and cons endlessly. The point is to decide, justify, and name what you are paying.

## Clarify

- What user or operator outcome matters most for this decision?
- Which resource is scarce: latency budget, correctness budget, cost budget, or operator time?
- Is the decision for all traffic or just one path or tenant class?
- What failure mode makes one option riskier?

## Requirements

### Functional

- Choose an option quickly.
- Tie the decision to one or two concrete requirements.
- Name one rejected alternative and why it lost.
- Explain one metric or failure signal affected by the choice.

### Non-functional

- Avoid buzzword-only comparisons.
- Keep the trade-off connected to workload shape.
- Make the downside explicit.
- Stay concise enough for real interview pacing.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Time per round | 2 to 4 min | enforces clarity |
| Requirements named | 1 to 2 | enough to justify a decision |
| Alternatives considered | 2 or 3 | keeps reasoning bounded |
| Metrics named | at least 1 | proves the trade-off is measurable |
| Prompt count | 5 to 10 per session | builds repetition without fatigue |

## Architecture

Trade-off loop:

```text
state the decision
  -> name the winning requirement
  -> choose the option
  -> say what it costs
  -> say what metric or failure mode will reveal if it was wrong
```

Useful prompt families:

1. latency versus correctness
2. throughput versus operator simplicity
3. cost versus freshness
4. isolation versus utilization
5. rollout speed versus safety

## Data Model & APIs

Trade-off template:

```text
tradeoff_round(
  decision,
  top_requirement,
  chosen_option,
  rejected_option,
  downside,
  validating_metric
)
```

Questions to keep nearby:

- what requirement wins?
- what gets worse?
- what signal proves the decision was wrong?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| answer never chooses | many pros and cons but no decision | force a winner within one minute |
| downside is hidden | answer sounds unrealistic | require one explicit cost |
| no metric is named | trade-off stays abstract | attach one validation signal |
| decision ignores workload shape | same answer used everywhere | restate the dominant requirement first |

## Observability

- metric: percentage of rounds with a clear decision and explicit downside
- metric: percentage of rounds that name a validating metric or failure mode
- metric: prompt family diversity across practice
- log: chosen requirement and rejected alternative per round
- trace: requirement -> decision -> downside -> metric
- SLO: produce a credible trade-off answer in under four minutes

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| forced early choice | stronger communication | less room for broad exploration | endless comparison |
| one validating metric | makes the trade measurable | can oversimplify a complex system | purely verbal justification |
| short rounds | better interview pacing | less detail | long case-study format |

## Interview It

**Google framing:** Expect the interviewer to care whether your choice follows from stated priorities rather than from memorized best practices.

**Cloudflare framing:** Expect more pressure on operational failure, propagation safety, routing behavior, and cost under partial outages.

**Suggested rounds:**
1. Active-active or active-passive for a write-heavy control plane?
2. Push fanout or pull fanout for a high-skew feed?
3. Fail-open or fail-closed for policy evaluation?
4. Global strong consistency or bounded staleness for metadata reads?

## Ship It

- `outputs/tradeoff-flashcards-lightning-rounds.md`
- `outputs/wrapup-worksheet-lightning-tradeoffs.md`

## Exercises

1. **Easy** — Do three rounds and force yourself to choose within 60 seconds.
2. **Medium** — Add one metric to each answer that would prove the choice was wrong.
3. **Hard** — Re-answer the same trade-off twice with different top requirements and compare the outcome.

## Further Reading

- [Google SRE](https://sre.google/books/) — useful examples of explicit operational trade-offs  
- [Designing Data-Intensive Applications](https://dataintensive.net/) — excellent source for concrete distributed-systems trade-offs  
