package msp

import "time"

type Period interface {
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetIdentifier() string
}
