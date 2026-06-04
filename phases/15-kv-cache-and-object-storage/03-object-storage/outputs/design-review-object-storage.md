# Design Review Prompt — Object Storage

Use this when reviewing an object storage answer or draft design.

## Clarify
- Are objects immutable, versioned, or mutable?
- What metadata queries must be fast without touching blob bytes?
- Which retention and legal-hold constraints change deletion behavior?

## Core checks
- Blob durability and metadata visibility are separate planes.
- Multipart uploads have explicit staged states and abandonment cleanup.
- Integrity verification is not hand-waved; checksums exist on upload and retrieval.
- Lifecycle policy execution is observable and retention-aware.

## Failure probes
- What happens if parts upload successfully but finalize fails?
- Can metadata point to missing or corrupt bytes?
- How do delete and retention policies avoid orphaning or premature purge?
- How is repair throttled during ingest peaks?
