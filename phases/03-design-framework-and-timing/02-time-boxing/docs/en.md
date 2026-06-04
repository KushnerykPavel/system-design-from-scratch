# How to Spend 45 Minutes Wisely

> Strong candidates budget time before the interviewer budgets it for them.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Learn to allocate interview time deliberately so clarification, sizing, architecture, deep dive, and wrap-up all get airtime.  
**Prerequisites:** `03-design-framework-and-timing/01-four-step-interview-loop`  
**Estimated time:** ~60 min  
**Primary artifact:** time-box worksheet + interview card  

## The Problem

Many system design interviews go wrong because the candidate treats time as infinite. They spend 20 minutes on a diagram, then rush the deep dive, skip observability, and never summarize trade-offs.

This lesson turns pacing into an explicit design decision. You are not only designing the system. You are designing how the conversation unfolds.

## Clarify

- Is this a 35-minute, 45-minute, or 60-minute design round?
- Does the interviewer want broad architecture first or is there a known hotspot to deep dive early?
- Are we optimizing for product reasoning, infrastructure realism, or a specific company style?
- Should I assume time for redesign follow-ups and failure discussion at the end?

If the interviewer gives no pacing signals, default to a five-stage plan and narrate your checkpoints.

## Requirements

### Functional

- Reach clarified scope before proposing architecture.
- Produce at least rough sizing before major component choices.
- Leave enough time for one meaningful deep dive and a wrap-up.

### Non-functional

- Maintain conversation control without sounding robotic.
- Adapt the plan if the interviewer interrupts or changes constraints.
- Preserve time for trade-offs, observability, and risks.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Interview length | 45 min default | determines how much depth fits |
| Clarification budget | 5-7 min | prevents drift into the wrong problem |
| Sizing budget | 4-6 min | anchors topology and cost |
| Deep dive budget | 12-15 min | shows prioritization and senior judgment |
| Wrap-up reserve | 5-7 min | protects failure modes and trade-offs |

## Architecture

A default 45-minute pacing plan:

```text
00-06  clarify scope and priorities
06-12  size traffic, storage, and peaks
12-24  propose high-level design and get buy-in
24-38  deep dive on the main risk
38-45  failure modes, observability, trade-offs, wrap-up
```

The important property is not the exact minute split. It is that you deliberately reserve time for the later parts that weaker answers never reach.

## Data Model & APIs

Treat your interview plan like a sequence of stages:

- `clarify`
- `size`
- `high_level_design`
- `deep_dive`
- `wrap_up`

The code artifact models those stages and flags when you drift beyond the time budget.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Architecture starts before scope is clear | no assumptions or priorities stated | pause and restate the problem in one sentence |
| Sizing gets skipped | component choices have no numbers attached | force one QPS and one storage estimate before continuing |
| Deep dive expands uncontrollably | more than one subsystem starts competing for time | choose the critical path and defer side topics |
| Wrap-up disappears | less than 3 minutes remain unexpectedly | start summarizing risks and trade-offs immediately |

## Observability

- metric: minutes spent in each stage
- metric: number of explicit assumption checks
- metric: number of times the plan was rebalanced mid-interview
- alert: no wrap-up reached in practice drills

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| fixed stage budget | keeps the answer complete | can feel scripted | unstructured pacing that often overruns |
| reserved wrap-up time | surfaces risks and maturity | slightly less architecture detail | using every minute on core design only |
| one primary deep dive | shows prioritization | reduces breadth | shallow coverage of several topics |

## Interview It

**Google framing:** "Design a notification system in 45 minutes." The interviewer wants to see whether you can structure the answer and still pivot when they push into reliability or delivery guarantees.

**Cloudflare framing:** "Design an edge configuration rollout service." The interviewer wants evidence that you can budget time for propagation risk, rollback, and observability rather than talking only about control-plane APIs.

**Follow-ups:**
1. How would your time plan change in a 30-minute screen?
2. When is it correct to shorten clarification and go deeper on design?
3. How do you recover if the interviewer pulls you into a low-value detail early?
4. What should you cut first if you are running late?
5. How does time-boxing change when the system is operationally heavy?

## Ship It

- `outputs/interview-card-time-boxing.md`
- `outputs/time-boxing-worksheet.md`

## Exercises

1. **Easy** — Build a 20-minute pacing plan for a rate limiter prompt.  
2. **Medium** — Reallocate the default budget for a compliance-heavy audit log system.  
3. **Hard** — Design a pacing plan for a 60-minute senior interview where the interviewer will force a redesign halfway through.  

## Further Reading

- [System design notes: design a system in 4 steps](https://github.com/liquidslr/system-design-notes) — classic pacing baseline  
- [Google SRE book](https://sre.google/books/) — strong reference for reserving time for reliability trade-offs  
