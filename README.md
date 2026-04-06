# Most Specific Period

[![Go Reference](https://pkg.go.dev/badge/github.com/taigrr/most-specific-period.svg)](https://pkg.go.dev/github.com/taigrr/most-specific-period)
[![Go Report Card](https://goreportcard.com/badge/github.com/taigrr/most-specific-period)](https://goreportcard.com/report/github.com/taigrr/most-specific-period)

A Go library for selecting the narrowest time period containing a given
timestamp from a set of potentially overlapping periods.

## Installation

```bash
go get github.com/taigrr/most-specific-period
```

## What is a Most Specific Period?

Given overlapping time periods, the MSP algorithm picks the most precise one:

- Given a single valid period containing the timestamp, that period is chosen.
- Given two overlapping periods of different lengths, the shorter one wins.
- Given two periods of equal length, the one that started more recently wins.
- Given two periods with the same duration and start time, the lexicographically
  last identifier wins (e.g. "B" over "A").
- Periods that haven't started yet or have already ended are ignored.

## Usage

Implement the `Period` interface or use the built-in `TimeWindow`:

```go
package main

import (
	"fmt"
	"time"

	"github.com/taigrr/most-specific-period/msp"
)

func main() {
	now := time.Now()

	periods := []msp.Period{
		msp.TimeWindow{
			StartTime:  now.Add(-24 * time.Hour),
			EndTime:    now.Add(24 * time.Hour),
			Identifier: "this-week",
		},
		msp.TimeWindow{
			StartTime:  now.Add(-1 * time.Hour),
			EndTime:    now.Add(1 * time.Hour),
			Identifier: "this-morning",
		},
	}

	id, err := msp.MostSpecificPeriod(now, periods...)
	if err != nil {
		fmt.Println("No matching period")
		return
	}
	fmt.Printf("MSP: %s\n", id) // "this-morning"
}
```

### Period Interface

```go
type Period interface {
	GetStartTime()  time.Time // Inclusive start time
	GetEndTime()    time.Time // Exclusive end time
	GetIdentifier() string
}
```

### Additional Functions

- `GenerateTimeline(periods...)` — Flatten overlapping periods into a
  non-overlapping timeline.
- `GetChangeOvers(periods...)` — Get timestamps where the MSP changes.
- `GetNextChangeOver(t, periods...)` — Get the next changeover after time `t`.
- `FlattenPeriods(periods...)` — Get ordered identifiers at each changeover.
- `ValidTimePeriods(ts, periods...)` — Filter periods valid at timestamp `ts`.
- `GetDuration(start, end)` — Calculate duration between two times.

## CLI

A demo CLI is included. It reads periods from stdin (one per three lines:
identifier, start time, end time in RFC 3339) and displays the timeline
and MSP.

```bash
go run . -d 2024-06-15T12:00:00Z <<EOF
summer
2024-06-01T00:00:00Z
2024-09-01T00:00:00Z
june
2024-06-01T00:00:00Z
2024-07-01T00:00:00Z
EOF
```

## License

0BSD — See [LICENSE](LICENSE) for details.
