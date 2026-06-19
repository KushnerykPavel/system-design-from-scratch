package main

import (
	"errors"
	"testing"
	"time"
)

func newTrip(state TripState) *Trip {
	return &Trip{
		ID:            "trip_test",
		State:         state,
		DriverID:      "driver_1",
		RiderID:       "rider_1",
		LastHeartbeat: time.Now(),
	}
}

// TestValidTransition verifies a simple permitted transition succeeds.
func TestValidTransition(t *testing.T) {
	trip := newTrip(StateRequested)
	if err := Transition(trip, StateMatching); err != nil {
		t.Fatalf("expected no error for REQUESTED → MATCHING, got %v", err)
	}
	if trip.State != StateMatching {
		t.Errorf("expected state MATCHING, got %s", trip.State)
	}
}

// TestInvalidTransition verifies that skipping states returns an error.
func TestInvalidTransition(t *testing.T) {
	trip := newTrip(StateRequested)
	err := Transition(trip, StateCompleted) // must skip multiple states
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("expected ErrInvalidTransition, got %v", err)
	}
	// State must not have changed.
	if trip.State != StateRequested {
		t.Errorf("state should remain REQUESTED after invalid transition, got %s", trip.State)
	}
}

// TestSameStateIsNoOp verifies idempotency: transitioning to current state returns nil.
func TestSameStateIsNoOp(t *testing.T) {
	trip := newTrip(StateInProgress)
	if err := Transition(trip, StateInProgress); err != nil {
		t.Fatalf("same-state transition should be a no-op, got %v", err)
	}
	if trip.State != StateInProgress {
		t.Errorf("state should remain IN_PROGRESS, got %s", trip.State)
	}
}

// TestIdempotentCompletedTransition verifies that re-sending COMPLETED is a no-op.
func TestIdempotentCompletedTransition(t *testing.T) {
	trip := newTrip(StateInProgress)
	// First transition to COMPLETED.
	if err := Transition(trip, StateCompleted); err != nil {
		t.Fatalf("first COMPLETED transition failed: %v", err)
	}
	// Second transition — must be idempotent.
	if err := Transition(trip, StateCompleted); err != nil {
		t.Fatalf("second COMPLETED transition should be no-op, got %v", err)
	}
}

// TestTerminalStateBlocksFurtherTransitions verifies no transition from terminal state.
func TestTerminalStateBlocksFurtherTransitions(t *testing.T) {
	terminalStatesSlice := []TripState{
		StateCompleted,
		StateCancelledByRider,
		StateCancelledByDriver,
		StateCancelledBySystem,
	}
	for _, terminal := range terminalStatesSlice {
		trip := newTrip(terminal)
		err := Transition(trip, StateRequested) // try any other state
		if !errors.Is(err, ErrInvalidTransition) {
			t.Errorf("expected ErrInvalidTransition from terminal state %s, got %v", terminal, err)
		}
	}
}

// TestDeadTripDetectionExceedsThreshold verifies a trip with old heartbeat is dead.
func TestDeadTripDetectionExceedsThreshold(t *testing.T) {
	trip := Trip{
		ID:            "trip_dead",
		State:         StateInProgress,
		LastHeartbeat: time.Now().Add(-11 * time.Minute),
	}
	if !IsDeadTrip(trip, 10*time.Minute, time.Now()) {
		t.Error("trip with 11m old heartbeat should be detected as dead")
	}
}

// TestDeadTripDetectionWithinThreshold verifies a trip with recent heartbeat is alive.
func TestDeadTripDetectionWithinThreshold(t *testing.T) {
	trip := Trip{
		ID:            "trip_alive",
		State:         StateInProgress,
		LastHeartbeat: time.Now().Add(-5 * time.Minute),
	}
	if IsDeadTrip(trip, 10*time.Minute, time.Now()) {
		t.Error("trip with 5m old heartbeat should NOT be detected as dead")
	}
}

// TestDeadTripNotInProgress verifies non-IN_PROGRESS trips are never dead.
func TestDeadTripNotInProgress(t *testing.T) {
	states := []TripState{StateRequested, StateMatching, StateDriverAssigned, StateCompleted}
	for _, s := range states {
		trip := Trip{
			ID:            "trip_x",
			State:         s,
			LastHeartbeat: time.Now().Add(-30 * time.Minute),
		}
		if IsDeadTrip(trip, 10*time.Minute, time.Now()) {
			t.Errorf("state %s should never be a dead trip", s)
		}
	}
}

// TestHappyPathFullTransition exercises the complete happy path.
func TestHappyPathFullTransition(t *testing.T) {
	trip := newTrip(StateRequested)
	path := []TripState{
		StateMatching,
		StateDriverAssigned,
		StateDriverEnRoute,
		StateDriverArrived,
		StateInProgress,
		StateCompleted,
	}
	for _, next := range path {
		if err := Transition(trip, next); err != nil {
			t.Fatalf("happy path transition to %s failed: %v", next, err)
		}
	}
	if trip.State != StateCompleted {
		t.Errorf("expected COMPLETED at end of happy path, got %s", trip.State)
	}
}
