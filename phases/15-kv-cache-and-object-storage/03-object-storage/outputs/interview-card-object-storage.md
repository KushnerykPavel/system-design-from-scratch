---
lesson: 03-object-storage
focus: balanced
---

## Clarify first
- Upload size distribution and need for multipart
- Metadata query requirements
- Durability, integrity, and retention expectations

## Must-size numbers
- Upload throughput and backup-window spikes
- Object size distribution
- Metadata query QPS
- Storage footprint by class

## Core design
- Separate blob plane and metadata plane
- Staged upload then finalize pattern
- Async replication, scanning, repair, and lifecycle jobs

## Failure probes
- Finalize fails after bytes land
- Checksums disagree across replicas or parts
- Delete collides with retention or legal hold

## Trade-off summary
- Simpler hot path vs stronger inline guarantees
- Storage cost vs durability class
- Operational flexibility vs more state machine complexity
