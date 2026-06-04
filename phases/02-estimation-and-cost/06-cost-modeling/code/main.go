package main

type CostItem struct {
	Name     string
	UnitCost float64
	Units    float64
}

func TotalCost(items []CostItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.UnitCost * item.Units
	}
	return total
}

func MostExpensiveItem(items []CostItem) string {
	maxName := ""
	maxValue := -1.0
	for _, item := range items {
		value := item.UnitCost * item.Units
		if value > maxValue {
			maxValue = value
			maxName = item.Name
		}
	}
	return maxName
}

func main() {}
