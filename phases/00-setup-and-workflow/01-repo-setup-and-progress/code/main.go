package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

var allowedStatuses = []string{"not_started", "in_progress", "done", "assisted"}
var allowedModes = []string{"learn", "practice", "interview"}
var allowedConfidence = []string{"low", "medium", "high"}
var allowedFeedbackSessionTypes = []string{"lesson", "check_understanding", "design_review", "mock_interview"}
var allowedFeedbackDimensions = []string{
	"clarification",
	"requirements",
	"sizing",
	"architecture",
	"deep_dive",
	"failure_modes",
	"observability",
	"trade_offs",
	"communication",
}
var allowedMistakeTags = []string{
	"skips_clarification",
	"weak_requirements",
	"no_sizing",
	"component_soup",
	"bad_deep_dive_choice",
	"weak_consistency_reasoning",
	"shallow_failure_modes",
	"missing_observability",
	"weak_tradeoffs",
	"weak_communication",
	"rushed_architecture",
	"over_indexes_storage",
	"no_operational_story",
	"does_not_tie_back_to_requirements",
}
var mistakeTagToDrill = map[string]string{
	"skips_clarification":               "Run a clarification-only drill and force five scope questions before architecture.",
	"weak_requirements":                 "Rewrite the prompt as prioritized functional and non-functional requirements.",
	"no_sizing":                         "Run a 10-minute capacity drill and defend one number that changes the design.",
	"component_soup":                    "Redraw the design with only major boundaries and justify each component by access pattern.",
	"bad_deep_dive_choice":              "Practice choosing one deep dive explicitly and explain why it matters most.",
	"weak_consistency_reasoning":        "Compare strong vs eventual consistency for one write path and name the user-visible trade-off.",
	"shallow_failure_modes":             "List three likely failures, degraded behavior, and one mitigation for each.",
	"missing_observability":             "Attach one metric, one alert, and one dashboard view to each critical path.",
	"weak_tradeoffs":                    "Do a trade-off drill comparing two valid architectures and name what each sacrifices.",
	"weak_communication":                "Run a compressed summary drill: restate scope, priorities, design, and trade-offs in two minutes.",
	"rushed_architecture":               "Pause after sizing and state the top priorities before naming any datastore or queue.",
	"over_indexes_storage":              "Practice starting from workload shape instead of database choice.",
	"no_operational_story":              "Add rollout, incident response, and failover steps to the design closeout.",
	"does_not_tie_back_to_requirements": "After each design choice, state which requirement it serves and what it hurts.",
}
var mistakeTagToPracticeKind = map[string]string{
	"skips_clarification":               "clarification_drill",
	"weak_requirements":                 "requirements_drill",
	"no_sizing":                         "capacity_drill",
	"component_soup":                    "redesign_drill",
	"bad_deep_dive_choice":              "redesign_drill",
	"weak_consistency_reasoning":        "tradeoff_drill",
	"shallow_failure_modes":             "failure_drill",
	"missing_observability":             "failure_drill",
	"weak_tradeoffs":                    "tradeoff_drill",
	"weak_communication":                "communication_drill",
	"rushed_architecture":               "clarification_drill",
	"over_indexes_storage":              "redesign_drill",
	"no_operational_story":              "failure_drill",
	"does_not_tie_back_to_requirements": "requirements_drill",
}

const currentSchemaVersion = 1
const dateLayout = "2006-01-02"

type LessonProgress struct {
	Lesson             string            `json:"lesson"`
	Status             string            `json:"status"`
	LastUpdated        string            `json:"last_updated"`
	NotesPath          string            `json:"notes_path"`
	ArtifactPaths      []string          `json:"artifact_paths"`
	QuizScore          *int              `json:"quiz_score,omitempty"`
	QuizHistory        []QuizAttempt     `json:"quiz_history,omitempty"`
	Confidence         string            `json:"confidence,omitempty"`
	MistakeTags        []string          `json:"mistake_tags,omitempty"`
	LastReviewedAt     string            `json:"last_reviewed_at,omitempty"`
	NextReviewAt       string            `json:"next_review_at,omitempty"`
	ReviewIntervalDays *int              `json:"review_interval_days,omitempty"`
	ReviewEase         *float64          `json:"review_ease,omitempty"`
	LapseCount         *int              `json:"lapse_count,omitempty"`
	ModeHistory        []SessionRecord   `json:"mode_history,omitempty"`
	FeedbackHistory    []SessionFeedback `json:"feedback_history,omitempty"`
}

type QuizAttempt struct {
	Score       int    `json:"score"`
	CompletedAt string `json:"completed_at"`
	Stage       string `json:"stage,omitempty"`
}

