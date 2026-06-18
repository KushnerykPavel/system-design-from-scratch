# Feature Flag Service

> A feature flag looks like an if-statement. At scale it is a distributed configuration system with sub-millisecond evaluation latency, targeting rules, kill switches, and A/B experiment integration.

**Type:** Learn  
**Company focus:** Cloudflare  
**Learning goal:** Design a feature flag service that evaluates targeting rules in under 1 ms, propagates flag changes to all service instances within seconds, supports kill switches and gradual rollouts, and integrates with an A/B experimentation platform without coupling flag evaluation to experiment assignment.  
**Prerequisites:** `06-caching-and-invalidation/02-freshness-models`, `14-rate-limiters-ids-and-hashing/04-consistent-hashing`, `13-multi-region-cdn-and-edge-traffic/05-edge-compute`  
**Estimated time:** ~60 min  
**Primary artifact:** capacity sheet

## The Problem

Design a feature flag service used by 500 microservices to gate feature rollouts, run A/B experiments, and implement kill switches. Each service evaluates flags on every request, so evaluation must be under 1 millisecond — slower than a local cache lookup is not acceptable. Flags must propagate to all service instances globally within 30 seconds of a change. A kill switch must be able to disable a feature for all traffic within 10 seconds.

Mid-level candidates describe a database-backed flag store with an API. Senior candidates recognize three distinct operational problems. First, evaluation latency: no service can afford a network round-trip to evaluate a flag on every request. Flags must be evaluated against a local copy of the flag state. Second, propagation speed: when an engineer changes a flag, they expect it to take effect quickly. The local copy must be kept fresh via a push notification, not polling. Third, the kill switch guarantee: if a feature is causing an outage, the engineer needs to disable it within 10 seconds, not wait for the next polling interval.

The experimentation integration is also a source of design confusion. A/B experiments and feature flags are related but not the same. Flag evaluation determines whether a feature is on or off for a user. Experiment assignment determines which variant a user sees. These two concerns should be separated in the data model even though they share the same evaluation infrastructure.

## Clarify

- What is the expected number of flags in the system, and what is the typical complexity of targeting rules (simple user percentage vs complex attribute matching)?
- Is the flag service responsible for experiment assignment and statistical analysis, or only for feature gating?
- What is the acceptable flag propagation latency — how quickly must a flag change be visible to all clients?

If the interviewer does not specify, assume up to 10K flags, targeting rules based on user attributes and percentage rollout, experiment assignment is separate but uses the same SDK, and propagation target of under 30 seconds (kill switch under 10 seconds).

## Requirements

### Functional

- Create, update, and delete feature flags with targeting rules (percentage rollout, user segment, attribute matching).
- Evaluate flags for a given user context (user_id, attributes, environment) in the calling service's process.
- Propagate flag state changes to all service instances within 30 seconds (10 seconds for kill switches).
- Provide an audit log of all flag changes with actor, timestamp, and before/after state.
- Support flag environments: production, staging, development — each with independent flag state.
- Integrate with an A/B experiment framework to associate flag variants with experiment assignments.

### Non-functional

- Flag evaluation latency: under 1 ms in-process (requires local flag state copy, not network call).
- Flag propagation latency: under 30 seconds for normal changes; under 10 seconds for kill switches.
- Flag store availability: 99.99% (flag evaluation must work even when the central flag store is unavailable).
- SDK footprint: small enough to embed in every service without significant memory or CPU overhead.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Active flags | 10K flags, each with targeting rules | total state size: ~50 MB if rules are compact; fits in each service's memory |
| Services | 500 services × 100 instances = 50K SDK instances | propagation fan-out target; push notification must reach 50K receivers within 30s |
| Evaluations per second | 50K instances × 5K req/s each = 250M evaluations/s total | must be in-process; any network call at this rate would be infeasible |
| Flag changes per day | 1K changes/day (engineers modifying flags) | low write rate; the flag store is not write-heavy |
| Kill switch events | 10 per day on average; latency-critical | sized as a priority path in the propagation system |

## Architecture

```text
flag management UI / API
  -> engineer creates or updates a flag
  -> writes to flag store (primary DB: PostgreSQL or similar)
  -> writes audit log record
  -> publishes flag_changed event to propagation bus

propagation bus (pub/sub)
  <- flag_changed event
  -> flag streaming service (WebSocket or Server-Sent Events endpoint)
  -> or: push to Redis pub/sub; SDK instances subscribe via long-lived connection

SDK (embedded in each service instance)
  -> on startup: fetch full flag snapshot from flag streaming service
  -> maintains in-memory flag state (sorted by flag_id)
  -> receives incremental updates via WebSocket / SSE / Redis pub/sub
  -> on flag_changed: update in-memory state atomically
  -> evaluation: in-process rule evaluation against local state

flag evaluation (in SDK, no network call)
  -> input: {flag_key, user_context: {user_id, attributes, environment}}
  -> lookup flag in local state by flag_key + environment
  -> evaluate targeting rules in order:
     1. kill switch (is_enabled = false) -> return OFF
     2. user_id in allowlist -> return ON
     3. user_id in denylist -> return OFF
     4. percentage rollout: hash(flag_key + user_id) mod 100 < rollout_percent
  -> return variant assignment

kill switch path (priority propagation)
  -> engineer marks flag as kill_switch
  -> propagation bus delivers to dedicated kill-switch topic with lower timeout
  -> SDK processes kill_switch updates before other flag updates in its queue
```

## Data Model & APIs

