package main

import "fmt"

// Partition represents a single Kafka partition with its offset state.
type Partition struct {
	ID              int
	LatestOffset    int64 // offset of the most recently produced record
	CommittedOffset int64 // offset the consumer group has committed
}

// ConsumerLag returns the number of unconsumed records in a single partition.
// Lag = LatestOffset - CommittedOffset.
// Returns 0 if the consumer has committed ahead of (or equal to) the latest offset.
func ConsumerLag(p Partition) int64 {
	lag := p.LatestOffset - p.CommittedOffset
	if lag < 0 {
		return 0
	}
	return lag
}

// GroupLag returns the total consumer lag across all partitions in a consumer group.
func GroupLag(partitions []Partition) int64 {
	var total int64
	for _, p := range partitions {
		total += ConsumerLag(p)
	}
	return total
}

// IsLagAlarm returns true if the total lag exceeds the given threshold.
// Example: threshold = 100_000 means "alert if more than 100K unprocessed records".
func IsLagAlarm(totalLag int64, threshold int64) bool {
	return totalLag > threshold
}

// SimulateConsumerBehind simulates one step of consumer processing.
// Each call advances LatestOffset by messageRate (new records produced)
// and CommittedOffset by processRate (records consumed in this step).
// CommittedOffset never exceeds LatestOffset.
func SimulateConsumerBehind(partitions []Partition, messageRate, processRate int) []Partition {
	updated := make([]Partition, len(partitions))
	for i, p := range partitions {
		p.LatestOffset += int64(messageRate)
		p.CommittedOffset += int64(processRate)
		if p.CommittedOffset > p.LatestOffset {
			p.CommittedOffset = p.LatestOffset
		}
		updated[i] = p
	}
	return updated
}

func main() {
	// Simulate a 3-partition consumer group.
	// The consumer starts caught up, then falls behind as message rate > process rate,
	// then recovers as process rate increases.
	partitions := []Partition{
		{ID: 0, LatestOffset: 1000, CommittedOffset: 1000},
		{ID: 1, LatestOffset: 1200, CommittedOffset: 1200},
		{ID: 2, LatestOffset: 980, CommittedOffset: 980},
	}

	const alarmThreshold int64 = 5000

	fmt.Println("Step | TotalLag | Alarm | Phase")
	fmt.Println("-----|----------|-------|------")

	// Phase 1 (steps 1-5): consumer falls behind (produce 500/step, process 200/step)
	for step := 1; step <= 5; step++ {
		partitions = SimulateConsumerBehind(partitions, 500, 200)
		lag := GroupLag(partitions)
		alarm := IsLagAlarm(lag, alarmThreshold)
		fmt.Printf("  %2d | %8d | %-5v | falling behind\n", step, lag, alarm)
	}

	// Phase 2 (steps 6-10): consumer catches up (produce 500/step, process 800/step)
	for step := 6; step <= 10; step++ {
		partitions = SimulateConsumerBehind(partitions, 500, 800)
		lag := GroupLag(partitions)
		alarm := IsLagAlarm(lag, alarmThreshold)
		fmt.Printf("  %2d | %8d | %-5v | catching up\n", step, lag, alarm)
	}
}
