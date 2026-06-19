package main

import "fmt"

// LinkedInScorecard holds scores (0–100) for each rubric dimension of a
// LinkedIn PYMK mock interview.
type LinkedInScorecard struct {
	GraphDesign          int // 20 points weight
	PipelineArchitecture int // 20 points weight
	RankingMLDesign      int // 15 points weight
	PrivacyCompliance    int // 15 points weight
	FailureRecovery      int // 15 points weight
	Observability        int // 15 points weight
}

// MockResult summarizes the outcome of evaluating a scorecard.
type MockResult struct {
	WeightedScore int      // 0–100, overall weighted score
	HireSignal    string   // "Strong Hire", "Hire", "No Hire (yet)", "No Hire"
	WeakAreas     []string // dimensions scoring below 60 (i.e., contributing less than 60% of their weight)
}

// dimensionWeights maps each scorecard field name to its percentage weight.
// Total weight = 100.
var dimensionWeights = []struct {
	name   string
	score  func(LinkedInScorecard) int
	weight int
}{
	{"Graph Design", func(c LinkedInScorecard) int { return c.GraphDesign }, 20},
	{"Pipeline Architecture", func(c LinkedInScorecard) int { return c.PipelineArchitecture }, 20},
	{"Ranking/ML Design", func(c LinkedInScorecard) int { return c.RankingMLDesign }, 15},
	{"Privacy/Compliance", func(c LinkedInScorecard) int { return c.PrivacyCompliance }, 15},
	{"Failure Recovery", func(c LinkedInScorecard) int { return c.FailureRecovery }, 15},
	{"Observability", func(c LinkedInScorecard) int { return c.Observability }, 15},
}

// EvaluateLinkedInMock computes the weighted score and hire signal for a
// completed LinkedIn mock interview scorecard.
//
// Weighted score formula: sum(score_i * weight_i / 100) for each dimension.
// A dimension is flagged as a weak area if its raw score < 60.
//
// Hire thresholds:
//   - 85+:   Strong Hire
//   - 70–84: Hire
//   - 55–69: No Hire (yet)
//   - <55:   No Hire
func EvaluateLinkedInMock(card LinkedInScorecard) MockResult {
	total := 0
	var weak []string

	for _, d := range dimensionWeights {
		rawScore := d.score(card)
		if rawScore < 0 {
			rawScore = 0
		}
		if rawScore > 100 {
			rawScore = 100
		}
		// Each dimension contributes (rawScore * weight / 100) to the total.
		total += rawScore * d.weight / 100
		if rawScore < 60 {
			weak = append(weak, d.name)
		}
	}

	signal := hireSignal(total)

	return MockResult{
		WeightedScore: total,
		HireSignal:    signal,
		WeakAreas:     weak,
	}
}

func hireSignal(score int) string {
	switch {
	case score >= 85:
		return "Strong Hire"
	case score >= 70:
		return "Hire"
	case score >= 55:
		return "No Hire (yet)"
	default:
		return "No Hire"
	}
}

func main() {
	candidates := []struct {
		name string
		card LinkedInScorecard
	}{
		{
			name: "Staff-level candidate",
			card: LinkedInScorecard{
				GraphDesign:          92,
				PipelineArchitecture: 88,
				RankingMLDesign:      85,
				PrivacyCompliance:    90,
				FailureRecovery:      80,
				Observability:        85,
			},
		},
		{
			name: "Solid mid-level candidate",
			card: LinkedInScorecard{
				GraphDesign:          75,
				PipelineArchitecture: 70,
				RankingMLDesign:      65,
				PrivacyCompliance:    60,
				FailureRecovery:      55,
				Observability:        50,
			},
		},
		{
			name: "Candidate who missed pipeline design",
			card: LinkedInScorecard{
				GraphDesign:          80,
				PipelineArchitecture: 30, // weak — only designed real-time traversal
				RankingMLDesign:      70,
				PrivacyCompliance:    75,
				FailureRecovery:      65,
				Observability:        60,
			},
		},
		{
			name: "Junior candidate — gaps across the board",
			card: LinkedInScorecard{
				GraphDesign:          45,
				PipelineArchitecture: 40,
				RankingMLDesign:      35,
				PrivacyCompliance:    30,
				FailureRecovery:      25,
				Observability:        20,
			},
		},
	}

	for _, c := range candidates {
		result := EvaluateLinkedInMock(c.card)
		fmt.Printf("Candidate: %s\n", c.name)
		fmt.Printf("  Weighted Score: %d/100\n", result.WeightedScore)
		fmt.Printf("  Hire Signal:    %s\n", result.HireSignal)
		if len(result.WeakAreas) > 0 {
			fmt.Printf("  Weak Areas:     %v\n", result.WeakAreas)
		} else {
			fmt.Printf("  Weak Areas:     none\n")
		}
		fmt.Println()
	}
}
