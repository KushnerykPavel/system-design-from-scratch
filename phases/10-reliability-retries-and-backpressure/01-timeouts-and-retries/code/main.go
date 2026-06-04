package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type RetryProfile struct {
	Fanout              int  `json:"fanout"`
	TimeoutMS           int  `json:"timeout_ms"`
	P99LatencyMS        int  `json:"p99_latency_ms"`
	MaxAttempts         int  `json:"max_attempts"`
	HasJitter           bool `json:"has_jitter"`
	PropagatesDeadline  bool `json:"propagates_deadline"`
	OperationIdempotent bool `json:"operation_idempotent"`
}

type Assessment struct {
	AmplificationFactor float64  `json:"amplification_factor"`
	Risk                string   `json:"risk"`
	Warnings            []string `json:"warnings"`
}

func AssessPolicy(p RetryProfile) Assessment {
	warnings := make([]string, 0, 4)
	amplification := 1 + float64(max(0, p.MaxAttempts-1)*max(1, p.Fanout))
	score := 0

	if p.TimeoutMS > 0 && p.P99LatencyMS > 0 && p.TimeoutMS < int(float64(p.P99LatencyMS)*0.75) {
		score += 2
		warnings = append(warnings, "timeout is likely below real tail latency")
	}
	if p.MaxAttempts > 2 && p.Fanout >= 3 {
		score += 2
		warnings = append(warnings, "retry amplification is high for the fanout width")
	}
	if !p.HasJitter {
		score++
		warnings = append(warnings, "backoff jitter is missing")
	}
	if !p.PropagatesDeadline {
		score++
		warnings = append(warnings, "deadline is not propagated across hops")
	}
	if !p.OperationIdempotent && p.MaxAttempts > 1 {
		score += 2
		warnings = append(warnings, "non-idempotent operation should not be retried blindly")
	}

	risk := "low"
	if score >= 5 {
		risk = "high"
	} else if score >= 2 {
		risk = "medium"
	}

	return Assessment{
		AmplificationFactor: amplification,
		Risk:                risk,
		Warnings:            warnings,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	profile := RetryProfile{
		Fanout:              4,
		TimeoutMS:           120,
		P99LatencyMS:        180,
		MaxAttempts:         3,
		HasJitter:           true,
		PropagatesDeadline:  true,
		OperationIdempotent: true,
	}

	encoded, err := json.MarshalIndent(AssessPolicy(profile), "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(encoded))
}
