package main

type IndexPlan struct {
	BaseWriteCost   int
	IndexCount      int
	ReplicationFactor int
}

func (p IndexPlan) WriteUnitsPerLogicalWrite() int {
	if p.ReplicationFactor < 1 {
		return 0
	}
	return p.BaseWriteCost * (1 + p.IndexCount) * p.ReplicationFactor
}

func main() {}
