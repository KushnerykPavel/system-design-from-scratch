# Chat System

> Chat design is about delivery contracts, connection state, and backpressure more than about storing strings.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Design a real-time chat backend by separating connection handling, message durability, ordering guarantees, and offline delivery from the UI story.  
**Prerequisites:** `13-multi-region-cdn-and-edge-traffic/04-traffic-steering`, `07-queues-streams-and-workflows/02-delivery-semantics`, `16-application-backends/02-news-feed`  
**Estimated time:** ~90 min  
**Primary artifact:** delivery topology validator  

## The Problem

Design a chat system that supports one-to-one and small-group conversations, low-latency delivery to online users, and reliable catch-up for offline users. Interview answers often skip the most important question: what exactly does "message delivered" mean?

This lesson forces precision around ordering, retries, connection fanout, attachments, read receipts, and how presence interacts with delivery.

## Clarify

- Are we designing direct messages, large channels, or both?
- What delivery semantics matter: at least once, exactly once at the UI, or durable send acknowledgment?
- Are users expected to roam across devices and regions while staying connected?

If open-ended, assume direct messages and small groups, durable send acknowledgment before "sent" is shown, and best-effort per-conversation ordering with client-side dedupe.

## Requirements

### Functional

- Send and receive messages in active conversations.
- Sync missed messages to offline or reconnecting devices.
- Show delivery state such as sent, delivered, and optionally read.
- Maintain lightweight presence for routing and UX hints.

### Non-functional

- Keep median end-to-end delivery under a few hundred milliseconds for online users.
- Avoid message loss after durable acknowledgment.
- Prevent one noisy conversation from starving others.
- Make reconnection and replay safe under retries.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Concurrent connections | 20M online sockets | shapes connection gateway fleet and sharding |
| Message ingress | 700K messages/s peak | drives durable log and fanout throughput |
| Average conversation size | 2 to 20 participants | helps keep the core problem bounded |
| Offline sync window | 30 days | affects message retention and compaction |
| Attachment ratio | 15% of messages | separates metadata path from blob storage |

## Architecture

```text
client
  -> connection gateway
  -> chat session router
  -> conversation log / durable append
  -> delivery fanout
     -> online device push
     -> offline inbox / sync cursor
  -> presence service
  -> attachment store
```

Key choices:

1. Acknowledge sends after the durable conversation log commit, not after every device receives the message.
2. Treat online push and offline replay as separate delivery paths.
3. Use client-generated message IDs plus server dedupe to survive retries cleanly.

## Data Model & APIs

Core entities:

```text
conversation(conversation_id, members, policy)
message(message_id, conversation_id, sender_id, body_ref, created_at, seq_no)
device_cursor(device_id, conversation_id, last_seen_seq_no)
presence(user_id, device_id, region, expires_at)
```

APIs:

- `POST /v1/conversations/{id}/messages`
- `GET /v1/conversations/{id}/messages?cursor=...`
- `POST /v1/presence/heartbeat`
- `POST /v1/conversations/{id}/acks`

If exactly-once UI semantics are required, be careful to describe dedupe keys and replay behavior instead of claiming the network magically provides exactly once.

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| connection gateway partition isolates some clients | connection drop surge and region-local disconnect rate | reconnect to alternate gateway, cursor-based replay |
| presence is stale and routes delivery to dead sockets | heartbeat expiry skew and undelivered push count | short TTL presence with fallback to offline inbox |
| one group chat bursts and overloads fanout workers | per-conversation enqueue lag and worker saturation | shard delivery queues, backpressure noisy conversations |
| duplicate sends during mobile retry | duplicate message ID rate | idempotent server insert keyed by client message ID |

## Observability

- metric: durable-send acknowledgment latency
- metric: online delivery latency and reconnect replay latency
- metric: duplicate message ID suppression count
- metric: active connections and heartbeat expiry drift
- log: message state transitions for sampled conversations
- trace: send path from gateway to durable append to device fanout
- SLO: acknowledged messages are durably stored and available for replay within the documented recovery window

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| ack on durable log commit | clean send contract and strong replay safety | recipient may not see it instantly | ack only after recipient device receipt |
| soft ordering per conversation | practical and scalable | cross-device edge cases remain | global strict ordering |
| short-lived presence with heartbeats | simple routing hint | stale state during failures | fully authoritative presence for delivery correctness |

## Interview It

**Google framing:** "Design a chat backend for direct messages and small groups." Expect follow-ups on durable acknowledgment, replay, dedupe, and mobile reconnect behavior.

**Cloudflare framing:** "Design a globally distributed real-time messaging path with region-local sockets." Expect pressure on connection routing, failover, and how much state can live close to the edge.

**Follow-ups:**
1. What changes when large channels are added?
2. How do you handle message edits or deletes?
3. What if legal policy requires regional data residency?
4. How do you keep presence from becoming a correctness dependency?
5. What changes if attachments dominate traffic?

## Ship It

- `outputs/failure-checklist-chat-system.md`

## Exercises

1. **Easy** — Define exactly what the sender sees when a message is "sent."
2. **Medium** — Design how device cursors work across three active devices for one user.
3. **Hard** — Extend the design to end-to-end encrypted group messaging with server-side fanout still intact.

## Further Reading

- [The Log: What every software engineer should know about real-time data's unifying abstraction](https://engineering.linkedin.com/distributed-systems/log-what-every-software-engineer-should-know-about-real-time-datas-unifying) — strong mental model for durable append and replay  
- [Designing Data-Intensive Applications](https://dataintensive.net/) — good framing for ordering, idempotency, and delivery semantics  
