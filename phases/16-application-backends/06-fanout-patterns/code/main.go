package main

import (
	"encoding/json"
	"flag"
	"os"
)

type WorkloadShape struct {
	Name                string `json:"name"`
	ReadsPerWrite       int    `json:"reads_per_write"`
	RecipientsPerEvent  int    `json:"recipients_per_event"`
	FreshnessSeconds    int    `json:"freshness_seconds"`
	SkewedAudience      bool   `json:"skewed_audience"`
	SupportsPartialRead bool   `json:"supports_partial_read"`
}

func RecommendFanout(shape WorkloadShape) string {
	if shape.SkewedAudience && shape.RecipientsPerEvent > 10000 {
		return "mixed"
	}
	if shape.RecipientsPerEvent <= 20 {
		return "push"
	}
	if shape.ReadsPerWrite < 1 && shape.RecipientsPerEvent > 1000 {
		return "pull"
	}
	if shape.FreshnessSeconds <= 5 && shape.ReadsPerWrite > 10 {
		return "push"
	}
	return "mixed"
}

func ExplainFanout(shape WorkloadShape) []string {
	var notes []string
	mode := RecommendFanout(shape)
	if mode == "push" && shape.RecipientsPerEvent > 1000 {
		notes = append(notes, "push fanout may create heavy write amplification for this audience size")
	}
	if mode == "pull" && shape.ReadsPerWrite > 10 {
		notes = append(notes, "pull fanout may make the read path too expensive for the read-to-write ratio")
	}
	if mode == "mixed" && !shape.SupportsPartialRead {
		notes = append(notes, "mixed fanout needs a clear correctness contract if partial reads are not acceptable")
	}
	return notes
}

func main() {
	name := flag.String("name", "generic-workload", "workload name")
	flag.Parse()

	shape := WorkloadShape{
		Name:                *name,
		ReadsPerWrite:       50,
		RecipientsPerEvent:  100,
		FreshnessSeconds:    10,
		SkewedAudience:      true,
		SupportsPartialRead: true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"workload": shape,
		"fanout":   RecommendFanout(shape),
		"notes":    ExplainFanout(shape),
	})
}
