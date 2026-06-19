package main

import (
	"fmt"
	"time"
)

// SenderProfile captures attributes used to score InMail spam risk and
// compute sending eligibility.
type SenderProfile struct {
	AccountAgeDays  int     // how many days since account was created
	ResponseRate    float64 // fraction of sent InMails that received a reply (0..1)
	ReportRate      float64 // fraction of sent InMails that were reported as spam (0..1)
	DailySentCount  int     // InMails already sent today (before this one)
	IsRecruiter     bool    // paid recruiter seat
}

// SpamScore returns a score in [0, 1] where 0 = clean and 1 = likely spam.
// The score is a weighted combination of risk signals:
//   - Low response rate increases risk (0.35 weight)
//   - High report rate increases risk (0.40 weight)
//   - New account increases risk (0.25 weight, linearly decays over 90 days)
func SpamScore(profile SenderProfile) float64 {
	// Response rate contribution: low response rate → high risk.
	// Invert so that 0% response rate contributes the full weight.
	responseRisk := (1.0 - profile.ResponseRate) * 0.35

	// Report rate contribution: 1% report rate is already very high.
	// Clamp to 1.0 to avoid overflow.
	reportRisk := profile.ReportRate * 0.40
	if reportRisk > 0.40 {
		reportRisk = 0.40
	}

	// Account age contribution: accounts younger than 90 days carry extra risk.
	ageRisk := 0.0
	if profile.AccountAgeDays < 90 {
		ageRisk = (1.0 - float64(profile.AccountAgeDays)/90.0) * 0.25
	}

	score := responseRisk + reportRisk + ageRisk
	if score > 1.0 {
		score = 1.0
	}
	return score
}

// DailyLimit returns the maximum number of InMails the sender is allowed to
// send per day based on their account type.
//
//   - Free member:  10 InMails/day
//   - Recruiter:    100 InMails/day
func DailyLimit(profile SenderProfile) int {
	if profile.IsRecruiter {
		return 100
	}
	return 10
}

// CanSend determines whether the sender is currently allowed to send an InMail.
// Returns (true, "") if sending is permitted, or (false, reason) if blocked.
//
// Blocking conditions (checked in order):
//  1. New account throttle: accounts < 30 days old are limited to 5 InMails/day.
//  2. Daily send limit: DailySentCount >= DailyLimit(profile).
//  3. Spam score threshold: SpamScore > 0.7 indicates likely spam.
func CanSend(profile SenderProfile, sentToday int) (bool, string) {
	// New account hard throttle
	newAccountLimit := 5
	if profile.AccountAgeDays < 30 && sentToday >= newAccountLimit {
		return false, fmt.Sprintf("new account throttle: accounts under 30 days are limited to %d InMails/day", newAccountLimit)
	}

	// Daily limit check
	limit := DailyLimit(profile)
	if sentToday >= limit {
		return false, fmt.Sprintf("daily limit reached: %d of %d InMails used today", sentToday, limit)
	}

	// Spam score check
	score := SpamScore(profile)
	if score > 0.7 {
		return false, fmt.Sprintf("spam score too high: %.2f (threshold 0.70) — improve response rate and reduce reports", score)
	}

	return true, ""
}

func main() {
	_ = time.Now() // import used in production for TTL calculations

	profiles := []struct {
		name       string
		profile    SenderProfile
		sentToday  int
	}{
		{
			name: "High-quality recruiter",
			profile: SenderProfile{
				AccountAgeDays: 365,
				ResponseRate:   0.55,
				ReportRate:     0.01,
				DailySentCount: 20,
				IsRecruiter:    true,
			},
			sentToday: 20,
		},
		{
			name: "New free account",
			profile: SenderProfile{
				AccountAgeDays: 8,
				ResponseRate:   0.0,
				ReportRate:     0.0,
				DailySentCount: 5,
				IsRecruiter:    false,
			},
			sentToday: 5,
		},
		{
			name: "Over daily limit (free member)",
			profile: SenderProfile{
				AccountAgeDays: 200,
				ResponseRate:   0.30,
				ReportRate:     0.02,
				DailySentCount: 10,
				IsRecruiter:    false,
			},
			sentToday: 10,
		},
		{
			name: "High spam score",
			profile: SenderProfile{
				AccountAgeDays: 45,
				ResponseRate:   0.05,
				ReportRate:     0.15,
				DailySentCount: 3,
				IsRecruiter:    false,
			},
			sentToday: 3,
		},
	}

	for _, p := range profiles {
		allowed, reason := CanSend(p.profile, p.sentToday)
		spam := SpamScore(p.profile)
		limit := DailyLimit(p.profile)
		fmt.Printf("%-35s | spam=%.2f | limit=%3d | allowed=%-5v | %s\n",
			p.name, spam, limit, allowed, reason)
	}
}
