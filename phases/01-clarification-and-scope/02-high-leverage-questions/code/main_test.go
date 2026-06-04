package main

import "testing"

func TestScoreQuestionsRewardsHighLeveragePivots(t *testing.T) {
	t.Parallel()

	set := QuestionSet{
		Questions: []Question{
			{Text: "Is the workload read-heavy or write-heavy?", Pivot: "workload", ChangesDesign: true},
			{Text: "Do we need strong read-after-write semantics?", Pivot: "consistency", ChangesDesign: true},
		},
	}

	score, issues := ScoreQuestions(set)
	if score != 4 {
		t.Fatalf("ScoreQuestions() score = %d, want 4", score)
	}
	if len(issues) != 0 {
		t.Fatalf("ScoreQuestions() issues = %v, want none", issues)
	}
}

func TestScoreQuestionsFlagsLowSignalQuestion(t *testing.T) {
	t.Parallel()

	set := QuestionSet{
		Questions: []Question{
			{Text: "Should the button be blue?", Pivot: "ui", ChangesDesign: false},
		},
	}

	_, issues := ScoreQuestions(set)
	if len(issues) == 0 {
		t.Fatal("ScoreQuestions() returned no issues for a low-signal question")
	}
}

func TestScoreQuestionsFlagsOverBudgetQuestionSet(t *testing.T) {
	t.Parallel()

	set := QuestionSet{
		Questions: []Question{
			{Text: "q1", Pivot: "workload", ChangesDesign: true},
			{Text: "q2", Pivot: "scope", ChangesDesign: true},
			{Text: "q3", Pivot: "failure", ChangesDesign: true},
			{Text: "q4", Pivot: "cost", ChangesDesign: true},
			{Text: "q5", Pivot: "security", ChangesDesign: true},
			{Text: "q6", Pivot: "consistency", ChangesDesign: true},
		},
	}

	_, issues := ScoreQuestions(set)
	if len(issues) == 0 {
		t.Fatal("ScoreQuestions() returned no issues for too many opening questions")
	}
}
