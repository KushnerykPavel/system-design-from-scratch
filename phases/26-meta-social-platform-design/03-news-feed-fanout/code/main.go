package main

import (
	"fmt"
)

// FanoutStrategy represents the fanout approach for an account.
type FanoutStrategy int

const (
	// PUSH precomputes the feed by writing to each follower's feed store at post time.
	PUSH FanoutStrategy = iota
	// PULL computes the feed at read time by querying the author's post store.
	PULL
	// HYBRID pushes to online followers only and skips offline users; celebrities use pull.
	HYBRID
)

func (s FanoutStrategy) String() string {
	switch s {
	case PUSH:
		return "PUSH"
	case PULL:
		return "PULL"
	case HYBRID:
		return "HYBRID"
	default:
		return "UNKNOWN"
	}
}

// Account represents a social network account with a follower count.
type Account struct {
	ID            int64
	Name          string
	FollowerCount int
}

// SelectFanoutStrategy returns the appropriate fanout strategy for an account
// given the configured celebrity threshold.
//
// Below threshold: PUSH fanout — feed workers write to each follower's cache.
// Above threshold: PULL fanout — celebrity posts are fetched at read time.
// Exactly at threshold: PULL (the threshold marks the boundary where push cost becomes unacceptable).
func SelectFanoutStrategy(account Account, threshold int) FanoutStrategy {
	if account.FollowerCount < threshold {
		return PUSH
	}
	return PULL
}

// FanoutResult summarises how accounts are distributed across strategies.
type FanoutResult struct {
	PushCount   int
	PullCount   int
	HybridCount int // reserved for future semi-celebrity tier
	Total       int
}

// SimulateFanout classifies each account and returns a breakdown by strategy.
func SimulateFanout(accounts []Account, threshold int) FanoutResult {
	var result FanoutResult
	result.Total = len(accounts)
	for _, a := range accounts {
		switch SelectFanoutStrategy(a, threshold) {
		case PUSH:
			result.PushCount++
		case PULL:
			result.PullCount++
		case HYBRID:
			result.HybridCount++
		}
	}
	return result
}

func main() {
	// Sample account set mixing normal users and celebrity-tier accounts.
	// Facebook uses ~50,000 followers as the push/pull threshold.
	const threshold = 50_000

	accounts := []Account{
		{ID: 1, Name: "alice", FollowerCount: 150},
		{ID: 2, Name: "bob", FollowerCount: 3_200},
		{ID: 3, Name: "carol", FollowerCount: 49_999},
		{ID: 4, Name: "celebrity_a", FollowerCount: 50_000},
		{ID: 5, Name: "celebrity_b", FollowerCount: 1_200_000},
		{ID: 6, Name: "mega_celebrity", FollowerCount: 100_000_000},
		{ID: 7, Name: "dave", FollowerCount: 800},
		{ID: 8, Name: "eve", FollowerCount: 25_000},
	}

	fmt.Printf("Fanout strategy simulation (threshold = %d followers)\n\n", threshold)
	fmt.Printf("%-20s %12s %s\n", "Account", "Followers", "Strategy")
	fmt.Printf("%-20s %12s %s\n", "-------", "---------", "--------")

	for _, a := range accounts {
		strategy := SelectFanoutStrategy(a, threshold)
		fmt.Printf("%-20s %12d %s\n", a.Name, a.FollowerCount, strategy)
	}

	result := SimulateFanout(accounts, threshold)
	fmt.Printf("\nBreakdown:\n")
	fmt.Printf("  PUSH  : %d accounts (%.0f%%)\n", result.PushCount, 100*float64(result.PushCount)/float64(result.Total))
	fmt.Printf("  PULL  : %d accounts (%.0f%%)\n", result.PullCount, 100*float64(result.PullCount)/float64(result.Total))
	fmt.Printf("  Total : %d accounts\n", result.Total)
	fmt.Printf("\nInsight: %d celebrity account(s) avoided fanout storms by using PULL.\n", result.PullCount)
	fmt.Printf("         Without hybrid routing, a single post from mega_celebrity\n")
	fmt.Printf("         would trigger %d Memcached writes.\n", accounts[5].FollowerCount)
}
