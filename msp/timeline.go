package msp

import (
	"fmt"
	"time"
)

// TimeWindow is a concrete implementation of the Period interface.
type TimeWindow struct {
	StartTime  time.Time
	EndTime    time.Time
	Identifier string
}

// GetIdentifier returns the period's identifier string.
func (p TimeWindow) GetIdentifier() string {
	return p.Identifier
}

// GetEndTime returns the period's exclusive end time.
func (p TimeWindow) GetEndTime() time.Time {
	return p.EndTime
}

// GetStartTime returns the period's inclusive start time.
func (p TimeWindow) GetStartTime() time.Time {
	return p.StartTime
}

// String returns a tab-separated representation of the time window.
func (t TimeWindow) String() string {
	return fmt.Sprintf("%s\t%s\t%s",
		t.GetIdentifier(),
		t.GetStartTime(),
		t.GetEndTime())
}

// GenerateTimeline produces a flattened timeline of non-overlapping periods
// by splitting overlapping input periods at changeover points.
func GenerateTimeline(periods ...Period) (out []Period) {
	if len(periods) == 0 {
		return out
	}
	periodsByID := make(map[string]Period)
	ids := FlattenPeriods(periods...)
	for _, val := range periods {
		id := val.GetIdentifier()
		periodsByID[id] = val
	}
	start := periodsByID[ids[0]].GetStartTime()
	for _, val := range ids {
		next, err := GetNextChangeOver(start, periods...)
		if err == nil {
			if next.Equal(periodsByID[val].GetStartTime()) {
				start = periodsByID[val].GetStartTime()
				next = periodsByID[val].GetEndTime()
			}
			out = append(out, TimeWindow{StartTime: start, EndTime: next, Identifier: val})
			start = next
		}
	}
	return out
}
