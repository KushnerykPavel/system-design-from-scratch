---
lesson: 07-messaging-drill
focus: balanced
---

# Scoring Rubric: Messaging Drill

| Dimension | Strong | Weak |
|-----------|--------|------|
| Primitive choice | matched to ownership and replay needs | picked by habit or product name |
| Delivery semantics | honest boundary and duplicate policy | vague exactly-once claims |
| Partitioning | ordering scope and skew discussed | no key or parallelism reasoning |
| Recovery | DLQ and replay are controlled | failures dumped with no plan |
| Backpressure | lag visibility and overload policy clear | queue treated as infinite safety |
| Trade-offs | explicit cost and complexity discussion | component list without judgment |
