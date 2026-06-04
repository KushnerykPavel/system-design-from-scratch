package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type ReframedPrompt struct {
	System        string `json:"system"`
	CoreWorkflow  string `json:"core_workflow"`
	ScopeCut      string `json:"scope_cut"`
	Priority      string `json:"priority"`
	Assumption    string `json:"assumption"`
	WorkloadShape string `json:"workload_shape"`
}

func MissingIngredients(prompt ReframedPrompt) []string {
	var missing []string
	if prompt.System == "" {
		missing = append(missing, "system")
	}
	if prompt.CoreWorkflow == "" {
		missing = append(missing, "core_workflow")
	}
	if prompt.ScopeCut == "" {
		missing = append(missing, "scope_cut")
	}
	if prompt.Priority == "" {
		missing = append(missing, "priority")
	}
	if prompt.Assumption == "" {
		missing = append(missing, "assumption")
	}
	if prompt.WorkloadShape == "" {
		missing = append(missing, "workload_shape")
	}
	return missing
}

func ScoreReframing(prompt ReframedPrompt) int {
	return 6 - len(MissingIngredients(prompt))
}

func main() {
	var path string
	flag.StringVar(&path, "prompt", "", "path to a reframed prompt JSON file")
	flag.Parse()

	if path == "" {
		fmt.Fprintln(os.Stderr, "missing -prompt path")
		os.Exit(2)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read prompt: %v\n", err)
		os.Exit(1)
	}

	var prompt ReframedPrompt
	if err := json.Unmarshal(data, &prompt); err != nil {
		fmt.Fprintf(os.Stderr, "decode prompt: %v\n", err)
		os.Exit(1)
	}

	missing := MissingIngredients(prompt)
	if len(missing) > 0 {
		fmt.Printf("missing: %v\n", missing)
		os.Exit(1)
	}

	fmt.Printf("reframing score: %d\n", ScoreReframing(prompt))
}
