# InMail & Messaging Platform

> InMail is professional communication — spam it and you lose your sending privilege permanently.

**Type:** Build
**Company focus:** LinkedIn
**Learning goal:** Design LinkedIn's messaging system that handles 1B+ messages/month with spam prevention, read receipts, and recruiter workflow integration.
**Prerequisites:** `16-application-backends/03-chat-system`, `12-security-abuse-and-multitenancy/03-abuse-prevention`
**Estimated time:** ~75 min
**Primary artifact:** InMail delivery pipeline + spam scoring design

## The Problem

Design LinkedIn's messaging platform. It must:

1. **Separate InMail from regular messaging** — InMail is paid cold outreach (recruiters pay per credit), regular messaging is free for connections; they have different delivery paths, rate limits, and spam rules.
2. **Prevent spam at scale** — 1B messages/month means even a 0.1% spam rate is 1M unwanted messages; reputation-based sending limits must be enforced in real time.
3. **Credit system integrity** — an InMail credit is a purchased unit of value; credits must not be double-spent, and must be refunded atomically if the recipient does not respond within 7 days.
4. **Real-time delivery** — members expect to see new messages within seconds; WebSocket or MQTT connections must be maintained for DAU (~30M members online simultaneously at peak).
5. **Recruiter workflow integration** — recruiters manage hundreds of InMail threads; the system must support bulk-send, response rate dashboards, and CRM-like tagging.

## Clarify

- What message types are in scope? Text-only, or rich media (files, images)?
- Is group messaging in scope? (LinkedIn Groups have up to 10K members — fan-out design is different)
- Should read receipts be mandatory? (Member-controlled opt-out allowed)
- What is the acceptable delivery latency? (Target: <500ms for WebSocket delivery to an online recipient)
- Are end-to-end encrypted messages in scope? (Out of scope for this design — LinkedIn is a professional platform where compliance monitoring requires message access)

## Requirements

### Functional Requirements

- Send InMail (cold outreach, credit-based) and direct message (connection-based, free).
- Enforce per-sender daily limits: 10 InMails/day (free), 100 InMails/day (recruiter).
- Refund InMail credit automatically if recipient does not respond within 7 days.
- Spam score each outbound InMail; block if score exceeds threshold.
- Deliver messages in real time to online recipients (<500ms); queue for offline recipients.
- Track delivery states: SENT → DELIVERED → READ.
- Display sender's response rate % on recruiter profile page.

### Non-functional Requirements

- 1B messages/month ≈ 400 messages/sec average; 4,000/sec peak.
- 30M daily active users; maintain WebSocket connections for online members.
- Message storage: Cassandra append-only log, partitioned by thread_id.
- Credit operations: atomic (no double-spend); credit refund idempotent.
- Spam model: score computed in <50ms per outbound InMail.

## Capacity Model

| Dimension | Estimate | Detail |
|-----------|----------|--------|
| Messages/sec average | ~400 | 1B/month ÷ 2.6M seconds |
| Messages/sec peak | ~4,000 | 10× evening burst for recruiter sends |
| Online members (peak) | ~5M WebSocket connections | WebSocket pool behind load balancer |
| Cassandra partitions | per thread_id | each thread = one partition key |
| Unread counts | Redis | O(1) increment/decrement per thread |
| Spam model inference | <50ms | TensorFlow Serving, pre-loaded model |

## Architecture

### InMail vs LinkedIn Messaging

| Dimension | InMail | Direct Messaging |
|-----------|--------|-----------------|
| Eligibility | Any member can receive; sender must have credits | Both parties must be 1st-degree connections |
| Cost | 1 credit per send (recruiter pays) | Free |
| Credit refund | Yes — if no reply within 7 days | N/A |
| Daily limit | 10/day (free), 100/day (recruiter) | No hard limit (spam scoring still applies) |
| Spam risk | Higher — cold outreach | Lower — established connection |
| Content rules | No phone numbers or URLs in first message | More permissive |

### Message Delivery Pipeline

```
Sender → POST /inmail
    ↓
[1. Auth & Eligibility]
   — sender authenticated?
   — is recipient in sender's network or open to InMail?
    ↓
[2. Spam Scoring]  (<50ms, synchronous)
   — SenderProfile: account age, response rate, report rate, daily velocity
   — ML classifier: content + profile features → SpamScore (0..1)
   — If SpamScore > 0.7: REJECT with explanation
    ↓
[3. Rate Limit Check]
   — Redis INCR on key: inmail:{sender_id}:{date}
   — If count > DailyLimit: REJECT "daily limit reached"
    ↓
[4. Credit Deduction]  (atomic, Espresso transaction)
   — Deduct 1 credit from sender's account
   — Write credit_event: {inmail_id, sender_id, credits_deducted=1, refund_at=now+7d}
    ↓
[5. Message Persistence]  (Cassandra)
   — INSERT INTO messages (thread_id, timestamp, sender_id, body, state=SENT)
   — UPDATE unread_counts (Redis INCR for recipient)
    ↓
[6. Real-time Delivery]
   — Check presence service: is recipient online?
   — If online: push via WebSocket → recipient's browser/app
   — If offline: enqueue to Kafka topic: message-delivery
   — Push notification service: send mobile push
    ↓
[7. Delivery Acknowledgment]
   — Client ACKs receipt: state SENT → DELIVERED
   — Client opens thread: state DELIVERED → READ (if read receipts enabled)
    ↓
[8. Credit Refund Job]  (async, Samza consumer)
   — Consumes Kafka topic: inmail-credit-refund-due
   — If no READ/REPLY event within 7 days: refund credit atomically
```

