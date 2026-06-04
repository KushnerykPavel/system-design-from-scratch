---
lesson: 07-outbox-and-cdc
focus: balanced
---

# Outbox Decision Guide

## Choose Transactional Outbox When

- The service owns the write transaction
- Domain events should be shaped in application code
- Per-aggregate ordering matters
- You want clear replay and publish ownership

## Choose CDC When

- Another team or platform owns the database transaction path
- You can access a reliable change stream
- Low-level row changes can be transformed safely downstream
- You accept more infra and envelope complexity

## Review Questions

- What is the event visibility lag target?
- How will duplicate publishes be handled?
- What is the retention and cleanup plan?
- How will relay backlog affect the primary database?
