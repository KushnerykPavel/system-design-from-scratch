# Debugging Narrative Worksheet

## Symptom

- What user harm is visible?
- What changed recently?

## Plausible cause buckets

- deploy / config
- dependency
- capacity / saturation
- skew / abuse
- data / correctness

## Highest-signal next checks

- Which metric narrows scope fastest?
- Which comparison separates local from global?
- Which trace or log query tests the leading hypothesis?

## Mitigation

- Safest reversible action:
- Risk if wrong:

## Validation

- Which user-facing signal should recover?
- Which bottleneck signal should move if the hypothesis was correct?
