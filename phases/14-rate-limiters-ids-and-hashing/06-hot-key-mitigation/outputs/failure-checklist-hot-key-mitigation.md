# Failure Checklist — Hot-Key Mitigation

- Did one key or tenant dominate traffic long before fleet-wide saturation?
- Was the hotspot read-heavy or write-heavy, and did mitigation match that shape?
- Did replication create stale-read exposure that was not monitored?
- Did request coalescing or isolation activate quickly enough?
- Did the mitigation stay targeted, or did it raise cost for the whole fleet?

Healthy answer pattern:

- detect
- classify
- mitigate selectively
- preserve explainability
- remove mitigation cleanly after the event
