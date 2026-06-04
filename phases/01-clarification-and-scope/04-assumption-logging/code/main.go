package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Assumption struct {
	Statement  string `json:"statement"`
	Category   string `json:"category"`
	Impact     string `json:"impact"`
	Reversible bool   `json:"reversible"`
}

type AssumptionLog struct {
	Assumptions []Assumption `json:"assumptions"`
}

var validAssumptionCategories = map[string]bool{
	"scale":       true,
	"consistency": true,
	"geography":   true,
	"scope":       true,
	"abuse":       true,
}

func ValidateAssumptionLog(log AssumptionLog) []string {
	var issues []string
	if len(log.Assumptions) == 0 {
		issues = append(issues, "record at least one explicit assumption")
	}

	for _, assumption := range log.Assumptions {
		if assumption.Statement == "" {
			issues = append(issues, "assumption statement is required")
		}
		if !validAssumptionCategories[assumption.Category] {
			issues = append(issues, fmt.Sprintf("assumption %q has invalid category %q", assumption.Statement, assumption.Category))
		}
		if assumption.Impact == "" {
			issues = append(issues, fmt.Sprintf("assumption %q must explain design impact", assumption.Statement))
		}
	}

	return issues
}

func CountExpensiveReversals(log AssumptionLog) int {
	count := 0
	for _, assumption := range log.Assumptions {
		if !assumption.Reversible {
			count++
		}
	}
	return count
}

func main() {
	var path string
	flag.StringVar(&path, "assumptions", "", "path to an assumption log JSON file")
	flag.Parse()

	if path == "" {
		fmt.Fprintln(os.Stderr, "missing -assumptions path")
		os.Exit(2)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read assumptions: %v\n", err)
		os.Exit(1)
	}

	var log AssumptionLog
	if err := json.Unmarshal(data, &log); err != nil {
		fmt.Fprintf(os.Stderr, "decode assumptions: %v\n", err)
		os.Exit(1)
	}

	if issues := ValidateAssumptionLog(log); len(issues) > 0 {
		for _, issue := range issues {
			fmt.Println(issue)
		}
		os.Exit(1)
	}

	fmt.Printf("expensive reversals: %d\n", CountExpensiveReversals(log))
}
