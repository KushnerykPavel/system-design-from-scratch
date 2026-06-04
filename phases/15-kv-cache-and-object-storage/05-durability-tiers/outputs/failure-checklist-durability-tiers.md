# Failure Checklist — Durability Tiers

- Durability claims reference real failure domains, not only copy counts.
- Repair deadlines are documented by tier and observable.
- Rebuild traffic is capacity-limited so it does not starve reads and writes.
- Geo-redundant tiers are tested for real restore, not just replica placement.
- Default tier selection cannot quietly underspecify critical data.
- Tier migration is idempotent and auditable.
- Cost dashboards are tied to logical bytes and durability class.
