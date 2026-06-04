# Final Capstone Review Panel

> The capstone is not one more mock; it is the moment to prove you can take feedback, compare alternatives, and defend a design under sustained review.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Run a final review-panel style capstone where multiple perspectives pressure the same design on product fit, scale, reliability, and operational realism.  
**Prerequisites:** `23-mixed-mocks-and-redesign-drills/01-consumer-backend-mock`, `23-mixed-mocks-and-redesign-drills/02-infra-platform-mock`, `23-mixed-mocks-and-redesign-drills/06-lightning-tradeoffs`  
**Estimated time:** ~90 min  
**Primary artifact:** capstone panel scorecard validator + review panel rubric  

## The Problem

Choose one substantial prompt and defend it through a simulated panel review. The panel can play three roles:

- product or user-value reviewer
- systems or scalability reviewer
- reliability and operations reviewer

Example prompts:

- design a global collaboration platform
- design a configuration and rollout service for edge workloads
- design a marketplace backend with realtime inventory and notifications
- design a multi-tenant search and analytics platform

This capstone is about more than finishing a design. It is about comparing options, absorbing criticism, and improving the answer without losing structure.

## Clarify

- Which user or platform contract will the panel judge most heavily?
- Which reviewer is most likely to disagree with your first design?
- Which trade-off are you willing to defend most strongly?
- What redesign pressure should you expect near the end?

## Requirements

### Functional

- Present a coherent full-loop answer.
- Survive multi-angle follow-up from a review panel.
- Compare alternatives instead of defending the first idea blindly.
- Close with a revised design, not just a summary.

### Non-functional

- Stay organized under cross-cutting feedback.
- Keep sizing, failure handling, and trade-offs visible throughout.
- Show humility without becoming directionless.
- Make observability and rollout part of the final design, not an afterthought.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Opening design window | 12 to 15 min | enough to establish scope and architecture |
| Panel challenge window | 15 to 20 min | tests whether the answer survives stress |
| Revised design window | 5 to 8 min | forces synthesis rather than repetition |
| Review dimensions | 8 scored categories | makes strengths and gaps concrete |
| Alternative count | 2 or 3 compared explicitly | proves judgment rather than attachment |

## Architecture

Capstone flow:

```text
clarify and size
  -> present architecture
  -> choose deep dive
  -> panel challenges assumptions
  -> compare alternatives
  -> revise the design
  -> close with final trade-offs and rollout
```

Strong capstones usually demonstrate:

1. one clear system boundary
2. one deliberate deep dive
3. one explicit alternative rejected
4. one meaningful redesign after critique

## Data Model & APIs

Panel review model:

```text
panel_review(
  clarification,
  sizing,
  architecture,
  deep_dive,
  failure_modes,
  observability,
  tradeoffs,
  communication
)
```

Each dimension should be scored from 1 to 4 with evidence.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| candidate becomes defensive | alternatives and revisions disappear | require one revised design decision |
| panel feedback breaks answer structure | no clear final architecture remains | return to priorities and restate the contract |
| operational concerns get postponed | observability and rollout are missing at the end | reserve explicit close-out time |
| comparison stays shallow | rejected alternative lacks evidence | force one requirement-based comparison |

## Observability

- metric: score by dimension across repeated capstone runs
- metric: whether critique led to a concrete design revision
- metric: whether each major risk had a detection signal
- log: original assumptions, panel objections, and final changes
- trace: opening answer -> critique -> redesign -> final defense
- SLO: complete a coherent revised design under multi-angle review pressure

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| panel-style review | surfaces weak spots quickly | feels more demanding than a normal mock | single-interviewer format only |
| scored dimensions with evidence | clearer learning loop | more evaluation overhead | vague overall feedback |
| mandatory revised design | proves adaptability | less time for polishing the original answer | defending the first draft unchanged |

## Interview It

**Google framing:** Expect a collaborative but demanding review where the interviewer wants to see whether feedback sharpens the design instead of derailing it.

**Cloudflare framing:** Expect more pressure on operational safety, traffic behavior, rollout, and how regional or tenant-specific incidents are contained.

**Suggested capstone prompts:**
1. Design a global collaboration platform.
2. Design an edge configuration and rollout service.
3. Design a marketplace backend with realtime inventory.
4. Design a multi-tenant analytics platform.

## Ship It

- `outputs/interview-card-final-capstone-review.md`
- `outputs/review-panel-rubric.md`
- `outputs/skill-final-capstone-review.md`

## Exercises

1. **Easy** — Run a panel review where only one reviewer challenges the design.
2. **Medium** — Add a second reviewer who pushes on cost or incident safety.
3. **Hard** — Re-run the capstone after the panel changes the top requirement midway through the mock.

## Further Reading

- [Google SRE Workbook](https://sre.google/workbook/table-of-contents/) — useful for redesign, rollout, and review-panel style follow-ups  
- [Designing Data-Intensive Applications](https://dataintensive.net/) — strong grounding for comparing architectural alternatives under pressure  