type SessionRecord struct {
	Mode        string `json:"mode"`
	CompletedAt string `json:"completed_at"`
}

type SessionFeedback struct {
	SessionType                string              `json:"session_type"`
	CompletedAt                string              `json:"completed_at"`
	Summary                    string              `json:"summary,omitempty"`
	Strengths                  []string            `json:"strengths,omitempty"`
	Gaps                       []string            `json:"gaps,omitempty"`
	HighestLeverageImprovement string              `json:"highest_leverage_improvement,omitempty"`
	Dimensions                 []DimensionFeedback `json:"dimensions,omitempty"`
}

type DimensionFeedback struct {
	Dimension  string `json:"dimension"`
	Score      int    `json:"score"`
	Evidence   string `json:"evidence,omitempty"`
	NextAction string `json:"next_action,omitempty"`
}

type ProgressSnapshot struct {
	SchemaVersion int              `json:"schema_version,omitempty"`
	Lessons       []LessonProgress `json:"lessons"`
}

type ReviewItem struct {
	Lesson     string
	DueAt      string
	State      string
	DaysLate   int
	Confidence string
	Priority   int
	TopMistake string
}

type Recommendation struct {
	Action      string
	Lesson      string
	Reason      string
	DrillKind   string
	DrillPrompt string
}

type DimensionTrend struct {
	Dimension string
	Delta     float64
}

type MistakeTrend struct {
	Tag   string
	Count int
	Drill string
}

type ProgressStats struct {
	TrackedLessons   int
	DoneLessons      int
	AssistedLessons  int
	InProgress       int
	NotStarted       int
	CompletionRate   float64
	AssistedRate     float64
	RecentQuizAvg    float64
	ActiveDays7d     int
	CurrentStreak    int
	DueToday         int
	Overdue          int
	StrongestArea    string
	StrongestScore   float64
	WeakestArea      string
	WeakestScore     float64
	WeeklySessions   int
	WeeklyLessons    int
	WeeklyReviewsDue int
}

