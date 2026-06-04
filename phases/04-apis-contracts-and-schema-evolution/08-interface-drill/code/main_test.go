package main

import "testing"

func TestScoreAndStrong(t *testing.T) {
	d := DrillScore{
		InterfaceChoice: true,
		RetrySafety:     true,
		QuerySafety:     true,
		Compatibility:   false,
	}
	if got := Score(d); got != 3 {
		t.Fatalf("expected score 3, got %d", got)
	}
	if !Strong(d) {
		t.Fatal("expected drill to count as strong")
	}
}
