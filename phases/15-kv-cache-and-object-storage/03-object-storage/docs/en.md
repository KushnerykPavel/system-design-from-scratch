# Object Storage and Blob Metadata

> Durable bytes, searchable metadata, and lifecycle control belong in related systems, not in one overloaded abstraction.

**Type:** Build  
**Company focus:** Balanced  
**Learning goal:** Design object storage with explicit separation between blob durability, metadata visibility, integrity verification, and lifecycle workflows.  
**Prerequisites:** `05-storage-indexing-and-access-patterns/05-blob-metadata-separation`, `13-multi-region-cdn-and-edge-traffic/03-cdn-layering`, `15-kv-cache-and-object-storage/01-distributed-kv-store`  
**Estimated time:** ~90 min  
**Primary artifact:** upload-policy validator + design review prompt  

## The Problem

Design an object storage platform for large immutable objects such as images, backups, logs, or user uploads. Clients need scalable writes, cheap durable bytes, strong integrity checks, and low-latency metadata queries for listing or policy decisions.

The interview signal is not "use S3." It is whether you can explain upload staging, metadata finalize, multipart behavior, durability claims, and lifecycle operations without pretending object storage is a magical black box.

## Clarify

- Are objects mutable, versioned, or write-once?
- Is metadata listing or prefix search a first-class product requirement?
- What integrity guarantees are required on upload and retrieval?
- Are geo-redundancy, legal hold, or lifecycle transitions part of the core product?

If the interviewer keeps it broad, assume immutable objects, multipart uploads for large files, metadata listing by tenant and prefix, and async lifecycle management.

## Requirements

### Functional

- Upload, retrieve, and delete large objects.
- Support multipart or resumable uploads.
- Expose searchable metadata separately from object bytes.
- Apply lifecycle policies such as storage-class transitions and retention controls.

### Non-functional

- Keep object durability high without forcing hot metadata queries through blob paths.
- Provide end-to-end integrity verification.
- Support background repair, replication, and policy jobs without breaking serving latency.

## Capacity Model

| Dimension | Estimate | Why it matters |
|-----------|----------|----------------|
| Upload throughput | 25 GB/s regional peak | shapes ingest fan-in and multipart orchestration |
| Average object size | 12 MB, with a long tail into multi-GB | forces separation between metadata and byte durability |
| Metadata query QPS | 60K req/s | metadata plane must scale independently |
| Storage footprint | 30 PB logical, multiple durability classes | cost and lifecycle policy dominate economics |
| Peak factor | 6x during backup windows | ingest spikes change temporary capacity and repair scheduling |

## Architecture

```text
client
  -> upload session / auth service
  -> multipart ingest path
  -> blob placement + checksum verification
  -> metadata finalize
  -> async replication, scanning, lifecycle, and repair
```

Separate planes:

1. **Data plane** stores object chunks or full blobs.
2. **Metadata plane** serves search, listing, policy state, and authorization context.
3. **Control jobs** handle repair, replication, lifecycle transitions, and orphan cleanup.

## Data Model & APIs

Object metadata:

```text
object_id
tenant_id
bucket
blob_uri
content_length
checksum
storage_class
state
version_id
retention_policy
```

Useful APIs:

- `CreateMultipartUpload`
- `UploadPart`
- `CompleteMultipartUpload`
- `HeadObject`
- `ListObjects(bucket, prefix, cursor)`
- `DeleteObject(object_id)`

## Failure Modes

| Failure | Detection | Mitigation |
|---------|-----------|------------|
| parts uploaded but finalize fails | orphan upload metrics and stale sessions | staged states plus sweeper jobs |
| metadata says complete but checksum chain is wrong | checksum validation failures on complete or read | reject finalize and quarantine corrupted parts |
| lifecycle job deletes bytes before retention check | policy violation audit and missing object reads | retention-aware state machine and delayed physical delete |
| listing path accidentally depends on blob store availability | metadata API latency coupled to object path incidents | keep listings entirely in metadata systems |

## Observability

- metric: multipart completion success rate and abandon rate
- metric: object integrity failures by storage class and region
- metric: metadata finalize latency separate from blob write latency
- metric: lifecycle backlog by action type and age
- log: object state transitions, retention exceptions, and delete actors
- trace: upload session to finalize path with part counts and checksum verification

## Trade-offs

| Decision | Benefit | Cost | Alternative rejected |
|----------|---------|------|----------------------|
| separate metadata plane | fast search and control workflows | extra consistency and finalize logic | serving listings from the blob layer |
| immutable object versions | simpler durability and caching semantics | extra storage until cleanup | in-place mutation of large blobs |
| async lifecycle transitions | cheap hot path and flexible policy | eventual policy application | fully synchronous class transition at write time |

## Interview It

**Google framing:** "Design durable object storage for internal backup or user content." Expect questions about metadata search, integrity, and lifecycle automation.

**Cloudflare framing:** "Design an asset storage platform feeding edge delivery." Expect questions about multipart ingest, object versioning, and how metadata and delivery layers stay decoupled.

**Follow-ups:**
1. How do you support resumable multi-GB uploads from unstable clients?
2. What changes if metadata visibility must happen before antivirus scanning completes?
3. How do you enforce legal hold without freezing the whole bucket?
4. What if one storage class is erasure-coded and slower to repair?
5. How do you backfill checksums for legacy objects?

## Ship It

- `outputs/design-review-object-storage.md`
- `outputs/failure-checklist-object-storage.md`
- `outputs/interview-card-object-storage.md`

## Exercises

1. **Easy** — Design the metadata record for user-uploaded videos with legal hold support.
2. **Medium** — Explain the state transitions for a multipart upload that partially fails.
3. **Hard** — Redesign the object storage service for cross-region disaster recovery without hiding the cost trade-offs.

## Further Reading

- [Amazon S3 multipart upload overview](https://docs.aws.amazon.com/AmazonS3/latest/userguide/mpuoverview.html) — concrete multipart workflow semantics  
- [Google Cloud Storage object metadata](https://cloud.google.com/storage/docs/metadata) — useful examples of control-plane metadata concerns  
