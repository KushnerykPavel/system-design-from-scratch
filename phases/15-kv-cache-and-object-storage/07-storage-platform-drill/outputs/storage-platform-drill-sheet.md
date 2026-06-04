# Storage Platform Drill Sheet

## Clarify
- Upload versus read versus listing priority
- Durability tiers exposed to users
- Freshness targets for metadata and derived assets
- Retention, delete, and legal-hold constraints

## Must-size numbers
- Upload QPS and average object size
- Metadata QPS and filter shapes
- Derived-asset hit-rate goal
- Footprint by durability class

## Minimum architecture
- Object plane
- Metadata plane
- Cache layer for hot derivatives
- Repair and lifecycle control plane

## Failure pushes
- Missing metadata after successful upload
- Cache flush under read spike
- Under-durable objects after component failures
- Retention-aware delete bug
