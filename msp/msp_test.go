package msp

import (
	"testing"
	"time"
)

func TestGetDuration(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		testID string
		start  time.Time
		end    time.Time
		dur    time.Duration
		err    error
	}{
		{
			testID: "Normal duration",
			start:  now,
			end:    now.Add(5 * time.Minute),
			dur:    5 * time.Minute,
			err:    nil,
		},
		{
			testID: "Zero duration",
			start:  now,
			end:    now,
			dur:    0,
			err:    nil,
		},
		{
			testID: "Start after end",
			start:  now.Add(5 * time.Minute),
			end:    now,
			dur:    -5 * time.Minute,
			err:    ErrEndAfterStart,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testID, func(t *testing.T) {
			dur, err := GetDuration(tc.start, tc.end)
			if dur != tc.dur {
				t.Errorf("Duration %v does not match expected %v", dur, tc.dur)
			}
			if err != tc.err {
				t.Errorf("Error %v does not match expected %v", err, tc.err)
			}
		})
	}
}

func TestValidTimePeriods(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		testID  string
		ts      time.Time
		periods []Period
		count   int
	}{
		{
			testID:  "No periods",
			ts:      now,
			periods: []Period{},
			count:   0,
		},
		{
			testID: "One valid period",
			ts:     now,
			periods: []Period{
				TimeWindow{
					StartTime:  now.Add(-time.Minute),
					EndTime:    now.Add(time.Minute),
					Identifier: "A",
				},
			},
			count: 1,
		},
		{
			testID: "Period in the past",
			ts:     now,
			periods: []Period{
				TimeWindow{
					StartTime:  now.Add(-2 * time.Minute),
					EndTime:    now.Add(-time.Minute),
					Identifier: "A",
				},
			},
			count: 0,
		},
		{
			testID: "Period in the future",
			ts:     now,
			periods: []Period{
				TimeWindow{
					StartTime:  now.Add(time.Minute),
					EndTime:    now.Add(2 * time.Minute),
					Identifier: "A",
				},
			},
			count: 0,
		},
		{
			testID: "Timestamp equals start time (inclusive)",
			ts:     now,
			periods: []Period{
				TimeWindow{
					StartTime:  now,
					EndTime:    now.Add(time.Minute),
					Identifier: "A",
				},
			},
			count: 1,
		},
		{
			testID: "Timestamp equals end time (exclusive)",
			ts:     now,
			periods: []Period{
				TimeWindow{
					StartTime:  now.Add(-time.Minute),
					EndTime:    now,
					Identifier: "A",
				},
			},
			count: 0,
		},
		{
			testID: "Mixed valid and invalid",
			ts:     now,
			periods: []Period{
				TimeWindow{
					StartTime:  now.Add(-time.Minute),
					EndTime:    now.Add(time.Minute),
					Identifier: "A",
				},
				TimeWindow{
					StartTime:  now.Add(-3 * time.Minute),
					EndTime:    now.Add(-2 * time.Minute),
					Identifier: "B",
				},
				TimeWindow{
					StartTime:  now.Add(-5 * time.Minute),
					EndTime:    now.Add(5 * time.Minute),
					Identifier: "C",
				},
			},
			count: 2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testID, func(t *testing.T) {
			valid := ValidTimePeriods(tc.ts, tc.periods...)
			if len(valid) != tc.count {
				t.Errorf("Got %d valid periods, expected %d", len(valid), tc.count)
			}
		})
	}
}

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
		t.Run(tc.testID, func(t *testing.T) {
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
