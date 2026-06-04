package main

type BandwidthModel struct {
	PeakQPS       float64
	ResponseKB    float64
	CacheHitRatio float64
	CostPerGB     float64
}

type BandwidthEstimate struct {
	TotalGBPerSecond  float64
	OriginGBPerSecond float64
	MonthlyOriginGB   float64
	MonthlyOriginCost float64
}

func EstimateBandwidth(model BandwidthModel) BandwidthEstimate {
	totalGBps := model.PeakQPS * model.ResponseKB / (1024 * 1024)
	originGBps := totalGBps * (1 - model.CacheHitRatio)
	monthlyOriginGB := originGBps * 86400 * 30
	return BandwidthEstimate{
		TotalGBPerSecond:  totalGBps,
		OriginGBPerSecond: originGBps,
		MonthlyOriginGB:   monthlyOriginGB,
		MonthlyOriginCost: monthlyOriginGB * model.CostPerGB,
	}
}

func main() {}
