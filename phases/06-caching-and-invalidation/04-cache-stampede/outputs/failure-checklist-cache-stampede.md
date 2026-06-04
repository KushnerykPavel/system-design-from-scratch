# Failure Checklist: Cache Stampede

- What happens when the hottest key expires?
- How many origin fetches can one expired key trigger?
- Can multiple waiters join one in-flight refresh?
- Are TTLs synchronized in a way that creates periodic refill spikes?
- Can slightly stale data be served while refresh is happening?
- What happens if refresh itself is slow or failing?
- Which metrics reveal miss amplification before the origin fails?
