package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type WorkloadProfile struct {
	Journey     string `json:"journey"`
	ReadQPS     int    `json:"read_qps"`
	WriteQPS    int    `json:"write_qps"`
	Fanout      int    `json:"fanout"`
	BurstFactor int    `json:"burst_factor"`
	HotKeyRisk  bool   `json:"hot_key_risk"`
}

func ValidateWorkloadProfile(profile WorkloadProfile) []string {
	var issues []string
	if profile.Journey == "" {
		issues = append(issues, "workload profile needs a user journey")
	}
	if profile.ReadQPS < 0 || profile.WriteQPS < 0 {
		issues = append(issues, "read and write QPS must be non-negative")
	}
	if profile.ReadQPS == 0 && profile.WriteQPS == 0 {
		issues = append(issues, "workload profile should include read or write traffic")
	}
	if profile.BurstFactor < 1 {
		issues = append(issues, "burst factor should be at least 1")
	}
	if profile.Fanout < 1 {
		issues = append(issues, "fanout should be at least 1")
	}
	return issues
}

func DominantPath(profile WorkloadProfile) string {
	totalRead := profile.ReadQPS * profile.BurstFactor
	totalWrite := profile.WriteQPS * profile.BurstFactor * profile.Fanout

	switch {
	case totalRead > totalWrite:
		return "read"
	case totalWrite > totalRead:
		return "write"
	default:
		return "balanced"
	}
}

func main() {
	var path string
	flag.StringVar(&path, "workload", "", "path to a workload profile JSON file")
	flag.Parse()

	if path == "" {
		fmt.Fprintln(os.Stderr, "missing -workload path")
		os.Exit(2)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read workload: %v\n", err)
		os.Exit(1)
	}

	var profile WorkloadProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		fmt.Fprintf(os.Stderr, "decode workload: %v\n", err)
		os.Exit(1)
	}

	if issues := ValidateWorkloadProfile(profile); len(issues) > 0 {
		for _, issue := range issues {
			fmt.Println(issue)
		}
		os.Exit(1)
	}

	fmt.Printf("dominant path: %s\n", DominantPath(profile))
}
