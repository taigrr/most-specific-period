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

func warnMessage(w io.Writer) error {
	_, err := fmt.Fprint(w, "Please type your date formats as follows, hit return between each field (RFC 3339), and hit Control+D to signal you are complete: \nIdentifier: id\nStartTime: 2019-10-12T07:20:50.52Z\nEndTime: 2019-10-12T07:20:50.52Z\n")
	return err
}

func promptForField(w io.Writer, field int) error {
	switch field {
	case 0:
		_, err := fmt.Fprint(w, "Identifier: ")
		return err
	case 1:
		_, err := fmt.Fprint(w, "StartTime: ")
		return err
	case 2:
		_, err := fmt.Fprint(w, "EndTime: ")
		return err
	}
	return nil
}

func parsePeriods(r io.Reader, prompt func(int) error) ([]msp.Period, error) {
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
			if err := prompt(field); err != nil {
				return nil, err
			}
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
					if _, writeErr := fmt.Fprintf(stdout, "most-specific-period %s (%s)\n", Version, Commit); writeErr != nil {
						return writeErr
					}
				} else {
					if _, writeErr := fmt.Fprintf(stdout, "most-specific-period %s\n", Version); writeErr != nil {
						return writeErr
					}
				}
				return nil
			}

			if userDate != "" {
				t, err := time.Parse(time.RFC3339, userDate)
				if err != nil {
					if _, writeErr := fmt.Fprintln(stdout, "Please enter the date using the YYYY-MM-DDT00:00:00.00Z"); writeErr != nil {
						return writeErr
					}
					return err
				}
				start = t
			} else {
				start = time.Now()
			}

			if terminal {
				if err := warnMessage(stdout); err != nil {
					return err
				}
				if err := promptForField(stdout, 0); err != nil {
					return err
				}
			}

			var prompt func(int) error
			var promptErr error
			if terminal {
				prompt = func(field int) error {
					promptErr = promptForField(stdout, field)
					return promptErr
				}
			}

			periods, err := parsePeriods(stdin, prompt)
			if promptErr != nil {
				return promptErr
			}
			if err != nil {
				if errors.Is(err, errIncompletePeriod) {
					if _, writeErr := fmt.Fprintln(stdout, "ERROR: each period must include identifier, start time, and end time"); writeErr != nil {
						return writeErr
					}
				} else {
					if _, writeErr := fmt.Fprintf(stdout, "ERROR: %v\n", err); writeErr != nil {
						return writeErr
					}
				}
				return err
			}

			vals := msp.GenerateTimeline(periods...)
			if _, writeErr := fmt.Fprint(stdout, "\nTimeline of changeovers:\n"); writeErr != nil {
				return writeErr
			}
			for _, val := range vals {
				if _, writeErr := fmt.Fprintln(stdout, val); writeErr != nil {
					return writeErr
				}
			}
			m, err := msp.MostSpecificPeriod(start, periods...)
			if err != nil {
				if _, writeErr := fmt.Fprintln(stdout, "No significant period found"); writeErr != nil {
					return writeErr
				}
				return err
			}
			if terminal {
				if _, writeErr := fmt.Fprint(stdout, "\nThe MSP from the list was: "); writeErr != nil {
					return writeErr
				}
			}
			if _, writeErr := fmt.Fprintf(stdout, "%s\n", m); writeErr != nil {
				return writeErr
			}
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
