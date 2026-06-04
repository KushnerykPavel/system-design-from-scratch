# Google Senior / Staff System Design Rubric

Score each dimension from 1-4. Total: 32.

| Dimension | 1 (Poor) | 2 (Weak) | 3 (Senior) | 4 (Staff) |
|-----------|----------|----------|------------|-----------|
| **Clarification quality** | Asks almost nothing | Asks shallow questions | Surfaces real design constraints early | Reframes prompt to expose hidden priorities |
| **Requirement prioritization** | Lists requirements randomly | Mentions some trade-offs | Chooses and defends priorities clearly | Aligns trade-offs to business and reliability goals |
| **Sizing and capacity** | Skips estimation | Rough numbers but little impact | Sizing is directionally sound and used in design | Uses sizing to drive bottleneck, cost, and topology decisions |
| **Architecture decomposition** | Component soup | Mostly coherent | Clean boundaries and data flow | Strong modularity plus clear control/data-plane separation |
| **Deep dive quality** | Superficial | Some detail, little rigor | Good reasoning in the critical area | Anticipates interviewer pushback and second-order effects |
| **Failure-mode reasoning** | Ignores failure | Lists common failures | Detection and mitigation are explicit | Degradation, rollout, and blast-radius thinking are strong |
| **Trade-off quality** | Buzzwords only | Hand-wavy | Concrete trade-offs tied to requirements | Evaluates multiple valid paths and chooses intentionally |
| **Communication / leadership** | Hard to follow | Understandable but reactive | Clear, structured, collaborative | Senior/staff presence: aligns, summarizes, and leads the discussion |

## Strong-positive signals

- Clarifies scope before drawing boxes.
- Sizes first and uses the numbers later.
- Says what is out of scope for v1.
- Names one deep dive and explains why.
- Mentions rollout or migration without prompting.
- Distinguishes availability, durability, and consistency clearly.

## Strong-negative signals

- Opens with a database choice before clarifying workload.
- Says "eventual consistency" without defining impact.
- Never mentions observability or incident handling.
- Treats cost and complexity as irrelevant.
- Keeps adding components instead of defending priorities.
