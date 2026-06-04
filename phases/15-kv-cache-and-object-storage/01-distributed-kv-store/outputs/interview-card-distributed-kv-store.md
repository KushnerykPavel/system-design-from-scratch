---
lesson: 01-distributed-kv-store
focus: balanced
---

## Clarify first
- Value size, access skew, and whether clients need freshness or just durability
- Failure domains for replicas: process, host, rack, or AZ
- Whether prefix scans or secondary lookups are required

## Must-size numbers
- Peak read and write QPS
- Hot-key share of total traffic
- Logical data size and replication factor
- Repair and rebalance throughput budget

## Core design
- Hash partitioning with replica sets across independent failure domains
- Write-ahead log plus LSM-style storage for write-heavy efficiency
- Client-selectable consistency or clearly documented service tiers

## Failure probes
- What happens during a network partition?
- How are hinted writes and anti-entropy bounded?
- How do hot keys avoid melting a single owner set?

## Trade-off summary
- Availability vs freshness
- Write cost vs repair complexity
- Simpler routing vs smarter hotspot mitigation
