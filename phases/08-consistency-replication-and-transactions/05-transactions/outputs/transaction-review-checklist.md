# Transaction Review Checklist

- What invariant must never be broken?
- Which anomaly are we actually preventing?
- What is the smallest boundary that protects the invariant?
- Which entities or keys are likely hotspots?
- Can optimistic concurrency work, or is serialization needed?
- Where should the transaction stop and async coordination begin?
- Which metrics prove contention is becoming the real limit?