Core entities:

```text
Flag        { flag_key, description, environment, is_enabled, rollout_percent,
              targeting_rules[], default_variant, created_by, created_at, updated_at }
TargetingRule { rule_id, flag_key, priority, condition_type, condition_value, variant }
Variant     { variant_key, flag_key, payload }
AuditEvent  { event_id, flag_key, environment, actor, change_type, before_state,
              after_state, timestamp }
ExperimentAssignment { user_id, flag_key, variant, assigned_at, experiment_id }
```

Key APIs (control plane — not on hot path):

- `POST /v1/flags` — create a flag with initial state
- `PUT /v1/flags/{flag_key}` — update flag state, rules, or rollout percentage
- `DELETE /v1/flags/{flag_key}` — archive flag (not hard delete; audit retention)
- `GET /v1/flags` — list all flags for an environment (used by SDK on startup)
- `GET /v1/flags/{flag_key}/audit` — returns full change history for the flag

SDK interface (in-process, no network):

```text
flagsdk.Evaluate(flagKey string, userCtx UserContext) (Variant, error)
flagsdk.IsEnabled(flagKey string, userCtx UserContext) bool
flagsdk.OnChange(flagKey string, callback func(Variant))
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Flag streaming service is unavailable | SDK cannot receive updates; stale flag state served | SDK uses last-known flag state from memory; evaluations continue with stale state; alert on-call when staleness exceeds 60s |
| Propagation bus drops a flag_changed event | flag state diverges between instances; inconsistency alert from periodic reconciliation | SDK periodically fetches a full snapshot (every 5 min) as a reconciliation mechanism |
| SDK bug causes incorrect evaluation for some users | feature metrics spike in A/B experiment; flag evaluations logged with unexpected distribution | canary SDK version before fleet-wide rollout; flag evaluation sampling with expected vs actual distribution |
| Kill switch activated but takes 30s to propagate | outage continues longer than expected | dedicated kill-switch topic with priority delivery; push rather than pull for kill-switch flags |

## Observability

- metric: flag propagation latency — time from flag_changed API call to SDK in-memory update across the fleet (sampled per flag, per environment)
- metric: SDK evaluation count per flag per second — identifies highest-traffic flags and aids debugging
- metric: flag staleness (time since last update received) per SDK instance — alerts when instances fall behind
- log: every flag state change with full before/after in the audit log
- SLO: 99% of flag changes propagate to 95% of SDK instances within 30 seconds; kill switches propagate within 10 seconds

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| In-process evaluation from local flag state | sub-millisecond evaluation; no network dependency on hot path | stale flag state if propagation is delayed; requires SDK in every language | RPC call to central flag service per evaluation — adds 1-5ms per request; creates flag service as critical dependency on every service's hot path |
| Push-based propagation (WebSocket/SSE) | low propagation latency (seconds); flag changes are immediately pushed | WebSocket connections at scale (50K instances) require connection management infrastructure | polling on a fixed interval — cheaper infrastructure but 30-second polling means 30-second staleness by default |
| Consistent hashing for percentage rollouts (hash(flag_key + user_id)) | the same user always gets the same variant for the same flag; no sticky session needed | changing the hash seed invalidates historical experiment assignments | random assignment at evaluation time — user sees different variants on different requests |

## Interview It

**Google framing:** "Design the internal feature flag service used by all Google product teams. It must evaluate flags on every search request without adding measurable latency." Expect pushback on SDK footprint, propagation guarantees, and experimentation integration.

**Cloudflare framing:** "Design a feature flag system that works at the edge. Flag evaluation happens at 200 POPs globally, and a kill switch must be effective at all POPs within 10 seconds." Expect questions on edge state management, push vs pull propagation at the edge, and consistency across POPs.

**Follow-ups:**
1. How would you support targeting rules based on real-time user attributes that change every request (e.g., current request's country, device type)?
2. How would you deprecate and garbage-collect old flags that are no longer referenced in code?
3. What happens if an engineer accidentally sets a kill switch on a flag that disables authentication for all users?
4. How would you add flag evaluation to a mobile app where the SDK cannot maintain a persistent connection?
5. How do you ensure that two flags used together in an A/B experiment do not produce contradictory targeting for the same user?

## Ship It

- `outputs/capacity-sheet-feature-flag-service.md`

## Exercises

1. **Easy** — Design the consistent hash function for percentage rollouts. Given flag_key = "new_checkout" and user_id = "user_12345", how do you deterministically assign a bucket 0–99 without a database lookup?
2. **Medium** — Design the reconciliation mechanism that detects and corrects divergence between SDK in-memory state and the flag store. How often does it run, and what is the minimum data needed to detect divergence efficiently?
3. **Hard** — Add multi-variate experimentation support. A flag now has 4 variants (A, B, C, D) with a 25% split. The experiment system must guarantee mutual exclusivity (a user cannot be in two conflicting experiments simultaneously). Design the assignment and audit model.

## Further Reading

- https://launchdarkly.com/blog/why-are-feature-flags-important/ — LaunchDarkly's engineering blog on feature flag architecture and the streaming propagation model they use at scale
- https://martinfowler.com/articles/feature-toggles.html — Martin Fowler's definitive taxonomy of feature flags (release, experiment, ops, permission toggles) — essential vocabulary for any interview
- https://engineering.fb.com/2023/06/26/developer-tools/introducing-gk/ — Facebook's Gatekeeper system, one of the early large-scale feature flag systems, with details on targeting and consistency
