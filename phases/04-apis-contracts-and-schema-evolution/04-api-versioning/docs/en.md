# API Versioning and Compatibility

> Good versioning is less about adding `/v2` and more about preventing accidental client breakage.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Explain how API changes stay compatible over time, when to version explicitly, and how to roll out contract changes without stranding clients.  
**Prerequisites:** `04-apis-contracts-and-schema-evolution/01-http-vs-grpc-vs-events`, `04-apis-contracts-and-schema-evolution/03-pagination-and-filtering`  
**Estimated time:** ~75 min  
**Primary artifact:** compatibility matrix + rollout checklist  

## The Problem

Many interview answers treat versioning as a naming detail: "we'll just make `/v2`." Senior answers go deeper:

- what kinds of changes are backward compatible
- how long old clients live
- what signals show that migrations are incomplete
- how rollout and deprecation interact with reliability

## Clarify

- Is the API public, partner-facing, or fully internal?
- What is the real client upgrade curve: hours, weeks, or years?
- Which changes are expected: added fields, renamed fields, semantics changes, or removed endpoints?
- Can the server translate between versions, or must consumers migrate directly?

## Requirements

### Functional

- Introduce changes without breaking supported clients.
- Communicate deprecation and migration paths clearly.
- Detect old-version traffic and incomplete rollouts.

### Non-functional

- Minimize parallel-version operational burden.
- Keep compatibility rules simple enough to enforce consistently.
- Avoid hidden semantic drift across versions.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Active clients | 500K SDK users + internal services | determines long-tail upgrade pain |
| Version retention window | 6-18 months | drives how many versions stay live |
| Peak QPS on old versions | 20% during migration | old traffic must still be observable |
| Schema churn | monthly additive changes, rare breaking changes | suggests compatibility-first design |
| Rough cost | multiple docs, translation logic, test matrix | version sprawl becomes a real tax |

## Architecture

Strong default:

1. Prefer additive, backward-compatible changes.
2. Reserve explicit version breaks for semantic incompatibility.
3. Instrument traffic by contract version.
4. Run deprecation like a rollout, not a documentation note.

Possible versioning levers:

- URI versioning for clear public breaks
- header or media type versioning for more flexible negotiation
- protobuf field evolution rules for gRPC

## Data Model & APIs

Compatibility examples:

- adding optional response fields is usually safe
- removing fields, changing meaning, or changing sort semantics is usually breaking
- reusing a field name for new semantics is worse than adding a new field

A good design calls out:

- supported version list
- deprecation policy
- migration timeline
- telemetry by version

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| additive-looking change is semantically breaking | client error spike without schema mismatch | document semantics, not only shape |
| version sprawl | growing traffic on many live versions | enforce support window and retirement policy |
| untracked old clients | deprecation date passes but hidden clients remain | version-tagged request metrics and outreach |
| server translation layer drifts | version-specific bugs rise | test matrix and contract golden cases |

## Observability

- metric: request volume by API version
- metric: error rate and latency by version
- metric: deprecation header seen vs acted upon
- log: version, client type, and sunset phase
- trace: translation layer path when serving legacy versions
- SLO: supported legacy versions remain stable until declared sunset dates

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| compatibility-first additive changes | fewer forced migrations | contract surface grows over time | frequent hard breaks |
| explicit version for semantic breaks | cleaner expectations | extra operational support burden | silent semantic changes in-place |
| translation layer for legacy clients | smoother migration | more code and test complexity | forcing immediate client rewrites |

## Interview It

**Google framing:** "Design an API evolution strategy for a public service used by mobile clients." The signal is whether you understand client lag and backward compatibility pressure.

**Cloudflare framing:** "Roll out a contract change across globally distributed infrastructure clients." The signal is whether you treat version telemetry and deprecation as operational work.

**Follow-ups:**
1. What if mobile clients upgrade slowly for months?
2. What if the new API changes semantics, not just field names?
3. How many versions would you support at once?
4. What if partner clients ignore deprecation notices?
5. How do you test a translation layer without exploding the matrix?

## Ship It

- `outputs/compatibility-matrix-api-versioning.md`
- `outputs/rollout-checklist-api-versioning.md`

## Exercises

1. **Easy** — Classify five example changes as compatible or breaking.  
2. **Medium** — Write a deprecation plan for a field that should stop being used in six months.  
3. **Hard** — Redesign a partner API where `/v1` and `/v2` have already drifted semantically.  

## Further Reading

- [Google API Design Guide](https://cloud.google.com/apis/design/versioning) — strong baseline on versioning decisions  
- [Protocol Buffers update rules](https://protobuf.dev/programming-guides/proto3/) — practical compatibility constraints for typed RPC contracts  
