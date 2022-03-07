package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/taigrr/most-specific-period/msp"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}
func warnMessage() {
	fmt.Print("Please type your date formats as follows, hit return between each field (RFC 3339), and hit Control+D to signal you are complete: \nIdentifier: id\nStartTime: 2019-10-12T07:20:50.52Z\nEndTime: 2019-10-12T07:20:50.52Z\n")

}

func main() {
	flag.Parse()
	terminal := false
	fi, _ := os.Stdin.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		// this is a file being read in, no need to print the prompt just yet
	} else {
		// this is a terminal, let's help the user out
		terminal = true
		warnMessage()
	}
	s := bufio.NewScanner(os.Stdin)
	count := 1

	if terminal {
		fmt.Print("Identifier: ")
	}

	periods := []msp.Period{}

	for s.Scan() {
		input := s.Text()
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		if count%3 == 0 {
			t, err := time.Parse(time.RFC3339, input)
			if err != nil {
				fmt.Printf("ERROR: Invalid timestamp: %v", t)
				os.Exit(1)
			}
			periods[(count-1)/3].EndTime = t

			if terminal {
				fmt.Print("Identifier: ")
			}
		}
		if count%3 == 1 {
			periods = append(periods, msp.Period{Identifier: s.Text()})
			if terminal {
				fmt.Print("StartTime: ")
			}
		}
		if count%3 == 2 {
			t, err := time.Parse(time.RFC3339, input)
			if err != nil {
				fmt.Printf("ERROR: Invalid timestamp: %v", t)
				os.Exit(1)
			}
			periods[(count-1)/3].StartTime = t

			if terminal {
				fmt.Print("EndTime: ")
			}
		}
		count++
	}

	m, err := msp.MostSpecificPeriod(time.Now(), periods...)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	if terminal {
		fmt.Printf("The MSP from the list was: ")
	}
	fmt.Printf("%s\n", m)
}
