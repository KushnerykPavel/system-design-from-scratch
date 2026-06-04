# Job Scheduler and Retry Platform

> Schedulers fail less from cron syntax than from thundering herds, stale jobs, and unclear retry ownership.

**Type:** Build
**Company focus:** Balanced
**Learning goal:** Design a scheduler platform that handles periodic and one-shot jobs with shard ownership, retry policy, catch-up behavior, and tenant isolation as first-class concerns.
**Prerequisites:** `02-estimation-and-cost/05-burstiness`, `10-reliability-retries-and-backpressure/06-retry-budgets`, `16-application-backends/04-notification-system`
**Estimated time:** ~75 min
**Primary artifact:** scheduler-plan validator + observability checklist

## The Problem

Design a job scheduler used for cron-like tasks, delayed tasks, and retry orchestration across many internal teams. Jobs may fire at fixed intervals, after specific delays, or as deferred retries. The platform must avoid synchronized load spikes, handle missed runs after outages, and provide enough control for operators to inspect failures without turning every retry into custom application code.

This lesson matters because scheduler answers often sound simple until the interviewer asks what happens at midnight UTC, after a control-plane outage, or when millions of retries align on the same minute.

## Clarify

- Are jobs latency-sensitive or is minute-level skew acceptable?
- Should missed runs be replayed, skipped, or coalesced after downtime?
- Are jobs idempotent, or must the scheduler help prevent duplicates?
- Do tenants have different priority classes or quotas?

If the prompt is broad, assume mixed periodic and delayed jobs, at-least-once triggering, configurable catch-up policy, and strong pressure to avoid synchronized spikes.

## Requirements

### Functional

- Register periodic and one-time jobs with retry policy.
- Trigger jobs close to schedule while tolerating node failure.
- Support catch-up, skip, or coalesce policies for missed runs.
- Track attempt history and expose operator controls.
- Isolate failed or toxic jobs after bounded retries.

### Non-functional

- Smooth load to avoid burst amplification from aligned schedules.
- Recover shard ownership quickly after scheduler-node failure.
- Prevent one tenant or job family from consuming all dispatch capacity.
- Make stale jobs and expired deadlines visible.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Active scheduled jobs | 500M | storage and shard scans must scale |
| Dispatches per minute | 40M peak | shapes lease ownership and batching |
| Retry amplification | up to 3x during incidents | retry design can dominate steady-state load |
| Deadline-sensitive jobs | 8% of workload | prioritization matters more than average dispatch rate |
| Peak factor | 10x at common wall-clock boundaries | jitter and spread are core architecture choices |

## Architecture

```text
job API
  -> scheduler control plane
  -> sharded schedule store
  -> shard lease manager
  -> dispatch queue
  -> execution workers
  -> retry / dlq path
```

Design notes:

1. Shard by next-run bucket plus tenant or job family so scans remain bounded.
2. Use shard leases to avoid duplicate dispatch after node loss.
3. Add jitter and schedule spreading for large aligned cohorts.
4. Make deadline-aware dropping or deprioritization explicit for jobs that are worthless after a delay.

## Data Model & APIs

Core records:

```text
job_id
tenant_id
schedule_spec
next_run_time
deadline_time
retry_policy
catch_up_policy
attempt_count
last_result
```

Useful interfaces:

- `CreateJob(job_spec, schedule, retry_policy)`
- `PauseJob(job_id)`
- `ResumeJob(job_id)`
- `DispatchBatch(shard_id, lease_token, now)`
- `ReportAttempt(job_id, result, next_retry_time)`

Strong answers distinguish schedule ownership, dispatch, and actual job execution.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| midnight cron storm overloads workers | dispatch burst metrics and worker-queue lag | jitter, spread windows, and priority lanes |
| scheduler node dies holding shard ownership | stale lease age and overdue dispatch count | short leases with safe steal and replay checks |
| outage creates huge missed-run backlog | missed-run queue growth and stale-job count | explicit catch-up policy, coalescing, or skip semantics |
| retries flood the platform during dependency incident | retry dispatch ratio and dependency-error correlation | bounded retry budgets, exponential backoff, and tenant caps |

## Observability

- metric: dispatch lag and overdue-job count by priority class
- metric: retry rate, dead-letter rate, and missed-run volume
- metric: shard-scan duration and lease-steal frequency
- metric: stale-job drops versus successful catch-up runs
- log: operator schedule changes and manual replays with actor identity
- trace: schedule evaluation to dispatch to execution result for sampled jobs
- SLO: 99% of standard jobs are dispatched within the scheduler target of their due time under normal load

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| jittered dispatch | avoids herd effects | less exact wall-clock fire time | perfectly aligned firing |
| explicit catch-up policies | predictable outage recovery | more product complexity | implicit replay of everything |
| scheduler-owned retry policy | centralized safety and visibility | less per-team freedom | every consumer invents retries alone |

## Interview It

**Google framing:** "Design a large internal scheduler for cron jobs and delayed retries." Expect follow-ups on missed runs, worker overload, and retry amplification during dependency failure.

**Cloudflare framing:** "Design a platform scheduler for many tenants and operational jobs." Expect questions on isolation, control-plane failure, and protecting the fleet from synchronized work.

**Follow-ups:**
1. What should happen to a job whose deadline has already passed?
2. How would you schedule millions of hourly jobs without an hourly spike?
3. When should retries stay in the scheduler versus in the application?
4. How would you let operators replay only one tenant's failed jobs?
5. What changes for multi-region active-active scheduling?

## Ship It

- `outputs/observability-checklist-job-scheduler.md`

## Exercises

1. **Easy** — Choose a catch-up policy for nightly reports after a two-hour outage.
2. **Medium** — Redesign the scheduler for one tenant that owns 30% of all jobs.
3. **Hard** — Add regional failover while avoiding duplicate dispatch of the same due job.

## Further Reading

- [Google SRE: Addressing Cascading Failures](https://sre.google/sre-book/addressing-cascading-failures/) — useful for retry and overload thinking
- [Quartz scheduler misfire instructions](https://www.quartz-scheduler.org/documentation/) — practical context for missed-run semantics
