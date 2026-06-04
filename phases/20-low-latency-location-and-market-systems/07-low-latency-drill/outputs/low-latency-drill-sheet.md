# Low-Latency Drill Sheet

## Opening

- What is the primary user-visible latency or freshness promise?
- Which request or event path determines whether that promise is met?
- What is the first scope cut if the prompt is too broad?

## Sizing

- What is steady-state QPS or event rate?
- What is the worst hotspot, burst, or fanout multiplier?
- Which path becomes saturated first?

## Architecture

- What must happen synchronously before success is returned?
- What can be asynchronous?
- Which one subsystem deserves the deep dive?

## Reliability

- What is the degraded mode?
- Which metric proves the hot path is healthy?
- Which metric proves one hotspot is melting the design?
