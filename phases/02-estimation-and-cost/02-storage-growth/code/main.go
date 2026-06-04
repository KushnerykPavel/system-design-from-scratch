package main

type StoragePlan struct {
	EventsPerDay       float64
	BytesPerEvent      float64
	ReplicationFactor  float64
	RetentionDays      float64
	IndexOverheadRatio float64
}

type StorageEstimate struct {
	DailyRawGB      float64
	DailyDurableGB  float64
	RetainedTotalGB float64
}

func EstimateStorage(plan StoragePlan) StorageEstimate {
	dailyRawGB := plan.EventsPerDay * plan.BytesPerEvent / (1024 * 1024 * 1024)
	dailyDurableGB := dailyRawGB * plan.ReplicationFactor * (1 + plan.IndexOverheadRatio)
	return StorageEstimate{
		DailyRawGB:      dailyRawGB,
		DailyDurableGB:  dailyDurableGB,
		RetainedTotalGB: dailyDurableGB * plan.RetentionDays,
	}
}

func main() {}
