package main

import "testing"

func TestBuildCapacitySheet(t *testing.T) {
	t.Parallel()

	input := Inputs{
		DailyActiveUsers: 100000,
		RequestsPerUser:  12,
		PeakFactor:       4,
		AverageBytes:     2048,
		WriteRatio:       0.25,
		RetentionDays:    30,
	}

	sheet := BuildCapacitySheet(input)
	if sheet.AverageQPS <= 0 {
		t.Fatalf("AverageQPS = %f, want positive", sheet.AverageQPS)
	}
	if sheet.PeakQPS <= sheet.AverageQPS {
		t.Fatalf("PeakQPS = %f, want greater than AverageQPS %f", sheet.PeakQPS, sheet.AverageQPS)
	}
	if sheet.DailyStorage <= 0 {
		t.Fatalf("DailyStorage = %f, want positive", sheet.DailyStorage)
	}
}

func TestValidateInputs(t *testing.T) {
	t.Parallel()

	input := Inputs{
		DailyActiveUsers: 0,
		RequestsPerUser:  -1,
		PeakFactor:       0.5,
		AverageBytes:     0,
		WriteRatio:       2,
		RetentionDays:    0,
	}

	if issues := ValidateInputs(input); len(issues) != 6 {
		t.Fatalf("ValidateInputs() returned %d issues, want 6: %v", len(issues), issues)
	}
}
