# Interview Card — Storage Growth

## Default formula

- raw ingest per day = events per day x bytes per event
- durable ingest per day = raw ingest x replicas x overhead
- retained size = durable ingest per day x retention days

## Prompts

- "I want to separate raw payload, indexes, and replicas before I choose storage."
- "Retention matters as much as ingest rate here."
- "If only a short window needs fast reads, I’d tier the rest down aggressively."

## Common misses

- forgetting backups
- forgetting secondary indexes
- ignoring deletion lag
- keeping cold data on expensive hot disks
