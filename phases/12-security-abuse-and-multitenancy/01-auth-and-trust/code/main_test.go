package main

import "testing"

func TestValidateTrustPlan(t *testing.T) {
	tests := []struct {
		name      string
		plan      TrustPlan
		minIssues int
	}{
		{
			name: "strong plan",
			plan: TrustPlan{
				ExternalIdentityValidated: true,
				ResourceOwnerAuthz:        true,
				ShortLivedInternalToken:   true,
				BreakGlassAudited:         true,
				AsyncActorPropagation:     true,
			},
			minIssues: 0,
		},
		{
			name: "weak plan",
			plan: TrustPlan{
				ExternalIdentityValidated: true,
			},
			minIssues: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := ValidateTrustPlan(tt.plan)
			if len(issues) < tt.minIssues {
				t.Fatalf("expected at least %d issues, got %v", tt.minIssues, issues)
			}
			if tt.minIssues == 0 && len(issues) != 0 {
				t.Fatalf("expected no issues, got %v", issues)
			}
		})
	}
}
