package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// UseCase describes the properties of a messaging use case.
type UseCase struct {
	// NeedOrdering requires strict message ordering within a group.
	NeedOrdering bool
	// NeedFanOut requires delivering the same message to multiple independent consumers.
	NeedFanOut bool
	// NeedPersistence requires messages to survive consumer restarts and support replay.
	NeedPersistence bool
	// RetentionHours is the number of hours messages must be retained for replay.
	// 0 means no replay needed beyond the current session.
	RetentionHours int
	// ConsumerCount is the number of independent consumer groups that need the same event.
	ConsumerCount int
	// IsMultiStepWorkflow is true when the use case is a sequence of dependent steps
	// with per-step retry and a durable audit trail requirement.
	IsMultiStepWorkflow bool
}

// Recommendation holds the service name and the primary reason for the choice.
type Recommendation struct {
	Service string `json:"service"`
	Reason  string `json:"reason"`
}

// RecommendService returns the most appropriate AWS messaging service for the given UseCase.
// Decision logic:
//  1. Multi-step workflow with audit trail → StepFunctions
//  2. Fan-out to multiple consumers:
//     - With content filtering need (ConsumerCount > 3) → EventBridge
//     - Otherwise → SNS (paired with SQS for durability)
//  3. Ordering + persistence + multiple consumers → Kinesis
//  4. Ordering + single consumer → SQS FIFO
//  5. Persistence/replay without ordering → Kinesis
//  6. Default point-to-point → SQS
func RecommendService(uc UseCase) Recommendation {
	switch {
	case uc.IsMultiStepWorkflow:
		return Recommendation{
			Service: "StepFunctions",
			Reason:  "multi-step workflow requires durable state, per-step retry, and audit trail; Step Functions Standard preserves progress across step failures",
		}

	case uc.NeedFanOut && uc.ConsumerCount > 3:
		return Recommendation{
			Service: "EventBridge",
			Reason:  fmt.Sprintf("fan-out to %d consumers benefits from content-based filtering rules so each consumer receives only relevant event types; EventBridge schema registry reduces contract drift", uc.ConsumerCount),
		}

	case uc.NeedFanOut && !uc.NeedPersistence:
		return Recommendation{
			Service: "SNS",
			Reason:  fmt.Sprintf("fan-out to %d consumers with push delivery; pair each SNS subscription with an SQS queue for per-consumer durability", uc.ConsumerCount),
		}

	case uc.NeedFanOut && uc.NeedPersistence:
		return Recommendation{
			Service: "SNS",
			Reason:  fmt.Sprintf("fan-out to %d consumers; use SNS → SQS subscriptions so each consumer queue buffers events independently during consumer downtime", uc.ConsumerCount),
		}

	case uc.NeedOrdering && uc.NeedPersistence && uc.ConsumerCount > 1:
		return Recommendation{
			Service: "Kinesis",
			Reason:  fmt.Sprintf("ordered stream with %d independent consumer groups and %dh retention; each consumer tracks its own shard iterator for replay", uc.ConsumerCount, uc.RetentionHours),
		}

	case uc.NeedOrdering:
		return Recommendation{
			Service: "SQS",
			Reason:  "single consumer with ordering and deduplication requirement; use SQS FIFO queue with message group ID per logical entity",
		}

	case uc.NeedPersistence && uc.RetentionHours > 24:
		return Recommendation{
			Service: "Kinesis",
			Reason:  fmt.Sprintf("%dh retention exceeds SQS practical replay capability; Kinesis retains records independently of consumption and supports multiple consumer groups", uc.RetentionHours),
		}

	default:
		return Recommendation{
			Service: "SQS",
			Reason:  "point-to-point command or event with single consumer; SQS Standard provides at-least-once delivery, DLQ support, and backpressure via queue depth",
		}
	}
}

func main() {
	useCases := []struct {
		name string
		uc   UseCase
	}{
		{
			name: "Order placement → fulfillment (single consumer, ordered)",
			uc:   UseCase{NeedOrdering: true, ConsumerCount: 1},
		},
		{
			name: "Order placed broadcast to 10 downstream services",
			uc:   UseCase{NeedFanOut: true, ConsumerCount: 10},
		},
		{
			name: "Analytics clickstream (ordered, 3-day replay, 3 consumers)",
			uc:   UseCase{NeedOrdering: true, NeedPersistence: true, RetentionHours: 72, ConsumerCount: 3},
		},
		{
			name: "Order fulfillment pipeline (payment → inventory → notify)",
			uc:   UseCase{IsMultiStepWorkflow: true},
		},
		{
			name: "Catalog update fan-out to 2 consumers, no persistence",
			uc:   UseCase{NeedFanOut: true, NeedPersistence: false, ConsumerCount: 2},
		},
		{
			name: "Fraud detection async job (single consumer, no ordering)",
			uc:   UseCase{ConsumerCount: 1},
		},
		{
			name: "Compliance event replay (7 days retention, no ordering)",
			uc:   UseCase{NeedPersistence: true, RetentionHours: 168, ConsumerCount: 1},
		},
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	for _, tc := range useCases {
		rec := RecommendService(tc.uc)
		fmt.Printf("Use case: %s\n", tc.name)
		_ = enc.Encode(rec)
		fmt.Println()
	}
}
