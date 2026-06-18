# Code Deployment Pipeline

> Anyone can push code to production. The interview asks how you do it safely at 1000 engineers deploying 10 times per day, with automatic rollback, multi-region coordination, and zero planned downtime.

**Type:** Learn  
**Company focus:** Google  
**Learning goal:** Design a production deployment pipeline covering artifact storage, progressive rollout with automatic rollback triggers, multi-region deployment coordination, and deployment observability — not just "run kubectl apply."  
**Prerequisites:** `11-observability-slos-and-debugging/01-sli-slo-error-budget`, `13-multi-region-cdn-and-edge-traffic/01-active-active-vs-passive`, `10-reliability-retries-and-backpressure/03-circuit-breakers`  
**Estimated time:** ~75 min  
**Primary artifact:** capacity sheet

## The Problem

Design the deployment pipeline for a large software organization where 1000 engineers deploy code to thousands of production services, multiple times per day, across multiple geographic regions. The pipeline must store build artifacts durably, deploy progressively to catch regressions before they affect all traffic, roll back automatically when key metrics degrade, and coordinate multi-region rollouts without causing cross-region inconsistency.

Mid-level candidates describe a CI/CD pipeline that runs tests and then deploys. Senior candidates recognize that the hard problems are not in the build phase — they are in the rollout phase. A bug that causes 2% of requests to fail is invisible in test environments; progressive rollout with real traffic exposes it on 1% of production before it reaches 100%. The rollback decision must be automatic (humans are too slow), and the rollback trigger must use the right signal (error rate, not just deployment health check).

The multi-region coordination problem also catches most candidates. If you deploy US-East before US-West and the new version breaks a cross-region API contract, US-East may be serving new code while US-West still calls old-format payloads. Forward and backward compatibility in APIs becomes a hard operational requirement, not a best practice.

## Clarify

- Is this a monolithic service or a microservices environment? The number of independent deployable units changes the coordination model dramatically.
- What is the rollback time target — how quickly must automatic rollback complete after a regression is detected?
- Are there regulatory requirements for change management (e.g., financial services requiring 4-eyes approval before production deployment)?

If the interviewer does not specify, assume microservices at scale (1000+ services), rollback time target of under 5 minutes from detection to 100% rollback, and no strict regulatory change management (but audit logs are required).

## Requirements

### Functional

- Build and store immutable deployment artifacts (container images or binaries) with content-addressed storage.
- Deploy progressively: 1% of traffic → 10% → 50% → 100%, with automatic advancement only when health criteria pass.
- Monitor deployment health automatically using error rate, latency percentiles, and custom SLOs.
- Roll back automatically when health criteria fail during a rollout stage.
- Coordinate multi-region deployments with configurable ordering and blast radius limits.
- Maintain a full audit log of every deployment event: who initiated, what artifact, what stage, and what outcome.

### Non-functional

- Artifact storage: immutable, content-addressed, with 7-year retention for compliance.
- Deployment initiation to first traffic: under 5 minutes.
- Automatic rollback completion: under 5 minutes from regression detection.
- Audit log availability: 99.99% (legal requirement).

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Deployments per day | 10K deployments/day (1000 engineers × 10 deploys each) | sizes pipeline capacity and artifact storage write throughput |
| Artifact size | 500 MB avg container image; 1000 unique images/day = 500 GB/day new artifacts | determines artifact storage growth and CDN requirements for image pulls |
| Progressive rollout stages | 4 stages per deployment, 5 min bake time each | minimum 20 min from deploy to 100% rollout at healthy pace |
| Rollback event rate | 5% of deployments roll back = 500/day | sizes rollback signal processing and rollback execution capacity |
| Audit log volume | 10K deployments × 20 events each = 200K events/day | append-only log; fits in any time-series store |

## Architecture

```text
engineer
  -> git push -> CI system (build, test, lint)
  -> on success: publish artifact to artifact registry
     artifact registry: content-addressed (image digest = SHA256 of layers)
     stored in object store (S3-compatible), indexed by service/version/digest

deployment controller
  -> receives deploy request {service, artifact_digest, rollout_config}
  -> creates deployment record in deployment DB
  -> initiates rollout stage 1: canary (1% of traffic)
     -> updates load balancer routing weights
     -> waits N minutes (bake time)
     -> evaluates health criteria:
        - error rate delta (new vs baseline p-value test)
        - p99 latency regression
        - custom SLO burn rate
     -> if healthy: advance to stage 2 (10%)
     -> if unhealthy: trigger rollback, alert on-call

multi-region coordinator
  -> for each region in deployment plan:
     -> wait for previous region to complete (or fail fast)
     -> enforce minimum soak time between regions
     -> monitor cross-region health signals during deployment
     -> halt pipeline if any region fails

rollback executor
  -> receives rollback signal {service, target_artifact_digest}
  -> immediately reroutes 100% of traffic to previous artifact
  -> terminates new instances
  -> marks deployment record as rolled-back
  -> posts audit event

audit service
  -> append-only log of all deployment events
  -> queryable by service, time range, engineer, artifact, and outcome
```

## Data Model & APIs

Core entities:

