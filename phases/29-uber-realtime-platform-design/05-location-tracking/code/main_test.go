package main

import (
	"errors"
	"testing"
	"time"
)

var baseNow = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

// validEvent returns a well-formed LocationEvent relative to baseNow.
func validEvent(driverID string, lat, lon float64, offset time.Duration) LocationEvent {
	return LocationEvent{
		DriverID:       driverID,
		Lat:            lat,
		Lon:            lon,
		Speed:          30,
		AccuracyMeters: 10,
		Timestamp:      baseNow.Add(offset),
	}
}

// TestValidEventIsStored verifies that a passing event lands in the store.
func TestValidEventIsStored(t *testing.T) {
	store := make(LocationStore)
	ev := validEvent("d1", 40.71, -74.00, 0)
	err := ProcessEvent(store, ev, baseNow)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, ok := store["d1"]; !ok {
		t.Fatal("expected event to be stored in LocationStore")
	}
}

// TestLowAccuracyRejected verifies that events with accuracy > 50m are dropped.
func TestLowAccuracyRejected(t *testing.T) {
	store := make(LocationStore)
	ev := validEvent("d2", 40.71, -74.00, 0)
	ev.AccuracyMeters = 75
	err := ProcessEvent(store, ev, baseNow)
	if !errors.Is(err, ErrLowAccuracy) {
		t.Fatalf("expected ErrLowAccuracy, got %v", err)
	}
	if _, ok := store["d2"]; ok {
		t.Fatal("rejected event should not be stored")
	}
}

// TestStaleEventRejected verifies that events older than 5 seconds are dropped.
func TestStaleEventRejected(t *testing.T) {
	store := make(LocationStore)
	ev := validEvent("d3", 40.71, -74.00, -10*time.Second)
	err := ProcessEvent(store, ev, baseNow)
	if !errors.Is(err, ErrStaleEvent) {
		t.Fatalf("expected ErrStaleEvent, got %v", err)
	}
}

// TestTeleportationRejected verifies that an impossibly fast position change is dropped.
func TestTeleportationRejected(t *testing.T) {
	store := make(LocationStore)

	// First valid event: Manhattan.
	first := validEvent("d4", 40.7128, -74.0060, 0)
	if err := ProcessEvent(store, first, baseNow); err != nil {
		t.Fatalf("first event failed: %v", err)
	}

	// Second event 4 seconds later: London — ~5500 km away.
	second := LocationEvent{
		DriverID:       "d4",
		Lat:            51.5074,
		Lon:            -0.1278,
		Speed:          0,
		AccuracyMeters: 10,
		Timestamp:      baseNow.Add(4 * time.Second),
	}
	err := ProcessEvent(store, second, baseNow.Add(4*time.Second))
	if !errors.Is(err, ErrTeleportation) {
		t.Fatalf("expected ErrTeleportation, got %v", err)
	}

	// Previous good position should still be in store.
	if stored, ok := store["d4"]; !ok || stored.Lat != first.Lat {
		t.Fatal("teleportation rejection should preserve the last good position")
	}
}

// TestStaleDriverDetection verifies IsStale returns true when position is old.
func TestStaleDriverDetection(t *testing.T) {
	store := make(LocationStore)
	ev := validEvent("d5", 40.71, -74.00, 0)
	store["d5"] = ev

	// 31 seconds later — past the 30s threshold.
	future := baseNow.Add(31 * time.Second)
	if !IsStale(store, "d5", 30*time.Second, future) {
		t.Error("driver with 31s old position should be considered stale")
	}
}

// TestFreshDriverNotStale verifies IsStale returns false for a recent position.
func TestFreshDriverNotStale(t *testing.T) {
	store := make(LocationStore)
	ev := validEvent("d6", 40.71, -74.00, 0)
	store["d6"] = ev

	// 5 seconds later — well within the 30s threshold.
	slightly := baseNow.Add(5 * time.Second)
	if IsStale(store, "d6", 30*time.Second, slightly) {
		t.Error("driver with 5s old position should not be stale")
	}
}

// TestUnknownDriverIsStale verifies that a driver with no entry is reported stale.
func TestUnknownDriverIsStale(t *testing.T) {
	store := make(LocationStore)
	if !IsStale(store, "unknown", 30*time.Second, baseNow) {
		t.Error("driver with no location entry should be considered stale")
	}
}

// TestBoundaryAccuracy verifies that exactly 50m accuracy is accepted.
func TestBoundaryAccuracy(t *testing.T) {
	store := make(LocationStore)
	ev := validEvent("d7", 40.71, -74.00, 0)
	ev.AccuracyMeters = 50
	if err := ProcessEvent(store, ev, baseNow); err != nil {
		t.Fatalf("accuracy exactly 50m should be accepted, got %v", err)
	}
}