func ValidateSnapshot(snapshot ProgressSnapshot) []string {
	var issues []string
	seen := map[string]bool{}

	if snapshot.SchemaVersion < 0 {
		issues = append(issues, "schema_version cannot be negative")
	}
	if snapshot.SchemaVersion > currentSchemaVersion {
		issues = append(
			issues,
			fmt.Sprintf(
				"schema_version %d is newer than supported version %d; migrate the validator before using this file",
				snapshot.SchemaVersion,
				currentSchemaVersion,
			),
		)
	}

	for _, lesson := range snapshot.Lessons {
		if lesson.Lesson == "" {
			issues = append(issues, "lesson slug is required")
		}
		if lesson.Lesson != "" {
			if seen[lesson.Lesson] {
				issues = append(issues, fmt.Sprintf("duplicate lesson entry %q", lesson.Lesson))
			}
			seen[lesson.Lesson] = true
		}
		if !slices.Contains(allowedStatuses, lesson.Status) {
			issues = append(issues, fmt.Sprintf("lesson %q uses invalid status %q", lesson.Lesson, lesson.Status))
		}
		if lesson.LastUpdated == "" {
			issues = append(issues, fmt.Sprintf("lesson %q is missing last_updated", lesson.Lesson))
		}
		if lesson.NotesPath == "" {
			issues = append(issues, fmt.Sprintf("lesson %q is missing notes_path", lesson.Lesson))
		}
		if lesson.QuizScore != nil && (*lesson.QuizScore < 0 || *lesson.QuizScore > 8) {
			issues = append(issues, fmt.Sprintf("lesson %q has quiz_score %d outside 0..8", lesson.Lesson, *lesson.QuizScore))
		}
		if lesson.Confidence != "" && !slices.Contains(allowedConfidence, lesson.Confidence) {
			issues = append(issues, fmt.Sprintf("lesson %q uses invalid confidence %q", lesson.Lesson, lesson.Confidence))
		}
		if lesson.ReviewIntervalDays != nil && *lesson.ReviewIntervalDays < 0 {
			issues = append(issues, fmt.Sprintf("lesson %q has negative review_interval_days", lesson.Lesson))
		}
		if lesson.ReviewEase != nil && *lesson.ReviewEase < 1.0 {
			issues = append(issues, fmt.Sprintf("lesson %q has review_ease %.2f below 1.0", lesson.Lesson, *lesson.ReviewEase))
		}
		if lesson.LapseCount != nil && *lesson.LapseCount < 0 {
			issues = append(issues, fmt.Sprintf("lesson %q has negative lapse_count", lesson.Lesson))
		}
		if lesson.NextReviewAt != "" {
			if _, err := time.Parse(dateLayout, lesson.NextReviewAt); err != nil {
				issues = append(issues, fmt.Sprintf("lesson %q has invalid next_review_at %q", lesson.Lesson, lesson.NextReviewAt))
			}
		}
		if lesson.LastReviewedAt != "" {
			if _, err := time.Parse(dateLayout, lesson.LastReviewedAt); err != nil {
				issues = append(issues, fmt.Sprintf("lesson %q has invalid last_reviewed_at %q", lesson.Lesson, lesson.LastReviewedAt))
			}
		}
		for _, tag := range lesson.MistakeTags {
			if strings.TrimSpace(tag) == "" {
				issues = append(issues, fmt.Sprintf("lesson %q has an empty mistake tag", lesson.Lesson))
				continue
			}
			if !slices.Contains(allowedMistakeTags, tag) {
				issues = append(issues, fmt.Sprintf("lesson %q uses invalid mistake tag %q", lesson.Lesson, tag))
			}
		}
		for _, attempt := range lesson.QuizHistory {
			if attempt.Score < 0 || attempt.Score > 8 {
				issues = append(issues, fmt.Sprintf("lesson %q has quiz_history score %d outside 0..8", lesson.Lesson, attempt.Score))
			}
			if attempt.CompletedAt == "" {
				issues = append(issues, fmt.Sprintf("lesson %q has quiz_history entry missing completed_at", lesson.Lesson))
			}
		}
		for _, session := range lesson.ModeHistory {
			if !slices.Contains(allowedModes, session.Mode) {
				issues = append(issues, fmt.Sprintf("lesson %q has invalid mode_history mode %q", lesson.Lesson, session.Mode))
			}
			if session.CompletedAt == "" {
				issues = append(issues, fmt.Sprintf("lesson %q has mode_history entry missing completed_at", lesson.Lesson))
			}
		}
		for _, feedback := range lesson.FeedbackHistory {
			if !slices.Contains(allowedFeedbackSessionTypes, feedback.SessionType) {
				issues = append(issues, fmt.Sprintf("lesson %q has invalid feedback_history session_type %q", lesson.Lesson, feedback.SessionType))
			}
			if feedback.CompletedAt == "" {
				issues = append(issues, fmt.Sprintf("lesson %q has feedback_history entry missing completed_at", lesson.Lesson))
			}
			if strings.TrimSpace(feedback.HighestLeverageImprovement) == "" {
				issues = append(issues, fmt.Sprintf("lesson %q has feedback_history entry missing highest_leverage_improvement", lesson.Lesson))
			}
			for _, strength := range feedback.Strengths {
				if strings.TrimSpace(strength) == "" {
					issues = append(issues, fmt.Sprintf("lesson %q has feedback_history entry with empty strength", lesson.Lesson))
				}
			}
			for _, gap := range feedback.Gaps {
				if strings.TrimSpace(gap) == "" {
					issues = append(issues, fmt.Sprintf("lesson %q has feedback_history entry with empty gap", lesson.Lesson))
				}
			}
			for _, dimension := range feedback.Dimensions {
				if !slices.Contains(allowedFeedbackDimensions, dimension.Dimension) {
					issues = append(issues, fmt.Sprintf("lesson %q has invalid feedback dimension %q", lesson.Lesson, dimension.Dimension))
				}
				if dimension.Score < 1 || dimension.Score > 4 {
					issues = append(issues, fmt.Sprintf("lesson %q has feedback dimension %q with score %d outside 1..4", lesson.Lesson, dimension.Dimension, dimension.Score))
				}
			}
		}
	}

	return issues
}

func RecommendNext(snapshot ProgressSnapshot) string {
	for _, status := range []string{"not_started", "in_progress", "assisted"} {
		for _, lesson := range snapshot.Lessons {
			if lesson.Status == status {
				return lesson.Lesson
			}
		}
	}
	return ""
}

func RecommendReviews(snapshot ProgressSnapshot, today string) ([]ReviewItem, error) {
	referenceDay, err := time.Parse(dateLayout, today)
	if err != nil {
		return nil, fmt.Errorf("parse today: %w", err)
	}

	var items []ReviewItem
	recurringMistakes := recurringMistakeCounts(snapshot)
	for _, lesson := range snapshot.Lessons {
		if lesson.NextReviewAt == "" {
			continue
		}
		dueAt, err := time.Parse(dateLayout, lesson.NextReviewAt)
		if err != nil {
			return nil, fmt.Errorf("parse next_review_at for %s: %w", lesson.Lesson, err)
		}
		days := int(referenceDay.Sub(dueAt).Hours() / 24)
		state := "upcoming"
		if days > 0 {
			state = "overdue"
		} else if days == 0 {
			state = "due_today"
		}
		if state == "upcoming" {
			continue
		}
		topMistake, priorityBoost := lessonMistakePriority(lesson, recurringMistakes)
		items = append(items, ReviewItem{
			Lesson:     lesson.Lesson,
			DueAt:      lesson.NextReviewAt,
			State:      state,
			DaysLate:   max(days, 0),
			Confidence: lesson.Confidence,
			Priority:   priorityBoost,
			TopMistake: topMistake,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].State != items[j].State {
			return items[i].State == "overdue"
		}
		if items[i].Priority != items[j].Priority {
			return items[i].Priority > items[j].Priority
		}
		if items[i].DaysLate != items[j].DaysLate {
			return items[i].DaysLate > items[j].DaysLate
		}
		return items[i].Lesson < items[j].Lesson
	})

	return items, nil
}

