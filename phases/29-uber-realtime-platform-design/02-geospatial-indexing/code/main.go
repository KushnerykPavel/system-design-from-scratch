package main

import (
	"fmt"
	"strconv"
	"strings"
)

// LatLon represents a geographic coordinate.
type LatLon struct {
	Lat float64
	Lon float64
}

// CellID maps a lat/lon to a simplified grid cell ID of the form "R{row}C{col}".
// The 100x100 grid covers lat [-90, 90] and lon [-180, 180].
// resolution is unused in this simplified model but mirrors the H3 API signature.
func CellID(lat, lon float64, resolution int) string {
	// Clamp to valid ranges.
	if lat < -90 {
		lat = -90
	}
	if lat > 90 {
		lat = 90
	}
	if lon < -180 {
		lon = -180
	}
	if lon > 180 {
		lon = 180
	}
	row := int((lat + 90) / 1.8)   // 180 degrees / 100 rows = 1.8 deg per row
	col := int((lon + 180) / 3.6)  // 360 degrees / 100 cols = 3.6 deg per col
	if row >= 100 {
		row = 99
	}
	if col >= 100 {
		col = 99
	}
	return fmt.Sprintf("R%dC%d", row, col)
}

// parseCell extracts row and col from a cell ID string "R{row}C{col}".
func parseCell(cellID string) (int, int, bool) {
	if !strings.HasPrefix(cellID, "R") {
		return 0, 0, false
	}
	cIdx := strings.Index(cellID, "C")
	if cIdx < 0 {
		return 0, 0, false
	}
	row, err1 := strconv.Atoi(cellID[1:cIdx])
	col, err2 := strconv.Atoi(cellID[cIdx+1:])
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return row, col, true
}

// KRing returns all cell IDs within k steps of the center cell (including center).
// For a grid, k-ring(k) returns (2k+1)^2 cells (clipped at grid boundaries).
func KRing(cellID string, k int) []string {
	row, col, ok := parseCell(cellID)
	if !ok {
		return nil
	}
	var result []string
	for dr := -k; dr <= k; dr++ {
		for dc := -k; dc <= k; dc++ {
			r := row + dr
			c := col + dc
			if r < 0 || r >= 100 || c < 0 || c >= 100 {
				continue
			}
			result = append(result, fmt.Sprintf("R%dC%d", r, c))
		}
	}
	return result
}

// DriverStore maps H3 cell ID to a set of driver IDs.
type DriverStore map[string]map[string]bool

// NewDriverStore creates an empty driver store.
func NewDriverStore() DriverStore {
	return make(DriverStore)
}

// UpdateDriverLocation moves a driver from oldCell to newCell.
// Pass empty string for oldCell if the driver has no previous cell.
func UpdateDriverLocation(store DriverStore, driverID, oldCell, newCell string) {
	if oldCell != "" {
		if drivers, ok := store[oldCell]; ok {
			delete(drivers, driverID)
			if len(drivers) == 0 {
				delete(store, oldCell)
			}
		}
	}
	if _, ok := store[newCell]; !ok {
		store[newCell] = make(map[string]bool)
	}
	store[newCell][driverID] = true
}

// FindNearbyDrivers returns all driver IDs within k steps of riderCell.
func FindNearbyDrivers(store DriverStore, riderCell string, k int) []string {
	cells := KRing(riderCell, k)
	seen := make(map[string]bool)
	var result []string
	for _, cell := range cells {
		for driverID := range store[cell] {
			if !seen[driverID] {
				seen[driverID] = true
				result = append(result, driverID)
			}
		}
	}
	return result
}

func main() {
	store := NewDriverStore()

	// Place five drivers at various locations.
	drivers := []struct {
		id  string
		lat float64
		lon float64
	}{
		{"driver_1", 40.71, -74.01},  // Manhattan
		{"driver_2", 40.72, -74.00},  // nearby
		{"driver_3", 40.73, -73.99},  // same neighborhood
		{"driver_4", 51.50, -0.12},   // London — far away
		{"driver_5", 40.70, -74.02},  // Manhattan
	}

	cells := make(map[string]string)
	for _, d := range drivers {
		cell := CellID(d.lat, d.lon, 9)
		cells[d.id] = cell
		UpdateDriverLocation(store, d.id, "", cell)
		fmt.Printf("Driver %s → cell %s\n", d.id, cell)
	}

	// Simulate a rider requesting a trip from Manhattan.
	riderLat, riderLon := 40.715, -74.005
	riderCell := CellID(riderLat, riderLon, 9)
	fmt.Printf("\nRider cell: %s\n", riderCell)

	nearby := FindNearbyDrivers(store, riderCell, 1)
	fmt.Printf("Drivers within k-ring(1): %v\n", nearby)

	// Simulate driver_1 moving to a new cell.
	newCell := CellID(40.80, -73.95, 9)
	oldCell := cells["driver_1"]
	UpdateDriverLocation(store, "driver_1", oldCell, newCell)
	fmt.Printf("\ndriver_1 moved from %s to %s\n", oldCell, newCell)

	nearby2 := FindNearbyDrivers(store, riderCell, 1)
	fmt.Printf("Drivers within k-ring(1) after driver_1 moved: %v\n", nearby2)
}
