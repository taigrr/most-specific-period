package msp

import (
	"fmt"
	"sort"
	"time"
)

func GetChangeOvers(periods ...Period) (changeovers []time.Time) {
	timeStamps := []time.Time{}
	for _, x := range periods {
		timeStamps = append(timeStamps, x.GetEndTime())
		timeStamps = append(timeStamps, x.GetStartTime())
	}
	if len(timeStamps) == 0 {
		return
	}
	sort.Slice(timeStamps, func(i, j int) bool {
		return timeStamps[i].Before(timeStamps[j])
	})
	// timeStamps is sorted, so this will always result in an unused time
	// struct, as it's before the first
	previousTs := timeStamps[0].Add(-10 * time.Nanosecond)
	for _, ts := range timeStamps {
		if ts.Equal(previousTs) {
			continue
		}
		previousTs = ts
		before := ts.Add(-1 * time.Nanosecond)
		after := ts.Add(1 * time.Nanosecond)
		from, _ := MostSpecificPeriod(before, periods...)
		to, _ := MostSpecificPeriod(after, periods...)
		if from == to {
			continue
		}
		changeovers = append(changeovers, ts)
	}
	return
}

func GetNextChangeOver(t time.Time, periods ...Period) (ts time.Time, err error) {
	changeOvers := GetChangeOvers(periods...)
	for _, ts := range changeOvers {
		if ts.After(t) {
			return ts, nil
		}
	}
	return time.Time{}, ErrNoNextChangeover
}

func FlattenPeriods(periods ...Period) (ids []string) {
	changeovers := GetChangeOvers(periods...)
	for _, c := range changeovers {
		id, err := MostSpecificPeriod(c, periods...)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return
}

func GenerateTimeline(periods ...Period) (out []string) {
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
			frame := fmt.Sprintf("%s\t%s\t%s\n", val, start, next)
			out = append(out, frame)
			start = next
		}
	}
	return out
}
