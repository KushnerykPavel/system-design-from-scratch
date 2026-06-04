package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
)

type Priority struct {
	Requirement string `json:"requirement"`
	Rank        int    `json:"rank"`
	Rationale   string `json:"rationale"`
}

type PrioritySet struct {
	Priorities []Priority `json:"priorities"`
}

func ValidatePrioritySet(set PrioritySet) []string {
	var issues []string
	if len(set.Priorities) == 0 {
		issues = append(issues, "record at least one priority")
		return issues
	}
	if len(set.Priorities) > 5 {
		issues = append(issues, "too many priorities; keep the ranking focused")
	}

	seenRanks := map[int]bool{}
	for _, priority := range set.Priorities {
		if priority.Requirement == "" {
			issues = append(issues, "priority requirement is required")
		}
		if priority.Rank <= 0 {
			issues = append(issues, fmt.Sprintf("priority %q must have positive rank", priority.Requirement))
		}
		if seenRanks[priority.Rank] {
			issues = append(issues, fmt.Sprintf("duplicate rank %d creates ambiguous ordering", priority.Rank))
		}
		seenRanks[priority.Rank] = true
		if priority.Rationale == "" {
			issues = append(issues, fmt.Sprintf("priority %q should include rationale", priority.Requirement))
		}
	}

	return issues
}

func SortedPriorities(set PrioritySet) []Priority {
	sorted := append([]Priority(nil), set.Priorities...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Rank < sorted[j].Rank
	})
	return sorted
}

func main() {
	var path string
	flag.StringVar(&path, "priorities", "", "path to a priority set JSON file")
	flag.Parse()

	if path == "" {
		fmt.Fprintln(os.Stderr, "missing -priorities path")
		os.Exit(2)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read priorities: %v\n", err)
		os.Exit(1)
	}

	var set PrioritySet
	if err := json.Unmarshal(data, &set); err != nil {
		fmt.Fprintf(os.Stderr, "decode priorities: %v\n", err)
		os.Exit(1)
	}

	if issues := ValidatePrioritySet(set); len(issues) > 0 {
		for _, issue := range issues {
			fmt.Println(issue)
		}
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(SortedPriorities(set))
}
