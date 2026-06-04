package main

import "testing"

func TestEstimateQueue(t *testing.T) {
	got := EstimateQueue(QueueModel{
		ArrivalRate:         60000,
		ServiceRate:         40000,
		BurstSeconds:        600,
		RecoveryServiceRate: 40000,
		PostBurstArrival:    30000,
	})

	if got.BacklogGrowthPerSec != 20000 {
		t.Fatalf("unexpected backlog growth: %.2f", got.BacklogGrowthPerSec)
	}
	if got.BacklogItems != 12000000 {
		t.Fatalf("unexpected backlog items: %.2f", got.BacklogItems)
	}
	if got.DrainSeconds != 1200 {
		t.Fatalf("unexpected drain seconds: %.2f", got.DrainSeconds)
	}
}
