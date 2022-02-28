package msp

import "time"

type Period struct {
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Identifier string    `json:"identifier"`
}
