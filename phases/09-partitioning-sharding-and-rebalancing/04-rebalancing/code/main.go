package main

import "fmt"

type Move struct {
	ID          string
	SizeGB      int
	BandwidthGb int
	Risk        int
}

type Plan struct {
	Batches [][]Move
}

func planMoves(moves []Move, maxConcurrent int, maxBandwidthGb int, maxRisk int) Plan {
	var plan Plan
	var current []Move
	currentBandwidth := 0
	currentRisk := 0

	flush := func() {
		if len(current) == 0 {
			return
		}
		batch := make([]Move, len(current))
		copy(batch, current)
		plan.Batches = append(plan.Batches, batch)
		current = nil
		currentBandwidth = 0
		currentRisk = 0
	}

	for _, move := range moves {
		if len(current) >= maxConcurrent || currentBandwidth+move.BandwidthGb > maxBandwidthGb || currentRisk+move.Risk > maxRisk {
			flush()
		}
		current = append(current, move)
		currentBandwidth += move.BandwidthGb
		currentRisk += move.Risk
	}
	flush()

	return plan
}

func totalSize(batch []Move) int {
	total := 0
	for _, move := range batch {
		total += move.SizeGB
	}
	return total
}

func main() {
	moves := []Move{
		{ID: "r1", SizeGB: 800, BandwidthGb: 2, Risk: 2},
		{ID: "r2", SizeGB: 600, BandwidthGb: 2, Risk: 2},
		{ID: "r3", SizeGB: 900, BandwidthGb: 3, Risk: 3},
		{ID: "r4", SizeGB: 300, BandwidthGb: 1, Risk: 1},
	}

	plan := planMoves(moves, 2, 4, 4)
	for i, batch := range plan.Batches {
		fmt.Printf("batch %d: moves=%d size_gb=%d\n", i+1, len(batch), totalSize(batch))
	}
}
