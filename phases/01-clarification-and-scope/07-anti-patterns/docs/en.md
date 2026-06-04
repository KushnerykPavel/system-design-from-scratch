# Senior-Level Clarification Anti-Patterns

> Most weak designs are not killed by a missing database. They are killed by a bad first five minutes.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Learn to recognize and avoid the clarification mistakes that make senior answers feel shallow, evasive, or operationally naive.  
**Prerequisites:** `01-clarification-and-scope/06-workload-shape`  
**Estimated time:** ~45 min  
**Primary artifact:** anti-pattern checklist + recovery prompts  

## The Problem

By this point in the phase, the learner knows what good clarification looks like. This lesson turns that inside out by naming the failure modes that repeatedly sink otherwise strong candidates:

- architecture-first without scoping
- too many low-signal questions
- silent assumptions
- vague "for simplicity" cuts
- no requirement ranking
- no workload-shape translation
- no explicit failure or observability angle

Senior interview performance improves faster when the learner can spot these mistakes in real time and recover before the answer drifts too far.

## Clarify

- Which anti-pattern is most likely for this learner under time pressure: rushing, over-questioning, or hand-wavy assumptions?
- What recovery move is smallest but highest leverage once the anti-pattern appears?
- Which anti-pattern destroys the most downstream reasoning if left uncorrected?

## Requirements

### Functional

- Identify anti-patterns quickly during live practice.
- Link each anti-pattern to a concrete recovery action.
- Build a mental checklist that fits within the interview pace.

### Non-functional

- Keep the checklist short enough for real use.
- Focus on behaviors that change the quality of the design, not presentation nitpicks.
- Preserve learner confidence by making recovery practical, not punitive.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Critical anti-patterns | 5 to 7 is enough | too many warnings become unusable live |
| Recovery time | 10 to 30 seconds each | good recovery must be small and immediate |
| Downstream damage | high for early mistakes | anti-patterns compound as the interview continues |
| Peak factor | highest in the first 10 minutes | early discipline has outsized leverage |
| Rough cost | low checklist cost, high saved interview signal | small corrections prevent large architecture drift |

## Architecture

Treat clarification anti-patterns like system failure modes:

- detect them early
- stop the spread
- apply the least disruptive correction

Example recoveries:

- **rushed into architecture** -> pause and restate scope plus top priorities
- **too many trivia questions** -> ask one pivot question and move on
- **silent assumption** -> say the assumption aloud and tie it to the current branch
- **no workload shape** -> summarize reads, writes, and bursts before adding more components

The lesson’s code artifact scans a practice session summary for missing behaviors and suggests recovery prompts.

## Data Model & APIs

Represent a practice summary with boolean signals:

- `clarified_scope`
- `ranked_requirements`
- `logged_assumptions`
- `named_workload_shape`
- `stated_scope_cut`

Useful review prompts:

- Which missing behavior would hurt the rest of the answer most?
- Which anti-pattern appeared first?
- What is the smallest recovery that restores coherence?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| anti-pattern recognized too late | half the interview already built on weak framing | insert a compact recovery summary immediately |
| over-correction | recovery takes longer than the original mistake | use one sentence, not a restart |
| checklist too long to use | learner cannot remember it under pressure | keep only the highest-signal anti-patterns |
| learner treats checklist mechanically | behaviors are named but not tied to design quality | connect each item to architecture consequences |

## Observability

- metric: number of anti-patterns detected during practice debrief
- metric: average time between detecting a mistake and applying a recovery
- metric: number of answers reaching clear scope and workload framing before architecture
- SLO: each common anti-pattern should have a one-sentence recovery move

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| short anti-pattern checklist | usable in live interviews | less exhaustive coverage | long review rubric |
| recovery-focused framing | preserves confidence and momentum | may feel repetitive in practice | pure criticism with no correction path |
| behavior-level signals | easy to observe in mocks | abstracts away some nuance | fully bespoke debrief every time |

## Interview It

**Google framing:** "You start designing before clarifying scale and consistency. How do you recover in under 20 seconds without losing credibility?"

**Cloudflare framing:** "In an edge design round, what clarification anti-pattern is most likely to produce a weak architecture, and why?"

**Follow-ups:**
1. Which anti-pattern causes the most expensive downstream redesign?
2. Which one is easiest to miss while sounding confident?
3. What recovery would you use if the interviewer interrupts your clarification early?
4. Which anti-pattern is most common when the prompt is globally distributed?
5. How do you keep the recovery brief enough to preserve momentum?

## Ship It

- `outputs/anti-pattern-checklist.md`
- `outputs/recovery-prompts.md`

## Exercises

1. **Easy** — Read a short design opening and identify two clarification anti-patterns.  
2. **Medium** — Write one-sentence recovery moves for five common anti-patterns.  
3. **Hard** — Diagnose a mock answer that sounds polished but skips scope, assumptions, and workload shape.  

## Further Reading

- [Google SRE Book](https://sre.google/sre-book/table-of-contents/) — useful for seeing how omitted assumptions and priorities create operational blind spots
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline framework that these anti-patterns often violate
