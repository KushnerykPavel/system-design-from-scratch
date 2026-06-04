package main

import "testing"

func TestEstimateTraffic(t *testing.T) {
	got := EstimateTraffic(TrafficModel{
		DAU:                   5000000,
		RequestsPerUserPerDay: 24,
		ReadRatio:             0.8,
		PeakFactor:            6,
	})

	if got.AverageQPS < 1388 || got.AverageQPS > 1390 {
		t.Fatalf("unexpected average qps: %.2f", got.AverageQPS)
	}
	if got.PeakQPS < 8332 || got.PeakQPS > 8335 {
		t.Fatalf("unexpected peak qps: %.2f", got.PeakQPS)
	}
	if got.PeakReadQPS < 6665 || got.PeakReadQPS > 6668 {
		t.Fatalf("unexpected peak read qps: %.2f", got.PeakReadQPS)
	}
}
