package msp

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestEmptyIdentifierChangeovers(t *testing.T) {
	start := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	for _, separated := range []bool{false, true} {
		name := "single window"
		periods := []Period{TimeWindow{StartTime: start, EndTime: end}}
		wantTimes := []time.Time{start, end}
		wantTimeline := []Period{TimeWindow{StartTime: start, EndTime: end}}
		wantIDs := []string{""}
		if separated {
			name = "separated windows"
			nextStart, nextEnd := end.Add(time.Hour), end.Add(2*time.Hour)
			periods = append(periods, TimeWindow{StartTime: nextStart, EndTime: nextEnd})
			wantTimes = append(wantTimes, nextStart, nextEnd)
			wantTimeline = append(wantTimeline, TimeWindow{StartTime: nextStart, EndTime: nextEnd})
			wantIDs = append(wantIDs, "")
		}
		t.Run(name, func(t *testing.T) {
			if got := GetChangeOvers(periods...); !reflect.DeepEqual(got, wantTimes) {
				t.Errorf("GetChangeOvers = %v, want %v", got, wantTimes)
			}
			if got := GenerateTimeline(periods...); !reflect.DeepEqual(got, wantTimeline) {
				t.Errorf("GenerateTimeline = %v, want %v", got, wantTimeline)
			}
			if got := FlattenPeriods(periods...); !reflect.DeepEqual(got, wantIDs) {
				t.Errorf("FlattenPeriods = %q, want %q", got, wantIDs)
			}
			for index, timestamp := range wantTimes {
				previous := timestamp.Add(-time.Nanosecond)
				if index > 0 {
					previous = wantTimes[index-1]
				}
				got, err := GetNextChangeOver(previous, periods...)
				if err != nil || !got.Equal(timestamp) {
					t.Errorf("GetNextChangeOver(%v) = %v, %v; want %v", previous, got, err, timestamp)
				}
			}
			if _, err := GetNextChangeOver(wantTimes[len(wantTimes)-1], periods...); !errors.Is(err, ErrNoNextChangeover) {
				t.Errorf("after final boundary: %v", err)
			}
		})
	}
}
