# Blob Metadata Separation

> Large objects want cheap durable bytes; user-facing queries want small indexed metadata.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Design systems where large blobs live separately from the metadata and indexes that drive discovery, policy, and request routing.  
**Prerequisites:** `01-storage-models`, `04-hot-and-cold-data`  
**Estimated time:** ~60 min  
**Primary artifact:** design review prompt + failure checklist  

## The Problem

Beginners often place binary objects and rich queryable metadata in the same storage path. That usually creates the worst of both worlds: expensive scans over large objects, awkward consistency, and poor lifecycle control.

This lesson trains the standard pattern:

- store **blob bytes** in an object or blob storage system
- store **metadata and indexes** in a smaller queryable store
- keep references, integrity checks, and lifecycle policy explicit

## Clarify

- Are objects immutable after upload, or can they be updated in place?
- What metadata needs low-latency filtering or listing?
- What is the allowed delay between object upload and metadata visibility?
- Are deletion, retention, or legal hold requirements stricter for metadata, blobs, or both?

## Requirements

### Functional

- Upload and retrieve large objects reliably.
- Query metadata by owner, type, state, or time window.
- Support deletion and lifecycle changes without losing track of stored bytes.

### Non-functional

- Keep metadata queries fast without reading full blobs.
- Make object durability and integrity visible.
- Avoid orphaned objects and broken metadata pointers.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Upload rate | 12K objects/s peak | drives write-path decoupling and async processing |
| Average blob size | 8 MB | blobs should not live in OLTP indexes |
| Metadata query rate | 40K req/s | metadata store must be optimized separately |
| Peak factor | 5x during ingest spikes | tests queueing and object finalize flow |
| Rough cost | object storage + metadata DB + lifecycle jobs | separate cost planes justify separation |

## Architecture

Canonical flow:

```text
client
  -> upload session service
  -> blob store
  -> finalize metadata record
  -> async scanners / thumbnails / policy jobs
```

The important rule is that metadata should record:

- object identifier
- owner or tenant
- size, checksum, content type
- lifecycle and visibility state

## Data Model & APIs

Example metadata record:

```text
object_id
tenant_id
blob_uri
size_bytes
checksum
status
created_at
retention_class
```

Useful APIs:

- `CreateUploadSession`
- `FinalizeObjectUpload`
- `ListObjects(tenant_id, prefix, cursor)`
- `DeleteObject(object_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| blob uploaded but metadata finalize fails | orphan-object sweeps find unreferenced blobs | use staged upload state and reconciliation jobs |
| metadata exists but blob is missing or corrupt | checksum mismatch and object HEAD failures | verify on finalize and run periodic integrity audits |
| listing path accidentally reads blob store directly | latency and request cost spike | keep all list and filter paths on metadata service |
| delete removes metadata but leaves durable bytes | storage drift and compliance risk | two-phase delete with tombstones and audit logs |

## Observability

- metric: finalize success rate after upload completion
- metric: orphan object count and reconciliation age
- metric: metadata query latency separate from blob retrieval latency
- metric: checksum or integrity failure count
- log: object state transitions and delete reasons
- SLO: metadata listing remains fast while upload finalize and delete workflows stay reliable

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| separate blob and metadata stores | fast queries plus cheap durable object storage | extra consistency workflow | storing binaries inline with queryable rows |
| staged finalize flow | safer lifecycle tracking | more states to reason about | treating upload as one opaque write |
| async post-processing | keeps upload path narrow | eventual visibility for derived artifacts | doing scanning and transformation inline |

## Interview It

**Google framing:** "Design document storage for a collaborative workspace." The signal is whether you keep object bytes and rich metadata concerns separate.

**Cloudflare framing:** "Design a large-asset storage path with metadata searchable from the control plane." The signal is whether you distinguish durable blob storage from operationally queryable metadata.

**Follow-ups:**
1. What if uploads can be multipart and resumed?
2. What if metadata must become visible before async virus scanning completes?
3. How do you handle legal hold on deleted objects?
4. What if object retrieval is much rarer than metadata search?
5. How do you detect and clean up orphan blobs safely?

## Ship It

- `outputs/design-review-blob-metadata-separation.md`
- `outputs/failure-checklist-blob-metadata-separation.md`

## Exercises

1. **Easy** — Sketch metadata fields for a video-upload product.  
2. **Medium** — Design the delete workflow for blobs under legal hold.  
3. **Hard** — Explain the upload, finalize, and reconciliation flow for multipart objects across regions.  

## Further Reading

- [Amazon S3 data consistency model](https://docs.aws.amazon.com/AmazonS3/latest/userguide/Welcome.html) — useful reference for blob semantics and lifecycle discussions  
- [Google Cloud Storage object metadata](https://cloud.google.com/storage/docs/metadata) — concrete examples of the blob and metadata split  
