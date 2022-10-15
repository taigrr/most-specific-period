package msp

import (
	"fmt"
	"testing"
	"time"
)

func slicesEqual[K comparable](x []K, y []K) bool {
	if len(x) != len(y) {
		return false
	}
	for i, w := range x {
		if w != y[i] {
			return false
		}
	}
	return true
}

func TestGetChangeOvers(t *testing.T) {
	// use a static timestamp to make sure tests don't fail on slower systems or during a process pause
	now := time.Now()
	testCases := []struct {
		ts      time.Time
		testID  string
		result  []time.Time
		periods []Period
	}{
		{
			testID:  "No choices",
			ts:      now,
			result:  []time.Time{},
			periods: []Period{},
		},
		{
			testID: "Two Choices, shorter is second",
			ts:     now,
			result: []time.Time{now.Add(-5 * time.Minute), now.Add(-2 * time.Minute), now.Add(time.Minute)},
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
			result: []time.Time{now.Add(-1 * time.Hour * 24 * 365), now.Add(-5 * time.Minute), now.Add(time.Minute)},
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
			result: []time.Time{now.Add(-5 * time.Minute), now.Add(-2 * time.Minute), now.Add(time.Minute)},
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
			result: []time.Time{now.Add(-2 * time.Minute), now.Add(-time.Minute), now.Add(time.Minute)},
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
			result: []time.Time{now.Add(-2 * time.Minute), now.Add(time.Minute)},
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
			result: []time.Time{now.Add(-time.Minute), now.Add(time.Minute)},
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
			result: []time.Time{now.Add(-time.Minute), now.Add(time.Minute)},
			periods: []Period{TimeWindow{
				StartTime:  now.Add(-time.Minute),
				EndTime:    now.Add(time.Minute),
				Identifier: "A",
			}},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.testID), func(t *testing.T) {
			changeovers := GetChangeOvers(tc.periods...)
			if !slicesEqual(changeovers, tc.result) {
				t.Errorf("Expected %v but got %v", tc.result, changeovers)
			}
		})
	}
}

func TestFlattenPeriods(t *testing.T) {
	// use a static timestamp to make sure tests don't fail on slower systems or during a process pause
	now := time.Now()
	testCases := []struct {
		ts      time.Time
		testID  string
		result  []string
		err     error
		periods []Period
	}{
		{
			testID:  "No choices",
			ts:      now,
			result:  []string{},
			err:     ErrNoValidPeriods,
			periods: []Period{},
		},
		{
			testID: "Two Choices, shorter is second",
			ts:     now,
			result: []string{"A", "B"},
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
			result: []string{"A", "B"},
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
			result: []string{"B", "A"},
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
			result: []string{"B", "A"},
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
			result: []string{"B"},
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
			result: []string{"B"},
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
			testID: "Triple Nested Periods",
			ts:     now,
			result: []string{"A", "B", "C", "B", "A"},
			err:    nil,
			periods: []Period{
				TimeWindow{
					StartTime:  now.Add(-15 * time.Minute),
					EndTime:    now.Add(15 * time.Minute),
					Identifier: "A",
				},
				TimeWindow{
					StartTime:  now.Add(-5 * time.Minute),
					EndTime:    now.Add(5 * time.Minute),
					Identifier: "C",
				},
				TimeWindow{
					StartTime:  now.Add(-10 * time.Minute),
					EndTime:    now.Add(10 * time.Minute),
					Identifier: "B",
				},
			},
		},
		{
			testID: "One choice",
			ts:     now,
			result: []string{"A"},
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
			changeovers := FlattenPeriods(tc.periods...)
			if !slicesEqual(changeovers, tc.result) {
				t.Errorf("Expected %v but got %v", tc.result, changeovers)
			}
		})
	}
}

func TestGetNextChangeOver(t *testing.T) {
	// use a static timestamp to make sure tests don't fail on slower systems or during a process pause
	now := time.Now()
	testCases := []struct {
		ts      time.Time
		testID  string
		result  time.Time
		err     error
		periods []Period
	}{
		{
			testID:  "No choices",
			ts:      now,
			result:  time.Time{},
			err:     ErrNoNextChangeover,
			periods: []Period{},
		},
		{
			testID: "Two Choices, shorter is second",
			ts:     now,
			result: now.Add(time.Minute),
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
			result: now.Add(time.Minute),
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
			result: now.Add(time.Minute),
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
			result: now.Add(time.Minute),
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
			result: now.Add(time.Minute),
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
			result: now.Add(time.Minute),
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
			result: now.Add(time.Minute),
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
			ts, err := GetNextChangeOver(now, tc.periods...)
			if tc.err != err {
				t.Errorf("Error %v does not match expected %v", tc.err, err)
			}
			if ts != tc.result {
				t.Errorf("Got %v but expected %v", ts, tc.result)
			}
		})
	}
}
