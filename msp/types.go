package msp

import "time"

// Compile-time interface check.
var _ Period = TimeWindow{}

type Period interface {
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetIdentifier() string
}
