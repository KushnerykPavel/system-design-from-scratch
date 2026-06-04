# Location Pipeline Operator Checklist

## Freshness

- What is the current publish-to-serve latency by region?
- What percentage of accepted updates exceed the freshness budget?
- Are stale drops rising because of real lag or because client clocks are bad?

## Correctness

- Is dedupe-hit rate elevated because of a retry storm?
- Are out-of-order discards concentrated in one device cohort or region?
- Are impossible-speed or GPS-jump filters overfiring?

## Capacity

- Which partitions have the highest ingress skew?
- Which tenants or device classes are generating the most retries?
- Is live-state publication lagging even when ingest looks healthy?

## Degradation

- Which products should fail data stale instead of serving old positions?
- What operator controls exist for raising or lowering stale-age thresholds?
- How is regional isolation preserved during one-region backlog incidents?
