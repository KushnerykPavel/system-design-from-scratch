package main

import (
	"fmt"
	"strings"
)

var requiredTopics = map[string][]string{
	"shard key":        {"shard key", "shard by", "sharding by"},
	"hot partition":    {"hot partition", "hot partitions", "hotspot"},
	"tenant isolation": {"tenant isolation", "noisy neighbor", "dedicated pool"},
	"directory":        {"directory", "routing layer", "lookup"},
	"rebalancing":      {"rebalancing", "rebalance", "move tenant"},
	"resharding":       {"resharding", "reshard", "split tenant"},
	"cross-shard":      {"cross-shard", "scatter-gather", "derived view"},
	"observability":    {"observability", "metrics", "slo"},
	"trade-off":        {"trade-off", "tradeoff", "cost"},
}

func coverage(answer string) map[string]bool {
	answer = strings.ToLower(answer)
	result := make(map[string]bool, len(requiredTopics))
	for topic, aliases := range requiredTopics {
		for _, alias := range aliases {
			if strings.Contains(answer, alias) {
				result[topic] = true
				break
			}
		}
	}
	return result
}

func score(answer string) int {
	total := 0
	for _, covered := range coverage(answer) {
		if covered {
			total++
		}
	}
	return total
}

func main() {
	answer := "Shard key by tenant, keep a directory for moves, handle hot partitions, plan rebalancing and resharding, and keep cross-shard dashboards derived with observability and trade-off discussion."
	fmt.Printf("coverage score: %d/%d\n", score(answer), len(requiredTopics))
}
