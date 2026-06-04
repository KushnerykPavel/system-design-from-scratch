# Failure Checklist — Log Aggregation Pipeline

- Collectors can buffer briefly when indexers or archives are slow.
- Parser versions are tracked so bad deploys can be rolled back and replayed.
- PII redaction happens before indexing into broad-search systems.
- Hot index lag and archive backlog are measured independently.
- Critical log classes have reserved capacity or priority handling.
- Tenant quotas and routing isolation protect shared infrastructure.
- Replay from durable storage is possible after parser or enrichment bugs.
