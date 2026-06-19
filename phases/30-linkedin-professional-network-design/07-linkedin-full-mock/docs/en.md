# LinkedIn Full Mock Loop

**Type:** Mock Interview
**Company focus:** LinkedIn
**Learning goal:** Execute a complete 45-minute LinkedIn senior-level system design interview for PYMK (People You May Know) and self-score against a structured rubric.
**Prerequisites:** All lessons 01–06 in this phase
**Estimated time:** 45 min (interview) + 15 min (self-score)
**Primary artifact:** PYMK architecture design + scored rubric

## The Interview Prompt

> "Design LinkedIn's People You May Know (PYMK) feature. PYMK recommends up to 10 connections per member that they don't yet know but likely should. It appears on the LinkedIn home feed, the My Network tab, and in email digests. How would you design this system?"

This is a real interview question asked at LinkedIn for senior and staff engineer roles. Budget 45 minutes.

## Milestone Map

Use this to pace yourself during the mock. An interviewer is watching the clock.

| Time | Milestone | What you must cover |
|------|-----------|---------------------|
| 0–5 min | **Clarify** | Scope, scale, freshness, privacy constraints |
| 5–15 min | **Design the graph + capacity** | Edge store, capacity math, 2nd-degree traversal |
| 15–25 min | **Batch + online pipeline** | Spark batch job, Venice features, Espresso storage |
| 25–35 min | **Ranking + privacy** | Feature scoring model, blocked members, visibility |
| 35–45 min | **Failure + observability** | Connection event invalidation, Kafka pipeline, A/B |

## Clarify (minutes 0–5)

Questions to ask before designing:

