package main

import "fmt"

type Policy string

const (
	PolicyFIFO Policy = "fifo"
	PolicyLRU  Policy = "lru"
	PolicyLFU  Policy = "lfu"
)

type Result struct {
	Policy   Policy
	Capacity int
	Hits     int
	Misses   int
}

func (r Result) HitRate() float64 {
	total := r.Hits + r.Misses
	if total == 0 {
		return 0
	}
	return float64(r.Hits) / float64(total)
}

func Simulate(trace []string, capacity int, policy Policy) Result {
	switch policy {
	case PolicyFIFO:
		return simulateFIFO(trace, capacity)
	case PolicyLRU:
		return simulateLRU(trace, capacity)
	case PolicyLFU:
		return simulateLFU(trace, capacity)
	default:
		return Result{Policy: policy, Capacity: capacity, Misses: len(trace)}
	}
}

func main() {
	trace := []string{"a", "b", "c", "a", "b", "d", "a", "b", "c", "d"}
	for _, policy := range []Policy{PolicyFIFO, PolicyLRU, PolicyLFU} {
		result := Simulate(trace, 3, policy)
		fmt.Printf("%s hits=%d misses=%d hit_rate=%.2f\n", result.Policy, result.Hits, result.Misses, result.HitRate())
	}
}

func simulateFIFO(trace []string, capacity int) Result {
	result := Result{Policy: PolicyFIFO, Capacity: capacity}
	if capacity <= 0 {
		result.Misses = len(trace)
		return result
	}

	cache := make(map[string]struct{}, capacity)
	order := make([]string, 0, capacity)

	for _, key := range trace {
		if _, ok := cache[key]; ok {
			result.Hits++
			continue
		}

		result.Misses++
		if len(order) == capacity {
			evicted := order[0]
			order = order[1:]
			delete(cache, evicted)
		}
		cache[key] = struct{}{}
		order = append(order, key)
	}

	return result
}

func simulateLRU(trace []string, capacity int) Result {
	result := Result{Policy: PolicyLRU, Capacity: capacity}
	if capacity <= 0 {
		result.Misses = len(trace)
		return result
	}

	cache := make(map[string]int, capacity)
	order := make([]string, 0, capacity)

	for _, key := range trace {
		if idx, ok := cache[key]; ok {
			result.Hits++
			order = moveToEnd(order, idx)
			reindex(order, cache)
			continue
		}

		result.Misses++
		if len(order) == capacity {
			evicted := order[0]
			order = order[1:]
			delete(cache, evicted)
		}
		order = append(order, key)
		reindex(order, cache)
	}

	return result
}

func simulateLFU(trace []string, capacity int) Result {
	result := Result{Policy: PolicyLFU, Capacity: capacity}
	if capacity <= 0 {
		result.Misses = len(trace)
		return result
	}

	type entry struct {
		freq      int
		timestamp int
	}

	cache := make(map[string]entry, capacity)
	now := 0

	for _, key := range trace {
		now++
		if item, ok := cache[key]; ok {
			item.freq++
			item.timestamp = now
			cache[key] = item
			result.Hits++
			continue
		}

		result.Misses++
		if len(cache) == capacity {
			var victim string
			first := true
			for existingKey, item := range cache {
				if first || item.freq < cache[victim].freq || (item.freq == cache[victim].freq && item.timestamp < cache[victim].timestamp) {
					victim = existingKey
					first = false
				}
			}
			delete(cache, victim)
		}
		cache[key] = entry{freq: 1, timestamp: now}
	}

	return result
}

func moveToEnd(order []string, idx int) []string {
	key := order[idx]
	copy(order[idx:], order[idx+1:])
	order[len(order)-1] = key
	return order
}

func reindex(order []string, cache map[string]int) {
	for idx, key := range order {
		cache[key] = idx
	}
}
