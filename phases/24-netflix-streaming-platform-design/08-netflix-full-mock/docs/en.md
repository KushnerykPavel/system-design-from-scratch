# Netflix Full Mock Loop

> A 45-minute Netflix interview is not a design session. It is a live demonstration that you can navigate ambiguity, communicate trade-offs clearly, and know when to go deep vs stay broad.

**Type:** Mock Interview  
**Company focus:** Netflix  
**Learning goal:** Simulate a full 45-minute Netflix-style system design interview. Practice the milestone cadence, common follow-up pivots, and the failure modes that separate strong hires from misses.  
**Prerequisites:** All lessons 01–07 in this phase  
**Estimated time:** ~90 min  
**Primary artifact:** completed mock design doc + self-assessment against rubric  

---

## Interview Prompt

**"Design Netflix's video streaming system end to end."**

*The interviewer gives you nothing more. The clarification round is yours to drive.*

---

## How to Use This Lesson

1. Set a 45-minute timer.
2. Work through the prompt independently using only a whiteboard or blank document.
3. After the timer ends, compare your design against the milestone map and failure mode list below.
4. Score yourself against the rubric.
5. Read the follow-up pivots and answer each one in writing (5 minutes each).

---

## Expected Strong-Hire Milestone Map

| Time mark | Expected progress | Failure mode if missed |
|-----------|-------------------|------------------------|
| **5 min** | Asked ≥3 clarifying questions. Gave rough capacity numbers (subscribers, streams, CDN bandwidth). Scoped the problem to one primary question. | Jumped straight to architecture without clarifying. Numbers never appeared. |
| **15 min** | Sketched high-level system with named components: client → Open Connect CDN → origin → encoding pipeline. Identified the three hardest sub-problems. | Drew a generic "CDN + S3" diagram with no Netflix-specific context. Did not name sub-problems. |
| **25 min** | Drilled into the video pipeline or CDN layer (whichever the interviewer is most interested in). Named: encoding ladder, HLS segmentation, ABR algorithm, or OCA fill hierarchy. Discussed one failure mode with a specific mitigation. | Stayed at high level the entire time. No depth on any subsystem. No failure modes. |
| **35 min** | Handled the interviewer's deliberate pivot (see below). Adjusted design without losing track of overall structure. Showed awareness of trade-offs in the new direction. | Got flustered by the pivot. Lost thread of overall design. Said "that's a good question" without answering. |
| **45 min** | Summarized: top 3 trade-off decisions made and why. Named what to instrument first. Named one thing you would do differently with 2x the time. | No summary. Ended mid-thought. Did not reflect on trade-offs. |

---

## Clarifying Questions a Strong Candidate Asks First

1. "Is this VOD only, or does it include live streaming? They have very different latency requirements."
2. "What is the device matrix? Different devices support different codecs, which affects the encoding pipeline."
3. "Is the primary concern the delivery layer, the encoding pipeline, or the recommendation layer? I can go deep on any of these."
4. "What does 'end to end' mean here — from studio delivery, or from subscriber pressing play?"
5. "What is the failure tolerance? Is brief degradation acceptable (720p instead of 4K), or is playback interruption never acceptable?"

---

## Design Sketch: What to Cover

### Layer 1: Ingestion & Encoding

```text
studio file (raw video)
  -> ingest service (checksum validation)
  -> job scheduler (splits into per-segment jobs)
  -> encoding workers (parallel, per codec)
     - encoding ladder: 10 profiles from 235 kbps to 5800 kbps
     - codecs: H.264, H.265, AV1
  -> packager (HLS segmentation, DRM encryption, manifest generation)
  -> S3-like object store (multi-region replication)
```

### Layer 2: CDN Delivery (Open Connect)

```text
S3 origin
  -> Netflix regional POP (reactive fill + active-active)
  -> OCA servers embedded in ISP networks
     - proactive prefill: top-K titles pushed during nightly window
     - reactive fill: cache miss escalates to regional POP
     - cache eviction: LFU per title, ISP-local popularity ranking
  -> subscriber device
```

### Layer 3: Client ABR

```text
client player
  -> fetch master manifest (codec, resolution list)
  -> fetch variant manifest (segment URLs)
  -> measure bandwidth (EWMA over recent segment download times)
  -> select bitrate profile (conservative up, aggressive down)
  -> fetch segment, decode, render
  -> loop: re-evaluate bitrate every N segments
```

### Layer 4: Recommendation & Personalization

```text
playback events
  -> Kafka -> Flink (short-term signal extraction)
  -> EVCache feature store (subscriber preferences updated in minutes)
  -> daily Spark + GPU retraining of two-tower model
  -> ANN index updated
  -> recommendation API (200ms SLA, EVCache fallback, trending fallback)
  -> homepage row ranking
```

### Layer 5: Reliability

```text
per-service: circuit breakers (Hystrix/Resilience4j), bulkheads (per-dependency thread pools)
fallback chains: personalized → pre-computed cache → trending → static top-50
chaos engineering: Chaos Monkey (instance kill), FIT (dependency injection), Chaos Kong (region kill)
multi-region: traffic routing via DNS/Anycast; region failure does not require manual intervention
```

---

## Failure Modes to Demonstrate

Name at least three of these without prompting:

