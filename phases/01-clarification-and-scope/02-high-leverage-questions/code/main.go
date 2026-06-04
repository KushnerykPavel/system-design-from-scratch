package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Question struct {
	Text          string `json:"text"`
	Pivot         string `json:"pivot"`
	ChangesDesign bool   `json:"changes_design"`
}

type QuestionSet struct {
	Questions []Question `json:"questions"`
}

var highLeveragePivots = map[string]bool{
	"workload":    true,
	"consistency": true,
	"scope":       true,
	"failure":     true,
	"cost":        true,
	"security":    true,
}

func ScoreQuestions(set QuestionSet) (int, []string) {
	score := 0
	var issues []string

	if len(set.Questions) == 0 {
		issues = append(issues, "at least one clarifying question is required")
		return score, issues
	}
	if len(set.Questions) > 5 {
		issues = append(issues, "too many opening clarification questions; keep the first set tight")
	}

	for _, question := range set.Questions {
		if !highLeveragePivots[question.Pivot] {
			issues = append(issues, fmt.Sprintf("question %q uses low-signal or unknown pivot %q", question.Text, question.Pivot))
			continue
		}
		score++
		if question.ChangesDesign {
			score++
		} else {
			issues = append(issues, fmt.Sprintf("question %q does not state how it changes the design", question.Text))
		}
	}

	return score, issues
}

func main() {
	var path string
	flag.StringVar(&path, "questions", "", "path to a question set JSON file")
	flag.Parse()

	if path == "" {
		fmt.Fprintln(os.Stderr, "missing -questions path")
		os.Exit(2)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read questions: %v\n", err)
		os.Exit(1)
	}

	var set QuestionSet
	if err := json.Unmarshal(data, &set); err != nil {
		fmt.Fprintf(os.Stderr, "decode questions: %v\n", err)
		os.Exit(1)
	}

	score, issues := ScoreQuestions(set)
	if len(issues) > 0 {
		for _, issue := range issues {
			fmt.Println(issue)
		}
	}
	fmt.Printf("question score: %d\n", score)
	if len(issues) > 0 {
		os.Exit(1)
	}
}
