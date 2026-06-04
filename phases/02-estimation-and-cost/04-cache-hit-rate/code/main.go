package main

type CacheModel struct {
	ReadQPS           float64
	HitRatio          float64
	MissLatencyMillis float64
}

type CacheEstimate struct {
	CacheServedQPS      float64
	OriginQPS           float64
	AverageLatencyDelta float64
}

func EstimateCache(model CacheModel) CacheEstimate {
	cacheServed := model.ReadQPS * model.HitRatio
	origin := model.ReadQPS - cacheServed
	return CacheEstimate{
		CacheServedQPS:      cacheServed,
		OriginQPS:           origin,
		AverageLatencyDelta: (origin / model.ReadQPS) * model.MissLatencyMillis,
	}
}

func main() {}
