package main

import "testing"

func TestValidateSnapshot(t *testing.T) {
	t.Parallel()

	quizScore := 7
	reviewInterval := 4
	reviewEase := 2.3
	lapseCount := 0

	snapshot := ProgressSnapshot{
		SchemaVersion: 1,
		Lessons: []LessonProgress{
			{
				Lesson:             "00-setup-and-workflow/01-repo-setup-and-progress",
				Status:             "done",
				LastUpdated:        "2026-06-03",
				NotesPath:          "notes/01.md",
				ArtifactPaths:      []string{"outputs/checklist.md"},
				QuizScore:          &quizScore,
				Confidence:         "medium",
				MistakeTags:        []string{"no_sizing"},
				LastReviewedAt:     "2026-06-03",
				NextReviewAt:       "2026-06-07",
				ReviewIntervalDays: &reviewInterval,
				ReviewEase:         &reviewEase,
				LapseCount:         &lapseCount,
				QuizHistory: []QuizAttempt{
					{Score: 6, CompletedAt: "2026-06-01", Stage: "post"},
					{Score: 7, CompletedAt: "2026-06-03", Stage: "post"},
				},
				ModeHistory: []SessionRecord{
					{Mode: "learn", CompletedAt: "2026-06-01"},
					{Mode: "practice", CompletedAt: "2026-06-03"},
				},
				FeedbackHistory: []SessionFeedback{
					{
						SessionType:                "lesson",
						CompletedAt:                "2026-06-03",
						Summary:                    "Sizing improved, but failure handling stayed shallow.",
						Strengths:                  []string{"Clarified the write/read split before choosing storage."},
						Gaps:                       []string{"Did not explain how alerts detect backlog growth."},
						HighestLeverageImprovement: "Practice tying one concrete metric and one alert to each failure mode.",
						Dimensions: []DimensionFeedback{
							{Dimension: "clarification", Score: 3, Evidence: "Asked about read/write ratio and durability."},
							{Dimension: "sizing", Score: 3, Evidence: "Estimated QPS and storage growth before architecture."},
							{Dimension: "observability", Score: 2, Evidence: "Mentioned metrics but not alert thresholds.", NextAction: "Add one alert and one dashboard per critical path."},
						},
					},
				},
			},
		},
	}

	if issues := ValidateSnapshot(snapshot); len(issues) != 0 {
		t.Fatalf("ValidateSnapshot() returned issues: %v", issues)
	}
}

func TestValidateSnapshotRejectsBadState(t *testing.T) {
	t.Parallel()

	quizScore := 11
	reviewInterval := -1
	reviewEase := 0.5
	lapseCount := -2

	snapshot := ProgressSnapshot{
		SchemaVersion: 99,
		Lessons: []LessonProgress{
			{
				Lesson:             "a",
				Status:             "wat",
				LastUpdated:        "",
				NotesPath:          "",
				QuizScore:          &quizScore,
				Confidence:         "certain",
				MistakeTags:        []string{"", "made_up_tag"},
				ReviewIntervalDays: &reviewInterval,
				ReviewEase:         &reviewEase,
				LapseCount:         &lapseCount,
				QuizHistory:        []QuizAttempt{{Score: 9}},
				ModeHistory:        []SessionRecord{{Mode: "coast"}},
				FeedbackHistory: []SessionFeedback{
					{
						SessionType:                "brainstorm",
						CompletedAt:                "",
						HighestLeverageImprovement: "",
						Strengths:                  []string{""},
						Gaps:                       []string{""},
						Dimensions: []DimensionFeedback{
							{Dimension: "mystery", Score: 7},
						},
					},
				},
			},
			{Lesson: "a", Status: "done", LastUpdated: "2026-06-03", NotesPath: "notes.md"},
		},
	}

	if issues := ValidateSnapshot(snapshot); len(issues) < 3 {
		t.Fatalf("ValidateSnapshot() returned %d issues, want at least 3: %v", len(issues), issues)
	}
}

func TestRecommendNext(t *testing.T) {
	t.Parallel()

	snapshot := ProgressSnapshot{
		Lessons: []LessonProgress{
			{Lesson: "01", Status: "done"},
			{Lesson: "02", Status: "not_started"},
			{Lesson: "03", Status: "in_progress"},
		},
	}

	if got := RecommendNext(snapshot); got != "02" {
		t.Fatalf("RecommendNext() = %q, want %q", got, "02")
	}
}

