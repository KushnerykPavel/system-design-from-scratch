# Design Review Prompts — Infra Platform Mock

## Ask These During Practice

- Who is the real user of this platform?
- What is the control plane and what is the data plane?
- Which state is authoritative, and how does it propagate?
- What is the smallest failure domain you can promise?
- What is the rollback path for a bad config or rollout?
- Which metric tells you tenants are getting different behavior than intended?

## Use Them To Stress-Test

- stale control-plane state
- noisy-neighbor traffic
- unsafe rollout waves
- tenant-specific SLA upgrades
- region-specific partial failure
