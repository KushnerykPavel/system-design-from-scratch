package main

import "testing"

func TestEstimateStorage(t *testing.T) {
	got := EstimateStorage(StoragePlan{
		EventsPerDay:       400000000,
		BytesPerEvent:      1200,
		ReplicationFactor:  3,
		RetentionDays:      30,
		IndexOverheadRatio: 0,
	})

	if got.DailyRawGB < 446 || got.DailyRawGB > 448 {
		t.Fatalf("unexpected daily raw gb: %.2f", got.DailyRawGB)
	}
	if got.DailyDurableGB < 1339 || got.DailyDurableGB > 1342 {
		t.Fatalf("unexpected daily durable gb: %.2f", got.DailyDurableGB)
	}
	if got.RetainedTotalGB < 40170 || got.RetainedTotalGB > 40260 {
		t.Fatalf("unexpected retained gb: %.2f", got.RetainedTotalGB)
	}
}
