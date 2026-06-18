package main

import "testing"

func TestEvaluateMockStrongHire(t *testing.T) {
	scores := make([]Score, 8)
	for i := range scores {
		scores[i] = Score{Points: 2}
	}
	result := EvaluateMock(scores)
	if result.TotalPoints != 16 {
		t.Fatalf("expected 16 total points, got %d", result.TotalPoints)
	}
	if result.Signal != "strong-hire" {
		t.Fatalf("expected strong-hire, got %s", result.Signal)
	}
}

func TestEvaluateMockNoHire(t *testing.T) {
	scores := make([]Score, 8)
	result := EvaluateMock(scores)
	if result.TotalPoints != 0 {
		t.Fatalf("expected 0 total points, got %d", result.TotalPoints)
	}
	if result.Signal != "no-hire" {
		t.Fatalf("expected no-hire, got %s", result.Signal)
	}
}

func TestEvaluateMockHireThreshold(t *testing.T) {
	// 10 points = hire
	scores := []Score{
		{Points: 2}, {Points: 2}, {Points: 2}, {Points: 2}, {Points: 2},
		{Points: 0}, {Points: 0}, {Points: 0},
	}
	result := EvaluateMock(scores)
	if result.TotalPoints != 10 {
		t.Fatalf("expected 10 points, got %d", result.TotalPoints)
	}
	if result.Signal != "hire" {
		t.Fatalf("expected hire at 10 points, got %s", result.Signal)
	}
}

func TestEvaluateMockMaxPoints(t *testing.T) {
	scores := make([]Score, 6)
	result := EvaluateMock(scores)
	if result.MaxPoints != 12 {
		t.Fatalf("expected max 12 for 6 dimensions, got %d", result.MaxPoints)
	}
}
