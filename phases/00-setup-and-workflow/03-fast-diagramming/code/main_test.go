package main

import "testing"

func TestValidateFirstPassAcceptsReasonableDiagram(t *testing.T) {
	t.Parallel()

	diagram := Diagram{
		Components: []Component{
			{Name: "client"},
			{Name: "gateway", Critical: true},
			{Name: "service"},
			{Name: "db"},
		},
		Edges: []Edge{
			{From: "client", To: "gateway"},
			{From: "gateway", To: "service"},
			{From: "service", To: "db"},
		},
	}

	if issues := ValidateFirstPass(diagram); len(issues) != 0 {
		t.Fatalf("ValidateFirstPass() returned issues: %v", issues)
	}
}

func TestValidateFirstPassRejectsOverdrawnDiagram(t *testing.T) {
	t.Parallel()

	diagram := Diagram{
		Components: []Component{
			{Name: "a"}, {Name: "b"}, {Name: "c"}, {Name: "d"},
			{Name: "e"}, {Name: "f"}, {Name: "g"}, {Name: "h"}, {Name: "i"},
		},
	}

	if issues := ValidateFirstPass(diagram); len(issues) == 0 {
		t.Fatal("ValidateFirstPass() returned no issues for an over-budget diagram")
	}
}

func TestComplexityScore(t *testing.T) {
	t.Parallel()

	diagram := Diagram{
		Components: []Component{{Name: "a"}, {Name: "b"}},
		Edges:      []Edge{{From: "a", To: "b"}},
	}

	if got := ComplexityScore(diagram); got != 5 {
		t.Fatalf("ComplexityScore() = %d, want 5", got)
	}
}
