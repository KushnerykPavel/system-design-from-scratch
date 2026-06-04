package main

import (
	"encoding/json"
	"flag"
	"os"
)

type TrustPlan struct {
	ExternalIdentityValidated bool `json:"external_identity_validated"`
	ResourceOwnerAuthz        bool `json:"resource_owner_authz"`
	ShortLivedInternalToken   bool `json:"short_lived_internal_token"`
	BreakGlassAudited         bool `json:"break_glass_audited"`
	AsyncActorPropagation     bool `json:"async_actor_propagation"`
}

func ValidateTrustPlan(plan TrustPlan) []string {
	var issues []string
	if !plan.ExternalIdentityValidated {
		issues = append(issues, "external identity must be validated at the ingress boundary")
	}
	if !plan.ResourceOwnerAuthz {
		issues = append(issues, "resource-owning service should enforce authorization")
	}
	if !plan.ShortLivedInternalToken {
		issues = append(issues, "internal identity should be short-lived and scoped")
	}
	if !plan.BreakGlassAudited {
		issues = append(issues, "break-glass access requires explicit audit coverage")
	}
	if !plan.AsyncActorPropagation {
		issues = append(issues, "async jobs should preserve actor and tenant context")
	}
	return issues
}

func main() {
	fail := flag.Bool("fail", false, "emit an intentionally weak plan")
	flag.Parse()

	plan := TrustPlan{
		ExternalIdentityValidated: true,
		ResourceOwnerAuthz:        true,
		ShortLivedInternalToken:   true,
		BreakGlassAudited:         true,
		AsyncActorPropagation:     true,
	}
	if *fail {
		plan.ResourceOwnerAuthz = false
		plan.AsyncActorPropagation = false
	}

	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"plan":   plan,
		"issues": ValidateTrustPlan(plan),
	})
}
