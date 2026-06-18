package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
)

const numBuckets = 10000

// Treatment defines a bucket range and its name within an experiment.
type Treatment struct {
	Name        string `json:"name"`
	BucketStart int    `json:"bucket_start"`
	BucketEnd   int    `json:"bucket_end"`
}

// Experiment defines an A/B experiment with its treatments.
type Experiment struct {
	ID         string      `json:"id"`
	Treatments []Treatment `json:"treatments"`
}

// bucket computes a deterministic bucket index for a user within an experiment.
func bucket(userID, experimentID string) int {
	h := fnv.New32a()
	_, _ = fmt.Fprintf(h, "%s:%s", userID, experimentID)
	return int(h.Sum32()) % numBuckets
}

// Assign returns the treatment name for userID in the given experiment.
// Returns "unassigned" if the bucket falls outside all treatment ranges.
func Assign(userID string, exp Experiment) string {
	b := bucket(userID, exp.ID)
	for _, t := range exp.Treatments {
		if b >= t.BucketStart && b < t.BucketEnd {
			return t.Name
		}
	}
	return "unassigned"
}

// AssignAll resolves all active experiments for a user in one pass.
func AssignAll(userID string, experiments []Experiment) map[string]string {
	result := make(map[string]string, len(experiments))
	for _, exp := range experiments {
		result[exp.ID] = Assign(userID, exp)
	}
	return result
}

func main() {
	experiments := []Experiment{
		{
			ID: "rec-algo-v2",
			Treatments: []Treatment{
				{Name: "control", BucketStart: 0, BucketEnd: 5000},
				{Name: "treatment_A", BucketStart: 5000, BucketEnd: 10000},
			},
		},
		{
			ID: "thumbnail-size",
			Treatments: []Treatment{
				{Name: "control", BucketStart: 0, BucketEnd: 6000},
				{Name: "large", BucketStart: 6000, BucketEnd: 10000},
			},
		},
	}

	users := []string{"user-001", "user-002", "user-003"}
	results := make([]map[string]any, 0, len(users))
	for _, u := range users {
		assignments := AssignAll(u, experiments)
		results = append(results, map[string]any{
			"user_id":     u,
			"assignments": assignments,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
