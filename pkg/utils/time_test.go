package utils

import (
	"fmt"
	"testing"
)

func TestProcessTimeRange(t *testing.T) {
	cases := []string{
		"09:15:00-09:30:00,13:00:00-14:57:00",
		"09:15:00-09:30:00",
		"09:15:00-09:30:00,13:00:00-14:57:00,15:00:00-15:30:00",
	}

	for _, c := range cases {
		ranges := ProcessTimeRange(c)
		if len(ranges) == 0 {
			t.Errorf("ProcessTimeRange(%s) = %v, want non-empty", c, ranges)
		}
		t.Logf("ProcessTimeRange(%s) = %v", c, ranges)
	}
}

func TestParseIntervalString(t *testing.T) {
	s := "[0-500):12:20"
	a, b, c, d, err := ParseIntervalString(s)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(a, b, c, d)
}
