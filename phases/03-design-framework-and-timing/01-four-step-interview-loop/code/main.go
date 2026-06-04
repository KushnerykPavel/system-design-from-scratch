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

type Plan struct {
	Prompt string  `json:"prompt"`
	Stages []Stage `json:"stages"`
}

func DefaultPlan(prompt string) Plan {
	return Plan{
		Prompt: prompt,
		Stages: []Stage{
			{Name: "clarify", Minutes: 7},
			{Name: "size", Minutes: 5},
			{Name: "high_level_design", Minutes: 12},
			{Name: "deep_dive", Minutes: 14},
			{Name: "wrap_up", Minutes: 7},
		},
	}
}

func TotalMinutes(plan Plan) int {
	total := 0
	for _, stage := range plan.Stages {
		total += stage.Minutes
	}
	return total
}

func ValidatePlan(plan Plan) []string {
	var issues []string
	if len(plan.Stages) < 5 {
		issues = append(issues, "plan should contain clarify, size, high-level design, deep dive, and wrap-up")
	}

	required := []string{"clarify", "size", "high_level_design", "deep_dive", "wrap_up"}
	seen := map[string]bool{}
	for _, stage := range plan.Stages {
		if stage.Minutes <= 0 {
			issues = append(issues, fmt.Sprintf("stage %q must have positive minutes", stage.Name))
		}
		seen[stage.Name] = true
	}

	for _, name := range required {
		if !seen[name] {
			issues = append(issues, fmt.Sprintf("missing required stage %q", name))
		}
	}

	if total := TotalMinutes(plan); total > 45 {
		issues = append(issues, fmt.Sprintf("plan uses %d minutes; must fit within 45", total))
	}

	return issues
}

func main() {
	prompt := flag.String("prompt", "Design a rate limiter", "interview prompt to build the plan for")
	flag.Parse()

	plan := DefaultPlan(*prompt)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(plan)
}
