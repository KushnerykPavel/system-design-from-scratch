---
lesson: 05-durability-tiers
focus: balanced
---

## Clarify first
- Data criticality classes
- Failure domains that count
- Repair-time guarantees and regional disaster expectations

## Must-size numbers
- Volume by tier
- Expected degraded-component rate
- Repair bandwidth budget
- Regional restore objectives

## Core design
- Multiple documented durability classes
- Placement and coding policy per class
- Auditor and repair scheduler coupled to those promises

## Failure probes
- Slow repair extends vulnerability window
- Geo copies are not actually independent
- Repair work competes with serving traffic

## Trade-off summary
- Cost vs safety
- Fast reads vs coding efficiency
- Fewer tiers vs product-fit flexibility
