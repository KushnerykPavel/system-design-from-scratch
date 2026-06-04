package main

type OutboxRecord struct {
	EventID          string
	AggregateID      string
	AggregateVersion int
	Published        bool
}

func ValidateRecords(records []OutboxRecord) []string {
	var problems []string
	lastVersionByAggregate := map[string]int{}
	seenEventIDs := map[string]bool{}

	for _, record := range records {
		if record.EventID == "" {
			problems = append(problems, "missing event id")
		}
		if seenEventIDs[record.EventID] {
			problems = append(problems, "duplicate event id")
		}
		seenEventIDs[record.EventID] = true

		if record.AggregateVersion <= lastVersionByAggregate[record.AggregateID] {
			problems = append(problems, "non-increasing aggregate version")
		}
		lastVersionByAggregate[record.AggregateID] = record.AggregateVersion
	}

	return problems
}

func Unpublished(records []OutboxRecord) []OutboxRecord {
	var result []OutboxRecord
	for _, record := range records {
		if !record.Published {
			result = append(result, record)
		}
	}
	return result
}

func main() {}
