package msp

import (
	"sort"
	"time"
)

func MostSpecificPeriod(ts time.Time, periods ...Period) (id string, err error) {
	// Filter to get only valid periods here
	periods = ValidTimePeriods(ts, periods...)
	if len(periods) == 0 {
		return "", ErrNoValidPeriods
	}
	// find the shortest duration
	d, err := GetDuration(periods[0].GetStartTime(), periods[0].GetEndTime())
	for _, x := range periods {
		p, err := GetDuration(x.GetStartTime(), x.GetEndTime())
		if err == nil && p < d {
			d = p
		}
	}
	// find all periods with this shortest duration
	var matchingDurations []Period
	for _, x := range periods {
		p, err := GetDuration(x.GetStartTime(), x.GetEndTime())
		if err == nil && p == d {
			matchingDurations = append(matchingDurations, x)
		}
	}
	// Find the newest time a period starts
	newest := matchingDurations[0].GetStartTime()
	for _, x := range matchingDurations {
		if x.GetStartTime().After(newest) {
			newest = x.GetStartTime()
		}
	}
	// Determine whichever of these periods have the same start time in addtion to duration
	var matchingDurationsAndStartTimes []Period
	for _, x := range matchingDurations {
		if x.GetStartTime() == newest {
			matchingDurationsAndStartTimes = append(matchingDurationsAndStartTimes, x)
		}
	}
	// Finally, return the period with the 'last' name lexicographically
	var identifiers []string
	for _, x := range matchingDurationsAndStartTimes {
		identifiers = append(identifiers, x.GetIdentifier())
	}
	sort.Strings(identifiers)
	return identifiers[len(identifiers)-1], nil
}

func GetDuration(start time.Time, end time.Time) (dur time.Duration, err error) {
	if start.After(end) {
		err = ErrEndAfterStart
	}
	dur = end.Sub(start)
	return dur, err
}

func ValidTimePeriods(ts time.Time, periods ...Period) []Period {
	var valid []Period
	for _, p := range periods {
		start := p.GetStartTime()
		end := p.GetEndTime()
		if (start.Before(ts) || start.Equal(ts)) && (end.After(ts)) {
			valid = append(valid, p)
		}
	}
	return valid
}
