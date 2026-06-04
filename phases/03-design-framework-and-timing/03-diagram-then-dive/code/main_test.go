package main

import "testing"

func TestReadyForDeepDive(t *testing.T) {
	review := Review{
		HasHighLevelDiagram: true,
		HasCriticalPath:     true,
		HasDeepDive:         false,
	}
	if !ReadyForDeepDive(review) {
		t.Fatal("expected review to be ready for deep dive")
	}

	missing := MissingForDeepDive(review)
	if len(missing) != 1 || missing[0] != "deep_dive" {
		t.Fatalf("unexpected missing items: %#v", missing)
	}
}
