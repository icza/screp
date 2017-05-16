// This file contains latencies.

package repcmd

import "github.com/icza/screp/rep/repcore"

// Latency describes the latency.
type Latency struct {
	repcore.Enum

	// ID as it appears in replays
	ID byte
}

// Latencies is an enumeration of the possible latencies.
var Latencies = []*Latency{
	{e("Low"), 0x00},
	{e("High"), 0x01},
	{e("Extra High"), 0x02},
}

// LatencyTypeByID returns the Latency for a given ID.
// A new Latency with Unknown name is returned if one is not found
// for the given ID (preserving the unknown ID).
func LatencyTypeByID(ID byte) *Latency {
	if int(ID) < len(Latencies) {
		return Latencies[ID]
	}
	return &Latency{repcore.UnknownEnum(ID), ID}
}
