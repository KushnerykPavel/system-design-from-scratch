# Trade-off Matrix — Rough Cost Modeling

| Choice | Lower cost path | Higher cost path | Reason to pay more |
|-------|------------------|------------------|--------------------|
| regions | single region | multi-region | better resilience and latency |
| storage | cheaper cold tiers | always-hot premium storage | faster queries |
| cache | smaller cache | larger cache | lower origin cost and latency |
| platform | self-managed | managed service | less operational complexity |
