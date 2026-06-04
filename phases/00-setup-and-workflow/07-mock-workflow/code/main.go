package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Stage struct {
	Name    string `json:"name"`
	Minutes int    `json:"minutes"`
}

type SessionPlan struct {
	DurationMinutes int     `json:"duration_minutes"`
	Stages          []Stage `json:"stages"`
}

func TotalMinutes(plan SessionPlan) int {
	total := 0
	for _, stage := range plan.Stages {
		total += stage.Minutes
	}
	return total
}

func ValidateSessionPlan(plan SessionPlan) []string {
	var issues []string
	if plan.DurationMinutes <= 0 {
		issues = append(issues, "duration_minutes must be positive")
	}
	if len(plan.Stages) == 0 {
		issues = append(issues, "session requires stages")
		return issues
	}

	required := map[string]bool{
		"pre_brief": false,
		"live":      false,
		"feedback":  false,
		"debrief":   false,
	}

	for _, stage := range plan.Stages {
		if stage.Minutes <= 0 {
			issues = append(issues, fmt.Sprintf("stage %q must have positive minutes", stage.Name))
		}
		if _, ok := required[stage.Name]; ok {
			required[stage.Name] = true
		}
	}

	for name, present := range required {
		if !present {
			issues = append(issues, fmt.Sprintf("missing required stage %q", name))
		}
	}

	if total := TotalMinutes(plan); total > plan.DurationMinutes {
		issues = append(issues, fmt.Sprintf("session uses %d minutes but only %d are available", total, plan.DurationMinutes))
	}

	return issues
}

func main() {
	var path string
	flag.StringVar(&path, "plan", "", "path to a mock session plan JSON file")
	flag.Parse()

	if path == "" {
		fmt.Fprintln(os.Stderr, "missing -plan path")
		os.Exit(2)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read plan: %v\n", err)
		os.Exit(1)
	}

	var plan SessionPlan
	if err := json.Unmarshal(data, &plan); err != nil {
		fmt.Fprintf(os.Stderr, "decode plan: %v\n", err)
		os.Exit(1)
	}

	issues := ValidateSessionPlan(plan)
	if len(issues) > 0 {
		for _, issue := range issues {
			fmt.Println(issue)
		}
		os.Exit(1)
	}

	fmt.Printf("mock session validated in %d minutes\n", TotalMinutes(plan))
}
