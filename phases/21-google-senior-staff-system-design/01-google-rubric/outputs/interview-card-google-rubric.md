# Interview Card - Google Rubric

## Before You Draw

- clarify primary user journey
- prioritize one latency, throughput, or correctness goal
- estimate request volume, storage, and growth
- choose one likely deep dive

## Strong Signals

- assumptions are explicit and reasonable
- architecture is tied to the workload
- trade-offs are named without prompting
- failure modes and observability appear before the end
- redesign changes something real

## Red Flags

- giant diagram before scope
- no quantitative estimates
- every component is "for scalability"
- no detection story for failures
- magical consistency or zero-cost multi-region claims
