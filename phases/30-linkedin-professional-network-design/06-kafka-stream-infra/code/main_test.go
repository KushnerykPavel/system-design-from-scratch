package main

import "testing"

func TestConsumerLag_CaughtUp(t *testing.T) {
	p := Partition{ID: 0, LatestOffset: 1000, CommittedOffset: 1000}
	if lag := ConsumerLag(p); lag != 0 {
		t.Errorf("expected 0 lag when caught up, got %d", lag)
	}
}

func TestConsumerLag_BehindByKnownAmount(t *testing.T) {
	p := Partition{ID: 0, LatestOffset: 5000, CommittedOffset: 3000}
	if lag := ConsumerLag(p); lag != 2000 {
		t.Errorf("expected lag of 2000, got %d", lag)
	}
}

func TestConsumerLag_NeverNegative(t *testing.T) {
	// CommittedOffset ahead of latest (shouldn't happen in practice, but must not panic)
	p := Partition{ID: 0, LatestOffset: 100, CommittedOffset: 150}
	if lag := ConsumerLag(p); lag != 0 {
		t.Errorf("expected 0 for committed ahead of latest, got %d", lag)
	}
}

func TestGroupLag_SumsAcrossPartitions(t *testing.T) {
	partitions := []Partition{
		{ID: 0, LatestOffset: 1000, CommittedOffset: 800},  // lag=200
		{ID: 1, LatestOffset: 2000, CommittedOffset: 1700}, // lag=300
		{ID: 2, LatestOffset: 500, CommittedOffset: 500},   // lag=0
	}
	total := GroupLag(partitions)
	if total != 500 {
		t.Errorf("expected total group lag of 500, got %d", total)
	}
}

func TestGroupLag_EmptyPartitions(t *testing.T) {
	if lag := GroupLag([]Partition{}); lag != 0 {
		t.Errorf("expected 0 for empty partition list, got %d", lag)
	}
}

func TestIsLagAlarm_BelowThreshold(t *testing.T) {
	if IsLagAlarm(4999, 5000) {
		t.Error("expected no alarm when lag is below threshold")
	}
}

func TestIsLagAlarm_AtThreshold(t *testing.T) {
	// Alarm triggers strictly above threshold
	if IsLagAlarm(5000, 5000) {
		t.Error("expected no alarm when lag equals threshold (not strictly greater)")
	}
}

func TestIsLagAlarm_AboveThreshold(t *testing.T) {
	if !IsLagAlarm(5001, 5000) {
		t.Error("expected alarm when lag exceeds threshold")
	}
}

func TestSimulateConsumerBehind_LagGrowsWhenBehind(t *testing.T) {
	partitions := []Partition{
		{ID: 0, LatestOffset: 1000, CommittedOffset: 1000},
	}
	initialLag := GroupLag(partitions)

	// Produce 500, process 200 — consumer falls behind
	updated := SimulateConsumerBehind(partitions, 500, 200)
	newLag := GroupLag(updated)

	if newLag <= initialLag {
		t.Errorf("expected lag to grow when message rate > process rate, got %d -> %d", initialLag, newLag)
	}
}

func TestSimulateConsumerBehind_LagShrinksDuringCatchUp(t *testing.T) {
	partitions := []Partition{
		{ID: 0, LatestOffset: 2000, CommittedOffset: 1000}, // 1000 lag
	}

	// Produce 100, process 400 — consumer catches up
	updated := SimulateConsumerBehind(partitions, 100, 400)
	newLag := GroupLag(updated)

	if newLag >= 1000 {
		t.Errorf("expected lag to shrink during catch-up, got %d", newLag)
	}
}

func TestSimulateConsumerBehind_CommittedNeverExceedsLatest(t *testing.T) {
	partitions := []Partition{
		{ID: 0, LatestOffset: 1000, CommittedOffset: 1000},
	}

	// Process rate (2000) far exceeds message rate (50)
	updated := SimulateConsumerBehind(partitions, 50, 2000)

	if updated[0].CommittedOffset > updated[0].LatestOffset {
		t.Errorf("committed offset (%d) should never exceed latest offset (%d)",
			updated[0].CommittedOffset, updated[0].LatestOffset)
	}
}
