package main

import (
	"fmt"
	"testing"
)

func TestAssignIsStable(t *testing.T) {
	exp := Experiment{
		ID: "stability-test",
		Treatments: []Treatment{
			{Name: "control", BucketStart: 0, BucketEnd: 5000},
			{Name: "treatment", BucketStart: 5000, BucketEnd: 10000},
		},
	}
	first := Assign("user-xyz", exp)
	// Same call must return the same treatment.
	for i := 0; i < 100; i++ {
		got := Assign("user-xyz", exp)
		if got != first {
			t.Fatalf("assignment not stable: first=%s got=%s on iteration %d", first, got, i)
		}
	}
}

func TestAssignExperimentsAreIndependent(t *testing.T) {
	// Two experiments with the same user should be independently assigned.
	expA := Experiment{
		ID:         "exp-A",
		Treatments: []Treatment{{Name: "control", BucketStart: 0, BucketEnd: 10000}},
	}
	expB := Experiment{
		ID:         "exp-B",
		Treatments: []Treatment{{Name: "control", BucketStart: 0, BucketEnd: 10000}},
	}
	// Both return control since 100% of buckets are control — just verify no panic.
	a := Assign("user-1", expA)
	b := Assign("user-1", expB)
	if a != "control" || b != "control" {
		t.Fatalf("expected both control, got a=%s b=%s", a, b)
	}
}

func TestAssignReturnsUnassignedForGap(t *testing.T) {
	exp := Experiment{
		ID: "partial-traffic",
		Treatments: []Treatment{
			{Name: "treatment", BucketStart: 0, BucketEnd: 1000}, // only 10% traffic
		},
	}
	// Force a bucket outside the treatment range by iterating users.
	unassignedCount := 0
	for i := 0; i < 100; i++ {
		userID := fmt.Sprintf("user-%d", i)
		if Assign(userID, exp) == "unassigned" {
			unassignedCount++
		}
	}
	if unassignedCount == 0 {
		t.Fatal("expected some users to be unassigned in a 10% traffic experiment")
	}
}

