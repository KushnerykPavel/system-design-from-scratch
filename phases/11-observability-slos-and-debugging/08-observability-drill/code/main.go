package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type DrillAnswer struct {
	HasSLI        bool `json:"has_sli"`
	HasSLO        bool `json:"has_slo"`
	HasMetrics    bool `json:"has_metrics"`
	HasLogsTraces bool `json:"has_logs_traces"`
	HasDashboards bool `json:"has_dashboards"`
	HasAlerts     bool `json:"has_alerts"`
	HasRunbook    bool `json:"has_runbook"`
	HasDebugStory bool `json:"has_debug_story"`
	HasTradeoffs  bool `json:"has_tradeoffs"`
}

type DrillAssessment struct {
	Score   int      `json:"score"`
	Level   string   `json:"level"`
	Missing []string `json:"missing"`
}

func AssessDrill(answer DrillAnswer) DrillAssessment {
	score := 0
	missing := make([]string, 0, 9)

	add := func(ok bool, label string, points int) {
		if ok {
			score += points
			return
		}
		missing = append(missing, label)
	}

	add(answer.HasSLI, "user-meaningful SLI", 2)
	add(answer.HasSLO, "explicit SLO target", 2)
	add(answer.HasMetrics, "diagnostic metrics", 2)
	add(answer.HasLogsTraces, "logs and traces strategy", 2)
	add(answer.HasDashboards, "dashboard and cardinality plan", 2)
	add(answer.HasAlerts, "paging policy", 2)
	add(answer.HasRunbook, "runbook / first response workflow", 2)
	add(answer.HasDebugStory, "debugging narrative", 2)
	add(answer.HasTradeoffs, "trade-off discussion", 2)

	level := "weak"
	switch {
	case score >= 16:
		level = "strong"
	case score >= 11:
		level = "developing"
	}

	return DrillAssessment{
		Score:   score,
		Level:   level,
		Missing: missing,
	}
}

func main() {
	answer := DrillAnswer{
		HasSLI:        true,
		HasSLO:        true,
		HasMetrics:    true,
		HasLogsTraces: true,
		HasDashboards: true,
		HasAlerts:     true,
		HasRunbook:    false,
		HasDebugStory: true,
		HasTradeoffs:  false,
	}

	encoded, err := json.MarshalIndent(AssessDrill(answer), "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(encoded))
}
