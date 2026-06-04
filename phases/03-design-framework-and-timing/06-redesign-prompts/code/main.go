package main

type ConstraintChange struct {
	Name     string
	Severity int
}

func RedesignPressure(changes []ConstraintChange) int {
	total := 0
	for _, change := range changes {
		total += change.Severity
	}
	return total
}

func RequiresTopologyChange(changes []ConstraintChange) bool {
	return RedesignPressure(changes) >= 7
}

func main() {}
