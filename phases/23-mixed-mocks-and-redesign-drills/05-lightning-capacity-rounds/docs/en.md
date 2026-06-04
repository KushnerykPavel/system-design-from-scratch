# Lightning Capacity Rounds

> Senior estimation is not about perfect arithmetic; it is about finding the one or two numbers that make the design choice obvious.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Practice very short estimation rounds that force fast, high-signal capacity modeling before architecture discussion.  
**Prerequisites:** `02-estimation-and-cost/07-bottleneck-math`, `02-estimation-and-cost/08-uncertain-inputs`, `23-mixed-mocks-and-redesign-drills/03-ten-x-redesign`  
**Estimated time:** ~60 min  
**Primary artifact:** capacity round validator + worksheet  

## The Problem

You have three to five minutes per prompt. Your task is not to finish the whole design. It is to extract the smallest useful set of numbers that changes the answer:

- QPS and request mix
- peak factor
- fanout or amplification
- storage growth
- bandwidth or egress
- one likely bottleneck

The goal is to make estimation feel automatic enough that it survives stress in real mocks.

## Clarify

- Which path is the dominant hot path for this prompt?
- Is the traffic steady, bursty, or event-driven?
- Which number would most change the architecture if it were 10x wrong?
- Does the prompt care more about latency, storage growth, or network amplification?

## Requirements

### Functional

- Produce one useful estimate set in under five minutes.
- Connect each estimate to an architecture consequence.
- Identify one likely bottleneck or cost term.
- Stay concise enough to use inside a full mock.

### Non-functional

- Prefer directionally useful numbers over false precision.
- Keep units and assumptions explicit.
- Avoid long arithmetic detours.
- Make uncertainty visible instead of hiding it.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Time per round | 3 to 5 min | forces prioritization |
| Core numbers | 3 to 5 | enough to guide design without stalling |
| Error tolerance | within a rough order of magnitude | interview value comes from direction, not exactness |
| Bottleneck count | at least 1 | every round should expose a limiting resource |
| Prompt classes | serving, storage, pipeline, edge | keeps practice varied |

## Architecture

Lightning loop:

```text
clarify hot path
  -> estimate traffic and amplification
  -> estimate storage or bandwidth if relevant
  -> name one bottleneck
  -> say how the architecture changes because of it
```

Useful prompt families:

1. low-latency serving
2. high-fanout notifications
3. storage-heavy retention
4. event pipeline bursts
5. regional failover amplification

## Data Model & APIs

Round template:

```text
capacity_round(
  hot_path,
  qps,
  peak_factor,
  amplification,
  storage_or_network,
  bottleneck
)
```

Validator checks:

- did the round estimate traffic?
- did it include burst or amplification?
- did it name a bottleneck?
- did it tie numbers back to a design choice?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| too many numbers and no conclusion | round ends without a design implication | cap the round at one bottleneck and one consequence |
| no peak or skew factor | average-case math dominates | add one amplification or burst estimate |
| arithmetic is precise but irrelevant | estimates do not change architecture | ask which number would most alter the answer |
| uncertainty is hidden | numbers sound fake or brittle | state ranges and assumptions explicitly |

## Observability

- metric: round completion time
- metric: whether each round included a bottleneck and one consequence
- metric: prompt family coverage across practice sessions
- log: assumptions and confidence range per round
- trace: clarify -> estimate -> bottleneck -> consequence
- SLO: produce a useful architecture-shaping estimate set inside five minutes

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| estimate only key numbers | faster and clearer | less detailed precision | exhaustive modeling in mock time |
| range-based assumptions | more honest and adaptable | less tidy-looking output | false precision |
| one bottleneck per round | high signal | may miss secondary constraints | diffuse discussion of many limits |

## Interview It

**Google framing:** Expect the interviewer to reward fast useful math more than perfectly polished arithmetic.

**Cloudflare framing:** Expect extra attention to burstiness, failover amplification, cache hit rates, and egress consequences.

**Suggested rounds:**
1. Estimate a notification burst fanout.
2. Estimate global API failover load.
3. Estimate chat message storage growth.
4. Estimate object metadata read pressure.

## Ship It

- `outputs/capacity-worksheet-lightning-rounds.md`
- `outputs/interview-card-lightning-capacity-rounds.md`

## Exercises

1. **Easy** — Do one round and name only the bottleneck.
2. **Medium** — Do three rounds from different prompt families.
3. **Hard** — Redo one round after the peak factor changes from 3x to 20x.

## Further Reading

- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) — good mental model for capacity tied to user-visible risk  
- [System Design Primer](https://github.com/donnemartin/system-design-primer) — broad reference for common estimation categories and trade-offs  
