package main

import "fmt"

type Burst struct {
	Requests        int
	ArrivalWindowMS int
	FetchLatencyMS  int
}

type OriginLoad struct {
	OriginFetches int
	Waiters       int
}

func EstimateWithoutCoalescing(b Burst) OriginLoad {
	if b.Requests <= 0 {
		return OriginLoad{}
	}
	return OriginLoad{OriginFetches: b.Requests}
}

func EstimateWithCoalescing(b Burst) OriginLoad {
	if b.Requests <= 0 {
		return OriginLoad{}
	}
	if b.FetchLatencyMS <= 0 {
		return OriginLoad{OriginFetches: 1}
	}
	if b.ArrivalWindowMS >= b.FetchLatencyMS {
		return OriginLoad{OriginFetches: 1, Waiters: b.Requests - 1}
	}

	batches := ceilDiv(b.FetchLatencyMS, max(1, b.ArrivalWindowMS))
	if batches > b.Requests {
		batches = b.Requests
	}
	return OriginLoad{
		OriginFetches: batches,
		Waiters:       b.Requests - batches,
	}
}

func main() {
	burst := Burst{Requests: 1000, ArrivalWindowMS: 50, FetchLatencyMS: 200}
	fmt.Printf("without coalescing: %+v\n", EstimateWithoutCoalescing(burst))
	fmt.Printf("with coalescing: %+v\n", EstimateWithCoalescing(burst))
}

func ceilDiv(a, b int) int {
	return (a + b - 1) / b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
