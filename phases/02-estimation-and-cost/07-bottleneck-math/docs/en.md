# Bottleneck Math for CPU, Disk, and Network

> The system is limited by the resource that saturates first, not by the component you spent the most time describing.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Identify the first limiting resource using simple throughput math across CPU, disk, and network.  
**Prerequisites:** `02-estimation-and-cost/03-bandwidth-and-egress`, `05-storage-indexing-and-access-patterns/03-indexes`  
**Estimated time:** ~75 min  
**Primary artifact:** bottleneck checklist + worksheet  

## The Problem

System design answers often say "the database becomes the bottleneck" without explaining whether that means CPU, IOPS, storage bandwidth, or network saturation. Strong answers are more specific.

This lesson teaches quick math to find the first constrained resource.

## Clarify

- Is the workload CPU-heavy, disk-heavy, or network-heavy?
- Are requests sequentially blocked on storage, or mostly memory-resident?
- Are we reasoning about one node, one shard, or the whole fleet?
- Which saturation matters first: throughput, latency, or cost?

## Requirements

### Functional

- Estimate max request rate per resource dimension.
- Identify the first limiting resource.
- Show how changing payload or query cost shifts the bottleneck.

### Non-functional

- Keep the model short and directional.
- Avoid vague "the DB is slow" language.
- Make scaling implications explicit.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| CPU budget | 32 cores x 4K req/s/core | one possible ceiling |
| Disk budget | 80K IOPS / 4 ops per req | storage ceiling |
| Network budget | 12.5 GB/s / 100 KB resp | network ceiling |
| Effective max req/s | min of all three | true bottleneck |
| Headroom target | 30% | avoids cliff-edge operation |

## Architecture

For each resource:

1. Estimate resource cost per request.
2. Divide total capacity by per-request cost.
3. Compare resulting max request rates.
4. The smallest limit is the first bottleneck.

If one node supports:

- 128K req/s by CPU
- 20K req/s by disk
- 131K req/s by network

Then disk is the real bottleneck, regardless of how powerful the CPUs look.

## Data Model & APIs

The code artifact models:

```text
BottleneckModel {
  CPUReqPerSecond
  DiskReqPerSecond
  NetworkReqPerSecond
}
```

Output:

- bottleneck resource
- max sustainable requests per second

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| one resource assumed without math | wrong scaling plan | compute all major ceilings quickly |
| per-request disk ops underestimated | IOPS saturates early | count indexes, fanout reads, and writes honestly |
| network payload ignored | links saturate at high response size | size bytes/request, not only QPS |
| no headroom reserved | latency spikes near saturation | plan operating margin, not hard max |

## Observability

- metric: CPU saturation and run queue
- metric: disk IOPS, throughput, and queue depth
- metric: network throughput and packet drops
- metric: requests per second vs headroom on each tier
- SLO: no critical tier should operate near saturation during normal peak

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| model three resources explicitly | clearer scaling plan | more arithmetic | vague component-level bottleneck claims |
| keep headroom | better resilience | lower utilization | running at theoretical max |
| optimize the first constraint | highest immediate gain | may expose next bottleneck | scattered micro-optimizations |

## Interview It

**Google framing:** "Design a heavy read API backed by a relational store." The signal is whether you know whether CPU, disk, or network fails first.

**Cloudflare framing:** "Design a globally served object metadata lookup path." The signal is whether you identify which layer actually saturates under byte and request pressure.

**Follow-ups:**
1. What if a cache removes most disk reads?
2. What if response size doubles?
3. What if the write path adds two extra index writes per request?
4. What if the fastest improvement is sharding, not vertical scaling?

## Ship It

- `outputs/bottleneck-worksheet.md`
- `outputs/bottleneck-checklist.md`

## Exercises

1. **Easy** — Identify the bottleneck for a CPU-heavy stateless service.  
2. **Medium** — Recompute when response size doubles.  
3. **Hard** — Explain how fixing disk IOPS may expose network as the next limit.  

## Further Reading

- [Google SRE book](https://sre.google/books/) — useful background on capacity headroom  
- [System design notes](https://github.com/liquidslr/system-design-notes) — baseline sizing context  
