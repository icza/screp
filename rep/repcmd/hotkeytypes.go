// This file contains hotkey types.

package repcmd

import "github.com/icza/screp/rep/repcore"

// HotkeyType describes the hotkey type.
type HotkeyType struct {
	repcore.Enum

	// ID as it appears in replays
	ID byte
}

// HotkeyTypes is an enumeration of the possible hotkey types.
var HotkeyTypes = []*HotkeyType{
	{e("Assign"), 0x00},
	{e("Select"), 0x01},
	{e("Add"), 0x02},
}

// HotkeyTypeByID returns the HotkeyType for a given ID.
// A new HotkeyType with Unknown name is returned if one is not found
// for the given ID (preserving the unknown ID).
func HotkeyTypeByID(ID byte) *HotkeyType {
	if int(ID) < len(HotkeyTypes) {
		return HotkeyTypes[ID]
	}
	return &HotkeyType{repcore.UnknownEnum(ID), ID}
}
