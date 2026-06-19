package main

import (
	"math"
	"testing"
)

func approxEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

func TestSkillMatchScore_ExactMatch(t *testing.T) {
	member := Member{Skills: []string{"Go", "Python", "SQL"}}
	job := Job{RequiredSkills: []string{"Go", "Python", "SQL"}}

	score := SkillMatchScore(member, job)
	if !approxEqual(score, 1.0, 0.001) {
		t.Errorf("expected 1.0 for exact match, got %.3f", score)
	}
}

func TestSkillMatchScore_PartialMatch(t *testing.T) {
	member := Member{Skills: []string{"Go", "SQL"}}
	job := Job{RequiredSkills: []string{"Go", "Python", "SQL"}}

	score := SkillMatchScore(member, job)
	// 2 of 3 required skills matched
	if !approxEqual(score, 2.0/3.0, 0.001) {
		t.Errorf("expected %.3f for partial match, got %.3f", 2.0/3.0, score)
	}
}

func TestSkillMatchScore_NoMatch(t *testing.T) {
	member := Member{Skills: []string{"Java", "Spring"}}
	job := Job{RequiredSkills: []string{"Go", "Python"}}

	score := SkillMatchScore(member, job)
	if score != 0.0 {
		t.Errorf("expected 0.0 for no match, got %.3f", score)
	}
}

func TestSkillMatchScore_NoRequirements(t *testing.T) {
	member := Member{Skills: []string{"Go"}}
	job := Job{RequiredSkills: []string{}}

	score := SkillMatchScore(member, job)
	if !approxEqual(score, 1.0, 0.001) {
		t.Errorf("expected 1.0 for job with no requirements, got %.3f", score)
	}
}

func TestSeniorityFit_InRange(t *testing.T) {
	member := Member{SeniorityLevel: 3}
	job := Job{SeniorityMin: 2, SeniorityMax: 4}

	fit := SeniorityFit(member, job)
	if !approxEqual(fit, 1.0, 0.001) {
		t.Errorf("expected 1.0 for in-range seniority, got %.3f", fit)
	}
}

func TestSeniorityFit_AdjacentBelow(t *testing.T) {
	member := Member{SeniorityLevel: 2}
	job := Job{SeniorityMin: 3, SeniorityMax: 4}

	fit := SeniorityFit(member, job)
	if !approxEqual(fit, 0.5, 0.001) {
		t.Errorf("expected 0.5 for one level below minimum, got %.3f", fit)
	}
}

func TestSeniorityFit_AdjacentAbove(t *testing.T) {
	member := Member{SeniorityLevel: 5}
	job := Job{SeniorityMin: 2, SeniorityMax: 4}

	fit := SeniorityFit(member, job)
	if !approxEqual(fit, 0.5, 0.001) {
		t.Errorf("expected 0.5 for one level above maximum, got %.3f", fit)
	}
}

func TestSeniorityFit_TooJunior(t *testing.T) {
	member := Member{SeniorityLevel: 1}
	job := Job{SeniorityMin: 4, SeniorityMax: 5}

	fit := SeniorityFit(member, job)
	if fit != 0.0 {
		t.Errorf("expected 0.0 for too-junior member, got %.3f", fit)
	}
}

func TestMatchScore_EasyApplyBonus(t *testing.T) {
	member := Member{
		Skills:         []string{"Go"},
		LocationID:     "us-ca-sf",
		SeniorityLevel: 3,
	}

	jobWithEasyApply := Job{
		RequiredSkills: []string{"Go"},
		LocationID:     "us-ca-sf",
		SeniorityMin:   3,
		SeniorityMax:   3,
		EasyApply:      true,
	}
	jobWithoutEasyApply := Job{
		RequiredSkills: []string{"Go"},
		LocationID:     "us-ca-sf",
		SeniorityMin:   3,
		SeniorityMax:   3,
		EasyApply:      false,
	}

	withBonus := MatchScore(member, jobWithEasyApply)
	withoutBonus := MatchScore(member, jobWithoutEasyApply)

	if withBonus <= withoutBonus {
		t.Errorf("Easy Apply job should score higher: %.3f vs %.3f", withBonus, withoutBonus)
	}
	if !approxEqual(withBonus-withoutBonus, 0.10, 0.001) {
		t.Errorf("expected Easy Apply bonus of 0.10, got diff %.3f", withBonus-withoutBonus)
	}
}

func TestFindTopJobs_ReturnsCorrectCount(t *testing.T) {
	member := Member{
		Skills:         []string{"Python"},
		LocationID:     "us-ca-sf",
		SeniorityLevel: 3,
	}

	jobs := []Job{
		{ID: "j1", RequiredSkills: []string{"Python"}, LocationID: "us-ca-sf", SeniorityMin: 2, SeniorityMax: 4, EasyApply: false},
		{ID: "j2", RequiredSkills: []string{"Java"}, LocationID: "us-ca-sf", SeniorityMin: 2, SeniorityMax: 4, EasyApply: false},
		{ID: "j3", RequiredSkills: []string{"Python"}, LocationID: "us-ny", SeniorityMin: 3, SeniorityMax: 3, EasyApply: true},
	}

	top := FindTopJobs(member, jobs, 2)
	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
}

func TestFindTopJobs_SeniorityMismatchPenalized(t *testing.T) {
	member := Member{
		Skills:         []string{"Go", "Python"},
		LocationID:     "us-ca-sf",
		SeniorityLevel: 2,
	}

	// goodFit: member is in range, same location, matching skills
	goodFit := Job{
		ID:             "good",
		RequiredSkills: []string{"Go", "Python"},
		LocationID:     "us-ca-sf",
		SeniorityMin:   2,
		SeniorityMax:   3,
		EasyApply:      false,
	}
	// badFit: member is far below required seniority
	badFit := Job{
		ID:             "bad",
		RequiredSkills: []string{"Go", "Python"},
		LocationID:     "us-ca-sf",
		SeniorityMin:   5,
		SeniorityMax:   5,
		EasyApply:      false,
	}

	top := FindTopJobs(member, []Job{badFit, goodFit}, 2)

	if top[0].ID != "good" {
		t.Errorf("expected good-fit job first, got %s", top[0].ID)
	}
}
