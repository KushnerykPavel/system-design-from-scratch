# Messaging Platform Drill

> The hard part of messaging interviews is choosing one coherent guarantee set and defending it under changed constraints.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Practice full-loop messaging and job-platform interviews with explicit semantics, sizing, deep-dive choice, failure reasoning, and redesign under pressure.
**Prerequisites:** `03-design-framework-and-timing/08-full-loop-drill`, `18-messaging-and-job-platforms/01-distributed-message-queue`, `18-messaging-and-job-platforms/06-exactly-once-myths`
**Estimated time:** ~60 min
**Primary artifact:** drill scorecard validator + practice sheet

## The Problem

Run a timed drill for prompts like "design a distributed queue," "design a workflow engine," "design pub/sub for service events," or "design a scheduler for retries and cron jobs." The point is to produce a senior-level answer quickly, not a maximal list of technologies.

This lesson exists to turn the phase into interview muscle memory. You should clarify the guarantee boundary, size the workload, choose one deep dive, and still leave room for observability, failure handling, and redesign.

## Clarify

- Is the prompt centered on queues, workflows, scheduling, or pub/sub fanout?
- What guarantee matters most: ordering, durability, latency, replayability, or isolation?
- Is the biggest risk correctness, cost, or shared-platform abuse?
- Which component should get the deep dive if time is short?

If the prompt stays broad, make assumptions explicitly and choose the deepest pressure point rather than trying to cover everything.

## Requirements

### Functional

- Clarify scope and guarantee boundaries.
- Estimate throughput, storage, and retry or replay amplification.
- Present a high-level design and one deliberate deep dive.
- Explain failure handling, observability, and operator controls.
- Redesign after a changed requirement.

### Non-functional

- Stay concrete under ambiguity.
- Avoid vague exactly-once or "infinite replay" claims.
- Keep the design operationally credible.
- Show why one trade-off was chosen over another.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Clarification time | 3 to 5 minutes | guarantee scope must be set early |
| Sizing time | 5 to 7 minutes | anchors partition, retry, and retention choices |
| Initial design time | 10 to 12 minutes | enough to get buy-in before deep dive |
| Deep dive time | 8 to 10 minutes | where senior judgment becomes visible |
| Redesign time | 5 minutes | proves adaptability under new constraints |

## Architecture

Use the interview loop:

```text
clarify guarantee and scope
  -> size throughput, retention, and failure amplification
  -> propose architecture
  -> deep dive one critical subsystem
  -> close with trade-offs, observability, and redesign
```

Strong deep-dive options for this phase:

- partitioning and lag recovery in a queue
- timer durability and activity retries in a workflow engine
- catch-up policy and retry storms in a scheduler
- subscriber isolation and replay in pub/sub

## Data Model & APIs

For the drill, your answer structure is the key interface:

- guarantee statement
- capacity sheet
- component boundaries
- one deep-dive data model or API surface
- failure and observability plan
- redesign trigger and response

If you cannot explain one core record clearly, such as offsets, workflow state, or replay scope, the answer is probably still too generic.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| candidate gives vague semantics | interviewer cannot tell what is guaranteed | state the guarantee boundary in the opening |
| retries or replay are ignored in sizing | capacity claims fail under incident conditions | add amplification factors to estimates |
| deep dive is chosen too late | answer stays shallow across everything | choose one critical subsystem early |
| redesign changes nothing meaningful | architecture feels memorized | say exactly which component or policy changes |

## Observability

- metric: did the answer define delivery or execution guarantees clearly?
- metric: did the sizing include retention, retries, or replay amplification?
- metric: did each failure mode include a detection signal?
- log: note where assumptions changed during the drill
- trace: follow the answer from scope to deep dive to redesign
- SLO: complete a coherent, defensible first-pass messaging-platform design in interview time

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| one strong deep dive | proves judgment and depth | less breadth elsewhere | shallow coverage of all subsystems |
| bounded guarantees | honest and defensible | may sound less flashy | absolute marketing claims |
| explicit operator controls | operational realism | more platform complexity | black-box broker answers |

## Interview It

**Google framing:** "Design a distributed queue for internal async jobs, then explain how replay and retries behave during a downstream outage." Expect pressure on semantics, scaling, and safe recovery.

**Cloudflare framing:** "Design a multi-tenant pub/sub or job platform for edge and platform services." Expect pressure on noisy neighbors, control-plane safety, and protecting live traffic during incidents.

**Follow-ups:**
1. What if one tenant generates 50% of traffic?
2. What if the interviewer now requires replay for seven days?
3. What if exactly-once side effects are requested?
4. What if the control plane is down but data plane must keep working?
5. What if cost becomes more important than latency?

## Ship It

- `outputs/skill-messaging-platform-drill.md`

## Exercises

1. **Easy** — Give a three-minute opening for "design a distributed message queue."
2. **Medium** — Choose the right deep dive for "design a workflow engine" and justify it.
3. **Hard** — Redesign a pub/sub answer after the interviewer adds premium replay windows and strict tenant isolation.

## Further Reading

- [liquidslr/system-design-notes](https://github.com/liquidslr/system-design-notes) — useful reminder of the four-step system-design rhythm
- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) — strong source of operational follow-up pressure
