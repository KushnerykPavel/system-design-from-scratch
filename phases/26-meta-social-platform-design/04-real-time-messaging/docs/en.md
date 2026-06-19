# Real-Time Messaging — WhatsApp & Messenger

> A message is either delivered exactly once or the user retries — there is no in-between.

**Type:** Build
**Company focus:** Meta
**Learning goal:** Design a real-time messaging system that handles 100B messages/day with end-to-end delivery guarantees and offline storage.
**Prerequisites:** `07-queues-streams-and-workflows/02-delivery-semantics`, `16-application-backends/03-chat-system`
**Estimated time:** ~90 min
**Primary artifact:** message delivery state machine

## The Problem

WhatsApp handles over 100 billion messages per day across 2 billion users. Messenger handles 20 billion messages per day. Both must guarantee that every message is either delivered exactly once to the recipient device or is durably stored until the device reconnects — no silent loss, no duplicates visible to the user.

At Meta's scale this means sustaining roughly 1 million concurrent WebSocket connections per datacenter, handling offline message queues for billions of users, and propagating delivery acknowledgements back to senders in real time — all while running end-to-end encryption so the server never sees plaintext message content.

Design this system.

## Clarify

- Is end-to-end encryption required? (Affects what the server can store and index.)
- What is the maximum message size? (WhatsApp: 100 MB for media, unlimited text in practice.)
- What offline retention period is required before a message is considered undeliverable?
- Are group messages supported? What is the maximum group size?
- Is presence (online/offline, last-seen) a required feature?

## Requirements

### Functional

- Send a message from one user to another in real time (sub-second delivery when both are online).
- Store messages durably when the recipient is offline; deliver on reconnect.
- Report delivery status back to the sender: SENT, DELIVERED, READ.
- Support group messaging (fan-out to all members).
- Show presence: online/offline and last-seen timestamp.

### Non-functional

