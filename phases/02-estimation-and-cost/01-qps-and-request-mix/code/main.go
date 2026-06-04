package main

type TrafficModel struct {
	DAU                   float64
	RequestsPerUserPerDay float64
	ReadRatio             float64
	PeakFactor            float64
}

type TrafficEstimate struct {
	AverageQPS   float64
	PeakQPS      float64
	PeakReadQPS  float64
	PeakWriteQPS float64
}

func EstimateTraffic(model TrafficModel) TrafficEstimate {
	averageQPS := (model.DAU * model.RequestsPerUserPerDay) / 86400.0
	peakQPS := averageQPS * model.PeakFactor
	return TrafficEstimate{
		AverageQPS:   averageQPS,
		PeakQPS:      peakQPS,
		PeakReadQPS:  peakQPS * model.ReadRatio,
		PeakWriteQPS: peakQPS * (1 - model.ReadRatio),
	}
}

func main() {}
