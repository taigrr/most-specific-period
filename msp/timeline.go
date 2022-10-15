package msp

import (
	"fmt"
	"time"
)

type TimeWindow struct {
	StartTime  time.Time
	EndTime    time.Time
	Identifier string
}

func (p TimeWindow) GetIdentifier() string {
	return p.Identifier
}

func (p TimeWindow) GetEndTime() time.Time {
	return p.EndTime
}

func (p TimeWindow) GetStartTime() time.Time {
	return p.StartTime
}

func (t TimeWindow) String() string {
	return fmt.Sprintf("%s\t%s\t%s",
		t.GetIdentifier(),
		t.GetStartTime(),
		t.GetEndTime())
}

// Outputs a formatted timeline of periods
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
