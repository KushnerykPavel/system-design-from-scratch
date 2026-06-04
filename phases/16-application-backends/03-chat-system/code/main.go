package main

import (
	"encoding/json"
	"flag"
	"os"
)

type ChatPlan struct {
	Name                    string `json:"name"`
	AckMode                 string `json:"ack_mode"`
	UsesClientMessageID     bool   `json:"uses_client_message_id"`
	HasOfflineReplay        bool   `json:"has_offline_replay"`
	UsesPresenceHeartbeats  bool   `json:"uses_presence_heartbeats"`
	PresenceTTLSeconds      int    `json:"presence_ttl_seconds"`
	PerConversationOrdering bool   `json:"per_conversation_ordering"`
	BackpressureControls    bool   `json:"backpressure_controls"`
}

func ValidateChatPlan(plan ChatPlan) []string {
	var issues []string
	if plan.AckMode != "durable_log" && plan.AckMode != "recipient_device" {
		issues = append(issues, "ack_mode must be durable_log or recipient_device")
	}
	if !plan.UsesClientMessageID {
		issues = append(issues, "uses_client_message_id should be true for safe retry dedupe")
	}
	if !plan.HasOfflineReplay {
		issues = append(issues, "has_offline_replay should be true for reconnecting devices")
	}
	if !plan.UsesPresenceHeartbeats {
		issues = append(issues, "uses_presence_heartbeats should be true for routing hints")
	}
	if plan.PresenceTTLSeconds <= 0 || plan.PresenceTTLSeconds > 300 {
		issues = append(issues, "presence_ttl_seconds should be a realistic heartbeat expiry bound")
	}
	if !plan.PerConversationOrdering {
		issues = append(issues, "per_conversation_ordering should be explicit even if global ordering is not")
	}
	if !plan.BackpressureControls {
		issues = append(issues, "backpressure_controls should be enabled so one noisy conversation does not starve the fleet")
	}
	return issues
}

func main() {
	name := flag.String("name", "chat-core", "plan name")
	flag.Parse()

	plan := ChatPlan{
		Name:                    *name,
		AckMode:                 "durable_log",
		UsesClientMessageID:     true,
		HasOfflineReplay:        true,
		UsesPresenceHeartbeats:  true,
		PresenceTTLSeconds:      30,
		PerConversationOrdering: true,
		BackpressureControls:    true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateChatPlan(plan),
	})
}
