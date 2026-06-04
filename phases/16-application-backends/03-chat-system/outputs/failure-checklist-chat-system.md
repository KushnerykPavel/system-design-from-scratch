# Failure Checklist — Chat System

- What exactly is acknowledged to the sender: durable append, recipient delivery, or read?
- Can reconnecting clients replay from stable cursors after gateway loss?
- Is stale presence allowed to misroute messages without losing them?
- Can one noisy conversation be throttled without delaying the whole fleet?
- Are duplicate retries deduped by client message ID?
- Are attachment fetch failures isolated from text-message delivery?
