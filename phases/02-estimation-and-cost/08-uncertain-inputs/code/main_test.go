package main

import "testing"

func TestSpreadRatioAndIsWide(t *testing.T) {
	r := Range{Low: 20, Base: 60, High: 180}
	if got := SpreadRatio(r); got != 9 {
		t.Fatalf("unexpected spread ratio: %.2f", got)
	}
	if !IsWide(r, 5) {
		t.Fatal("expected range to be wide")
	}
}
