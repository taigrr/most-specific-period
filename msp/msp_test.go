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

func TestMostSpecificPeriodEqualStartTimesAcrossLocations(t *testing.T) {
	utc := time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC)
	est := time.FixedZone("EST", -5*60*60)
	startInEST := utc.In(est)
	ts := utc.Add(30 * time.Minute)
	end := ts.Add(30 * time.Minute)

	id, err := MostSpecificPeriod(ts,
		TimeWindow{StartTime: utc, EndTime: end, Identifier: "A"},
		TimeWindow{StartTime: startInEST, EndTime: end, Identifier: "B"},
	)
	if err != nil {
		t.Fatalf("MostSpecificPeriod returned unexpected error: %v", err)
	}
	if id != "B" {
		t.Fatalf("expected lexicographically last identifier for equal start times, got %q", id)
	}
}

func TestMostSpecificPeriodLongDurations(t *testing.T) {
	start := time.Date(1000, time.January, 1, 0, 0, 0, 0, time.UTC)
	shortEnd := start.AddDate(400, 0, 0)
	for _, extra := range []time.Duration{time.Nanosecond, time.Second, time.Hour} {
		t.Run(extra.String(), func(t *testing.T) {
			short := TimeWindow{StartTime: start, EndTime: shortEnd, Identifier: "A"}
			long := TimeWindow{StartTime: start, EndTime: shortEnd.Add(extra), Identifier: "Z"}
			for _, periods := range [][]Period{{short, long}, {long, short}} {
				id, err := MostSpecificPeriod(start, periods...)
				if err != nil || id != "A" {
					t.Fatalf("expected shorter multi-century period A, got %q, %v", id, err)
				}
			}
		})
	}
}

func TestMostSpecificPeriodLongDurationOrdering(t *testing.T) {
	start := time.Date(1000, time.January, 1, 0, 0, 0, 1, time.UTC)
	end := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	for _, testCase := range []struct {
		name  string
		other TimeWindow
		want  string
	}{
		{"shorter beats later start", TimeWindow{StartTime: start.AddDate(1, 0, 0), EndTime: end.AddDate(2, 0, 0), Identifier: "Z"}, "A"},
		{"equal duration prefers later start", TimeWindow{StartTime: start.Add(time.Second), EndTime: end.Add(time.Second), Identifier: "Z"}, "Z"},
		{"equal duration and start prefers identifier", TimeWindow{StartTime: start, EndTime: end, Identifier: "Z"}, "Z"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			first := TimeWindow{StartTime: start, EndTime: end, Identifier: "A"}
			for _, periods := range [][]Period{{first, testCase.other}, {testCase.other, first}} {
				id, err := MostSpecificPeriod(start.AddDate(5, 0, 0), periods...)
				if err != nil || id != testCase.want {
					t.Fatalf("expected %q, got %q, %v", testCase.want, id, err)
				}
			}
		})
	}
}
