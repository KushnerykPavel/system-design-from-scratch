package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Pipeline struct {
	IngressPerSecond int     `json:"ingress_per_second"`
	ConsumePerSecond int     `json:"consume_per_second"`
	QueueTTLSeconds  int     `json:"queue_ttl_seconds"`
	RetryMultiplier  float64 `json:"retry_multiplier"`
	ProducerFeedback bool    `json:"producer_feedback"`
	SeparateClasses  bool    `json:"separate_classes"`
}

type PipelineAssessment struct {
	BacklogGrowthPerSecond int      `json:"backlog_growth_per_second"`
	Risk                   string   `json:"risk"`
	Recommendations        []string `json:"recommendations"`
}

func AssessPipeline(p Pipeline) PipelineAssessment {
	recommendations := make([]string, 0, 4)
	effectiveIngress := int(float64(p.IngressPerSecond) * p.RetryMultiplier)
	backlogGrowth := effectiveIngress - p.ConsumePerSecond
	score := 0

	if backlogGrowth > 0 {
		score += 2
		recommendations = append(recommendations, "consumers are falling behind effective ingress")
	}
	if backlogGrowth > 0 && !p.ProducerFeedback {
		score++
		recommendations = append(recommendations, "producers are not receiving backpressure")
	}
	if p.QueueTTLSeconds == 0 {
		score++
		recommendations = append(recommendations, "message lifetime is unbounded")
	}
	if !p.SeparateClasses {
		score++
		recommendations = append(recommendations, "mixed-value traffic shares one backlog")
	}

	risk := "low"
	if score >= 4 {
		risk = "high"
	} else if score >= 2 {
		risk = "medium"
	}

	return PipelineAssessment{
		BacklogGrowthPerSecond: backlogGrowth,
		Risk:                   risk,
		Recommendations:        recommendations,
	}
}

func main() {
	assessment := AssessPipeline(Pipeline{
		IngressPerSecond: 400000,
		ConsumePerSecond: 280000,
		QueueTTLSeconds:  600,
		RetryMultiplier:  1.2,
		ProducerFeedback: true,
		SeparateClasses:  true,
	})

	encoded, err := json.MarshalIndent(assessment, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(string(encoded))
}
