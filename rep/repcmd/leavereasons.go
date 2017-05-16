// This file contains hotkey types.

package repcmd

import "github.com/icza/screp/rep/repcore"

// LeaveReason describes the leave reason.
type LeaveReason struct {
	repcore.Enum

	// ID as it appears in replays
	ID byte
}

// LeaveReasons is an enumeration of the possible leave reasons.
var LeaveReasons = []*LeaveReason{
	{e("Quit"), 0x00},
	{e("Defeat"), 0x01},
	{e("Victory"), 0x02},
	{e("Finished"), 0x03},
	{e("Draw"), 0x04},
	{e("Dropped"), 0x05},
}

// LeaveReasonByID returns the LeaveReason for a given ID.
// A new LeaveReason with Unknown name is returned if one is not found
// for the given ID (preserving the unknown ID).
func LeaveReasonByID(ID byte) *LeaveReason {
	if int(ID) < len(LeaveReasons) {
		return LeaveReasons[ID]
	}
	return &LeaveReason{repcore.UnknownEnum(ID), ID}
}
