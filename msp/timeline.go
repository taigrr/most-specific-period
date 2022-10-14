package msp

import "fmt"

// Outputs a formatted timeline of periods
func GenerateTimeline(periods ...Period) (out []string) {
	if len(periods) == 0 {
		out = []string{}
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
			frame := fmt.Sprintf("%s\t%s\t%s\n", val, start, next)
			out = append(out, frame)
			start = next
		}
	}
	return out
}
