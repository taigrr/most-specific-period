package msp

import (
	"fmt"
	"testing"
	"time"
)

// (periods ...Period) (id string, err error) {
func TestMostSpecificPeriod(t *testing.T) {
	// use a static timestamp to make sure tests don't fail on slower systems or during a process pause
	now := time.Now()
	testCases := []struct {
		ts      time.Time
		testID  string
		result  string
		err     error
		periods []Period
	}{
		{
			testID:  "No choices",
			ts:      now,
			result:  "",
			err:     ErrNoValidPeriods,
			periods: []Period{},
		},
		{
			testID: "Two Choices, shorter is second",
			ts:     now,
			result: "B",
			err:    nil,
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
			result: "B",
			err:    nil,
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
			result: "A",
			err:    nil,
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
			result: "A",
			err:    nil,
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
			result: "B",
			err:    nil,
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
			result: "B",
			err:    nil,
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
			result: "A",
			err:    nil,
			periods: []Period{TimeWindow{
				StartTime:  now.Add(-time.Minute),
				EndTime:    now.Add(time.Minute),
				Identifier: "A",
			}},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.testID), func(t *testing.T) {
			id, err := MostSpecificPeriod(tc.ts, tc.periods...)
			if id != tc.result {
				t.Errorf("ID '%s' does not match expected '%s'", id, tc.result)
			}
			if err != tc.err {
				t.Errorf("Error '%v' does not match expected '%v'", err, tc.err)
			}
		})
	}
}
