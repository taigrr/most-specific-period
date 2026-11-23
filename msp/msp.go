package msp

import (
	"sort"
	"time"
)

// MostSpecificPeriod returns the identifier of the shortest-duration period
// that contains timestamp ts. When multiple periods share the shortest
// duration, the one with the latest start time wins; if start times also
// match, the lexicographically last identifier is returned.
func MostSpecificPeriod(ts time.Time, periods ...Period) (id string, err error) {
	// Filter to get only valid periods here
	periods = ValidTimePeriods(ts, periods...)
	if len(periods) == 0 {
		return "", ErrNoValidPeriods
	}
	// find the shortest duration
	d := exactDuration(periods[0])
	for _, x := range periods {
		p := exactDuration(x)
		if p.seconds < d.seconds || (p.seconds == d.seconds && p.nanoseconds < d.nanoseconds) {
			d = p
		}
	}
	// find all periods with this shortest duration
	var matchingDurations []Period
	for _, x := range periods {
		p := exactDuration(x)
		if p == d {
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
	// Determine whichever of these periods have the same start time in addition to duration
	var matchingDurationsAndStartTimes []Period
	for _, x := range matchingDurations {
		if x.GetStartTime().Equal(newest) {
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

// periodDuration represents a nonnegative span without time.Duration's
// approximately 292-year saturation limit.
type periodDuration struct {
	seconds     uint64
	nanoseconds int
}

// exactDuration is called only for periods containing the query timestamp,
// so end is strictly after start. Unsigned subtraction also handles spans
// crossing the Unix epoch without signed integer overflow.
func exactDuration(period Period) periodDuration {
	start, end := period.GetStartTime(), period.GetEndTime()
	seconds := uint64(end.Unix()) - uint64(start.Unix())
	nanoseconds := end.Nanosecond() - start.Nanosecond()
	if nanoseconds < 0 {
		seconds--
		nanoseconds += int(time.Second)
	}
	return periodDuration{seconds: seconds, nanoseconds: nanoseconds}
}

// GetDuration returns the duration between start and end. If start is after
// end, ErrEndAfterStart is returned alongside the (negative) duration.
func GetDuration(start time.Time, end time.Time) (dur time.Duration, err error) {
	if start.After(end) {
		err = ErrEndAfterStart
	}
	dur = end.Sub(start)
	return dur, err
}

// ValidTimePeriods filters periods to those whose start time is at or before
// ts and whose end time is strictly after ts.
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
