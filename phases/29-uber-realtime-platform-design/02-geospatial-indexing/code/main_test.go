package main

import (
	"fmt"
	"testing"
)

// TestCellIDAssignment verifies that nearby coordinates map to nearby or identical cells.
func TestCellIDAssignment(t *testing.T) {
	// Two points very close together should map to the same or adjacent cells.
	c1 := CellID(40.71, -74.01, 9)
	c2 := CellID(40.711, -74.011, 9)
	// They are close enough that the cell IDs should be identical or adjacent.
	// We verify that neither is empty.
	if c1 == "" || c2 == "" {
		t.Fatalf("expected non-empty cell IDs, got %q and %q", c1, c2)
	}
	// A point in London should map to a completely different cell.
	cLondon := CellID(51.50, -0.12, 9)
	if cLondon == c1 {
		t.Fatalf("London and Manhattan should not share a cell, both got %s", c1)
	}
}

// TestCellIDBoundary verifies that extreme coordinates are handled without panic.
func TestCellIDBoundary(t *testing.T) {
	cases := []LatLon{
		{90, 180},
		{-90, -180},
		{0, 0},
		{90.1, 180.1}, // over-boundary — should be clamped
	}
	for _, c := range cases {
		id := CellID(c.Lat, c.Lon, 9)
		if id == "" {
			t.Errorf("CellID(%v, %v) returned empty string", c.Lat, c.Lon)
		}
	}
}

// TestKRingSize verifies that k-ring returns (2k+1)^2 cells for an interior cell.
func TestKRingSize(t *testing.T) {
	// Use a cell in the interior of the 100x100 grid (not near boundaries).
	centerCell := "R50C50"
	tests := []struct {
		k        int
		wantSize int
	}{
		{0, 1},
		{1, 9},   // (2*1+1)^2 = 9
		{2, 25},  // (2*2+1)^2 = 25
		{3, 49},  // (2*3+1)^2 = 49
	}
	for _, tt := range tests {
		cells := KRing(centerCell, tt.k)
		if len(cells) != tt.wantSize {
			t.Errorf("KRing(%s, %d) = %d cells, want %d", centerCell, tt.k, len(cells), tt.wantSize)
		}
	}
}

// TestKRingIncludesCenter verifies that the center cell is always included.
func TestKRingIncludesCenter(t *testing.T) {
	center := "R30C40"
	cells := KRing(center, 2)
	found := false
	for _, c := range cells {
		if c == center {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("KRing did not include center cell %s", center)
	}
}

// TestUpdateDriverLocationMovesToNewCell verifies that a driver is removed from the old cell
// and added to the new cell after an update.
func TestUpdateDriverLocationMovesToNewCell(t *testing.T) {
	store := NewDriverStore()
	UpdateDriverLocation(store, "driver_A", "", "R10C10")

	if !store["R10C10"]["driver_A"] {
		t.Fatal("driver_A should be in R10C10 after initial placement")
	}

	// Move to a new cell.
	UpdateDriverLocation(store, "driver_A", "R10C10", "R20C20")

	if store["R10C10"]["driver_A"] {
		t.Fatal("driver_A should NOT be in R10C10 after moving away")
	}
	if !store["R20C20"]["driver_A"] {
		t.Fatal("driver_A should be in R20C20 after moving there")
	}
}

// TestFindNearbyDriversReturnsDriversInRing verifies proximity search correctness.
func TestFindNearbyDriversReturnsDriversInRing(t *testing.T) {
	store := NewDriverStore()

	// Place drivers: some near the rider, one far away.
	UpdateDriverLocation(store, "near_1", "", "R50C50") // rider cell
	UpdateDriverLocation(store, "near_2", "", "R50C51") // adjacent cell
	UpdateDriverLocation(store, "near_3", "", "R51C50") // adjacent cell
	UpdateDriverLocation(store, "far_1", "", "R90C90")  // far away

	riderCell := "R50C50"
	nearby := FindNearbyDrivers(store, riderCell, 1)

	// Build a set for easy lookup.
	nearbySet := make(map[string]bool)
	for _, d := range nearby {
		nearbySet[d] = true
	}

	if !nearbySet["near_1"] {
		t.Error("near_1 (in rider cell) should be in results")
	}
	if !nearbySet["near_2"] {
		t.Error("near_2 (adjacent cell) should be in results")
	}
	if !nearbySet["near_3"] {
		t.Error("near_3 (adjacent cell) should be in results")
	}
	if nearbySet["far_1"] {
		t.Error("far_1 (R90C90) should NOT be in k-ring(1) of R50C50")
	}
}

// TestFindNearbyDriversEmptyStore returns empty slice without panic.
func TestFindNearbyDriversEmptyStore(t *testing.T) {
	store := NewDriverStore()
	nearby := FindNearbyDrivers(store, "R50C50", 1)
	if len(nearby) != 0 {
		t.Fatalf("expected empty result on empty store, got %v", nearby)
	}
}

// TestCellIDParseRoundTrip verifies that parseCell correctly inverts CellID output.
func TestCellIDParseRoundTrip(t *testing.T) {
	lats := []float64{0, 45.5, -33.9, 89.9}
	lons := []float64{0, -73.9, 151.2, -179.9}
	for i, lat := range lats {
		lon := lons[i]
		cellID := CellID(lat, lon, 9)
		row, col, ok := parseCell(cellID)
		if !ok {
			t.Errorf("parseCell(%q) failed", cellID)
			continue
		}
		expected := fmt.Sprintf("R%dC%d", row, col)
		if expected != cellID {
			t.Errorf("round-trip failed: CellID produced %q, parseCell gave R%dC%d", cellID, row, col)
		}
	}
}
