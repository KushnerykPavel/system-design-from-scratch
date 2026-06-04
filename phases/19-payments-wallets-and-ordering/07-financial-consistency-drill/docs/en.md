# Financial Consistency Drill

> The point of this drill is not to name components; it is to defend one correctness boundary under pressure.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Practice a full interview loop on financial consistency problems by forcing explicit clarification, sizing, invariants, failure handling, and redesign under changed constraints.
**Prerequisites:** `03-design-framework-and-timing/08-full-loop-drill`, `19-payments-wallets-and-ordering/01-payment-ledger`, `19-payments-wallets-and-ordering/06-fraud-hooks`
**Estimated time:** ~60 min
**Primary artifact:** drill scorecard validator + practice sheet

## The Problem

Run a timed interview drill for prompts such as "design a payment ledger," "design a wallet with holds," "design inventory reservation for flash sales," or "design order recovery after partial payment success." Your goal is to produce a coherent correctness story quickly.

This lesson exists to turn Phase 19 into interview reflex. You should clarify the money or stock invariant, size retry amplification, choose one deep dive, and still leave time for observability, operator recovery, and redesign.

## Clarify

- What is the core invariant: no money loss, no oversell, no double-spend, or no ambiguous recovery?
- Which subsystem is the source of truth?
- Which retries or callbacks are most likely to duplicate work?
- Where should the deep dive go if time is limited?

If the prompt stays broad, choose one invariant and organize the answer around that boundary instead of trying to solve every financial subsystem at once.

## Requirements

### Functional

- Clarify scope and invariant before naming infrastructure.
- Size steady-state load and incident-time amplification.
- Present a high-level design with one deliberate deep dive.
- Explain failure handling, observability, and operator workflows.
- Redesign after a changed requirement without losing correctness.

### Non-functional

- Stay concrete under ambiguity.
- Avoid pretending exactly-once or global transactions are free.
- Show where auditability and compliance influence the design.
- Keep the design operationally credible.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Clarification time | 3 to 5 minutes | sets the invariant and scope early |
| Sizing time | 5 to 7 minutes | retries and batch fanout must be quantified |
| Initial design time | 10 to 12 minutes | enough to establish trust before the deep dive |
| Deep dive time | 8 to 10 minutes | where senior judgment becomes visible |
| Redesign time | 5 minutes | tests whether the design has real constraints |

## Architecture

Use the same loop every time:

```text
clarify invariant and source of truth
  -> size steady-state and retry amplification
  -> propose architecture and guarantee boundaries
  -> deep dive one critical subsystem
  -> close with failure handling, observability, and redesign
```

Good deep-dive choices for this phase:

- ledger posting invariants and reversal model
- hold lifecycle and overspend protection
- hot-SKU reservation control
- ambiguous order recovery and compensation
- fraud-gate degradation policy

## Data Model & APIs

For the drill, your answer structure is the main interface:

- invariant statement
- source-of-truth ownership
- idempotency boundary
- one core record or state machine
- observability and operator controls
- redesign trigger and response

If you cannot name the authoritative record or transition model, the answer is probably still too abstract.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| candidate jumps into components before invariant | unclear guarantee boundary in first minutes | force a one-sentence correctness statement |
| retries are ignored in sizing | architecture fails under incident amplification | add duplicate and backlog factors to estimates |
| answer has no operator recovery | stuck-state handling is hand-waved | require a recovery queue or runbook |
| redesign preserves all guarantees magically | trade-offs stay hidden | state exactly what becomes slower, costlier, or less strict |

## Observability

- metric: did the answer define the invariant and source of truth clearly?
- metric: did the sizing include retries, holds, or reconciliation amplification?
- metric: did each failure mode name a detection signal?
- log: note where assumptions changed during the drill
- trace: scope -> design -> deep dive -> redesign
- SLO: produce a coherent, defensible financial-consistency answer inside interview time

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| one explicit invariant | focused, defensible answer | less breadth elsewhere | vague "consistency is important" answer |
| one strong deep dive | shows judgment and depth | less time for edge topics | shallow coverage of everything |
| honest bounded guarantees | operational credibility | fewer flashy claims | magical exactly-once story |

## Interview It

**Google framing:** "Design a payment or ordering subsystem and explain how correctness survives retries and partial failure." Expect pushback on invariants, compensation, and operational debugging.

**Cloudflare framing:** "Design a global billing or ordering platform where shared systems cannot be destabilized by one tenant or one incident." Expect pressure on control-plane safety, auditability, and graceful degradation.

**Follow-ups:**
1. What if traffic grows 10x during a launch?
2. What if regulators require stronger audit retention?
3. What if the source of truth must become multi-region?
4. What if fraud review must become synchronous for premium tenants?
5. What if cost becomes more important than freshness?

## Ship It

- `outputs/skill-financial-consistency-drill.md`

## Exercises

1. **Easy** — Give a three-minute opening for "design a payment ledger."
2. **Medium** — Pick the best deep dive for "design a wallet with holds" and justify it.
3. **Hard** — Redesign an order-recovery answer after the interviewer adds global active-active writes and stricter audit retention.

## Further Reading

- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) — good source of operational follow-up pressure
- [liquidslr/system-design-notes](https://github.com/liquidslr/system-design-notes) — useful reminder of the four-step interview rhythm
