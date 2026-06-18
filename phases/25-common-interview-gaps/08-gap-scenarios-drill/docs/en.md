# Gap Scenarios Drill

> Speed under pressure. This lesson is not about learning new patterns — it is about surfacing the non-obvious constraint fast enough to make the interviewer confident you have done this before.

**Type:** Learn  
**Company focus:** Balanced  
**Learning goal:** Practice the senior-level design instinct across all seven gap scenarios in this phase: identify the non-obvious constraint within 2 minutes, sketch the critical subsystem, state the key trade-off, and articulate the failure mode — all within a 45-minute interview time budget.  
**Prerequisites:** `25-common-interview-gaps/01-file-sync-service`, `25-common-interview-gaps/02-ad-click-aggregation`, `25-common-interview-gaps/03-ride-sharing-dispatch`, `25-common-interview-gaps/04-live-streaming`, `25-common-interview-gaps/05-booking-and-reservation`, `25-common-interview-gaps/06-code-deployment-pipeline`, `25-common-interview-gaps/07-feature-flag-service`  
**Estimated time:** ~90 min  
**Primary artifact:** interview card

## The Problem

This drill lesson revisits all seven scenarios from Phase 25 under interview time pressure. The goal is not to produce a perfect design — it is to demonstrate that you can rapidly identify the constraint that separates a senior answer from a mid-level answer, communicate a coherent architecture under time pressure, and defend trade-offs without hedging.

The scenarios in this phase share a pattern: there is always one constraint that mid-level candidates miss, and it is always the constraint that makes the design non-trivial. File sync → conflict resolution and block deduplication. Ad clicks → deduplication and late events. Ride sharing → atomic driver assignment and ETA accuracy. Live streaming → ingest-to-edge pipeline and chat fan-out at 1M viewers. Booking → seat hold atomicity and payment-reservation decoupling. Deployment → progressive rollout with automatic rollback triggers. Feature flags → in-process evaluation with push propagation.

This lesson trains you to surface that constraint in the first 2 minutes, so the rest of the interview is a productive deep-dive rather than a grope toward the real problem.

## Clarify

- For each scenario, ask only the one clarifying question that most changes the architecture: (e.g., "Is overbooking permitted?" for booking, or "What is the acceptable latency for flag propagation?" for feature flags.)
- State your assumption out loud if the interviewer does not answer, then proceed — do not wait for every detail.
- If the interviewer pivots to a failure mode or trade-off, follow them — do not finish your original thought at the cost of ignoring their cue.

## Requirements

### Functional

- Complete a recognizable high-level design for any of the seven scenarios within 15 minutes.
- Identify and verbalize the non-obvious constraint within the first 2 minutes of each scenario.
- Sketch the critical subsystem in ASCII or boxes-and-arrows within 5 minutes.
- State at least one key trade-off with two sides compared, not just "it depends."
- Identify at least one failure mode and its mitigation before being asked.

### Non-functional

- No more than 5 minutes on clarification and requirements gathering.
- No more than 10 minutes on capacity estimation (keep it rough; the goal is orders of magnitude).
- No less than 10 minutes on deep-dive into the hardest subsystem.
- Wrap up with trade-offs and failure modes before the 40-minute mark.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Scenarios to drill | 7 scenarios | covers full Phase 25 scope; one drill session per scenario |
| Time per full scenario | 45 min target | matches real interview duration |
| Time to surface non-obvious constraint | 2 min target | most differentiating signal for senior candidates |
| Trade-offs to state per session | 2–3 | quality signal: can you compare two sides without prompting? |
| Failure modes to identify | 1–2 unprompted | senior candidates mention failure modes before the interviewer asks |

## Architecture

The drill format for each scenario:

```text
[0:00 - 2:00] Clarify and identify non-obvious constraint
  -> Ask the one pivotal clarifying question
  -> State the key non-obvious constraint aloud:
     "The hard part here is _____, which most candidates miss."

[2:00 - 7:00] Rough capacity estimate
  -> Users / QPS / data volume — one number per dimension
  -> Flag which number drives the critical design decision

[7:00 - 17:00] High-level architecture
  -> Sketch the full flow end-to-end (ingest → process → store → serve)
  -> Label every component with its role
  -> Mark the one component that is the hardest to scale or get right

[17:00 - 30:00] Deep-dive: the hardest subsystem
  -> Go deep on the non-obvious constraint
  -> Show the mechanism, not just the name
     ("I use a CAS on the driver status field" not "atomic assignment")
  -> State the correctness guarantee the mechanism provides

[30:00 - 38:00] Trade-offs and alternatives rejected
  -> Pick two decisions and compare both sides
  -> Name the alternative you rejected and why

[38:00 - 42:00] Failure modes
  -> One failure mode you are most worried about
  -> Detection metric and mitigation

[42:00 - 45:00] Wrap up
  -> Summarize: "If I had 5 more minutes I would dig into ___."
```

