package msp

import (
	"sort"
	"time"
)

// GetChangeOvers returns the sorted list of timestamps where the most
// specific period changes from one identifier to another.
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
	var previousTs time.Time
	seenPrevious := false
	for _, ts := range timeStamps {
		if seenPrevious && ts.Equal(previousTs) {
			continue
		}
		previousTs = ts
		seenPrevious = true
		from, to := adjacentPeriodIDs(ts, periods...)
		if from == to {
			continue
		}
		changeovers = append(changeovers, ts)
	}
	return
}

func adjacentPeriodIDs(ts time.Time, periods ...Period) (from string, to string) {
	if before, ok := addTime(ts, -1*time.Nanosecond); ok {
		from, _ = MostSpecificPeriod(before, periods...)
	}
	if after, ok := addTime(ts, time.Nanosecond); ok {
		to, _ = MostSpecificPeriod(after, periods...)
	}
	return from, to
}

func addTime(ts time.Time, duration time.Duration) (out time.Time, ok bool) {
	defer func() {
		if recover() != nil {
			out = time.Time{}
			ok = false
		}
	}()
	return ts.Add(duration), true
}

// GetNextChangeOver returns the first changeover timestamp strictly after t.
// If no such changeover exists, ErrNoNextChangeover is returned.
func GetNextChangeOver(t time.Time, periods ...Period) (ts time.Time, err error) {
	changeOvers := GetChangeOvers(periods...)
	for _, ts := range changeOvers {
		if ts.After(t) {
			return ts, nil
		}
	}
	return time.Time{}, ErrNoNextChangeover
}

// FlattenPeriods returns an ordered list of period identifiers representing
// the most specific period at each changeover point.
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
