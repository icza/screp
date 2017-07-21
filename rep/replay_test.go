package rep

import (
	"math"
	"testing"
)

func TestAngleToClock(t *testing.T) {
	cases := []struct {
		angle float64
		clock int
	}{
		{0, 3},
		{math.Pi / 2, 12},
		{math.Pi, 9},
		{math.Pi * 3 / 2, 6},

		{0 + math.Pi*6, 3},
		{0 - math.Pi*6, 3},

		{math.Pi/2 + math.Pi/13, 12},
		{math.Pi/2 - math.Pi/13, 12},
	}

	for _, c := range cases {
		if got := angleToClock(c.angle); got != c.clock {
			t.Errorf("Expected: %v, got: %v", c.clock, got)
		}
	}
}