### Spam Prevention Architecture

**Sender Reputation Score** is a weighted function of:
- Response rate (last 90 days): high response rate → lower spam risk
- Report rate (reported as spam by recipients): high report rate → higher risk
- Account age: accounts < 30 days old are throttled by default
- Daily send velocity: sending 50 messages in 1 hour is suspicious

**ML Spam Classifier** is trained on:
- Positive class: InMails reported as spam by recipients
- Negative class: InMails that received responses

Features include: message body length, presence of phone numbers or URLs, keyword patterns, sender profile completeness, ratio of first-InMails to total sends.

**Content Rules** (hard-coded, not ML):
- No phone numbers in first InMail (reduces off-platform redirection)
- No external URLs in first InMail (phishing vector)
- Message body must be >20 characters (avoids "click this link" stub messages)

**Account Quality Score Gating:**
New recruiter accounts are limited to 5 InMails/day for the first 30 days regardless of paid tier, preventing bulk spam from newly created accounts.

### Message Storage Design

```
Cassandra table: messages
  partition key: thread_id  (UUID, derived from sorted pair of member_ids)
  clustering key: timestamp  (DESC — newest first)
  columns: message_id, sender_id, body, state, read_receipt_enabled
```

Thread ID is deterministic: `thread_id = SHA256(min(member_a, member_b) + max(member_a, member_b))`. This ensures two members always map to the same thread_id regardless of who queries first.

**Unread counts:** Redis `INCR thread:{recipient_id}:{thread_id}:unread` on delivery, `DEL` on read. This allows the notification badge count to be served from Redis without reading Cassandra.

### Read Receipts

Read receipts are opt-in/opt-out per member (privacy setting). When enabled:
1. Client sends `ACK{message_id, state=READ}` when message is visible on screen.
2. Server updates Cassandra row state column.
3. Server pushes READ event to sender's WebSocket connection.

When disabled by the recipient, the sender sees only SENT/DELIVERED states.

**False positive risk:** Email clients that auto-open preview panes can trigger the email notification link, which some LinkedIn clients misinterpret as a READ event. Mitigation: require explicit client ACK (web/app), not link-click-based read tracking.

### Group Messages

LinkedIn Groups support up to 10K members. Fan-out strategy:
- **Small groups (<100 members):** synchronous fan-out — one Cassandra write per recipient thread.
- **Large groups (>100 members):** fan-out via Kafka. One write to group message log; consumer fan-out writes to individual recipient threads asynchronously.

### Recruiter InMail Analytics

Recruiters see their InMail response rate on their profile. This is computed by Pinot in near real-time:
- Numerator: InMails that received a reply within 14 days
- Denominator: InMails sent in last 90 days
- Displayed as a percentage on the recruiter's LinkedIn profile

This creates a natural incentive for quality outreach — a low response rate is publicly visible.

## Failure Modes

| Mode | Cause | Mitigation |
|------|-------|------------|
| Spam burst from compromised account | Credentials stolen; attacker sends bulk InMail | Rate limit + Account Quality Score gate; velocity alarm triggers account suspension |
| Message delivery during outage | WebSocket service restarts; offline queue lost | Kafka-backed delivery queue; consumer drains on recovery; messages not ACKd re-delivered |
| Credit double-spend | Race condition between two concurrent InMail sends | Optimistic lock on credit balance in Espresso; conditional write with version check |
| Read receipt false positive | Email client auto-opens notification email | Read state updated only on explicit client ACK, not on notification click |
| Consumer group rebalance on high-volume day | Kafka consumer group rebalance during recruiter mass-send event | Sticky partition assignment; pre-warm consumers before known high-volume periods (e.g., January job market surge) |

## Interview Trade-offs to Discuss

- **Cassandra vs. relational DB for messages:** Cassandra's append-only model and partition-per-thread gives linear write scaling and fast range reads per thread. A relational DB would require complex sharding for 1B messages/month. Trade-off: no cross-thread queries (need Elasticsearch for search).
- **WebSocket vs. polling for real-time delivery:** WebSocket maintains a persistent connection — lower latency, higher server resource per user. Polling is simpler but adds 2–30 second delivery delay. At 30M DAU, WebSocket connection management (sticky load balancing, reconnect handling) is significant operational complexity.
- **Synchronous spam scoring:** Adds 50ms latency to every send. Alternative: async post-delivery scoring with message recall. Synchronous is better for user trust — spammed messages never reach the inbox.
- **InMail credit refund window (7 days):** This is a business policy, not a technical constraint. Shorter windows (3 days) increase refund rate but incentivize low-quality first messages. Longer windows (30 days) make the credit system feel extractive.
