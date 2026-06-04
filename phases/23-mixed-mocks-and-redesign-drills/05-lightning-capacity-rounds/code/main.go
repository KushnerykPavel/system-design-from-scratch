package main

import (
	"encoding/json"
	"flag"
	"os"
)

type CapacityRound struct {
	Name               string `json:"name"`
	HasQPS             bool   `json:"has_qps"`
	HasPeakFactor      bool   `json:"has_peak_factor"`
	HasStorageOrEgress bool   `json:"has_storage_or_egress"`
	HasAmplification   bool   `json:"has_amplification"`
	NamesBottleneck    bool   `json:"names_bottleneck"`
	LinksToDesign      bool   `json:"links_to_design"`
}

func ScoreCapacityRound(round CapacityRound) []string {
	var issues []string
	if !round.HasQPS {
		issues = append(issues, "has_qps should be true so the round has a traffic anchor")
	}
	if !round.HasPeakFactor {
		issues = append(issues, "has_peak_factor should be true so average-only math does not hide burst risk")
	}
	if !round.HasStorageOrEgress {
		issues = append(issues, "has_storage_or_egress should be true when persistence or network cost shapes the answer")
	}
	if !round.HasAmplification {
		issues = append(issues, "has_amplification should be true so fanout, replication, or retry multiplication is visible")
	}
	if !round.NamesBottleneck {
		issues = append(issues, "names_bottleneck should be true so the math reaches a useful conclusion")
	}
	if !round.LinksToDesign {
		issues = append(issues, "links_to_design should be true so the estimate changes architecture instead of staying decorative")
	}
	return issues
}

func main() {
	name := flag.String("name", "lightning-capacity-round", "round name")
	flag.Parse()

	round := CapacityRound{
		Name:               *name,
		HasQPS:             true,
		HasPeakFactor:      true,
		HasStorageOrEgress: true,
		HasAmplification:   true,
		NamesBottleneck:    true,
		LinksToDesign:      true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"round":  round,
		"issues": ScoreCapacityRound(round),
	})
}
