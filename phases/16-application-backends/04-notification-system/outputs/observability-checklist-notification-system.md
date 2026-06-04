# Observability Checklist — Notification System

- Can we see accepted, suppressed, retried, and permanently failed notifications separately?
- Do we measure queue lag by priority class and by channel?
- Are duplicate sends detectable by dedupe key and provider response ambiguity?
- Can we prove user preferences and quiet hours were honored?
- Are provider timeout spikes visible before a hard outage is declared?
- Do traces show policy evaluation separately from transport latency?
