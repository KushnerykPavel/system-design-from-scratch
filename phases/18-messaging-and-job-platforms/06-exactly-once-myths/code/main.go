package main

import (
	"encoding/json"
	"flag"
	"os"
)

type ExactlyOnceClaim struct {
	Name                   string `json:"name"`
	DefinesBoundary        bool   `json:"defines_boundary"`
	HasIdempotentConsumer  bool   `json:"has_idempotent_consumer"`
	HasDedupKey            bool   `json:"has_dedup_key"`
	IncludesSideEffects    bool   `json:"includes_side_effects"`
	ExplainsFailureCase    bool   `json:"explains_failure_case"`
	UsesTransactionalStore bool   `json:"uses_transactional_store"`
	NamesResidualRisk      bool   `json:"names_residual_risk"`
}

func ValidateExactlyOnceClaim(claim ExactlyOnceClaim) []string {
	var issues []string
	if !claim.DefinesBoundary {
		issues = append(issues, "defines_boundary should be true because exactly-once claims need scope")
	}
	if !claim.HasIdempotentConsumer {
		issues = append(issues, "has_idempotent_consumer should be true because duplicate attempts still happen operationally")
	}
	if !claim.HasDedupKey {
		issues = append(issues, "has_dedup_key should be true if duplicate detection is part of the claim")
	}
	if !claim.ExplainsFailureCase {
		issues = append(issues, "explains_failure_case should be true so ambiguity after partial success is addressed")
	}
	if !claim.NamesResidualRisk {
		issues = append(issues, "names_residual_risk should be true because interview answers should not oversell guarantees")
	}
	if claim.IncludesSideEffects && !claim.UsesTransactionalStore {
		issues = append(issues, "uses_transactional_store should be true when side effects are included in the exactly-once boundary")
	}
	return issues
}

func main() {
	name := flag.String("name", "exactly-once-claim", "claim name")
	flag.Parse()

	claim := ExactlyOnceClaim{
		Name:                   *name,
		DefinesBoundary:        true,
		HasIdempotentConsumer:  true,
		HasDedupKey:            true,
		IncludesSideEffects:    false,
		ExplainsFailureCase:    true,
		UsesTransactionalStore: false,
		NamesResidualRisk:      true,
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"claim":  claim,
		"issues": ValidateExactlyOnceClaim(claim),
	})
}
