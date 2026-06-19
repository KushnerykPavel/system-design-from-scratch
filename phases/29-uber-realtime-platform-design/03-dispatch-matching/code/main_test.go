package main

import (
	"math"
	"testing"
)

func TestDistanceZero(t *testing.T) {
	d := Distance(40.71, -74.01, 40.71, -74.01)
	if d != 0.0 {
		t.Fatalf("expected 0 distance for identical points, got %f", d)
	}
}

func TestDistanceSymmetry(t *testing.T) {
	d1 := Distance(40.71, -74.01, 40.72, -74.00)
	d2 := Distance(40.72, -74.00, 40.71, -74.01)
	if math.Abs(d1-d2) > 1e-10 {
		t.Fatalf("distance should be symmetric: d1=%f d2=%f", d1, d2)
	}
}

// TestFindBestMatchNearestDriver verifies that the closest available driver is selected
// when all drivers have the same rating.
func TestFindBestMatchNearestDriver(t *testing.T) {
	rider := Rider{ID: "r1", Lat: 0.0, Lon: 0.0}
	drivers := []Driver{
		{ID: "far", Lat: 5.0, Lon: 5.0, Available: true, Rating: 4.5},
		{ID: "near", Lat: 0.1, Lon: 0.1, Available: true, Rating: 4.5},
		{ID: "medium", Lat: 2.0, Lon: 2.0, Available: true, Rating: 4.5},
	}
	best, found := FindBestMatch(drivers, rider)
	if !found {
		t.Fatal("expected a match to be found")
	}
	if best.ID != "near" {
		t.Fatalf("expected nearest driver 'near', got %s", best.ID)
	}
}

// TestFindBestMatchExcludesUnavailable verifies that unavailable drivers are never selected.
func TestFindBestMatchExcludesUnavailable(t *testing.T) {
	rider := Rider{ID: "r1", Lat: 0.0, Lon: 0.0}
	drivers := []Driver{
		{ID: "closest_but_busy", Lat: 0.01, Lon: 0.01, Available: false, Rating: 5.0},
		{ID: "available_farther", Lat: 1.0, Lon: 1.0, Available: true, Rating: 4.0},
	}
	best, found := FindBestMatch(drivers, rider)
	if !found {
		t.Fatal("expected a match to be found")
	}
	if best.ID != "available_farther" {
		t.Fatalf("expected 'available_farther', got %s", best.ID)
	}
}

// TestFindBestMatchNoAvailableDrivers returns false when all drivers are unavailable.
func TestFindBestMatchNoAvailableDrivers(t *testing.T) {
	rider := Rider{ID: "r1", Lat: 0.0, Lon: 0.0}
	drivers := []Driver{
		{ID: "d1", Lat: 0.1, Lon: 0.1, Available: false, Rating: 4.5},
		{ID: "d2", Lat: 0.2, Lon: 0.2, Available: false, Rating: 4.8},
	}
	_, found := FindBestMatch(drivers, rider)
	if found {
		t.Fatal("expected no match when all drivers are unavailable")
	}
}

// TestFindBestMatchRatingTieBreak verifies that when two drivers are at the same position,
// the higher-rated driver wins.
func TestFindBestMatchRatingTieBreak(t *testing.T) {
	rider := Rider{ID: "r1", Lat: 0.0, Lon: 0.0}
	drivers := []Driver{
		{ID: "low_rating", Lat: 1.0, Lon: 1.0, Available: true, Rating: 3.5},
		{ID: "high_rating", Lat: 1.0, Lon: 1.0, Available: true, Rating: 5.0},
	}
	best, found := FindBestMatch(drivers, rider)
	if !found {
		t.Fatal("expected a match to be found")
	}
	if best.ID != "high_rating" {
		t.Fatalf("expected 'high_rating' to win tie-break, got %s", best.ID)
	}
}

// TestScoreMatchHigherRatingLowerScore verifies that a higher-rated driver scores better
// at the same distance.
func TestScoreMatchHigherRatingLowerScore(t *testing.T) {
	rider := Rider{ID: "r1", Lat: 0.0, Lon: 0.0}
	d1 := Driver{ID: "d1", Lat: 1.0, Lon: 1.0, Available: true, Rating: 4.0}
	d2 := Driver{ID: "d2", Lat: 1.0, Lon: 1.0, Available: true, Rating: 5.0}
	s1 := ScoreMatch(d1, rider)
	s2 := ScoreMatch(d2, rider)
	if s2 >= s1 {
		t.Fatalf("higher-rated driver should have lower score: s1(rating=4)=%f, s2(rating=5)=%f", s1, s2)
	}
}

// TestFindBestMatchEmptyDrivers returns false with empty driver slice.
func TestFindBestMatchEmptyDrivers(t *testing.T) {
	rider := Rider{ID: "r1", Lat: 0.0, Lon: 0.0}
	_, found := FindBestMatch([]Driver{}, rider)
	if found {
		t.Fatal("expected no match with empty driver list")
	}
}
