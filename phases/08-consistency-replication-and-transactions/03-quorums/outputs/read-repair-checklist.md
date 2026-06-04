# Read Repair Checklist

- What metadata identifies the newest or most correct version?
- Can the application tolerate last-write-wins, or is a domain merge needed?
- Does repair happen inline on reads, in the background, or both?
- How is repair throttled so incidents do not create a second outage?
- Are delete tombstones retained long enough to prevent resurrection?
- Which metrics show growing divergence before users notice?