- **What signals define "likely to know"?** (shared connections, same company, same school, same geography)
- **Is freshness important?** (Does a new hire at the member's old company appear in PYMK within minutes, or can it be a daily batch job?)
- **Privacy?** (Blocked members must never appear; visibility settings must be respected — some members hide their connections list)
- **What surfaces show PYMK?** (Home feed, My Network tab, email digest — each has different volume/latency requirements)
- **Scale?** (950M members, average 500 connections per member = 475B edges)

## Capacity Model

| Dimension | Estimate | Detail |
|-----------|----------|--------|
| Members | 950M | global |
| Avg connections per member | 500 | = 475B directed edges in the graph |
| 2nd-degree candidates per member | up to 250K | 500 connections × 500 each |
| PYMK impressions/day | ~30M | DAU seeing My Network or home feed |
| Pre-computed PYMK results | top-10 per member stored in Espresso | served at API layer without real-time computation |
| Batch job frequency | nightly (Spark) | full recompute of candidates |
| Online invalidation latency | <5 minutes | new connections → PYMK cache invalidation via Kafka |

## Architecture

### Graph Representation

LinkedIn's professional graph has ~950M members and ~475B edges. This does not fit in memory on a single machine.

**Storage layer:**
- **Espresso (document store):** each member_id maps to a document containing their adjacency list (sorted array of connection_ids). This is the authoritative edge store.
- **Voldemort (key-value store):** used for fast adjacency-list reads during graph traversal.
- **In-memory graph cache (heap-based):** the top 1% of members by connection count ("super nodes") have their adjacency lists cached in memory on graph processing workers.

**2nd-degree traversal:**
```
For member M:
  1st_degree = Espresso.get(M).connections  // up to 500
  candidates = {}
  for each friend F in 1st_degree:
    2nd_degree_of_F = Espresso.get(F).connections
    for each C in 2nd_degree_of_F:
      if C != M and C not in 1st_degree(M):
        candidates[C].shared_connections += 1
```

This naive traversal is O(degree²) — for a member with 500 connections where each connection also has 500, that's 250K candidate lookups. This is batch-computed by Spark, not done at request time.

### Batch Pipeline (Spark + Hadoop)

```
Nightly Spark Job:
  Input: full graph snapshot (Espresso dump → HDFS)
  Step 1: For each member M, compute 2nd-degree candidate set
  Step 2: For each candidate pair (M, C), compute features:
     — shared_connection_count (strongest signal)
     — same_company (current or past)
     — same_school (same graduation period ± 3 years)
     — same_location (same metro)
     — profile_similarity_score (Venice feature)
  Step 3: Score each candidate with ranking model
  Step 4: Keep top-100 candidates per member
  Output: write top-100 per member to Espresso
```

### Online Ranking (Request Time)

When a member visits My Network:

```
API request (member_id)
  → Read top-100 candidates from Espresso (precomputed)
  → Fetch live Venice features (member's current company, recent activity)
  → Re-rank top-100 with latest features
  → Apply privacy filters:
      — remove blocked members
      — remove members who hid themselves from PYMK
      — respect mutual visibility settings
  → Return top-10
```

### Venice Feature Store

Venice stores precomputed features for each member used in PYMK ranking:
- `pymk_shared_connection_score` — normalized count of shared connections
- `pymk_company_overlap_score` — current company match + past company match (weighted)
- `pymk_school_overlap_score` — school + graduation year proximity
- `pymk_engagement_probability` — predicted probability of connection acceptance based on historical data

These features are updated daily by the batch Spark job and propagated to Venice via Brooklin CDC.

### Kafka Event Pipeline for Invalidation

When member A connects with member B, this event invalidates both A's and B's PYMK cache (the new connection should not appear as a PYMK suggestion anymore):

```
Connection accepted
  → Kafka topic: connection-events
  → Samza consumer: pymk-invalidation
      — Remove B from A's PYMK candidate list in Espresso
      — Remove A from B's PYMK candidate list in Espresso
      — Trigger partial recompute for both A and B (online incremental update)
  → Latency: <5 minutes for invalidation
```

### Privacy and Compliance

**Blocked members:** LinkedIn maintains a block list per member. Members who have blocked each other must never appear in each other's PYMK. Block lists are loaded into memory on API servers (small enough — average member blocks <10 accounts).

**Visibility settings:** Members can set "hide my connections from other members" and "exclude me from PYMK." These are enforced at the online ranking layer as hard filters before the top-10 is returned.

**GDPR right-to-be-forgotten:** When a member deletes their account:
1. Espresso deletes their adjacency list (graph node + edges).
2. A Kafka event triggers a cleanup job that removes them from all other members' precomputed PYMK candidate lists in Espresso.
3. The cleanup must complete within 30 days (GDPR Article 17).

**Mutual visibility:** If member A has hidden their connections from the public, and B's 2nd-degree path to C goes through A's connections, C should not be recommended to B via that path. This requires tracking the "path source" during candidate generation and filtering paths through private connection lists.

### A/B Testing PYMK

PYMK model changes are tested via the experiment framework:
- **Treatment group (10% of members):** new ranking model weights
- **Control group:** existing model
- **Primary metric:** connection acceptance rate (did the member click Connect on the PYMK suggestion?)
- **Guardrail metric:** 7-day retention (did engagement with LinkedIn decrease for treatment members?)
- **Graduation criteria:** +2% connection acceptance rate at p<0.05 over 14 days

## Strong-Hire vs. Weak-Hire Patterns

| Behavior | Strong-Hire | Weak-Hire |
|----------|-------------|-----------|
| Clarification | Asks about privacy, scale, freshness, blocked members | Jumps to design without clarifying |
| Capacity | Derives 475B edges, explains super-node problem | States "big graph" without numbers |
| Pipeline | Distinguishes batch vs. online, explains why each | Designs only one layer |
| Privacy | Proactively mentions blocked members and GDPR | Only mentions if prompted |
| Failure | Discusses PYMK cache staleness, Kafka consumer lag | Doesn't discuss failure modes |
| Observability | Proposes connection acceptance rate + CTR as metrics | No metrics mentioned |
| LinkedIn-specifics | Mentions Venice, Espresso, Kafka, Samza by name with context | Uses generic "database" and "cache" |

## Scoring Rubric

| Dimension | Weight | Full Credit (100pts) |
|-----------|--------|---------------------|
| Graph Design | 20 | Edge store choice, 2nd-degree traversal algorithm, super-node handling |
| Pipeline Architecture | 20 | Batch Spark job + online re-ranking, Venice features, Espresso storage |
| Ranking/ML Design | 15 | Feature taxonomy, model update cadence, how freshness is handled |
| Privacy/Compliance | 15 | Blocked members, visibility settings, GDPR delete flow |
| Failure Recovery | 15 | Kafka invalidation latency, stale cache handling, connection event ordering |
| Observability | 15 | Connection acceptance rate, Pinot analytics, lag monitoring |

**Total: 100 points**

Score thresholds:
- **85+:** Strong Hire — ready for senior/staff role
- **70–84:** Hire — strong fundamentals, some gaps
- **55–69:** No Hire (yet) — needs improvement in 1-2 dimensions
- **<55:** No Hire — significant gaps across multiple dimensions

## Follow-up Questions

**1. How do you handle GDPR right-to-be-forgotten?**
Expected answer: Delete member's adjacency list in Espresso (graph node), publish a Kafka deletion event, trigger async cleanup job to remove them from all other members' PYMK candidate lists. Complete within 30 days. Audit log in the deletion pipeline.

**2. How do you A/B test PYMK model changes?**
Expected answer: Experiment framework assigns members to treatment/control at login. Both groups get PYMK computed with their respective model. Measure connection acceptance rate as primary metric over 14 days. Graduate if statistically significant improvement at p<0.05. Use guardrail metric (7-day retention) to prevent optimizing against LinkedIn's core engagement.

**3. What if a member has 50,000 connections (super-node)?**
Expected answer: 50K connections × 500 average = 25M 2nd-degree candidates. Naive batch computation is too expensive. Solutions: (a) cap 2nd-degree traversal at the top-1K connections by affinity score; (b) pre-cache super-node adjacency lists in memory; (c) use approximate nearest-neighbor search on member embeddings instead of graph traversal for super-nodes.

**4. How do you handle the cold-start problem for a new member with 0 connections?**
Expected answer: No 2nd-degree candidates exist. Fall back to: (a) same company members from their profile, (b) same school members, (c) geographic proximity, (d) popular members in their industry vertical. These are computed without graph traversal and served as "people from your industry."

**5. How does PYMK scale to 950M members?**
Expected answer: Partition the batch Spark job by member_id hash — each partition handles a slice of the member space. Espresso is horizontally sharded by member_id. Venice feature reads scale with more read replicas. The nightly batch job is the primary scaling concern — parallelism across hundreds of Spark executors, output written to HDFS and then bulk-loaded to Espresso.