func TestSyncReviewSchedule(t *testing.T) {
	t.Parallel()

	quizScore := 8
	snapshot := ProgressSnapshot{
		Lessons: []LessonProgress{
			{
				Lesson:      "clarity",
				Status:      "done",
				LastUpdated: "2026-06-03",
				NotesPath:   "notes/clarity.md",
				QuizScore:   &quizScore,
				Confidence:  "high",
				MistakeTags: []string{"weak_tradeoffs"},
				FeedbackHistory: []SessionFeedback{
					{
						SessionType:                "lesson",
						CompletedAt:                "2026-06-03",
						HighestLeverageImprovement: "Keep the same cadence.",
						Dimensions: []DimensionFeedback{
							{Dimension: "clarification", Score: 4},
							{Dimension: "sizing", Score: 3},
						},
					},
				},
			},
			{
				Lesson:      "repair",
				Status:      "assisted",
				LastUpdated: "2026-06-03",
				NotesPath:   "notes/repair.md",
				MistakeTags: []string{"missing_observability", "shallow_failure_modes"},
				FeedbackHistory: []SessionFeedback{
					{
						SessionType:                "lesson",
						CompletedAt:                "2026-06-03",
						HighestLeverageImprovement: "Retest tomorrow.",
						Dimensions: []DimensionFeedback{
							{Dimension: "observability", Score: 2},
						},
					},
				},
			},
		},
	}

	if err := SyncReviewSchedule(&snapshot, "2026-06-03"); err != nil {
		t.Fatalf("SyncReviewSchedule() error = %v", err)
	}

	first := snapshot.Lessons[0]
	if got := derefInt(first.ReviewIntervalDays, 0); got != 13 {
		t.Fatalf("first interval = %d, want 13", got)
	}
	if first.NextReviewAt != "2026-06-16" {
		t.Fatalf("first NextReviewAt = %q, want %q", first.NextReviewAt, "2026-06-16")
	}

	second := snapshot.Lessons[1]
	if got := derefInt(second.ReviewIntervalDays, 0); got != 1 {
		t.Fatalf("second interval = %d, want 1", got)
	}
	if second.NextReviewAt != "2026-06-04" {
		t.Fatalf("second NextReviewAt = %q, want %q", second.NextReviewAt, "2026-06-04")
	}
	if got := derefInt(second.LapseCount, 0); got != 1 {
		t.Fatalf("second lapse_count = %d, want 1", got)
	}
}

