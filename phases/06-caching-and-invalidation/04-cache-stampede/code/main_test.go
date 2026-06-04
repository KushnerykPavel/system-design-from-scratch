package main

import "testing"

func TestEstimateWithCoalescing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		burst Burst
		load  OriginLoad
	}{
		{
			name:  "no requests",
			burst: Burst{},
			load:  OriginLoad{},
		},
		{
			name:  "all requests join one refresh",
			burst: Burst{Requests: 10, ArrivalWindowMS: 200, FetchLatencyMS: 50},
			load:  OriginLoad{OriginFetches: 1, Waiters: 9},
		},
		{
			name:  "multiple batches while refresh is still in flight",
			burst: Burst{Requests: 10, ArrivalWindowMS: 50, FetchLatencyMS: 200},
			load:  OriginLoad{OriginFetches: 4, Waiters: 6},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := EstimateWithCoalescing(tt.burst)
			if got != tt.load {
				t.Fatalf("EstimateWithCoalescing() = %+v, want %+v", got, tt.load)
			}
		})
	}
}

func TestEstimateWithoutCoalescing(t *testing.T) {
	t.Parallel()

	got := EstimateWithoutCoalescing(Burst{Requests: 8})
	if got.OriginFetches != 8 || got.Waiters != 0 {
		t.Fatalf("EstimateWithoutCoalescing() = %+v, want 8 origin fetches and 0 waiters", got)
	}
}
