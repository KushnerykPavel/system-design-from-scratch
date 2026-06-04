# Prompt Reframing Drill

> The first senior deliverable is often not a diagram. It is a better version of the question.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Practice restating an ambiguous prompt into a crisp, designable v1 with explicit scope, assumptions, priorities, and workload shape before architecture begins.  
**Prerequisites:** `01-clarification-and-scope/07-anti-patterns`  
**Estimated time:** ~45 min  
**Primary artifact:** reframing worksheet + interview card  

## The Problem

By the end of this phase, the learner should be able to hear a vague prompt such as "Design Dropbox" or "Design a global rate limiter" and quickly transform it into:

- the core workflow
- in-scope and out-of-scope features
- dominant non-functional requirements
- key assumptions
- rough workload shape

That reframed version becomes the contract for the rest of the answer. Without it, later design discussion drifts.

This lesson is a drill, not just a concept lesson. The learner repeatedly converts ambiguous prompts into a usable opening statement.

## Clarify

- What is the smallest credible v1 that still answers the prompt?
- Which assumption or priority must be made explicit before architecture starts?
- Which workload-shape summary belongs in the reframed prompt?

## Requirements

### Functional

- Restate the prompt in two to five sentences.
- Preserve the identity of the original problem.
- Include one explicit scope boundary and one explicit priority statement.

### Non-functional

- Keep the reframing compact enough to deliver live.
- Make the reframed version specific enough to guide architecture.
- Leave room for follow-up redesign when constraints change.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Reframing length | 2 to 5 sentences | too short misses signal, too long burns time |
| Required ingredients | scope, priorities, assumptions, workload | these anchor the next design step |
| Follow-up flexibility | at least one branch left explicit | reframing should not lock the system prematurely |
| Peak factor | most important in minute 1 | strong openings have disproportionate leverage |
| Rough cost | low upfront cost, high downstream payoff | a clean reframing reduces later confusion and rework |

## Architecture

The reframing itself is pre-architecture architecture work.

A strong pattern is:

1. restate the system and primary workflow
2. name the v1 scope cut
3. rank the dominant requirement
4. state the key assumption
5. summarize the workload shape

Example:

"I’ll design a v1 photo-sharing backend focused on upload and read paths for personal accounts. I’m excluding advanced search and collaboration for now. I’ll prioritize read latency over rich write-side processing, assuming a mostly read-heavy workload with bursty celebrity hotspots. If later we require global active-active writes, the replication design changes materially."

The code artifact scores a reframed prompt for the presence of these core ingredients.

## Data Model & APIs

Represent a reframed prompt with:

- `system`
- `core_workflow`
- `scope_cut`
- `priority`
- `assumption`
- `workload_shape`

Useful review prompts:

- Is the prompt now designable?
- Did the reframing preserve the original intent?
- What future follow-up is most naturally enabled by this framing?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| reframing is just repetition | original ambiguity remains unchanged | add scope, priorities, and assumptions explicitly |
| reframing overcommits too early | no room for follow-up redesign | keep one branch explicit rather than pretending certainty |
| reframing is too long | architecture is delayed | keep to a few sentences |
| reframing distorts the prompt | interviewer would not recognize the original ask | preserve the core system identity and user journey |

## Observability

- metric: percentage of reframed prompts containing all required ingredients
- metric: average time to produce a first reframing statement
- metric: number of later architecture corrections caused by weak reframing
- SLO: each reframed prompt should create a clear path into sizing and architecture within one minute

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| compact reframing template | easy to use live | less room for nuance in the opening | free-form opening monologue |
| explicit scope and assumptions | architecture becomes grounded | small up-front time cost | implied context only |
| preserve one open branch | makes redesign follow-ups cleaner | less initial certainty | pretending the prompt is fully pinned down |

## Interview It

**Google framing:** "You get the prompt 'Design Google Calendar notifications.' Reframe it before you describe any components."

**Cloudflare framing:** "You get 'Design a customer-facing edge firewall.' Reframe the prompt into a designable v1 with clear priorities and assumptions."

**Follow-ups:**
1. Which sentence in your reframing carries the most architectural weight?
2. What did you intentionally leave flexible for later redesign?
3. How do you know your reframing did not dodge the prompt?
4. Which ingredient is easiest to forget under pressure?
5. How would the reframing change for a stricter compliance environment?

## Ship It

- `outputs/reframing-worksheet.md`
- `outputs/interview-card-prompt-reframing.md`

## Exercises

1. **Easy** — Reframe "Design a URL shortener" into a crisp v1 in under three sentences.  
2. **Medium** — Reframe "Design Slack" while preserving one explicit open question for later redesign.  
3. **Hard** — Reframe a Cloudflare-style edge product prompt so the control plane, data plane, and abuse model all remain discussable.  

## Further Reading

- [System design notes](https://github.com/liquidslr/system-design-notes) — useful baseline on structuring an opening answer
- [Google SRE workbook](https://sre.google/workbook/table-of-contents/) — good reminder that strong framing sets up strong operational reasoning
