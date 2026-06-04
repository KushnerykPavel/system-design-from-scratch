package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type IDPlan struct {
	Name                   string `json:"name"`
	Strategy               string `json:"strategy"`
	PeakWritesPerSecond    int    `json:"peak_writes_per_second"`
	RequiresSortability    bool   `json:"requires_sortability"`
	RequiresGuessability   bool   `json:"requires_guessability"`
	Regions                int    `json:"regions"`
	ClientsGenerateIDs     bool   `json:"clients_generate_ids"`
	ClockDiscipline        bool   `json:"clock_discipline"`
	StorageLocalityMatters bool   `json:"storage_locality_matters"`
}

func ValidatePlan(plan IDPlan) []string {
	var issues []string
	switch plan.Strategy {
	case "db_sequence", "snowflake", "random":
	default:
		issues = append(issues, "strategy must be db_sequence, snowflake, or random")
	}
	if plan.PeakWritesPerSecond <= 0 {
		issues = append(issues, "peak_writes_per_second must be positive")
	}
	if plan.Strategy == "db_sequence" && plan.PeakWritesPerSecond > 50000 {
		issues = append(issues, "db_sequence may bottleneck at this write rate without range allocation or sharding")
	}
	if plan.Strategy == "snowflake" && !plan.ClockDiscipline {
		issues = append(issues, "snowflake requires clock discipline or explicit rollback handling")
	}
	if plan.Strategy == "db_sequence" && plan.ClientsGenerateIDs {
		issues = append(issues, "db_sequence does not fit offline or client-side generation")
	}
	if plan.Strategy == "random" && plan.RequiresSortability {
		issues = append(issues, "random IDs do not provide useful temporal ordering")
	}
	if plan.Strategy == "random" && plan.StorageLocalityMatters {
		issues = append(issues, "random IDs can hurt storage locality when used as primary write keys")
	}
	if plan.Strategy == "snowflake" && plan.Regions > 1 && !plan.RequiresSortability && !plan.StorageLocalityMatters {
		issues = append(issues, "snowflake may add needless clock and worker complexity when simple uniqueness is enough")
	}
	if plan.RequiresGuessability {
		issues = append(issues, "IDs should generally not be predictable from the outside")
	}
	return issues
}

func main() {
	name := flag.String("name", "orders-id-plan", "plan name")
	flag.Parse()

	plan := IDPlan{
		Name:                   *name,
		Strategy:               "snowflake",
		PeakWritesPerSecond:    250000,
		RequiresSortability:    true,
		RequiresGuessability:   false,
		Regions:                3,
		ClientsGenerateIDs:     false,
		ClockDiscipline:        true,
		StorageLocalityMatters: true,
	}

	payload := map[string]any{
		"plan":   plan,
		"issues": ValidatePlan(plan),
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
