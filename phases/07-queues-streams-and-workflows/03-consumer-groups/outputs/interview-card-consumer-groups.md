---
lesson: 07-consumer-groups
focus: balanced
---

# Interview Card: Consumer Groups

## Clarify First

- What must stay ordered together?
- How skewed is the workload?
- What is the acceptable lag during incidents and deploys?
- How often do members join and leave?

## Must Mention

- Partition key
- Partition count as a parallelism ceiling
- Offset ownership and rebalance behavior
- Hot partition mitigation

## Failure Probes

- Hot tenant or celebrity account
- Rolling deploy rebalance stall
- Too few partitions after growth
- Slow downstream dependency causing uneven lag