func RecommendAction(snapshot ProgressSnapshot, today string) (Recommendation, error) {
	reviews, err := RecommendReviews(snapshot, today)
	if err != nil {
		return Recommendation{}, err
	}
	if len(reviews) > 0 {
		top := reviews[0]
		reason := "review is due"
		if top.State == "overdue" {
			reason = fmt.Sprintf("review is overdue by %d day(s)", top.DaysLate)
		}
		if top.TopMistake != "" {
			reason = fmt.Sprintf("%s and recurring mistake %s is still active", reason, top.TopMistake)
			if drillKind, drillPrompt := practiceDrillForLesson(findLesson(snapshot, top.Lesson)); drillKind != "" {
				return Recommendation{
					Action:      "practice",
					Lesson:      top.Lesson,
					Reason:      reason,
					DrillKind:   drillKind,
					DrillPrompt: drillPrompt,
				}, nil
			}
		}
		return Recommendation{
			Action: "review",
			Lesson: top.Lesson,
			Reason: reason,
		}, nil
	}

	if trends, err := WeeklyMistakeSummary(snapshot, today); err == nil && len(trends) > 0 && trends[0].Count >= 2 {
		return Recommendation{
			Action:      "practice",
			Reason:      fmt.Sprintf("recurring mistake %s appeared %d times this week", trends[0].Tag, trends[0].Count),
			DrillKind:   mistakeTagToPracticeKind[trends[0].Tag],
			DrillPrompt: trends[0].Drill,
		}, nil
	}

	for _, lesson := range snapshot.Lessons {
		if lesson.Status == "assisted" {
			return Recommendation{
				Action: "revisit",
				Lesson: lesson.Lesson,
				Reason: "last attempt was assisted and should be retried unassisted",
			}, nil
		}
	}

	for _, lesson := range snapshot.Lessons {
		if lesson.Status == "in_progress" {
			if weak := weakestDimensionForLesson(lesson); weak != "" {
				return Recommendation{
					Action: "resume",
					Lesson: lesson.Lesson,
					Reason: fmt.Sprintf("session is in progress and still weak on %s", weak),
				}, nil
			}
			return Recommendation{
				Action: "resume",
				Lesson: lesson.Lesson,
				Reason: "session is already in progress",
			}, nil
		}
	}

	next := RecommendNext(snapshot)
	if next != "" {
		return Recommendation{
			Action: "start",
			Lesson: next,
			Reason: "next unfinished lesson in roadmap order",
		}, nil
	}

	return Recommendation{
		Action: "idle",
		Reason: "no due reviews and no unfinished lessons found",
	}, nil
}

func WeakestTrend(snapshot ProgressSnapshot) DimensionTrend {
	scores := averageDimensionScores(snapshot)
	weakest := DimensionTrend{}
	found := false
	for dimension, score := range scores {
		if !found || score < weakest.Delta {
			weakest = DimensionTrend{Dimension: dimension, Delta: score}
			found = true
		}
	}
	return weakest
}

func ImprovingTrend(snapshot ProgressSnapshot) DimensionTrend {
	deltas := map[string]float64{}
	counts := map[string]int{}

	for _, lesson := range snapshot.Lessons {
		if len(lesson.FeedbackHistory) < 2 {
			continue
		}
		prev := feedbackDimensionMap(lesson.FeedbackHistory[len(lesson.FeedbackHistory)-2])
		curr := feedbackDimensionMap(lesson.FeedbackHistory[len(lesson.FeedbackHistory)-1])
		for dimension, currentScore := range curr {
			previousScore, ok := prev[dimension]
			if !ok {
				continue
			}
			deltas[dimension] += float64(currentScore - previousScore)
			counts[dimension]++
		}
	}

	best := DimensionTrend{}
	found := false
	for dimension, total := range deltas {
		avg := total / float64(counts[dimension])
		if avg <= 0 {
			continue
		}
		if !found || avg > best.Delta {
			best = DimensionTrend{Dimension: dimension, Delta: avg}
			found = true
		}
	}
	return best
}

