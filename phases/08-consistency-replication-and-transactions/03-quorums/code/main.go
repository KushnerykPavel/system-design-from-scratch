package main

type Plan struct {
	Replicas     int
	ReadQuorum   int
	WriteQuorum  int
	LatencyBias  string
	ConflictRisk string
}

type Evaluation struct {
	Overlaps     bool
	Availability string
	Warnings     []string
}

func Evaluate(p Plan) Evaluation {
	e := Evaluation{
		Overlaps: p.ReadQuorum+p.WriteQuorum > p.Replicas,
	}

	if p.ReadQuorum <= 0 || p.WriteQuorum <= 0 || p.Replicas <= 0 {
		e.Warnings = append(e.Warnings, "replicas and quorums must be positive")
		return e
	}

	if p.ReadQuorum > p.Replicas || p.WriteQuorum > p.Replicas {
		e.Warnings = append(e.Warnings, "quorum cannot exceed replica set size")
	}

	if p.ReadQuorum == 1 && p.WriteQuorum == 1 {
		e.Availability = "very high"
		e.Warnings = append(e.Warnings, "tiny quorums maximize availability but make stale reads and divergence more likely")
		return e
	}

	if e.Overlaps {
		e.Availability = "balanced"
	} else {
		e.Availability = "high"
		e.Warnings = append(e.Warnings, "non-overlapping quorums weaken freshness confidence")
	}

	if p.LatencyBias == "read" && p.ReadQuorum > 1 {
		e.Warnings = append(e.Warnings, "read-biased latency goal conflicts with larger read quorum")
	}
	if p.ConflictRisk == "high" && !e.Overlaps {
		e.Warnings = append(e.Warnings, "high conflict risk usually needs stronger overlap or different ownership")
	}

	return e
}

func main() {}
