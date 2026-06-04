package main

type QueueModel struct {
	ArrivalRate         float64
	ServiceRate         float64
	BurstSeconds        float64
	RecoveryServiceRate float64
	PostBurstArrival    float64
}

type QueueEstimate struct {
	BacklogGrowthPerSec float64
	BacklogItems        float64
	DrainSeconds        float64
}

func EstimateQueue(model QueueModel) QueueEstimate {
	growth := model.ArrivalRate - model.ServiceRate
	if growth < 0 {
		growth = 0
	}
	backlog := growth * model.BurstSeconds
	spareCapacity := model.RecoveryServiceRate - model.PostBurstArrival
	drainSeconds := 0.0
	if backlog > 0 && spareCapacity > 0 {
		drainSeconds = backlog / spareCapacity
	}
	return QueueEstimate{
		BacklogGrowthPerSec: growth,
		BacklogItems:        backlog,
		DrainSeconds:        drainSeconds,
	}
}

func main() {}
