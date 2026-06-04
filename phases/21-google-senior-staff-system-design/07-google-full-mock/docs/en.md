# Google Full Mock Loop

> The full mock is where all the earlier lessons stop being ideas and become timing, pressure, judgment, and recovery.

**Type:** Build
**Company focus:** Google
**Learning goal:** Run a full Google-style system design mock that forces structured clarification, sizing, architecture, deep-dive choice, risk communication, and redesign under changed constraints.
**Prerequisites:** `03-design-framework-and-timing/08-full-loop-drill`, `21-google-senior-staff-system-design/01-google-rubric`, `21-google-senior-staff-system-design/06-staff-deep-dive`
**Estimated time:** ~90 min
**Primary artifact:** full mock scorecard validator + practice script

## The Problem

This lesson is a timed full-loop drill for prompts that look like a real Google senior or staff interview. The goal is not to produce the perfect architecture. The goal is to produce a coherent, scoped, workload-aware answer that reveals strong judgment in the available time.

Example prompts:

- design a large-scale file sync system
- design a metadata service for a distributed storage platform
- design a global rate-limited API serving layer
- design a search indexing pipeline with freshness constraints

## Clarify

- What is the core user journey and main workload shape?
- Which guarantee matters most if there is conflict: latency, correctness, or availability?
- Which deep dive would most change trust in the design?
- What changed constraint are you most likely to be tested on near the end?

If the prompt stays open, choose one user journey, one main bottleneck, and one system risk to organize the whole answer around.

## Requirements

### Functional

- Clarify and bound the prompt quickly.
- Produce a useful capacity model.
- Present a high-level design before going deep.
- Choose and execute one strong deep dive.
- Close with failure modes, observability, rollout, and redesign.

### Non-functional

- Stay calm and structured under pressure.
- Avoid spending the full mock on architecture without trade-offs.
- Keep the answer operationally credible.
- Adapt gracefully when the interviewer changes a key assumption.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Clarification window | 3 to 5 min | early structure earns trust |
| Estimation window | 5 to 7 min | numbers should shape topology |
| High-level design window | 10 to 12 min | enough to establish system boundaries |
| Deep-dive window | 8 to 12 min | tests judgment, not just completeness |
| Close and redesign window | 5 to 8 min | proves ownership and flexibility |

## Architecture

Use the full loop:

```text
clarify
  -> prioritize
  -> size
  -> propose high-level architecture
  -> choose deep dive
  -> explain failure and observability
  -> redesign under changed constraints
```

Good mock prompts for this phase should force at least one of:

1. serving versus storage consistency trade-offs
2. high-scale read/write asymmetry
3. migration or rollout safety
4. hotspot or partitioning pressure
5. multi-region latency versus correctness tension

## Data Model & APIs

For the mock, the main artifact is the answer plan:

```text
mock_answer(
  scope,
  top_requirements,
  sizing,
  architecture,
  deep_dive,
  risks,
  observability,
  redesign
)
```

Scorecard fields:

- scope clarity
- requirement prioritization
- quantitative sizing
- high-level architecture quality
- deep-dive quality
- failure and observability coverage
- redesign quality

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| answer spends too long clarifying | little architecture progress by minute 10 | cap questions and state assumptions |
| answer jumps into detail too early | no shared high-level design exists | force a top-level architecture before depth |
| answer lacks wrap-up | no failure or rollout narrative near the end | reserve close-out time explicitly |
| changed constraint causes collapse | redesign discards earlier reasoning | anchor the answer in assumptions and priorities |

## Observability

- metric: time spent in each interview phase
- metric: whether the answer contained at least one useful estimate
- metric: whether the deep dive targeted the main system risk
- metric: whether one detection signal was named for each major failure mode
- log: assumptions, corrections, and redesign triggers
- trace: clarify -> size -> design -> deep dive -> wrap-up -> redesign
- SLO: complete a coherent full-loop system design answer inside mock time

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| one disciplined loop | steady progress and better trust | less room for improvisational wandering | unstructured answer |
| one strongest deep dive | higher signal | reduced breadth elsewhere | many shallow dives |
| explicit redesign step | shows adaptability | forces trade-offs to stay visible | pretending the first design solved everything |

## Interview It

**Google framing:** Run this as if the interviewer will reward clarity, modeling, and judgment more than ornamental detail. A good performance sounds collaborative, concrete, and resilient under follow-up pressure.

**Suggested prompts:**
1. Design a globally distributed metadata service.
2. Design a large-scale document collaboration backend.
3. Design a storage-backed search indexing pipeline.
4. Design a low-latency abuse-detection serving system.

## Ship It

- `outputs/skill-google-full-mock.md`
- `outputs/interview-card-google-full-mock.md`

## Exercises

1. **Easy** - Time-box a 45-minute answer into six stages.
2. **Medium** - Pick the best deep dive for "design a globally distributed metadata service."
3. **Hard** - Run a full mock where the interviewer changes the top priority from latency to correctness at minute 30.

## Further Reading

- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) - good operational pressure source for follow-up questions
- [liquidslr/system-design-notes](https://github.com/liquidslr/system-design-notes) - useful reminder of the four-step interview rhythm this mock extends
