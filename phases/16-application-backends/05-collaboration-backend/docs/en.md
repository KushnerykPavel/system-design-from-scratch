# Collaborative Document / Presence Backend

> Real-time collaboration is a state-convergence problem first and a text-editing problem second.

**Type:** Build  
**Company focus:** Google  
**Learning goal:** Design a collaborative editing backend by reasoning about convergence, session routing, snapshotting, and presence without overselling impossible consistency guarantees.  
**Prerequisites:** `08-consistency-replication-and-transactions/01-consistency-spectrum`, `07-queues-streams-and-workflows/04-dlq-and-replay`, `16-application-backends/03-chat-system`  
**Estimated time:** ~75 min  
**Primary artifact:** collaboration session validator  

## The Problem

Design the backend for a collaborative document editor with live updates, cursors or presence, and durable document history. Many answers stay at "use WebSockets." Strong answers explain the document operation model, how edits converge, when snapshots are created, and what the system promises during reconnects or region failures.

The purpose of this lesson is not to turn every learner into a CRDT researcher. It is to force clarity on whether the system uses operational transforms, CRDTs, or a simpler single-writer coordinator and what trade-offs that choice creates.

## Clarify

- Is the product optimized for text documents, structured docs, whiteboards, or code?
- Are simultaneous edits from many users common, or are most sessions small?
- Do we need offline editing with merge on reconnect, or only live collaborative sessions?

If unspecified, assume text or structured docs, small active sessions, and live collaboration with short reconnect windows rather than days of offline divergence.

## Requirements

### Functional

- Multiple users can edit the same document concurrently.
- Participants can observe presence and near-real-time updates.
- The service preserves document history and durable snapshots.
- Reconnecting users can catch up without downloading every historical operation forever.

### Non-functional

- Keep local collaboration latency subjectively real-time.
- Prevent document corruption after acknowledged operations.
- Bound replay cost for long-lived active documents.
- Support graceful degradation when live collaboration is impaired.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Concurrent active docs | 1M | drives session routing and coordinator load |
| Peak ops/sec | 2M operations/s | shapes append log and merge pipeline |
| Typical collaborators | 2 to 8 per doc | keeps the common case grounded |
| Snapshot cadence | every 500 ops or 30 seconds | balances replay against storage cost |
| Presence heartbeat | every 10 to 20 seconds | affects ephemeral state scale |

## Architecture

```text
editor client
  -> session gateway
  -> doc coordinator / shard owner
  -> operation log
  -> transform / merge engine
  -> snapshot store
  -> presence service
```

Useful framing:

1. Live operations flow through a session owner or shard that imposes a convergence rule.
2. Durable history is append-first; snapshots cap replay cost.
3. Presence is ephemeral and must not be confused with durable document truth.

## Data Model & APIs

Core entities:

```text
document(doc_id, latest_snapshot_ref, version)
operation(doc_id, op_id, author_id, base_version, payload, applied_version)
presence(doc_id, user_id, cursor, expires_at)
snapshot(doc_id, version, blob_ref, created_at)
```

APIs:

- `POST /v1/docs/{id}/ops`
- `GET /v1/docs/{id}/sync?since_version=...`
- `POST /v1/docs/{id}/presence`
- `GET /v1/docs/{id}/snapshot`

Candidates should explicitly say whether version conflicts are transformed, merged commutatively, or rejected for retry.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| doc coordinator fails mid-session | session reconnect surge and ownership failover metrics | lease-based ownership, replay from operation log |
| operation log grows too large for hot docs | replay duration and snapshot staleness | regular snapshots and archival compaction |
| stale presence shows ghost collaborators | heartbeat expiry drift and stale-user reports | short TTL presence and explicit disconnect cleanup |
| bad transform logic corrupts doc state | divergence checks against snapshot replay | deterministic validation and emergency rollback to known snapshot |

## Observability

- metric: end-to-end operation apply latency
- metric: snapshot age and replay length per active document
- metric: session reconnect recovery time
- metric: presence heartbeat expiry skew
- log: rejected or transformed operations with version context
- trace: operation from client submit through durable append to fanout
- SLO: acknowledged operations remain durable and replayable while active sessions meet the collaboration latency target

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| session-owner coordination | simpler convergence reasoning | hotspot risk on very active docs | fully decentralized multi-writer merge with no central sequencing |
| regular snapshots | bounded replay cost | extra storage and background work | replay full operation history forever |
| ephemeral presence store | cheap and scalable presence | occasional ghost states | strongly consistent presence for every cursor update |

## Interview It

**Google framing:** "Design a collaborative document backend with real-time editing." Expect questions on convergence choice, snapshotting, and reconnect replay.

**Cloudflare framing:** "Design a low-latency collaboration path with region-local sessions." Expect pressure on session steering, cross-region failover, and what state must remain authoritative.

**Follow-ups:**
1. What changes if the product adds offline edits?
2. How would you support comments, suggestions, or change review?
3. What if one document becomes massively popular in a live event?
4. How do you verify the transform engine is not corrupting documents?
5. What changes for code editing where operation semantics are trickier?

## Ship It

- `outputs/design-review-collaboration-backend.md`

## Exercises

1. **Easy** — Pick between OT, CRDT, or coordinated sequencing for the default document type and justify it.
2. **Medium** — Design snapshot and compaction rules for docs that stay active for days.
3. **Hard** — Extend the system for offline editing with week-long divergence windows.

## Further Reading

- [Operational Transformation](https://research.google/pubs/pub36726/) — foundational intuition for concurrent text editing  
- [CRDTs: Consistency without concurrency control](https://hal.inria.fr/inria-00555588/document) — useful for understanding commutative merge approaches  
