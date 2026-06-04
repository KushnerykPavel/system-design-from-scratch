package main

import "sort"

type Assignment struct {
	Member     string
	Partitions []int
}

func AssignEvenly(partitions int, members []string) []Assignment {
	if partitions <= 0 || len(members) == 0 {
		return nil
	}

	sorted := append([]string(nil), members...)
	sort.Strings(sorted)

	assignments := make([]Assignment, len(sorted))
	for i, member := range sorted {
		assignments[i] = Assignment{Member: member}
	}

	for partition := 0; partition < partitions; partition++ {
		idx := partition % len(assignments)
		assignments[idx].Partitions = append(assignments[idx].Partitions, partition)
	}

	return assignments
}

func MaxLoad(assignments []Assignment) int {
	max := 0
	for _, assignment := range assignments {
		if len(assignment.Partitions) > max {
			max = len(assignment.Partitions)
		}
	}
	return max
}

func NeedsMorePartitions(partitions int, members int) bool {
	return members > partitions
}

func main() {}
