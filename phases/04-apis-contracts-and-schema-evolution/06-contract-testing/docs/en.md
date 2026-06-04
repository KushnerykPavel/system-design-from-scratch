# Ownership Boundaries and Contract Testing

> Shared contracts fail when nobody is clearly responsible for keeping them honest.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Define ownership boundaries around service contracts and explain how contract tests catch breaking changes before deployment.  
**Prerequisites:** `04-apis-contracts-and-schema-evolution/04-api-versioning`, `04-apis-contracts-and-schema-evolution/05-event-schema-evolution`  
**Estimated time:** ~75 min  
**Primary artifact:** contract-test checklist + ownership rubric  

## The Problem

Compatibility rules are not enough if teams cannot tell who owns the promise or how they verify it. In multi-team systems, the contract itself becomes a product surface with producers, consumers, rollout order, and failure blast radius.

This lesson connects two ideas:

- boundaries need explicit ownership
- compatibility should be verified continuously, not discovered in production

## Clarify

- Is this a public API, an internal platform contract, or an event consumed by many teams?
- Which team owns the source of truth for schema and behavior?
- Are consumer expectations documented only in code, or in shared contracts?
- What is the deployment order when provider and consumer both change?

## Requirements

### Functional

- Define who owns the contract and change review.
- Validate that providers still satisfy consumer expectations.
- Catch breaking changes before they reach production.

### Non-functional

- Keep test feedback fast enough for normal development.
- Avoid brittle end-to-end-only validation.
- Limit cross-team coordination cost while preserving trust.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Service pairs | 30 providers, 120 consumers | test fanout grows fast with shared platforms |
| Contract changes | dozens per month | manual review alone will miss regressions |
| CI budget | minutes, not hours | tests must be scoped and automatable |
| Blast radius | one provider can break many consumers | strong ownership and compatibility gating matter |
| Rough cost | schema tooling + CI matrix + review process | cheaper than production contract outages |

## Architecture

Healthy pattern:

1. Provider owns the contract publication.
2. Consumers publish expectations or generated contract fixtures.
3. CI verifies provider compatibility against consumer expectations.
4. Rollout includes observability for real traffic behavior.

```text
provider change -> contract validation -> consumer expectation tests -> deploy
```

Contract tests are not a substitute for integration tests. They are a cheaper, sharper guardrail against accidental compatibility drift.

## Data Model & APIs

Boundary questions:

- what fields and semantics are guaranteed
- who approves breaking changes
- what compatibility window exists
- how request and response examples are versioned

Provider tests should answer:

- can I still serve existing consumers?
- did I change status codes, ordering, defaults, or required fields?

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| provider changes behavior without contract update | contract tests fail or consumer errors spike | block deploy on compatibility violations |
| nobody owns consumer expectation updates | stale fixtures and weak trust | assign clear provider and consumer responsibilities |
| over-reliance on end-to-end tests | regressions found late and noisily | add focused contract tests in CI |
| consumer expectations too loose | tests pass while production breaks semantically | include behavior-level assertions, not just shape |

## Observability

- metric: contract test failures by provider and change type
- metric: version skew between provider and consumers
- metric: production errors after contract-related deploys
- log: deploys tagged with contract change identifiers
- trace: request path including version and feature-flag context
- SLO: supported consumers continue to pass contract checks across routine provider releases

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| provider-owned contract with consumer assertions | clear accountability and practical safety | coordination overhead | vague shared ownership |
| CI contract gating | faster breakage detection | more setup and fixture maintenance | waiting for production canaries only |
| focused contract tests | cheap and specific | not full-system validation | end-to-end tests as the only defense |

## Interview It

**Google framing:** "How would you stop a platform team from breaking dozens of clients with a schema change?" The signal is whether you talk about ownership and automated checks together.

**Cloudflare framing:** "How do you evolve shared edge service contracts across many internal consumers?" The signal is whether you manage blast radius and rollout order, not just schema files.

**Follow-ups:**
1. What if consumers cannot all publish formal expectations?
2. What if a provider change is intentionally breaking?
3. How do you test semantics like default behavior or ordering?
4. Where do end-to-end tests still matter?
5. What changes when the contract is public to customers instead of internal teams?

## Ship It

- `outputs/contract-test-checklist.md`
- `outputs/ownership-rubric-contract-testing.md`

## Exercises

1. **Easy** — Name the provider and consumer responsibilities for a shared profile service.  
2. **Medium** — Design a contract test for an API whose status code behavior changed accidentally.  
3. **Hard** — Redesign a team boundary where "everyone owns the API" and breakages are recurring.  

## Further Reading

- [Pact contract testing](https://docs.pact.io/) — practical introduction to consumer-driven contract testing  
- [Google API Improvement Proposals](https://google.aip.dev/) — good examples of contract ownership and evolution discipline  
