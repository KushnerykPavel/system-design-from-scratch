package main

type IsolationLevel string

const (
	ReadCommitted  IsolationLevel = "read-committed"
	RepeatableRead IsolationLevel = "repeatable-read"
	Serializable   IsolationLevel = "serializable"
)

type TransactionProfile struct {
	InvariantCritical bool
	ConflictRate      string
	HotspotRisk       string
	CrossService      bool
}

type PlanResult struct {
	Level    IsolationLevel
	UseSaga  bool
	Warnings []string
}

func PlanTransaction(p TransactionProfile) PlanResult {
	r := PlanResult{Level: ReadCommitted}

	if p.InvariantCritical {
		r.Level = RepeatableRead
	}

	if p.InvariantCritical && p.ConflictRate == "high" {
		r.Level = Serializable
	}

	if p.HotspotRisk == "high" {
		r.Warnings = append(r.Warnings, "hotspot risk may serialize throughput and require redesign")
	}

	if p.CrossService {
		r.UseSaga = true
		r.Warnings = append(r.Warnings, "cross-service workflow likely needs local transaction plus saga or outbox")
	}

	if p.LevelNeedsDowngrade() {
		r.Warnings = append(r.Warnings, "lighter isolation may allow anomalies that violate the invariant")
	}

	return r
}

func (p TransactionProfile) LevelNeedsDowngrade() bool {
	return p.InvariantCritical && p.ConflictRate == "high" && p.CrossService
}

func main() {}
