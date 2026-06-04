package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type ScopePlan struct {
	CoreWorkflows         []string `json:"core_workflows"`
	DeferredFeatures      []string `json:"deferred_features"`
	Reason                string   `json:"reason"`
	PreservesPromptIntent bool     `json:"preserves_prompt_intent"`
}

func ValidateScopePlan(plan ScopePlan) []string {
	var issues []string

	if len(plan.CoreWorkflows) == 0 {
		issues = append(issues, "keep at least one core workflow in scope")
	}
	if len(plan.DeferredFeatures) == 0 {
		issues = append(issues, "name at least one deferred feature so the cut is explicit")
	}
	if plan.Reason == "" {
		issues = append(issues, "scope cuts should include a reason")
	}
	if !plan.PreservesPromptIntent {
		issues = append(issues, "scope cut no longer preserves the original prompt intent")
	}

	return issues
}

func ComplexityReduction(plan ScopePlan) int {
	return len(plan.DeferredFeatures) - len(plan.CoreWorkflows)
}

func main() {
	var path string
	flag.StringVar(&path, "plan", "", "path to a scope plan JSON file")
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

	var plan ScopePlan
	if err := json.Unmarshal(data, &plan); err != nil {
		fmt.Fprintf(os.Stderr, "decode plan: %v\n", err)
		os.Exit(1)
	}

	if issues := ValidateScopePlan(plan); len(issues) > 0 {
		for _, issue := range issues {
			fmt.Println(issue)
		}
		os.Exit(1)
	}

	fmt.Printf("complexity reduction score: %d\n", ComplexityReduction(plan))
}
