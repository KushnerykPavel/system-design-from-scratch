package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Component struct {
	Name     string `json:"name"`
	Critical bool   `json:"critical"`
}

type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Diagram struct {
	Components []Component `json:"components"`
	Edges      []Edge      `json:"edges"`
}

func ComplexityScore(diagram Diagram) int {
	return len(diagram.Components)*2 + len(diagram.Edges)
}

func ValidateFirstPass(diagram Diagram) []string {
	var issues []string
	if len(diagram.Components) < 3 {
		issues = append(issues, "diagram should show at least client, service, and state boundaries")
	}
	if len(diagram.Components) > 8 {
		issues = append(issues, "diagram exceeds the recommended first-pass box budget of 8")
	}
	if len(diagram.Edges) == 0 {
		issues = append(issues, "diagram should show at least one request path")
	}

	critical := 0
	for _, component := range diagram.Components {
		if component.Critical {
			critical++
		}
	}
	if critical == 0 {
		issues = append(issues, "diagram should mark at least one critical boundary for deep dive")
	}

	return issues
}

func main() {
	var path string
	flag.StringVar(&path, "diagram", "", "path to a diagram JSON file")
	flag.Parse()

	if path == "" {
		fmt.Fprintln(os.Stderr, "missing -diagram path")
		os.Exit(2)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read diagram: %v\n", err)
		os.Exit(1)
	}

	var diagram Diagram
	if err := json.Unmarshal(data, &diagram); err != nil {
		fmt.Fprintf(os.Stderr, "decode diagram: %v\n", err)
		os.Exit(1)
	}

	issues := ValidateFirstPass(diagram)
	if len(issues) > 0 {
		for _, issue := range issues {
			fmt.Println(issue)
		}
		os.Exit(1)
	}

	fmt.Printf("diagram complexity score: %d\n", ComplexityScore(diagram))
}
