package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

const secondsPerDay = 86400
const daysPerYear = 365

type Inputs struct {
	DailyActiveUsers int64   `json:"daily_active_users"`
	RequestsPerUser  float64 `json:"requests_per_user"`
	PeakFactor       float64 `json:"peak_factor"`
	AverageBytes     int64   `json:"average_bytes"`
	WriteRatio       float64 `json:"write_ratio"`
	RetentionDays    int64   `json:"retention_days"`
}

type CapacitySheet struct {
	AverageQPS       float64 `json:"average_qps"`
	PeakQPS          float64 `json:"peak_qps"`
	DailyWrites      float64 `json:"daily_writes"`
	DailyStorage     float64 `json:"daily_storage_bytes"`
	AnnualStorage    float64 `json:"annual_storage_bytes"`
	PeakBandwidthBps float64 `json:"peak_bandwidth_bytes_per_second"`
}

func BuildCapacitySheet(in Inputs) CapacitySheet {
	totalRequestsPerDay := float64(in.DailyActiveUsers) * in.RequestsPerUser
	averageQPS := totalRequestsPerDay / secondsPerDay
	peakQPS := averageQPS * in.PeakFactor
	dailyWrites := totalRequestsPerDay * in.WriteRatio
	dailyStorage := dailyWrites * float64(in.AverageBytes)

	return CapacitySheet{
		AverageQPS:       averageQPS,
		PeakQPS:          peakQPS,
		DailyWrites:      dailyWrites,
		DailyStorage:     dailyStorage,
		AnnualStorage:    dailyStorage * daysPerYear,
		PeakBandwidthBps: peakQPS * float64(in.AverageBytes),
	}
}

func ValidateInputs(in Inputs) []string {
	var issues []string
	if in.DailyActiveUsers <= 0 {
		issues = append(issues, "daily_active_users must be positive")
	}
	if in.RequestsPerUser <= 0 {
		issues = append(issues, "requests_per_user must be positive")
	}
	if in.PeakFactor < 1 {
		issues = append(issues, "peak_factor must be at least 1")
	}
	if in.AverageBytes <= 0 {
		issues = append(issues, "average_bytes must be positive")
	}
	if in.WriteRatio < 0 || in.WriteRatio > 1 {
		issues = append(issues, "write_ratio must be between 0 and 1")
	}
	if in.RetentionDays <= 0 {
		issues = append(issues, "retention_days must be positive")
	}
	return issues
}

func main() {
	var path string
	flag.StringVar(&path, "input", "", "path to an inputs JSON file")
	flag.Parse()

	if path == "" {
		fmt.Fprintln(os.Stderr, "missing -input path")
		os.Exit(2)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read input: %v\n", err)
		os.Exit(1)
	}

	var in Inputs
	if err := json.Unmarshal(data, &in); err != nil {
		fmt.Fprintf(os.Stderr, "decode input: %v\n", err)
		os.Exit(1)
	}

	issues := ValidateInputs(in)
	if len(issues) > 0 {
		for _, issue := range issues {
			fmt.Println(issue)
		}
		os.Exit(1)
	}

	sheet := BuildCapacitySheet(in)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(sheet)
}
