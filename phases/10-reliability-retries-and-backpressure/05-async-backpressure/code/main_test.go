package main

import "testing"

func TestAssessPipelineFlagsGrowingBacklogWithoutFeedback(t *testing.T) {
	got := AssessPipeline(Pipeline{
		IngressPerSecond: 300000,
		ConsumePerSecond: 200000,
		QueueTTLSeconds:  0,
		RetryMultiplier:  1.3,
		ProducerFeedback: false,
		SeparateClasses:  false,
	})

	if got.Risk != "high" {
		t.Fatalf("risk = %q, want high", got.Risk)
	}
}

func TestAssessPipelineApprovesBalancedPipeline(t *testing.T) {
	got := AssessPipeline(Pipeline{
		IngressPerSecond: 200000,
		ConsumePerSecond: 240000,
		QueueTTLSeconds:  300,
		RetryMultiplier:  1.0,
		ProducerFeedback: true,
		SeparateClasses:  true,
	})

	if got.Risk != "low" {
		t.Fatalf("risk = %q, want low", got.Risk)
	}
}
