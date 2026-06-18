# Chaos Engineering — Resilience by Design

> You do not discover that your fallbacks are broken during a chaos experiment. You discover it during a production incident. Chaos engineering just controls when.

**Type:** Concept + Build  
**Company focus:** Netflix  
**Learning goal:** Design a system to be resilient against deliberate fault injection. Understand the Simian Army (Chaos Monkey, Chaos Kong, Latency Monkey), FIT (Fault Injection Testing), bulkheads, and fallback chains. Learn how to design fallbacks before chaos reveals gaps.  
**Prerequisites:** `10-reliability-retries-and-backpressure/03-circuit-breakers`, `01-netflix-rubric`  
**Estimated time:** ~75 min  
**Primary artifact:** resilience design doc + fallback chain map  

## The Problem

Netflix runs in AWS across multiple regions. Any dependency can fail at any time. The system must continue serving subscribers through service failures, regional outages, elevated latency, and network partitions. To verify resilience before incidents find gaps, Netflix invented Chaos Engineering — deliberately injecting failures into production to test system behavior.

Design both the chaos engineering platform and the resilience patterns a system must implement to survive it.

## Clarify

- What types of failures should the chaos platform inject? (process kill, latency injection, network partition, region kill)
- What is the blast radius limit for each experiment? (single instance? single AZ? single region?)
- How do you decide when it is safe to run a chaos experiment? (traffic thresholds, time windows)
- What is the steady-state hypothesis? (what metrics prove the system is healthy?)
- How do you abort an experiment when the steady state degrades unexpectedly?

## Requirements

### Functional

- Kill random service instances (Chaos Monkey).
- Kill an entire AWS region (Chaos Kong).
- Inject artificial latency into service calls (Latency Monkey).
- Inject faults into specific dependency paths (FIT — Fault Injection Testing).
- Automatically abort if the steady-state hypothesis is violated.
- Record experiment results and correlate with system health metrics.

### Non-functional

- Blast radius: configurable per experiment (instance / AZ / region).
- Safety: automatic abort within 60 seconds of steady-state violation.
- Observability: full metric correlation between experiment events and system behavior.
- Scheduling: chaos runs during business hours (deliberately) and not during planned releases.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Services in fleet | 1,000+ microservices | experiment targeting complexity |
| Instances per service | 10–1,000 depending on service | blast radius calculation |
| Experiments per week | dozens to hundreds | scheduling and safety coordination |
| MTTR target | under 5 minutes for region failover | drives failover readiness requirements |
| Steady-state SLO | playback start success >99.9% | abort threshold for experiments |

## Architecture: The Simian Army

```text
chaos control plane
  -> experiment scheduler (selects targets, checks safety gates)
  -> blast radius calculator (instance / AZ / region scope)
  -> fault injector
     -> Chaos Monkey: kills random EC2 instances or containers
     -> Chaos Kong: disables an entire AWS region
     -> Latency Monkey: adds synthetic delay to service calls
     -> FIT: injects failures at specific request paths
  -> steady-state monitor (watches SLOs, triggers abort)
  -> result recorder (stores experiment metadata + outcome)
```

### Chaos Monkey

Chaos Monkey randomly terminates EC2 instances or containers within a service:

```text
target = random_select(instances(service=target_service, az=any))
terminate(target)
observe: does the service recover within SLO?
```

Key requirement: **services must be stateless or handle sudden termination gracefully** (no sessions stored in-process, no unsaved state).

### Chaos Kong

Chaos Kong simulates an entire region failure:

```text
disable_traffic_routing(region=us-east-1)
observe: does traffic shift to us-west-2 within N minutes?
observe: do subscribers notice? (measure playback errors)
restore after experiment duration
```

This requires multi-region active-active or active-standby readiness. Netflix's service mesh must reroute traffic automatically.

### Latency Monkey

Latency Monkey injects artificial latency into service-to-service calls:

```text
intercept calls: service_A -> service_B
add delay: +500ms to 50% of calls
observe: does service_A's circuit breaker open appropriately?
observe: does service_A's fallback activate?
```

This reveals whether circuit breakers and timeouts are correctly tuned.

### FIT — Fault Injection Testing

FIT injects failures at a specific path through the request graph, scoped to a percentage of traffic:

```text
inject: recommendation service returns 500 for 5% of homepage requests
observe: does the homepage fall back to trending content?
observe: does the playback start success rate remain above SLO?
```

FIT is more surgical than Chaos Monkey — it tests specific fallback paths without killing infrastructure.

## Bulkheads

