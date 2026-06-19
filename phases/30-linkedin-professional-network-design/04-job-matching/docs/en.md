# Job Matching — Search & Skills Graph

> The best job recommendation is not the most popular listing — it's the one the member is most likely to apply for and get.

**Type:** Build
**Company focus:** LinkedIn
**Learning goal:** Design LinkedIn's job matching system that recommends 15M active listings to 950M members with sub-second relevance ranking.
**Prerequisites:** `17-search-crawl-and-monitoring-systems/02-search-autocomplete`, `30-linkedin-professional-network-design/02-professional-graph`
**Estimated time:** ~75 min
**Primary artifact:** job matching relevance model design

## The Problem

Design LinkedIn's job matching system. It must:

1. **Handle the scale mismatch** — 15M active listings × 950M members makes full pairwise scoring impossible; candidate generation must prune to <1K listings per member before the heavy ranker runs.
2. **Normalize skills** — members write "JS", recruiters write "JavaScript"; the system must unify these into a canonical skills graph.
3. **Rank with multi-signal relevance** — required-skills overlap, seniority fit, location preference, company preference, social proof, and historical apply rates all matter.
4. **Surface listings before they fill** — a job posted 2 hours ago should outrank an identical listing posted 3 weeks ago.
5. **Support Easy Apply** — jobs that allow one-click applying through LinkedIn convert at 3× the rate; this signal must be incorporated.

## Clarify

- What is the target latency for job recommendations? (Target: <500ms p99 for the recommendations API)
- Is location matching hard-filter or soft-preference? (Treat as soft preference with configurable radius)
- Are remote jobs in scope? (Yes — "Remote" is a location option, not a filter exclusion)
- Do we surface only open listings or also recent-closed ones? (Active listings only; TTL removes closed ones)
- Is personalization based on the member's current profile only, or also their browsing and apply history?

## Requirements

### Functional Requirements

- Return top-N ranked job listings for a given member in <500ms.
- Normalize member skills and job required-skills against the LinkedIn Skills Graph.
- Score each candidate listing using: skills intersection, seniority fit, location, company size preference, Easy Apply flag, and social proof (connections at that company).
- Support freshness boost: listings newer than 48 hours receive a recency multiplier.
- Expose an admin signal-weight configuration without code redeployment.

### Non-functional Requirements

- 15M active job listings, 950M members.
- Recommendations triggered for 30M DAU; peak load ~3,000 recommendation requests/sec.
- Index freshness: new job listing must appear in search within 60 seconds of publication.
- Listing expiry: expired listings removed from candidate set within 5 minutes of TTL.
- A/B testing: new ranking signal variants testable on 10% of members with apply-rate + hire-rate measurement.

## Capacity Model

| Dimension | Estimate | Detail |
|-----------|----------|--------|
| Active listings | 15M | partitioned by location + industry in Elasticsearch |
| Recommendation requests/sec | ~350 avg, ~3,000 peak | morning job-search window |
| Candidate generation output | top-1,000 listings/member | Elasticsearch query, inverted index on skills + location |
| Ranking model input | 1,000 candidates | feature vectors from Venice |
| Final recommendations | top-15 per surface | home feed + Jobs tab |
| Skills Graph nodes | ~36,000 canonical skills | maintained by LinkedIn editorial + ML taxonomy |
| Venice feature reads/request | ~200 | member features + job features |

## Architecture

### Two-Stage Retrieval + Ranking

```
Member Request (member_id)
    ↓
[1. Feature Extraction]  — pull member profile from Venice
   — skills (normalized), title, location, seniority_level
    ↓
[2. Candidate Generation]  — Elasticsearch query
   — inverted index: required_skills[], location, industry
   — boost: Easy Apply flag, listing age < 48h
   — returns top-1,000 job IDs
    ↓
[3. Feature Hydration]  — Venice batch lookup
   — for each job_id: required_skills, seniority_range,
     salary, company_size, apply_rate_for_similar_profiles
    ↓
[4. Ranking Model]  — weighted score per listing
   — skills_score (intersection / job requirements)
   — seniority_fit (in range: 1.0, adjacent: 0.5, out: 0.0)
   — location_score (same city: 1.0, same region: 0.7, remote: 0.9)
   — easy_apply_bonus (+0.15 if EasyApply=true)
   — social_proof (connections at company: log-scaled)
   — historical_apply_rate (members with similar profile applied here)
    ↓
[5. Business Rules Layer]
   — deduplicate same company (max 3 per company in top-15)
   — inject "promoted" listings at position 3, 7 (ad system)
    ↓
[6. Response]  — top-15 ranked listings returned
```

