# Backend Product Drill

> Senior product-backend answers win by structuring ambiguity, not by naming the most components.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Synthesize the phase by practicing how to clarify, size, choose deep dives, and defend trade-offs across common application-backend prompts.  
**Prerequisites:** `16-application-backends/01-url-shortener`, `16-application-backends/02-news-feed`, `16-application-backends/03-chat-system`, `16-application-backends/04-notification-system`, `16-application-backends/05-collaboration-backend`, `16-application-backends/06-fanout-patterns`  
**Estimated time:** ~60 min  
**Primary artifact:** backend answer scorecard  

## The Problem

This drill is the capstone for the phase. Instead of introducing one more product, it teaches a reusable answer pattern for application-backend interviews where the interviewer wants both product intuition and systems depth.

The learner should be able to handle prompts like feed, chat, notifications, collaboration, or short links without giving the same generic architecture every time.

## Clarify

- What single product outcome matters most: freshness, delivery confidence, ranking quality, or cost?
- Which part of the system is read-heavy versus write-heavy?
- Where is the likely skew: hot objects, celebrity producers, reconnect storms, or urgent priority traffic?

If an interviewer is intentionally vague, choose one crisp assumption set, state it, and keep the rest of the answer consistent with it.

## Requirements

### Functional

- Extract the product-specific contract before drawing the system.
- Produce rough sizing before selecting architecture.
- Choose one or two deep dives that match the real risk of the prompt.
- Explain degraded behavior, observability, and rollout.

### Non-functional

- Keep the answer organized enough to survive follow-up pressure.
- Avoid copy-pasting one design into every product domain.
- Tie trade-offs to metrics and failure modes, not just preferences.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Clarification budget | 2 to 5 questions | too few causes shallow assumptions; too many wastes time |
| Sizing budget | 2 to 4 key numbers | enough to guide architecture without turning into bookkeeping |
| Deep dives | 1 or 2 | breadth without losing depth |
| Failure areas | at least 3 | shows operational maturity |
| Time budget | 35 to 45 minutes typical | forces prioritization |

## Architecture

Recommended answer flow:

```text
clarify product contract
  -> define requirements
  -> size dominant read/write paths
  -> choose high-level architecture
  -> deep dive on the true bottleneck
  -> describe failure modes and observability
  -> close with trade-offs and redesign options
```

This is not just interview theater. It is how you keep application backends from collapsing into the same generic queue-plus-cache diagram every time.

## Data Model & APIs

For this drill, the "API" is your answer structure:

- `ClarifyPrompt()`
- `EstimateLoad()`
- `ChooseWorkPlacement()`
- `NameFailureModes()`
- `CloseWithTradeoffs()`

If the answer cannot explain the core API or state boundary of the product, the design is probably still too hand-wavy.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| candidate jumps to architecture without contract | vague requirements and weak trade-offs | force explicit clarification and stated assumptions |
| sizing is skipped entirely | component choices feel arbitrary | estimate at least traffic, skew, and retention |
| deep dive focuses on the wrong subsystem | follow-ups expose shallow reasoning | choose the area with real amplification or correctness risk |
| observability is omitted | incident handling remains magical | attach metrics and SLOs to the main user-visible outcomes |

## Observability

- metric: does the answer define one clear latency or freshness SLO for the product?
- metric: does each deep dive name at least one detection signal?
- metric: are skew and amplification measured explicitly?
- log: record which assumption set the learner chose for a drill
- trace: map user-visible action to backend pipeline stages during practice
- SLO: every strong answer should connect its main product promise to one measurable outcome

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| one crisp assumption set | keeps the answer coherent | may miss alternative interpretations | trying to solve every possible product variant at once |
| two intentional deep dives | shows real depth | leaves less time for minor components | broad but shallow component inventory |
| degraded-mode explanation | demonstrates operational maturity | requires more discipline and time | ideal-path-only answer |

## Interview It

**Google framing:** "I will give you a backend product prompt; structure the answer like a senior engineer." Expect pressure on sizing first, then one deliberate deep dive.

**Cloudflare framing:** "Design the backend, then tell me how it behaves under real operational stress." Expect stronger focus on traffic shape, routing, and incident behavior.

**Follow-ups:**
1. Which clarification question would most change your design?
2. Which metric would tell you you picked the wrong fanout strategy?
3. What if the same product must now work in multiple regions?
4. How would you roll out a redesign with minimal user impact?
5. What is the most important trade-off you would say out loud in the interview?

## Ship It

- `outputs/skill-backend-product-drill.md`

## Exercises

1. **Easy** — Pick one of the phase prompts and write only the clarification and sizing sections.
2. **Medium** — For the same prompt, choose one deep dive and write the failure-mode table.
3. **Hard** — Practice answering two different prompts back to back without reusing the same architecture story.

## Further Reading

- [System Design Interview - An Insider's Guide](https://github.com/liquidslr/system-design-notes) — helpful for interview flow and pacing, even though this course goes deeper operationally  
- [Google SRE](https://sre.google/books/) — strong reference for turning abstract architectures into measurable operational systems  
