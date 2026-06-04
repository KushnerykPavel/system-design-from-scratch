package main

type WrapUp struct {
	Risks          int
	TradeOffs      int
	Observability  bool
	RolloutMention bool
}

func CoverageScore(summary WrapUp) int {
	score := summary.Risks + summary.TradeOffs
	if summary.Observability {
		score++
	}
	if summary.RolloutMention {
		score++
	}
	return score
}

func StrongWrapUp(summary WrapUp) bool {
	return CoverageScore(summary) >= 5
}

func main() {}
