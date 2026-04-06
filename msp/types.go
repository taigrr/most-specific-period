// Package msp implements the Most Specific Period algorithm, which selects
// the narrowest (shortest-duration) time period containing a given timestamp
// from a set of potentially overlapping periods.
package msp

import "time"

// Compile-time interface check.
var _ Period = TimeWindow{}

// Period represents a named time window with inclusive start and exclusive end.
type Period interface {
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetIdentifier() string
}
