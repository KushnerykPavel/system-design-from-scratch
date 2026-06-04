---
lesson: 05-bot-mitigation
---

| Score band | Default action | Safer rollout action | Notes |
|------------|----------------|----------------------|-------|
| low risk | allow | allow | keep latency minimal |
| medium risk | challenge or throttle | log-only or challenge | best for staged rollout |
| high risk | block | challenge then block | use strongest reason-code visibility |
