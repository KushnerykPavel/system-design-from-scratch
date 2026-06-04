# Stock Exchange Order Path Review

## Correctness Boundary

- What event makes an order officially accepted?
- What has already been persisted at that moment?
- What sequencing point defines price-time priority?

## Core Isolation

- Is matching isolated per symbol or per product class?
- What happens when one symbol gets dramatically hotter?
- Can market-data lag hurt execution latency?

## Recovery

- What is the authoritative journal?
- How often are snapshots taken?
- How do you verify replay produces the same book?

## Senior Follow-Ups

- What is the cancel path latency target?
- Where do shared account-risk checks live?
- Which guarantees remain if latency budgets tighten further?
