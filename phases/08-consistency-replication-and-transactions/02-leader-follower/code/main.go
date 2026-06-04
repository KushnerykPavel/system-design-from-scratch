package main

type AckPolicy string

const (
	AckLeaderLocal AckPolicy = "leader-local"
	AckLeaderPlus1 AckPolicy = "leader-plus-one-follower"
	AckMajority    AckPolicy = "majority"
)

type Topology struct {
	Followers             int
	AckPolicy             AckPolicy
	CriticalFollowerReads bool
	MaxFollowerLagSeconds int
	AutoFailover          bool
	FencingEnabled        bool
	CandidateCaughtUp     bool
}

type Assessment struct {
	SafeForCriticalReads bool
	FailoverRisk         string
	Warnings             []string
}

func Assess(t Topology) Assessment {
	a := Assessment{
		SafeForCriticalReads: !t.CriticalFollowerReads,
		FailoverRisk:         "moderate",
	}

	if t.Followers == 0 {
		a.Warnings = append(a.Warnings, "no followers means no read scale or replica failover target")
	}

	if t.CriticalFollowerReads {
		if t.MaxFollowerLagSeconds > 0 {
			a.Warnings = append(a.Warnings, "critical reads on lagging followers risk stale results")
		} else {
			a.SafeForCriticalReads = true
		}
	}

	if t.AckPolicy == AckLeaderLocal {
		a.Warnings = append(a.Warnings, "leader-local ack risks losing recent writes during failover")
		a.FailoverRisk = "high"
	}

	if t.AutoFailover && !t.FencingEnabled {
		a.Warnings = append(a.Warnings, "automatic failover without fencing risks split brain")
		a.FailoverRisk = "high"
	}

	if !t.CandidateCaughtUp {
		a.Warnings = append(a.Warnings, "promotion candidate is not caught up enough for safe failover")
		a.FailoverRisk = "high"
	}

	if t.AckPolicy == AckMajority && t.FencingEnabled && t.CandidateCaughtUp {
		a.FailoverRisk = "low"
	}

	return a
}

func main() {}
