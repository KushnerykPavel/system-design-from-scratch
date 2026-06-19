package main

import "testing"

func seedOrigin() map[string][]Association {
	return map[string][]Association{
		"alice:friend": {
			{ID1: "alice", AType: "friend", ID2: "bob", Time: 1700000000},
		},
	}
}

func TestAssocQueryCacheHit(t *testing.T) {
	cache := NewTaoCache(seedOrigin())

	// First query populates caches from origin.
	r1 := cache.AssocQuery("alice", "friend")
	if r1.ServedFrom != CacheLevelOrigin {
		t.Fatalf("expected first query to be served from origin, got %s", r1.ServedFrom)
	}

	// Second query must hit follower cache.
	r2 := cache.AssocQuery("alice", "friend")
	if r2.ServedFrom != CacheLevelFollower {
		t.Fatalf("expected second query to be served from follower, got %s", r2.ServedFrom)
	}
	if cache.Stats.FollowerHits != 1 {
		t.Fatalf("expected 1 follower hit, got %d", cache.Stats.FollowerHits)
	}
}

func TestAssocQueryCacheMissFallsBackToLeader(t *testing.T) {
	cache := NewTaoCache(seedOrigin())

	// Populate leader cache only (bypass follower by writing directly).
	key := cacheKey("alice", "friend")
	cache.leaderCache[key] = seedOrigin()[key]

	// Follower cache is empty, leader cache is populated.
	r := cache.AssocQuery("alice", "friend")
	if r.ServedFrom != CacheLevelLeader {
		t.Fatalf("expected query to fall back to leader, got %s", r.ServedFrom)
	}
	if cache.Stats.LeaderHits != 1 {
		t.Fatalf("expected 1 leader hit, got %d", cache.Stats.LeaderHits)
	}
	// Follower should now be populated after leader hit.
	if _, ok := cache.followerCache[key]; !ok {
		t.Fatal("expected follower cache to be populated after leader hit")
	}
}

func TestInvalidateAndRefetchFromOrigin(t *testing.T) {
	cache := NewTaoCache(seedOrigin())

	// Warm both caches.
	cache.AssocQuery("alice", "friend")
	// Add a new association via Write (which invalidates caches).
	cache.Write(Association{ID1: "alice", AType: "friend", ID2: "carol", Time: 1700000020})

	// Next query must go to origin because caches were invalidated.
	r := cache.AssocQuery("alice", "friend")
	if r.ServedFrom != CacheLevelOrigin {
		t.Fatalf("expected query after invalidation to be served from origin, got %s", r.ServedFrom)
	}
	// Result must include both bob and carol.
	if len(r.Assocs) != 2 {
		t.Fatalf("expected 2 associations after write, got %d", len(r.Assocs))
	}
}