func TestRecommendReviews(t *testing.T) {
	t.Parallel()

	snapshot := ProgressSnapshot{
		Lessons: []LessonProgress{
			{
				Lesson:       "overdue",
				Status:       "done",
				LastUpdated:  "2026-06-03",
				NotesPath:    "notes/overdue.md",
				NextReviewAt: "2026-06-01",
				Confidence:   "low",
				MistakeTags:  []string{"missing_observability"},
			},
			{
				Lesson:       "today",
				Status:       "done",
				LastUpdated:  "2026-06-03",
				NotesPath:    "notes/today.md",
				NextReviewAt: "2026-06-03",
				Confidence:   "medium",
				MistakeTags:  []string{"missing_observability"},
			},
			{
				Lesson:       "later",
				Status:       "done",
				LastUpdated:  "2026-06-03",
				NotesPath:    "notes/later.md",
				NextReviewAt: "2026-06-10",
			},
		},
	}

	items, err := RecommendReviews(snapshot, "2026-06-03")
	if err != nil {
		t.Fatalf("RecommendReviews() error = %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("RecommendReviews() len = %d, want 2", len(items))
	}
	if items[0].Lesson != "overdue" || items[0].State != "overdue" || items[0].DaysLate != 2 {
		t.Fatalf("first item = %+v, want overdue lesson 2 days late", items[0])
	}
	if items[0].TopMistake != "missing_observability" {
		t.Fatalf("first item top mistake = %q, want missing_observability", items[0].TopMistake)
	}
	if items[1].Lesson != "today" || items[1].State != "due_today" {
		t.Fatalf("second item = %+v, want due_today lesson", items[1])
	}
}

func TestRecommendActionPrefersDueReview(t *testing.T) {
	t.Parallel()

	snapshot := ProgressSnapshot{
		Lessons: []LessonProgress{
			{
				Lesson:       "due-review",
				Status:       "done",
				LastUpdated:  "2026-06-03",
				NotesPath:    "notes/due.md",
				NextReviewAt: "2026-06-02",
				MistakeTags:  []string{"no_sizing", "missing_observability"},
			},
			{
				Lesson:      "in-flight",
				Status:      "in_progress",
				LastUpdated: "2026-06-03",
				NotesPath:   "notes/in-flight.md",
				MistakeTags: []string{"no_sizing"},
			},
		},
	}

	got, err := RecommendAction(snapshot, "2026-06-03")
	if err != nil {
		t.Fatalf("RecommendAction() error = %v", err)
	}
	if got.Action != "practice" || got.Lesson != "due-review" {
		t.Fatalf("RecommendAction() = %+v, want practice recommendation for due review", got)
	}
	if got.DrillKind != "capacity_drill" {
		t.Fatalf("RecommendAction() drill kind = %q, want capacity_drill", got.DrillKind)
	}
}

func TestRecommendActionFallsBackToAssistedThenInProgress(t *testing.T) {
	t.Parallel()

	snapshot := ProgressSnapshot{
		Lessons: []LessonProgress{
			{
				Lesson:      "retry-me",
				Status:      "assisted",
				LastUpdated: "2026-06-03",
				NotesPath:   "notes/retry.md",
				MistakeTags: []string{"weak_tradeoffs"},
			},
			{
				Lesson:      "resume-me",
				Status:      "in_progress",
				LastUpdated: "2026-06-03",
				NotesPath:   "notes/resume.md",
			},
		},
	}

	got, err := RecommendAction(snapshot, "2026-06-03")
	if err != nil {
		t.Fatalf("RecommendAction() error = %v", err)
	}
	if got.Action != "revisit" || got.Lesson != "retry-me" {
		t.Fatalf("RecommendAction() = %+v, want assisted retry recommendation", got)
	}
}

func TestTrendHelpers(t *testing.T) {
	t.Parallel()

	snapshot := ProgressSnapshot{
		Lessons: []LessonProgress{
			{
				Lesson:      "lesson-a",
				Status:      "done",
				LastUpdated: "2026-06-03",
				NotesPath:   "notes/a.md",
				FeedbackHistory: []SessionFeedback{
					{
						SessionType:                "lesson",
						CompletedAt:                "2026-06-01",
						HighestLeverageImprovement: "Keep going.",
						Dimensions: []DimensionFeedback{
							{Dimension: "clarification", Score: 2},
							{Dimension: "observability", Score: 1},
						},
					},
					{
						SessionType:                "lesson",
						CompletedAt:                "2026-06-03",
						HighestLeverageImprovement: "Keep going.",
						Dimensions: []DimensionFeedback{
							{Dimension: "clarification", Score: 4},
							{Dimension: "observability", Score: 2},
						},
					},
				},
			},
			{
				Lesson:      "lesson-b",
				Status:      "done",
				LastUpdated: "2026-06-03",
				NotesPath:   "notes/b.md",
				FeedbackHistory: []SessionFeedback{
					{
						SessionType:                "lesson",
						CompletedAt:                "2026-06-03",
						HighestLeverageImprovement: "Keep going.",
						Dimensions: []DimensionFeedback{
							{Dimension: "clarification", Score: 3},
							{Dimension: "observability", Score: 1},
						},
					},
				},
			},
		},
	}

	weak := WeakestTrend(snapshot)
	if weak.Dimension != "observability" {
		t.Fatalf("WeakestTrend() = %+v, want observability as weakest", weak)
	}

	improving := ImprovingTrend(snapshot)
	if improving.Dimension != "clarification" {
		t.Fatalf("ImprovingTrend() = %+v, want clarification as improving", improving)
	}
}

func TestWeeklyMistakeSummary(t *testing.T) {
	t.Parallel()

	snapshot := ProgressSnapshot{
		Lessons: []LessonProgress{
			{
				Lesson:      "lesson-a",
				Status:      "done",
				LastUpdated: "2026-06-03",
				NotesPath:   "notes/a.md",
				MistakeTags: []string{"no_sizing", "missing_observability"},
				FeedbackHistory: []SessionFeedback{
					{SessionType: "lesson", CompletedAt: "2026-06-03", HighestLeverageImprovement: "x"},
				},
			},
			{
				Lesson:      "lesson-b",
				Status:      "done",
				LastUpdated: "2026-06-02",
				NotesPath:   "notes/b.md",
				MistakeTags: []string{"missing_observability"},
				FeedbackHistory: []SessionFeedback{
					{SessionType: "lesson", CompletedAt: "2026-06-02", HighestLeverageImprovement: "x"},
				},
			},
			{
				Lesson:      "lesson-old",
				Status:      "done",
				LastUpdated: "2026-05-20",
				NotesPath:   "notes/old.md",
				MistakeTags: []string{"no_sizing"},
				FeedbackHistory: []SessionFeedback{
					{SessionType: "lesson", CompletedAt: "2026-05-20", HighestLeverageImprovement: "x"},
				},
			},
		},
	}

	trends, err := WeeklyMistakeSummary(snapshot, "2026-06-07")
	if err != nil {
		t.Fatalf("WeeklyMistakeSummary() error = %v", err)
	}
	if len(trends) < 2 {
		t.Fatalf("WeeklyMistakeSummary() len = %d, want at least 2", len(trends))
	}
	if trends[0].Tag != "missing_observability" || trends[0].Count != 2 {
		t.Fatalf("top trend = %+v, want missing_observability x2", trends[0])
	}
	if trends[0].Drill == "" {
		t.Fatalf("top trend should include a drill recommendation")
	}
}

func TestRecommendActionFallsBackToWeeklyPractice(t *testing.T) {
	t.Parallel()

	snapshot := ProgressSnapshot{
		Lessons: []LessonProgress{
			{
				Lesson:      "a",
				Status:      "done",
				LastUpdated: "2026-06-03",
				NotesPath:   "notes/a.md",
				MistakeTags: []string{"missing_observability"},
				FeedbackHistory: []SessionFeedback{
					{SessionType: "lesson", CompletedAt: "2026-06-03", HighestLeverageImprovement: "x"},
				},
			},
			{
				Lesson:      "b",
				Status:      "done",
				LastUpdated: "2026-06-02",
				NotesPath:   "notes/b.md",
				MistakeTags: []string{"missing_observability"},
				FeedbackHistory: []SessionFeedback{
					{SessionType: "lesson", CompletedAt: "2026-06-02", HighestLeverageImprovement: "x"},
				},
			},
		},
	}

	got, err := RecommendAction(snapshot, "2026-06-07")
	if err != nil {
		t.Fatalf("RecommendAction() error = %v", err)
	}
	if got.Action != "practice" {
		t.Fatalf("RecommendAction() action = %q, want practice", got.Action)
	}
	if got.DrillKind != "failure_drill" {
		t.Fatalf("RecommendAction() drill kind = %q, want failure_drill", got.DrillKind)
	}
}

func TestBuildProgressStats(t *testing.T) {
	t.Parallel()

	quizSeven := 7
	quizFive := 5
	reviews := []ReviewItem{
		{Lesson: "a", State: "due_today"},
		{Lesson: "b", State: "overdue"},
	}
	snapshot := ProgressSnapshot{
		Lessons: []LessonProgress{
			{
				Lesson:      "a",
				Status:      "done",
				LastUpdated: "2026-06-04",
				NotesPath:   "notes/a.md",
				QuizScore:   &quizSeven,
				FeedbackHistory: []SessionFeedback{
					{
						SessionType:                "lesson",
						CompletedAt:                "2026-06-04",
						HighestLeverageImprovement: "x",
						Dimensions: []DimensionFeedback{
							{Dimension: "clarification", Score: 4},
							{Dimension: "observability", Score: 2},
						},
					},
				},
			},
			{
				Lesson:      "b",
				Status:      "assisted",
				LastUpdated: "2026-06-03",
				NotesPath:   "notes/b.md",
				QuizScore:   &quizFive,
				FeedbackHistory: []SessionFeedback{
					{
						SessionType:                "lesson",
						CompletedAt:                "2026-06-03",
						HighestLeverageImprovement: "x",
						Dimensions: []DimensionFeedback{
							{Dimension: "clarification", Score: 3},
							{Dimension: "observability", Score: 1},
						},
					},
				},
			},
			{
				Lesson:      "c",
				Status:      "in_progress",
				LastUpdated: "2026-06-01",
				NotesPath:   "notes/c.md",
			},
			{
				Lesson:      "d",
				Status:      "not_started",
				LastUpdated: "2026-05-20",
				NotesPath:   "notes/d.md",
			},
		},
	}

	stats, err := BuildProgressStats(snapshot, "2026-06-04", reviews)
	if err != nil {
		t.Fatalf("BuildProgressStats() error = %v", err)
	}
	if stats.TrackedLessons != 4 || stats.DoneLessons != 1 || stats.AssistedLessons != 1 || stats.InProgress != 1 || stats.NotStarted != 1 {
		t.Fatalf("unexpected lesson counts: %+v", stats)
	}
	if stats.DueToday != 1 || stats.Overdue != 1 {
		t.Fatalf("unexpected review backlog counts: %+v", stats)
	}
	if stats.CurrentStreak != 2 {
		t.Fatalf("current streak = %d, want 2", stats.CurrentStreak)
	}
	if stats.StrongestArea != "clarification" {
		t.Fatalf("strongest area = %q, want clarification", stats.StrongestArea)
	}
	if stats.WeakestArea != "observability" {
		t.Fatalf("weakest area = %q, want observability", stats.WeakestArea)
	}
}
