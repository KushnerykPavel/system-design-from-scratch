package main

import (
	"fmt"
	"math"
	"sort"
)

// Driver represents an available Uber driver.
type Driver struct {
	ID        string
	Lat       float64
	Lon       float64
	Available bool
	Rating    float64 // 1.0–5.0
}

// Rider represents a trip request from a rider.
type Rider struct {
	ID  string
	Lat float64
	Lon float64
}

// Distance computes Euclidean distance between two lat/lon points (simplified —
// not Haversine; suitable for small areas or relative comparisons).
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	dlat := lat1 - lat2
	dlon := lon1 - lon2
	return math.Sqrt(dlat*dlat + dlon*dlon)
}

// ScoreMatch computes a match score for a driver-rider pair.
// Lower score is better. Distance is the primary factor; rating adjusts it so
// that a higher-rated driver is preferred when distances are similar.
//
// Formula: score = distance / rating_normalized
// rating_normalized = rating / 5.0, so a 5.0-rated driver gets full weight
// and a 1.0-rated driver gets 5x the effective distance penalty.
func ScoreMatch(driver Driver, rider Rider) float64 {
	dist := Distance(driver.Lat, driver.Lon, rider.Lat, rider.Lon)
	ratingNorm := driver.Rating / 5.0
	if ratingNorm <= 0 {
		ratingNorm = 0.01 // guard against division by zero
	}
	return dist / ratingNorm
}

// MatchCandidate pairs a driver with their computed match score.
type MatchCandidate struct {
	Driver Driver
	Score  float64
}

// FindBestMatch returns the best available driver for a rider and a boolean
// indicating whether a match was found. Drivers with Available=false are excluded.
func FindBestMatch(drivers []Driver, rider Rider) (Driver, bool) {
	var candidates []MatchCandidate
	for _, d := range drivers {
		if !d.Available {
			continue
		}
		score := ScoreMatch(d, rider)
		candidates = append(candidates, MatchCandidate{Driver: d, Score: score})
	}
	if len(candidates) == 0 {
		return Driver{}, false
	}
	// Sort by score ascending (lower = better).
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Score != candidates[j].Score {
			return candidates[i].Score < candidates[j].Score
		}
		// Tie-break by rating descending.
		return candidates[i].Driver.Rating > candidates[j].Driver.Rating
	})
	return candidates[0].Driver, true
}

func main() {
	drivers := []Driver{
		{ID: "driver_A", Lat: 40.71, Lon: -74.01, Available: true, Rating: 4.8},
		{ID: "driver_B", Lat: 40.72, Lon: -74.00, Available: true, Rating: 4.2},
		{ID: "driver_C", Lat: 40.70, Lon: -74.02, Available: false, Rating: 4.9}, // unavailable
		{ID: "driver_D", Lat: 40.75, Lon: -73.98, Available: true, Rating: 3.8},
		{ID: "driver_E", Lat: 40.71, Lon: -74.01, Available: true, Rating: 5.0}, // same position as A, higher rating
	}

	rider := Rider{ID: "rider_1", Lat: 40.715, Lon: -74.005}

	fmt.Printf("Rider %s at (%.3f, %.3f)\n\n", rider.ID, rider.Lat, rider.Lon)

	for _, d := range drivers {
		dist := Distance(d.Lat, d.Lon, rider.Lat, rider.Lon)
		score := 0.0
		status := "UNAVAILABLE"
		if d.Available {
			score = ScoreMatch(d, rider)
			status = "available"
		}
		fmt.Printf("  %-10s dist=%.4f rating=%.1f score=%.4f [%s]\n",
			d.ID, dist, d.Rating, score, status)
	}

	best, found := FindBestMatch(drivers, rider)
	if !found {
		fmt.Println("\nNo available drivers found.")
		return
	}
	fmt.Printf("\nBest match: %s (rating=%.1f, score=%.4f)\n",
		best.ID, best.Rating, ScoreMatch(best, rider))
}
