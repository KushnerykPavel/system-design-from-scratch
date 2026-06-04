package main

import "testing"

func TestSimulate(t *testing.T) {
	t.Parallel()

	trace := []string{"a", "b", "a", "c", "a", "b", "d", "a", "b"}
	tests := []struct {
		name     string
		policy   Policy
		capacity int
		hits     int
		misses   int
	}{
		{name: "fifo capacity 2", policy: PolicyFIFO, capacity: 2, hits: 1, misses: 8},
		{name: "lru capacity 2", policy: PolicyLRU, capacity: 2, hits: 2, misses: 7},
		{name: "lfu capacity 2", policy: PolicyLFU, capacity: 2, hits: 3, misses: 6},
		{name: "zero capacity", policy: PolicyLRU, capacity: 0, hits: 0, misses: len(trace)},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Simulate(trace, tt.capacity, tt.policy)
			if got.Hits != tt.hits || got.Misses != tt.misses {
				t.Fatalf("Simulate() = hits %d misses %d, want hits %d misses %d", got.Hits, got.Misses, tt.hits, tt.misses)
			}
		})
	}
}

func TestHitRate(t *testing.T) {
	t.Parallel()

	got := Result{Hits: 3, Misses: 1}.HitRate()
	if got != 0.75 {
		t.Fatalf("HitRate() = %v, want 0.75", got)
	}
}
