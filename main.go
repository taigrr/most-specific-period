package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
	"github.com/taigrr/most-specific-period/msp"
)

var errIncompletePeriod = errors.New("incomplete period input")

type cliPeriod struct {
	EndTime    time.Time
	StartTime  time.Time
	Identifier string
}

func (p cliPeriod) GetEndTime() time.Time {
	return p.EndTime
}

func (p cliPeriod) GetStartTime() time.Time {
	return p.StartTime
}

func (p cliPeriod) GetIdentifier() string {
	return p.Identifier
}

func warnMessage(w io.Writer) {
	fmt.Fprint(w, "Please type your date formats as follows, hit return between each field (RFC 3339), and hit Control+D to signal you are complete: \nIdentifier: id\nStartTime: 2019-10-12T07:20:50.52Z\nEndTime: 2019-10-12T07:20:50.52Z\n")
}

func promptForField(w io.Writer, field int) {
	switch field {
	case 0:
		fmt.Fprint(w, "Identifier: ")
	case 1:
		fmt.Fprint(w, "StartTime: ")
	case 2:
		fmt.Fprint(w, "EndTime: ")
	}
}

func parsePeriods(r io.Reader, prompt func(int)) ([]msp.Period, error) {
	scanner := bufio.NewScanner(r)
	periods := []msp.Period{}
	currentPeriod := cliPeriod{}
	field := 0

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		switch field {
		case 0:
			currentPeriod = cliPeriod{Identifier: input}
		case 1:
			start, err := time.Parse(time.RFC3339, input)
			if err != nil {
				return nil, fmt.Errorf("invalid start timestamp %q: %w", input, err)
			}
			currentPeriod.StartTime = start
		case 2:
			end, err := time.Parse(time.RFC3339, input)
			if err != nil {
				return nil, fmt.Errorf("invalid end timestamp %q: %w", input, err)
			}
			currentPeriod.EndTime = end
			periods = append(periods, currentPeriod)
		}

		field = (field + 1) % 3
		if prompt != nil {
			prompt(field)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if field != 0 {
		return nil, errIncompletePeriod
	}

	return periods, nil
}

func newRootCommand(stdin io.Reader, stdout io.Writer, terminal bool) *cobra.Command {
	var start time.Time
	var showVersion bool
	var userDate string

	root := &cobra.Command{
		Use:           "most-specific-period",
		Short:         "Generate a timeline of periods and calculate the most specific period.",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				if Commit != "" {
					fmt.Fprintf(stdout, "most-specific-period %s (%s)\n", Version, Commit)
				} else {
					fmt.Fprintf(stdout, "most-specific-period %s\n", Version)
				}
				return nil
			}

			if userDate != "" {
				t, err := time.Parse(time.RFC3339, userDate)
				if err != nil {
					fmt.Fprintln(stdout, "Please enter the date using the YYYY-MM-DDT00:00:00.00Z")
					return err
				}
				start = t
			} else {
				start = time.Now()
			}

			if terminal {
				warnMessage(stdout)
				promptForField(stdout, 0)
			}

			var prompt func(int)
			if terminal {
				prompt = func(field int) {
					promptForField(stdout, field)
				}
			}

			periods, err := parsePeriods(stdin, prompt)
			if err != nil {
				if errors.Is(err, errIncompletePeriod) {
					fmt.Fprintln(stdout, "ERROR: each period must include identifier, start time, and end time")
				} else {
					fmt.Fprintf(stdout, "ERROR: %v\n", err)
				}
				return err
			}

			vals := msp.GenerateTimeline(periods...)
			fmt.Fprint(stdout, "\nTimeline of changeovers:\n")
			for _, val := range vals {
				fmt.Fprintln(stdout, val)
			}
			m, err := msp.MostSpecificPeriod(start, periods...)
			if err != nil {
				fmt.Fprintln(stdout, "No significant period found")
				return err
			}
			if terminal {
				fmt.Fprint(stdout, "\nThe MSP from the list was: ")
			}
			fmt.Fprintf(stdout, "%s\n", m)
			return nil
		},
	}
	root.Flags().BoolVarP(&showVersion, "version", "v", false, "print version and exit")
	root.Flags().StringVarP(&userDate, "date", "d", "", "use a custom RFC 3339 date to calculate MSP")
	return root
}

func main() {
	terminal := false
	fi, err := os.Stdin.Stat()
	if err == nil && (fi.Mode()&os.ModeCharDevice) != 0 {
		terminal = true
	}

	root := newRootCommand(os.Stdin, os.Stdout, terminal)
	if err := fang.Execute(context.Background(), root, fang.WithoutVersion()); err != nil {
		os.Exit(1)
	}
}
