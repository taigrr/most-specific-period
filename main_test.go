package main

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestParsePeriods(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.May, 5, 10, 30, 0, 0, time.UTC)
	later := now.Add(30 * time.Minute)

	tests := []struct {
		name      string
		input     string
		want      []Period
		wantErr   error
		errText   string
		wantSteps []int
	}{
		{
			name: "single period",
			input: strings.Join([]string{
				"work",
				now.Format(time.RFC3339),
				later.Format(time.RFC3339),
			}, "\n"),
			want:      []Period{{Identifier: "work", StartTime: now, EndTime: later}},
			wantSteps: []int{1, 2, 0},
		},
		{
			name: "ignores blank lines",
			input: strings.Join([]string{
				"",
				"work",
				"",
				now.Format(time.RFC3339),
				later.Format(time.RFC3339),
				"",
			}, "\n"),
			want:      []Period{{Identifier: "work", StartTime: now, EndTime: later}},
			wantSteps: []int{1, 2, 0},
		},
		{
			name: "invalid start time",
			input: strings.Join([]string{
				"work",
				"not-a-time",
				later.Format(time.RFC3339),
			}, "\n"),
			errText:   `invalid start timestamp "not-a-time"`,
			wantSteps: []int{1},
		},
		{
			name: "invalid end time",
			input: strings.Join([]string{
				"work",
				now.Format(time.RFC3339),
				"not-a-time",
			}, "\n"),
			errText:   `invalid end timestamp "not-a-time"`,
			wantSteps: []int{1, 2},
		},
		{
			name: "incomplete period",
			input: strings.Join([]string{
				"work",
				now.Format(time.RFC3339),
			}, "\n"),
			wantErr:   errIncompletePeriod,
			wantSteps: []int{1, 2},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var gotSteps []int
			got, err := parsePeriods(strings.NewReader(tc.input), func(field int) {
				gotSteps = append(gotSteps, field)
			})
			if len(gotSteps) != len(tc.wantSteps) {
				t.Fatalf("expected prompt steps %v, got %v", tc.wantSteps, gotSteps)
			}
			for index, step := range tc.wantSteps {
				if gotSteps[index] != step {
					t.Fatalf("expected prompt steps %v, got %v", tc.wantSteps, gotSteps)
				}
			}
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
				return
			}
			if tc.errText != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tc.errText) {
					t.Fatalf("expected error containing %q, got %q", tc.errText, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("parsePeriods returned error: %v", err)
			}
			if len(got) != len(tc.want) {
				t.Fatalf("expected %d periods, got %d", len(tc.want), len(got))
			}
			for index := range tc.want {
				period, ok := got[index].(Period)
				if !ok {
					t.Fatalf("period %d had unexpected type %T", index, got[index])
				}
				if period != tc.want[index] {
					t.Fatalf("expected %+v, got %+v", tc.want[index], period)
				}
			}
		})
	}
}
