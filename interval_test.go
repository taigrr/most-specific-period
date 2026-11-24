package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommandRejectsInvalidIntervals(t *testing.T) {
	t.Parallel()
	for _, end := range []string{"2026-05-05T09:00:00Z", "2026-05-05T10:00:00Z", "2026-05-05T12:00:00+02:00"} {
		t.Run(end, func(t *testing.T) {
			input := "valid\n2026-05-05T00:00:00Z\n2026-05-06T00:00:00Z\ninvalid\n2026-05-05T10:00:00Z\n" + end + "\n"
			periods, err := parsePeriods(strings.NewReader(input), nil)
			if err == nil || !strings.Contains(err.Error(), `period "invalid": end time must be after start time`) {
				t.Fatalf("expected interval error, got %v", err)
			}
			if periods != nil {
				t.Fatalf("returned partial periods: %v", periods)
			}
			var output bytes.Buffer
			command := newRootCommand(strings.NewReader(input), &output, false)
			command.SetArgs([]string{"--date", "2026-05-05T11:00:00Z"})
			if err := command.Execute(); err == nil {
				t.Fatal("invalid interval accepted")
			}
			if strings.Contains(output.String(), "Timeline") {
				t.Fatalf("emitted partial timeline: %s", &output)
			}
		})
	}
}

func TestParsePeriodsAcceptsOrderedOffsetInterval(t *testing.T) {
	t.Parallel()
	// The end's wall-clock hour is earlier, but its instant is later.
	input := "offset\n2026-05-05T10:00:00+02:00\n2026-05-05T09:00:00Z\n"
	periods, err := parsePeriods(strings.NewReader(input), nil)
	if err != nil || len(periods) != 1 {
		t.Fatalf("valid offset interval rejected: periods=%v, err=%v", periods, err)
	}
}
