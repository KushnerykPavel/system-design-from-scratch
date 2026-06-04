package main

type DeliverySemantics string

const (
	AtMostOnce  DeliverySemantics = "at-most-once"
	AtLeastOnce DeliverySemantics = "at-least-once"
	ExactlyOnce DeliverySemantics = "exactly-once"
)

type DeliveryProfile struct {
	LossTolerant         bool
	DuplicateTolerant    bool
	ReplayRequired       bool
	ExternalSideEffect   bool
	ConsumerIdempotent   bool
	BrokerTransaction    bool
	DeduplicationStore   bool
	AckAfterDurableWrite bool
}

type Recommendation struct {
	Semantics DeliverySemantics
	Warnings  []string
}

func Recommend(profile DeliveryProfile) Recommendation {
	rec := Recommendation{}

	if profile.LossTolerant && !profile.ReplayRequired && profile.DuplicateTolerant {
		rec.Semantics = AtMostOnce
	} else {
		rec.Semantics = AtLeastOnce
	}

	if !profile.DuplicateTolerant && profile.ConsumerIdempotent && profile.DeduplicationStore && profile.AckAfterDurableWrite {
		rec.Semantics = ExactlyOnce
	}

	if rec.Semantics == ExactlyOnce && profile.ExternalSideEffect {
		rec.Warnings = append(rec.Warnings, "exactly-once does not fully cover external side effects")
	}
	if rec.Semantics == AtLeastOnce && !profile.ConsumerIdempotent {
		rec.Warnings = append(rec.Warnings, "at-least-once without idempotent consumer risks duplicate side effects")
	}
	if rec.Semantics == AtMostOnce && !profile.LossTolerant {
		rec.Warnings = append(rec.Warnings, "at-most-once chosen for a loss-intolerant workflow")
	}
	if rec.Semantics == ExactlyOnce && !profile.BrokerTransaction {
		rec.Warnings = append(rec.Warnings, "exactly-once claim depends on application boundary, not broker transaction alone")
	}

	return rec
}

func main() {}
