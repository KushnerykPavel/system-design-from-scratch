package main

import "testing"

func TestAssignEvenly(t *testing.T) {
	got := AssignEvenly(5, []string{"c", "a"})
	if len(got) != 2 {
		t.Fatalf("len(assignments) = %d, want 2", len(got))
	}
	if got[0].Member != "a" || got[1].Member != "c" {
		t.Fatalf("members not sorted: %+v", got)
	}
	if len(got[0].Partitions) != 3 || len(got[1].Partitions) != 2 {
		t.Fatalf("unexpected partition distribution: %+v", got)
	}
}

func TestMaxLoad(t *testing.T) {
	assignments := []Assignment{
		{Member: "a", Partitions: []int{0, 1}},
		{Member: "b", Partitions: []int{2}},
	}
	if got := MaxLoad(assignments); got != 2 {
		t.Fatalf("MaxLoad() = %d, want 2", got)
	}
}

func TestNeedsMorePartitions(t *testing.T) {
	if !NeedsMorePartitions(3, 5) {
		t.Fatal("NeedsMorePartitions(3, 5) = false, want true")
	}
	if NeedsMorePartitions(5, 3) {
		t.Fatal("NeedsMorePartitions(5, 3) = true, want false")
	}
}