- Message delivery latency: p99 under 500ms when both users are online.
- Offline storage durability: messages retained for 30 days before expiry.
- Scale: 100B messages/day (WhatsApp); ~1M concurrent WebSocket connections per datacenter.
- End-to-end encryption: server stores only encrypted blobs; plaintext never leaves devices.
- Exactly-once delivery semantics (from the user's perspective): no visible duplicates.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Messages per day | 100B (WhatsApp) | throughput requirement for storage and gateway |
| Messages per second | ~1.2M peak | gateway and message store write throughput |
| Concurrent WebSocket connections | 1M per datacenter | gateway server sizing |
| Avg message size (text) | 200 bytes | storage: ~20TB raw/day before replication |
| Offline queue depth per user | up to 30 days of messages | HBase/Cassandra partition sizing |
| Group size limit (WhatsApp) | 1,024 members | fan-out upper bound per group message |

## Architecture

```text
sender device
  -> WebSocket gateway (Ejabberd/Erlang node or similar, stateful)
       -> message store (HBase/Cassandra: durable queue per user)
       -> recipient gateway lookup (which gateway node owns recipient's connection?)
            |-- recipient online: push to recipient gateway -> recipient device
            |-- recipient offline: message stays in store; device polls on reconnect
       -> delivery ack path:
            recipient device -> ack -> recipient gateway -> message store (mark DELIVERED)
            message store -> Iris/MQTT event -> sender gateway -> sender device (tick icon)

read receipt path:
  recipient opens conversation -> READ ack -> message store (mark READ)
  -> event -> sender device (double tick)
```

### Connection Layer

Each user device holds a persistent WebSocket (or MQTT on mobile) connection to an assigned gateway node. The gateway is stateless with respect to message content but stateful with respect to connection routing:

- **Ejabberd** (Erlang/OTP) is WhatsApp's original gateway technology. Each node handles hundreds of thousands of concurrent connections using lightweight BEAM processes.
- A **connection registry** (e.g., a distributed key-value store backed by ZooKeeper or a custom ring) maps `user_id -> gateway_node_id`. When a message arrives for a user, the sending gateway looks up this registry to find the recipient's current node and forwards the message over an internal RPC.
- If the connection registry lookup misses (user is offline), the message is stored in the offline queue and the registry is checked again on the user's next connection.

### Message Delivery State Machine

```
                  [post created on sender device]
                           |
                         QUEUED
                           |  (sender gateway receives message)
                         SENT    <-- server-side ack; sender gets single tick
                           |  (recipient device receives & acks)
                        DELIVERED <-- recipient gateway acks; sender gets double tick
                           |  (recipient opens conversation)
                          READ    <-- sender gets double blue tick
                           |
                       (terminal)

Failure path:
  QUEUED -> FAILED after retry exhaustion (network error, 30-day expiry)
```

### Offline Storage

When the recipient device is offline, the message is stored in a per-user offline queue:

```
HBase table: offline_messages
  row key:    user_id + message_id (reverse timestamp for recency ordering)
  columns:    sender_id, encrypted_payload, timestamp, ttl_expiry

Cassandra alternative:
  partition key: recipient_user_id
  clustering key: message_id ASC
  TTL: 30 days (Cassandra native TTL support)
```

On reconnect, the device fetches all pending messages since its last-seen timestamp, acks each one, and the server deletes (or marks delivered) the offline queue entries.

### Group Messages

WhatsApp performs server-side fan-out for group messages:

```text
sender sends one message to group_id
  -> server expands group membership list (HBase: group_members table)
  -> for each member:
       if online: push to member's gateway
       if offline: write to member's offline_messages queue
  -> delivery acks collected per-member
  -> group delivery status = DELIVERED when all online members have acked
```

Group size is capped at 1,024 members (WhatsApp) to bound the fan-out cost. Larger communities use broadcast channels (one-way, no fan-out ack tracking).

### End-to-End Encryption (Signal Protocol)

WhatsApp and Messenger's end-to-end encryption uses the Signal Protocol:

- Each device generates a long-term identity key pair and a set of one-time prekeys.
- Prekeys are published to the server's key distribution service.
- When Alice wants to message Bob, she fetches one of Bob's prekeys from the server and performs a Diffie-Hellman key exchange to derive a shared secret — without the server ever knowing the secret.
- The server stores only the encrypted ciphertext. It cannot read, index, or scan message content.
- Key implications for system design:
  - Message search is impossible server-side; search is local on-device only.
  - Content moderation relies on user reports (which include a small decrypted sample) rather than server-side scanning.
  - Multi-device support requires per-device key sessions (each device gets its own copy of each message, encrypted separately).

### Presence

Last-seen timestamps are stored in a fast KV store (Redis/Memcached):

```
key:   presence:{user_id}
value: {status: "online"|"offline", last_seen: timestamp}
TTL:   5 minutes (refreshed by heartbeat; expiry = went offline)
```

Presence updates are published to TAO graph edges connecting the user to their contacts. To avoid a "presence storm" when a popular user reconnects (millions of contacts receiving an update simultaneously), updates are throttled to at most one per 60 seconds and fan-out is rate-limited by the presence service.

### Meta-Specific Infrastructure

- **Iris**: Meta's internal messaging bus. Iris is the layer that routes messages between gateway nodes, the message store, and notification systems.
- **MQTT**: WhatsApp uses MQTT (not raw WebSocket) for mobile connections. MQTT has a smaller protocol overhead than WebSocket and handles lossy mobile networks better through its quality-of-service levels.
- **Scribe**: Meta's distributed log collection system. All message events (sent, delivered, read, failed) are written to Scribe for analytics, debugging, and compliance. Scribe does not see message content — only metadata.

## Data Model

```
-- Message store (Cassandra, simplified)
message(
  conversation_id UUID,
  message_id      BIGINT,  -- time-ordered, device-generated
  sender_id       BIGINT,
  encrypted_blob  BYTES,
  state           SMALLINT, -- QUEUED=0, SENT=1, DELIVERED=2, READ=3, FAILED=4
  created_at      TIMESTAMP,
  PRIMARY KEY (conversation_id, message_id)
) WITH CLUSTERING ORDER BY (message_id ASC)
  AND default_time_to_live = 2592000; -- 30 days

-- Offline queue (HBase)
-- row key: recipient_user_id + reverse(message_id)
-- columns: sender_id, encrypted_payload, conversation_id, ttl_expiry

-- Group membership
group_member(group_id BIGINT, user_id BIGINT, joined_at TIMESTAMP)

-- Presence
presence(user_id BIGINT, status VARCHAR, last_seen TIMESTAMP)
```

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| Gateway node crash | Connection drops; device reconnects to a different node | Offline queue ensures no message loss; reconnect triggers offline sync |
| Message store partition | Write failure on message persist | Retry with exponential backoff; message remains QUEUED on sender until ack |
| Delivery ack lost in transit | No DELIVERED ack received | Sender retries delivery; recipient deduplicates using message_id |
| Presence storm on reconnect | Presence fan-out latency spike | Rate-limit presence updates to 1/min; batch contact notifications |
| Group fan-out lag | Member A sees message 5s after Member B | Acceptable under SLA; show per-member delivery status independently |
| MQTT connection flap on mobile | Device disconnects and reconnects rapidly | Exponential backoff on client; gateway rate-limits reconnect storms |

## Observability

- metric: message_delivery_latency_ms p50/p99 (sender send to recipient ack)
- metric: offline_queue_depth per user (alert on p99 > 1000 messages)
- metric: gateway_connections_active per node
- metric: delivery_ack_rate (SENT -> DELIVERED transitions per second)
- metric: failed_delivery_rate (messages reaching FAILED state)
- metric: presence_fan_out_latency_ms
- log: per-message state transition with message_id and user_ids (no content)
- trace: send -> store -> gateway-route -> device-ack latency breakdown

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| WebSocket/MQTT persistent connection | Sub-second delivery; server push | Each connection consumes a file descriptor and memory; gateway must scale to millions of connections | HTTP polling simpler but adds 1–30s latency and wastes battery on mobile |
| Server-side group fan-out | Single API call from sender; server handles expansion | Server fan-out cost scales with group size; group size must be capped | Client-side fan-out (sender addresses each member) shifts complexity and cost to client; breaks E2E encryption key management |
| HBase/Cassandra for offline queue | Scales to billions of user queues; TTL support | Cassandra fan-out latency for large offline queues on reconnect | SQL databases cannot handle per-user queue isolation at this scale |
| Signal Protocol E2E encryption | Server never sees plaintext; strong privacy guarantee | No server-side content search or moderation; multi-device key management complexity | Server-side encryption simpler operationally but cannot claim E2E privacy |
| MQTT over WebSocket for mobile | Lower protocol overhead; QoS levels handle lossy networks natively | Separate client library; less browser support than WebSocket | Raw WebSocket works well on desktop but suffers on unreliable mobile networks |

## Interview It

**Meta framing:** "Design WhatsApp." Strong answers cover the WebSocket connection layer, the offline message queue, the delivery state machine (SENT/DELIVERED/READ), and at least one of: E2E encryption implications, group fan-out, or presence. Weak answers stop at "store messages in a database and poll."

**Follow-ups:**

1. Both Alice and Bob are online. Walk me through the exact path of a message from Alice's device to Bob's screen and back as a delivery ack.
2. Bob goes offline for 3 days and then reconnects. How does his device get all missed messages without duplicates?
3. How does the system prevent a gateway crash from silently dropping in-flight messages?
4. Alice is in a group with 1,024 members. 900 are online when she sends a message. What happens to the 124 offline members?
5. WhatsApp rolls out multi-device support (same account, 4 devices). What changes in the message delivery and encryption model?

## Ship It

- `outputs/design-doc-real-time-messaging.md`
- `outputs/message-delivery-state-machine.md`
- `outputs/interview-card-real-time-messaging.md`

## Exercises

1. **Easy** — Draw the delivery state machine as a diagram. Label every valid transition and every terminal state.
2. **Medium** — Design the deduplication mechanism that prevents Bob's device from displaying a message twice if the delivery ack was lost and the server retransmitted.
3. **Hard** — Extend the design to support disappearing messages (auto-delete after 7 days). What changes in the storage layer, the state machine, and the client?

## Further Reading

- [WhatsApp Engineering Blog — End-to-End Encryption](https://engineering.fb.com/2016/07/06/security/end-to-end-encryption-in-the-real-world/)
- [Signal Protocol technical specification](https://signal.org/docs/)
- [Erlang/OTP and Ejabberd for XMPP at scale](https://www.erlang.org/)
- [MQTT specification (OASIS)](https://mqtt.org/mqtt-specification/)
- [Iris — Meta's internal messaging infrastructure (Meta Engineering)](https://engineering.fb.com/2014/10/09/production-engineering/building-mobile-first-infrastructure-for-messenger/)
