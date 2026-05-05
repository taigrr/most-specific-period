package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/taigrr/most-specific-period/msp"
)

var errIncompletePeriod = errors.New("incomplete period input")

type Period struct {
	EndTime    time.Time
	StartTime  time.Time
	Identifier string
}

func (p Period) GetEndTime() time.Time {
	return p.EndTime
}

func (p Period) GetStartTime() time.Time {
	return p.StartTime
}

func (p Period) GetIdentifier() string {
	return p.Identifier
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func warnMessage() {
	fmt.Print("Please type your date formats as follows, hit return between each field (RFC 3339), and hit Control+D to signal you are complete: \nIdentifier: id\nStartTime: 2019-10-12T07:20:50.52Z\nEndTime: 2019-10-12T07:20:50.52Z\n")
}

func helpMessage() {
	fmt.Print("\nmost-specific-period [-h][-d]\n\nGenerates a timeline of periods and will provide a most specific period if available.\n\n-h\tShows this help menu\n-d\tProvide an RFC 3339 time to provide an alternate point for calculating MSP.")
}

func promptForField(field int) {
	switch field {
	case 0:
		fmt.Print("Identifier: ")
	case 1:
		fmt.Print("StartTime: ")
	case 2:
		fmt.Print("EndTime: ")
	}
}

func parsePeriods(r io.Reader, prompt func(int)) ([]msp.Period, error) {
	scanner := bufio.NewScanner(r)
	periods := []msp.Period{}
	currentPeriod := Period{}
	field := 0

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		switch field {
		case 0:
			currentPeriod = Period{Identifier: input}
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

func main() {
	var start time.Time
	help := flag.Bool("h", false, "displays help command")
	userDate := flag.String("d", "", "use a custom date to calculate MSP")
	flag.Parse()
	if *help {
		helpMessage()
		os.Exit(0)
	}

	if userDate != nil && *userDate != "" {
		t, err := time.Parse(time.RFC3339, *userDate)
		if err != nil {
			fmt.Println("Please enter the date using the YYYY-MM-DDT00:00:00.00Z")
			os.Exit(1)
		}
		start = t
	} else {
		start = time.Now()
	}
	terminal := false
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		// this is a file being read in, no need to print the prompt just yet
	} else {
		// this is a terminal, let's help the user out
		terminal = true
		warnMessage()
		promptForField(0)
	}

	var prompt func(int)
	if terminal {
		prompt = promptForField
	}

	periods, err := parsePeriods(os.Stdin, prompt)
	if err != nil {
		if errors.Is(err, errIncompletePeriod) {
			fmt.Println("ERROR: each period must include identifier, start time, and end time")
		} else {
			fmt.Printf("ERROR: %v\n", err)
		}
		os.Exit(1)
	}

	vals := msp.GenerateTimeline(periods...)
	fmt.Print("\nTimeline of changeovers:\n")
	for _, val := range vals {
		fmt.Println(val)
	}
	m, err := msp.MostSpecificPeriod(start, periods...)
	if err != nil {
		fmt.Printf("No significant period found\n")
		os.Exit(1)
	}
	if terminal {
		fmt.Printf("\nThe MSP from the list was: ")
	}
	fmt.Printf("%s\n", m)
}