func WeeklyMistakeSummary(snapshot ProgressSnapshot, today string) ([]MistakeTrend, error) {
	referenceDay, err := time.Parse(dateLayout, today)
	if err != nil {
		return nil, fmt.Errorf("parse today: %w", err)
	}

	counts := map[string]int{}
	for _, lesson := range snapshot.Lessons {
		sessionDate, ok := latestSessionDate(lesson)
		if !ok {
			continue
		}
		if referenceDay.Sub(sessionDate).Hours() > 24*7 {
			continue
		}
		for _, tag := range lesson.MistakeTags {
			counts[tag]++
		}
	}

	trends := make([]MistakeTrend, 0, len(counts))
	for tag, count := range counts {
		trends = append(trends, MistakeTrend{
			Tag:   tag,
			Count: count,
			Drill: mistakeTagToDrill[tag],
		})
	}

	sort.Slice(trends, func(i, j int) bool {
		if trends[i].Count != trends[j].Count {
			return trends[i].Count > trends[j].Count
		}
		return trends[i].Tag < trends[j].Tag
	})

	return trends, nil
}

func BuildProgressStats(snapshot ProgressSnapshot, today string, reviews []ReviewItem) (ProgressStats, error) {
	referenceDay, err := time.Parse(dateLayout, today)
	if err != nil {
		return ProgressStats{}, fmt.Errorf("parse today: %w", err)
	}

	stats := ProgressStats{
		TrackedLessons: len(snapshot.Lessons),
	}
	activeDaySet := map[string]bool{}
	weeklyLessonSet := map[string]bool{}

	for _, lesson := range snapshot.Lessons {
		switch lesson.Status {
		case "done":
			stats.DoneLessons++
		case "assisted":
			stats.AssistedLessons++
		case "in_progress":
			stats.InProgress++
		case "not_started":
			stats.NotStarted++
		}

		if score, ok := latestQuizScore(lesson); ok {
			stats.RecentQuizAvg += float64(score)
		}

		if sessionDate, ok := latestSessionDate(lesson); ok {
			dayKey := sessionDate.Format(dateLayout)
			activeDaySet[dayKey] = true
			if referenceDay.Sub(sessionDate).Hours() <= 24*7 {
				weeklyLessonSet[lesson.Lesson] = true
				stats.WeeklySessions++
			}
		}
	}

	quizCount := 0
	for _, lesson := range snapshot.Lessons {
		if _, ok := latestQuizScore(lesson); ok {
			quizCount++
		}
	}
	if stats.TrackedLessons > 0 {
		stats.CompletionRate = float64(stats.DoneLessons+stats.AssistedLessons) / float64(stats.TrackedLessons)
	}
	completedAttempts := stats.DoneLessons + stats.AssistedLessons
	if completedAttempts > 0 {
		stats.AssistedRate = float64(stats.AssistedLessons) / float64(completedAttempts)
	}
	if quizCount > 0 {
		stats.RecentQuizAvg = stats.RecentQuizAvg / float64(quizCount)
	}
	stats.ActiveDays7d = countActiveDaysSince(activeDaySet, referenceDay, 7)
	stats.CurrentStreak = currentStreak(activeDaySet, referenceDay)
	stats.WeeklyLessons = len(weeklyLessonSet)

	for _, review := range reviews {
		if review.State == "due_today" {
			stats.DueToday++
			stats.WeeklyReviewsDue++
		}
		if review.State == "overdue" {
			stats.Overdue++
			stats.WeeklyReviewsDue++
		}
	}

	averages := averageDimensionScores(snapshot)
	for dimension, score := range averages {
		if stats.StrongestArea == "" || score > stats.StrongestScore {
			stats.StrongestArea = dimension
			stats.StrongestScore = score
		}
		if stats.WeakestArea == "" || score < stats.WeakestScore {
			stats.WeakestArea = dimension
			stats.WeakestScore = score
		}
	}

	return stats, nil
}

func practiceDrillForLesson(lesson *LessonProgress) (string, string) {
	if lesson == nil {
		return "", ""
	}
	for _, tag := range lesson.MistakeTags {
		kind := mistakeTagToPracticeKind[tag]
		prompt := mistakeTagToDrill[tag]
		if kind != "" && prompt != "" {
			return kind, prompt
		}
	}

	weakest := weakestDimensionForLesson(*lesson)
	switch weakest {
	case "sizing":
		return "capacity_drill", "Run a 10-minute estimate for QPS, storage, bandwidth, and cost before touching architecture."
	case "trade_offs":
		return "tradeoff_drill", "Compare two viable designs and state what each sacrifices."
	case "failure_modes", "observability":
		return "failure_drill", "Name three likely failures, one signal for each, and how the system degrades."
	case "clarification", "requirements":
		return "clarification_drill", "Force five scope and priority questions before designing."
	default:
		return "", ""
	}
}

