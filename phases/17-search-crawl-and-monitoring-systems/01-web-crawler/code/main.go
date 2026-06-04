package main

import (
	"encoding/json"
	"flag"
	"os"
)

type CrawlPlan struct {
	Name                string `json:"name"`
	Fetchers            int    `json:"fetchers"`
	PerHostQPS          int    `json:"per_host_qps"`
	QueueLeaseSeconds   int    `json:"queue_lease_seconds"`
	RobotsEnforced      bool   `json:"robots_enforced"`
	FrontierReplicas    int    `json:"frontier_replicas"`
	ContentDedupEnabled bool   `json:"content_dedup_enabled"`
	RecrawlTiers        int    `json:"recrawl_tiers"`
}

func ValidateCrawlPlan(plan CrawlPlan) []string {
	var issues []string
	if plan.Fetchers < 10 {
		issues = append(issues, "fetchers should be at least 10 for a distributed crawler")
	}
	if plan.PerHostQPS <= 0 {
		issues = append(issues, "per_host_qps must be positive to encode politeness")
	}
	if plan.QueueLeaseSeconds < 30 {
		issues = append(issues, "queue_lease_seconds should allow retries without immediate duplicate claims")
	}
	if !plan.RobotsEnforced {
		issues = append(issues, "robots_enforced should be true for public-web crawling")
	}
	if plan.FrontierReplicas < 3 {
		issues = append(issues, "frontier_replicas should be at least 3 for durable scheduling")
	}
	if !plan.ContentDedupEnabled {
		issues = append(issues, "content_dedup_enabled should be true to control duplicate fetch waste")
	}
	if plan.RecrawlTiers < 2 {
		issues = append(issues, "recrawl_tiers should separate at least hot and cold freshness classes")
	}
	return issues
}

func main() {
	name := flag.String("name", "public-web-crawler", "crawl plan name")
	flag.Parse()

	plan := CrawlPlan{
		Name:                *name,
		Fetchers:            200,
		PerHostQPS:          2,
		QueueLeaseSeconds:   120,
		RobotsEnforced:      true,
		FrontierReplicas:    3,
		ContentDedupEnabled: true,
		RecrawlTiers:        3,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateCrawlPlan(plan),
	})
}