| Failure | What strong candidates say |
|---------|---------------------------|
| CDN miss during major launch | "Proactive prefill prevents cold cache, but if we miss a viral moment, the fill path is regional POP → origin. The origin can handle the load because it is not serving segments directly — only OCA misses do." |
| Recommendation service slow | "We have pre-computed recommendations in EVCache with sub-millisecond latency. Below that, we fall back to regional trending. The homepage always loads." |
| Encoding worker crashes mid-job | "Jobs are idempotent and checkpointed at the segment level. The scheduler reschedules the failed segment; we never re-encode the full title." |
| Region goes dark | "Traffic reroutes via DNS health checks to the nearest healthy region. We have verified this with Chaos Kong. The main risk is single-region data dependencies — we audit for these quarterly." |
| A/B experiment leaks | "The assignment service uses deterministic hashing. Holdout groups detect contamination. We enforce mutual exclusion for experiments on the same UI surface." |

---

## Interviewer Pivots (Follow-Up Challenges)

Practice answering each in under 5 minutes:

**Pivot 1: Live streaming**  
"Now assume 20% of Netflix's catalog is live events. What changes?"

*What to cover:* Segment duration drops to 2 seconds. No proactive prefill (content is being created). Origin is a live encoder, not S3. ABR parameters change (lower buffer target). DVR window is a ring buffer at edge. Multi-CDN is more critical (redundancy for live failover).

**Pivot 2: Cost optimization**  
"Your CDN costs are 3x budget. How do you reduce them without degrading quality?"

*What to cover:* Per-title encoding ladder (reduce bitrate for simple content). AV1 codec (40% smaller files at same quality, but higher encoding cost). Better prefill (fewer reactive fetches from origin). Longer segment TTLs (reduce re-fetch frequency). Reduce overprovision in lower-traffic regions.

**Pivot 3: Personalized thumbnails**  
"We want to show each subscriber a different thumbnail for the same title. How does the system change?"

*What to cover:* Thumbnails are a CDN-cached asset — they must be personalized without breaking CDN cacheability. Options: (a) personalized CDN URL with user-specific segment (breaks cache), (b) thumbnail selection at the edge (Cloudflare Worker-style), (c) thumbnail metadata served from recommendation API (client picks the right URL from a set). Discuss cache hit rate trade-off.

**Pivot 4: Sub-second recommendation refresh**  
"After a user finishes a show, we want recommendations to update in under 1 second. What changes?"

*What to cover:* Current path is minutes (Kafka → Flink → EVCache). For sub-second: synchronous recommendation re-fetch triggered by playback-end event on the client. Or: push update from server via WebSocket/SSE after processing the event. Trade-off: higher infrastructure cost, more complexity, but better user experience at season-ending moments.

**Pivot 5: Global content regulation**  
"A government requires that content in Country X be removed within 1 hour. How does the system handle it?"

*What to cover:* Invalidation must propagate to all OCA servers in Country X within 1 hour. Push invalidation message to all OCAs tagged with the country. OCAs acknowledge deletion. Manifests are invalidated at the CDN layer (versioned URLs expired, cache purge). Content catalog blocks the title for affected regions. Audit log for compliance.

---

## Rubric: Score Yourself

| Dimension | Strong hire (2 pts) | Hire (1 pt) | No hire (0 pts) |
|-----------|---------------------|-------------|-----------------|
| Clarification | ≥3 targeted questions, numbers provided | 1-2 questions, partial numbers | Jumped to architecture |
| Breadth | Named all 4 major layers with correct flow | Named 3 of 4 | Only 1-2 layers |
| Depth | Went deep on at least 1 subsystem with data model and failure modes | Went deep on 1 subsystem but no failure modes | Stayed high-level throughout |
| Failure modes | Named ≥3 failure modes with specific mitigations | Named 1-2 failure modes | No failure modes |
| Pivot handling | Answered pivot question without losing design thread | Answered pivot but lost track of overall design | Could not answer pivot question |
| Trade-offs | Named ≥2 trade-offs with rejected alternatives | Named 1 trade-off | No trade-offs mentioned |
| Netflix vocabulary | Used ≥2 Netflix-specific terms correctly | Used 1 correctly | No Netflix context |
| Summary | Summarized top trade-offs and next steps | Partial summary | No summary |

**Maximum score: 16 points**  
- 14–16: Strong hire  
- 10–13: Hire  
- 6–9: Mixed — likely needs another round  
- 0–5: No hire  

---

## Common Mistakes to Avoid

1. **Describing a generic CDN without explaining Open Connect.** Netflix's CDN is the differentiator in the design.
2. **Skipping capacity numbers.** Interviewers use these to gauge engineering calibration.
3. **One layer only.** Spending 40 minutes on encoding and never discussing delivery or recommendations.
4. **No failure modes.** Netflix interviewers care deeply about this.
5. **Getting stuck on the pivot.** The pivot is deliberate. Have a mental model flexible enough to adapt.
6. **Saying "we can scale horizontally" without explaining what that means for the specific component.**
7. **Not summarizing.** A crisp closing that ties trade-offs together is a strong signal.

---

## Ship It

- `outputs/mock-design-doc-netflix-full.md`
- `outputs/self-assessment-rubric-netflix-full.md`
- `outputs/pivot-answers-netflix-full.md`

## Further Reading

- [Netflix Tech Blog](https://netflixtechblog.com/) — review any 3 architecture posts before the real interview  
- [System Design Interview Vol. 2 (Alex Xu)](https://www.amazon.com/System-Design-Interview-Insiders-Guide/dp/1736049119) — Chapter on video streaming  
- [Designing Data-Intensive Applications (Kleppmann)](https://www.oreilly.com/library/view/designing-data-intensive-applications/9781491903063/) — foundational for stream processing and data pipelines  
