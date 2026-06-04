package main

import "testing"

func TestValidateChatPlanHealthy(t *testing.T) {
	plan := ChatPlan{
		Name:                    "healthy",
		AckMode:                 "durable_log",
		UsesClientMessageID:     true,
		HasOfflineReplay:        true,
		UsesPresenceHeartbeats:  true,
		PresenceTTLSeconds:      30,
		PerConversationOrdering: true,
		BackpressureControls:    true,
	}
	if issues := ValidateChatPlan(plan); len(issues) != 0 {
		t.Fatalf("ValidateChatPlan returned issues: %v", issues)
	}
}

func TestValidateChatPlanRejectsWeakPlan(t *testing.T) {
	plan := ChatPlan{
		Name:               "weak",
		AckMode:            "unknown",
		PresenceTTLSeconds: 1000,
	}
	if issues := ValidateChatPlan(plan); len(issues) < 5 {
		t.Fatalf("ValidateChatPlan returned too few issues: %v", issues)
	}
}