## Data Model & APIs

Use the following as a quick-reference cheat sheet of the pivotal data model element for each scenario:

```text
File sync:       FileVersion { version_id, file_id, block_list[], parent_version_id }
Ad clicks:       AggregateRecord { campaign_id, window_start, click_count, correction_version }
Ride sharing:    Driver { driver_id, status: available|pending:{ride_id}|in_trip }
Live streaming:  Segment { stream_id, rendition, sequence_num, storage_key }
Booking:         Seat { seat_id, status: available|held|confirmed, held_for, hold_expires_at }
Deployment:      Stage { stage_id, deployment_id, traffic_percent, health_snapshot, status }
Feature flags:   Flag { flag_key, is_enabled, rollout_percent, targeting_rules[] }
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Spending too much time on requirements, not enough on design | interviewer becomes quiet; no design on the whiteboard after 10 min | set a personal timer; state "I'll make this assumption and move to design" |
| Describing a component by name without explaining the mechanism | interviewer asks "how exactly does that work?" | always follow a component name with "which works by ___" |
| Missing the non-obvious constraint and designing for the obvious one | interviewer redirects to a harder problem | practice recognizing the pattern: "what would break at 100x scale or on the third retry?" |
| Hedging on trade-offs without committing to a position | interviewer pushes back with "but which would you choose?" | state a position with one reason; you can qualify it but never leave it as "it depends" alone |

## Observability

- metric: time-to-constraint (how quickly you identify the non-obvious constraint) — self-measured during practice sessions
- metric: follow-up question rate (how often the interviewer must prompt you vs you volunteer information) — lower is better
- metric: trade-off commitment rate (fraction of trade-off discussions where you state a position) — target 100%
- log: written post-drill notes on which constraints you missed and why — the most valuable output of each drill session
- SLO: achieve a complete high-level design sketch within 15 minutes for any of the seven scenarios after three full drill sessions

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Drill all seven scenarios in one session | exposes cross-scenario patterns; builds stamina | 7 × 45 min = 5+ hours; better as seven separate sessions | drill each scenario once and never revisit — misses the pattern recognition that comes from repeated exposure |
| Focus on the non-obvious constraint first | earns senior signal early; redirects the interview | risks superficial coverage of basic requirements | cover all requirements before going deep — safe but does not differentiate from mid-level |
| State failure modes before being asked | demonstrates operational maturity; impresses interviewers | takes time away from architecture depth | wait for interviewer to ask about failures — reactive posture signals less experience |

## Interview It

**Google framing:** "We'll do a rapid-fire session: I'll give you a system design prompt and you have 45 minutes per scenario. I'll interrupt if I want to go deeper somewhere." Practice treating each interruption as the interviewer revealing what they care about, not as a derailment.

**Cloudflare framing:** "For each system, tell me: what breaks at 10x scale, and what would you change?" Practice leading with the failure mode rather than the architecture.

**Follow-ups across all scenarios:**
1. Which of the seven scenarios has the hardest consistency problem? Defend your answer.
2. Which scenario would you redesign if the latency requirement tightened by 10x?
3. Where across these seven systems does eventual consistency cause the most user-visible harm?
4. Which system is most sensitive to a third-party dependency failure (payment processor, maps API, fraud service)?
5. If you could only add one metric to each system, what would it be and why?

## Ship It

- `outputs/interview-card-gap-scenarios-drill.md`

## Exercises

1. **Easy** — Pick any two scenarios from this phase and write down the non-obvious constraint for each in one sentence without looking at the lesson. Compare with the lesson text.
2. **Medium** — Set a 45-minute timer and complete a full design session for the ride-sharing dispatch scenario. Record which minute you identified the atomic driver assignment constraint. Target: under 2 minutes.
3. **Hard** — Design a system that incorporates three scenarios from this phase simultaneously: a booking platform where drivers are also available for ride-sharing, and every interaction is feature-flagged. Identify where the data models conflict and how you resolve the conflicts.

## Further Reading

- https://www.hellointerview.com/learn/system-design/in-a-hurry/introduction — Hello Interview's speed-focused system design guide; complements this drill with scored practice sessions
- https://github.com/donnemartin/system-design-primer — System Design Primer; useful for cross-referencing patterns across the scenarios in this phase
- https://sre.google/sre-book/table-of-contents/ — Google SRE Book; the operational mindset behind the failure modes and observability patterns practiced in this drill