### LinkedIn Skills Graph

The Skills Graph is a canonical knowledge graph of ~36,000 professional skills maintained by LinkedIn's editorial and ML teams.

**Normalization:** "JS" → "JavaScript", "javascript" → "JavaScript", "node.js" → "Node.js". All member skills and job requirements are stored and queried in canonical form.

**Adjacency (inference):** The graph encodes which skills are adjacent. Python → {NumPy, Pandas, scikit-learn, PyTorch, Keras}. A member with Python + NumPy is inferred to be learnable-for ML-adjacent roles even without an explicit ML label.

**Seniority inference:** Title sequence ("Software Engineer" → "Senior Software Engineer" → "Staff Engineer") maps to seniority levels 1–5. The graph encodes typical title progressions per industry so seniority_level can be imputed for profiles that lack explicit labels.

**Drift handling:** New technologies (e.g., a new framework announced at a conference) enter the graph via an editorial ingestion pipeline that monitors LinkedIn's own skills section self-reports. Until a new skill reaches 10K member endorsements, it is treated as a synonym of its closest existing parent node to prevent cold-start underscoring.

### Storage Components

| Component | Role |
|-----------|------|
| Elasticsearch | Inverted index on skills, title, location, industry; candidate generation; TTL on listing expiry |
| Espresso | Primary job listing store (job_id → full listing JSON); source of truth |
| Venice | Precomputed feature vectors: member skills, historical apply rates, job normalized features |
| Pinot | Real-time OLAP: "top jobs in your metro this week", apply-rate analytics per listing |
| Kafka | Listing publish events → Elasticsearch sync (via Brooklin CDC), notification events |

### Notification Pipeline

```
New listing published
    → Kafka topic: job-listing-events
    → Samza consumer: matches against member saved-search preferences
    → Kafka topic: notification-candidates
    → Email service consumer: "Jobs matching your profile"
    → Push notification service
```

### A/B Testing

New ranking signals are gated behind an experiment framework. 10% of members are assigned to the treatment model variant. Primary metrics measured: apply-through rate (did they click Apply?), and hire signal (recruiter marked candidate as hired — fed back via recruiter workflow API). A signal is graduated to 100% rollout only if both metrics improve at p<0.05 over a 2-week window.

## Failure Modes

| Mode | Cause | Mitigation |
|------|-------|------------|
| Skills graph drift | New technology not yet in graph; member underscored for relevant roles | Synonym fallback to parent node; editorial fast-track for >1K self-reports |
| Listing freshness staleness | Expired listings still ranked (recruiter forgot to close) | TTL on Elasticsearch doc; Espresso CDC triggers deletion on status change |
| Cold start — new member | No skills data; fallback needed | Default to popular jobs in their location + industry; prompt to complete profile |
| Venice feature staleness | Offline job features not refreshed | Feature TTL alert; daily Spark job refreshes job features; alert on lag >24h |
| Elasticsearch shard hotspot | All new listings land on one shard (monotonic ID) | Partition listing index by (location_hash XOR industry_hash) |
| Ranking model silent degradation | New signal hurts hire rate but improves apply rate | Track both metrics; use hire rate as primary guardrail metric |

## Interview Trade-offs to Discuss

- **Candidate generation strategy:** Why Elasticsearch (inverted index on skills) instead of nearest-neighbor vector search? Answer: Skills are discrete categorical labels, not continuous embeddings; inverted index is faster and more explainable. Vector search adds value for semantic similarity ("project manager" ≈ "program manager") but adds latency.
- **Skills normalization vs. exact match:** Normalization improves recall (finds more relevant jobs) but reduces precision (may surface tangentially-related jobs). LinkedIn tunes this via the adjacency weight — adjacent-skill matches contribute 0.6× the score of exact matches.
- **Seniority as a hard filter vs. soft signal:** Hard-filtering by seniority range (SeniorityMin..SeniorityMax) improves precision but can exclude members who are slightly under/over-qualified and still good candidates. LinkedIn uses soft scoring (1.0 in-range, 0.5 adjacent) to allow adjacency.
- **Why Venice for feature serving and not Redis?** Venice is LinkedIn's own offline-first feature store optimized for batch writes from Spark + fast read serving. Redis would require manual TTL management and doesn't integrate with the Spark feature computation pipeline.
