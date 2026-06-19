package main

import (
	"strings"
	"testing"
)

// TestSQSRecommendation verifies that a single-consumer ordered use case recommends SQS.
func TestSQSRecommendation(t *testing.T) {
	uc := UseCase{
		NeedOrdering:  true,
		ConsumerCount: 1,
	}
	rec := RecommendService(uc)
	if rec.Service != "SQS" {
		t.Fatalf("expected SQS for ordering+single consumer, got %s", rec.Service)
	}
	if !strings.Contains(strings.ToLower(rec.Reason), "fifo") {
		t.Errorf("expected FIFO mentioned in reason, got: %s", rec.Reason)
	}
}

// TestSNSRecommendationFanOutNoPersistence verifies SNS for fan-out without persistence.
func TestSNSRecommendationFanOutNoPersistence(t *testing.T) {
	uc := UseCase{
		NeedFanOut:      true,
		NeedPersistence: false,
		ConsumerCount:   2,
	}
	rec := RecommendService(uc)
	if rec.Service != "SNS" {
		t.Fatalf("expected SNS for fan-out without persistence, got %s", rec.Service)
	}
}

// TestKinesisRecommendation verifies Kinesis for ordered stream with multiple consumers and persistence.
func TestKinesisRecommendation(t *testing.T) {
	uc := UseCase{
		NeedOrdering:    true,
		NeedPersistence: true,
		RetentionHours:  72,
		ConsumerCount:   3,
	}
	rec := RecommendService(uc)
	if rec.Service != "Kinesis" {
		t.Fatalf("expected Kinesis for ordering+persistence+multiple consumers, got %s", rec.Service)
	}
}

// TestStepFunctionsRecommendation verifies Step Functions for multi-step workflows.
func TestStepFunctionsRecommendation(t *testing.T) {
	uc := UseCase{
		IsMultiStepWorkflow: true,
	}
	rec := RecommendService(uc)
	if rec.Service != "StepFunctions" {
		t.Fatalf("expected StepFunctions for multi-step workflow, got %s", rec.Service)
	}
}

// TestEventBridgeRecommendation verifies EventBridge for large fan-out (>3 consumers).
func TestEventBridgeRecommendation(t *testing.T) {
	uc := UseCase{
		NeedFanOut:    true,
		ConsumerCount: 10,
	}
	rec := RecommendService(uc)
	if rec.Service != "EventBridge" {
		t.Fatalf("expected EventBridge for fan-out to 10 consumers, got %s", rec.Service)
	}
}

// TestDefaultSQSRecommendation verifies SQS for a basic point-to-point use case.
func TestDefaultSQSRecommendation(t *testing.T) {
	uc := UseCase{
		ConsumerCount: 1,
	}
	rec := RecommendService(uc)
	if rec.Service != "SQS" {
		t.Fatalf("expected SQS for default single-consumer use case, got %s", rec.Service)
	}
}

// TestKinesisLongRetention verifies Kinesis for long-retention replay use cases.
func TestKinesisLongRetention(t *testing.T) {
	uc := UseCase{
		NeedPersistence: true,
		RetentionHours:  168, // 7 days
		ConsumerCount:   1,
	}
	rec := RecommendService(uc)
	if rec.Service != "Kinesis" {
		t.Fatalf("expected Kinesis for 7-day retention, got %s", rec.Service)
	}
}
