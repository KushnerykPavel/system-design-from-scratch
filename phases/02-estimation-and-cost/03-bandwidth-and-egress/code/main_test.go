package main

import "testing"

func TestEstimateBandwidth(t *testing.T) {
	got := EstimateBandwidth(BandwidthModel{
		PeakQPS:       120000,
		ResponseKB:    220,
		CacheHitRatio: 0.85,
		CostPerGB:     0.02,
	})

	if got.TotalGBPerSecond < 25.1 || got.TotalGBPerSecond > 25.3 {
		t.Fatalf("unexpected total gbps: %.2f", got.TotalGBPerSecond)
	}
	if got.OriginGBPerSecond < 3.75 || got.OriginGBPerSecond > 3.8 {
		t.Fatalf("unexpected origin gbps: %.2f", got.OriginGBPerSecond)
	}
	if got.MonthlyOriginCost <= 0 {
		t.Fatalf("expected positive cost, got %.2f", got.MonthlyOriginCost)
	}
}
