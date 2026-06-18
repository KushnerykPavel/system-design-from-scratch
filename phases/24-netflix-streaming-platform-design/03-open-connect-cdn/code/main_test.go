package main

import "testing"

func TestOCACacheHitAndMiss(t *testing.T) {
	cache := NewOCACache(1024 * 1024 * 1024)
	cache.Put("tt-001", "h264/720p", 100*1024*1024)
	if !cache.Get("tt-001", "h264/720p") {
		t.Fatal("expected cache hit for tt-001")
	}
	if cache.Get("tt-999", "h264/720p") {
		t.Fatal("expected cache miss for tt-999")
	}
}

func TestOCACacheEvictsLeastPopular(t *testing.T) {
	// 200 MB cache: fill with two 100 MB entries, then add a third.
	cache := NewOCACache(200 * 1024 * 1024)
	cache.Put("tt-low", "h264/720p", 100*1024*1024)   // access count 1
	cache.Put("tt-high", "h264/1080p", 100*1024*1024) // access count 1
	// Boost tt-high popularity.
	cache.Get("tt-high", "h264/1080p")
	cache.Get("tt-high", "h264/1080p")
	// Adding tt-new should evict tt-low (lower access count).
	cache.Put("tt-new", "h264/480p", 100*1024*1024)
	if cache.Get("tt-low", "h264/720p") {
		t.Fatal("expected tt-low to be evicted")
	}
	if !cache.Get("tt-high", "h264/1080p") {
		t.Fatal("expected tt-high to survive eviction")
	}
}

func TestOCACacheStatsConsistent(t *testing.T) {
	cache := NewOCACache(500 * 1024 * 1024)
	cache.Put("tt-a", "h264/720p", 200*1024*1024)
	cache.Put("tt-b", "h264/480p", 150*1024*1024)
	stats := cache.Stats()
	if stats["entry_count"].(int) != 2 {
		t.Fatalf("expected 2 entries, got %v", stats["entry_count"])
	}
	if stats["fill_percent"].(int) != 70 {
		t.Fatalf("expected 70%% fill, got %v", stats["fill_percent"])
	}
}
