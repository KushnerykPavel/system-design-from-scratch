package main

import (
	"fmt"
	"sort"
)

// Member represents a LinkedIn member's profile attributes used for job matching.
type Member struct {
	ID             string
	Skills         []string
	Title          string
	LocationID     string
	SeniorityLevel int // 1=Entry, 2=Associate, 3=Mid, 4=Senior, 5=Staff/Principal
}

// Job represents an active job listing.
type Job struct {
	ID             string
	Title          string
	LocationID     string
	RequiredSkills []string
	SeniorityMin   int // minimum seniority level accepted
	SeniorityMax   int // maximum seniority level accepted
	EasyApply      bool
}

// SkillMatchScore returns the fraction of the job's required skills that the
// member possesses. Score = |member.Skills ∩ job.RequiredSkills| / |job.RequiredSkills|.
// Returns 0 if the job has no required skills.
func SkillMatchScore(member Member, job Job) float64 {
	if len(job.RequiredSkills) == 0 {
		return 1.0 // no requirements means everyone qualifies
	}

	memberSkillSet := make(map[string]bool, len(member.Skills))
	for _, s := range member.Skills {
		memberSkillSet[s] = true
	}

	matched := 0
	for _, required := range job.RequiredSkills {
		if memberSkillSet[required] {
			matched++
		}
	}

	return float64(matched) / float64(len(job.RequiredSkills))
}

// SeniorityFit returns how well the member's seniority level fits the job's
// accepted range.
//   - 1.0  if member.SeniorityLevel is within [job.SeniorityMin, job.SeniorityMax]
//   - 0.5  if member.SeniorityLevel is exactly one level below SeniorityMin or
//     one level above SeniorityMax (adjacent)
//   - 0.0  otherwise (too far from the accepted range)
func SeniorityFit(member Member, job Job) float64 {
	level := member.SeniorityLevel
	switch {
	case level >= job.SeniorityMin && level <= job.SeniorityMax:
		return 1.0
	case level == job.SeniorityMin-1 || level == job.SeniorityMax+1:
		return 0.5
	default:
		return 0.0
	}
}

// MatchScore computes a weighted relevance score for a (member, job) pair.
// Weights:
//   - Skills intersection:   0.45
//   - Seniority fit:         0.30
//   - Same location:         0.15 (binary: same LocationID)
//   - Easy Apply bonus:      +0.10 added to the final score (not a multiplier)
func MatchScore(member Member, job Job) float64 {
	skillScore := SkillMatchScore(member, job)
	seniorityScore := SeniorityFit(member, job)

	locationScore := 0.0
	if member.LocationID == job.LocationID {
		locationScore = 1.0
	}

	score := skillScore*0.45 + seniorityScore*0.30 + locationScore*0.15

	if job.EasyApply {
		score += 0.10
	}

	return score
}

// FindTopJobs returns the top-N jobs for the given member, sorted by MatchScore
// descending. If topN is larger than the number of jobs, all jobs are returned.
func FindTopJobs(member Member, jobs []Job, topN int) []Job {
	type scored struct {
		job   Job
		score float64
	}

	results := make([]scored, len(jobs))
	for i, j := range jobs {
		results[i] = scored{job: j, score: MatchScore(member, j)}
	}

	sort.Slice(results, func(i, k int) bool {
		if results[i].score == results[k].score {
			return results[i].job.ID < results[k].job.ID
		}
		return results[i].score > results[k].score
	})

	if topN > len(results) {
		topN = len(results)
	}

	top := make([]Job, topN)
	for i := 0; i < topN; i++ {
		top[i] = results[i].job
	}
	return top
}

func main() {
	member := Member{
		ID:             "member-1",
		Skills:         []string{"Go", "Python", "Kubernetes", "SQL"},
		Title:          "Software Engineer",
		LocationID:     "us-ca-sf",
		SeniorityLevel: 3,
	}

	jobs := []Job{
		{
			ID:             "job-001",
			Title:          "Senior Backend Engineer",
			LocationID:     "us-ca-sf",
			RequiredSkills: []string{"Go", "Kubernetes", "SQL"},
			SeniorityMin:   3,
			SeniorityMax:   4,
			EasyApply:      true,
		},
		{
			ID:             "job-002",
			Title:          "ML Engineer",
			LocationID:     "us-ca-sf",
			RequiredSkills: []string{"Python", "PyTorch", "Kubernetes"},
			SeniorityMin:   3,
			SeniorityMax:   5,
			EasyApply:      false,
		},
		{
			ID:             "job-003",
			Title:          "Principal Architect",
			LocationID:     "us-ny",
			RequiredSkills: []string{"Go", "Python", "Kubernetes", "SQL", "Terraform"},
			SeniorityMin:   5,
			SeniorityMax:   5,
			EasyApply:      false,
		},
		{
			ID:             "job-004",
			Title:          "Junior Data Analyst",
			LocationID:     "us-ca-sf",
			RequiredSkills: []string{"SQL", "Excel"},
			SeniorityMin:   1,
			SeniorityMax:   2,
			EasyApply:      true,
		},
	}

	top := FindTopJobs(member, jobs, 3)

	fmt.Printf("Top jobs for %s (seniority=%d, location=%s):\n", member.ID, member.SeniorityLevel, member.LocationID)
	for rank, job := range top {
		score := MatchScore(member, job)
		fmt.Printf("  %d. %s [%s] score=%.3f easyApply=%v\n", rank+1, job.Title, job.ID, score, job.EasyApply)
	}
}
