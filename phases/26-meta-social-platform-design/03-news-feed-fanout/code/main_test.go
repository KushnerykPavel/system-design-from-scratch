package main

import "testing"

func TestSelectFanoutStrategy_PushBelowThreshold(t *testing.T) {
	account := Account{ID: 1, Name: "alice", FollowerCount: 49_999}
	strategy := SelectFanoutStrategy(account, 50_000)
	if strategy != PUSH {
		t.Errorf("expected PUSH for follower count below threshold, got %s", strategy)
	}
}

func TestSelectFanoutStrategy_PullAtThreshold(t *testing.T) {
	account := Account{ID: 2, Name: "celebrity", FollowerCount: 50_000}
	strategy := SelectFanoutStrategy(account, 50_000)
	if strategy != PULL {
		t.Errorf("expected PULL for follower count at threshold, got %s", strategy)
	}
}

func TestSelectFanoutStrategy_PullAboveThreshold(t *testing.T) {
	account := Account{ID: 3, Name: "mega_celebrity", FollowerCount: 100_000_000}
	strategy := SelectFanoutStrategy(account, 50_000)
	if strategy != PULL {
		t.Errorf("expected PULL for follower count above threshold, got %s", strategy)
	}
}

func TestSelectFanoutStrategy_ZeroFollowers(t *testing.T) {
	account := Account{ID: 4, Name: "new_user", FollowerCount: 0}
	strategy := SelectFanoutStrategy(account, 50_000)
	if strategy != PUSH {
		t.Errorf("expected PUSH for zero followers, got %s", strategy)
	}
}

func TestSimulateFanout_CorrectCounts(t *testing.T) {
	accounts := []Account{
		{ID: 1, FollowerCount: 100},      // PUSH
		{ID: 2, FollowerCount: 30_000},   // PUSH
		{ID: 3, FollowerCount: 49_999},   // PUSH
		{ID: 4, FollowerCount: 50_000},   // PULL
		{ID: 5, FollowerCount: 500_000},  // PULL
		{ID: 6, FollowerCount: 10_000_000}, // PULL
	}

	result := SimulateFanout(accounts, 50_000)

	if result.PushCount != 3 {
		t.Errorf("expected 3 PUSH accounts, got %d", result.PushCount)
	}
	if result.PullCount != 3 {
		t.Errorf("expected 3 PULL accounts, got %d", result.PullCount)
	}
	if result.Total != 6 {
		t.Errorf("expected total 6 accounts, got %d", result.Total)
	}
}

func TestSimulateFanout_AllPush(t *testing.T) {
	accounts := []Account{
		{ID: 1, FollowerCount: 0},
		{ID: 2, FollowerCount: 1},
		{ID: 3, FollowerCount: 49_999},
	}

	result := SimulateFanout(accounts, 50_000)

	if result.PushCount != 3 {
		t.Errorf("expected all 3 as PUSH, got %d", result.PushCount)
	}
	if result.PullCount != 0 {
		t.Errorf("expected 0 PULL, got %d", result.PullCount)
	}
}

func TestSimulateFanout_AllPull(t *testing.T) {
	accounts := []Account{
		{ID: 1, FollowerCount: 50_000},
		{ID: 2, FollowerCount: 1_000_000},
	}

	result := SimulateFanout(accounts, 50_000)

	if result.PullCount != 2 {
		t.Errorf("expected all 2 as PULL, got %d", result.PullCount)
	}
	if result.PushCount != 0 {
		t.Errorf("expected 0 PUSH, got %d", result.PushCount)
	}
}

func TestSimulateFanout_EmptyAccounts(t *testing.T) {
	result := SimulateFanout([]Account{}, 50_000)
	if result.Total != 0 || result.PushCount != 0 || result.PullCount != 0 {
		t.Errorf("expected all zeros for empty account list, got %+v", result)
	}
}

func TestFanoutStrategy_String(t *testing.T) {
	cases := []struct {
		strategy FanoutStrategy
		want     string
	}{
		{PUSH, "PUSH"},
		{PULL, "PULL"},
		{HYBRID, "HYBRID"},
	}
	for _, tc := range cases {
		if got := tc.strategy.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}
