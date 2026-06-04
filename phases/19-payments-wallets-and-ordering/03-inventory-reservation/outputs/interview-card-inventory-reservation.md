---
lesson: 03-inventory-reservation
focus: balanced
---

## Clarify first

- Is stock fungible or uniquely identified?
- How long can a reservation live?
- Is temporary oversell acceptable?

## Must-size numbers

- Peak reservation attempts per second
- Hottest SKU traffic concentration
- Reservation TTL and abandonment rate

## Failure probes

- What happens after payment success arrives late?
- How are leaked reservations detected?
- How is one hot SKU isolated from platform-wide damage?
