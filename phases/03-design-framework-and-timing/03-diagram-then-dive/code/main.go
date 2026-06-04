package main

type Review struct {
	HasHighLevelDiagram bool
	HasCriticalPath     bool
	HasDeepDive         bool
}

func ReadyForDeepDive(review Review) bool {
	return review.HasHighLevelDiagram && review.HasCriticalPath
}

func MissingForDeepDive(review Review) []string {
	var missing []string
	if !review.HasHighLevelDiagram {
		missing = append(missing, "high_level_diagram")
	}
	if !review.HasCriticalPath {
		missing = append(missing, "critical_path")
	}
	if !review.HasDeepDive {
		missing = append(missing, "deep_dive")
	}
	return missing
}

func main() {}
