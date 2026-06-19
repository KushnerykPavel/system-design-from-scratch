package main

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// LocationEvent represents a single GPS update from a driver's mobile device.
type LocationEvent struct {
	DriverID       string
	Lat            float64
	Lon            float64
	Speed          float64 // km/h, as reported by device
	AccuracyMeters float64 // GPS accuracy radius (68% confidence)
	Timestamp      time.Time
}

// LocationStore maps driver IDs to their most recent accepted location event.
type LocationStore map[string]LocationEvent

// Validation error sentinels.
var (
	ErrLowAccuracy   = errors.New("GPS accuracy too low (> 50m)")
	ErrStaleEvent    = errors.New("event timestamp too old (> 5s)")
	ErrTeleportation = errors.New("implied speed exceeds 300 km/h (teleportation detected)")
)

// haversineKm returns the great-circle distance in kilometres between two lat/lon points.
func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

// ValidateEvent checks three quality gates:
//  1. Accuracy must be <= 50m.
//  2. Event must not be older than 5 seconds relative to now.
//  3. Implied speed from the last stored position must not exceed 300 km/h.
//
// The store parameter provides the last known position for velocity check.
// Pass a nil/empty store to skip the velocity check (e.g., first event from driver).
func ValidateEvent(store LocationStore, event LocationEvent, now time.Time) error {
	// Gate 1: accuracy.
	if event.AccuracyMeters > 50 {
		return ErrLowAccuracy
	}

	// Gate 2: staleness.
	if now.Sub(event.Timestamp) > 5*time.Second {
		return ErrStaleEvent
	}

	// Gate 3: teleportation (only if we have a previous position).
	if prev, ok := store[event.DriverID]; ok {
		elapsed := event.Timestamp.Sub(prev.Timestamp).Hours()
		if elapsed > 0 {
			distKm := haversineKm(prev.Lat, prev.Lon, event.Lat, event.Lon)
			impliedSpeedKmh := distKm / elapsed
			if impliedSpeedKmh > 300 {
				return ErrTeleportation
			}
		}
	}

	return nil
}

// ProcessEvent validates and, if valid, stores the location event.
// Returns nil on success, or a validation error on failure.
// Invalid events are silently dropped (the last good position is preserved).
func ProcessEvent(store LocationStore, event LocationEvent, now time.Time) error {
	if err := ValidateEvent(store, event, now); err != nil {
		return err
	}
	store[event.DriverID] = event
	return nil
}

// IsStale returns true if the driver's last known position is older than threshold,
// or if the driver has no position in the store.
func IsStale(store LocationStore, driverID string, threshold time.Duration, now time.Time) bool {
	event, ok := store[driverID]
	if !ok {
		return true
	}
	return now.Sub(event.Timestamp) > threshold
}

func main() {
	store := make(LocationStore)

	// Use a fixed simulation clock so output is reproducible.
	// Events are created with timestamps relative to t0.
	t0 := time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC)

	type scenario struct {
		ev  LocationEvent
		now time.Time // server arrival time for this event
	}

	scenarios := []scenario{
		// Valid first event for driver_1 — arrives immediately.
		{
			ev:  LocationEvent{DriverID: "driver_1", Lat: 40.7128, Lon: -74.0060, Speed: 30, AccuracyMeters: 8, Timestamp: t0},
			now: t0,
		},
		// Valid second event for driver_1 — 4 seconds later, small movement.
		{
			ev:  LocationEvent{DriverID: "driver_1", Lat: 40.7135, Lon: -74.0055, Speed: 32, AccuracyMeters: 10, Timestamp: t0.Add(4 * time.Second)},
			now: t0.Add(4 * time.Second),
		},
		// Low accuracy — should be rejected.
		{
			ev:  LocationEvent{DriverID: "driver_2", Lat: 40.7200, Lon: -74.0100, Speed: 10, AccuracyMeters: 80, Timestamp: t0},
			now: t0,
		},
		// Stale event — device clock drifted; event is 10s behind server time.
		{
			ev:  LocationEvent{DriverID: "driver_3", Lat: 40.7300, Lon: -73.9900, Speed: 20, AccuracyMeters: 15, Timestamp: t0.Add(-10 * time.Second)},
			now: t0,
		},
		// Valid event for driver_4.
		{
			ev:  LocationEvent{DriverID: "driver_4", Lat: 40.6900, Lon: -74.0200, Speed: 45, AccuracyMeters: 5, Timestamp: t0},
			now: t0,
		},
		// Teleportation — driver_1 "jumps" to London 4 seconds after last event.
		{
			ev:  LocationEvent{DriverID: "driver_1", Lat: 51.5074, Lon: -0.1278, Speed: 0, AccuracyMeters: 10, Timestamp: t0.Add(8 * time.Second)},
			now: t0.Add(8 * time.Second),
		},
	}

	for _, s := range scenarios {
		err := ProcessEvent(store, s.ev, s.now)
		status := "accepted"
		if err != nil {
			status = "REJECTED: " + err.Error()
		}
		fmt.Printf("%-10s lat=%.4f lon=%.4f acc=%.0fm  → %s\n",
			s.ev.DriverID, s.ev.Lat, s.ev.Lon, s.ev.AccuracyMeters, status)
	}

	fmt.Println()
	checkNow := t0.Add(35 * time.Second)
	staleThreshold := 30 * time.Second
	for _, id := range []string{"driver_1", "driver_2", "driver_3", "driver_4"} {
		stale := IsStale(store, id, staleThreshold, checkNow)
		fmt.Printf("%-10s stale=%v (checked 35s after t0)\n", id, stale)
	}
}
