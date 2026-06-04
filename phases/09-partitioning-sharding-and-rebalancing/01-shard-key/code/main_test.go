package main

import "testing"

func TestRankCandidatesPrefersTenantWhenLocalityDominates(t *testing.T) {
	workload := Workload{
		SingleShardReadRatio:   0.9,
		SingleShardWriteRatio:  0.9,
		CrossShardQueryRatio:   0.1,
		HotTenantRisk:          0.5,
		MigrationComplexity:    0.3,
		PlacementControlNeeded: 0.4,
	}

	results := rankCandidates(workload, []Candidate{
		{
			Name:                  "tenant_id",
			LocalityScore:         0.95,
			DistributionScore:     0.45,
			QueryAlignmentScore:   0.95,
			IsolationScore:        0.95,
			MigrationFlexibility:  0.60,
			PlacementControlScore: 0.70,
		},
		{
			Name:                  "random_object_id",
			LocalityScore:         0.20,
			DistributionScore:     0.95,
			QueryAlignmentScore:   0.20,
			IsolationScore:        0.10,
			MigrationFlexibility:  0.50,
			PlacementControlScore: 0.20,
		},
	})

	if results[0].Name != "tenant_id" {
		t.Fatalf("expected tenant_id first, got %s", results[0].Name)
	}
}

func TestRankCandidatesPrefersFlexiblePlacementWhenMigrationMatters(t *testing.T) {
	workload := Workload{
		SingleShardReadRatio:   0.6,
		SingleShardWriteRatio:  0.7,
		CrossShardQueryRatio:   0.25,
		HotTenantRisk:          0.8,
		MigrationComplexity:    0.9,
		PlacementControlNeeded: 0.9,
	}

	results := rankCandidates(workload, []Candidate{
		{
			Name:                  "tenant_id",
			LocalityScore:         0.95,
			DistributionScore:     0.30,
			QueryAlignmentScore:   0.90,
			IsolationScore:        0.95,
			MigrationFlexibility:  0.40,
			PlacementControlScore: 0.50,
		},
		{
			Name:                  "regional_tenant_bucket",
			LocalityScore:         0.80,
			DistributionScore:     0.80,
			QueryAlignmentScore:   0.75,
			IsolationScore:        0.75,
			MigrationFlexibility:  0.90,
			PlacementControlScore: 0.95,
		},
	})

	if results[0].Name != "regional_tenant_bucket" {
		t.Fatalf("expected regional_tenant_bucket first, got %s", results[0].Name)
	}
}
