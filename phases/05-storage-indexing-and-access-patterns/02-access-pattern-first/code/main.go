package main

import "sort"

type AccessPattern struct {
	Name              string
	ReadQPS           int
	WriteQPS          int
	LatencyCritical   bool
	RequiresStrongTxn bool
	SupportsRevenue   bool
}

type RankedPattern struct {
	AccessPattern
	Score int
}

func RankPatterns(patterns []AccessPattern) []RankedPattern {
	ranked := make([]RankedPattern, 0, len(patterns))
	for _, pattern := range patterns {
		score := pattern.ReadQPS*3 + pattern.WriteQPS*2
		if pattern.LatencyCritical {
			score += 1200
		}
		if pattern.RequiresStrongTxn {
			score += 900
		}
		if pattern.SupportsRevenue {
			score += 700
		}
		ranked = append(ranked, RankedPattern{
			AccessPattern: pattern,
			Score:         score,
		})
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].Score == ranked[j].Score {
			return ranked[i].Name < ranked[j].Name
		}
		return ranked[i].Score > ranked[j].Score
	})

	return ranked
}

func main() {}
