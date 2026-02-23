package msp

import (
	"errors"
)

var (
	// ErrEndAfterStart occurs when a period's start time is after its end time
	ErrEndAfterStart = errors.New("error: start time is after end time")
	// ErrNoValidPeriods occurs when an empty set of periods is passed or when all periods are invalid
	ErrNoValidPeriods = errors.New("error: no valid periods available")
	// ErrNoNextChangeover occurs when GetNextChangeover is called but there are no changeovers after t
	ErrNoNextChangeover = errors.New("error: no valid changeovers available")
)
