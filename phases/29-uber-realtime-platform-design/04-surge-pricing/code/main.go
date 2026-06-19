package main

import "fmt"

// Zone represents a geographic surge zone with current supply and demand counts.
type Zone struct {
	ID               string
	AvailableDrivers int
	OpenRequests     int
}

// PricingTier maps a minimum demand/supply ratio to a surge multiplier.
type PricingTier struct {
	MinRatio   float64
	Multiplier float64
}

// pricingTiers defines the surge multiplier step function in ascending ratio order.
// The last tier acts as the cap.
var pricingTiers = []PricingTier{
	{MinRatio: 4.0, Multiplier: 8.0},
	{MinRatio: 2.0, Multiplier: 4.0},
	{MinRatio: 1.5, Multiplier: 2.0},
	{MinRatio: 1.25, Multiplier: 1.5},
	{MinRatio: 1.0, Multiplier: 1.2},
	{MinRatio: 0.0, Multiplier: 1.0},
}

// SurgeMultiplier computes the surge multiplier for a zone based on the
// demand/supply ratio. Supply is floored at 1 to prevent division by zero.
func SurgeMultiplier(zone Zone) float64 {
	supply := zone.AvailableDrivers
	if supply < 1 {
		supply = 1
	}
	demand := zone.OpenRequests
	if demand < 0 {
		demand = 0
	}
	ratio := float64(demand) / float64(supply)

	for _, tier := range pricingTiers {
		if ratio >= tier.MinRatio {
			return tier.Multiplier
		}
	}
	return 1.0
}

// RequiresConfirmation returns true when the multiplier is high enough that
// explicit rider acknowledgement is required before booking.
// Uber requires confirmation at 2.0× and above.
func RequiresConfirmation(multiplier float64) bool {
	return multiplier >= 2.0
}

func main() {
	scenarios := []Zone{
		{ID: "zone_downtown", AvailableDrivers: 20, OpenRequests: 10},
		{ID: "zone_airport", AvailableDrivers: 4, OpenRequests: 5},
		{ID: "zone_stadium", AvailableDrivers: 2, OpenRequests: 5},
		{ID: "zone_concert", AvailableDrivers: 1, OpenRequests: 8},
		{ID: "zone_disaster", AvailableDrivers: 0, OpenRequests: 15},
		{ID: "zone_quiet", AvailableDrivers: 30, OpenRequests: 3},
	}

	fmt.Printf("%-20s  %8s  %8s  %8s  %s\n",
		"Zone", "Drivers", "Requests", "Multiplier", "Confirmation")
	fmt.Println("---------------------------------------------------------------")

	for _, z := range scenarios {
		m := SurgeMultiplier(z)
		confirm := ""
		if RequiresConfirmation(m) {
			confirm = "REQUIRED"
		}
		fmt.Printf("%-20s  %8d  %8d  %8.1fx  %s\n",
			z.ID, z.AvailableDrivers, z.OpenRequests, m, confirm)
	}
}
