package main

import "fmt"

type Cohort struct {
	ID               string
	BackfillDone     bool
	DualWriteReady   bool
	ParityMismatch   int
	ClientCompatible bool
}

type Step struct {
	CohortID string
	Action   string
}

func nextSteps(cohorts []Cohort) []Step {
	var steps []Step
	for _, cohort := range cohorts {
		switch {
		case !cohort.BackfillDone:
			steps = append(steps, Step{CohortID: cohort.ID, Action: "start_backfill"})
		case !cohort.DualWriteReady:
			steps = append(steps, Step{CohortID: cohort.ID, Action: "enable_dual_write"})
		case !cohort.ClientCompatible:
			steps = append(steps, Step{CohortID: cohort.ID, Action: "hold_for_client_compat"})
		case cohort.ParityMismatch > 0:
			steps = append(steps, Step{CohortID: cohort.ID, Action: "investigate_parity"})
		default:
			steps = append(steps, Step{CohortID: cohort.ID, Action: "cutover"})
		}
	}
	return steps
}

func main() {
	cohorts := []Cohort{
		{ID: "c1", BackfillDone: false},
		{ID: "c2", BackfillDone: true, DualWriteReady: true, ClientCompatible: true},
	}

	for _, step := range nextSteps(cohorts) {
		fmt.Printf("%s -> %s\n", step.CohortID, step.Action)
	}
}
