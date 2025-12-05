// This file contains hotkey types.

package repcmd

import "github.com/icza/screp/rep/repcore"

// CheatCode describes a cheat code.
type CheatCode struct {
	repcore.Enum

	// Bitmask of the enabled cheat codes as it appears in replays.
	BitMask uint32
}

// CheatCodes is an enumeration of the possible cheat codes.
var CheatCodes = []*CheatCode{
	{e("Black Sheep Wall"), 0x01},              // Reveals entire map [toggleable]
	{e("Operation CWAL"), 0x02},                // Build Faster (effects enemy teams) [toggleable]
	{e("Power Overwhelming"), 0x04},            // Grants Invincibility [toggleable]
	{e("Something for Nothing"), 0x08},         // Grants Free Upgrades
	{e("Show Me The Money"), 0x10},             // Gain 10,000 Minerals and Vespene Gas [repeatable]
	{e("Game Over Man"), 0x40},                 // Instantly Lose the Mission
	{e("There is no cow level"), 0x80},         // Instantly Win Mission
	{e("Staying Alive"), 0x100},                // Can Continue Playing After Winning [toggleable]
	{e("Ophelia"), 0x200},                      // Enable Mission Select (enter terran#, zerg#, or protoss# after this cheat)
	{e("The Gathering"), 0x800},                // Unit Spells and Abilities Are Free [toggleable]
	{e("Medieval Man"), 0x1000},                // Unlocks all Research Abilities [toggleable]
	{e("Modify The Phase Variance"), 0x2000},   // Unlocks all Buildings [toggleable]
	{e("War Aint What it Used to be"), 0x4000}, // Disabled Fog of War [toggleable]
	{e("Food For Thought"), 0x20000},           // Eliminates Unit Supply Limit [toggleable]
	{e("Whats Mine is Mine"), 0x40000},         // Gain 500 Crystals [repeatable]
	{e("Breathe Deep"), 0x80000},               // Gain 500 Vespene Gas [repeatable]
	{e("Zoom Zoom"), 0x100000},                 // Allows for further zoom out and zoom in (added in StarCraft: Remastered, usable in both versions)
	{e("Noglues"), 0x20000000},                 // Disables chapter and exposition screens in between campaign levels

	// "[race][#]" jump to given mission, not recorded
	// "radio free zerg" plays secret zerg theme song (only works in Brood War while playing as zerg), not recorded
}

var bitMaskCheatCodes = map[uint32]*CheatCode{}

func init() {
	for _, cc := range CheatCodes {
		bitMaskCheatCodes[cc.BitMask] = cc
	}
}

// CheatCodesByBitMap returns the CheatCodes identified by the given bitmap.
// A new CheatCode with Unknown name is added to the slice if one is not found
// for a given bit (preserving the unknown bitmask).
func CheatCodesByBitMap(bitmap uint32) (ccs []*CheatCode) {
	for bitMask := uint32(0x01); bitmap != 0; bitMask <<= 1 {
		if bitmap&bitMask == 0 {
			continue
		}

		if cc := bitMaskCheatCodes[bitMask]; cc != nil {
			ccs = append(ccs, cc)
		} else {
			ccs = append(ccs, &CheatCode{repcore.UnknownEnum(bitMask), bitMask})
		}

		// Clear the processed bit so we can stop when no more bits left:
		bitmap &^= bitMask
		// Note: bitmap and bitMask are both uint32, so eventually bitmap will be 0.
	}

	return
}