func SyncReviewSchedule(snapshot *ProgressSnapshot, today string) error {
	referenceDay, err := time.Parse(dateLayout, today)
	if err != nil {
		return fmt.Errorf("parse today: %w", err)
	}

	for i := range snapshot.Lessons {
		updateReviewSchedule(&snapshot.Lessons[i], referenceDay)
	}

	return nil
}

func updateReviewSchedule(lesson *LessonProgress, reviewedAt time.Time) {
	if lesson.Status == "not_started" {
		return
	}

	quality := deriveReviewQuality(*lesson)
	if quality == "" {
		return
	}

	currentInterval := derefInt(lesson.ReviewIntervalDays, 0)
	currentEase := derefFloat(lesson.ReviewEase, 2.3)
	lapses := derefInt(lesson.LapseCount, 0)
	recurringPenalty := lessonRecurringMistakePenalty(*lesson)

	nextInterval := nextReviewIntervalDays(quality, currentInterval, currentEase)
	nextEase := nextReviewEase(quality, currentEase)
	if quality == "again" {
		lapses++
	}
	if recurringPenalty > 0 {
		nextInterval = max(1, nextInterval-recurringPenalty)
	}

	lesson.LastReviewedAt = reviewedAt.Format(dateLayout)
	lesson.NextReviewAt = reviewedAt.AddDate(0, 0, nextInterval).Format(dateLayout)
	lesson.ReviewIntervalDays = intPtr(nextInterval)
	lesson.ReviewEase = floatPtr(nextEase)
	lesson.LapseCount = intPtr(lapses)
}

func deriveReviewQuality(lesson LessonProgress) string {
	weakDimensions := countWeakDimensions(lesson)

	if lesson.Status == "assisted" {
		return "again"
	}

	score, hasScore := latestQuizScore(lesson)
	switch {
	case hasScore && score <= 5:
		return "again"
	case hasScore && score <= 7:
		if weakDimensions > 0 {
			return "hard"
		}
		return "good"
	case hasScore && score == 8:
		if weakDimensions > 1 || lesson.Confidence == "low" {
			return "hard"
		}
		if lesson.Confidence == "high" && weakDimensions == 0 {
			return "easy"
		}
		return "good"
	case weakDimensions > 0:
		return "hard"
	case lesson.Status == "done":
		return "good"
	default:
		return ""
	}
}

func latestQuizScore(lesson LessonProgress) (int, bool) {
	if len(lesson.QuizHistory) > 0 {
		return lesson.QuizHistory[len(lesson.QuizHistory)-1].Score, true
	}
	if lesson.QuizScore != nil {
		return *lesson.QuizScore, true
	}
	return 0, false
}

func countWeakDimensions(lesson LessonProgress) int {
	if len(lesson.FeedbackHistory) == 0 {
		return 0
	}
	latest := lesson.FeedbackHistory[len(lesson.FeedbackHistory)-1]
	count := 0
	for _, dimension := range latest.Dimensions {
		if dimension.Score < 3 {
			count++
		}
	}
	return count
}

func weakestDimensionForLesson(lesson LessonProgress) string {
	if len(lesson.FeedbackHistory) == 0 {
		if len(lesson.MistakeTags) > 0 {
			return lesson.MistakeTags[0]
		}
		return ""
	}
	latest := lesson.FeedbackHistory[len(lesson.FeedbackHistory)-1]
	lowestScore := 5
	weakest := ""
	for _, dimension := range latest.Dimensions {
		if dimension.Score < lowestScore {
			lowestScore = dimension.Score
			weakest = dimension.Dimension
		}
	}
	if weakest == "" && len(lesson.MistakeTags) > 0 {
		return lesson.MistakeTags[0]
	}
	return weakest
}

func recurringMistakeCounts(snapshot ProgressSnapshot) map[string]int {
	counts := map[string]int{}
	for _, lesson := range snapshot.Lessons {
		seen := map[string]bool{}
		for _, tag := range lesson.MistakeTags {
			if seen[tag] {
				continue
			}
			counts[tag]++
			seen[tag] = true
		}
	}
	return counts
}

func lessonMistakePriority(lesson LessonProgress, recurring map[string]int) (string, int) {
	topTag := ""
	topCount := 0
	for _, tag := range lesson.MistakeTags {
		if recurring[tag] > topCount {
			topTag = tag
			topCount = recurring[tag]
		}
	}
	if topCount <= 1 {
		return topTag, 0
	}
	return topTag, topCount - 1
}

