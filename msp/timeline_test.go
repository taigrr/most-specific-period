package msp

import (
	"fmt"
	"testing"
	"time"
)

// (periods ...Period) (id string, err error) {
func TestGenerateTime(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		ts      time.Time
		testID  string
		result  []string
		periods []Period
	}{
		{
			testID:  "No choices",
			ts:      now,
			result:  []string{},
			periods: []Period{},
		},
		{
			testID: "Two Choices, shorter is second",
			ts:     now,
			result: []string{
				fmt.Sprintf("A\t%s\t%s\n", now.Add(-5*time.Minute), now.Add(-2*time.Minute)),
				fmt.Sprintf("B\t%s\t%s\n", now.Add(-2*time.Minute), now.Add(time.Minute)),
			},
			periods: []Period{
				TimeWindow{
					StartTime:  now.Add(-5 * time.Minute),
					EndTime:    now.Add(time.Minute),
					Identifier: "A",
				},
				TimeWindow{
					StartTime:  now.Add(-2 * time.Minute),
					EndTime:    now.Add(time.Minute),
					Identifier: "B",
				},
			},
		},
		{
			testID: "Two Choices, one is a year, other a minute",
			ts:     now,
			result: []string{
				fmt.Sprintf("A\t%s\t%s\n", now.Add(-1*time.Hour*24*365), now.Add(-5*time.Minute)),
				fmt.Sprintf("B\t%s\t%s\n", now.Add(-5*time.Minute), now.Add(time.Minute)),
			},
			periods: []Period{
				TimeWindow{
					StartTime:  now.Add(-1 * time.Hour * 24 * 365),
					EndTime:    now.Add(time.Minute),
					Identifier: "A",
				},
				TimeWindow{
					StartTime:  now.Add(-5 * time.Minute),
					EndTime:    now.Add(time.Minute),
					Identifier: "B",
				},
			},
		},
		{
			testID: "Two Choices, shorter is first",
			ts:     now,
			result: []string{
				fmt.Sprintf("B\t%s\t%s\n", now.Add(-5*time.Minute), now.Add(-2*time.Minute)),
				fmt.Sprintf("A\t%s\t%s\n", now.Add(-2*time.Minute), now.Add(time.Minute)),
			},
			periods: []Period{
				TimeWindow{
					StartTime:  now.Add(-2 * time.Minute),
					EndTime:    now.Add(time.Minute),
					Identifier: "A",
				},
				TimeWindow{
					StartTime:  now.Add(-5 * time.Minute),
					EndTime:    now.Add(time.Minute),
					Identifier: "B",
				},
			},
		},
		{
			testID: "Two Choices, one in the past",
			ts:     now,
			result: []string{
				fmt.Sprintf("B\t%s\t%s\n", now.Add(-2*time.Minute), now.Add(-time.Minute)),
				fmt.Sprintf("A\t%s\t%s\n", now.Add(-time.Minute), now.Add(time.Minute)),
			},
			periods: []Period{
				TimeWindow{
					StartTime:  now.Add(-time.Minute),
					EndTime:    now.Add(time.Minute),
					Identifier: "A",
				},
				TimeWindow{
					StartTime:  now.Add(-2 * time.Minute),
					EndTime:    now.Add(-time.Minute),
					Identifier: "B",
				},
			},
		},
		{
			testID: "Two Choices, one invalid",
			ts:     now,
			result: []string{
				fmt.Sprintf("B\t%s\t%s\n", now.Add(-2*time.Minute), now.Add(time.Minute)),
			},
			periods: []Period{
				TimeWindow{
					StartTime:  now.Add(time.Minute),
					EndTime:    now.Add(-time.Minute),
					Identifier: "A",
				},
				TimeWindow{
					StartTime:  now.Add(-2 * time.Minute),
					EndTime:    now.Add(time.Minute),
					Identifier: "B",
				},
			},
		},
		{
			testID: "Two Choices, Identical periods",
			ts:     now,
			result: []string{
				fmt.Sprintf("B\t%s\t%s\n", now.Add(-time.Minute), now.Add(time.Minute)),
			},
			periods: []Period{
				TimeWindow{
					StartTime:  now.Add(-time.Minute),
					EndTime:    now.Add(time.Minute),
					Identifier: "A",
				},
				TimeWindow{
					StartTime:  now.Add(-time.Minute),
					EndTime:    now.Add(time.Minute),
					Identifier: "B",
				},
			},
		},
		{
			testID: "One choice",
			ts:     now,
			result: []string{
				fmt.Sprintf("A\t%s\t%s\n", now.Add(-time.Minute), now.Add(time.Minute)),
			},
			periods: []Period{TimeWindow{
				StartTime:  now.Add(-time.Minute),
				EndTime:    now.Add(time.Minute),
				Identifier: "A",
			}},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.testID), func(t *testing.T) {
			timeline := GenerateTimeline(tc.periods...)
			if len(timeline) != len(tc.result) {
				t.Fatalf("Time line had %d results, expected %d", len(timeline), len(tc.result))
			}
			for idx, period := range timeline {
				if period != tc.result[idx] {
					t.Errorf("Expected:\t%s\t\tHad:\t%s", period, tc.result[idx])
				}
			}
		})
	}
}
