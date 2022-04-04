package msp

import (
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
	return time.Unix(0, 0), ErrNoNextChangeover
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