func lessonRecurringMistakePenalty(lesson LessonProgress) int {
	penalty := 0
	seen := map[string]bool{}
	for _, tag := range lesson.MistakeTags {
		if seen[tag] {
			continue
		}
		if strings.TrimSpace(tag) != "" {
			penalty++
		}
		seen[tag] = true
	}
	if penalty > 2 {
		return 2
	}
	if penalty > 0 {
		return 1
	}
	return 0
}

func findLesson(snapshot ProgressSnapshot, slug string) *LessonProgress {
	for i := range snapshot.Lessons {
		if snapshot.Lessons[i].Lesson == slug {
			return &snapshot.Lessons[i]
		}
	}
	return nil
}

func averageDimensionScores(snapshot ProgressSnapshot) map[string]float64 {
	total := map[string]float64{}
	count := map[string]int{}
	for _, lesson := range snapshot.Lessons {
		if len(lesson.FeedbackHistory) == 0 {
			continue
		}
		latest := lesson.FeedbackHistory[len(lesson.FeedbackHistory)-1]
		for _, dimension := range latest.Dimensions {
			total[dimension.Dimension] += float64(dimension.Score)
			count[dimension.Dimension]++
		}
	}

	averages := map[string]float64{}
	for dimension, sum := range total {
		averages[dimension] = sum / float64(count[dimension])
	}
	return averages
}

func feedbackDimensionMap(feedback SessionFeedback) map[string]int {
	values := map[string]int{}
	for _, dimension := range feedback.Dimensions {
		values[dimension.Dimension] = dimension.Score
	}
	return values
}

func latestSessionDate(lesson LessonProgress) (time.Time, bool) {
	if len(lesson.FeedbackHistory) > 0 {
		last := lesson.FeedbackHistory[len(lesson.FeedbackHistory)-1]
		if last.CompletedAt != "" {
			value, err := time.Parse(dateLayout, last.CompletedAt)
			if err == nil {
				return value, true
			}
		}
	}
	if lesson.LastUpdated != "" {
		value, err := time.Parse(dateLayout, lesson.LastUpdated)
		if err == nil {
			return value, true
		}
	}
	return time.Time{}, false
}

func countActiveDaysSince(activeDays map[string]bool, referenceDay time.Time, windowDays int) int {
	count := 0
	for i := 0; i < windowDays; i++ {
		day := referenceDay.AddDate(0, 0, -i).Format(dateLayout)
		if activeDays[day] {
			count++
		}
	}
	return count
}

func currentStreak(activeDays map[string]bool, referenceDay time.Time) int {
	streak := 0
	for i := 0; ; i++ {
		day := referenceDay.AddDate(0, 0, -i).Format(dateLayout)
		if !activeDays[day] {
			break
		}
		streak++
	}
	return streak
}

func nextReviewIntervalDays(quality string, currentInterval int, currentEase float64) int {
	if currentInterval <= 0 {
		switch quality {
		case "again":
			return 1
		case "hard":
			return 3
		case "good":
			return 7
		case "easy":
			return 14
		default:
			return 0
		}
	}

	switch quality {
	case "again":
		return 1
	case "hard":
		return max(3, int(float64(currentInterval)*1.2+0.5))
	case "good":
		return max(7, int(float64(currentInterval)*currentEase+0.5))
	case "easy":
		return max(14, int(float64(currentInterval)*(currentEase+0.3)+0.5))
	default:
		return currentInterval
	}
}

func nextReviewEase(quality string, currentEase float64) float64 {
	switch quality {
	case "again":
		return maxFloat(1.3, currentEase-0.2)
	case "hard":
		return maxFloat(1.5, currentEase-0.05)
	case "good":
		return currentEase
	case "easy":
		return currentEase + 0.1
	default:
		return currentEase
	}
}

func derefInt(value *int, fallback int) int {
	if value == nil {
		return fallback
	}
	return *value
}

func derefFloat(value *float64, fallback float64) float64 {
	if value == nil {
		return fallback
	}
	return *value
}

func intPtr(value int) *int {
	return &value
}

