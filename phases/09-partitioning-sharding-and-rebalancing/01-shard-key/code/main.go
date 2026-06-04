package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Workload struct {
	SingleShardReadRatio   float64 `json:"single_shard_read_ratio"`
	SingleShardWriteRatio  float64 `json:"single_shard_write_ratio"`
	CrossShardQueryRatio   float64 `json:"cross_shard_query_ratio"`
	HotTenantRisk          float64 `json:"hot_tenant_risk"`
	MigrationComplexity    float64 `json:"migration_complexity"`
	PlacementControlNeeded float64 `json:"placement_control_needed"`
}

type Candidate struct {
	Name                  string  `json:"name"`
	LocalityScore         float64 `json:"locality_score"`
	DistributionScore     float64 `json:"distribution_score"`
	QueryAlignmentScore   float64 `json:"query_alignment_score"`
	IsolationScore        float64 `json:"isolation_score"`
	MigrationFlexibility  float64 `json:"migration_flexibility"`
	PlacementControlScore float64 `json:"placement_control_score"`
}

type Evaluation struct {
	Name    string  `json:"name"`
	Score   float64 `json:"score"`
	Summary string  `json:"summary"`
}

func scoreCandidate(w Workload, c Candidate) Evaluation {
	score := 0.0
	score += c.LocalityScore * (w.SingleShardReadRatio*0.20 + w.SingleShardWriteRatio*0.15)
	score += c.QueryAlignmentScore * (1 - w.CrossShardQueryRatio) * 0.18
	score += c.DistributionScore * (0.10 + w.HotTenantRisk*0.14)
	score += c.IsolationScore * (0.08 + w.HotTenantRisk*0.12)
	score += c.MigrationFlexibility * (0.05 + w.MigrationComplexity*0.18)
	score += c.PlacementControlScore * (0.05 + w.PlacementControlNeeded*0.18)

	summary := "balanced choice"
	if c.DistributionScore < 0.45 && w.HotTenantRisk > 0.6 {
		summary = "high hotspot risk under skew"
	} else if c.QueryAlignmentScore < 0.45 && w.CrossShardQueryRatio > 0.3 {
		summary = "fanout-heavy for the current query mix"
	} else if c.MigrationFlexibility > 0.75 && w.MigrationComplexity > 0.5 {
		summary = "strong option when future moves are likely"
	}

	return Evaluation{Name: c.Name, Score: score, Summary: summary}
}

func rankCandidates(w Workload, candidates []Candidate) []Evaluation {
	results := make([]Evaluation, 0, len(candidates))
	for _, candidate := range candidates {
		results = append(results, scoreCandidate(w, candidate))
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Name < results[j].Name
		}
		return results[i].Score > results[j].Score
	})

	return results
}

func main() {
	workload := Workload{
		SingleShardReadRatio:   0.8,
		SingleShardWriteRatio:  0.9,
		CrossShardQueryRatio:   0.15,
		HotTenantRisk:          0.7,
		MigrationComplexity:    0.6,
		PlacementControlNeeded: 0.6,
	}

	candidates := []Candidate{
		{
			Name:                  "tenant_id",
			LocalityScore:         0.95,
			DistributionScore:     0.45,
			QueryAlignmentScore:   0.90,
			IsolationScore:        0.95,
			MigrationFlexibility:  0.65,
			PlacementControlScore: 0.80,
		},
		{
			Name:                  "hashed_tenant_bucket",
			LocalityScore:         0.75,
			DistributionScore:     0.90,
			QueryAlignmentScore:   0.80,
			IsolationScore:        0.70,
			MigrationFlexibility:  0.85,
			PlacementControlScore: 0.85,
		},
		{
			Name:                  "random_object_id",
			LocalityScore:         0.25,
			DistributionScore:     0.95,
			QueryAlignmentScore:   0.20,
			IsolationScore:        0.15,
			MigrationFlexibility:  0.40,
			PlacementControlScore: 0.30,
		},
	}

	encoded, err := json.MarshalIndent(rankCandidates(workload, candidates), "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(encoded))
}
