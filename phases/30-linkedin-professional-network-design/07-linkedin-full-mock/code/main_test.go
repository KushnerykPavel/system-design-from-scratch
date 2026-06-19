package main

import (
	"testing"
)

func TestEvaluateLinkedInMock_StrongHire(t *testing.T) {
	card := LinkedInScorecard{
		GraphDesign:          90,
		PipelineArchitecture: 90,
		RankingMLDesign:      90,
		PrivacyCompliance:    90,
		FailureRecovery:      90,
		Observability:        90,
	}

	result := EvaluateLinkedInMock(card)

	if result.HireSignal != "Strong Hire" {
		t.Errorf("expected 'Strong Hire' for all-90 scorecard, got %q (score=%d)", result.HireSignal, result.WeightedScore)
	}
	if result.WeightedScore < 85 {
		t.Errorf("expected weighted score >= 85 for all-90 scorecard, got %d", result.WeightedScore)
	}
}

func TestEvaluateLinkedInMock_HireBand(t *testing.T) {
	card := LinkedInScorecard{
		GraphDesign:          75,
		PipelineArchitecture: 75,
		RankingMLDesign:      75,
		PrivacyCompliance:    75,
		FailureRecovery:      75,
		Observability:        75,
	}

	result := EvaluateLinkedInMock(card)

	if result.HireSignal != "Hire" {
		t.Errorf("expected 'Hire' for all-75 scorecard, got %q (score=%d)", result.HireSignal, result.WeightedScore)
	}
}

func TestEvaluateLinkedInMock_NoHireYet(t *testing.T) {
	card := LinkedInScorecard{
		GraphDesign:          65,
		PipelineArchitecture: 60,
		RankingMLDesign:      60,
		PrivacyCompliance:    60,
		FailureRecovery:      60,
		Observability:        60,
	}

	result := EvaluateLinkedInMock(card)

	if result.HireSignal != "No Hire (yet)" {
		t.Errorf("expected 'No Hire (yet)' band, got %q (score=%d)", result.HireSignal, result.WeightedScore)
	}
}

func TestEvaluateLinkedInMock_NoHire(t *testing.T) {
	card := LinkedInScorecard{
		GraphDesign:          40,
		PipelineArchitecture: 40,
		RankingMLDesign:      40,
		PrivacyCompliance:    40,
		FailureRecovery:      40,
		Observability:        40,
	}

	result := EvaluateLinkedInMock(card)

	if result.HireSignal != "No Hire" {
		t.Errorf("expected 'No Hire' for all-40 scorecard, got %q (score=%d)", result.HireSignal, result.WeightedScore)
	}
	if result.WeightedScore >= 55 {
		t.Errorf("expected weighted score < 55 for all-40 scorecard, got %d", result.WeightedScore)
	}
}

func TestEvaluateLinkedInMock_WeightedScoreIsCorrect(t *testing.T) {
	// All dimensions score 100 → weighted score must be 100
	card := LinkedInScorecard{
		GraphDesign:          100,
		PipelineArchitecture: 100,
		RankingMLDesign:      100,
		PrivacyCompliance:    100,
		FailureRecovery:      100,
		Observability:        100,
	}

	result := EvaluateLinkedInMock(card)
	if result.WeightedScore != 100 {
		t.Errorf("expected weighted score of 100 for perfect scorecard, got %d", result.WeightedScore)
	}
}

func TestEvaluateLinkedInMock_WeakAreaDetected(t *testing.T) {
	// PipelineArchitecture is a weak area (score < 60)
	card := LinkedInScorecard{
		GraphDesign:          85,
		PipelineArchitecture: 30, // weak
		RankingMLDesign:      80,
		PrivacyCompliance:    80,
		FailureRecovery:      75,
		Observability:        75,
	}

	result := EvaluateLinkedInMock(card)

	found := false
	for _, area := range result.WeakAreas {
		if area == "Pipeline Architecture" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'Pipeline Architecture' in weak areas, got %v", result.WeakAreas)
	}
}

func TestEvaluateLinkedInMock_NoWeakAreasWhenAllStrong(t *testing.T) {
	card := LinkedInScorecard{
		GraphDesign:          80,
		PipelineArchitecture: 80,
		RankingMLDesign:      75,
		PrivacyCompliance:    75,
		FailureRecovery:      70,
		Observability:        70,
	}

	result := EvaluateLinkedInMock(card)

	if len(result.WeakAreas) != 0 {
		t.Errorf("expected no weak areas for all-70+ scorecard, got %v", result.WeakAreas)
	}
}

func TestEvaluateLinkedInMock_GraphDesignAndPipelineCarryMoreWeight(t *testing.T) {
	// High Graph+Pipeline with low others
	highCore := LinkedInScorecard{
		GraphDesign:          100,
		PipelineArchitecture: 100,
		RankingMLDesign:      0,
		PrivacyCompliance:    0,
		FailureRecovery:      0,
		Observability:        0,
	}
	// Low Graph+Pipeline with high others
	highOthers := LinkedInScorecard{
		GraphDesign:          0,
		PipelineArchitecture: 0,
		RankingMLDesign:      100,
		PrivacyCompliance:    100,
		FailureRecovery:      100,
		Observability:        100,
	}

	coreResult := EvaluateLinkedInMock(highCore)
	othersResult := EvaluateLinkedInMock(highOthers)

	// Graph (20) + Pipeline (20) = 40pts; ML (15) + Privacy (15) + Failure (15) + Obs (15) = 60pts
	// highCore should score 40, highOthers should score 60
	if coreResult.WeightedScore != 40 {
		t.Errorf("expected core-only score of 40, got %d", coreResult.WeightedScore)
	}
	if othersResult.WeightedScore != 60 {
		t.Errorf("expected others-only score of 60, got %d", othersResult.WeightedScore)
	}
	if coreResult.WeightedScore >= othersResult.WeightedScore {
		t.Errorf("others dimensions (60pts) should outweigh core dimensions (40pts)")
	}
}
