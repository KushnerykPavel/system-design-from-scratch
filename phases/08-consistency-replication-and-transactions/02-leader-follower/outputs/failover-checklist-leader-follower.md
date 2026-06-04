# Failover Checklist - Leader-Follower

- Is the current leader fenced so only one writer can accept commits?
- What is the candidate replica's commit index and applied index?
- Which recent writes were acknowledged under the current commit policy?
- Are follower reads temporarily restricted during promotion?
- Will clients need leader discovery or epoch refresh after failover?
- Are lag and promotion decisions logged with replica identity and epoch?
- What product paths must degrade rather than serve stale data?
