# Common Interviewer Moves and How to Respond

> Interviewer pressure is not a derailment. It is part of the prompt.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Learn to recognize common interviewer interventions and respond in a way that preserves structure, collaboration, and senior signal.  
**Prerequisites:** `03-design-framework-and-timing/06-redesign-prompts`  
**Estimated time:** ~60 min  
**Primary artifact:** interviewer-move response card + recovery checklist  

## The Problem

Senior interviewers rarely sit silently while you present. They interrupt, tighten constraints, challenge assumptions, ask for trade-offs, and redirect you toward the area they care about most.

Candidates who interpret these moves as failure tend to get flustered. Candidates who treat them as useful signals usually recover and gain control.

## Clarify

- Is the interviewer changing the prompt or probing the existing design?
- Are they testing trade-off reasoning, detail depth, or composure under ambiguity?
- Does the interruption require a redesign or just a brief clarification?
- What part of my answer should remain stable after I respond?

If unsure, say what you think they are asking and answer that directly.

## Requirements

### Functional

- Recognize common interviewer moves quickly.
- Answer the move without abandoning the answer structure.
- Use interruptions to improve the relevance of the design.

### Non-functional

- Stay calm and collaborative.
- Avoid defensive or overly verbose responses.
- Preserve forward momentum after interruptions.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Major interruptions | 2-5 per interview | normal at senior levels |
| Recovery time | under 60 seconds | keeps pacing intact |
| Clarification overhead | 1 short sentence | avoids drift |
| Trade-off branches | 1-2 per probe | keeps response tight |
| Redesign frequency | occasional | some probes change the architecture materially |

## Architecture

Common interviewer moves:

1. Scale changes
2. Trade-off probes
3. Assumption challenges
4. Constraint tightening
5. Detail pulls into implementation

Strong response pattern:

- acknowledge the move
- classify it
- explain what changes first
- answer tightly
- reconnect to the broader design

## Data Model & APIs

The code artifact maps interviewer moves to default response strategies. It is deliberately simple, because the main goal is building a reusable conversational reflex.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Defensive response | candidate argues with the prompt change | acknowledge and adapt instead |
| Full derailment | answer never returns to the main structure | explicitly reconnect after the probe |
| Over-answering | one follow-up consumes many minutes | cap the response and summarize |
| Misclassified move | detail probe treated as full redesign | restate the interpreted question before answering |

## Observability

- metric: time from interruption to structured response
- metric: whether the answer returned to the main flow afterward
- metric: number of trade-off probes answered with explicit options
- log: interviewer moves that caused pacing collapse

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| concise interruption handling | preserves pace | may omit some nuance | long defensive explanations |
| explicit move classification | increases clarity | costs a sentence upfront | reacting without framing |
| reconnecting to main flow | keeps the answer coherent | requires discipline | letting each probe create a new branch |

## Interview It

**Google framing:** "Design a metrics platform." The interviewer may interrupt to ask about cost, cardinality, or failure isolation. Each move is a hint about where the real signal is.

**Cloudflare framing:** "Design global edge bot mitigation." The interviewer may abruptly shift to propagation speed, false positives, or failure blast radius. A strong response adapts without losing the control-plane/data-plane structure.

**Follow-ups:**
1. Which interviewer moves usually require true redesign?
2. How do you answer trade-off probes without monologuing?
3. What is the best way to handle challenged assumptions?
4. When should you ask the interviewer a clarification question back?
5. How do you preserve a calm tone under rapid interruptions?

## Ship It

- `outputs/interviewer-move-checklist.md`
- `outputs/interview-card-interviewer-moves.md`

## Exercises

1. **Easy** — Practice answering a scale-change interruption in under 30 seconds.  
2. **Medium** — Respond to a trade-off challenge for a distributed cache design and reconnect to the main answer.  
3. **Hard** — Handle three consecutive interviewer moves during a mock design for a global API gateway.  

## Further Reading

- [Google SRE book](https://sre.google/sre-book/table-of-contents/) — good mental model for responding to operational probes  
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline interview structure that becomes more dynamic under interruptions  
