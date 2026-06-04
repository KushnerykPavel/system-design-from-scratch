package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type PracticeSummary struct {
	ClarifiedScope     bool `json:"clarified_scope"`
	RankedRequirements bool `json:"ranked_requirements"`
	LoggedAssumptions  bool `json:"logged_assumptions"`
	NamedWorkloadShape bool `json:"named_workload_shape"`
	StatedScopeCut     bool `json:"stated_scope_cut"`
}

func DetectAntiPatterns(summary PracticeSummary) []string {
	var issues []string
	if !summary.ClarifiedScope {
		issues = append(issues, "architecture-first without scope clarification")
	}
	if !summary.RankedRequirements {
		issues = append(issues, "requirements were listed but not prioritized")
	}
	if !summary.LoggedAssumptions {
		issues = append(issues, "silent assumptions left ambiguity unmanaged")
	}
	if !summary.NamedWorkloadShape {
		issues = append(issues, "workload shape was not translated from the user journey")
	}
	if !summary.StatedScopeCut {
		issues = append(issues, "scope cuts were implicit or absent")
	}
	return issues
}

func main() {
	var path string
	flag.StringVar(&path, "summary", "", "path to a practice summary JSON file")
	flag.Parse()

	if path == "" {
		fmt.Fprintln(os.Stderr, "missing -summary path")
		os.Exit(2)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read summary: %v\n", err)
		os.Exit(1)
	}

	var summary PracticeSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		fmt.Fprintf(os.Stderr, "decode summary: %v\n", err)
		os.Exit(1)
	}

	issues := DetectAntiPatterns(summary)
	for _, issue := range issues {
		fmt.Println(issue)
	}
	if len(issues) > 0 {
		os.Exit(1)
	}
}
