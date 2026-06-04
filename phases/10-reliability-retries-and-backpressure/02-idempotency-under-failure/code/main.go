package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type IdempotencyPolicy struct {
	DurableStore            bool `json:"durable_store"`
	ReserveBeforeSideEffect bool `json:"reserve_before_side_effect"`
	StoresRequestHash       bool `json:"stores_request_hash"`
	ReturnsStoredResponse   bool `json:"returns_stored_response"`
	TTLHours                int  `json:"ttl_hours"`
}

type PolicyAssessment struct {
	Safe     bool     `json:"safe"`
	Risk     string   `json:"risk"`
	Warnings []string `json:"warnings"`
}

func AssessIdempotencyPolicy(policy IdempotencyPolicy, requiredWindowHours int) PolicyAssessment {
	warnings := make([]string, 0, 4)
	score := 0

	if !policy.DurableStore {
		score += 2
		warnings = append(warnings, "dedupe state is not durable")
	}
	if !policy.ReserveBeforeSideEffect {
		score += 3
		warnings = append(warnings, "side effects can happen before the key is reserved")
	}
	if !policy.StoresRequestHash {
		score++
		warnings = append(warnings, "payload mismatch on duplicate reuse cannot be detected")
	}
	if policy.TTLHours < requiredWindowHours {
		score += 2
		warnings = append(warnings, "dedupe TTL is shorter than the required retry window")
	}
	if !policy.ReturnsStoredResponse {
		score++
		warnings = append(warnings, "duplicate callers will not get the original response body")
	}

	risk := "low"
	if score >= 3 {
		risk = "high"
	} else if score >= 2 {
		risk = "medium"
	}

	return PolicyAssessment{
		Safe:     score < 4,
		Risk:     risk,
		Warnings: warnings,
	}
}

func main() {
	policy := IdempotencyPolicy{
		DurableStore:            true,
		ReserveBeforeSideEffect: true,
		StoresRequestHash:       true,
		ReturnsStoredResponse:   true,
		TTLHours:                24,
	}

	encoded, err := json.MarshalIndent(AssessIdempotencyPolicy(policy, 12), "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(encoded))
}
