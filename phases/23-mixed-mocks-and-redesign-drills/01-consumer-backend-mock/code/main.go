package main

import (
	"encoding/json"
	"flag"
	"os"
)

type ConsumerMockScorecard struct {
	Name                    string `json:"name"`
	ClarifiesUserPromise    bool   `json:"clarifies_user_promise"`
	SizesReadWriteAsymmetry bool   `json:"sizes_read_write_asymmetry"`
	ChoosesStateBoundary    bool   `json:"chooses_state_boundary"`
	CoversSkewOrFanout      bool   `json:"covers_skew_or_fanout"`
	NamesDegradedMode       bool   `json:"names_degraded_mode"`
	CoversObservability     bool   `json:"covers_observability"`
	HandlesRedesign         bool   `json:"handles_redesign"`
}

func ValidateConsumerMock(card ConsumerMockScorecard) []string {
	var issues []string
	if !card.ClarifiesUserPromise {
		issues = append(issues, "clarifies_user_promise should be true so the answer anchors on one user-visible contract")
	}
	if !card.SizesReadWriteAsymmetry {
		issues = append(issues, "sizes_read_write_asymmetry should be true so fanout and storage decisions are workload-aware")
	}
	if !card.ChoosesStateBoundary {
		issues = append(issues, "chooses_state_boundary should be true so correctness claims have a source of truth")
	}
	if !card.CoversSkewOrFanout {
		issues = append(issues, "covers_skew_or_fanout should be true so hotspot pressure is not hand-waved away")
	}
	if !card.NamesDegradedMode {
		issues = append(issues, "names_degraded_mode should be true so the system has a credible incident story")
	}
	if !card.CoversObservability {
		issues = append(issues, "covers_observability should be true so the main promise is measurable")
	}
	if !card.HandlesRedesign {
		issues = append(issues, "handles_redesign should be true so changed constraints lead to a concrete update")
	}
	return issues
}

func main() {
	name := flag.String("name", "consumer-backend-mock", "mock name")
	flag.Parse()

	card := ConsumerMockScorecard{
		Name:                    *name,
		ClarifiesUserPromise:    true,
		SizesReadWriteAsymmetry: true,
		ChoosesStateBoundary:    true,
		CoversSkewOrFanout:      true,
		NamesDegradedMode:       true,
		CoversObservability:     true,
		HandlesRedesign:         true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"scorecard": card,
		"issues":    ValidateConsumerMock(card),
	})
}