```text
Artifact    { digest, service_id, version, build_time, build_url, size_bytes, storage_key }
Deployment  { deployment_id, service_id, artifact_digest, initiated_by, started_at,
              current_stage, status, rollout_config, region_plan }
Stage       { stage_id, deployment_id, stage_name, traffic_percent, started_at,
              ended_at, status, health_snapshot }
AuditEvent  { event_id, deployment_id, event_type, timestamp, actor, metadata }
RolloutConfig { canary_percent, stages[], bake_time_minutes, health_criteria,
                auto_advance, auto_rollback }
```

Key APIs:

- `POST /v1/deployments` — body: `{service_id, artifact_digest, rollout_config}`; returns `{deployment_id}`
- `GET /v1/deployments/{deployment_id}` — returns full deployment status including current stage and health
- `POST /v1/deployments/{deployment_id}/advance` — manually advance to next stage (if auto_advance is disabled)
- `POST /v1/deployments/{deployment_id}/rollback` — trigger immediate rollback to previous artifact
- `GET /v1/deployments/{deployment_id}/audit` — returns full audit log for the deployment

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Health metric spike during canary stage | error rate delta exceeds threshold within bake window | auto-rollback triggers within 60s of threshold breach; on-call paged simultaneously |
| Deployment controller crashes during a rollout | deployment record stuck in intermediate stage; health check times out | deployment controller is stateless; deployment DB is the source of truth; restart resumes from current stage |
| New artifact breaks cross-region API contract | services in advanced region fail when calling services in regions still on old version | deploy in backward-compatible stages; never deploy both sides of a breaking API change in the same deployment window |
| Artifact registry unavailable during peak deploy time | image pull fails for new instances; deployment stalls | artifact registry has multi-region replication; deploy from nearest replica; on-failure, halt new deployments but do not affect running instances |

## Observability

- metric: deployment stage duration and health check outcome per service — tracks rollout health and identifies slow-to-converge services
- metric: rollback rate per service and team — rising rollback rate signals code quality or testing coverage problems
- metric: time from commit to 100% production rollout (deployment cycle time) — key engineering velocity metric
- metric: artifact pull duration from registry during deployment — slow pulls extend deployment time and can cause instance start timeouts
- log: every stage transition with deployment_id, stage name, traffic percentage, health criteria results, and actor (auto or human)
- SLO: 99% of healthy deployments (no rollback triggered) complete within 30 minutes from initiation to 100% rollout

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| Automatic rollback on metric threshold | catches regressions before full traffic exposure; reduces MTTR to < 5 min | false positives from normal metric variance can cause unnecessary rollbacks | human-triggered rollback only — humans respond in 5–30 minutes while automatic detects in seconds |
| Content-addressed artifact storage (digest = hash of content) | immutable by definition; deduplication across services that share layers; rollback is always to a known artifact | immutability means you cannot patch an artifact in place; must rebuild and redeploy | mutable artifact tags (e.g. "latest") — a rollback is ambiguous if the tag was reassigned |
| Progressive rollout with automatic stage advancement | limits blast radius to current traffic percentage; healthy rollouts are fully automated | adds 20-30 minutes to total deployment time vs direct cutover | big-bang deployment — cheaper in engineering cost but risks 100% of traffic on every bug |

## Interview It

**Google framing:** "Design the internal deployment system for a large tech company where 1000 engineers each deploy multiple times per day. How do you prevent a bad deploy from taking down production?" Expect pushback on rollback trigger design, false positive rate, and how you handle multi-service breaking changes.

**Cloudflare framing:** "How would you deploy configuration changes to 200 edge POPs globally? What is your blast radius control model?" Expect questions on deployment ordering, POP-level health signals, and how to roll back a configuration change that is already propagated to 100 POPs.

**Follow-ups:**
1. How would you detect that a deployment caused a regression in a metric that takes 10 minutes to manifest (e.g., memory leak)?
2. How would you handle a deployment that must be rolled out to 10 interdependent services simultaneously?
3. What happens if the deployment controller itself needs to be deployed?
4. How would you implement feature branch deployments so engineers can test in production with real traffic?
5. What is the minimum change to your pipeline to comply with SOC 2 change management requirements?

## Ship It

- `outputs/capacity-sheet-code-deployment-pipeline.md`

## Exercises

1. **Easy** — Design the content-addressed artifact storage layer. What is the naming scheme? How do you handle garbage collection of artifacts more than 90 days old that are not referenced by any active deployment?
2. **Medium** — Design the health criteria evaluation algorithm. Given error rate baseline = 0.5% and canary error rate = 1.2%, is this a regression? What statistical test should you use, and what is the minimum canary duration before the test is meaningful?
3. **Hard** — A breaking API change must be deployed across 50 microservices. Design the rollout ordering that prevents any service from calling an incompatible version of a dependency during the rollout.

## Further Reading

- https://sre.google/sre-book/release-engineering/ — Google SRE book chapter on release engineering; authoritative source on progressive rollout and rollback patterns
- https://martinfowler.com/bliki/CanaryRelease.html — Martin Fowler's canonical definition of canary releases and their role in risk management
- https://spinnaker.io/concepts/ — Spinnaker's deployment pipeline concepts, showing how progressive delivery is implemented in practice at Netflix and Google
