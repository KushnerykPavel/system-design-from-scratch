package main

type Candidate struct {
	Name       string
	Risk       int
	Scale      int
	Novelty    int
	Dependency int
}

func Score(candidate Candidate) int {
	return candidate.Risk + candidate.Scale + candidate.Novelty + candidate.Dependency
}

func BestCandidate(candidates []Candidate) Candidate {
	best := candidates[0]
	bestScore := Score(best)
	for _, candidate := range candidates[1:] {
		if score := Score(candidate); score > bestScore {
			best = candidate
			bestScore = score
		}
	}
	return best
}

func main() {}
