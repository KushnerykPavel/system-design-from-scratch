package main

type AnswerCoverage struct {
	Clarify       bool
	Sizing        bool
	HighLevel     bool
	DeepDive      bool
	FailureModes  bool
	Observability bool
	TradeOffs     bool
	Redesign      bool
}

func CoveredCount(answer AnswerCoverage) int {
	count := 0
	if answer.Clarify {
		count++
	}
	if answer.Sizing {
		count++
	}
	if answer.HighLevel {
		count++
	}
	if answer.DeepDive {
		count++
	}
	if answer.FailureModes {
		count++
	}
	if answer.Observability {
		count++
	}
	if answer.TradeOffs {
		count++
	}
	if answer.Redesign {
		count++
	}
	return count
}

func IsFullLoop(answer AnswerCoverage) bool {
	return CoveredCount(answer) == 8
}

func main() {}
