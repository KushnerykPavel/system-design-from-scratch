package main

import (
	"errors"
	"fmt"
	"time"
)

// TripState represents one phase in a trip's lifecycle.
type TripState string

const (
	StateRequested        TripState = "REQUESTED"
	StateMatching         TripState = "MATCHING"
	StateDriverAssigned   TripState = "DRIVER_ASSIGNED"
	StateDriverEnRoute    TripState = "DRIVER_EN_ROUTE"
	StateDriverArrived    TripState = "DRIVER_ARRIVED"
	StateInProgress       TripState = "IN_PROGRESS"
	StateCompleted        TripState = "COMPLETED"
	StateCancelledByRider  TripState = "CANCELLED_BY_RIDER"
	StateCancelledByDriver TripState = "CANCELLED_BY_DRIVER"
	StateCancelledBySystem TripState = "CANCELLED_BY_SYSTEM"
)

// terminalStates is the set of states from which no further transition is allowed.
var terminalStates = map[TripState]bool{
	StateCompleted:         true,
	StateCancelledByRider:  true,
	StateCancelledByDriver: true,
	StateCancelledBySystem: true,
}

// validTransitions defines all permitted state transitions.
// Any (from, to) pair not in this map is rejected.
var validTransitions = map[TripState][]TripState{
	StateRequested: {
		StateMatching,
		StateCancelledBySystem,
	},
	StateMatching: {
		StateDriverAssigned,
		StateCancelledBySystem,
	},
	StateDriverAssigned: {
		StateDriverEnRoute,
		StateCancelledByRider,
		StateCancelledByDriver,
	},
	StateDriverEnRoute: {
		StateDriverArrived,
		StateCancelledByRider,
		StateCancelledByDriver,
	},
	StateDriverArrived: {
		StateInProgress,
		StateCancelledByRider,
	},
	StateInProgress: {
		StateCompleted,
		StateCancelledBySystem,
	},
}

// Trip holds the mutable state for a single ride.
type Trip struct {
	ID            string
	State         TripState
	DriverID      string
	RiderID       string
	LastHeartbeat time.Time
	AssignedAt    time.Time // set when entering DRIVER_ASSIGNED
}

// ErrInvalidTransition is returned when a requested transition is not permitted.
var ErrInvalidTransition = errors.New("invalid state transition")

// isValidTransition returns true if transitioning from current to next is permitted.
func isValidTransition(current, next TripState) bool {
	for _, allowed := range validTransitions[current] {
		if allowed == next {
			return true
		}
	}
	return false
}

// Transition attempts to move the trip to newState.
//   - If the trip is already in newState, it is a no-op (idempotency).
//   - If the transition is not valid, ErrInvalidTransition is returned.
//   - If the trip is in a terminal state and newState differs, ErrInvalidTransition is returned.
func Transition(trip *Trip, newState TripState) error {
	// Idempotency: same state is always a no-op.
	if trip.State == newState {
		return nil
	}

	// Terminal states cannot transition further.
	if terminalStates[trip.State] {
		return fmt.Errorf("%w: %s is a terminal state", ErrInvalidTransition, trip.State)
	}

	if !isValidTransition(trip.State, newState) {
		return fmt.Errorf("%w: %s → %s", ErrInvalidTransition, trip.State, newState)
	}

	trip.State = newState
	return nil
}

// IsDeadTrip returns true when the trip is IN_PROGRESS and has not received
// a heartbeat update within the given threshold duration.
func IsDeadTrip(trip Trip, threshold time.Duration, now time.Time) bool {
	if trip.State != StateInProgress {
		return false
	}
	return now.Sub(trip.LastHeartbeat) > threshold
}

func main() {
	// --- Happy path simulation ---
	fmt.Println("=== Happy Path ===")
	trip := &Trip{
		ID:            "trip_001",
		DriverID:      "driver_42",
		RiderID:       "rider_99",
		State:         StateRequested,
		LastHeartbeat: time.Now(),
	}

	transitions := []TripState{
		StateMatching,
		StateDriverAssigned,
		StateDriverEnRoute,
		StateDriverArrived,
		StateInProgress,
		StateCompleted,
	}

	for _, next := range transitions {
		if err := Transition(trip, next); err != nil {
			fmt.Printf("  FAILED %s → %s: %v\n", trip.State, next, err)
		} else {
			fmt.Printf("  OK     → %s\n", trip.State)
		}
	}

	// --- Idempotency check ---
	fmt.Println("\n=== Idempotency: re-send COMPLETED ===")
	err := Transition(trip, StateCompleted)
	if err == nil {
		fmt.Println("  OK (no-op, idempotent)")
	} else {
		fmt.Printf("  ERROR: %v\n", err)
	}

	// --- Invalid transition from terminal state ---
	fmt.Println("\n=== Invalid: COMPLETED → DRIVER_ASSIGNED ===")
	err = Transition(trip, StateDriverAssigned)
	if err != nil {
		fmt.Printf("  Rejected as expected: %v\n", err)
	}

	// --- Dead trip detection ---
	fmt.Println("\n=== Dead Trip Detection ===")
	liveTrip := &Trip{
		ID:            "trip_002",
		State:         StateInProgress,
		LastHeartbeat: time.Now().Add(-5 * time.Minute),
	}
	deadTrip := &Trip{
		ID:            "trip_003",
		State:         StateInProgress,
		LastHeartbeat: time.Now().Add(-12 * time.Minute),
	}

	threshold := 10 * time.Minute
	now := time.Now()
	fmt.Printf("  trip_002 (last heartbeat 5m ago)  dead=%v\n", IsDeadTrip(*liveTrip, threshold, now))
	fmt.Printf("  trip_003 (last heartbeat 12m ago) dead=%v\n", IsDeadTrip(*deadTrip, threshold, now))
}
