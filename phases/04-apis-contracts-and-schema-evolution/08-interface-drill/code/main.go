package main

type DrillScore struct {
	InterfaceChoice bool
	RetrySafety     bool
	QuerySafety     bool
	Compatibility   bool
}

func Score(d DrillScore) int {
	score := 0
	if d.InterfaceChoice {
		score++
	}
	if d.RetrySafety {
		score++
	}
	if d.QuerySafety {
		score++
	}
	if d.Compatibility {
		score++
	}
	return score
}

func Strong(d DrillScore) bool {
	return Score(d) >= 3
}

func main() {}
