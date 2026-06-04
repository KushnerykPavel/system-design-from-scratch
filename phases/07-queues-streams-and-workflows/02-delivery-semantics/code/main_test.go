package main

import "testing"

func TestRecommend(t *testing.T) {
	tests := []struct {
		name     string
		profile  DeliveryProfile
		want     DeliverySemantics
		warnings int
	}{
		{
			name: "low value telemetry",
			profile: DeliveryProfile{
				LossTolerant:      true,
				DuplicateTolerant: true,
			},
			want: AtMostOnce,
		},
		{
			name: "safe default",
			profile: DeliveryProfile{
				DuplicateTolerant:    false,
				ConsumerIdempotent:   true,
				DeduplicationStore:   true,
				AckAfterDurableWrite: true,
				BrokerTransaction:    true,
			},
			want: ExactlyOnce,
		},
		{
			name: "external side effect warning",
			profile: DeliveryProfile{
				DuplicateTolerant:    false,
				ConsumerIdempotent:   true,
				DeduplicationStore:   true,
				AckAfterDurableWrite: true,
				BrokerTransaction:    true,
				ExternalSideEffect:   true,
			},
			want:     ExactlyOnce,
			warnings: 1,
		},
		{
			name: "at least once without idempotency",
			profile: DeliveryProfile{
				DuplicateTolerant: false,
			},
			want:     AtLeastOnce,
			warnings: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Recommend(tt.profile)
			if got.Semantics != tt.want {
				t.Fatalf("Recommend() semantics = %q, want %q", got.Semantics, tt.want)
			}
			if len(got.Warnings) != tt.warnings {
				t.Fatalf("Recommend() warnings = %d, want %d", len(got.Warnings), tt.warnings)
			}
		})
	}
}
