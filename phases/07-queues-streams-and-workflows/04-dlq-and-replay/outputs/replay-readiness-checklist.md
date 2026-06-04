---
lesson: 07-dlq-and-replay
focus: balanced
---

# Replay Readiness Checklist

## Before You Dead-Letter

- Transient retry policy is distinct from quarantine policy
- Failure class is recorded, not just raw error text
- Original topic, partition, offset, and key are preserved

## Before You Replay

- Root cause is understood or a mitigation is deployed
- Replay rate limit is chosen
- Live traffic protection is in place
- Ordering requirements are checked
- Duplicate side effects are considered

## During Replay

- Success and repeat-failure rates are monitored
- Downstream error budget is watched
- Operator can stop replay quickly

## After Replay

- Residual dead letters are reviewed
- Root cause and detection gaps are documented
- Permanent fix avoids repeating the same class of failures
