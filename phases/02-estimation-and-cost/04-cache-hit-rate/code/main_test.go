package main

import "testing"

func TestEstimateCache(t *testing.T) {
	got := EstimateCache(CacheModel{
		ReadQPS:           90000,
		HitRatio:          0.92,
		MissLatencyMillis: 40,
	})

	if got.OriginQPS < 7199 || got.OriginQPS > 7201 {
		t.Fatalf("unexpected origin qps: %.2f", got.OriginQPS)
	}
	if got.CacheServedQPS < 82799 || got.CacheServedQPS > 82801 {
		t.Fatalf("unexpected cache served qps: %.2f", got.CacheServedQPS)
	}
	if got.AverageLatencyDelta < 3.1 || got.AverageLatencyDelta > 3.3 {
		t.Fatalf("unexpected latency delta: %.2f", got.AverageLatencyDelta)
	}
}
