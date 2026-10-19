package msp

import (
	"testing"
	"time"
)

func TestMostSpecificPeriodLatestStartTieBreak(t *testing.T) {
	t.Parallel()

	start := time.Date(2026, time.October, 19, 8, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		periods []Period
		want    string
	}{
		{
			name: "latest start beats identifier ordering",
			periods: []Period{
				TimeWindow{StartTime: start, EndTime: start.Add(2 * time.Hour), Identifier: "z-older"},
				TimeWindow{StartTime: start.Add(time.Hour), EndTime: start.Add(3 * time.Hour), Identifier: "a-newer"},
			},
			want: "a-newer",
		},
		{
			name: "shortest duration beats latest start",
			periods: []Period{
				TimeWindow{StartTime: start, EndTime: start.Add(2 * time.Hour), Identifier: "shorter"},
				TimeWindow{StartTime: start.Add(time.Hour), EndTime: start.Add(4 * time.Hour), Identifier: "longer"},
			},
			want: "shorter",
		},
		{
			name: "latest instant beats local clock ordering",
			periods: []Period{
				TimeWindow{StartTime: start.In(time.FixedZone("east", 2*60*60)), EndTime: start.Add(2 * time.Hour), Identifier: "z-older"},
				TimeWindow{StartTime: start.Add(time.Hour).In(time.FixedZone("west", -5*60*60)), EndTime: start.Add(3 * time.Hour), Identifier: "a-newer"},
			},
			want: "a-newer",
		},
		{
			name: "identifier only breaks ties among latest starts",
			periods: []Period{
				TimeWindow{StartTime: start, EndTime: start.Add(2 * time.Hour), Identifier: "z-older"},
				TimeWindow{StartTime: start.Add(time.Hour), EndTime: start.Add(3 * time.Hour), Identifier: "a-newer"},
				TimeWindow{StartTime: start.Add(time.Hour), EndTime: start.Add(3 * time.Hour), Identifier: "b-newer"},
			},
			want: "b-newer",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// Exercise every input permutation: tie-breaking must not depend
			// on caller ordering or reorder the caller's slice.
			var checkPermutations func(int)
			checkPermutations = func(index int) {
				if index < len(test.periods) {
					for candidate := index; candidate < len(test.periods); candidate++ {
						test.periods[index], test.periods[candidate] = test.periods[candidate], test.periods[index]
						checkPermutations(index + 1)
						test.periods[index], test.periods[candidate] = test.periods[candidate], test.periods[index]
					}
					return
				}
				before := append([]Period(nil), test.periods...)
				got, err := MostSpecificPeriod(start.Add(90*time.Minute), test.periods...)
				if err != nil || got != test.want {
					t.Fatalf("MostSpecificPeriod(%v) = %q, %v; want %q, nil", before, got, err, test.want)
				}
				for position := range before {
					if test.periods[position] != before[position] {
						t.Fatalf("input period %d changed", position)
					}
				}
			}
			checkPermutations(0)
		})
	}
}