A bulkhead isolates failure domains so that one failing dependency does not exhaust resources for all dependencies:

```text
thread_pool_A: dedicated to recommendation service calls (max 50 threads)
thread_pool_B: dedicated to playback service calls (max 200 threads)
```

Without bulkheads, a slow recommendation service can exhaust all HTTP client threads, blocking playback — which is far more critical.

Netflix uses Hystrix (now Resilience4j) for bulkhead and circuit breaker implementation.

## Fallback Chains

A fallback chain defines what the system should do when a dependency is unavailable:

```text
[homepage load]
  1. Personalized recommendations (primary)
  2. If slow (>150ms): serve pre-computed cache (EVCache)
  3. If cache miss: serve regional trending content
  4. If trending unavailable: serve global top-50 static list
  5. If everything fails: show search bar only
```

Each level must be tested independently with FIT before chaos experiments validate the chain end-to-end.

## Designing Fallbacks Before Chaos

The sequence matters:
1. **Define the steady-state hypothesis** before writing any code.
2. **Design the fallback chain** as part of service design, not as an afterthought.
3. **Write FIT tests** for each fallback level.
4. **Run Chaos Monkey** to verify instance recovery.
5. **Run Chaos Kong** to verify regional failover.

Teams that design fallbacks reactively (after chaos reveals gaps) are always late.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Chaos experiment breaks steady state unexpectedly | SLO monitor triggers abort | Automatic termination of experiment; alert on-call; post-mortem |
| Circuit breaker never opens during Latency Monkey | Fallback never activates despite injected latency | Timeout thresholds too high; retune circuit breaker sensitivity |
| Chaos Kong reveals single-region data dependency | Regional failover succeeds but some subscribers lose state | Identify and replicate the data dependency across regions |
| FIT reveals fallback returns wrong data | Fallback content is irrelevant or broken | Fix fallback implementation; add integration test for fallback path |
| Experiment runs during a planned release | Double failure: experiment + bad deploy | Safety gate: block chaos during active deployments |

## Observability

- metric: playback start success rate (primary steady-state signal)
- metric: homepage load latency at p99
- metric: circuit breaker state transitions (open/closed/half-open) per service
- metric: fallback activation rate per service and fallback level
- metric: chaos experiment duration, scope, and outcome
- log: every fault injection event with target, scope, and timestamp
- alert: steady-state violation during experiment → auto-abort + page on-call

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Run chaos in production, not staging | Tests the real system under real load | Risk of real subscriber impact | Staging environments never fully replicate production failure modes |
| Automatic abort on SLO violation | Limits blast radius of chaos experiments gone wrong | May abort valid experiments prematurely on noisy metrics | Manual abort requires human in the loop who may be slow |
| FIT over full Chaos Monkey for new systems | Surgical testing of specific paths without infrastructure disruption | More complex to set up per dependency | Infrastructure-level chaos reveals different failure modes; both are needed |
| Bulkheads per dependency | Failure isolation between service call paths | More thread pools to configure and monitor | Shared pools allow cascading failures |

## Interview It

**Netflix framing:** "How do you design a system to be resilient?" Strong answers discuss explicit fallback chains, circuit breakers, bulkheads, and the role of chaos experiments in validating those designs. Weak answers say "add redundancy" without discussing what the system does when a dependency actually fails.

**Follow-ups:**
1. How do you know when it is safe to run a Chaos Kong experiment?
2. What is the difference between a circuit breaker and a bulkhead?
3. How do you test your fallback chain without running chaos in production?
4. What metrics do you watch to know whether a chaos experiment is harming subscribers?
5. How do you get teams to design fallback chains before they are forced to?

## Ship It

- `outputs/design-doc-chaos-engineering.md`
- `outputs/fallback-chain-map.md`
- `outputs/interview-card-chaos-engineering.md`

## Exercises

1. **Easy** — Write a fallback chain for the Netflix search service with at least 3 levels.  
2. **Medium** — Design the safety gate that prevents chaos experiments from running during active deployments or when traffic is above 95th-percentile load.  
3. **Hard** — Design FIT injection for a specific path: recommendation service → feature store → EVCache. Write the steady-state hypothesis and abort conditions.  

## Further Reading

- [Principles of Chaos Engineering](https://principlesofchaos.org/)  
- [Netflix Chaos Monkey GitHub](https://github.com/Netflix/chaosmonkey)  
- [Chaos Engineering book (O'Reilly)](https://www.oreilly.com/library/view/chaos-engineering/9781492043850/)  
