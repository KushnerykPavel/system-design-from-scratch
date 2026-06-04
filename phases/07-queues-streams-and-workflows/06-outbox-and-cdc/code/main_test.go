package main

import "testing"

func TestValidateRecords(t *testing.T) {
	valid := []OutboxRecord{
		{EventID: "e1", AggregateID: "a", AggregateVersion: 1},
		{EventID: "e2", AggregateID: "a", AggregateVersion: 2},
		{EventID: "e3", AggregateID: "b", AggregateVersion: 1},
	}
	if problems := ValidateRecords(valid); len(problems) != 0 {
		t.Fatalf("ValidateRecords(valid) = %v, want no problems", problems)
	}

	invalid := []OutboxRecord{
		{EventID: "e1", AggregateID: "a", AggregateVersion: 1},
		{EventID: "e1", AggregateID: "a", AggregateVersion: 1},
	}
	if problems := ValidateRecords(invalid); len(problems) != 2 {
		t.Fatalf("ValidateRecords(invalid) returned %d problems, want 2", len(problems))
	}
}

func TestUnpublished(t *testing.T) {
	records := []OutboxRecord{
		{EventID: "e1", Published: true},
		{EventID: "e2", Published: false},
		{EventID: "e3", Published: false},
	}
	got := Unpublished(records)
	if len(got) != 2 {
		t.Fatalf("len(Unpublished(records)) = %d, want 2", len(got))
	}
	if got[0].EventID != "e2" || got[1].EventID != "e3" {
		t.Fatalf("unexpected unpublished records: %+v", got)
	}
}
