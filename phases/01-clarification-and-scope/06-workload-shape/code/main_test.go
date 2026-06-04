package main

import "testing"

func TestValidateWorkloadProfileAcceptsBasicProfile(t *testing.T) {
	t.Parallel()

	profile := WorkloadProfile{
		Journey:     "user reads feed",
		ReadQPS:     1000,
		WriteQPS:    100,
		Fanout:      1,
		BurstFactor: 4,
	}

	if issues := ValidateWorkloadProfile(profile); len(issues) != 0 {
		t.Fatalf("ValidateWorkloadProfile() returned issues: %v", issues)
	}
}

func TestValidateWorkloadProfileRejectsMissingTraffic(t *testing.T) {
	t.Parallel()

	profile := WorkloadProfile{
		Journey:     "unknown",
		Fanout:      1,
		BurstFactor: 1,
	}

	if issues := ValidateWorkloadProfile(profile); len(issues) == 0 {
		t.Fatal("ValidateWorkloadProfile() returned no issues for missing traffic")
	}
}

func TestDominantPathAccountsForFanout(t *testing.T) {
	t.Parallel()

	profile := WorkloadProfile{
		Journey:     "post to followers",
		ReadQPS:     1000,
		WriteQPS:    100,
		Fanout:      100,
		BurstFactor: 2,
	}

	if got, want := DominantPath(profile), "write"; got != want {
		t.Fatalf("DominantPath() = %q, want %q", got, want)
	}
}
