// This file contains general types.

package repcore

import (
	"fmt"
	"time"
)

// Frame is the basic time unit in StarCraft.
// There are approximately ~23.81 frames in a second;
// 1 frame = 0.042 second = 42 ms to be exact.
type Frame int32

// Seconds returns the time equivalent to the frames in seconds.
func (f Frame) Seconds() float64 {
	return float64(f.Milliseconds()) / 1000
}

// Milliseconds returns the time equivalent to the frames in milliseconds.
func (f Frame) Milliseconds() int64 {
	return int64(f) * 42
}

// Duration returns the frame as a time.Duration value.
func (f Frame) Duration() time.Duration {
	return time.Millisecond * time.Duration(f.Milliseconds())
}

// Point describes a point in the map.
type Point struct {
	// X and Y coordinates of the point
	// 1 Tile is 32 units (pixel)
	X, Y uint16
}

// String returns a string representation of the point in the format:
//     "x=X, y=Y"
func (p Point) String() string {
	return fmt.Sprint("x=", p.X, ", y=", p.Y)
}
