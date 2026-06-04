# Interview Card: Eviction Policies

## Strong answer shape

- describe the working set and traffic skew
- name the likely harmful pattern: scans, trend shifts, or memory pressure
- choose a policy for that trace, not by habit
- mention what metric proves the policy is working

## Senior signals

- "LRU is good if recency drives reuse, but scans can pollute it badly."
- "LFU helps on stable popularity, but old hot keys can linger unless we age counts."
- "If every policy is weak, the problem may be memory budget or admission, not eviction alone."
