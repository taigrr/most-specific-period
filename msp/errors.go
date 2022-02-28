package msp

import (
	"errors"
)

var (
	// ErrEndAfterStart occurs when a period given has an end time after its start time
	ErrEndAfterStart = errors.New("error: start time is after end time")
	// ErrNoValidPeriods occurs when an empty set of periods is passed or when ll periods are invalid
	ErrNoValidPeriods = errors.New("error: no valid periods available")
)
