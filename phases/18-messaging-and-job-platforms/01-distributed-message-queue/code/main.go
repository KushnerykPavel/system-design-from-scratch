package main

import (
	"encoding/json"
	"flag"
	"os"
)

type QueuePlan struct {
	Name                     string `json:"name"`
	HasDurableLog            bool   `json:"has_durable_log"`
	UsesVisibilityTimeout    bool   `json:"uses_visibility_timeout"`
	TracksConsumerOffsets    bool   `json:"tracks_consumer_offsets"`
	SupportsRedelivery       bool   `json:"supports_redelivery"`
	HasPartitionStrategy     bool   `json:"has_partition_strategy"`
	HasPoisonMessageHandling bool   `json:"has_poison_message_handling"`
	HasBackpressureControls  bool   `json:"has_backpressure_controls"`
}

func ValidateQueuePlan(plan QueuePlan) []string {
	var issues []string
	if !plan.HasDurableLog {
		issues = append(issues, "has_durable_log should be true so messages survive broker failure")
	}
	if !plan.UsesVisibilityTimeout {
		issues = append(issues, "uses_visibility_timeout should be true so abandoned deliveries can be retried")
	}
	if !plan.TracksConsumerOffsets {
		issues = append(issues, "tracks_consumer_offsets should be true so delivery progress is explicit")
	}
	if !plan.SupportsRedelivery {
		issues = append(issues, "supports_redelivery should be true because at-least-once delivery needs retry paths")
	}
	if !plan.HasPartitionStrategy {
		issues = append(issues, "has_partition_strategy should be true so ordering and throughput are grounded")
	}
	if !plan.HasPoisonMessageHandling {
		issues = append(issues, "has_poison_message_handling should be true so one bad message does not block a shard")
	}
	if !plan.HasBackpressureControls {
		issues = append(issues, "has_backpressure_controls should be true so slow consumers cannot destabilize the queue")
	}
	return issues
}

func main() {
	name := flag.String("name", "distributed-message-queue", "plan name")
	flag.Parse()

	plan := QueuePlan{
		Name:                     *name,
		HasDurableLog:            true,
		UsesVisibilityTimeout:    true,
		TracksConsumerOffsets:    true,
		SupportsRedelivery:       true,
		HasPartitionStrategy:     true,
		HasPoisonMessageHandling: true,
		HasBackpressureControls:  true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateQueuePlan(plan),
	})
}
