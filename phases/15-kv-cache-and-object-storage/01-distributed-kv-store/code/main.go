package main

import (
	"encoding/json"
	"flag"
	"os"
)

type KVTopology struct {
	Name              string `json:"name"`
	Replicas          int    `json:"replicas"`
	WriteQuorum       int    `json:"write_quorum"`
	ReadQuorum        int    `json:"read_quorum"`
	FailureDomains    int    `json:"failure_domains"`
	ConsistencyMode   string `json:"consistency_mode"`
	RepairEnabled     bool   `json:"repair_enabled"`
	HintedHandoff     bool   `json:"hinted_handoff"`
	HotKeyMitigation  bool   `json:"hot_key_mitigation"`
	ConditionalWrites bool   `json:"conditional_writes"`
}

func ValidateTopology(cfg KVTopology) []string {
	var issues []string
	if cfg.Replicas < 3 {
		issues = append(issues, "replicas should be at least 3 for durable regional service")
	}
	if cfg.WriteQuorum <= 0 || cfg.WriteQuorum > cfg.Replicas {
		issues = append(issues, "write_quorum must be between 1 and replicas")
	}
	if cfg.ReadQuorum <= 0 || cfg.ReadQuorum > cfg.Replicas {
		issues = append(issues, "read_quorum must be between 1 and replicas")
	}
	if cfg.FailureDomains < cfg.Replicas {
		issues = append(issues, "replicas should span at least as many failure domains as replica copies")
	}
	if cfg.ConsistencyMode != "eventual" && cfg.ConsistencyMode != "quorum" && cfg.ConsistencyMode != "session" {
		issues = append(issues, "consistency_mode must be eventual, quorum, or session")
	}
	if !cfg.RepairEnabled {
		issues = append(issues, "repair_enabled should be true for replicated storage")
	}
	if cfg.ConsistencyMode == "eventual" && !cfg.HintedHandoff {
		issues = append(issues, "eventual mode should define hinted handoff or another missed-write recovery path")
	}
	if !cfg.HotKeyMitigation {
		issues = append(issues, "hot_key_mitigation should be enabled for skewed real-world traffic")
	}
	if !cfg.ConditionalWrites {
		issues = append(issues, "conditional_writes should be supported for conflict-sensitive callers")
	}
	return issues
}

func main() {
	name := flag.String("name", "regional-config-kv", "name of the topology")
	flag.Parse()

	cfg := KVTopology{
		Name:              *name,
		Replicas:          3,
		WriteQuorum:       2,
		ReadQuorum:        2,
		FailureDomains:    3,
		ConsistencyMode:   "quorum",
		RepairEnabled:     true,
		HintedHandoff:     true,
		HotKeyMitigation:  true,
		ConditionalWrites: true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"topology": cfg,
		"issues":   ValidateTopology(cfg),
	})
}