func floatPtr(value float64) *float64 {
	return &value
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	var path string
	var today string
	var syncReviews bool
	flag.StringVar(&path, "progress", "", "path to a progress snapshot JSON file")
	flag.StringVar(&today, "today", time.Now().Format(dateLayout), "reference date in YYYY-MM-DD format")
	flag.BoolVar(&syncReviews, "sync-reviews", false, "update review metadata in the progress snapshot before reporting")
	flag.Parse()

	if path == "" {
		fmt.Fprintln(os.Stderr, "missing -progress path")
		os.Exit(2)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read snapshot: %v\n", err)
		os.Exit(1)
	}

	var snapshot ProgressSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		fmt.Fprintf(os.Stderr, "decode snapshot: %v\n", err)
		os.Exit(1)
	}

	issues := ValidateSnapshot(snapshot)
	if len(issues) > 0 {
		for _, issue := range issues {
			fmt.Println(issue)
		}
		os.Exit(1)
	}

	if syncReviews {
		if err := SyncReviewSchedule(&snapshot, today); err != nil {
			fmt.Fprintf(os.Stderr, "sync review schedule: %v\n", err)
			os.Exit(1)
		}

		encoded, err := json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "encode snapshot: %v\n", err)
			os.Exit(1)
		}
		encoded = append(encoded, '\n')
		if err := os.WriteFile(path, encoded, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write snapshot: %v\n", err)
			os.Exit(1)
		}
	}

	reviews, err := RecommendReviews(snapshot, today)
	if err != nil {
		fmt.Fprintf(os.Stderr, "recommend reviews: %v\n", err)
		os.Exit(1)
	}
	recommendation, err := RecommendAction(snapshot, today)
	if err != nil {
		fmt.Fprintf(os.Stderr, "recommend action: %v\n", err)
		os.Exit(1)
	}
	weakTrend := WeakestTrend(snapshot)
	improvingTrend := ImprovingTrend(snapshot)
	mistakeTrends, err := WeeklyMistakeSummary(snapshot, today)
	if err != nil {
		fmt.Fprintf(os.Stderr, "weekly mistake summary: %v\n", err)
		os.Exit(1)
	}
	stats, err := BuildProgressStats(snapshot, today, reviews)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build progress stats: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("tracked lessons: %d\n", stats.TrackedLessons)
	fmt.Printf("completion: %d done, %d assisted, %d in progress, %d not started (%.0f%% complete)\n",
		stats.DoneLessons,
		stats.AssistedLessons,
		stats.InProgress,
		stats.NotStarted,
		stats.CompletionRate*100,
	)
	fmt.Printf("assisted ratio: %.0f%%\n", stats.AssistedRate*100)
	if stats.RecentQuizAvg > 0 {
		fmt.Printf("recent quiz average: %.2f/8\n", stats.RecentQuizAvg)
	}
	fmt.Printf("consistency: %d active day(s) in last 7, current streak %d day(s)\n", stats.ActiveDays7d, stats.CurrentStreak)
	fmt.Printf("weekly summary: %d session(s), %d lesson(s) touched, %d review(s) due or overdue\n",
		stats.WeeklySessions,
		stats.WeeklyLessons,
		stats.WeeklyReviewsDue,
	)
	if stats.StrongestArea != "" {
		fmt.Printf("strongest area: %s (avg %.2f/4)\n", stats.StrongestArea, stats.StrongestScore)
	}
	fmt.Printf("next lesson: %s\n", RecommendNext(snapshot))
	fmt.Printf("due reviews: %d\n", len(reviews))
	fmt.Printf("review backlog: %d due today, %d overdue\n", stats.DueToday, stats.Overdue)
	for _, review := range reviews {
		suffix := ""
		if review.State == "overdue" {
			suffix = " (" + strconv.Itoa(review.DaysLate) + "d late)"
		}
		if review.Confidence != "" {
			suffix += " [" + review.Confidence + "]"
		}
		if review.TopMistake != "" {
			suffix += " {mistake:" + review.TopMistake + "}"
		}
		fmt.Printf("- %s: %s on %s%s\n", review.Lesson, review.State, review.DueAt, suffix)
	}
	fmt.Printf("recommended action: %s %s\n", recommendation.Action, recommendation.Lesson)
	fmt.Printf("recommendation reason: %s\n", recommendation.Reason)
	if recommendation.DrillKind != "" {
		fmt.Printf("recommended drill: %s\n", recommendation.DrillKind)
	}
	if recommendation.DrillPrompt != "" {
		fmt.Printf("drill prompt: %s\n", recommendation.DrillPrompt)
	}
	if weakTrend.Dimension != "" {
		fmt.Printf("weak trend: %s (avg %.2f/4)\n", weakTrend.Dimension, weakTrend.Delta)
	}
	if improvingTrend.Dimension != "" {
		fmt.Printf("improving trend: %s (+%.2f)\n", improvingTrend.Dimension, improvingTrend.Delta)
	}
	if len(mistakeTrends) > 0 {
		top := mistakeTrends[0]
		fmt.Printf("weekly mistake trend: %s (%dx)\n", top.Tag, top.Count)
		if top.Drill != "" {
			fmt.Printf("suggested drill: %s\n", top.Drill)
		}
	}
}
