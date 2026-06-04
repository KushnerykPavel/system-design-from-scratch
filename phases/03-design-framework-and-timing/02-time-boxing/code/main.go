package main

type Segment struct {
	Name    string
	Minutes int
}

func Total(segments []Segment) int {
	total := 0
	for _, segment := range segments {
		total += segment.Minutes
	}
	return total
}

func OverBudget(total, limit int) int {
	if total <= limit {
		return 0
	}
	return total - limit
}

func Default45MinutePlan() []Segment {
	return []Segment{
		{Name: "clarify", Minutes: 6},
		{Name: "size", Minutes: 6},
		{Name: "high_level_design", Minutes: 12},
		{Name: "deep_dive", Minutes: 14},
		{Name: "wrap_up", Minutes: 7},
	}
}

func main() {}
