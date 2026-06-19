package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// CacheLevel identifies which tier of the TAO two-level cache served a lookup.
type CacheLevel string

const (
	CacheLevelFollower CacheLevel = "follower"
	CacheLevelLeader   CacheLevel = "leader"
	CacheLevelOrigin   CacheLevel = "origin" // simulated MySQL shard
)

// Association represents a directed edge in the social graph.
type Association struct {
	ID1   string `json:"id1"`
	AType string `json:"atype"`
	ID2   string `json:"id2"`
	Time  int64  `json:"time"`
}

// LookupResult is returned by an AssocQuery call.
type LookupResult struct {
	ID1        string        `json:"id1"`
	AType      string        `json:"atype"`
	Assocs     []Association `json:"assocs"`
	ServedFrom CacheLevel    `json:"served_from"`
}

// Stats tracks cache hit/miss counts across the simulation.
type Stats struct {
	FollowerHits   int `json:"follower_hits"`
	FollowerMisses int `json:"follower_misses"`
	LeaderHits     int `json:"leader_hits"`
	LeaderMisses   int `json:"leader_misses"`
	OriginFetches  int `json:"origin_fetches"`
}

// HitRate returns the follower cache hit rate as a percentage.
func (s Stats) HitRate() int {
	total := s.FollowerHits + s.FollowerMisses
	if total == 0 {
		return 0
	}
	return (s.FollowerHits * 100) / total
}

// TaoCache simulates a TAO two-level cache (follower + leader) backed by a simulated origin.
type TaoCache struct {
	followerCache map[string][]Association // keyed by "id1:atype"
	leaderCache   map[string][]Association // keyed by "id1:atype"
	origin        map[string][]Association // simulated MySQL shard
	Stats         Stats
}

// NewTaoCache creates a TaoCache pre-populated with origin data.
func NewTaoCache(originData map[string][]Association) *TaoCache {
	return &TaoCache{
		followerCache: make(map[string][]Association),
		leaderCache:   make(map[string][]Association),
		origin:        originData,
	}
}

func cacheKey(id1, atype string) string {
	return id1 + ":" + atype
}

// AssocQuery performs a (id1, atype) association lookup against the TAO cache hierarchy.
// Read path: follower → leader → origin (simulated MySQL).
func (t *TaoCache) AssocQuery(id1, atype string) LookupResult {
	key := cacheKey(id1, atype)

	// Level 1: follower cache
	if assocs, ok := t.followerCache[key]; ok {
		t.Stats.FollowerHits++
		return LookupResult{ID1: id1, AType: atype, Assocs: assocs, ServedFrom: CacheLevelFollower}
	}
	t.Stats.FollowerMisses++

	// Level 2: leader cache
	if assocs, ok := t.leaderCache[key]; ok {
		t.Stats.LeaderHits++
		// Populate follower cache on leader hit (cache fill).
		t.followerCache[key] = assocs
		return LookupResult{ID1: id1, AType: atype, Assocs: assocs, ServedFrom: CacheLevelLeader}
	}
	t.Stats.LeaderMisses++

	// Level 3: origin (simulated MySQL shard)
	t.Stats.OriginFetches++
	assocs := t.origin[key] // empty slice if not found — correct behavior
	// Populate both leader and follower caches.
	t.leaderCache[key] = assocs
	t.followerCache[key] = assocs
	return LookupResult{ID1: id1, AType: atype, Assocs: assocs, ServedFrom: CacheLevelOrigin}
}

// Invalidate simulates a TAO leader invalidating follower and leader caches on a write.
// In real TAO: write → MySQL → leader cache update → follower invalidation messages.
func (t *TaoCache) Invalidate(id1, atype string) {
	key := cacheKey(id1, atype)
	delete(t.followerCache, key)
	delete(t.leaderCache, key)
}

// Write adds an association to the origin and invalidates caches (simulates a TAO write path).
func (t *TaoCache) Write(assoc Association) {
	key := cacheKey(assoc.ID1, assoc.AType)
	t.origin[key] = append(t.origin[key], assoc)
	t.Invalidate(assoc.ID1, assoc.AType)
}

func main() {
	// Seed origin with sample social graph data.
	origin := map[string][]Association{
		"alice:friend": {
			{ID1: "alice", AType: "friend", ID2: "bob", Time: 1700000000},
			{ID1: "alice", AType: "friend", ID2: "carol", Time: 1700000010},
		},
		"bob:friend": {
			{ID1: "bob", AType: "friend", ID2: "alice", Time: 1700000000},
		},
	}

	cache := NewTaoCache(origin)

	// First query: cold — should fetch from origin.
	r1 := cache.AssocQuery("alice", "friend")
	// Second query: warm — should hit follower cache.
	r2 := cache.AssocQuery("alice", "friend")
	// Write a new association: invalidates cache.
	cache.Write(Association{ID1: "alice", AType: "friend", ID2: "dave", Time: 1700000020})
	// Third query after write: cold again — should fetch from origin.
	r3 := cache.AssocQuery("alice", "friend")

	output := map[string]any{
		"query_1": r1,
		"query_2": r2,
		"query_3": r3,
		"stats": map[string]any{
			"follower_hits":   cache.Stats.FollowerHits,
			"follower_misses": cache.Stats.FollowerMisses,
			"leader_hits":     cache.Stats.LeaderHits,
			"leader_misses":   cache.Stats.LeaderMisses,
			"origin_fetches":  cache.Stats.OriginFetches,
			"follower_hit_rate_pct": cache.Stats.HitRate(),
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
