package main

import "testing"

func TestBestCandidate(t *testing.T) {
	candidates := []Candidate{
		{Name: "api_gateway", Risk: 2, Scale: 2, Novelty: 1, Dependency: 2},
		{Name: "fanout_pipeline", Risk: 4, Scale: 4, Novelty: 3, Dependency: 4},
		{Name: "admin_panel", Risk: 1, Scale: 1, Novelty: 1, Dependency: 1},
	}

	best := BestCandidate(candidates)
	if best.Name != "fanout_pipeline" {
		t.Fatalf("expected fanout_pipeline, got %s", best.Name)
	}
}
