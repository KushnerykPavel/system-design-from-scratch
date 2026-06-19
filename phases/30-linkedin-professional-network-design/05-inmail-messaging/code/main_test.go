package main

import (
	"strings"
	"testing"
)

func TestSpamScore_HighReputationSenderIsClean(t *testing.T) {
	profile := SenderProfile{
		AccountAgeDays: 500,
		ResponseRate:   0.60,
		ReportRate:     0.005,
		IsRecruiter:    true,
	}

	score := SpamScore(profile)
	if score > 0.3 {
		t.Errorf("expected low spam score for high-reputation sender, got %.2f", score)
	}
}

func TestSpamScore_NewAccountWithZeroResponseIsHighRisk(t *testing.T) {
	profile := SenderProfile{
		AccountAgeDays: 5,
		ResponseRate:   0.0,
		ReportRate:     0.0,
		IsRecruiter:    false,
	}

	score := SpamScore(profile)
	// New account (age risk) + zero response rate (response risk) should push score high
	if score < 0.5 {
		t.Errorf("expected high spam score for new account with zero response rate, got %.2f", score)
	}
}

func TestSpamScore_HighReportRateDominates(t *testing.T) {
	profile := SenderProfile{
		AccountAgeDays: 365,
		ResponseRate:   0.50,
		ReportRate:     1.0, // 100% report rate — maximum reportRisk contribution (0.40)
		IsRecruiter:    true,
	}

	score := SpamScore(profile)
	// responseRisk = (1-0.50)*0.35 = 0.175; reportRisk = 0.40 (capped); ageRisk = 0 → total = 0.575
	if score < 0.5 {
		t.Errorf("expected spam score dominated by high report rate to exceed 0.5, got %.2f", score)
	}
}

func TestDailyLimit_RecruiterGetsHigherLimit(t *testing.T) {
	free := SenderProfile{IsRecruiter: false}
	recruiter := SenderProfile{IsRecruiter: true}

	if DailyLimit(free) >= DailyLimit(recruiter) {
		t.Errorf("recruiter limit (%d) should be higher than free limit (%d)", DailyLimit(recruiter), DailyLimit(free))
	}
}

func TestDailyLimit_FreeAccountIs10(t *testing.T) {
	profile := SenderProfile{IsRecruiter: false}
	if DailyLimit(profile) != 10 {
		t.Errorf("expected free daily limit of 10, got %d", DailyLimit(profile))
	}
}

func TestDailyLimit_RecruiterIs100(t *testing.T) {
	profile := SenderProfile{IsRecruiter: true}
	if DailyLimit(profile) != 100 {
		t.Errorf("expected recruiter daily limit of 100, got %d", DailyLimit(profile))
	}
}

func TestCanSend_HighReputationSenderAllowed(t *testing.T) {
	profile := SenderProfile{
		AccountAgeDays: 400,
		ResponseRate:   0.55,
		ReportRate:     0.01,
		IsRecruiter:    true,
	}

	allowed, reason := CanSend(profile, 5)
	if !allowed {
		t.Errorf("expected high-reputation recruiter to be allowed, blocked with: %s", reason)
	}
}

func TestCanSend_NewAccountThrottled(t *testing.T) {
	profile := SenderProfile{
		AccountAgeDays: 10,
		ResponseRate:   0.0,
		ReportRate:     0.0,
		IsRecruiter:    false,
	}

	// sentToday = 5, which hits the new account throttle limit
	allowed, reason := CanSend(profile, 5)
	if allowed {
		t.Errorf("expected new account at 5 sends to be throttled")
	}
	if !strings.Contains(reason, "new account") {
		t.Errorf("expected 'new account' in reason, got: %s", reason)
	}
}

func TestCanSend_OverDailyLimitBlocked(t *testing.T) {
	profile := SenderProfile{
		AccountAgeDays: 300,
		ResponseRate:   0.40,
		ReportRate:     0.02,
		IsRecruiter:    false,
	}

	// Free member has limit of 10; sentToday = 10 means limit reached
	allowed, reason := CanSend(profile, 10)
	if allowed {
		t.Errorf("expected over-daily-limit sender to be blocked")
	}
	if !strings.Contains(reason, "daily limit") {
		t.Errorf("expected 'daily limit' in reason, got: %s", reason)
	}
}

func TestCanSend_HighSpamScoreBlocked(t *testing.T) {
	// AccountAgeDays=10: ageRisk = (1 - 10/90) * 0.25 ≈ 0.222
	// ResponseRate=0.02: responseRisk = 0.98 * 0.35 ≈ 0.343
	// ReportRate=0.50: reportRisk = 0.50 * 0.40 = 0.20
	// Total ≈ 0.765 > 0.70 threshold → blocked by spam score
	// sentToday=3 is below the new-account throttle limit of 5, so spam score check runs
	profile := SenderProfile{
		AccountAgeDays: 10,
		ResponseRate:   0.02,
		ReportRate:     0.50,
		IsRecruiter:    false,
	}

	allowed, reason := CanSend(profile, 3)
	if allowed {
		t.Errorf("expected high-spam-score sender to be blocked")
	}
	if !strings.Contains(reason, "spam score") {
		t.Errorf("expected 'spam score' in reason, got: %s", reason)
	}
}

func TestCanSend_RecruiterAtLimitBlocked(t *testing.T) {
	profile := SenderProfile{
		AccountAgeDays: 200,
		ResponseRate:   0.50,
		ReportRate:     0.01,
		IsRecruiter:    true,
	}

	// Recruiter has limit of 100
	allowed, reason := CanSend(profile, 100)
	if allowed {
		t.Errorf("expected recruiter at 100 sends to be blocked")
	}
	if !strings.Contains(reason, "daily limit") {
		t.Errorf("expected 'daily limit' in reason, got: %s", reason)
	}
}
