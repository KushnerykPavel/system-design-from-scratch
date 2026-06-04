package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
)

type Requirement struct {
	Text     string `json:"text"`
	Kind     string `json:"kind"`
	Priority int    `json:"priority"`
	Driver   bool   `json:"driver"`
}

type RequirementSet struct {
	Requirements []Requirement `json:"requirements"`
}

func ValidateRequirementSet(set RequirementSet) []string {
	var issues []string

	functional := 0
	nonFunctional := 0
	drivers := 0

	for _, requirement := range set.Requirements {
		switch requirement.Kind {
		case "functional":
			functional++
		case "non_functional":
			nonFunctional++
		default:
			issues = append(issues, fmt.Sprintf("requirement %q has invalid kind %q", requirement.Text, requirement.Kind))
		}

		if requirement.Priority <= 0 {
			issues = append(issues, fmt.Sprintf("requirement %q must have positive priority", requirement.Text))
		}
		if requirement.Driver {
			drivers++
		}
	}

	if functional == 0 {
		issues = append(issues, "at least one functional requirement is required")
	}
	if nonFunctional == 0 {
		issues = append(issues, "at least one non-functional requirement is required")
	}
	if drivers == 0 {
		issues = append(issues, "mark one dominant non-functional driver")
	}
	if drivers > 1 {
		issues = append(issues, "too many dominant drivers; choose one primary design driver")
	}

	return issues
}

func RankedNonFunctional(set RequirementSet) []Requirement {
	var ranked []Requirement
	for _, requirement := range set.Requirements {
		if requirement.Kind == "non_functional" {
			ranked = append(ranked, requirement)
		}
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Priority < ranked[j].Priority
	})

	return ranked
}

func main() {
	var path string
	flag.StringVar(&path, "requirements", "", "path to a requirement set JSON file")
	flag.Parse()

	if path == "" {
		fmt.Fprintln(os.Stderr, "missing -requirements path")
		os.Exit(2)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read requirements: %v\n", err)
		os.Exit(1)
	}

	var set RequirementSet
	if err := json.Unmarshal(data, &set); err != nil {
		fmt.Fprintf(os.Stderr, "decode requirements: %v\n", err)
		os.Exit(1)
	}

	if issues := ValidateRequirementSet(set); len(issues) > 0 {
		for _, issue := range issues {
			fmt.Println(issue)
		}
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(RankedNonFunctional(set))
}
