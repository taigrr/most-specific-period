package msp

import (
	"fmt"
	"testing"
)

//(periods ...Period) (id string, err error) {
func TestMostSpecificPeriod(t *testing.T) {
	testCases := []struct {
		testID  string
		result  string
		Periods []Period
	}{{testID: "test"}}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.testID), func(t *testing.T) {

		})
	}
}
